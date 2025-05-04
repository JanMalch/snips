package cmd

import (
	"errors"
	"fmt"
	"os"
	"snips/internal/add"
	"snips/internal/repositories"
	"snips/internal/snippets"

	"github.com/spf13/cobra"
)

var installCmdRepoQuery string

var installCmd = &cobra.Command{
	Use:     "install [flags] ids...",
	Aliases: []string{"i"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Adds specified snippet IDs from registered repositories",
	Long:    "Adds specified snippet IDs from registered repositories.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cobra.CheckErr("provided no IDs to install")
		}
		repos, err := repositories.ReadStored()
		cobra.CheckErr(err)
		if len(repos) == 0 {
			cobra.CheckErr("no repositories added")
		}
		if installCmdRepoQuery != "" {
			repo, err := repositories.FindStored(installCmdRepoQuery)
			cobra.CheckErr(err)
			repos = []repositories.Repository{repo}
		}
		installed := make([]snippets.Id, 0)
		notInstalled := make([]snippets.Id, 0)
		for _, arg := range args {
			ok := false
			id, err := snippets.NewId(arg)
			cobra.CheckErr(err)
			for _, repo := range repos {
				content, err := repo.Read(id)
				if err != nil {
					if errors.Is(err, repositories.ErrNoSuchSnippet) {
						continue
					} else {
						cobra.CheckErr(err)
					}
				}
				id, err = snippets.NewId(repo.Name() + "/" + id.String())
				cobra.CheckErr(err)
				ok, err = add.Create(id, content, repo.Name(), true, false, os.Stdout)
				cobra.CheckErr(err)
				if ok {
					installed = append(installed, id)
				}
				break
			}
			if !ok {
				notInstalled = append(notInstalled, id)
			}
		}
		fmt.Printf("%d / %d installed: %s\n", len(installed), len(args), installed)
		fmt.Printf("%d / %d not installed: %s\n", len(notInstalled), len(args), notInstalled)
	},
}

func init() {
	installCmd.Flags().StringVarP(&installCmdRepoQuery, "repository", "r", "", "Search in specified repository (name or index)")
}
