
package controller
import(
	"context"
	hello "cope/tools/cope/examplate/servers/generate"
)
//SayHelloController ..
type SayHelloController struct {
}


//CheckParams 检查请求参数，如果该函数返回错误，则Run函数不会执行
func (s *SayHelloController) CheckParams(ctx context.Context, r*hello.HelloRequest) (err error) {
	return
}

//Run SayHello函数的实现
func (s *SayHelloController) Run(ctx context.Context, r*hello.HelloRequest) (
	resp*hello.HelloResponse, err error) {
	resp = &hello.HelloResponse{
		Reply: "hehe",
	}
	return
}

