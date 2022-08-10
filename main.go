package main

import (
	"github.com/leicc520/go-disc-srv/app"
	"github.com/leicc520/go-disc-srv/app/migrate"
	"github.com/leicc520/go-disc-srv/app/service"
	"github.com/leicc520/go-gin-http"
	"github.com/leicc520/go-gin-http/micro"
	"github.com/leicc520/go-orm"
)

func main() {
	micro.CmdInit(service.InitConfig)
	defer func() { //资源的回收处理逻辑
		orm.GdbPoolSt.Release()
	}()
	migrate.InitDBCheck()
	/*.RegHandler(func(c *gin.Engine) {
			staticBox := packr.New("webv5","./webroot")
			c.StaticFS("webv5", staticBox).GET("/", func(context *gin.Context) {
				context.Header("Content-Type", "text/html; charset=utf-8")
				indexStr, _ := staticBox.FindString("index.html")
				context.String(http.StatusOK, indexStr)
			})
		})*/
	core.NewApp(&service.Config.App).RegHandler(app.WebRegister).RegHandler(app.MicroRegister).Start() //启动服务
}
