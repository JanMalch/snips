package core

import (
	"errors"
	"os"
	"path/filepath"

	yaml "github.com/goccy/go-yaml"
)

var (
	ErrEmptyInclude = errors.New("snips file found, but include isn't defined")
)

type SnipsSource struct {
	Include []string `yaml:"include"`
}

func SourceFromDirectory(dir string) (SnipsSource, error) {
	dat, err := os.ReadFile(filepath.Join(dir, "snips.yaml"))
	if err != nil {
		if !os.IsNotExist(err) {
			// snips.yaml exists, but we cannot read it
			return SnipsSource{}, err
		}
		// snips.yaml doesn't exist
		dat, err = os.ReadFile(filepath.Join(dir, "snips.yml"))
		if err != nil {
			if os.IsNotExist(err) {
				// snips.yml also doesn't exist
				return SnipsSource{
					Include: []string{"**/*"},
				}, nil
			} else {
				// snips.yml exist, but we cannot read it
				return SnipsSource{}, err
			}
		}
	}
	return SourceFromBytes(dat)
}

func SourceFromBytes(b []byte) (SnipsSource, error) {
	var source SnipsSource
	if err := yaml.Unmarshal(b, &source); err != nil {
		return SnipsSource{}, err
	}
	if len(source.Include) == 0 {
		return SnipsSource{}, ErrEmptyInclude
	}
	return source, nil
}
