package scheman

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lestrrat/go-test-mysqld"
	"log"
	"os/exec"
	"strings"
	"testing"
)

var (
	Mysqld *mysqltest.TestMysqld
)

func init() {
	Verbose = false
}

func TestStartMysqld(t *testing.T) {
	var err error

	log.Print("Starting mysqld")
	mysqld, err := mysqltest.NewMysqld(nil)
	if err != nil {
		log.Fatalf("Failed to start mysqld: %s", err)
	}

	Mysqld = mysqld
	db, err := sql.Open("mysql", Mysqld.Datasource("", "", "", 0))
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()

	_, err = db.Exec("SELECT 1")
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestThisFileIsFirstTestFile(t *testing.T) {
	cmd := exec.Command("go", "list", "-f", "{{.TestGoFiles}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err.Error())
	}

	tests := strings.Split(strings.Trim(string(output), " \n[]"), " ")
	lastTest := tests[0]

	if lastTest != "a_first_test.go" {
		t.Fatalf("expected last_test is a_first_test.go, but got %v", lastTest)
	}
}
