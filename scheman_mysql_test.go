package scheman

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestMySQLMigrate(t *testing.T) {
	var err error

	db := mysqlGetDatabase()

	defer db.Close()              // 2. close database
	defer mysqlDropTestDatabase() // 1. drop database

	migrator := requireMigrator(db, "testdata/migrations")

	if err = migrator.MigrateTo("20131103115446"); err != nil {
		panic(err)
	}
	AssertEqual(t, "20131103115446", migrator.Version)

	if err = migrator.MigrateTo("20131103115447"); err != nil {
		panic(err)
	}
	AssertEqual(t, "20131103115447", migrator.Version)

	if err = migrator.MigrateTo("20131103115446"); err != nil {
		panic(err)
	}
	AssertEqual(t, "20131103115446", migrator.Version)

	if err = migrator.MigrateTo("20131103115448"); err != nil {
		panic(err)
	}
	AssertEqual(t, "20131103115448", migrator.Version)
}

func TestMySQLRollbackMigration(t *testing.T) {
	var err error

	db := mysqlGetDatabase()

	defer db.Close()              // 2. close database
	defer mysqlDropTestDatabase() // 1. drop database

	migrator := requireMigrator(db, "testdata/migrations_20131103115449_invalid")

	if err = migrator.MigrateTo("20131103115446"); err != nil {
		panic(err)
	}

	err = migrator.MigrateTo("20131103115449")
	AssertNotEqual(t, nil, err)
	AssertEqual(t, "20131103115446", migrator.Version)
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

func mysqlGetDatabase() *sql.DB {
	_mysqlCreateTestDatabase()

	db, err := sql.Open("mysql", Mysqld.Datasource("mysqltest", "", "", 0))

	if err != nil {
		panic(err)
	}

	return db
}

func mysqlDropTestDatabase() {
	db, err := sql.Open("mysql", Mysqld.Datasource("", "", "", 0))

	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF EXISTS mysqltest")

	if err != nil {
		panic(err)
	}
}

func _mysqlCreateTestDatabase() {
	db, err := sql.Open("mysql", Mysqld.Datasource("", "", "", 0))

	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS mysqltest")

	if err != nil {
		panic(err)
	}
}
