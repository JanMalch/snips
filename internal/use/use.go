package use

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"snips/internal/snippets"
	"snips/internal/utils"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func Use(
	idArg string,
	fileArg string,
	allowDirty bool,
	overwrite bool,
	copy bool,
	stdout io.Writer,
) error {
	id, err := snippets.NewId(idArg)
	if err != nil {
		return err
	}
	content, err := snippets.Read(id)
	if err != nil {
		return err
	}

	if fileArg != "" {
		if !allowDirty {
			if utils.IsDirtyGitRepo() {
				cobra.CheckErr("git repository is dirty")
			}
		}
		outFile := fileArg
		if fileArg == "." {
			outFile, err = promptFileOutput(id)
			cobra.CheckErr(err)
		}
		cobra.CheckErr(applySnippetToFile(content, outFile, overwrite))
	} else if copy {
		cobra.CheckErr(clipboard.WriteAll(content))
	} else {
		fmt.Fprint(stdout, content)
	}
	return nil
}

func promptFileOutput(id snippets.Id) (string, error) {
	namePath, err := filepath.Localize(id.Name)
	if err != nil {
		return "", err
	}

	options := make([]huh.Option[string], 0)
	options = append(options, huh.NewOption(namePath, namePath))
	for _, b := range id.InverseBreadcrumbs() {
		idPath, err := filepath.Localize(b + "/" + id.Name)
		if err != nil {
			return "", err
		}
		options = append(options, huh.NewOption(idPath, idPath))
	}

	var path string
	form := huh.NewSelect[string]().
		Title("Suggested output paths").
		Options(options...).
		Value(&path)
	if err := form.Run(); err != nil {
		return "", err
	}
	return path, nil
}

func applySnippetToFile(content, path string, overwrite bool) error {
	err := os.MkdirAll(filepath.Dir(path), 0644)
	if err != nil {
		return err
	}
	if overwrite {
		return os.WriteFile(path, []byte(content), 0644)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
