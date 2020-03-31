package main

import (
	"cope/logs"
	"cope/server"
	pd "cope/tools/cope/outPut/generate"
	"cope/tools/cope/outPut/router"
)

var routerServer = &router.RouterServer{}

func main() {
	err := server.Init("hello")
	if err != nil {
		logs.Stop()
		return
	}
	pd.RegisterHelloServiceServer(server.GRPCServer(), routerServer)
	server.Run()
}
