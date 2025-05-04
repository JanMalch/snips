package repositories

import (
	"io/fs"
	"os"
	"path/filepath"
	"snips/internal/snippets"
)

type localDirRepository struct {
	name string
	dir  string
}

func Local(name, path string) Repository {
	return &localDirRepository{name: name, dir: path}
}

func (r *localDirRepository) String() string {
	return r.name + " @ " + r.dir
}

func (r *localDirRepository) Name() string {
	return r.name
}

func (r *localDirRepository) ListAll() (snippets.Ids, error) {
	result := make([]snippets.Id, 0)
	err := filepath.WalkDir(r.dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			id, err := snippets.IdFromPath(r.dir, path)
			if err != nil {
				return err
			}
			result = append(result, id)
		}
		return nil
	})
	if err != nil {
		return []snippets.Id{}, err
	}
	snippets.Sort(result)
	return result, nil
}

func (r *localDirRepository) Read(id snippets.Id) (string, error) {
	path := filepath.Join(r.dir, id.String())
	// abs mostly for better debugging. Shouldn't even be necessary
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrNoSuchSnippet
		}
		return "", err
	}
	return string(b), nil
}
