package snippets

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"snips/internal/cnfg"
	"strings"

	"github.com/spf13/viper"
)

var (
	ErrMalformedSnipsFile  = errors.New("malformed snips file")
	ErrBlankSnippetContent = errors.New("snippet content may not be blank")
	ErrNoSnippetDirectory  = errors.New("no snippet directory configured")
)

func IdFromPath(dir, file string) (Id, error) {
	rawID, err := filepath.Rel(dir, file)
	if err != nil {
		return Id{}, err
	}
	return NewId(filepath.Clean(rawID))
}

type Ids []Id

func (ids Ids) String(i int) string {
	return ids[i].String()
}

func (ids Ids) Len() int {
	return len(ids)
}

func Sort(ids Ids) {
	slices.SortFunc(ids, func(a, b Id) int {
		return cmp.Compare(a.String(), b.String())
	})
}

func ListAll() (Ids, error) {
	dir, err := Dir()
	if err != nil {
		return []Id{}, err
	}
	result := make([]Id, 0)
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			id, err := IdFromPath(dir, path)
			if err != nil {
				return err
			}
			result = append(result, id)
		}
		return nil
	})
	if err != nil {
		return []Id{}, err
	}
	slices.SortFunc(result, func(a, b Id) int {
		return cmp.Compare(a.String(), b.String())
	})
	return result, nil
}

func Dir() (string, error) {
	dir := viper.GetString(cnfg.KEY_SNIPPET_DIRECTORY)
	if dir == "" {
		return "", ErrNoSnippetDirectory
	}
	return dir, nil
}

func PathOf(id Id) (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, id.String()), nil
}

func Read(id Id) (string, error) {
	path, err := PathOf(id)
	if err != nil {
		return "", err
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func Write(id Id, content, source string) (string, error) {
	if source != "" {
		content = fmt.Sprintf("%s Source: %s %s\n%s", commentPrefix(id.Name), source, commentSuffix(id.Name), content)
	}
	path, err := PathOf(id)
	if err != nil {
		return path, err
	}
	err = os.MkdirAll(filepath.Dir(path), 0644)
	if err != nil {
		return path, err
	}
	return path, os.WriteFile(path, []byte(content), 0644)
}

func commentPrefix(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".go", ".java", ".js", ".ts", ".c", ".cpp", ".cs", ".swift", ".kt", ".scala", ".php", ".scss", ".sass":
		return "//"
	case ".py", ".rb", ".sh", ".bash", ".zsh":
		return "#"
	case ".sql":
		return "--"
	case ".html":
		return "<!--"
	case ".css":
		return "/*"
	default:
		return "//"
	}
}

func commentSuffix(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".html":
		return "-->"
	case ".css":
		return "*/"
	default:
		return ""
	}
}
