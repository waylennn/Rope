package middleware

import (
	"context"
	"cope/logs"

	"google.golang.org/grpc/metadata"
)

func PrepareMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {

		//处理traceId
		var traceId string
		//从ctx获取grpc的metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			vals, ok := md["trace_id"]
			if ok && len(vals) > 0 {
				traceId = vals[0]
			}
		}

		if len(traceId) == 0 {
			traceId = logs.GenTraceId()
		}

		ctx = logs.WithFieldContext(ctx)
		ctx = logs.WithTraceId(ctx, traceId)
		resp, err = next(ctx, req)
		return
	}
}