package helloc

import (
	"context"
	"cope/logs"
	"cope/meta"
	"cope/rpc"
	hello "cope/tools/cope/examplate/rpcClientEx/generate"
	"fmt"
)

type HelloClient struct {
	serviceName string
	client      *rpc.CopeClient
}

func NewHelloClient(serviceName string, opts ...rpc.RpcOptionFunc) *HelloClient {
	copeConf := &rpc.CopeConf{
		&rpc.Logs{
			ConsoloLog:  true,
			Path:        "./logs/",
			ServiceName: serviceName,
			Level:       "debug",
			ChanSize:    10000,
		},
	}
	c := &HelloClient{
		serviceName: serviceName,
		client:      rpc.NewCopeClient(serviceName, copeConf, opts...),
	}
	return c
}

func (s *HelloClient) SayHello(ctx context.Context, r *hello.HelloRequest) (resp *hello.HelloResponse, err error) {
	// middlewareFunc := rpc.BuildClientMiddleware()
	// mkResp, err := middlewareFunc(ctx, r)

	mkResp, err := s.client.Call(ctx, "SayHello", r, mwClientSayHello)
	resp, ok := mkResp.(*hello.HelloResponse)
	if !ok {
		err = fmt.Errorf("invalid resp, not *hello.HelloResponse")
		return nil, err
	}
	return resp, err
}

func mwClientSayHello(ctx context.Context, request interface{}) (resp interface{}, err error) {
	// conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	// if err != nil {
	// 	logs.Error(ctx, "grpc client start failed err :", err)
	// 	return
	// }
	// defer conn.Close()

	rpcMeta := meta.GetRpcMeta(ctx)
	c := hello.NewHelloServiceClient(rpcMeta.Conn)

	req := request.(*hello.HelloRequest)
	resp, err = c.SayHello(ctx, req)
	if err != nil {
		logs.Error(ctx, "grpc client start failed err :%v", err)
		return
	}
	return
}

func (s *HelloClient) SayHello2(ctx context.Context, r *hello.HelloRequest) (resp *hello.HelloResponse, err error) {
	// middlewareFunc := rpc.BuildClientMiddleware()
	// mkResp, err := middlewareFunc(ctx, r)

	mkResp, err := s.client.Call(ctx, "SayHello2", r, mwClientSayHello2)
	resp, ok := mkResp.(*hello.HelloResponse)
	if !ok {
		err = fmt.Errorf("invalid resp, not *hello.HelloResponse")
		return nil, err
	}
	return resp, err
}

func mwClientSayHello2(ctx context.Context, request interface{}) (resp interface{}, err error) {
	// conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	// if err != nil {
	// 	logs.Error(ctx, "grpc client start failed err :", err)
	// 	return
	// }
	// defer conn.Close()

	rpcMeta := meta.GetRpcMeta(ctx)
	c := hello.NewHelloServiceClient(rpcMeta.Conn)

	req := request.(*hello.HelloRequest)
	resp, err = c.SayHello2(ctx, req)
	if err != nil {
		logs.Error(ctx, "grpc client start failed err :", err)
		return
	}
	return
}
