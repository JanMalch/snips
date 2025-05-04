package repositories

import (
	"errors"
	"snips/internal/snippets"
)

type Repository interface {
	Name() string
	ListAll() (snippets.Ids, error)
	Read(id snippets.Id) (string, error)
}

var (
	ErrNoSuchSnippet = errors.New("failed to find snippet for given ID")
)
