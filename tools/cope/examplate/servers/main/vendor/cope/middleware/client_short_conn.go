package middleware

import (
	"context"
	"fmt"

	"cope/meta"

	"cope/errno.go"

	"google.golang.org/grpc"
)

//ClientShortConnMiddleware 连接服务端
func ClientShortConnMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {

		rpcMeta := meta.GetRpcMeta(ctx)
		if rpcMeta.CurNode == nil {
			err = errno.InvalidNode
			return
		}

		address := fmt.Sprintf("%s:%d", rpcMeta.CurNode.IP, rpcMeta.CurNode.Port)
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			err = errno.ConnFailed
			return
		}
		rpcMeta.Conn = conn
		defer conn.Close()

		resp, err = next(ctx, req)
		return
	}
}
