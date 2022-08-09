package main

import (
	"github.com/leicc520/go-gin-http"
	"github.com/leicc520/go-gin-http/micro"
	"github.com/leicc520/go-orm"
	"githunb.com/leicc520/go-disc-srv/app"
	"githunb.com/leicc520/go-disc-srv/app/service"
)

func main() {
	micro.CmdInit(service.InitConfig)
	micro.InitMicroHttp()//初始化http协议
	defer func() { //资源的回收处理逻辑
		orm.GdbPoolSt.Release()
	}()
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
