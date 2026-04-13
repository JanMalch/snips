package core

import (
	"errors"
	"path/filepath"
	"snips/internal/config"
	"strings"

	fzf "github.com/junegunn/fzf/src"
)

const unitSep = "\x1F"

var (
	ErrNoSnippetFound = errors.New("no snippet found")
	ErrNoMatches      = errors.New("no matching snippets found")
	ErrNoSources      = errors.New("no script sources set")
	ErrNoCmd          = errors.New("no cmd in source")
	ErrNoFzfOutput    = errors.New("fzf returned no output")
)

func FindSnippet(query string, dirs []string, cfg config.SnipsFzfConfig) (string, error) {
	if len(dirs) == 0 {
		return "", ErrNoSources
	}

	matches := make([]string, 0)

	for _, dir := range dirs {

		dir, err := filepath.Abs(dir)
		if err != nil {
			return "", err
		}
		source, err := SourceFromDirectory(dir)
		if err != nil {
			return "", err
		}

		for _, include := range source.Include {
			imatches, err := Doublestar(dir, include)
			if err != nil {
				return "", err
			}
			for _, s := range imatches {
				rels, err := filepath.Rel(dir, s)
				if err != nil {
					return "", err
				}
				matches = append(matches, s+unitSep+rels)
			}
		}
	}
	if len(matches) == 0 {
		return "", ErrNoMatches
	}

	inputChan := make(chan string)
	go func() {
		for _, m := range matches {
			inputChan <- m
		}
		close(inputChan)
	}()
	// TODO: use header
	return runFzf("", inputChan, query, cfg)
}

func runFzf(header string, input chan string, query string, cfg config.SnipsFzfConfig) (string, error) {
	// TODO: handle channel skill issues.. is buffering the best here?
	output := make(chan string, 1)
	defer close(output)

	// Automatically select the only match, exit immediately when there's no match.
	opts := []string{
		"--select-1",
		"--exit-0",
		"--style", "full",
		"--scheme", "path",
		"--delimiter", unitSep,
		"--with-nth", "2",
		"--input-label", " Query ",
		"--preview", cfg.Preview,
	}
	if cfg.PreviewLabel != "" {
		opts = append(opts, "--bind", "focus:transform-preview-label:"+cfg.PreviewLabel)
	}
	if cfg.ListLabel != "" {
		opts = append(opts, "--bind", "result:transform-list-label:"+cfg.ListLabel)
	}
	if header != "" {
		opts = append(opts, "--header-first", "--header", header)
	}
	if query != "" {
		opts = append(opts, "--query", query)
	}
	options, err := fzf.ParseOptions(cfg.UseEnv, opts)
	if err != nil {
		return "", err
	}
	options.Input = input
	options.Output = output

	_, err = fzf.Run(options)
	if err != nil {
		return "", err
	}

	select {
	case res := <-output:
		if res == "" {
			return "", ErrNoSnippetFound
		}
		return res[0:strings.Index(res, unitSep)], nil
	default:
		return "", ErrNoFzfOutput
	}
}
