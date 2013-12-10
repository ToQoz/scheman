package scheman

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

type Config struct {
	filename       string `json:"-`
	User           string
	Password       string
	Database       string
	MigrationsPath string
	Version        string
	Encoding       string
}

func NewConfig(filename string) *Config {
	cfg := &Config{filename: filename}

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

func (cfg *Config) Require(key string) {
	v := cfg.Get(key)

	if v == "" {
		fmt.Fprintf(os.Stderr, "<"+cfg.filename+"> "+key+" should not be empty")
		os.Exit(1)
	}
}

func (cfg *Config) Get(key string) string {
	return reflect.ValueOf(*cfg).FieldByName(key).String()
}
