package app

import (
	"regexp"
	"strings"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/leicc520/go-disc-srv/app/models"
	"github.com/leicc520/go-disc-srv/app/service"
	"github.com/leicc520/go-gin-http"
	"github.com/leicc520/go-orm"
)

type argsSysUserListSt struct {
	service.ArgsRequestList
	Query  string `form:"query"`
	Status int    `form:"status"`
}

type argsSysUserSt struct {
	Id       int64   `json:"id" form:"id"`
	Account  string  `json:"account"   form:"account"  binding:"required,min=3"`
	NickName string  `json:"nickname"  form:"nickname" binding:"required,min=1"`
	Email    string  `json:"email"     form:"email"    binding:"required,email"`
	Mobile   string  `json:"mobile"    form:"mobile"   binding:"required,regex=^1[3456789][\d]{9}$"`
	Status   int8    `json:"status"    form:"status"`
	Expire   string  `json:"expire"    form:"expire"   binding:"required,regex=^[\d]{4}\-[\d]{2}\-[\d]{2}$"`
	Loginpw  string  `json:"loginpw"   form:"loginpw"  binding:"omitempty"`
}

// @Summary 系统账号管理
// @Description 中台的系统账号管理
// @Tags 中台管理
// @Param SIGNATURE header string true "对接第三方登录的时候获取的token"
// @Param query formData string true "查询关键词-账号名称"
// @Param status formData int true "通过状态检索用户"
// @Param stime formData array true "开始/结束时间Y-m-d"
// @Param sorts formData object true "数据排序业务"
// @Param start formData int true "获取数据偏移量"
// @Param limit formData int true "获取列表记录数量"
// @Success 200 {json} HttpView {sid:xxx[要记录,以后请求放http请求头中],user:{用户基础信息}}
// @Router /api/core/user/list [post]
func sysUserList(c *gin.Context) {
	args := argsSysUserListSt{}
	if err := c.ShouldBind(&args); err != nil {
		core.PanicValidateHttpError(1001, err)
	}
	sorm    := models.NewSysUser()
	wHandle := sysUserListWhere(&args)
	total   := sorm.GetTotal(wHandle, "COUNT(1)").ToInt64()
	var list []orm.SqlMap = nil
	if total > 0 { //数据大于0 的情况
		sorm.Format(func(sm orm.SqlMap) {
			delete(sm, "loginpw")
		})
		list = sorm.GetList(args.Start, args.Limit, wHandle, "*")
	}
	core.NewHttpView(c).JsonDisplay(gin.H{"total": total, "list": list})
}

//获取查询条件的设定
func sysUserListWhere(args *argsSysUserListSt) orm.WHandler {
	return func(st *orm.QuerySt) string {
		if args.Status > 0 {
			st.Where("status", args.Status)
		}
		if len(args.Stime) > 0 {
			dtime := orm.DT2UnixTimeStamp(args.Stime[0], "2006-01-02")
			if dtime > 1 { //数据大于0的情况
				st.Where("stime", dtime, orm.OP_GE)
			}
		}
		if len(args.Stime) > 1 {
			dtime := orm.DT2UnixTimeStamp(args.Stime[1], "2006-01-02")
			if dtime > 1 { //数据大于0的情况
				st.Where("stime", dtime, orm.OP_LE)
			}
		}
		if strings.TrimSpace(args.Query) != "" {
			st.Where("(", "")
			st.Where("account", "%"+args.Query+"%", orm.OP_LIKE)
			st.Where("nickname", "%"+args.Query+"%", orm.OP_LIKE, orm.OP_OR)
			st.Where(")", "")
		}
		//设置排序处理逻辑
		if args.Sorts != nil && len(args.Sorts) > 0 {
			for field, orderBy := range args.Sorts {
				st.OrderBy(field, orderBy)
			}
		}
		return st.GetWheres()
	}
}

