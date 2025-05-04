package snippets

import (
	"errors"
	"path/filepath"
	"slices"
	"strings"
)

type Id struct {
	Group string
	Name  string
}

var (
	ErrInvalidSnippetId = errors.New("invalid snippet ID")
)

func NewId(id string) (Id, error) {
	if id == "" {
		return Id{}, ErrInvalidSnippetId
	}
	id = filepath.ToSlash(id)
	parts := slices.DeleteFunc(strings.Split(id, "/"), func(s string) bool { return s == "" || s == "." })
	if len(parts) == 0 {
		return Id{}, ErrInvalidSnippetId
	}
	if len(parts) == 1 {
		return Id{
			Group: "",
			Name:  id,
		}, nil
	}
	return Id{
		Group: strings.Join(parts[0:len(parts)-1], "/"),
		Name:  parts[len(parts)-1],
	}, nil
}

func (s Id) GroupSegments() []string {
	if s.Group == "" {
		return []string{}
	}
	return strings.Split(s.Group, "/")
}

func (s Id) Breadcrumbs() []string {
	if s.Group == "" {
		return []string{}
	}
	segments := strings.Split(s.String(), "/")
	if len(segments) <= 1 {
		return []string{}
	}

	var breadcrumbs []string
	var current string

	for i := 0; i < len(segments)-1; i++ {
		if i == 0 {
			current = segments[i]
		} else {
			current = current + "/" + segments[i]
		}
		breadcrumbs = append(breadcrumbs, current)
	}

	return breadcrumbs
}

func (s Id) InverseBreadcrumbs() []string {
	if s.Group == "" {
		return []string{}
	}

	parts := strings.Split(s.Group, "/")
	n := len(parts)
	result := make([]string, n)

	for i := 0; i < n; i++ {
		result[i] = strings.Join(parts[n-1-i:], "/")
	}

	return result
}

func (s Id) String() string {
	if s.Group == "" {
		return s.Name
	}
	return s.Group + "/" + s.Name
}
