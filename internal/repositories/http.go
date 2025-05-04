package repositories

import (
	"errors"
	"net/url"
	"snips/internal/snippets"
	"snips/internal/utils"
	"strings"
)

var (
	ErrInvalidUrl     = errors.New("invalid HTTP URL")
	ErrListUnspported = errors.New("cannot list HTTP repositories")
)

type httpRepository struct {
	name string
	url  string
}

func Http(name, rawURL string) (Repository, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, ErrInvalidUrl
	}
	if !strings.HasSuffix(rawURL, "/") {
		rawURL = rawURL + "/"
	}
	return &httpRepository{name: name, url: rawURL}, nil
}

func (r *httpRepository) String() string {
	return r.name + " @ " + r.url
}

func (r *httpRepository) Name() string {
	return r.name
}

func (r *httpRepository) ListAll() (snippets.Ids, error) {
	return nil, ErrListUnspported
}

func (r *httpRepository) Read(id snippets.Id) (string, error) {
	content, err := utils.FetchBody(r.url + id.String())
	if err != nil {
		return "", ErrNoSuchSnippet
	}
	return content, nil
}
