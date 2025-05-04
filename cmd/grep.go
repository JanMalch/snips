package cmd

import (
	"os"
	"regexp"
	"snips/internal/grep"
	"snips/internal/snippets"

	"github.com/spf13/cobra"
)

var grepCmdRegex bool
var grepCmdIgnoreCase bool

// grepCmd represents the grep command
var grepCmd = &cobra.Command{
	Use:   "grep [flags] query",
	Args:  cobra.ExactArgs(1),
	Short: "Searches through the contents of all added snippets",
	Long: `Searches through the contents of all added snippets.

Examples:
  snips grep -i pmap`,
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		dir, err := snippets.Dir()
		cobra.CheckErr(err)
		var opts *grep.SearchOptions
		debug := &grep.SearchDebug{
			Workers: 128,
		}
		if grepCmdRegex {
			r, err := regexp.Compile(query)
			cobra.CheckErr(err)
			opts = &grep.SearchOptions{
				Kind:   grep.REGEX,
				Lines:  false,
				Regex:  r,
				Finder: nil,
			}
		} else {
			opts = &grep.SearchOptions{
				Kind:   grep.LITERAL,
				Lines:  false,
				Regex:  nil,
				Finder: grep.MakeStringFinder(query, grepCmdIgnoreCase),
			}
		}
		grep.Search(dir, opts, debug, os.Stdout)
	},
}

func init() {
	grepCmd.Flags().BoolVarP(&grepCmdRegex, "regexp", "e", false, "treat query as a regex")
	grepCmd.Flags().BoolVarP(&grepCmdIgnoreCase, "ignore-case", "i", false, "ignore case distinctions")
	grepCmd.MarkFlagsMutuallyExclusive("regexp", "ignore-case")
}
