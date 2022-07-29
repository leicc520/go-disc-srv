package migrate

import (
	"fmt"
	"os"
	"testing"
)

func TestSqlite3(t *testing.T) {
	dir, _ := os.Getwd()
	dbname := dir + "/go.disc.srv.db"
	fmt.Println(dbname)
	sqliteInitialize(dbname)
}
