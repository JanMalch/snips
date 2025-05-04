package repository

import (
	"snips/internal/repositories"

	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:     "add name source",
	Aliases: []string{"a"},
	Args:    cobra.ExactArgs(2),
	Short:   "Adds the source as a repository under the given name",
	Long: `Adds the source as a repository under the given name.

Examples:
  snips repository add snipssrc 'https://github.com/JanMalch/snips'`,
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(repositories.AddStored(args[0], args[1]))
	},
}
