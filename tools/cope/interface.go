package main

import (
	"github.com/emicklei/proto"
)

//Generator .
type Generator interface {
	Run(opt *Option, metaData *ServiceMetaData) error
}

//ServiceMetaData 保存proto文件解析之后的元信息
type ServiceMetaData struct {
	Service *proto.Service
	Message []*proto.Message
	RPC     []*proto.RPC
	Package *proto.Package
	Prefix  string
}
