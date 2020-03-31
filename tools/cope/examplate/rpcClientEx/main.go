package main

import (
	"context"
	"cope/rpc"
	hello "cope/tools/cope/examplate/rpcClientEx/generate"
	helloc "cope/tools/cope/examplate/rpcClientEx/generate/client/hello"
	"fmt"
	"time"
)

func main() {
	helloClient := helloc.NewHelloClient("HelloService", rpc.WithLimitQPS(3))

	for {
		ctx := context.Background()
		res, err := helloClient.SayHello(ctx, &hello.HelloRequest{Name: "test my client"})
		if err != nil {
			//logs.Error(ctx, "employ sayhello failed err:%v", err)
			time.Sleep(time.Millisecond * 100)
			continue
		}
		fmt.Println(res)
		time.Sleep(time.Millisecond * 1000)

	}

}
