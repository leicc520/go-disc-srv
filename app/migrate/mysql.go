package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	
	_ "github.com/go-sql-driver/mysql"
	"github.com/leicc520/go-disc-srv/app/models"
	"github.com/leicc520/go-orm"
	"github.com/leicc520/go-orm/log"
)

//检测是否已经存在表信息了
func mysqlCheckExists(dataSource string) error {
	dbConfig := orm.DbConfig{Driver: "mysql", Host: dataSource}
	orm.InitDBPoolSt().Set("dbmaster", &dbConfig)
	orm.InitDBPoolSt().Set("dbslaver", &dbConfig)
	orm.IsDebug = true //打印sql语句的情况
	ttlYml := models.NewSysYaml().NoCache().GetTotal(nil, "COUNT(1)").ToInt64()
	ttlSrv := models.NewSysMsrv().NoCache().GetTotal(nil, "COUNT(1)").ToInt64()
	ttlUsr := models.NewSysUser().NoCache().GetTotal(nil, "COUNT(1)").ToInt64()
	ttlSafe:= models.NewSysSafe().NoCache().GetTotal(nil, "COUNT(1)").ToInt64()
	if ttlYml < 0 || ttlSrv < 0 || ttlUsr < 1 || ttlSafe < 0 {
		log.Write(log.ERROR, "mysql表初始化检测异常")
		return errors.New("mysql表初始化检测异常")
	}
	return nil
}

//初始化数据库
func mysqlInitialize(dataSource string) {
	fmt.Println("===============================开始数据库初始化===============================")
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Write(log.ERROR, "mysql数据库初始化失败", err)
		panic("mysql数据库初始化失败")
	}
	defer db.Close() //关闭连接
	astr := []string{
		"CREATE TABLE IF NOT EXISTS `sys_user` (  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '账号id',  `account` varchar(31) NOT NULL COMMENT '登录账号',  `nickname` varchar(63) NOT NULL COMMENT '昵称',  `loginpw` varchar(32) NOT NULL COMMENT '登录密码 要求客户端md5之后传到服务端做二次校验',  `email` varchar(63) NOT NULL COMMENT '电子邮箱',  `mobile` varchar(15) NOT NULL COMMENT '手机号码',  `regtime` int(11) NOT NULL COMMENT '注册时间',  `status` tinyint(1) NOT NULL COMMENT '状态1-正常 2-冻结',  `expire` int(11) NOT NULL COMMENT '账号过期时间 0-永不过期',  `stime` int(11) NOT NULL COMMENT '最后操作时间',  PRIMARY KEY (`id`),  UNIQUE KEY `idx_account` (`account`) USING BTREE) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;",
		"REPLACE INTO `sys_user` VALUES (1, 'admin', '超级管理员', '1f9abbabf9926d579a3c5d1140421be8', 'xxx@xxx.com', '11000000000', 1646036519, 1, 1727452800, 1646036519);",
		"CREATE TABLE IF NOT EXISTS `sys_msrv` (  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '记录ID',  `srv` varchar(127) DEFAULT '' COMMENT '服务地址',  `name` varchar(63) DEFAULT '' COMMENT '服务名称',  `version` varchar(15) DEFAULT '' COMMENT '版本号',  `proto` varchar(15) DEFAULT '' COMMENT '协议',  `status` tinyint(4) DEFAULT 0 COMMENT '状态 1-正常 0-失效',  `addtime` int(11) DEFAULT 0 COMMENT '记录时间',  `stime` int(11) DEFAULT 0 COMMENT '更新时间',  PRIMARY KEY (`id`),  UNIQUE KEY `idx_srv` (`srv`),  KEY `idx_name_proto_status` (`name`,`proto`,`status`) USING BTREE) ENGINE=InnoDB AUTO_INCREMENT=137 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;",
		"CREATE TABLE IF NOT EXISTS `sys_yaml` (  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '记录ID',  `name` varchar(31) DEFAULT NULL COMMENT '配置名称',  `status` tinyint(4) DEFAULT 0 COMMENT '状态 1-正常 0-副本 编辑的时候保存副本',  `yaml` text DEFAULT NULL COMMENT '配置内容',  `userid` int(11) DEFAULT 0 COMMENT '编辑的用户ID',  `calls` int(11) DEFAULT 0 COMMENT '调用统计',  `version` varchar(255) DEFAULT '' COMMENT '版本号数据',  `stime` int(11) DEFAULT 0 COMMENT '更新时间',  PRIMARY KEY (`id`),  KEY `idx_userid` (`userid`),  KEY `idx_name_status` (`name`,`status`) USING BTREE) ENGINE=InnoDB AUTO_INCREMENT=103 DEFAULT CHARSET=utf8mb4;",
		"CREATE TABLE IF NOT EXISTS `sys_safe` (  `userid` bigint(20) unsigned NOT NULL COMMENT '角色ID',  `sys` tinyint(1) NOT NULL COMMENT '系统别0-web 1-app',  `loginpw` varchar(63) DEFAULT NULL COMMENT '会员密码生成的Tocken',  `tocken` varchar(32) DEFAULT NULL COMMENT '随机码生成的Tocken',  `expire` int(10) unsigned DEFAULT NULL COMMENT '过期时间',  PRIMARY KEY (`userid`,`sys`)) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='会员的安全Tocken';",
	}
	for _, sqlStr := range astr {
		if _, err = db.Exec(sqlStr); err != nil {
			log.Write(log.ERROR, "mysql数据库执行初始化失败", err)
			panic("mysql数据库执行SQL初始化失败:"+sqlStr)
		}
	}
	if err = mysqlCheckExists(dataSource); err != nil {
		panic(err)
	}
	fmt.Println("=============================初始化"+dataSource+"完成=============================")
}
