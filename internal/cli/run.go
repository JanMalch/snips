package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"snips/internal/config"
	"snips/internal/core"
	"snips/internal/exe"

	"charm.land/huh/v2"
	"github.com/alecthomas/kong"
	"github.com/atotto/clipboard"
)

var (
	ErrNoEditor = errors.New("no EDITOR env var defined")
)

func Run(cli *CLI, ctx *kong.Context, cfg config.SnipsConfig) {
	if cli.Config {
		path, err := config.Path()
		ctx.FatalIfErrorf(err)
		if cli.Locate {
			fmt.Fprintln(ctx.Stdout, path)
			return
		}
		if cli.Edit {
			ctx.FatalIfErrorf(edit(path))
			return
		}
		dat, err := os.ReadFile(path)
		ctx.FatalIfErrorf(err)
		fmt.Fprintln(ctx.Stdout, string(dat))
		return
	}

	dirs := cli.Sources(cfg)

	print := !cli.Exec
	if cli.Print != nil {
		print = *cli.Print
	}

	snippet, err := core.FindSnippet(cli.Snippet, dirs, cfg.IncludeSourceName, cfg.Fzf)
	ctx.FatalIfErrorf(err)

	if cli.Locate {
		fmt.Fprintln(ctx.Stdout, snippet)
		return
	}
	if cli.Edit {
		ctx.FatalIfErrorf(edit(snippet))
		return
	}

	if !cli.Exec {
		dat, err := os.ReadFile(snippet)
		ctx.FatalIfErrorf(err)
		if cli.Copy {
			ctx.FatalIfErrorf(clipboard.WriteAll(string(dat)))
		}
		if print {
			fmt.Fprintln(ctx.Stdout, string(dat))
		}
		return
	}

	cmds := exe.DetermineCmds(snippet, cfg.Runners)
	if len(cmds) == 0 {
		ctx.Fatalf("Failed to determine any appropriate command for %s", snippet)
	}

	cmdIdx := -1
	options := make([]huh.Option[int], len(cmds))
	for i, c := range cmds {
		options[i] = huh.NewOption(c.String(), i)
	}
	actionTitle := "run"
	if cli.Copy {
		actionTitle = "copy"
	} else if print {
		actionTitle = "print"
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title(fmt.Sprintf("Pick a command to %s.", actionTitle)).
				Options(options...).
				Value(&cmdIdx)),
	).WithAccessible(os.Getenv("ACCESSIBLE") != "")

	ctx.FatalIfErrorf(form.Run())

	cmd := cmds[cmdIdx]

	if print || cli.Copy {
		if print {
			fmt.Fprintln(ctx.Stdout, cmd.String())
		}
		if cli.Copy {
			ctx.FatalIfErrorf(clipboard.WriteAll(cmd.String()))
		}
	} else {
		ctx.FatalIfErrorf(cmd.Run())
	}
}

func edit(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return ErrNoEditor
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
