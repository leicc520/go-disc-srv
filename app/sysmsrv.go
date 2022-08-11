package app

import (
	"github.com/leicc520/go-gin-http/micro"
	"strings"
	
	"github.com/gin-gonic/gin"
	"github.com/leicc520/go-disc-srv/app/models"
	"github.com/leicc520/go-disc-srv/app/service"
	"github.com/leicc520/go-gin-http"
	"github.com/leicc520/go-orm"
)

type argsSysMsrvListSt struct {
	service.ArgsRequestList
	Query string `form:"query"`
}

// @Summary 系统微服务管理
// @Description 中台的系统微服务管理
// @Tags 中台管理
// @Param SIGNATURE header string true "对接第三方登录的时候获取的token"
// @Param query formData string true "查询关键词,微服务名称"
// @Param stime formData array true "开始/结束时间Y-m-d"
// @Param sorts formData object true "数据排序业务"
// @Param start formData int true "获取数据偏移量"
// @Param limit formData int true "获取列表记录数量"
// @Success 200 {json} HttpView {sid:xxx[要记录,以后请求放http请求头中],user:{用户基础信息}}
// @Router /api/core/misrv/list [post]
func sysmsrvList(c *gin.Context) {
	args   := argsSysMsrvListSt{}
	if err := c.ShouldBind(&args); err != nil {
		core.PanicValidateHttpError(1001, err)
	}
	sorm    := models.NewSysMsrv()
	wHandle := sysMSrvListWhere(&args)
	var list interface{} = nil
	total := sorm.GetTotal(wHandle, "COUNT(1)").ToInt64()
	if total > 0 { //数据大于0 的情况
		list = sorm.GetList(args.Start, args.Limit, wHandle, "*")
	}
	core.NewHttpView(c).JsonDisplay(gin.H{"total": total, "list": list})
}

// @Summary 系统微服务状态设置
// @Description 中台的服务状态设置
// @Tags 中台管理
// @Param SIGNATURE header string true "对接第三方登录的时候获取的token"
// @Param id formData int true "服务ID记录"
// @Param status formData int true "服务状态1-上线 2-离线"
// @Success 200 {json} HttpView {sid:xxx[要记录,以后请求放http请求头中]}
// @Router /api/core/misrv/status [post]
func sysmsrvStatus(c *gin.Context) {
	args   := struct {
		Id     int64 `form:"id" binding:"required,min=1"`
		Status int8  `form:"status" binding:"required,min=0"`
	}{}
	if err := c.ShouldBind(&args); err != nil {
		core.PanicValidateHttpError(1001, err)
	}
	sorm   := models.NewSysMsrv()
	data   := models.SysMsrvSt{}
	if err := sorm.GetOne(args.Id).ToStruct(&data); err != nil || data.Id < 1 {
		core.PanicHttpError(1002, "请求的微服务不存在哟.")
	}
	regSrv := micro.NewRegSrvClient("")
	if args.Status == 1 { //上线注册服务  更新数据让业务代码自己做
		if !regSrv.Health(1, data.Proto, data.Srv) {
			core.PanicHttpError(1004, "该服务心跳检查异常,无法上架.")
		}
		regSrv.Register(data.Name, data.Srv, data.Proto, data.Version)
	} else { //注销离线的处理逻辑
		nSize := sorm.GetTotal(func(st *orm.QuerySt) string {
			st.Where("name", data.Name).Where("proto", data.Proto)
			st.Where("id", data.Id, orm.OP_NE)
			return st.GetWheres()
		}, "COUNT(1)").ToInt64()
		if nSize < 1 {//没有其他的可用服务的情况
			core.PanicHttpError(1003, "无其他可用服务,不允许下架啦.")
		}
		regSrv.UnRegister(data.Proto, data.Name, data.Srv)
	}
	core.NewHttpView(c).JsonDisplay(nil)
}

// @Summary 系统微服务状态重置
// @Description 中台的服务状态重置
// @Tags 中台管理
// @Param SIGNATURE header string true "对接第三方登录的时候获取的token"
// @Success 200 {json} HttpView {sid:xxx[要记录,以后请求放http请求头中]}
// @Router /api/core/misrv/reload [post]
func sysmsrvReload(c *gin.Context) {
	regSrv := micro.NewRegSrvClient("")
	if err := regSrv.Reload(); err != nil {
		core.PanicHttpError(500, err.Error())
	}
	core.NewHttpView(c).JsonDisplay(nil)
}

//获取查询条件的设定
func sysMSrvListWhere(args *argsSysMsrvListSt) orm.WHandler {
	return func(st *orm.QuerySt) string {
		if len(args.Stime) > 0 {
			dTime := orm.DT2UnixTimeStamp(args.Stime[0], "2006-01-02")
			if dTime > 1 { //数据大于0的情况
				st.Where("stime", dTime, orm.OP_GE)
			}
		}
		if len(args.Stime) > 1 {
			dTime := orm.DT2UnixTimeStamp(args.Stime[1], "2006-01-02")
			if dTime > 1 { //数据大于0的情况
				st.Where("stime", dTime, orm.OP_LE)
			}
		}
		if strings.TrimSpace(args.Query) != "" {
			st.Where("(", "")
			st.Where("name", args.Query, orm.OP_EQ)
			st.Where("srv", args.Query, orm.OP_EQ, orm.OP_OR)
			st.Where(")", "")
		}
		if args.Sorts != nil && len(args.Sorts) > 0 {
			for field, orderBy := range args.Sorts {
				st.OrderBy(field, orderBy)
			}
		}
		return st.GetWheres()
	}
}
