package discCTRL

import (
	"github.com/leicc520/go-orm"
	"net"
	"net/http"
	"regexp"
	"time"
	
	"git.cht-group.net/leicc/go-disc-srv/app/logic"
	"git.cht-group.net/leicc/go-disc-srv/app/models/sys"
	"git.cht-group.net/leicc/go-orm"
	"git.cht-group.net/leicc/go-orm/log"
	"github.com/gin-gonic/gin"
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
	args := micsrvNodeSt{}
	if err := c.ShouldBind(&args); err != nil {
		logic.PanicHttpError(500, err.Error())
	}
	//拼接真实的服务地址 只提交端口的情况 要取一下请求IP
	if ok, err := regexp.MatchString("^[\\d]+$", args.Srv); ok && err == nil {
		srvip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
		args.Srv = srvip + ":" + args.Srv
	}
	sorm := sys.NewSysMsrv()
	oldid := sorm.NewOneFromHandler(func(st *orm.QuerySt) *orm.QuerySt {
		st.Value("srv", args.Srv).Value("name", args.Name)
		st.Value("version", args.Version).Value("proto", args.Proto)
		st.Value("status", 1).Value("stime", time.Now().Unix())
		return st
	}, func(st *orm.QuerySt) interface{} {
		st.Duplicate("name", args.Name).Duplicate("version", args.Version).Duplicate("proto", args.Proto)
		st.Duplicate("status", 1).Duplicate("stime", time.Now().Unix())
		return nil
	})
	log.Write(log.INFO, oldid, "---->", args)
	//分别添加到对应的数据结构当中
	if args.Proto == "http" && gHttpPools != nil {
		gHttpPools.Put(args.Name, args.Srv, oldid)
	} else if args.Proto == "grpc" && gGrpcPools != nil {
		gGrpcPools.Put(args.Name, args.Srv, oldid)
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
		logic.PanicHttpError(500, err.Error())
	}
	sorm := sys.NewSysMsrv()
	sorm.MultiDelete(func(st *orm.QuerySt) string {
		st.Where("srv", args.Srv).Where("name", args.Name).Where("proto", args.Proto)
		return st.GetWheres()
	})
	if args.Proto == "http" && gHttpPools != nil {
		gHttpPools.Del(args.Name, args.Srv)
	} else if args.Proto == "grpc" && gGrpcPools != nil {
		gGrpcPools.Del(args.Name, args.Srv)
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
	name := c.Param("name")
	proto := c.Param("proto")
	if len(name) < 6 || len(proto) < 4 {
		logic.PanicHttpError(500, "discover server{"+name+"-->"+proto+"} error")
	}
	var list []string = nil
	if proto == "http" && gHttpPools != nil {
		list = getSrv(name, proto, gHttpPools)
	} else if proto == "grpc" && gGrpcPools != nil {
		list = getSrv(name, proto, gGrpcPools)
	}
	log.Write(log.INFO, name, "-->", proto, list)
	c.JSON(http.StatusOK, gin.H{"srvs": list, "code": 0, "msg": "OK"})
}

//获取服务数据资料信息
func getSrv(name, proto string, gsrv *MicSrvPoolSt) []string {
	list := gsrv.Get(name)
	if list == nil || len(list) == 0 { //没找到通过db查找
		sorm := sys.NewSysMsrv()
		datas := sorm.GetList(0, -1, func(st *orm.QuerySt) string {
			st.Where("name", name).Where("proto", proto)
			st.OrderBy("status", orm.ASC).OrderBy("stime", orm.DESC)
			return st.GetWheres()
		}, "id,name,srv,proto,status,version")
		node := micsrvNodeSt{}
		for _, msrv := range datas {
			if err := msrv.ToStruct(&node); err != nil || node.Id < 1 {
				log.Write(log.ERROR, err)
				continue
			}
			//心跳正常的情况加入 数据
			if regSrv.Health(1, node.Proto, node.Srv) {
				list = append(list, node.Srv)
				gsrv.Put(node.Name, node.Srv, node.Id)
				if node.Status != 1 {//更新数据记录状态
					sorm.Save(node.Id, orm.SqlMap{"status":1, "stime":time.Now().Unix()})
				}
			}
		}
	}
	return list
}

// @Summary 微服务重载
// @Description 微服务重载发现的服务信息
// @Tags 中台管理
// @Success 200 {string} OK
// @Router /micsrv/reload [get]
func doReload(c *gin.Context) {
	if gHttpPools != nil {
		gHttpPools.Load("http")
	}
	if gGrpcPools != nil {
		gGrpcPools.Load("grpc")
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
		logic.PanicHttpError(500, "grpc config server{"+name+"} error")
	}
	sorm := sys.NewSysYaml()
	yaml := sorm.GetValue(func(st *orm.QuerySt) string {
		st.Where("name", name).Where("status", 1)
		return st.GetWheres()
	}, "yaml")
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK", "yaml": yaml})
}
