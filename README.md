# Scheman

!!! Under development !!!

scheman is database schema migration tool.

## Testing

- Install scheman. `git get github.com/ToQoz/scheman`
- Install dependencies for test. `go list -f '{{.TestImports}}' github.com/ToQoz/scheman | sed 's/\[//g' | sed 's/\]//g' | xargs go get`
- Run tests. `go test github.com/ToQoz/scheman`

## Examples

```
$ mkdir migrations
$ scheman -name create_posts
create: migrations/20131103115446_create_posts_up.sql
create: migrations/20131103115446_create_posts_down.sql

# if you want to specify migrations directory
$ scheman -path ./sql -name create_posts
create: sql/20131103115446_create_posts_up.sql
create: sql/20131103115446_create_posts_down.sql
```

```go
package main

import (
	"database/sql"
	"github.com/ToQoz/scheman"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "user:passwd@/dbname")
	if err != nil {
		panic(err)
	}

	migrator, err := scheman.NewMigrator(db, "./migrations")
	if err != nil {
		panic(err)
	}

	err = migrator.MigrateTo("20131103115446")
	if err != nil {
		panic(err)
	}
}
```
