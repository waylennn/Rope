package middleware

import (
	"context"
	"cope/logs"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//Limiter 限流
type Limiter interface {
	Allow() bool
}

//RateLimitMiddleware 监控请求数，错误，延迟
func RateLimitMiddleware(l Limiter) Middleware {
	return func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (res interface{}, err error) {
			isAllow := l.Allow()
			// limiter.WaitN(ctx, 2)
			if !isAllow {
				err = status.Error(codes.ResourceExhausted, "rate limited")
				logs.Error(ctx, "rate limit :err:%v", err)
				return
			}
			res, err = next(ctx, req)

			return
		}
	}
}
