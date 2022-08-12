package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"github.com/leicc520/go-disc-srv/app"
	"github.com/leicc520/go-disc-srv/app/migrate"
	"github.com/leicc520/go-disc-srv/app/service"
	"github.com/leicc520/go-gin-http"
	"github.com/leicc520/go-gin-http/micro"
)

func main() {
	micro.CmdInit(service.InitConfig)
	migrate.InitDBCheck() //数据库初始化
	core.NewApp(&service.Config.App).RegHandler(func(c *gin.Engine) {
		staticBox := packr.New("webv5","./dist")
		c.StaticFS("webv5", staticBox).GET("/", func(context *gin.Context) {
			context.Header("Content-Type", "text/html; charset=utf-8")
			indexStr, _ := staticBox.FindString("index.html")
			context.String(200, indexStr)
		})
	}).RegHandler(app.Register).Start() //启动服务
}
