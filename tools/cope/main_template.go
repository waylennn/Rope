package main

var mainTemplate = `
package main
import (
	"cope/server"
	pd "{{.Prefix}}generate"
	"{{.Prefix}}router"
	"cope/logs"
)


var routerServer = &router.RouterServer{}
func main() {
	err := server.Init("{{.Package.Name}}")
	if err != nil {
		logs.Stop()		
		return
	}
	pd.Register{{ .Service.Name }}Server(server.GRPCServer(), routerServer)
	server.Run()
}
`
