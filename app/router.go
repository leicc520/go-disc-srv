package app

import (
	"os"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/leicc520/go-disc-srv/app/service"
	"github.com/leicc520/go-gin-http"
	"github.com/leicc520/go-orm/log"
)

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
func Register(weApp *gin.Engine)  {
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
	//微服务接口数据信息管理
	xRouter := weApp.Group("/micsrv").Use(XTCheck())
	xRouter.GET("/discover/:proto/:name", doDiscover)
	xRouter.GET("/config/:name", doConfig)
	xRouter.POST("/unregister", doUnRegister)
	xRouter.POST("/register", doRegister)
	xRouter.GET("/reload", doReload)
	time.AfterFunc(time.Millisecond*100, func() {//延迟加载数据
		if service.HttpPools != nil { //启动的时候加载
			service.HttpPools.Load("http")
			log.Write(log.INFO, "start http discover service")
		}
		if service.GrpcPools != nil { //启动的时候加载
			service.GrpcPools.Load("grpc")
			log.Write(log.INFO, "start grpc discover service")
		}
	})
}
