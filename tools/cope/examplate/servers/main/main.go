package main

import (
	"cope/logs"
	"cope/server"
	pd "cope/tools/cope/examplate/servers/generate"
	"cope/tools/cope/examplate/servers/router"
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
