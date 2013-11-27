package scheman

import (
	"database/sql"
	"testing"
)

func requireMigrator(db *sql.DB, path string) *Migrator {
	s, err := NewMigrator(db, path)

	if err != nil {
		panic(err)
	}

	return s
}

func AssertEqual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("expected <%s>, but got <%s>", expected, actual)
	}
}

func AssertNotEqual(t *testing.T, expected, actual interface{}) {
	if expected == actual {
		t.Errorf("not expected <%s>, but got <%s>", expected, actual)
	}
}

