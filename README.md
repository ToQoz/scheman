# Scheman

[![Build Status](https://travis-ci.org/ToQoz/scheman.png?branch=master)](https://travis-ci.org/ToQoz/scheman)

!!! Under development !!!

scheman is database schema migration tool.

## Testing

- Install scheman. `go get github.com/ToQoz/scheman`
- Install dependencies for test. `go list -f '{{.TestImports}}' github.com/ToQoz/scheman/... | sed 's/\[//g' | sed 's/\]//g' | xargs go get`
- Run tests. `go test github.com/ToQoz/scheman/...`

## Usage

### Generate migration

```
$ mkdir migrations
$ go get github.com/ToQoz/scheman/scheman-g
$ scheman-g create_posts
create: migrations/20131103115446_create_posts_up.sql
create: migrations/20131103115446_create_posts_down.sql

# if you want to specify migrations directory
$ scheman-g -path ./sql create_posts
create: sql/20131103115446_create_posts_up.sql
create: sql/20131103115446_create_posts_down.sql
```

### Execute migrations

#### If you use MySQL and want to use frontend that provided by scheman.

```
$ go get github.com/ToQoz/scheman/scheman-mysql
$ vi scheman.json
$ cat !$
cat scheman.json
{
  "User":           "mysql login user",
  "Password":       "mysql login password",
  "Database":       "mysql database name",
  "Encoding":       "mysql database encoding(default utf8)",
  "MigrationsPath": "migration files directory path",
  "Version":        "migration version"
}
$ scheman-mysql create
$ scheman-mysql migrate
```

See also [scheman-mysql's README](http://github.com/ToQoz/scheman/tree/master/scheman-mysql)

#### If you use other RDBMS or want to use your own frontend.

```
$ vi migrate.go # Write with reference to http://github.com/ToQoz/scheman/tree/master/scheman-mysql
$ go run migrate.go
```
