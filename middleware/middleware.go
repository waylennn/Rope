package middleware

import (
	"context"
)

//MiddlewareFunc ...
type MiddlewareFunc func(ctx context.Context, req interface{}) (res interface{}, err error)

//Middleware ...
type Middleware func(MiddlewareFunc) MiddlewareFunc

//Chain ...
func Chain(outer Middleware, others ...Middleware) Middleware {

	return func(next MiddlewareFunc) MiddlewareFunc {
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}

		return outer(next)
	}
}
