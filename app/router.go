package coreCTRL

import (
	"git.cht-group.net/leicc/go-core"
	"git.cht-group.net/leicc/go-disc-srv/app/logic"
	"github.com/gin-gonic/gin"
	"github.com/leicc520/go-core"
)


var regSrv *core.MicRegSrv = nil

//注册核心业务的http请求方法
func microRegister(wapp *gin.Engine) {
	regSrv = core.NewMicRegSrv(logic.GConfig.DiSrv)
	hRouter := wapp.Group("/micsrv").Use(XTCheck())
	hRouter.GET("/discover/:proto/:name", doDiscover)
	hRouter.GET("/config/:name", doConfig)
	hRouter.POST("/unregister", doUnRegister)
	hRouter.POST("/register", doRegister)
	hRouter.GET("/reload", doReload)
	if gHttpPools != nil { //启动的时候加载
		gHttpPools.Load("http")
		log.Write(log.INFO, "start http discover service")
	}
	if gGrpcPools != nil { //启动的时候加载
		gGrpcPools.Load("grpc")
		log.Write(log.INFO, "start grpc discover service")
	}
}

//做最简单的token校验即可  一般服务发现只开通内网
func XTCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-TOKEN")
		log.Write(log.INFO, "xtcheck token {"+token+"--->"+regSrv.JwtKey+"}")
		if token != regSrv.JwtKey {
			log.Write(log.ERROR, "jwt check error.")
			c.AbortWithStatus(403)
			return
		}
		c.Next()
	}
}

//注册核心业务的http请求方法
func appRegister(app *gin.Engine)  {
	go logic.DbLogService()//开启数据日志服务
	app.MaxMultipartMemory = 64 * 1024 * 1024 //上传文件-64MB
	if  logic.GConfig != nil && len(logic.GConfig.App.UpFileDir) > 1 {
		app.Static("/upfile", logic.GConfig.App.UpFileDir)
	}
	app.GET("/captcha", doCaptcha)
	app.GET("/ueditor", uEditorAction)
	app.POST("/ueditor", uEditorAction)
	app.POST("/captcha/json", doCaptchaJson)
	app.POST("/captcha/check", doCheckCaptcha)
	app.POST("/signin/check", signInCheck)
	hRouter := app.Group("/api/core")
	hRouter.POST("/config/dict", sysConfigDict)
	hRouter.Use(core.JWTACLCheck()).Use(logic.SSOACLCheck(0))
	hRouter.POST("/signin/recheck", signInReCheck)
	hRouter.POST("/upfile/image", sysUpfileImage)
	hRouter.POST("/share/tables", doTables)
	hRouter.POST("/share/sysorg", doSysOrg)
	hRouter.POST("/share/sysuser", doSysUser)
	hRouter.POST("/share/sysrole", doSysRole)
	hRouter.POST("/share/configs", doConfig)
	hRouter.POST("/share/module", doSysModule)
	hRouter.POST("/user/mine", sysUserMine)
	hRouter.POST("/log/list", syslogList)
	hRouter.Use(logic.ModuleACLCheck())//需要做权限验证的放这个后面
	hRouter.POST("/org/tree", sysOrgTree)
	hRouter.POST("/org/delete", sysOrgDelete)
	hRouter.POST("/org/update", sysOrgUpdate)
	hRouter.POST("/user/list", sysUserList)
	hRouter.POST("/user/update", sysUserUpdate)
	hRouter.POST("/user/roleids", sysUserRoleIds)
	hRouter.POST("/user/getaccess", sysUserGetAccess)
	hRouter.POST("/user/setaccess", sysUserSetAccess)
	hRouter.POST("/role/list", sysRoleList)
	hRouter.POST("/role/delete", sysRoleDelete)
	hRouter.POST("/role/update", sysRoleUpdate)
	hRouter.POST("/role/getaccess", sysRoleGetAccess)
	hRouter.POST("/role/setaccess", sysRoleSetAccess)
	hRouter.POST("/config/list", sysConfigList)
	hRouter.POST("/config/delete", sysConfigDelete)
	hRouter.POST("/config/update", sysConfigUpdate)
	hRouter.POST("/module/tree", sysModuleTrees)
	hRouter.POST("/module/all", sysSysModuleAllTree)
	hRouter.POST("/module/update", sysModuleUpdate)
	hRouter.POST("/module/delete", sysModuleDelete)
	hRouter.POST("/misrv/reload", sysmsrvReload)
	hRouter.POST("/misrv/status", sysmsrvStatus)
	hRouter.POST("/misrv/list", sysmsrvList)
	hRouter.POST("/yaml/list", sysYamlList)
	hRouter.POST("/yaml/update", sysYamlUpdate)
	hRouter.POST("/yaml/delete", sysYamlDelete)
}
