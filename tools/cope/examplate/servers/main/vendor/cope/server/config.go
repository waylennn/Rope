package server

import (
	"cope/util"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

//CopeCofig ...
type CopeCofig struct {
	Port        int            `yaml:"port"`
	Prometheus  PrometheusConf `yaml:"prometheus"`
	ServiceName string         `yaml:"service_name"`
	Limit       `yaml:"limit"`
	Logs        `yaml:"log"`
	Regiser     RegisterConf `yaml:"register"`
	Trace       TraceConf    `yaml:"trace"`

	//内部的配置项
	ConfigDir  string `yaml:"-"`
	RootDir    string `yaml:"-"`
	ConfigFile string `yaml:"-"`
}

//Logs ...
type Logs struct {
	Level      string `yaml:"level"`
	ChanSize   int    `yaml:"chan_size"`
	ConsoloLog bool   `yaml:"true"`
	Path       string `yaml:"path"`
}

//PrometheusConf ...
type PrometheusConf struct {
	SwitchOn bool `yaml:"switch_on"`
	Port     int  `yaml:"port"`
}

//Limit 限流相关
type Limit struct {
	SwitchOn bool `yaml:"switch_on"`
	QPSLimit int  `yaml:"qps"`
}

//TraceConf 分布式追踪相关
type TraceConf struct {
	SwitchOn   bool    `yaml:"switch_on"`
	ReportAddr string  `yaml:"report_addr"`
	SampleType string  `yaml:"sample_type"`
	SampleRate float64 `yaml:"sample_rate"`
}

//RegisterConf etcd相关
type RegisterConf struct {
	SwitchOn     bool          `yaml:"switch_on"`
	RegisterPath string        `yaml:"register_path"`
	Timeout      time.Duration `yaml:"timeout"`
	HeartBeat    int64         `yaml:"heart_beat"`
	RegisterName string        `yaml:"register_name"`
	RegisterAddr string        `yaml:"register_addr"`
}

var (
	copeConf = &CopeCofig{
		Port: 8080,
		Prometheus: PrometheusConf{
			SwitchOn: true,
			Port:     8081,
		},
		ServiceName: "koala_server",
	}
)

func initDir(serviceName string) (err error) {
	exeFilePath := os.Args[0]
	if runtime.GOOS == "windows" {
		exeFilePath = strings.Replace(exeFilePath, "\\", "/", -1)
	}

	lastIdx := strings.LastIndex(exeFilePath, "/")
	if lastIdx < 0 {
		err = fmt.Errorf("invalid exe path:%v", exeFilePath)
		return
	}

	copeConf.RootDir = path.Join(strings.ToLower(exeFilePath[0:lastIdx+1]), "..")
	copeConf.ConfigDir = path.Join(copeConf.RootDir, "./conf/", util.GetEnv())
	copeConf.ConfigFile = path.Join(copeConf.ConfigDir, fmt.Sprintf("%s.yaml", serviceName))
	return
}

//InitConfig 读取并解析配置文件转化为结构体
func InitConfig(serviceName string) (err error) {
	err = initDir(serviceName)
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(copeConf.ConfigFile)
	if err != nil {
		fmt.Printf("read config file failed err:%v", err)
		return
	}

	err = yaml.Unmarshal(data, copeConf)
	if err != nil {
		fmt.Printf("unmarshal config file failed err:%v", err)
		return
	}
	fmt.Printf("init koala conf succ, conf:%#v\n", copeConf)
	return
}

//GetConfigDir 获取配置文件所在目录
func GetConfigDir() string {
	return copeConf.ConfigDir
}

//GetRootDir 获取目录"....\output\"
func GetRootDir() string {
	return copeConf.RootDir
}

//GetServerPort 获取端口
func GetServerPort() int {
	return copeConf.Port
}

//GetConf 获取配置文件结构体
func GetConf() *CopeCofig {
	return copeConf
}
