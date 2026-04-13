package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

var (
	ErrNoSources         = errors.New("no sources defined")
	ErrNoFileAtConfigEnv = errors.New("failed to find config file specified by SNIPS_CONFIG")
	ErrNoFileAtDefault   = errors.New("failed to find config in default directory")
)

type SnipsFzfConfig struct {
	UseEnv       bool   `yaml:"use_env"`
	Preview      string `yaml:"preview"`
	PreviewLabel string `yaml:"preview_label"`
	ListLabel    string `yaml:"list_label"`
}

type SnipsConfig struct {
	Sources []string       `yaml:"sources"`
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
		if !filepath.IsAbs(s) {
			config.Sources[i] = filepath.Join(filepath.Dir(path), s)
		}
	}
	return config, err
}
