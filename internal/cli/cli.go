package cli

import (
	"fmt"
	"snips/internal/config"

	"github.com/alecthomas/kong"
)

// https://danielms.site/zet/2024/how-i-write-golang-cli-tools-today-using-kong/
type versionFlag string

func (v versionFlag) Decode(_ *kong.DecodeContext) error { return nil }
func (v versionFlag) IsBool() bool                       { return true }
func (v versionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

// TODO: https://github.com/alecthomas/kong?tab=readme-ov-file#configurationloader-paths---load-defaults-from-configuration-files ?

type CLI struct {
	Snippet string `arg:"" optional:"" name:"snippet" help:"Optional initial query for the snippet path, or content when using --grep/-g. The query is optional. When only one snippet matches, it is selected automatically."`
	Exec    bool   `name:"exec" default:"false" short:"x" help:"Execute the selected snippet after confirmation. Defaults to false."`
	// TODO: use env default?
	Copy  bool  `name:"copy" default:"false" short:"c" help:"Copies the selected snippet to the system clipboard. In --exec mode it copies the command instead of executing it. Defaults to false."`
	Print *bool `name:"print" short:"p" negatable:"" help:"Prints the selected snippet, and defaults to true. In --exec mode, it prints the command instead of executing it, and defaults to false."`
	// TODO: Locate & Here: improve names and shorts?
	Locate bool `name:"locate" default:"false" short:"l" help:"Only print the full absolute path of the selected snippet before exiting."`
	Here   bool `name:"here" default:"false" short:"w" help:"Use the current working directory as the only source. Defaults to false."`
	// TODO
	// Grep    bool   `name:"grep" default:"false" short:"g" help:"Grep on snippet contents instead. Requires ripgrep, grep or git."`
	// Typ     string `name:"type" short:"t" help:"Filter by a ripgrep file type. Requires ripgrep."`
	Edit   bool `name:"edit" default:"false" short:"e" help:"Open the selected snippet with the editor defined by the EDITOR environment variable. Defaults to false."`
	Config bool `name:"config" default:"false" help:"Print snips config. Works with --locate/-l and --edit/-e."`
	// Color never,always,auto
	Source *int `name:"source" short:"s" help:"Select a source by index from your global snips config file. You can also use -0 to -9."`
	// TODO: there must be a better way, right?
	Source0         bool        `short:"0" default:"false" hidden:""`
	Source1         bool        `short:"1" default:"false" hidden:""`
	Source2         bool        `short:"2" default:"false" hidden:""`
	Source3         bool        `short:"3" default:"false" hidden:""`
	Source4         bool        `short:"4" default:"false" hidden:""`
	Source5         bool        `short:"5" default:"false" hidden:""`
	Source6         bool        `short:"6" default:"false" hidden:""`
	Source7         bool        `short:"7" default:"false" hidden:""`
	Source8         bool        `short:"8" default:"false" hidden:""`
	Source9         bool        `short:"9" default:"false" hidden:""`
	Version         versionFlag `name:"version" short:"v" help:"Print version information and quit"`
	CheckForUpdates bool        `name:"updates" aliases:"up" default:"false" help:"Checks for updates for the snips CLI. Same as --up."`
}

func (c *CLI) Sources(cfg config.SnipsConfig) []string {
	if c.Here {
		return []string{"."}
	}

	dirs := cfg.Sources
	if c.Source != nil {
		return []string{dirs[*c.Source]}
	}
	if c.Source0 {
		return dirs[0:1]
	}
	if c.Source1 {
		return dirs[1:2]
	}
	if c.Source2 {
		return dirs[2:3]
	}
	if c.Source3 {
		return dirs[3:4]
	}
	if c.Source4 {
		return dirs[4:5]
	}
	if c.Source5 {
		return dirs[5:6]
	}
	if c.Source6 {
		return dirs[6:7]
	}
	if c.Source7 {
		return dirs[7:8]
	}
	if c.Source8 {
		return dirs[8:9]
	}
	if c.Source9 {
		return dirs[9:10]
	}

	return dirs
}
