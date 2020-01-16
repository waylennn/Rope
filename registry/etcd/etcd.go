package etcd

import (
	"context"
	"cope/registry"
	"encoding/json"
	"fmt"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"go.etcd.io/etcd/clientv3"
)

const (
	MaxServiceNum          = 8
	MaxSyncServiceInterval = time.Second * 10
)

//EtcdRegistry etcd register plugin
type EtcdRegistry struct {
	options   *registry.Options
	client    *clientv3.Client
	serviceCh chan *registry.Service

	registryServiceMap map[string]*RegisterService
	lock               sync.Mutex
	value              atomic.Value
}

var (
	etcdRegistry *EtcdRegistry = &EtcdRegistry{
		serviceCh:          make(chan *registry.Service, MaxServiceNum),
		registryServiceMap: make(map[string]*RegisterService, MaxServiceNum),
	}
)

// AllServiceInfo 存入value中 使用原子操作 避免加锁降低性能
type AllServiceInfo struct {
	serviceMap map[string]*registry.Service
}

//RegisterService ...
type RegisterService struct {
	id          clientv3.LeaseID
	service     *registry.Service
	registered  bool
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse
}

func init() {
	registry.RegisterPlugin(etcdRegistry)
	allServiceInfo := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	etcdRegistry.value.Store(allServiceInfo)

	go etcdRegistry.run()
}

//Name etcd name
func (e *EtcdRegistry) Name() string {
	return "etcd"
}

//Init etcd init
func (e *EtcdRegistry) Init(ctx context.Context, opts ...registry.Option) (err error) {

	e.options = &registry.Options{}
	for _, opt := range opts {
		opt(e.options)
	}

	e.client, err = clientv3.New(clientv3.Config{
		Endpoints:   e.options.Addrs,
		DialTimeout: e.options.TimeOut,
	})
	if err != nil {
		fmt.Println(err)
		err = fmt.Errorf("init etcd failed, err:%v", err)
		return
	}

	return
}

// Register service register
func (e *EtcdRegistry) Register(ctx context.Context, service *registry.Service) (err error) {

	select {
	case e.serviceCh <- service:
	default:
		err = fmt.Errorf("register chan is full")
		return
	}

	return
}

// Unregister service unregister
func (e *EtcdRegistry) Unregister(ctx context.Context, service *registry.Service) (err error) {
	return
}

