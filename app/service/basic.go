package service

import (
	"fmt"
	
	"github.com/gin-gonic/gin"
	"github.com/leicc520/go-gin-http"
	"github.com/leicc520/go-disc-srv/app/models"
)

//普通的列表检索业务逻辑基础
type ArgsBasicList struct {
	Sorts map[string]string `form:"sorts"`
	Start int64             `form:"start" binding:"min=0"`
	Limit int64             `form:"limit" binding:"required,max=500"`
}

type ArgsRequestList struct {
	ArgsBasicList
	Stime []string `form:"stime"`
}

//获取当前登录用户的信息
func JWTACLGetUserid(c *gin.Context) int64 {
	if aclUser, isExists := c.Get("user"); isExists {
		if signUser, ok := aclUser.(*core.JwtUser); ok {
			return signUser.Id
		}
	}
	return core.JWTACLUserid(c)
}

//获取指定用户ID
func GetSysOperator(id interface{}) string {
	sorm := models.NewSysUser()
	user := struct {
		NickName string `json:"nickname"`
		Account  string `json:"account"`
	}{}
	if err := sorm.GetOne(id).ToStruct(&user); err != nil {
		return "匿名[" + fmt.Sprintf("%v", id) + "]"
	}
	return user.NickName + "[" + user.Account + "]"
}
