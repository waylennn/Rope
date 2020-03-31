package main

import (
	"fmt"
	"os"
	"os/exec"
)

//GrpcGenerator 创建目录
type GrpcGenerator struct {
}

//Run .
func (d *GrpcGenerator) Run(opt *Option, metaData *ServiceMetaData) error {
	outputParams := fmt.Sprintf("plugins=grpc:%sgenerate", opt.Output)
	cmd := exec.Command("protoc", "--go_out", outputParams, opt.Proto3Filename)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Printf("cmd.Run faile :%s", err)
		return err
	}
	return nil
}

func init() {
	dir := &GrpcGenerator{}

	err := ServerRegister("grpc generator", dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ClientRegister("grpc generator", dir)
	if err != nil {
		fmt.Println(err)
		return
	}
}
