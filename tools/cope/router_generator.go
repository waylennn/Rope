package main

import (
	"fmt"
	"os"
	"path"
	"text/template"
)

//RouterGenerator ...
type RouterGenerator struct {
}

//Run .
func (r *RouterGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	fileName := path.Join(opt.Output, "router", fmt.Sprintf("router.go"))
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("w file failed :%v", err)
		return
	}
	defer file.Close()
	err = r.reder(file, routerTemplate, metaData)
	if err != nil {
		fmt.Printf("reder file template failed :%v", err)
		return
	}

	return
}

func (r *RouterGenerator) reder(file *os.File, data string, metaData *ServiceMetaData) (err error) {
	t := template.New("router")
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
	router := &RouterGenerator{}

	err := ServerRegister("router generator", router)
	if err != nil {
		fmt.Println(err)
		return
	}
}
