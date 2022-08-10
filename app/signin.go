package app

import (
	"strconv"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/leicc520/go-gin-http"
	"github.com/leicc520/go-orm"
	"github.com/leicc520/go-disc-srv/app/models"
	"github.com/leicc520/go-disc-srv/app/service"
)

//获取用户的基础信息
func _user_(user *models.SysUserSt) orm.SqlMap {
	data := orm.SqlMap{"id":user.Id, "account":user.Account, "nick":user.NickName,
		"email":user.Email, "mobile":user.Mobile, "expire":user.Expire}
	return data
}

// @Summary 中台的登录
// @Description 中台的登录管理平台
// @Tags 中台管理
// @Param account formData string true "登录账号"
// @Param loginpw formData string true "账号密码"
// @Param xtoken formData string true "验证码的验证token"
// @Success 200 {json} HttpView {sid:xxx[要记录,以后请求放http请求头中],user:{用户基础信息}}
// @Router /signin/check [post]
func signInCheck(c *gin.Context) {
	args := struct {
		Account string `form:"account" binding:"required,min=3"`
		Loginpw string `form:"loginpw" binding:"required,min=6"`
		Xtoken string `form:"xtoken" binding:"required,min=8"`
	}{}
	if err := c.ShouldBind(&args); err != nil {
		core.PanicValidateHttpError(1001, err)
	}
	if !core.Gcaptcha.CheckHash(c, args.Xtoken) {
		core.PanicHttpError(1010)
	}
	acl   := core.NewAcl(0, models.NewSysSafe().Instance())
	sTime := time.Now().Unix()
	fields:= "id,status,loginpw,expire,account,nickname,email,mobile"
	user  := models.SysUserSt{}
	err   := models.NewSysUser().GetItem(func(st *orm.QuerySt) string {
		st.Where("account", args.Account)
		return st.GetWheres()
	}, fields).ToStruct(&user)
	if err != nil || user.Id < 1 {
		core.PanicHttpError(1001, "请求的账号密码错误.")
	}
	if user.Status != 1 {
		core.PanicHttpError(1011)
	} else if user.Loginpw != acl.Crypt(args.Loginpw) {
		core.PanicHttpError(1012)
	} else if user.Expire > 0 && user.Expire < sTime {
		core.PanicHttpError(1013)
	}
	xToken := strconv.FormatInt(time.Now().Unix(), 10)
	xToken  = acl.SetToken(user.Id, user.Loginpw, xToken, 0)
	jToken := core.JwtToken(user.Id, c.Request.UserAgent(), xToken)
	data   := _user_(&user) //数据资料信息
	core.NewHttpView(c).JsonDisplay(gin.H{"sign":jToken, "user":data})
}

// @Summary 中台登录-获取用户信息
// @Description 中台的登录管理平台
// @Tags 中台管理
// @Param SIGNATURE header string true "对接第三方登录的时候获取的token"
// @Success 200 {json} HttpView
// @Router /api/core/signin/recheck [post]
func signInReCheck(c *gin.Context) {
	userId := service.JWTACLGetUserid(c)
	sorm   := models.NewSysUser()
	user   := models.SysUserSt{}
	if err := sorm.GetOne(userId).ToStruct(&user); err != nil {
		core.PanicHttpError(1001)
	}
	data   := _user_(&user) //数据资料信息
	core.NewHttpView(c).JsonDisplay(gin.H{"user":data})
}
