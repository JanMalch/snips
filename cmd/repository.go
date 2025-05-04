package cmd

import (
	"fmt"
	"snips/cmd/repository"
	"snips/internal/repositories"

	"github.com/spf13/cobra"
)

var repositoryCmdList bool

var repositoryCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"r"},
	Short:   "Parent command for repositories",
	Long: `Parent command for repositories.

Examples:
  snips repository --list`,
	Run: func(cmd *cobra.Command, args []string) {
		if repositoryCmdList {
			repos, err := repositories.ReadStored()
			cobra.CheckErr(err)
			for i, r := range repos {
				fmt.Printf("[%d] %s\n", i, r)
			}
		} else {
			cobra.CheckErr("no action specified")
		}
	},
}

func init() {
	repositoryCmd.AddCommand(repository.AddCmd)

	repositoryCmd.Flags().BoolVar(&repositoryCmdList, "list", false, "List all repositories")
}
