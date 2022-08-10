package migrate

import (
	"fmt"
	"github.com/leicc520/go-disc-srv/app/models"
	"github.com/leicc520/go-orm"
	"os"
	"testing"
	"time"
)

func TestSqlite3(t *testing.T) {
	dir, _ := os.Getwd()
	dbname := dir + "/go.disc.srv.db"
	fmt.Println(dbname)
	sqliteInitialize(dbname)
}

func TestUser(t *testing.T) {
	dir, _ := os.Getwd()
	dbname := dir + "/go.disc.srv.db"
	dbConfig := orm.DbConfig{Driver: "sqlite3", Host: dbname}
	orm.InitDBPoolSt().Set("dbmaster", &dbConfig)
	orm.InitDBPoolSt().Set("dbslaver", &dbConfig)
	user  := models.SysUserSt{}
	err   := models.NewSysUser().GetItem(func(st *orm.QuerySt) string {
		st.Where("account", "admin")
		return st.GetWheres()
	}, "*").ToStruct(&user)
	
	fmt.Println(user, user.Id, err, time.Now().Unix())
}
