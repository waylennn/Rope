package registry

import "context"

// Registry 服务注册插件的接口
type Registry interface {
	//插件的名字
	Name() string
	//初始化
	Init(ctx context.Context, opt ...Option) (err error)
	//服务注册
	Register(ctx context.Context, service *Service) (err error)
	//服务反注册
	Unregister(ctx context.Context, service *Service) (err error)
	//服务发现:通过服务的名字获取服务的位置信息(ip:port)
	GetService(ctx context.Context, name string) (service *Service, err error)
}
