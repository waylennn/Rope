package server

import (
	"context"
	"cope/logs"
	"cope/middleware"
	"cope/registry"
	_ "cope/registry/etcd"
	"cope/util"
	"fmt"
	"net"
	"net/http"

	"golang.org/x/time/rate"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport/zipkin"
)

//CopeServer ...
type CopeServer struct {
	*grpc.Server
	limiter  *rate.Limiter
	register registry.Registry
}

var copeServer = &CopeServer{
	Server: grpc.NewServer(),
}

//Init ....
func Init(serviceName string) (err error) {

	//初始化配置文件相关
	err = InitConfig(serviceName)
	if err != nil {
		return
	}
	//初始化日志
	err = initLogger()
	if err != nil {
		return
	}

	//初始化限流器
	if copeConf.Limit.SwitchOn {
		copeServer.limiter = rate.NewLimiter(rate.Limit(copeConf.Limit.QPSLimit),
			copeConf.Limit.QPSLimit)
	}

	//初始化etcd
	if copeConf.Regiser.SwitchOn {
		err = initRegisterService()
		if err != nil {
			return
		}
	}

	//初始化分布式追踪
	if copeConf.Trace.SwitchOn {
		err = initTrace(copeConf.ServiceName)
		if err != nil {
			return
		}
	}

	return
}

func initLogger() (err error) {
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

//initRegisterService 初始化服务注册
func initRegisterService() (err error) {
	if !copeConf.Regiser.SwitchOn {
		return
	}

	ctx := context.TODO()
	registryInit, err := registry.InitRegistry(
		ctx,
		copeConf.Regiser.RegisterName,
		registry.WithAddrs([]string{copeConf.Regiser.RegisterAddr}),
		registry.WithTimeout(copeConf.Regiser.Timeout),
		registry.WithRegistryPath(copeConf.Regiser.RegisterPath),
		registry.WithHeartBeat(copeConf.Regiser.HeartBeat),
	)

	if err != nil {
		logs.Error(ctx, "init registe etcd failed %v", err)
		return
	}

	copeServer.register = registryInit
	service := &registry.Service{
		Name: copeConf.ServiceName,
	}

	ip, err := util.GetLocalIP()
	if err != nil {
		logs.Error(ctx, "init registe etcd failed %v", err)
		return
	}
	service.Nodes = append(service.Nodes, &registry.Nodes{
		IP:   ip,
		Port: copeConf.Port,
	})

	err = registryInit.Register(ctx, service)
	if err != nil {
		logs.Error(ctx, "init registe etcd failed %v", err)
		return
	}
	return
}

//initTrace 初始化分布式追踪
func initTrace(serverName string) (err error) {
	transport, err := zipkin.NewHTTPTransport(
		copeConf.Trace.ReportAddr,
		zipkin.HTTPBatchSize(16),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		logs.Error(context.TODO(), "ERROR: cannot init Jaeger: %v\n", err)
	}

	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  copeConf.Trace.SampleType,
			Param: copeConf.Trace.SampleRate,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	r := jaeger.NewRemoteReporter(transport)
	tracer, closer, err := cfg.New(serverName,
		config.Logger(jaeger.StdLogger),
		config.Reporter(r))
	if err != nil {
		logs.Error(context.TODO(), "ERROR: cannot init Jaeger: %v\n", err)
	}

	_ = closer
	opentracing.SetGlobalTracer(tracer)
	return
}

//Run 这里启动服务相关,方便代码扩展
func Run() {

	if GetConf().Prometheus.SwitchOn {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			addr := fmt.Sprintf("0.0.0.0:%d", GetConf().Prometheus.Port)
			http.ListenAndServe(addr, nil)
		}()
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", GetConf().Port))
	if err != nil {
		return
	}
	copeServer.Serve(lis)
}

//GRPCServer 对外提供server参数
func GRPCServer() *grpc.Server {
	return copeServer.Server
}

//中间件相关
var userMiddleware []middleware.Middleware

//Use 注册中间件
func Use(m ...middleware.Middleware) {
	userMiddleware = append(userMiddleware, m...)
}

//BuildServerMiddleware 把中间件函数串起来
func BuildServerMiddleware(handleFunc middleware.MiddlewareFunc) middleware.MiddlewareFunc {

	var middleList []middleware.Middleware
	middleList = append(middleList, middleware.PrometheusServerMiddleware)

	if copeConf.Limit.SwitchOn {
		middleList = append(middleList, middleware.RateLimitMiddleware(copeServer.limiter))
	}

	if copeConf.Trace.SwitchOn {
		middleList = append(middleList, middleware.DistributeTraceMiddleware)
	}

	if len(userMiddleware) > 0 {
		middleList = append(middleList, userMiddleware...)
	}

	if len(middleList) > 0 {
		middleChain := middleware.Chain(middleList[0], middleList[1:]...)
		return middleChain(handleFunc)
	}
	return handleFunc
}
