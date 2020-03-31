package main

import (
	"cope/util"
	"fmt"
	"os"
	"path"
	"text/template"
)

//MainGenerator ...
type MainGenerator struct {
}

//Run .
func (d *MainGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	fileName := path.Join(opt.Output, "main", fmt.Sprintf("main.go"))
	exist := util.IsFileExist(fileName)
	if exist {
		return
	}
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("w file failed :%v", err)
		return
	}
	defer file.Close()
	err = d.reder(file, mainTemplate, metaData)
	if err != nil {
		fmt.Printf("reder file template failed :%v", err)
		return
	}

	return
}

func (d *MainGenerator) reder(file *os.File, data string, metaData *ServiceMetaData) (err error) {
	t := template.New("main")
	t, err = t.Parse(data)
	if err != nil {
		fmt.Printf("reder controller template failed :%v", err)
		return
	}
	err = t.Execute(file, metaData)
	if err != nil {
		fmt.Printf("reder controller template failed :%v", err)
		return
	}

	return
}

func init() {
	main := &MainGenerator{}

	err := ServerRegister("main generator", main)
	if err != nil {
		fmt.Println(err)
		return
	}
}
