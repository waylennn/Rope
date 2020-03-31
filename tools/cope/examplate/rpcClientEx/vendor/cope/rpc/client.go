package rpc

import (
	"context"
	"cope/loadbalance"
	"cope/logs"
	"cope/meta"
	"cope/middleware"
	"cope/registry"
	_ "cope/registry/etcd"
	"cope/registry/loadbanlance"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

var initRegistryOnce sync.Once

//CopeClient ...
type CopeClient struct {
	limiter  *rate.Limiter
	opts     *RpcOptions
	register registry.Registry
	balance  loadbalance.LoadBalance
}

//CopeConf 临时日志对象
type CopeConf struct {
	*Logs
}

//Logs ...
type Logs struct {
	ConsoloLog  bool
	Path        string
	ServiceName string
	Level       string
	ChanSize    int
}

func init() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		addr := fmt.Sprintf("0.0.0.0:%d", 8888)
		http.ListenAndServe(addr, nil)
	}()
}

//NewCopeClient ...
func NewCopeClient(serviceName string, copeConf *CopeConf, optfunc ...RpcOptionFunc) *CopeClient {
	ctx := context.Background()

	err := initLogger(copeConf)
	if err != nil {
		panic(err)
	}

	copeClient := &CopeClient{
		opts: &RpcOptions{
			ServiceName:       serviceName,
			RegisterName:      "etcd",
			RegisterAddr:      "127.0.0.1:2379",
			RegisterPath:      "/ibinarytree/koala/service/",
			TraceReportAddr:   "http://60.205.218.189:9411/api/v1/spans",
			TraceSampleType:   "const",
			TraceSampleRate:   1,
			ClientServiceName: "default",
		},
		balance: loadbanlance.NewRandomBalance(),
	}
	for _, opt := range optfunc {
		opt(copeClient.opts)
	}
	//连接etcd
	initRegistryOnce.Do(func() {
		registry, err := registry.InitRegistry(
			context.Background(),
			copeClient.opts.RegisterName,
			registry.WithAddrs([]string{copeClient.opts.RegisterAddr}),
			registry.WithTimeout(time.Second),
			registry.WithRegistryPath(copeClient.opts.RegisterPath),
			registry.WithHeartBeat(10),
		)
		if err != nil {
			logs.Error(ctx, "init registry failed, err:%v", err)
			return
		}
		copeClient.register = registry
	})

	if copeClient.opts.MaxLimitQps > 0 {
		copeClient.limiter = rate.NewLimiter(rate.Limit(copeClient.opts.MaxLimitQps),
			copeClient.opts.MaxLimitQps)
	}

	return copeClient
}

//BuildClientMiddleware 把中间件函数串起来
func (c *CopeClient) BuildClientMiddleware(handleFunc middleware.MiddlewareFunc) middleware.MiddlewareFunc {
	var mids []middleware.Middleware

	// mids = append(mids, middleware.RpcLogMiddleware)
	mids = append(mids, middleware.PrometheusClientMiddleware)
	if c.limiter != nil {
		mids = append(mids, middleware.RateLimitMiddleware(c.limiter))
	}
	mids = append(mids, middleware.HystrixMiddleware)
	mids = append(mids, middleware.NewDiscoveryMiddleware(c.register))
	mids = append(mids, middleware.NewLoadBalanceMiddleware(c.balance))
	mids = append(mids, middleware.ClientShortConnMiddleware)

	if len(mids) > 0 {
		middleChain := middleware.Chain(mids[0], mids[1:]...)
		return middleChain(handleFunc)
	}
	return handleFunc
}

//Call 对外提供的调用接口
func (c *CopeClient) Call(ctx context.Context, method string, req interface{}, handleFunc middleware.MiddlewareFunc) (res interface{}, err error) {
	caller := c.getCaller(ctx)
	ctx = meta.InitRpcMeta(ctx, c.opts.ServiceName, method, caller)
	middlewareFunc := c.BuildClientMiddleware(handleFunc)
	res, err = middlewareFunc(ctx, req)
	if err != nil {
		return
	}
	return
}

func (c *CopeClient) getCaller(ctx context.Context) string {
	serverMeta := meta.GetServerMeta(ctx)
	if serverMeta == nil {
		return ""
	}
	return serverMeta.ServiceName
}

func initLogger(copeConf *CopeConf) (err error) {
	if copeConf.Logs.ConsoloLog {
		logs.NewConsoleOutputer()
		logs.AddOutputer(logs.NewConsoleOutputer())
		return
	}

	filename := fmt.Sprintf("%s/%s.log", copeConf.Logs.Path, copeConf.ServiceName)
	outputer, err := logs.NewFileOutputer(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	level := logs.GetLogLevel(copeConf.Logs.Level)
	logs.InitLogger(level, copeConf.Logs.ChanSize, copeConf.ServiceName)
	logs.AddOutputer(outputer)
	return
}
