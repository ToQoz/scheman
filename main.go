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
	pathOpt       = flag.String("path", "migrations", "migrations path")
	nameOpt       = flag.String("name", "", "migration name.(e.g. create_posts)")
	versionLayout = "20060102150405"
)

func main() {
	var err error

	flag.Parse()

	if *nameOpt == "" {
		die(errors.New("-name is required"))
	}

	_, err = os.Stat(*pathOpt)
	if err != nil {
		die(err)
	}

	version := generateVersion()

	for _, kind := range []string{"up", "down"} {
		createMigration(filepath.Join(*pathOpt, version+"_"+*nameOpt+"_"+kind+".sql"))
	}
}

func createMigration(filename string) {
	if err := create(filename); err != nil {
		die(err)
	} else {
		fmt.Println("create: " + filename)
	}
}

func create(filename string) error {
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
	fmt.Println(err.Error())
	os.Exit(1)
}

func generateVersion() string {
	return time.Now().Format(versionLayout)
}
