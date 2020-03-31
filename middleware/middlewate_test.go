package middleware

import (
	"context"
	"fmt"
	"testing"
)

func TestMidlleware(t *testing.T) {

	outer := func(middle MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (res interface{}, err error) {
			fmt.Println("outer")
			res, err = middle(ctx, req)
			fmt.Println("outer")
			return
		}
	}

	middle1 := func(middle MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (res interface{}, err error) {
			fmt.Println("m1")
			res, err = middle(ctx, req)
			fmt.Println("m1")
			return
		}
	}

	middle2 := func(middle MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (res interface{}, err error) {
			fmt.Println("m2")
			res, err = middle(ctx, req)
			fmt.Println("m2")
			return
		}
	}

	test := func(ctx context.Context, req interface{}) (res interface{}, err error) {
		fmt.Println("func")
		return
	}

	allMiddles := Chain(outer, middle1, middle2)
	allMiddles(test)(context.Background(), "text")
}
