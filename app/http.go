package app

import (
	"net"
	"net/http"
	"regexp"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/leicc520/go-disc-srv/app/models"
	"github.com/leicc520/go-disc-srv/app/service"
	"github.com/leicc520/go-gin-http"
	"github.com/leicc520/go-orm"
	"github.com/leicc520/go-orm/log"
)

// @Summary 微服务注册
// @Description 中台的微服务注册管理
// @Tags 中台管理
// @Param name formData string true "服务名称"
// @Param srv  formData string true "服务注册地址"
// @Param proto formData string true "服务注册的协议"
// @Param version formData string true "服务的版本号"
// @Success 200 {string} OK
// @Router /micsrv/register [post]
func doRegister(c *gin.Context) {
	args   := service.MicSrvNodeSt{}
	if err := c.ShouldBind(&args); err != nil {
		core.PanicHttpError(500, err.Error())
	}
	//拼接真实的服务地址 只提交端口的情况 要取一下请求IP
	if ok, err := regexp.MatchString("^[\\d]+$", args.Srv); ok && err == nil {
		srvIp, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
		args.Srv = srvIp + ":" + args.Srv
	}
	sorm  := models.NewSysMsrv()
	oldid := sorm.GetValue(func(st *orm.QuerySt) string {
		st.Where("srv", args.Srv)
		return st.GetWheres()
	}, "id").ToInt64()
	if oldid < 1 {//记录不存在的情况新增
		sorm.NewOneFromHandler(func(st *orm.QuerySt) *orm.QuerySt {
			st.Value("srv", args.Srv).Value("name", args.Name).Value("status", 1)
			st.Value("version", args.Version).Value("proto", args.Proto)
			st.Value("addtime", time.Now().Unix()).Value("stime", time.Now().Unix())
			return st
		}, nil)
	} else {//否则更新数据记录信息
		sorm.Save(oldid, orm.SqlMap{"status":1, "name":args.Name,
			"version":args.Version, "proto":args.Proto, "stime":time.Now().Unix()})
	}
	log.Write(log.INFO, oldid, "---->", args)
	//分别添加到对应的数据结构当中
	if args.Proto == "http" && service.HttpPools != nil {
		service.HttpPools.Put(args.Name, args.Srv, oldid)
	} else if args.Proto == "grpc" && service.GrpcPools != nil {
		service.GrpcPools.Put(args.Name, args.Srv, oldid)
	}
	c.JSON(http.StatusOK, gin.H{"Code": 0, "Msg": "OK", "Srv": args.Srv})
}

// @Summary 微服务注销
// @Description 中台的微服务注销管理
// @Tags 中台管理
// @Param name formData string true "服务名称"
// @Param srv  formData string true "服务注册地址"
// @Success 200 {string} OK
// @Router /micsrv/unregister [post]
func doUnRegister(c *gin.Context) {
	args := struct {
		Name  string `json:"name" binding:"required"`
		Srv   string `json:"srv" binding:"required"`
		Proto string `json:"proto" binding:"required"`
	}{}
	if err := c.ShouldBind(&args); err != nil {
		core.PanicHttpError(500, err.Error())
	}
	sorm := models.NewSysMsrv()
	sorm.MultiDelete(func(st *orm.QuerySt) string {
		st.Where("srv", args.Srv).Where("name", args.Name)
		st.Where("proto", args.Proto)
		return st.GetWheres()
	})
	if args.Proto == "http" && service.HttpPools != nil {
		service.HttpPools.Del(args.Name, args.Srv)
	} else if args.Proto == "grpc" && service.GrpcPools != nil {
		service.GrpcPools.Del(args.Name, args.Srv)
	}
	c.JSON(http.StatusOK, gin.H{"Code": 0, "Msg": "OK"})
}

// @Summary 微服务发现
// @Description 中台的微服务发现管理
// @Tags 中台管理
// @Param name formData string true "服务名称"
// @Success 200 {string} OK
// @Router /micsrv/discover [get]
func doDiscover(c *gin.Context) {
	name  := c.Param("name")
	proto := c.Param("proto")
	if len(name) < 6 || len(proto) < 4 {
		core.PanicHttpError(500, "discover server{"+name+"-->"+proto+"} error")
	}
	sorm  := models.NewSysMsrv()
	list  := sorm.GetColumn(0, -1, func(st *orm.QuerySt) string {
		st.Where("name", name).Where("proto", proto)
		st.Where("status", 1).OrderBy("stime", orm.DESC)
		return st.GetWheres()
	}, "srv")
	log.Write(log.INFO, name, "-->", proto, list)
	c.JSON(http.StatusOK, gin.H{"srvs": list, "code": 0, "msg": "OK"})
}

// @Summary 微服务重载
// @Description 微服务重载发现的服务信息
// @Tags 中台管理
// @Success 200 {string} OK
// @Router /micsrv/reload [get]
func doReload(c *gin.Context) {
	if service.HttpPools != nil {
		service.HttpPools.Load("http")
	}
	if service.GrpcPools != nil {
		service.GrpcPools.Load("grpc")
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
}

// @Summary 微服务配置
// @Description 微服务配置 加载
// @Tags 中台管理
// @Success 200 {string} OK
// @Router /micsrv/config [get]
func doConfig(c *gin.Context) {
	name := c.Param("name")
	if len(name) < 6 {
		core.PanicHttpError(500, "grpc config server{"+name+"} error")
	}
	sorm := models.NewSysYaml()
	yaml := sorm.GetValue(func(st *orm.QuerySt) string {
		st.Where("name", name).Where("status", 1)
		return st.GetWheres()
	}, "yaml")
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK", "yaml": yaml})
}
