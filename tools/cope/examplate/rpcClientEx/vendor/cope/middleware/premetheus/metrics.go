package prometheus

import (
	"context"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/status"
)

// Metrics 服务端采样打点
type Metrics struct {
	requestCounter *prom.CounterVec
	codeCounter    *prom.CounterVec
	latencySummary *prom.SummaryVec
}

//NewServerMetrics ...
func NewServerMetrics() *Metrics {
	return &Metrics{
		requestCounter: promauto.NewCounterVec(
			prom.CounterOpts{
				Name: "cope_server_request_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure.",
			}, []string{"service", "method"}),
		codeCounter: promauto.NewCounterVec(
			prom.CounterOpts{
				Name: "cope_server_handled_code_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure.",
			}, []string{"service", "method", "grpc_code"}),
		latencySummary: promauto.NewSummaryVec(
			prom.SummaryOpts{
				Name:       "cope_proc_cost",
				Help:       "RPC latency distributions.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"service", "method"},
		),
	}
}

//NewRPCMetrics 生成实例
func NewRPCMetrics() *Metrics {
	return &Metrics{
		requestCounter: promauto.NewCounterVec(
			prom.CounterOpts{
				Name: "cope_rpc_call_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure.",
			}, []string{"service", "method"}),
		codeCounter: promauto.NewCounterVec(
			prom.CounterOpts{
				Name: "cope_rpc_code_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure.",
			}, []string{"service", "method", "grpc_code"}),
		latencySummary: promauto.NewSummaryVec(
			prom.SummaryOpts{
				Name:       "cope_rpc_cost",
				Help:       "RPC latency distributions.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"service", "method"},
		),
	}
}

//IncrRequest 统计请求数
func (m *Metrics) IncrRequest(ctx context.Context, serviceName, methodName string) {
	m.requestCounter.WithLabelValues(serviceName, methodName).Inc()
}

//IncrCode 统计错误数
func (m *Metrics) IncrCode(ctx context.Context, serviceName, methodName string, err error) {
	st, _ := status.FromError(err)
	m.codeCounter.WithLabelValues(serviceName, methodName, st.Code().String()).Inc()
}

//Latency 统计延迟
func (m *Metrics) Latency(ctx context.Context, serviceName, methodName string, us int64) {

	m.latencySummary.WithLabelValues(serviceName, methodName).Observe(float64(us))
}