func (e *EtcdRegistry) run() {
	ticker := time.NewTicker(MaxSyncServiceInterval)

	//把服务从管道里面取出来 放进map里面,然后进行注册或者保持心跳
	for {
		select {
		case service := <-e.serviceCh:
			serviceOld, ok := e.registryServiceMap[service.Name]
			if ok {
				for _, node := range service.Nodes {
					serviceOld.service.Nodes = append(serviceOld.service.Nodes, node)
				}
				// serviceOld.registered = false
				e.putInfoToEtcd(serviceOld, serviceOld.id)

				break
			}

			registryService := &RegisterService{
				service: service,
			}
			e.registryServiceMap[service.Name] = registryService
		case <-ticker.C:
			e.syncServiceFromEtcd()
		default:

			e.registerOrKeepAlive()
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (e *EtcdRegistry) registerOrKeepAlive() {

	for _, registryService := range e.registryServiceMap {

		if registryService.registered {
			e.keepAlive(registryService)
			continue
		}
		e.registerService(registryService)

	}

}

func (e *EtcdRegistry) keepAlive(registryService *RegisterService) {

	select {
	case resp := <-registryService.keepAliveCh:
		// fmt.Println(resp)
		if resp == nil {
			registryService.registered = false
			return
		}
	}

	return
}

//把服务信息put进etcd封装起来的一个函数
func (e *EtcdRegistry) putInfoToEtcd(registryService *RegisterService, respID clientv3.LeaseID) {
	for _, node := range registryService.service.Nodes {

		tmp := &registry.Service{
			Name:  registryService.service.Name,
			Nodes: []*registry.Nodes{node},
		}

		data, err := json.Marshal(tmp)
		if err != nil {
			continue
		}

		key := e.serviceNodePath(tmp)
		_, err = e.client.Put(context.TODO(), key, string(data), clientv3.WithLease(respID))
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	return
}

func (e *EtcdRegistry) registerService(registryService *RegisterService) {
	resp, err := e.client.Grant(context.TODO(), e.options.HeartBeat)
	if err != nil {
		fmt.Println(err)
		return
	}

	registryService.id = resp.ID
	//把信息注册到etcd里面
	e.putInfoToEtcd(registryService, registryService.id)

	// the key 'foo' will be kept forever
	ch, err := e.client.KeepAlive(context.TODO(), resp.ID)
	if err != nil {
		fmt.Println(err)
	}
	registryService.keepAliveCh = ch
	registryService.registered = true

	return
}

func (e *EtcdRegistry) serviceNodePath(service *registry.Service) string {

	nodeIP := fmt.Sprintf("%s:%d", service.Nodes[0].IP, service.Nodes[0].Port)
	return path.Join(e.options.RegistryPath, service.Name, nodeIP)
}

func (e *EtcdRegistry) servicePath(name string) string {
	//拼出服务的路径,查找出来的就是这个服务下的所有节点
	return path.Join(e.options.RegistryPath, name)
}

func (e *EtcdRegistry) getServiceFromCache(ctx context.Context, name string) (service *registry.Service, ok bool) {

	allServiceInfo := e.value.Load().(*AllServiceInfo)
	//一般情况下，都会从缓存中读取
	service, ok = allServiceInfo.serviceMap[name]
	return
}

//GetService ...
func (e *EtcdRegistry) GetService(ctx context.Context, name string) (service *registry.Service, err error) {

	// 进来的请求先从缓存里面取读取信息

	service, ok := e.getServiceFromCache(ctx, name)
	if ok {
		return
	}

	//没有则从etcd请求数据，防止etcd崩溃，加锁
	e.lock.Lock()
	defer e.lock.Unlock()
	// 再次检测是否加载成功，这个时候可能其他拿到锁的已经更新了缓存
	service, ok = e.getServiceFromCache(ctx, name)
	if ok {
		return
	}

	key := e.servicePath(name)
	res, err := e.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		fmt.Print(err)
		return
	}
	//构造一个需要返回的对象
	service = &registry.Service{
		Name: name,
	}
	if res.Kvs == nil {
		return
	}
	for _, kvs := range res.Kvs {

		var tmpService registry.Service

		err = json.Unmarshal(kvs.Value, &tmpService)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, node := range tmpService.Nodes {
			service.Nodes = append(service.Nodes, node)
		}
	}
	//更新缓存
	allServiceInfoOld := e.value.Load().(*AllServiceInfo)

	allServiceInfoNew := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}

	for key, val := range allServiceInfoOld.serviceMap {
		allServiceInfoNew.serviceMap[key] = val
	}

	allServiceInfoNew.serviceMap[name] = service
	e.value.Store(allServiceInfoNew)

	return
}

//定时器,同步etcd服务信息
func (e *EtcdRegistry) syncServiceFromEtcd() {

	allServiceInfoOld := e.value.Load().(*AllServiceInfo)

	allServiceInfoNew := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	//更新缓存中的节点，若缓存中的节点不存在 则不管它
	for name, service := range allServiceInfoOld.serviceMap {

		key := e.servicePath(name)
		res, err := e.client.Get(context.TODO(), key, clientv3.WithPrefix())
		if err != nil {
			fmt.Print(err)
			//若etcd不存在该节点 则用原来的节点
			allServiceInfoNew.serviceMap[name] = service
			continue
		}

		serviceNew := &registry.Service{
			Name: service.Name,
		}

		for _, kv := range res.Kvs {
			var tmpService registry.Service

			err = json.Unmarshal(kv.Value, &tmpService)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, node := range tmpService.Nodes {
				serviceNew.Nodes = append(serviceNew.Nodes, node)
			}
		}

		allServiceInfoNew.serviceMap[name] = serviceNew

	}

	e.value.Store(allServiceInfoNew)
}
