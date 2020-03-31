package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/emicklei/proto"
)

//AllDirList 需要创建的目录
var AllDirList []string = []string{
	"controller",
	"idl",
	"main",
	"scripts",
	"conf/test",
	"conf/product",
	"app/router",
	"app/config",
	"model",
	"router",
	"generate",
}

//GeneratorMgr ...
type GeneratorMgr struct {
	serverGeneratorMap map[string]Generator
	clientGeneratorMap map[string]Generator
	metaData           *ServiceMetaData
}

var genMgr *GeneratorMgr = &GeneratorMgr{
	serverGeneratorMap: make(map[string]Generator),
	clientGeneratorMap: make(map[string]Generator),
	metaData:           &ServiceMetaData{},
}

// Run .
func (g *GeneratorMgr) Run(opt *Option) (err error) {
	err = g.initFilePath(opt)
	if err != nil {
		fmt.Println("generator_mgr initDir  failed")
		return
	}

	g.metaData.Prefix = opt.Prefix

	err = g.parseMeta(opt)
	if err != nil {
		fmt.Println("generator_mgr parse meta data failed")
		return
	}

	err = g.createAllDir(opt)
	if err != nil {
		fmt.Println("generator_mgr create dir  failed")
		return
	}

	//如果存在 创建客户端不创建服务端
	if opt.GenClientCode {
		for _, gen := range g.clientGeneratorMap {
			err = gen.Run(opt, g.metaData)
			if err != nil {
				return
			}
		}
		return
	}

	//创建服务端
	for _, gen := range g.serverGeneratorMap {
		err = gen.Run(opt, g.metaData)
		if err != nil {
			return
		}
	}
	return
}

//初始化文件创建路径 和import 的绝对路径
func (g *GeneratorMgr) initFilePath(opt *Option) (err error) {
	goPath := os.Getenv("GOPATH")

	//如果用户指定了目录
	if len(opt.Prefix) > 0 {
		opt.Output = path.Join(goPath, "src", opt.Prefix)
		return
	}

	//获取当前执行文件的目录和相对src的目录
	curDir := os.Args[0]
	curDir, err = filepath.Abs(curDir)
	if err != nil {
		fmt.Printf("get curDir failed :%v", err)
		return
	}

	if runtime.GOOS == "windows" {
		curDir = strings.Replace(curDir, "\\", "/", -1)
		goPath = strings.Replace(goPath, "\\", "/", -1)
	}

	curDir = curDir[:strings.LastIndex(curDir, "/")+1]
	temp := fmt.Sprintf("%v/src/", goPath)
	srcPath := strings.Replace(curDir, temp, "", -1)

	// opt.Output = path.Join(curDir, "outPut")
	// opt.Prefix = path.Join(srcPath, "outPut") //相对路径
	opt.Output = curDir
	opt.Prefix = srcPath
	return
}

//创建文件夹
func (g *GeneratorMgr) createAllDir(opt *Option) error {
	if opt.GenClientCode {
		potoPath := fmt.Sprintf("%s%s", opt.Output, "generate")
		err := os.MkdirAll(potoPath, 0755)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}
	for _, dir := range AllDirList {
		fullDir := path.Join(opt.Output, dir)
		err := os.MkdirAll(fullDir, 0755)
		if err != nil {
			fmt.Printf("mkdir dir %s failed, err:%v\n", dir, err)
			return err
		}
	}

	return nil
}

// 解析元数据
func (g *GeneratorMgr) parseMeta(opt *Option) (err error) {
	//proto文件 而不是解析之后的文件
	reader, _ := os.Open(opt.Proto3Filename)
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		fmt.Printf("parse file failed:%s", err)
		return
	}
	proto.Walk(definition,
		proto.WithService(g.handleService),
		proto.WithMessage(g.handleMessage),
		proto.WithRPC(g.handleRPC),
		proto.WithPackage(g.handlePackage),
	)

	if err != nil {
		return
	}
	return
}

func (g *GeneratorMgr) handlePackage(p *proto.Package) {
	g.metaData.Package = p
}

func (g *GeneratorMgr) handleService(s *proto.Service) {
	g.metaData.Service = s
}

func (g *GeneratorMgr) handleMessage(m *proto.Message) {
	g.metaData.Message = append(g.metaData.Message, m)
}

func (g *GeneratorMgr) handleRPC(s *proto.RPC) {
	g.metaData.RPC = append(g.metaData.RPC, s)
}

// ServerRegister for genrator to regist
func ServerRegister(name string, gen Generator) (err error) {
	_, ok := genMgr.serverGeneratorMap[name]
	if ok {
		err = fmt.Errorf("generator is exists %s", name)
		return err
	}
	genMgr.serverGeneratorMap[name] = gen
	return
}

// ClientRegister for genrator to regist
func ClientRegister(name string, gen Generator) (err error) {
	_, ok := genMgr.clientGeneratorMap[name]
	if ok {
		err = fmt.Errorf("generator is exists %s", name)
		return err
	}
	genMgr.clientGeneratorMap[name] = gen
	return
}
