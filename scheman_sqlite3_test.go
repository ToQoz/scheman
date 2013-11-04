package scheman

import (
	"database/sql"
	"github.com/ToQoz/scheman/test_helpers"
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

	migrator := newMigrator(db, "./migrations")

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

func TestSQLite3RollbackMigration(t *testing.T) {
	var err error

	db := sqlite3GetDB()

	defer os.Remove(sqlite3DBName)
	defer db.Close()

	migrator := newMigrator(db, "./migrations_20131103115449_invalid")

	if err = migrator.MigrateTo("20131103115446"); err != nil {
		panic(err)
	}

	err = migrator.MigrateTo("20131103115449")
	test_helpers.AssertEqual(t, "20131103115446", migrator.Version)
	test_helpers.AssertNotEqual(t, nil, err)
}

func sqlite3GetDB() *sql.DB {
	db, err := sql.Open("sqlite3", sqlite3DBName)

	if err != nil {
		panic(err)
	}

	return db
}
