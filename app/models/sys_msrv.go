package models

import (
	"github.com/leicc520/go-orm"
	"reflect"
)

type SysMsrv struct {
	*orm.ModelSt
}

//结构体实例的结构说明
type SysMsrvSt struct {
	Id			uint		`json:"id"`
	Srv			string		`json:"srv"`
	Name		string		`json:"name"`		
	Version		string		`json:"version"`		
	Proto		string		`json:"proto"`		
	Status		int8		`json:"status"`		
	Stime		int			`json:"stime"`
}

//这里默认引用全局的连接池句柄
func NewSysMsrv() *SysMsrv {
	fields := map[string]reflect.Kind{
		"id":			reflect.Uint,		//记录ID
		"srv":			reflect.String,		//服务地址
		"name":			reflect.String,		//服务名称
		"version":		reflect.String,		//版本号
		"proto":		reflect.String,		//协议
		"status":		reflect.Int8,		//状态 1-正常 0-失效
		"stime":		reflect.Int,		//更新时间
	}
	
	args  := map[string]interface{}{
		"table":		"sys_msrv",
		"orgtable":		"sys_msrv",
		"prikey":		"id",
		"dbmaster":		"dbmaster",
		"dbslaver":		"dbslaver",
		"slot":			0,
	}

	data := &SysMsrv{&orm.ModelSt{}}
	data.Init(&orm.GdbPoolSt, args, fields)
	return data
}
