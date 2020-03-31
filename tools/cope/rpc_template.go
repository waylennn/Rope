package main

var rpcClientTemplate = `
package {{.Package.Name}}c

import (
	"context"
	"cope/logs"
	"cope/meta"
	"cope/rpc"
	{{.Package.Name}} "{{.Prefix}}generate"
	"fmt"
	"google.golang.org/grpc"
)


type {{Capitalize .Package.Name}}Client struct {
	serviceName string
	client      *rpc.CopeClient

}

func New{{Capitalize .Package.Name}}Client(serviceName string, opts ...rpc.RpcOptionFunc) *{{Capitalize .Package.Name}}Client {
	c :=  &{{Capitalize .Package.Name}}Client{
		serviceName: serviceName,
		client:      rpc.NewCopeClient(serviceName, opts...),

	}
	return c
}

{{range .RPC}}
func (s *{{Capitalize $.Package.Name}}Client) {{.Name}}(ctx context.Context, r*{{$.Package.Name}}.{{.RequestType}})(resp*{{$.Package.Name}}.{{.ReturnsType}}, err error){
	// middlewareFunc := rpc.BuildClientMiddleware()
	// mkResp, err := middlewareFunc(ctx, r)

	mkResp, err := s.client.Call(ctx, "{{.Name}}", r, mwClient{{.Name}})
	resp, ok := mkResp.(*{{$.Package.Name}}.{{.ReturnsType}})
	if !ok {
		err = fmt.Errorf("invalid resp, not *{{$.Package.Name}}.{{.ReturnsType}}")
		return nil, err
	}
	return resp, err
}


func mwClient{{.Name}}(ctx context.Context, request interface{}) (resp interface{}, err error) {
	// conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	// if err != nil {
	// 	logs.Error(ctx, "grpc client start failed err :", err)
	// 	return
	// }
	// defer conn.Close()

	rpcMeta := meta.GetRpcMeta(ctx)
	c := {{$.Package.Name}}.New{{$.Service.Name}}Client(rpcMeta.Conn)

	req := request.(*{{$.Package.Name}}.{{.RequestType}})
	resp, err = c.{{.Name}}(ctx, req)
	if err != nil {
		logs.Error(ctx, "grpc client start failed err :", err)
		return 
	}
	return
}
{{end}}
`
