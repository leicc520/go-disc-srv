package migrate

import (
	"database/sql"
	"github.com/leicc520/go-orm"
	"github.com/leicc520/go-orm/log"
	_ "github.com/mattn/go-sqlite3"
	"githunb.com/leicc520/go-disc-srv/app/models"
)

//初始化数据库--
func sqliteInitialize(dataSource string) {
	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		log.Write(log.ERROR, "sqlite3数据库初始化失败", err)
		panic("sqlite3数据库初始化失败")
	}

	str := `
	CREATE TABLE sys_user (
	  id      INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, --账号id
	  account VARCHAR(31) NOT NULL, --登录账号(工号)
	  loginpw VARCHAR(32) NOT NULL, --登录密码 要求客户端md5之后传到服务端做二次校验
	  email   VARCHAR(63) DEFAULT NULL, --电子邮箱
	  mobile  VARCHAR(15) DEFAULT NULL, --手机号码
	  regtime UNSIGNED INT NOT NULL, --注册时间
	  status  TINYINT DEFAULT '1', --状态 2-离职 1-在职
	  expire  UNSIGNED INT DEFAULT '0', --账号过期时间 0-永不过期
	  stime   UNSIGNED INT NOT NULL --最后操作时间
	);
	CREATE UNIQUE INDEX idx_account ON sys_user (account);
	INSERT INTO sys_user VALUES (1, 'admin', '1f9abbabf9926d579a3c5d1140421be8', 'xxx@xxx.com', '11000000000', 1646036519, 1, 1727452800, 1646036519);
	
	CREATE TABLE sys_msrv (
	  id 	  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, --记录ID
	  srv 	  VARCHAR(127) DEFAULT '',      --服务地址
	  name 	  VARCHAR(63) DEFAULT '' ,      --服务名称
	  version VARCHAR(15) DEFAULT '',       --版本号
	  proto   VARCHAR(15) DEFAULT '',       --协议http/grpc
	  status  TINYINT DEFAULT '0',          --状态 1-正常 0-失效
	  stime   UNSIGNED INT DEFAULT '0'     --更新时间
	);
	CREATE UNIQUE INDEX idx_srv ON sys_msrv (srv);
	CREATE INDEX idx_name_proto_status ON sys_msrv (name,proto,status);
	
	CREATE TABLE sys_yaml (
	  id 	  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, --记录ID
	  name 	  VARCHAR(31) DEFAULT NULL,     --配置名称
	  status  TINYINT DEFAULT '0',          --状态 1-正常 0-副本 编辑的时候保存副本
	  yaml 	  TEXT DEFAULT '',              --配置内容
	  userid  INT DEFAULT '0',              --编辑的用户ID
	  calls   INT DEFAULT '0',              --调用统计
	  version VARCHAR(255) DEFAULT '',      --版本号数据
	  stime   UNSIGNED INT DEFAULT '0'      --更新时间
	);
	CREATE INDEX idx_userid ON sys_yaml (userid);
	CREATE INDEX idx_name_status ON sys_yaml (name,status);
`
	if _, err = db.Exec(str); err != nil {
		log.Write(log.ERROR, "sqlite3数据库执行初始化失败", err)
		panic("sqlite3数据库执行SQL初始化失败")
	}

	dbConfig := orm.DbConfig{Driver: "sqlite3", Host: dataSource}
	orm.InitDBPoolSt().Set("dbmaster", &dbConfig)
	orm.InitDBPoolSt().Set("dbslaver", &dbConfig)
	orm.IsDebug = true //打印sql语句的情况
	ttlYml := models.NewSysYaml().GetTotal(nil, "COUNT(1)").ToInt64()
	ttlSrv := models.NewSysMsrv().GetTotal(nil, "COUNT(1)").ToInt64()
	ttlUsr := models.NewSysUser().GetTotal(nil, "COUNT(1)").ToInt64()
	if ttlYml < 0 || ttlSrv < 0 || ttlUsr < 1 {
		log.Write(log.ERROR, "sqlite3表初始化检测异常", err)
		panic("sqlite3表初始化检测异常")
	}
}
