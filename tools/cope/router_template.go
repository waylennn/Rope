package main

var routerTemplate = `package router
import (
	"context"
	{{.Package.Name}} "{{.Prefix}}generate"
	"{{.Prefix}}controller"
	"cope/server"
	"cope/meta"

)

type RouterServer struct{}

{{range .RPC}}
//{{.Name}} ...
func (s *RouterServer) {{.Name}}(ctx context.Context, r *{{$.Package.Name}}.{{.RequestType}})(resp*{{$.Package.Name}}.{{.ReturnsType}}, err error){
	
	ctx = meta.InitServerMeta(ctx,"{{$.Package.Name}}", "{{.Name}}")
	mwFunc := server.BuildServerMiddleware(mw{{.Name}})
	res, err := mwFunc(ctx, r)
	if err != nil {
		return
	}
	resp = res.(*{{$.Package.Name}}.{{.ReturnsType}})

	return
}

func mw{{.Name}}(ctx context.Context, request interface{}) (resp interface{}, err error) {
	r := request.(*{{$.Package.Name}}.{{.RequestType}})
	ctrl := &controller.{{.Name}}Controller{}
	if err = ctrl.CheckParams(ctx, r); err != nil {
		return
	}

	if resp, err = ctrl.Run(ctx, r); err != nil {
		return
	}

	return
}

{{end}}
`
