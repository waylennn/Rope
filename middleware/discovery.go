package middleware

import (
	"context"
	"cope/meta"
	"cope/registry"
	"fmt"
)

//NewDiscoveryMiddleware  服务发现
func NewDiscoveryMiddleware(r registry.Registry) Middleware {
	return func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			rpcMeta := meta.GetRpcMeta(ctx)
			if len(rpcMeta.AllNodes) > 0 {
				return next(ctx, req)
			}
			service, err := r.GetService(ctx, rpcMeta.ServiceName)
			if err != nil {
				fmt.Println(err)
				return
			}

			rpcMeta.AllNodes = service.Nodes
			resp, err = next(ctx, req)
			return
		}
	}
}
