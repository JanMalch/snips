package utils

import gap "github.com/muesli/go-app-paths"

var paths *gap.Scope

func init() {
	paths = gap.NewScope(gap.User, "snips")
}

func DefaultSnippetDirectory() string {
	s, _ := paths.DataPath("snippets")
	return s
}

func RepositoriesFile() (string, error) {
	return paths.DataPath("repositories.yaml")
}
