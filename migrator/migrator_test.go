package migrator

import (
	"database/sql"
)

func newMigrator(db *sql.DB, path string) *Migrator {
	migrator, err := New(db, path)

	if err != nil {
		panic(err)
	}

	return migrator
}
