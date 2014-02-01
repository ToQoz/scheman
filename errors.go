package scheman

import (
	"fmt"
)

// ----------------------------------------------------------------------------
// ErrInQuery
// ----------------------------------------------------------------------------

type ErrQuery struct {
	query  string
	detail string
}

func (e ErrQuery) Error() string {
	return fmt.Sprintf("Error in query `%s`: %s", e.query, e.detail)
}

// ----------------------------------------------------------------------------
// ErrInMigration
// ----------------------------------------------------------------------------

type ErrMigration struct {
	migration *migration
	detail    string
}

func (e ErrMigration) Error() string {
	return fmt.Sprintf("Error in migration `%s`: %s", e.migration.filepath, e.detail)
}

// ----------------------------------------------------------------------------
// ErrMigrationFileIsEmpty
// ----------------------------------------------------------------------------

type ErrMigrationFileIsEmpty struct {
	filepath string
}

func (e ErrMigrationFileIsEmpty) Error() string {
	return fmt.Sprintf("Migration file `%s` is empty.", e.filepath)
}
