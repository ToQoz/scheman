package scheman

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

type Config struct {
	filename       string `json:"-"`
	Host           string
	Port           string
	User           string
	Password       string
	Database       string
	MigrationsPath string
	Version        string
	Encoding       string
	Params         map[string]string
}

func NewConfig(filename string) *Config {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	cfg := newConfigFromData(f)
	cfg.filename = filename
	return cfg
}

func newConfigFromData(data []byte) *Config {
	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		panic(err)
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
