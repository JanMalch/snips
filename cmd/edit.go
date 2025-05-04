package cmd

import (
	"os"
	"snips/internal/snippets"
	"snips/internal/utils"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:     "edit id",
	Aliases: []string{"e"},
	Args:    cobra.ExactArgs(1),
	Short:   "Edit a snippet",
	Long: `Edit a snippet.

Examples:
  snips edit foo/bar.java`,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := snippets.NewId(args[0])
		cobra.CheckErr(err)
		path, err := snippets.PathOf(id)
		cobra.CheckErr(err)
		_, err = os.Stat(path)
		cobra.CheckErr(err)
		cobra.CheckErr(utils.OpenEditorWithViper(path))
	},
}
