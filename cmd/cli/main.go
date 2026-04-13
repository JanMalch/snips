package main

import (
	"snips/internal/cli"
	"snips/internal/config"

	"github.com/alecthomas/kong"
)

const VERSION = "0.2.0"

func main() {
	var app cli.CLI
	ctx := kong.Parse(&app,
		kong.Name("snips"),
		kong.Description(`CLI to manage and run scripts of any language.

Examples:
  snips foo		Searches a snippet by path, then prints the file's content.
  snips -c foo	Searches a snippet by path, then prints and copies the file's content.
  snips -x foo	Searches a snippet by path, then executes the selected command.
  snips -xc foo	Searches a snippet by path, then copies the selected command.
  snips -xp foo	Searches a snippet by path, then prints the selected command.
`),
		kong.UsageOnError(),
		kong.Vars{
			"version": VERSION,
		})

	if app.CheckForUpdates {
		ctx.FatalIfErrorf(cli.CheckForUpdate(VERSION))
		return
	}

	cfg, err := config.Load()
	ctx.FatalIfErrorf(err)
	cli.Run(&app, ctx, cfg)
}
