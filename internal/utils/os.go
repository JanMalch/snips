package utils

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"snips/internal/cnfg"

	"github.com/spf13/viper"
)

var (
	ErrReadStdinFailed = errors.New("failed to read from stdin")
)

func GetStdin() (string, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return "", nil
}

func OpenEditorWithViper(path string) error {
	editor := viper.GetString(cnfg.KEY_EDITOR)
	editorCmd := exec.Command(editor, path)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr
	return editorCmd.Run()
}
