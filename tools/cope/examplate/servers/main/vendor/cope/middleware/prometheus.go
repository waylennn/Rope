package middleware

import (
	"context"
	"cope/meta"
	prometheus "cope/middleware/premetheus"
	"time"
)

var (
	//DefaultServerMetrics 创建监控的服务实例
	DefaultServerMetrics = prometheus.NewServerMetrics()
)

// func init() {
// 	go func() {
// 		http.Handle("/metrics", promhttp.Handler())
// 		addr := fmt.Sprintf("0.0.0.0:%d", 8888)
// 		http.ListenAndServe(addr, nil)
// 	}()
// }

//PrometheusServerMiddleware 监控请求数，错误，延迟
func PrometheusServerMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		serverMeta := meta.GetServerMeta(ctx)
		resp, err = next(ctx, req)
		DefaultServerMetrics.IncrRequest(ctx, serverMeta.ServiceName, serverMeta.Method)
		startTime := time.Now()
		resp, err = next(ctx, req)

		DefaultServerMetrics.IncrCode(ctx, serverMeta.ServiceName, serverMeta.Method, err)
		DefaultServerMetrics.Latency(ctx, serverMeta.ServiceName,
			serverMeta.Method, time.Since(startTime).Nanoseconds()/1000)

		return
	}

}
