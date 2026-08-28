package main

import (
	"snips/internal/cli"
	"snips/internal/config"

	"github.com/alecthomas/kong"
)

const VERSION = "0.4.0"

func main() {
	var app cli.CLI
	ctx := kong.Parse(&app,
		kong.Name("snips"),
		kong.Description(`Usage: snips [flags] [<snippet>] [-- <snippet-args> ...]
		
CLI to help with snippets and scripts.

Examples:
  snips foo			Searches a snippet by path, then prints the file's content.
  snips -c foo		Searches a snippet by path, then prints and copies the file's content.
  snips -x foo		Searches a snippet by path, then executes the selected command.
  snips -xc foo		Searches a snippet by path, then copies the selected command.
  snips -xp foo		Searches a snippet by path, then prints the selected command.
  snips -x foo -- hello	Searches a snippet by path, then executes the selected command, passing "hello" as an argument.
`),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			NoAppSummary: true,
		}),
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
