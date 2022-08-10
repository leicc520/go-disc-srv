package app

import (
	"strings"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/leicc520/go-gin-http"
	"github.com/leicc520/go-orm"
	"github.com/leicc520/go-disc-srv/app/models"
	"github.com/leicc520/go-disc-srv/app/service"
)

type argsYamlSt struct {
	Id      int64  `form:"id"`
	Name    string `form:"name" binding:"required,min=3,max=31"`
	Version string `form:"version" binding:"required,min=1"`
	Yaml    string `form:"yaml" binding:"required,min=1"`
	Status  int8   `form:"status"`
	UserId  int64  `form:"userid"`
	Stime   int64  `form:"stime"`
}

type yamlQuerySt struct {
	service.ArgsRequestList
	UserId int64  `form:"userid"`
	Query  string `form:"query"`
}

// @Summary 系统服务配置
// @Description 中台的服务配置管理
// @Tags 中台管理
// @Param SIGNATURE header string true "对接第三方登录的时候获取的token"
// @Param userid formData int true "归属用户"
// @Param status formData int true "状态信息"
// @Param query formData string true "关键词内容检索"
// @Param start formData int true "获取数据偏移量"
// @Param limit formData int true "获取列表记录数量"
// @Success 200 {json} HttpView {sid:xxx[要记录,以后请求放http请求头中],user:{用户基础信息}}
// @Router /api/core/yaml/list [post]
func sysYamlList(c *gin.Context) {
	args   := yamlQuerySt{}
	if err := c.ShouldBind(&args); err != nil {
		core.PanicValidateHttpError(1001, err)
	}
	sorm   := models.NewSysYaml()
	var list []orm.SqlMap = nil
	wHandle:= sysYamlListWhere(&args)
	total  := sorm.GetTotal(wHandle, "COUNT(1)").ToInt64()
	if total > 0 { //获取列表数据
		sorm.Format(func(sm orm.SqlMap) {
			sm["user"] = service.GetSysOperator(sm["userid"])
		})
		list = sorm.GetList(args.Start, args.Limit, wHandle, "*")
	}
	core.NewHttpView(c).JsonDisplay(gin.H{"total": total, "list": list})
}

//设置查询条件的处理逻辑
func sysYamlListWhere(args *yamlQuerySt) orm.WHandler {
	return func(st *orm.QuerySt) string {
		if args.UserId > 0 {
			st.Where("userid", args.UserId)
		}
		args.Query = strings.TrimSpace(args.Query)
		if len(args.Query) > 0 {
			st.Where("(", "")
			st.Where("name", "%"+args.Query+"%", orm.OP_LIKE)
			st.Where("yaml", "%"+args.Query+"%", orm.OP_LIKE, orm.OP_OR)
			st.Where(")", "")
		} else { //默认取在线的记录
			st.Where("status", 1)
		}
		if args.Sorts != nil && len(args.Sorts) > 0 {
			for field, orderBy := range args.Sorts {
				st.OrderBy(field, orderBy)
			}
		}
		return st.GetWheres()
	}
}

// @Summary 服务配置删除
// @Description 中台的服务配置管理
// @Tags 中台管理
// @Param SIGNATURE header string true "对接第三方登录的时候获取的token"
// @Param id formData int true "角色ID记录"
// @Success 200 {json} HttpView {sid:xxx[要记录,以后请求放http请求头中],user:{用户基础信息}}
// @Router /api/core/yaml/delete [post]
func sysYamlDelete(c *gin.Context) {
	args := struct {
		Id int64 `form:"id" binding:"required,min=1"`
	}{}
	if err := c.ShouldBind(&args); err != nil {
		core.PanicValidateHttpError(1001, err)
	}
	sorm  := models.NewSysYaml()
	data  := models.SysYamlSt{}
	if err := sorm.GetOne(args.Id).ToStruct(&data); err != nil || data.Status == 1 {
		core.PanicHttpError(1002, "配置已经投入使用,无法删除.")
	}
	sorm.Delete(args.Id) //删除缓存
	core.NewHttpView(c).JsonDisplay(nil)
}

// @Summary 服务配置编辑
// @Description 中台的服务配置管理
// @Tags 中台管理
// @Param SIGNATURE header string true "对接第三方登录的时候获取的token"
// @Param id formData int true "不传或者0则新增，否则编辑"
// @Param name formData string true "配置名称"
// @Param isdev formData int8 true "配置环境"
// @Param version formData string true "配置版本号"
// @Param Yaml formData string true "配置值信息yaml"
// @Success 200 {json} HttpView {sid:xxx[要记录,以后请求放http请求头中],user:{用户基础信息}}
// @Router /api/core/yaml/update [post]
func sysYamlUpdate(c *gin.Context) {
	args   := argsYamlSt{}
	if err := c.ShouldBind(&args); err != nil {
		core.PanicValidateHttpError(1001, err)
	}
	sorm  := models.NewSysYaml()
	oldId := sorm.IsExists(func(st *orm.QuerySt) string {
		st.Where("name", args.Name)
		st.Where("status", 1)
		return st.GetWheres()
	}).ToInt64()
	if oldId > 0 && oldId != args.Id {
		core.PanicHttpError(1016)
	}
	//封装成事务，确保程序的一致性
	result := sorm.Query().Transaction(func(st *orm.QuerySt) bool {
		if args.Id > 0 { //将原来的迁移到无效的归档记录当中
			datas := sorm.GetOne(args.Id)
			if datas != nil && len(datas) > 0 {
				datas  = datas.Merge(orm.SqlMap{"status":2})
				xldid := sorm.Save(args.Id, datas)
				if xldid < 1 {
					return false
				}
			}
		}
		args.UserId = service.JWTACLGetUserid(c)
		datas := orm.SqlMap{"name":args.Name, "version":args.Version,
			"yaml":args.Yaml, "status":1, "userid":args.UserId, "stime":time.Now().Unix()}
		xldid := sorm.NewOne(datas, nil) //更新数据信息
		if xldid < 1 {
			return false
		}
		return true
	})
	if !result {//配置异常的情况
		core.PanicHttpError(1017, "请求更新配置异常")
	}
	core.NewHttpView(c).JsonDisplay(nil)
}
