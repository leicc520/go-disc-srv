package models

import (
	"github.com/leicc520/go-orm"
	"reflect"
)

type SysUser struct {
	*orm.ModelSt
}

//结构体实例的结构说明
type SysUserSt struct {
	Id		uint   	`json:"id"`
	Account	string	`json:"account"`
	Loginpw	string	`json:"loginpw"`
	Email	string	`json:"email"`
	Mobile	string	`json:"mobile"`
	Regtime	int		`json:"regtime"`
	Status	int8	`json:"status"`
	Expire	int		`json:"expire"`
	Isdup	int8	`json:"isdup"`
	Stime	int		`json:"stime"`
}

//这里默认引用全局的连接池句柄
func NewSysUser() *SysUser {
	fields := map[string]reflect.Kind{
		"id":			reflect.Uint,		//账号id
		"account":		reflect.String,		//登录账号
		"loginpw":		reflect.String,		//登录密码 要求客户端md5之后传到服务端做二次校验
		"email":		reflect.String,		//电子邮箱
		"mobile":		reflect.String,		//手机号码
		"regtime":		reflect.Int,		//注册时间
		"status":		reflect.Int8,		//状态1-正常 2-冻结
		"expire":		reflect.Int,		//账号过期时间 0-永不过期
		"isdup":		reflect.Int8,		//是否允许多终端登录 1-允许 2-不允许
		"stime":		reflect.Int,		//最后操作时间
	}
	
	args  := map[string]interface{}{
		"table":		"sys_user",
		"orgtable":		"sys_user",
		"prikey":		"id",
		"dbmaster":		"dbmaster",
		"dbslaver":		"dbslaver",
		"slot":			0,
	}

	data := &SysUser{&orm.ModelSt{}}
	data.Init(&orm.GdbPoolSt, args, fields)
	return data
}
