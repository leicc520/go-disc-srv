package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/leicc520/go-disc-srv/app/service"
)

//检测db业务逻辑并完成数据库的初始
func InitDBCheck() {
	if service.Config.DbMaster.Driver == "sqlite3" {//如果使用sqlite3数据库引擎的情况
		workDir, _ := os.Getwd()
		dbSource   := filepath.Join(workDir, service.Config.DbMaster.Host)
		if fileInfo, err := os.Stat(dbSource); err != nil || fileInfo.Size() < 1 {//不存在的话完成初始化
			sqliteInitialize(service.Config.DbMaster.Host)
		} else {
			fmt.Println("==================================数据库已完成初始化===============================")
			fmt.Println("================================="+dbSource+"==================================")
		}
	} else if service.Config.DbMaster.Driver == "mysql" {
	
	} else {
		panic("db 存储引擎不支持,无法完成初始化...")
	}
}
