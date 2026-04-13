package core

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func lsRecursively(dir string) ([]string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0)
	err = filepath.WalkDir(abs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			res = append(res, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Based on https://github.com/yargevad/filepathx/blob/d026a3f0c9b7bd7c107f6c0260996ef00688a9ca/filepathx.go
func Doublestar(dir, pattern string) ([]string, error) {
	if pattern == "**/*" || pattern == "**" {
		return lsRecursively(dir)
	}
	pattern = filepath.Join(dir, pattern)
	if !strings.Contains(pattern, "**") {
		return filepath.Glob(pattern)
	}
	globs := strings.Split(pattern, "**")
	matches := []string{""}
	for _, glob := range globs {
		hits := make([]string, 0)
		hitMap := make(map[string]bool)
		for _, match := range matches {
			paths, err := filepath.Glob(match + glob)
			if err != nil {
				return nil, err
			}
			for _, path := range paths {
				err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					// TODO: what about directories here?!
					// save deduped match from current iteration
					if _, ok := hitMap[path]; !ok {
						hits = append(hits, path)
						hitMap[path] = true
					}
					return nil
				})
				if err != nil {
					return nil, err
				}
			}
		}
		matches = hits
	}

	// fix up return value for nil input
	if globs == nil && len(matches) > 0 && matches[0] == "" {
		return matches[1:], nil
	}
	return matches, nil
}
