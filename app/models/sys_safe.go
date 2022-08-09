package models

import (
	"github.com/leicc520/go-orm"
	"reflect"
)

type SysSafe struct {
	*orm.ModelSt
}

//结构体实例的结构说明
type SysSafeSt struct {
	Userid		uint64		`json:"userid"`		
	Sys		    int8		`json:"sys"`
	Loginpw		string		`json:"loginpw"`		
	Tocken		string		`json:"tocken"`		
	Expire		uint		`json:"expire"`		
}

//这里默认引用全局的连接池句柄
func NewSysSafe() *SysSafe {
	fields := map[string]reflect.Kind{
		"userid":		reflect.Uint64,		//角色ID
		"sys":		    reflect.Int8,		//系统别0-web 1-app
		"loginpw":		reflect.String,		//会员密码生成的Tocken
		"tocken":		reflect.String,		//随机码生成的Tocken
		"expire":		reflect.Uint,		//过期时间
	}
	
	args  := map[string]interface{}{
		"table":		"sys_safe",
		"orgtable":		"sys_safe",
		"prikey":		"sys",
		"dbmaster":		"dbmaster",
		"dbslaver":		"dbslaver",
		"slot":			0,
	}

	data := &SysSafe{&orm.ModelSt{}}
	data.Init(&orm.GdbPoolSt, args, fields)
	return data
}
