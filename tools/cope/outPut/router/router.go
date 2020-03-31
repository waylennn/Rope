package router

import (
	"context"
	"cope/meta"
	"cope/server"
	"cope/tools/cope/outPut/controller"
	hello "cope/tools/cope/outPut/generate"
)

type RouterServer struct{}

//SayHello ...
func (s *RouterServer) SayHello(ctx context.Context, r *hello.HelloRequest) (resp *hello.HelloResponse, err error) {

	ctx = meta.InitServerMeta(ctx, "hello", "SayHello")
	mwFunc := server.BuildServerMiddleware(mwSayHello)
	res, err := mwFunc(ctx, r)
	if err != nil {
		return
	}
	resp = res.(*hello.HelloResponse)

	return
}

func mwSayHello(ctx context.Context, request interface{}) (resp interface{}, err error) {
	r := request.(*hello.HelloRequest)
	ctrl := &controller.SayHelloController{}
	if err = ctrl.CheckParams(ctx, r); err != nil {
		return
	}

	if resp, err = ctrl.Run(ctx, r); err != nil {
		return
	}

	return
}

//SayHello2 ...
func (s *RouterServer) SayHello2(ctx context.Context, r *hello.HelloRequest) (resp *hello.HelloResponse, err error) {

	ctx = meta.InitServerMeta(ctx, "hello", "SayHello2")
	mwFunc := server.BuildServerMiddleware(mwSayHello2)
	res, err := mwFunc(ctx, r)
	if err != nil {
		return
	}
	resp = res.(*hello.HelloResponse)
	return
}

func mwSayHello2(ctx context.Context, request interface{}) (resp interface{}, err error) {
	r := request.(*hello.HelloRequest)
	ctrl := &controller.SayHello2Controller{}
	if err = ctrl.CheckParams(ctx, r); err != nil {
		return
	}

	if resp, err = ctrl.Run(ctx, r); err != nil {
		return
	}

	return
}
