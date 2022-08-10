package app

import (
	"github.com/gin-gonic/gin"
	"github.com/leicc520/go-gin-http"
	"github.com/leicc520/go-orm/log"
	"github.com/leicc520/go-disc-srv/app/service"
	"os"
)


//注册核心业务的http请求方法
func MicroRegister(weApp *gin.Engine) {
	hRouter := weApp.Group("/micsrv").Use(XTCheck())
	hRouter.GET("/discover/:proto/:name", doDiscover)
	hRouter.GET("/config/:name", doConfig)
	hRouter.POST("/unregister", doUnRegister)
	hRouter.POST("/register", doRegister)
	hRouter.GET("/reload", doReload)
	if service.HttpPools != nil { //启动的时候加载
		service.HttpPools.Load("http")
		log.Write(log.INFO, "start http discover service")
	}
	if service.GrpcPools != nil { //启动的时候加载
		service.GrpcPools.Load("grpc")
		log.Write(log.INFO, "start grpc discover service")
	}
}

//做最简单的token校验即可  一般服务发现只开通内网
func XTCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token  := c.GetHeader("X-TOKEN")
		jwtKey := os.Getenv(core.DCJWT)
		log.Write(log.INFO, "xtcheck token {"+token+"--->"+jwtKey+"}")
		if token != jwtKey {
			log.Write(log.ERROR, "jwt check error.")
			c.AbortWithStatus(403)
			return
		}
		c.Next()
	}
}

//注册核心业务的http请求方法
func WebRegister(weApp *gin.Engine)  {
	weApp.MaxMultipartMemory = 64 * 1024 * 1024 //上传文件-64MB
	if  service.Config != nil && len(service.Config.App.UpFileDir) > 1 {
		weApp.Static("/upfile", service.Config.App.UpFileDir)
	}
	weApp.GET("/captcha", doCaptcha)
	weApp.POST("/captcha/json", doCaptchaJson)
	weApp.POST("/captcha/check", doCheckCaptcha)
	weApp.POST("/signin/check", signInCheck)
	hRouter := weApp.Group("/api/core").Use(core.GINJWTCheck())
	hRouter.POST("/user/list", sysUserList)
	hRouter.POST("/user/safe", sysUserSafe)
	hRouter.POST("/user/update", sysUserUpdate)
	hRouter.POST("/signin/recheck", signInReCheck)
	hRouter.POST("/misrv/reload", sysmsrvReload)
	hRouter.POST("/misrv/status", sysmsrvStatus)
	hRouter.POST("/misrv/list", sysmsrvList)
	hRouter.POST("/yaml/list", sysYamlList)
	hRouter.POST("/yaml/update", sysYamlUpdate)
	hRouter.POST("/yaml/delete", sysYamlDelete)
}
