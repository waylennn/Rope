package rpc

import "cope/middleware"

//BuildClientMiddleware 把中间件函数串起来
func BuildClientMiddleware(handleFunc middleware.MiddlewareFunc) middleware.MiddlewareFunc {
	var middleList []middleware.Middleware

	if len(middleList) > 0 {
		middleChain := middleware.Chain(middleList[0], middleList[1:]...)
		return middleChain(handleFunc)
	}
	return handleFunc
}
