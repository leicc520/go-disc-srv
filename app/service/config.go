package service

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/leicc520/go-gin-http"
	"github.com/leicc520/go-gin-http/micro"
)

const (
	NAME    = "go-disc-srv"
	VERSION = "v1.0.0"
)


var Config *micro.Config = nil

//初始化加载配置信息
func InitConfig() {
	conf := os.Getenv(core.DCSRV)
	if len(conf) < 1 {//通过环境变量设置加载配置
		basePath, _ := os.Getwd()
		conf = filepath.Join(basePath, fmt.Sprintf("/config/%s-%s.yml", conf, os.Getenv(core.DCENV)))
	}
	Config = &micro.Config{}
	Config.Load(conf, Config)
	Config.App.Name = NAME
	Config.App.Version = VERSION
}

