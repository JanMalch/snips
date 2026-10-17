package cli

import (
	"github.com/JanMalch/snips/internal/config"

	"github.com/alecthomas/kong"
)

const VERSION = "0.7.1"

func Run() {
	var app CLI
	ctx := kong.Parse(&app,
		kong.Name("snips"),
		kong.Description(`Usage: snips [flags] [<snippet>] [-- <snippet-args> ...]
		
CLI to help with snippets and scripts.

Examples:
  snips foo				Searches a snippet by path, then prints the file's content.
  snips -c foo			Searches a snippet by path, then prints and copies the file's content.
  snips -x foo			Searches a snippet by path, then executes the selected command.
  snips -xc foo			Searches a snippet by path, then copies the selected command.
  snips -xp foo			Searches a snippet by path, then prints the selected command.
  snips -x foo -- hello		Searches a snippet by path, then executes the selected command, passing "hello" as an argument.
  snips -lS0 | xargs grep -r 'hello'	Use grep to search for 'hello' in all files of the first source directory.
`),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			NoAppSummary: true,
		}),
		kong.Vars{
			"version": VERSION,
		})

	if app.CheckForUpdates {
		ctx.FatalIfErrorf(CheckForUpdate(VERSION))
		return
	}

	cfg, err := config.Load()
	ctx.FatalIfErrorf(err)
	start(&app, ctx, cfg)
}
