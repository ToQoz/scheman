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

	migrator := requireMigrator(db, "_test_data/migrations")

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

func TestSQLite3RollbackMigration(t *testing.T) {
	var err error

	db := sqlite3GetDB()

	defer os.Remove(sqlite3DBName)
	defer db.Close()

	migrator := requireMigrator(db, "_test_data/migrations_20131103115449_invalid")

	if err = migrator.MigrateTo("20131103115446"); err != nil {
		panic(err)
	}

	err = migrator.MigrateTo("20131103115449")
	AssertEqual(t, "20131103115446", migrator.Version)
	AssertNotEqual(t, nil, err)
}

func sqlite3GetDB() *sql.DB {
	db, err := sql.Open("sqlite3", sqlite3DBName)

	if err != nil {
		panic(err)
	}

	return db
}
