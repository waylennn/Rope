package main

import (
	"fmt"
	"os"
	"path"
	"text/template"
)

//RPCGenrator ...
type RPCGenrator struct {
}

//Run ...
func (r *RPCGenrator) Run(opt *Option, meta *ServiceMetaData) (err error) {
	clientPath := fmt.Sprintf("%s%s/%s/%s", opt.Output, "generate", "client", meta.Package.Name)
	err = os.MkdirAll(clientPath, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileName := path.Join(clientPath, "client.go")
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = r.reden(rpcClientTemplate, meta, file)

	return
}

func (r *RPCGenrator) reden(data string, meta *ServiceMetaData, file *os.File) (err error) {
	t := template.New("rpc").Funcs(templateFuncMap)

	t, err = t.Parse(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = t.Execute(file, meta)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
func init() {
	rpc := &RPCGenrator{}

	err := ClientRegister("rpc generator", rpc)
	if err != nil {
		fmt.Println(err)
		return
	}

}
