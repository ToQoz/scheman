package scheman

import (
	"database/sql"
)

func newMigrator(db *sql.DB, path string) *Migrator {
	migrator, err := NewMigrator(db, path)

	if err != nil {
		panic(err)
	}

	return migrator
}
