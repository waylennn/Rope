package main

import (
	"cope/util"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/emicklei/proto"
)

//CtrlGenerator 创建目录
type CtrlGenerator struct {
	service *proto.Service
	message []*proto.Message
	rpc     []*proto.RPC
}

//RPCMeta ..
type RPCMeta struct {
	Prefix      string
	RPC         *proto.RPC
	PackageName string
}

//Run .
func (d *CtrlGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	for _, rpc := range metaData.RPC {
		var file *os.File
		fileName := path.Join(opt.Output, "controller", fmt.Sprintf("%s.go", rpc.Name))
		exist := util.IsFileExist(fileName)
		if exist {
			return
		}
		file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Printf("w file failed :%v", err)
			return
		}
		defer file.Close()

		rpcMeta := &RPCMeta{
			Prefix:      opt.Prefix,
			RPC:         rpc,
			PackageName: metaData.Package.Name,
		}
		err = d.render(file, controllerTemplate, rpcMeta)
		if err != nil {
			return
		}
	}

	return

}

func (d *CtrlGenerator) render(file *os.File, data string, metaData *RPCMeta) (err error) {
	t := template.New("main")
	t, err = t.Parse(data)
	if err != nil {
		return
	}
	err = t.Execute(file, metaData)
	if err != nil {
		return
	}
	return
}

func init() {
	ctrl := &CtrlGenerator{}

	err := ServerRegister("contorller generator", ctrl)
	if err != nil {
		fmt.Println(err)
		return
	}
}
