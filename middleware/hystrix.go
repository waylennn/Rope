package middleware

import (
	"context"
	"cope/meta"

	"github.com/afex/hystrix-go/hystrix"
)

//HystrixMiddleware 熔断中间件
func HystrixMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (interface{}, error) {

		rpcMeta := meta.GetRpcMeta(ctx)
		var resp interface{}

		hystrixErr := hystrix.Do(rpcMeta.ServiceName, func() (err error) {
			resp, err = next(ctx, req)
			return err
		}, nil)

		if hystrixErr != nil {
			return nil, hystrixErr
		}

		return resp, hystrixErr
	}
}
