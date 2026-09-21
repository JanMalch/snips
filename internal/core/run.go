package core

import (
	"errors"
	"fmt"
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

func FindSnippet(query string, dirs []string, includeSourceName bool, cfg config.SnipsFzfConfig, grep bool) (string, error) {
	if len(dirs) == 0 {
		return "", ErrNoSources
	}

	matches := make([]string, 0)
	matchedDirs := make([]string, 0)
	useDirPrefix := includeSourceName && len(dirs) > 1

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
				if useDirPrefix {
					rels = filepath.Join(filepath.Base(dir), rels)
				}
				matches = append(matches, s+unitSep+rels)
				matchedDirs = append(matchedDirs, s)
			}
		}
	}
	if len(matches) == 0 {
		return "", ErrNoMatches
	}

	if grep {
		return runFzfWithRipgrep("", matchedDirs, query, cfg)
	}
	inputChan := make(chan string)
	go func() {
		for _, m := range matches {
			inputChan <- m
		}
		close(inputChan)
	}()
	// TODO: use header
	return runFzfByConfig("", inputChan, query, cfg)
}

func runFzfWithRipgrep(header string, directories []string, query string, cfg config.SnipsFzfConfig) (string, error) {
	qdir := make([]string, 0)
	for _, d := range directories {
		qdir = append(qdir, "'"+d+"'")
	}
	rgcmd := fmt.Sprintf("rg -l --smart-case {q} %s || echo %s", strings.Join(qdir, " "), unitSep)
	// TODO: use -i if all lowercase
	// FIXME: grep not working when no results
	// rgcmd := fmt.Sprintf("grep -l -i {q} %s || echo %s", strings.Join(qdir, " "), unitSep)
	fmt.Println(rgcmd)
	opts := []string{
		// https://junegunn.github.io/fzf/tips/ripgrep-integration/
		"--disabled",
		"--bind", "start:reload:" + rgcmd,
		"--bind", "change:reload:" + rgcmd,
		// Automatically select the only match, exit immediately when there's no match.
		"--select-1",
		"--exit-0",
		"--style", "full",
		"--delimiter", unitSep,
		"--input-label", " Ripgrep Query ",
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
	res, err := runFzf(opts, nil, cfg.UseEnv)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(res), nil
}

func runFzfByConfig(header string, input chan string, query string, cfg config.SnipsFzfConfig) (string, error) {
	opts := []string{
		// Automatically select the only match, exit immediately when there's no match.
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
	res, err := runFzf(opts, input, cfg.UseEnv)
	if err != nil {
		return "", err
	}
	return res[0:strings.Index(res, unitSep)], nil
}

func runFzf(opts []string, input chan string, useEnv bool) (string, error) {
	// TODO: handle channel skill issues.. is buffering the best here?
	output := make(chan string, 1)
	defer close(output)

	options, err := fzf.ParseOptions(useEnv, opts)
	if err != nil {
		return "", fmt.Errorf("invalid fzf options: %w", err)
	}
	if input != nil {
		options.Input = input
	}
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
		return res, nil
	default:
		return "", ErrNoFzfOutput
	}
}
