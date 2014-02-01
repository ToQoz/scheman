package scheman

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ----------------------------------------------------------------------------
// Migrator
// ----------------------------------------------------------------------------

type Migrator struct {
	Version        string
	db             *sql.DB
	targetVersion  string
	migrationsPath string
}

func NewMigrator(db *sql.DB, migrationsPath string) (*Migrator, error) {
	mr := &Migrator{
		db:             db,
		migrationsPath: migrationsPath,
	}

	err := mr.createMigrationTableIfNotExists()

	if err != nil {
		return nil, err
	}

	version, err := mr.version()

	if err != nil {
		return nil, err
	}

	mr.Version = version

	return mr, nil
}

func (mr *Migrator) MigrateTo(targetVersion string) error {
	var kind string

	mr.targetVersion = targetVersion

	if mr.targetVersion < mr.Version {
		kind = "down"
	} else {
		kind = "up"
	}

	ms, err := mr.NewMigrations(kind)

	if err != nil {
		return err
	}

	err = ms.migrate(mr.db)

	if err != nil {
		return err
	}

	// Update cache for current version
	mr.Version = mr.targetVersion

	return nil
}

func (mr *Migrator) NewMigrations(kind string) (migrations, error) {
	matcher := regexp.MustCompile(`\/(\d+)_(\w+)_` + kind + `.sql$`)
	migrations := migrations{}

	files, err := ioutil.ReadDir(mr.migrationsPath)
	if err != nil {
		return nil, err
	}

	uppedVersions, err := mr.uppedVersions()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fpath := filepath.Join(mr.migrationsPath, file.Name())

		if !matcher.MatchString(fpath) {
			continue
		}

		matched := matcher.FindStringSubmatch(fpath)
		m := &migration{
			version:  matched[1],
			name:     matched[2],
			filepath: fpath,
			kind:     kind,
		}

		if m.kind == "up" {
			upped := containString(uppedVersions, m.version)
			inRange := mr.Version <= m.version && mr.targetVersion >= m.version
			if !upped && inRange {
				migrations = append(migrations, m)
			}
		} else {
			downed := !containString(uppedVersions, m.version)
			inRange := mr.Version >= m.version && mr.targetVersion < m.version
			if !downed && inRange {
				migrations = append(migrations, m)
			}
		}
	}

	if kind == "up" {
		sort.Sort(migrations)
	} else {
		sort.Reverse(migrations)
	}

	return migrations, nil
}

func (mr *Migrator) createMigrationTableIfNotExists() error {
	q := "CREATE TABLE IF NOT EXISTS scheman_versions ("
	q += "  version CHAR(64) NOT NULL PRIMARY KEY"
	q += ");"

	_, err := mr.db.Exec(q)

	if err != nil {
		return ErrQuery{query: q, detail: err.Error()}
	}

	return nil
}

func (mr *Migrator) version() (string, error) {
	q := "SELECT version FROM scheman_versions ORDER BY version DESC LIMIT 1"

	var version string

	err := mr.db.QueryRow(q).Scan(&version)

	switch {
	case err == sql.ErrNoRows:
		return "0", nil
	case err != nil:
		return "", ErrQuery{query: q, detail: err.Error()}
	default:
		return version, nil
	}
}

func (mr *Migrator) uppedVersions() ([]string, error) {
	q := "SELECT version FROM scheman_versions"

	rows, err := mr.db.Query(q)

	if err != nil {
		return nil, ErrQuery{query: q, detail: err.Error()}
	}

	defer rows.Close()

	versions := []string{}

	for rows.Next() {
		var version string
		err = rows.Scan(&version)

		if err != nil {
			return nil, err
		}

		versions = append(versions, version)
	}

	return versions, err
}

// ----------------------------------------------------------------------------
// Migratoions
// ----------------------------------------------------------------------------

type migrations []*migration

func (ms migrations) migrate(db *sql.DB) error {
	if len(ms) == 0 {
		fmt.Println("Nothing to do")
		return nil
	}

	fmt.Println("\n=== migrations.Begin ===")
	tx, _ := db.Begin()

	for _, migration := range ms {
		fmt.Printf("%5s: %s %s\n", migration.kind, migration.version, migration.name)
		err := migration.migrate(tx)

		if err != nil {
			fmt.Println(err.Error())
			tx.Rollback()
			fmt.Println("=== migrations.Rollback!!! ===\n")
			return err
		}
	}

	tx.Commit()
	fmt.Println("=== migrations.End ===\n")

	return nil
}

func (ms migrations) Len() int {
	return len(ms)
}

func (ms migrations) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (ms migrations) Less(i, j int) bool {
	return ms[i].version < ms[j].version
}

// ----------------------------------------------------------------------------
// Migratoions
// ----------------------------------------------------------------------------

type migration struct {
	version  string
	name     string
	filepath string
	kind     string
}

func (m *migration) migrate(db *sql.Tx) error {
	buf, err := ioutil.ReadFile(m.filepath)

	if err != nil {
		return err
	}

	strbuf := removeComment(string(buf))
	fmt.Print(strbuf)

	// If multiple stmts given, this split stmts at ";".
	// WARNING: This is realized by very ugly and not good way.
	//          But probably this is not problem in general cases.
	//          Generally, DDL uses ";" as end of stmt, or comment.

	var stmt string

	for {
		if strings.Trim(strbuf, " \n") == "" {
			break
		}

		i := strings.Index(strbuf, ";")

		// "foo" ->
		//     stmt: "foo;"
		//     strbuf: ""
		// "foo;bar;foobar;" ->
		//     stmt: "foo;"
		//     strbuf: "bar;foobar;"
		if i == -1 {
			stmt = strbuf
			strbuf = ""
		} else {
			stmt = strbuf[:i+1]
			strbuf = strbuf[i+1:]
		}

		_, err = db.Exec(stmt)

		if err != nil {
			return ErrMigration{migration: m, detail: err.Error()}
		}
	}

	err = m.updateVersion(db)

	if err != nil {
		return err
	}

	return nil
}

func (m *migration) updateVersion(db *sql.Tx) error {
	var q string

	if m.kind == "up" {
		q = "INSERT INTO scheman_versions ( version ) VALUES ( ? );"
	} else {
		q = "DELETE FROM scheman_versions WHERE version = ?;"
	}

	_, err := db.Exec(q, m.version)

	if err != nil {
		return ErrQuery{query: q, detail: err.Error()}
	}

	return nil
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

func containString(slice []string, s string) bool {
	for _, sliceS := range slice {
		if sliceS == s {
			return true
		}
	}

	return false
}

func removeComment(q string) string {
	r := regexp.MustCompile(`--.*`)
	return r.ReplaceAllString(q, "")
}
