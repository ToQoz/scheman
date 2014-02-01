package scheman

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

var (
	sqlite3DBName = "./scheman_tester.sqlite3"
)

func TestSQLite3Migrate(t *testing.T) {
	var err error

	db := sqlite3GetDB()

	defer os.Remove(sqlite3DBName)
	defer db.Close()

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

func TestSQLite3RollbackMigration(t *testing.T) {
	var err error

	db := sqlite3GetDB()

	defer os.Remove(sqlite3DBName)
	defer db.Close()

	migrator, err := NewMigrator(db, "testdata/migrations_20131103115449_invalid")
	if err != nil {
		t.Errorf("unexpected err: %s", err)
	}

	err = migrator.MigrateTo("20131103115446")

	if err != nil {
		t.Errorf("unexpected err: %s", err)
	}

	err = migrator.MigrateTo("20131103115449")

	if err == nil {
		t.Errorf("validation error is expected, but not got it.")
	}

	if expected := "20131103115446"; migrator.Version != expected {
		t.Errorf("expected version %s, but got %s", expected, migrator.Version)
	}
}

func TestSQLite3MultipleStmt(t *testing.T) {
	var err error

	db := sqlite3GetDB()

	defer os.Remove(sqlite3DBName)
	defer db.Close()

	migrator, err := NewMigrator(db, "testdata/migrations_multiple_stmts")

	if err != nil {
		t.Errorf("unexpected err: %s", err)
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

func TestSQLite3EmptyMigrationFile(t *testing.T) {
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

func TestSQLite3ReversedMigrationFileIsNotFound(t *testing.T) {
	var err error

	db := mysqlGetDatabase()

	defer db.Close()              // 2. close database
	defer mysqlDropTestDatabase() // 1. drop database

	migrator, err := NewMigrator(db, "testdata/migrations_reverse_migration_is_empty")

	if err != nil {
		t.Error("Unexpected error, %s", err)
	}

	err = migrator.MigrateTo("1")

	if err == nil {
		t.Error("error expected, but not got.")
	}
}

func sqlite3GetDB() *sql.DB {
	db, err := sql.Open("sqlite3", sqlite3DBName)

	if err != nil {
		panic(err)
	}

	return db
}
