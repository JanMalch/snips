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
	ErrInvalidAutopick    = errors.New("'auto_pick' must be one of 'never', 'executable_shebang', 'shebang', or 'always'.")
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

type Autopick int

const (
	AutopickNever Autopick = iota
	AutopickExecutableShebang
	AutopickShebang
	AutopickAlways
)

var autopickName = map[Autopick]string{
	AutopickNever:             "never",
	AutopickExecutableShebang: "executable_shebang",
	AutopickShebang:           "shebang",
	AutopickAlways:            "always",
}

func (a Autopick) String() string {
	return autopickName[a]
}

func (a Autopick) Accept(other Autopick) bool {
	if other == AutopickNever {
		return false
	}
	switch a {
	case AutopickNever:
		return false
	case AutopickAlways:
		// other is Executable, Shebang, or Always
		return true
	case AutopickExecutableShebang:
		return other == AutopickExecutableShebang
	case AutopickShebang:
		return other == AutopickShebang
	}
	panic("Unknown Autopick type.")
}

type SnipsConfig struct {
	Sources           []string       `yaml:"sources"`
	IncludeSourceName bool           `yaml:"include_source_name"`
	RawAutopick       string         `yaml:"auto_pick"`
	Runners           []SnipsRunner  `yaml:"runners"`
	Fzf               SnipsFzfConfig `yaml:"fzf"`
}

func (c SnipsConfig) Autopick() (Autopick, error) {
	switch c.RawAutopick {
	case autopickName[AutopickAlways]:
		return AutopickAlways, nil
	case autopickName[AutopickExecutableShebang]:
		return AutopickExecutableShebang, nil
	case autopickName[AutopickShebang]:
		return AutopickShebang, nil
	case autopickName[AutopickNever]:
		return AutopickNever, nil
	default:
		return AutopickNever, ErrInvalidAutopick
	}
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

// Returns the path where the file is stored, which keeps the snippet history.
func HistoryPath() (string, error) {
	path, err := Path()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(path), ".snipshistory"), nil
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
		RawAutopick:       AutopickNever.String(),
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
	if _, err := config.Autopick(); err != nil {
		return SnipsConfig{}, err
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
