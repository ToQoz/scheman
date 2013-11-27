package scheman

import (
	"database/sql"
)

func newMigrator(db *sql.DB, path string) *Migrator {
	s, err := NewMigrator(db, path)

	if err != nil {
		panic(err)
	}

	return s
}
