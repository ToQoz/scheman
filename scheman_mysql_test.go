package scheman

import (
	"database/sql"
	"github.com/ToQoz/scheman/test_helpers"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"testing"
)

var (
	mysqlDBName   = "scheman_tester_test"
	mysqlUser     = "root"
	mysqlPassword = ""
)

func init() {
	u := os.Getenv("DB_USER")

	if u != "" {
		mysqlUser = u
	}

	p := os.Getenv("DB_PASSWORD")

	if p != "" {
		mysqlPassword = p
	}
}

func TestMySQLMigrate(t *testing.T) {
	var err error

	mysqlCreateTestDatabase()
	db := mysqlGetDatabase()

	defer db.Close()              // 2. close database
	defer mysqlDropTestDatabase() // 1. drop database

	migrator := requireMigrator(db, "_test_data/migrations")

	if err = migrator.MigrateTo("20131103115446"); err != nil {
		panic(err)
	}
	test_helpers.AssertEqual(t, "20131103115446", migrator.Version)

	if err = migrator.MigrateTo("20131103115447"); err != nil {
		panic(err)
	}
	test_helpers.AssertEqual(t, "20131103115447", migrator.Version)

	if err = migrator.MigrateTo("20131103115446"); err != nil {
		panic(err)
	}
	test_helpers.AssertEqual(t, "20131103115446", migrator.Version)

	if err = migrator.MigrateTo("20131103115448"); err != nil {
		panic(err)
	}
	test_helpers.AssertEqual(t, "20131103115448", migrator.Version)
}

func TestMySQLRollbackMigration(t *testing.T) {
	var err error

	mysqlCreateTestDatabase()
	db := mysqlGetDatabase()

	defer db.Close()              // 2. close database
	defer mysqlDropTestDatabase() // 1. drop database

	migrator := newMigrator(db, "_test_data/migrations_20131103115449_invalid")

	if err = migrator.MigrateTo("20131103115446"); err != nil {
		panic(err)
	}

	err = migrator.MigrateTo("20131103115449")
	test_helpers.AssertNotEqual(t, nil, err)
	test_helpers.AssertEqual(t, "20131103115446", migrator.Version)
}

func mysqlGetDatabase() *sql.DB {
	db, err := sql.Open("mysql", mysqlUser+":"+mysqlPassword+"@/"+mysqlDBName)

	if err != nil {
		panic(err)
	}

	return db
}

func mysqlCreateTestDatabase() {
	db, err := sql.Open("mysql", mysqlUser+":"+mysqlPassword+"@/")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + mysqlDBName)

	if err != nil {
		panic(err)
	}
}

func mysqlDropTestDatabase() {
	db, err := sql.Open("mysql", mysqlUser+":"+mysqlPassword+"@/")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF EXISTS " + mysqlDBName)

	if err != nil {
		panic(err)
	}
}
