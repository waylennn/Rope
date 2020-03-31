package main

var controllerTemplate = `
package controller
import(
	"context"
	{{.PackageName}} "{{.Prefix}}generate"
)
//{{.RPC.Name}}Controller ..
type {{.RPC.Name}}Controller struct {
}


//CheckParams 检查请求参数，如果该函数返回错误，则Run函数不会执行
func (s *{{.RPC.Name}}Controller) CheckParams(ctx context.Context, r*{{.PackageName}}.{{.RPC.RequestType}}) (err error) {
	return
}

//Run SayHello函数的实现
func (s *{{.RPC.Name}}Controller) Run(ctx context.Context, r*{{.PackageName}}.{{.RPC.RequestType}}) (
	resp*{{.PackageName}}.{{.RPC.ReturnsType}}, err error) {
	resp = &hello.HelloResponse{
		Reply: "hehe",
	}
	return
}

`
