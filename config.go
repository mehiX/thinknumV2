package thinknum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	filterColDateName = "as_of_date"
)

type Config struct {
	Hostname     string             `json:"hostname"`
	Version      string             `json:"version"`
	ClientID     string             `json:"client_id"`
	ClientSecret string             `json:"client_secret"`
	Workers      int                `json:"workers"`
	PageSize     int                `json:"page_size"`
	Searches     []SearchDefinition `json:"searches"`
}

// ConfigFromJSON Loads configuration data from a JSON file
func ConfigFromJSON(fn string) (*Config, error) {

	var cfg Config

	fi, err := os.Stat(fn)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() || !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("Config file %s is not a regular file", fi)
	}

	all, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(bytes.NewBuffer(all)).Decode(&cfg); err != nil {
		return nil, err
	}

	err = validate(cfg)

	return &cfg, err
}

// validate Validates that the configuration is vallid
func validate(cfg Config) error {
	for _, s := range cfg.Searches {
		// check that the output file is not a directory
		info, err := os.Stat(s.OutputFile)
		// only check if it is a directory for now. If it is a file it is probably missing since it will be created later
		if err == nil && info.IsDir() {
			return fmt.Errorf("outputFile for Search %s should not be a directory", s.Name)
		}

		// check that the parent directory of the output file is writable
		if !isDirWritable(filepath.Dir(s.OutputFile)) {
			return fmt.Errorf("cannot write to file: %s", s.OutputFile)
		}
	}

	return nil
}

// permission bits
const (
	GroupWrite fs.FileMode = 1 << (7 - 3*iota)
	UserWrite
	OthersWrite
)

// isDirWritable Checks that the directory specified as parameter is writable
func isDirWritable(dirPath string) bool {

	info, err := os.Stat(dirPath)
	if err != nil {
		log.Println(err)
		return false
	}

	mode := info.Mode().Perm()
	if mode&OthersWrite == 0 && mode&UserWrite == 0 && mode&GroupWrite == 0 {
		return false
	}

	return true
}
