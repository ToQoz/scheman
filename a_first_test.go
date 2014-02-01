package scheman

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lestrrat/go-test-mysqld"
	"log"
	"testing"
)

var (
	Mysqld *mysqltest.TestMysqld
)

func TestStartMysqld(t *testing.T) {
	var err error

	log.Print("Starting mysqld")

	mysqld, err := mysqltest.NewMysqld(nil)
	if err != nil {
		log.Fatalf("Failed to start mysqld: %s", err)
	}

	err = mysqld.Start()

	Mysqld = mysqld

	db, err := sql.Open("mysql", Mysqld.Datasource("", "", "", 0))

	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.Exec("SELECT 1")

	if err != nil {
		panic(err)
	}
}
