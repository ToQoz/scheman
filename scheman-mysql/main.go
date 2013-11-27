package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ToQoz/scheman/migrator"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"reflect"
)

var (
	cfg *config
)

type config struct {
	filename       string `json:"-`
	User           string
	Password       string
	Database       string
	MigrationsPath string
	Version        string
	Encoding       string
}

func newConfig(filename string) *config {
	cfg := &config{filename: filename}

	f, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(f, cfg)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	return cfg
}

func (cfg *config) Require(key string) {
	v := cfg.Get(key)

	if v == "" {
		fmt.Fprintf(os.Stderr, "<"+cfg.filename+"> "+key+" should not be empty")
		os.Exit(1)
	}
}

func (cfg *config) Get(key string) string {
	return reflect.ValueOf(*cfg).FieldByName(key).String()
}

func usage() {
	banner := `scheman-mysql is github.com/ToQoz/scheman's frontend for MySQL

Usage:

	scheman-mysql [sub-command]

The sub-commands are:

	create
	drop
	migrate
	reset

The options are:

`
	fmt.Fprintf(os.Stderr, banner)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func createDatabase() {
	db, err := sql.Open("mysql", cfg.User+":"+cfg.Password+"@/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Fprintf(os.Stderr, "Create database "+cfg.Database+" with caracter "+cfg.Encoding+"\n\n")
	_, err = db.Exec("CREATE DATABASE " + cfg.Database + " DEFAULT CHARACTER SET " + cfg.Encoding)
	if err != nil {
		panic(err)
	}
}

func dropDatabase() {
	db, err := sql.Open("mysql", cfg.User+":"+cfg.Password+"@/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Fprintf(os.Stderr, "Drop database "+cfg.Database+"\n\n")
	_, err = db.Exec("DROP DATABASE " + cfg.Database)
	if err != nil {
		panic(err)
	}
}

func migrateDatabase() {
	db, err := sql.Open("mysql", cfg.User+":"+cfg.Password+"@/"+cfg.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	migrator, err := migrator.New(db, cfg.MigrationsPath)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "Migrate database "+cfg.Database+" to "+cfg.Version+"\n\n")
	err = migrator.MigrateTo(cfg.Version)
	if err != nil {
		panic(err)
	}
}

func main() {
	config_filename := flag.String("c", "scheman.json", "scheman configuration json file")
	flag.Usage = usage
	flag.Parse()

	cfg = newConfig(*config_filename)

	cfg.Require("User")
	cfg.Require("Database")
	cfg.Require("MigrationsPath")
	cfg.Require("Version")

	if cfg.Encoding == "" {
		cfg.Encoding = "utf8"
	}

	subcmd := flag.Arg(0)

	switch subcmd {
	case "create":
		createDatabase()
	case "drop":
		dropDatabase()
	case "migrate":
		migrateDatabase()
	case "reset":
		dropDatabase()
		createDatabase()
		migrateDatabase()
	default:
		flag.Usage()
	}
}
