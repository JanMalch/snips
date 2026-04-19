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
	ErrNoFileAtConfigEnv  = errors.New("failed to find config file specified by SNIPS_CONFIG")
	ErrNoFileAtDefault    = errors.New("failed to find config in default directory")
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
	Sources []string       `yaml:"sources"`
	Runners []SnipsRunner  `yaml:"runners"`
	Fzf     SnipsFzfConfig `yaml:"fzf"`
}

func Load() (SnipsConfig, error) {
	path := os.Getenv("SNIPS_CONFIG")
	usesExplicitConfig := path != ""
	if path == "" {
		cfg, err := os.UserConfigDir()
		if err != nil {
			return SnipsConfig{}, err
		}
		path = filepath.Join(cfg, "snips", "config.yaml")
	}

	dat, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if usesExplicitConfig {
				return SnipsConfig{}, ErrNoFileAtConfigEnv
			} else {
				return SnipsConfig{}, ErrNoFileAtDefault
			}
		}
		return SnipsConfig{}, err
	}

	config := SnipsConfig{
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
