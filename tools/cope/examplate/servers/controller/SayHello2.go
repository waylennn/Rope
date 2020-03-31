
package controller
import(
	"context"
	hello "cope/tools/cope/examplate/servers/generate"
)
//SayHello2Controller ..
type SayHello2Controller struct {
}


//CheckParams 检查请求参数，如果该函数返回错误，则Run函数不会执行
func (s *SayHello2Controller) CheckParams(ctx context.Context, r*hello.HelloRequest) (err error) {
	return
}

//Run SayHello函数的实现
func (s *SayHello2Controller) Run(ctx context.Context, r*hello.HelloRequest) (
	resp*hello.HelloResponse, err error) {
	resp = &hello.HelloResponse{
		Reply: "hehe",
	}
	return
}

