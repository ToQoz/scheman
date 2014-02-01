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

	migrator, err := NewMigrator(db, "testdata/migrations")

	if err != nil {
		t.Errorf("unexpected err: %s", err)
	}

	err = migrator.MigrateTo("20131103115446")

	if err != nil {
		t.Errorf("unexpected err: %s", err)
	}

	if expected := "20131103115446"; migrator.Version != expected {
		t.Errorf("expected version %s, but got %s", expected, migrator.Version)
	}

	err = migrator.MigrateTo("20131103115447")

	if err != nil {
		t.Errorf("unexpected err: %s", err)
	}

	if expected := "20131103115447"; migrator.Version != expected {
		t.Errorf("expected version %s, but got %s", expected, migrator.Version)
	}

	err = migrator.MigrateTo("20131103115446")

	if err != nil {
		t.Errorf("unexpected err: %s", err)
	}

	if expected := "20131103115446"; migrator.Version != expected {
		t.Errorf("expected version %s, but got %s", expected, migrator.Version)
	}

	err = migrator.MigrateTo("20131103115448")

	if err != nil {
		t.Errorf("unexpected err: %s", err)
	}

	if expected := "20131103115448"; migrator.Version != expected {
		t.Errorf("expected version %s, but got %s", expected, migrator.Version)
	}
}

func TestMySQLRollbackMigration(t *testing.T) {
	var err error

	db := mysqlGetDatabase()

	defer db.Close()              // 2. close database
	defer mysqlDropTestDatabase() // 1. drop database

	migrator, err := NewMigrator(db, "testdata/migrations_20131103115449_invalid")
	if err != nil {
		panic(err)
	}

	if err = migrator.MigrateTo("20131103115446"); err != nil {
		panic(err)
	}

	err = migrator.MigrateTo("20131103115449")

	if err == nil {
		t.Error("expected error on migration, but not got it.")
	}

	if expected := "20131103115446"; migrator.Version != expected {
		t.Errorf("expected version %s, but got %s", expected, migrator.Version)
	}
}

func TestMySQLMultipleStmts(t *testing.T) {
	var err error

	db := mysqlGetDatabase()

	defer db.Close()              // 2. close database
	defer mysqlDropTestDatabase() // 1. drop database

	migrator, err := NewMigrator(db, "testdata/migrations_multiple_stmts")

	if err != nil {
		t.Error("expected error on migration, but not got it.")
	}

	err = migrator.MigrateTo("1")

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if expected := "1"; migrator.Version != expected {
		t.Errorf("expected version %s, but got %s", expected, migrator.Version)
	}

	err = migrator.MigrateTo("2")

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if expected := "2"; migrator.Version != expected {
		t.Errorf("expected version %s, but got %s", expected, migrator.Version)
	}
}

func TestMySQLEmptyMigrationFile(t *testing.T) {
	var err error

	db := mysqlGetDatabase()

	defer db.Close()              // 2. close database
	defer mysqlDropTestDatabase() // 1. drop database

	migrator, err := NewMigrator(db, "testdata/migrations_has_empty_sqlfile")

	if err != nil {
		t.Error("Unexpected error, %s", err)
	}

	err = migrator.MigrateTo("1")

	if err == nil {
		t.Error("error expected, but not got.")
	}
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

func mysqlGetDatabase() *sql.DB {
	mysqlCreateTestDatabase()

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

func mysqlCreateTestDatabase() {
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
