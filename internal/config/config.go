package config

import (
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
	"path/filepath"
)

var errMissingVersion = errors.New("no version number")
var errNoSpecFile = errors.New("no spec file")
var errNoGoFile = errors.New("no go file")
var errNoPackage = errors.New("no package name")
var errNoGeneration = errors.New("no generation method")
var errNoOutPath = errors.New("no output path")
var errUnknownVersion = errors.New("invalid version number")

const errMessageNoVersion = `The configuration file must have a version number.
Set the version to 1 at the top of abnf.json:

{
  "version": "1"
  ...
}
`

const errMessageUnknownVersion = `The configuration file has an invalid version number.
The only supported version is "1".
`

const errMessageNoGenerator = `No generation option was configured; this is a required configuration.
The value must be one of "operators" or "alternatives".
`

const errMessageNoPackages = `No Go package name was configured; this is a required configuration`

const errMessageNoGoFile = `No Go file name was configured; this is a required configuration`

type Config struct {
	Version     string `json:"version" yaml:"version"`
	SpecFile    string `json:"spec" yaml:"spec"`
	Generation  string `json:"generate" yaml:"generate"`
	PackageName string `json:"package" yaml:"package"`
	GoFileName  string `json:"gofile" yaml:"gofile"`
	OutputPath  string `json:"output" yaml:"output"`
	VerboseFlag bool   `json:"verbose" yaml:"verbose"`
}

type versionSetting struct {
	Number string `json:"version" yaml:"version"`
}

func parseConfig(rd io.Reader) (Config, error) {
	var buf bytes.Buffer
	var config Config
	var version versionSetting

	ver := io.TeeReader(rd, &buf)
	dec := yaml.NewDecoder(ver)
	if err := dec.Decode(&version); err != nil {
		return config, err
	}
	if version.Number == "" {
		return config, errMissingVersion
	}
	switch version.Number {
	case "1":
		return v1ParseConfig(&buf)
	default:
		return config, errUnknownVersion
	}
}

func v1ParseConfig(rd io.Reader) (Config, error) {
	dec := yaml.NewDecoder(rd)

	var config Config
	if err := dec.Decode(&config); err != nil {
		return config, err
	}
	if config.Version == "" {
		return config, errMissingVersion
	}
	if config.Version != "1" {
		return config, errUnknownVersion
	}
	if config.GoFileName == "" {
		return config, errNoGoFile
	}
	if config.SpecFile == "" {
		return config, errNoSpecFile
	}
	if config.Generation == "" {
		return config, errNoGeneration
	}
	if config.PackageName == "" {
		return config, errNoPackage
	}
	if config.OutputPath == "" {
		return config, errNoOutPath
	}

	return config, nil
}

func ReadConfig(stderr io.Writer, dir, filename string) (string, *Config, error) {
	configPath := ""
	if filename != "" {
		configPath = filepath.Join(dir, filename)
	} else {
		var yamlMissing, jsonMissing bool
		yamlPath := filepath.Join(dir, "abnf.yml")
		jsonPath := filepath.Join(dir, "abnf.json")

		if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
			yamlMissing = true
		}
		if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
			jsonMissing = true
		}

		if yamlMissing && jsonMissing {
			fmt.Fprintln(stderr, "error parsing configuration files. abnf.yml or abnf.json: file does not exist")
			return "", nil, errors.New("config file missing")
		}

		if !yamlMissing && !jsonMissing {
			fmt.Fprintln(stderr, "error: both abnf.json and abnf.yml files present")
			return "", nil, errors.New("abnf.json and abnf.yml present")
		}

		configPath = yamlPath
		if yamlMissing {
			configPath = jsonPath
		}
	}

	base := filepath.Base(configPath)
	blob, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "error parsing %s: file does not exist\n", base)
		return "", nil, err
	}

	conf, err := parseConfig(bytes.NewReader(blob))
	if err != nil {
		errMessage := fmt.Sprintf("error parsing %s: %s\n", base, err)
		switch err {
		case errMissingVersion:
			errMessage = errMessageNoVersion
		case errUnknownVersion:
			errMessage = errMessageUnknownVersion
		case errNoPackage:
			errMessage = errMessageNoPackages
		case errNoGeneration:
			errMessage = errMessageNoGenerator
		case errNoGoFile:
			errMessage = errMessageNoGoFile
		}

		_, ferr := fmt.Fprintf(stderr, errMessage)
		if ferr != nil {
			log.Fatalln("Error occurred", ferr)
		}
		return "", nil, err
	}

	return configPath, &conf, nil
}
