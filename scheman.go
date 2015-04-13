package scheman

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	Verbose              = true
	migrationFileMatcher = regexp.MustCompile(`\/(\d+)_(\w+)_(?:up|down).sql$`)
	commentMatcher       = regexp.MustCompile(`--.*`)
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

		fp := filepath.Join(mr.migrationsPath, file.Name())

		// Check all migrations(*_up.sql and *_down.sql).
		if migrationFileMatcher.MatchString(fp) {
			data, err := ioutil.ReadFile(fp)

			if err != nil {
				return nil, err
			}

			// if file is empty, return error
			if strings.Trim(string(data), " \n ") == "" {
				return nil, &ErrMigrationFileIsEmpty{fp}
			}

			// Check reverse migration existance
			var rfp string

			switch {
			case strings.HasSuffix(fp, "_up.sql"):
				rfp = strings.Replace(fp, "_up.sql", "_down.sql", 1)
			case strings.HasSuffix(fp, "_down.sql"):
				rfp = strings.Replace(fp, "_down.sql", "_up.sql", 1)
			}

			_, err = os.Stat(rfp)

			if err != nil {
				return nil, fmt.Errorf("%s is found, but %s is not found", fp, rfp)
			}

			// Don't match kind -> next
			if !strings.HasSuffix(fp, fmt.Sprintf("_%s.sql", kind)) {
				continue
			}
		} else {
			continue
		}

		captured := migrationFileMatcher.FindStringSubmatch(fp)
		m := &migration{
			version:  captured[1],
			name:     captured[2],
			filepath: fp,
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
		if Verbose {
			fmt.Println("Nothing to do")
		}
		return nil
	}

	if Verbose {
		fmt.Println("\n=== migrations.Begin ===")
	}

	tx, _ := db.Begin()

	for _, migration := range ms {
		if Verbose {
			fmt.Printf("%5s: %s %s\n", migration.kind, migration.version, migration.name)
		}
		err := migration.migrate(tx)

		if err != nil {
			if Verbose {
				fmt.Println(err.Error())
			}
			tx.Rollback()
			if Verbose {
				fmt.Println("=== migrations.Rollback!!! ===\n")
			}
			return err
		}
	}

	tx.Commit()

	if Verbose {
		fmt.Println("=== migrations.End ===\n")
	}

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

	stmts := commentMatcher.ReplaceAll(buf, []byte{})
	for _, stmt := range splitStmt(stmts) {
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

// If multiple stmts given, this split stmts at ";".
// WARNING: This is realized by very ugly and not good way.
//          But probably this is not problem in general cases.
//          Generally, DDL uses ";" as end of stmt, or comment.
func splitStmt(stmtsData []byte) (stmts []string) {
	stmtsStr := string(stmtsData)
	var stmt string

	for {
		stmtsStr = strings.Trim(stmtsStr, " \n")
		if stmtsStr == "" {
			break
		}

		i := strings.Index(stmtsStr, ";")

		// "foo" ->
		//     stmt: "foo;"
		//     stmtsStr: ""
		// "foo;bar;foobar;" ->
		//     stmt: "foo;"
		//     stmtsStr: "bar;foobar;"
		if i == -1 {
			stmt = stmtsStr
			stmtsStr = ""
		} else {
			stmt = stmtsStr[:i+1]
			stmtsStr = stmtsStr[i+1:]
		}

		stmts = append(stmts, stmt)
	}

	return stmts
}
