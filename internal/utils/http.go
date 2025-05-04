package utils

import (
	"errors"
	"io"
	"net/http"
)

var (
	ErrDownloadFailed = errors.New("failed to download snippet")
)

func FetchBody(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", ErrDownloadFailed
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
