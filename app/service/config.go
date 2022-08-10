package service

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/leicc520/go-gin-http"
	"github.com/leicc520/go-gin-http/micro"
	"github.com/leicc520/go-orm"
)

const (
	NAME    = "go-disc-srv"
	VERSION = "v1.0.0"
)

type ConfigSt struct {
	micro.Config
	DbMaster orm.DbConfig `yaml:"dbmaster"`
	DbSlaver orm.DbConfig `yaml:"dbslaver"`
}

var Config *ConfigSt = nil

//初始化加载配置信息
func InitConfig() {
	workDir, _ := os.Getwd()
	filePath := fmt.Sprintf("config/%s-%s.yml", NAME, os.Getenv(core.DCENV))
	filePath  = filepath.Join(workDir, filePath)
	Config    = &ConfigSt{}
	Config.Load(filePath, Config)
	Config.App.Name    = NAME
	Config.App.Version = VERSION
}