// @Summary 系统账号管理
// @Description 中台的系统账号管理
// @Tags 中台管理
// @Param SIGNATURE header string true "对接第三方登录的时候获取的token"
// @Param email formData string true "编辑邮箱信息"
// @Param mobile formData string true "编辑手机号码信息"
// @Param loginpw formData string true "编辑密码信息"
// @Param avatar formData string true "编辑个人头像"
// @Success 200 {json} HttpView {sid:xxx[要记录,以后请求放http请求头中],user:{用户基础信息}}
// @Router /api/core/user/safe [post]
func sysUserSafe(c *gin.Context) {
	args := struct {
		Email   string `form:"email" json:"email" binding:"required,email"`
		Mobile  string `form:"mobile" json:"mobile" binding:"required,regex=^1[3456789][\d]{9}$"`
		Loginpw string `form:"loginpw" json:"loginpw" binding:"omitempty"`
	}{}
	if err := c.ShouldBind(&args); err != nil {
		core.PanicValidateHttpError(1001, err)
	}
	acl    := core.NewAcl(0, models.NewSysSafe().Instance())
	data   := orm.SqlMap{"email":args.Email, "mobile":args.Mobile, "loginpw":args.Loginpw}
	userId := service.JWTACLGetUserid(c)
	args.Loginpw = strings.TrimSpace(args.Loginpw)
	if len(args.Loginpw) > 0 { //验证码密码的处理逻辑
		if ok, _ := regexp.MatchString("^[^\\s\\r\\n]{6,}$", args.Loginpw); !ok {
			core.PanicHttpError(1001, "密码必须6个以上字符组成,不能含空字符.")
		}
		args.Loginpw = acl.Crypt(args.Loginpw)
		data["loginpw"] = args.Loginpw
		acl.SetToken(userId, args.Loginpw, "", 0)
	} else { //删除密码的需改
		delete(data, "loginpw")
	}
	models.NewSysUser().Save(userId, data)
	core.NewHttpView(c).JsonDisplay(nil)
}

// @Summary 系统账号管理
// @Description 中台的系统账号管理
// @Tags 中台管理
// @Param SIGNATURE header string true "对接第三方登录的时候获取的token"
// @Param id formData int true "选传大于1说明编辑否则新增"
// @Param account formData string true "登录账号 唯一"
// @Param nickname formData string true "用户昵称信息"
// @Param email formData string true "联系邮箱"
// @Param mobile formData string true "	联系手机号码"
// @Param status formData int8 true "审核状态"
// @Param expire formData string true "过期时间"
// @Param loginpw formData string true "设置的话代表重置密码"
// @Success 200 {json} HttpView {sid:xxx[要记录,以后请求放http请求头中],user:{用户基础信息}}
// @Router /api/core/user/update [post]
func sysUserUpdate(c *gin.Context) {
	args   := argsSysUserSt{}
	if err := c.ShouldBind(&args); err != nil {
		core.PanicValidateHttpError(1001, err)
	}
	expire := orm.DT2UnixTimeStamp(args.Expire, orm.DATEBASICFormat)
	args.Loginpw = strings.TrimSpace(args.Loginpw)
	data   := orm.SqlMap{"account":args.Account, "nickname":args.NickName, "email":args.Email,
		"mobile":args.Mobile, "status":args.Status, "expire":expire, "loginpw":args.Loginpw, "stime":time.Now().Unix()}
	acl := core.NewAcl(0, models.NewSysSafe().Instance())
	if len(args.Loginpw) > 0 { //验证码密码的处理逻辑
		if ok, _ := regexp.MatchString("^[^\\s\\r\\n]{6,}$", args.Loginpw); !ok {
			core.PanicHttpError(1001, "密码必须6个以上字符组成,不能含空字符.")
		}
		data["loginpw"] = acl.Crypt(args.Loginpw)
	} else { //删除密码的需改
		delete(data, "loginpw")
	}
	sorm  := models.NewSysUser()
	oldId := sorm.IsExists(func(st *orm.QuerySt) string {
		st.Where("account", args.Account)
		return st.GetWheres()
	}).ToInt64()
	if oldId > 0 && oldId != args.Id {
		core.PanicHttpError(1001, "请求账号已存在,请无重复.")
	}
	if args.Id > 0 { //更新记录
		sorm.Save(args.Id, data)
	} else { //新增记录
		data["regtime"] = time.Now().Unix()
		args.Id = sorm.NewOne(data, nil)
	}
	core.NewHttpView(c).JsonDisplay(nil)
}
