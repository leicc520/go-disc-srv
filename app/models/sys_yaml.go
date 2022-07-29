package models

import (
	"github.com/leicc520/go-orm"
	"reflect"
)

type SysYaml struct {
	*orm.ModelSt
}

//结构体实例的结构说明
type SysYamlSt struct {
	Id		uint	`json:"id"`
	Name	string	`json:"name"`
	Status	int8	`json:"status"`
	Yaml	string	`json:"yaml"`
	Userid	int		`json:"userid"`
	Calls	int		`json:"calls"`
	Version	string	`json:"version"`
	Stime	int		`json:"stime"`
}

//这里默认引用全局的连接池句柄
func NewSysYaml() *SysYaml {
	fields := map[string]reflect.Kind{
		"id":			reflect.Uint,		//记录ID
		"name":			reflect.String,		//配置名称
		"status":		reflect.Int8,		//状态 1-正常 0-副本 编辑的时候保存副本
		"yaml":			reflect.String,		//配置内容
		"userid":		reflect.Int,		//编辑的用户ID
		"calls":		reflect.Int,		//调用统计
		"version":		reflect.String,		//版本号数据
		"stime":		reflect.Int,		//更新时间
	}
	
	args  := map[string]interface{}{
		"table":		"sys_yaml",
		"orgtable":		"sys_yaml",
		"prikey":		"id",
		"dbmaster":		"dbmaster",
		"dbslaver":		"dbslaver",
		"slot":			0,
	}

	data := &SysYaml{&orm.ModelSt{}}
	data.Init(&orm.GdbPoolSt, args, fields)
	return data
}
