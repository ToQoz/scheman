package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	versionLayout = "20060102150405"
)

func usage() {

	banner := `scheman-g is helper command for generating migration

Usage:

    scheman-g [migration-name]

The options are:

`
	fmt.Fprintf(os.Stderr, banner)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func main() {
	migrationsPath := flag.String("path", "migrations", "migrations path")
	flag.Usage = usage
	flag.Parse()

	name := flag.Arg(0)

	switch name {
	case "":
		flag.Usage()
	default:
		generateMigration(name, *migrationsPath)
	}
}

func generateMigration(name string, migrationsPath string) {
	_, err := os.Stat(migrationsPath)
	if err != nil {
		die(err)
	}

	version := generateVersion()

	for _, kind := range []string{"up", "down"} {
		createMigrationFile(filepath.Join(migrationsPath, version+"_"+name+"_"+kind+".sql"))
	}
}

func generateVersion() string {
	return time.Now().Format(versionLayout)
}

func createMigrationFile(filename string) {
	if err := createFile(filename); err != nil {
		die(err)
	} else {
		fmt.Println("create: " + filename)
	}
}

func createFile(filename string) error {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) { // create file
			file, err := os.OpenFile(filename, os.O_CREATE, 0666)

			if err != nil {
				return err
			}

			return file.Close()
		} else {
			return err
		}
	} else {
		return errors.New(filename + " already exists")
	}
}

func die(err error) {
	fmt.Fprintf(os.Stderr, err.Error())
	os.Exit(1)
}
