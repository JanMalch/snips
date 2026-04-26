package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

var (
	ErrNoSources          = errors.New("no sources defined")
	ErrNoExtensionDefined = errors.New("no 'ext' or 'exts' defined for runner")
)

type SnipsFzfConfig struct {
	UseEnv       bool   `yaml:"use_env"`
	Preview      string `yaml:"preview"`
	PreviewLabel string `yaml:"preview_label"`
	ListLabel    string `yaml:"list_label"`
}

type SnipsRunner struct {
	Ext  string   `yaml:"ext"`
	Exts []string `yaml:"exts"`
	Name string   `yaml:"name"`
	Args []string `yaml:"args"`
}

func (r SnipsRunner) Matches(fileext string) bool {
	if r.Ext != "" {
		if eqExt(r.Ext, fileext) {
			return true
		}
	}
	for _, ext := range r.Exts {
		if eqExt(ext, fileext) {
			return true
		}
	}
	return false
}

type SnipsConfig struct {
	Sources           []string       `yaml:"sources"`
	IncludeSourceName bool           `yaml:"include_source_name"`
	Runners           []SnipsRunner  `yaml:"runners"`
	Fzf               SnipsFzfConfig `yaml:"fzf"`
}

// Returns the path of the config
func Path() (string, error) {
	path := os.Getenv("SNIPS_CONFIG")
	if path == "" {
		cfg, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(cfg, "snips", "config.yaml")
	}
	return path, nil
}

func Load() (SnipsConfig, error) {
	path, err := Path()
	if err != nil {
		return SnipsConfig{}, err
	}
	dat, err := os.ReadFile(path)
	if err != nil {
		return SnipsConfig{}, err
	}

	config := SnipsConfig{
		IncludeSourceName: true,
		Fzf: SnipsFzfConfig{
			Preview:      "cat {1}",
			UseEnv:       true,
			PreviewLabel: fzfOptPreviewLabel,
			ListLabel:    fzfOptListLabel,
		},
	}
	if err := yaml.Unmarshal(dat, &config); err != nil {
		return SnipsConfig{}, err
	}
	for i, s := range config.Sources {
		if s[0] == '~' {
			home, err := os.UserHomeDir()
			if err != nil {
				return SnipsConfig{}, err
			}
			abs, err := filepath.Abs(strings.ReplaceAll(s, "~", home))
			if err != nil {
				return SnipsConfig{}, err
			}
			config.Sources[i] = abs
		} else if !filepath.IsAbs(s) {
			config.Sources[i] = filepath.Join(filepath.Dir(path), s)
		}
	}
	for _, r := range config.Runners {
		if r.Ext == "" && len(r.Exts) == 0 {
			return SnipsConfig{}, ErrNoExtensionDefined
		}
	}
	return config, err
}

func eqExt(actual, expected string) bool {
	if actual == "" {
		return false
	}
	if actual[0] == '.' {
		return actual == expected
	}
	return ("." + actual) == expected
}
