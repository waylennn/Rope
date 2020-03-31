package controller

import (
	"context"
	"cope/logs"
	hello "cope/tools/cope/outPut/generate"
)

//SayHelloController ..
type SayHelloController struct {
}

//CheckParams 检查请求参数，如果该函数返回错误，则Run函数不会执行
func (s *SayHelloController) CheckParams(ctx context.Context, r *hello.HelloRequest) (err error) {
	return
}

//Run SayHello函数的实现
func (s *SayHelloController) Run(ctx context.Context, r *hello.HelloRequest) (
	resp *hello.HelloResponse, err error) {
	logs.Debug(ctx, "%v", "debug")
	resp = &hello.HelloResponse{
		Reply: "hehe",
	}
	return
}
