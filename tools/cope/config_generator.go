package main

import (
	"cope/util"
	"fmt"
	"html/template"
	"os"
	"path"
)

//ConfigGenerator 配置文件
type ConfigGenerator struct {
}

//Run ...
func (c *ConfigGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	err = c.envConfig(util.TEST_ENV, opt, metaData)
	if err != nil {
		return
	}

	err = c.envConfig(util.PRODUCT_ENV, opt, metaData)
	if err != nil {
		return
	}
	return
}

func (c *ConfigGenerator) envConfig(env string, opt *Option, metaData *ServiceMetaData) (err error) {
	fileName := path.Join(opt.Output, "conf", env, fmt.Sprintf("%v.yaml", metaData.Package.Name))
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		fmt.Printf("open file:%s failed, err:%v\n", fileName, err)
		return
	}

	err = c.render(file, configTemplate, metaData)
	if err != nil {
		return
	}

	defer file.Close()
	return
}

func (c *ConfigGenerator) render(file *os.File, data string, metaData *ServiceMetaData) (err error) {
	t := template.New("main")
	t, err = t.Parse(data)
	if err != nil {
		fmt.Println("config template init failed")
		return
	}
	err = t.Execute(file, metaData)
	if err != nil {
		fmt.Println("config template init failed")
		return
	}
	return
}

func init() {

	config := &ConfigGenerator{}

	err := ServerRegister("configGenerator", config)
	if err != nil {
		fmt.Printf("config init failed err:%v", err)
		return
	}
}
