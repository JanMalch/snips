package cmd

import (
	"fmt"
	"snips/internal/snippets"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"l"},
	Short:   "Lists all installed snippets",
	Long:    "Lists all installed snippets.",
	Run: func(cmd *cobra.Command, args []string) {
		ids, err := snippets.ListAll()
		cobra.CheckErr(err)
		for _, id := range ids {
			fmt.Println(id.String())
		}
	},
}
