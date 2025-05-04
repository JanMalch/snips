package add

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"snips/internal/snippets"
	"snips/internal/utils"
	"strings"

	"github.com/charmbracelet/huh"
)

type confirmOutcome int

const (
	confirmYes confirmOutcome = iota
	confirmNo
	confirmEdit
)

type argType int

const (
	argInvalid argType = iota
	argFile
	argHttps
)

var (
	ErrNoContent           = errors.New("got no arguments or stdin")
	ErrNoNameForNew        = errors.New("cannot determine name for snippet")
	ErrInvalidNewArg       = errors.New("invalid argument")
	ErrEmptySnippet        = errors.New("downloaded empty snippet")
	ErrNameButMultipleArgs = errors.New("cannot use \"name\" option when adding multiple snippets")
)

func Add(
	args []string,
	stdinContent string,
	name string,
	group string,
	confirmContent bool,
	yesToAll bool,
	stdout io.Writer,
) error {
	if len(args) == 0 {
		if stdinContent == "" {
			return ErrNoContent
		}
		id, err := determineId(group, name, "")
		if err != nil {
			return err
		}
		// confirming doesn't work when using stdin ...?
		if confirmContent {
			fmt.Fprintln(stdout, "cannot confirm content when using stdin content")
		}
		_, err = Create(id, stdinContent, "", false, yesToAll, stdout)
		return err
	}
	if len(args) > 1 && name != "" {
		return ErrNameButMultipleArgs
	}
	for _, arg := range args {
		if arg == "" {
			return ErrInvalidNewArg
		}
		err := handleArg(group, name, arg, confirmContent, yesToAll, stdout)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleArg(group, name, arg string, confirmContent, yesToAll bool, stdout io.Writer) error {
	argType, err := determineArgType(arg)
	if err != nil {
		return err
	}
	id, err := determineId(group, name, arg)
	if err != nil {
		return err
	}

	content := ""
	source := arg
	switch argType {
	case argInvalid:
		return ErrInvalidNewArg
	case argFile:
		fileContent, err := os.ReadFile(arg)
		if err != nil {
			return err
		}
		if abs, err := filepath.Abs(arg); err == nil {
			source = abs
		}
		content = string(fileContent)
	case argHttps:
		httpContent, err := utils.FetchBody(arg)
		if err != nil {
			return err
		}
		content = httpContent
	}

	_, err = Create(id, content, source, confirmContent, yesToAll, stdout)
	return err
}

func Create(id snippets.Id, content, source string, confirmContent, yesToAll bool, stdout io.Writer) (bool, error) {
	if strings.TrimSpace(content) == "" {
		return false, ErrEmptySnippet
	}

	var err error
	outcome := confirmYes
	if !yesToAll && confirmContent {
		outcome, err = promptConfirmContent(content, stdout)
		if err != nil {
			return false, err
		}
		if outcome == confirmNo {
			return false, nil
		}
	}
	snippetPath, err := snippets.Write(id, content, source)
	if err != nil {
		return false, err
	}
	fmt.Fprintf(stdout, "snippet created successfully: %s\n", id)
	if outcome == confirmEdit {
		return true, utils.OpenEditorWithViper(snippetPath)
	}
	return true, nil
}

func determineArgType(arg string) (argType, error) {
	if arg == "" {
		return argInvalid, ErrInvalidNewArg
	}
	if strings.HasPrefix(arg, "https://") || strings.HasPrefix(arg, "http://") {
		_, err := url.Parse(arg)
		if err != nil {
			return argInvalid, ErrInvalidNewArg
		}
		return argHttps, nil
	}
	return argFile, nil
}

func promptConfirmContent(content string, stdout io.Writer) (confirmOutcome, error) {
	// TODO: fancy box & colors
	fmt.Fprintln(stdout, content+"\n")
	outcome := confirmYes
	form := huh.NewSelect[confirmOutcome]().
		Title("Install this snippet?").
		Options(
			huh.NewOption("Yes", confirmYes),
			huh.NewOption("Yes, and edit afterwards", confirmEdit),
			huh.NewOption("No", confirmNo),
		).
		Value(&outcome)
	if err := form.Run(); err != nil {
		return confirmNo, err
	}
	return outcome, nil
}

func determineId(group, name, arg string) (snippets.Id, error) {
	parts := []string{}
	if name == "" {
		if arg == "" {
			return snippets.Id{}, ErrNoNameForNew
		}
		urlFile := arg[strings.LastIndex(arg, "/")+1:]
		parts = append(parts, urlFile)
	} else {
		parts = append(parts, strings.Split(name, "/")...)
	}
	if group != "" {
		parts = append(strings.Split(group, "/"), parts...)
	}
	joined := strings.Join(parts, "/")
	return snippets.NewId(filepath.ToSlash(filepath.Clean(joined)))
}
