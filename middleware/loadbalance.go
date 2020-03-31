package middleware

import (
	"context"
	"cope/loadbalance"
	"cope/meta"

	"cope/errno.go"

	"fmt"
)

//NewLoadBalanceMiddleware 提供负载均衡功能
func NewLoadBalanceMiddleware(loadBalance loadbalance.LoadBalance) Middleware {
	return func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			//从ctx中获取rpc的metadata
			rpcMeta := meta.GetRpcMeta(ctx)
			lenNodes := len(rpcMeta.AllNodes)
			if lenNodes == 0 {
				fmt.Println("nodes not exsize")
				err = errno.AllNodeFailed
				return
			}
			//生成loadbalance的上下文,用来过滤已经选择的节点
			// ctx = loadbalance.WithBalanceContext(ctx)
			for {
				//实现负载均衡
				rpcMeta.CurNode, err = loadBalance.Select(ctx, rpcMeta.AllNodes)
				fmt.Println("cur node:", rpcMeta.CurNode)
				if err != nil {
					fmt.Printf("load balance failed err:%s", err)
					return
				}
				resp, err = next(ctx, req)
				if err != nil {
					if errno.IsConnectError(err) {
						continue
					}
					return
				}
				break
			}
			return
		}
	}
}
