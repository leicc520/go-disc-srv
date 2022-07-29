package main

import (
	"git.cht-group.net/leicc/go-core"
	"git.cht-group.net/leicc/go-disc-srv/app/coreCTRL"
	"git.cht-group.net/leicc/go-disc-srv/app/discCTRL"
	"git.cht-group.net/leicc/go-disc-srv/app/logic"
	"git.cht-group.net/leicc/go-orm"
	"git.cht-group.net/leicc/go-orm/log"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"github.com/leicc520/go-core"
	"github.com/leicc520/go-orm"
	"net/http"
	"os"
)

func main() {
	conf := os.Getenv("CONFIGYML")
	if len(conf) < 1 {//通过环境变量设置加载配置
		conf = "./config/go-disc-srv-loc.yml"
	}
	config := logic.NewConfig(conf)
	log.SetLogger(config.Logger.Init())
	orm.InitDBPoolSt().LoadDbConfig(config) //配置数据库结构注册到数据库调用配置当中
	log.Write(-1, "env", os.Getenv("DCENV"), os.Getenv("DCSRV"))
	defer func() { //资源的回收处理逻辑
		orm.GdbPoolSt.Release()
	}()
	orm.SetCachePrefix("dc")//独立缓存策略
	core.NewApp(&config.App).RegHandler(func(c *gin.Engine) {
		staticBox := packr.New("webv5","./webroot")
		c.StaticFS("webv5", staticBox).GET("/", func(context *gin.Context) {
			context.Header("Content-Type", "text/html; charset=utf-8")
			indexStr, _ := staticBox.FindString("index.html")
			context.String(http.StatusOK, indexStr)
		})
	}).RegHandler(coreCTRL.Register).RegHandler(discCTRL.Register).Start() //启动服务
}
