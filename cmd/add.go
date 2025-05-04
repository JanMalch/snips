package cmd

import (
	"os"
	"snips/internal/add"
	"snips/internal/cnfg"
	"snips/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addCmdName string
var addCmdGroup string
var addCmdYesToAll bool

var addCmd = &cobra.Command{
	Use:     "add [flags] [sources...]",
	Aliases: []string{"a"},
	Short:   "Adds snippets from the given sources",
	Long: `Adds snippets from the given sources.

Examples:
  snips add ./utils/strings.ts
  snips add 'https://raw.githubusercontent.com/sindresorhus/github-markdown-css/refs/heads/main/github-markdown.css' -n 'css/github-markdown.css'
  cat ./utils.kt | snips add -g kotlin -n utils
`,
	Run: func(cmd *cobra.Command, args []string) {
		stdinContent, _ := utils.GetStdin()
		cobra.CheckErr(add.Add(
			args,
			stdinContent,
			addCmdName,
			addCmdGroup,
			viper.GetBool(cnfg.KEY_NEW_CONFIRM_CONTENT),
			addCmdYesToAll,
			os.Stdout,
		))
	},
}

func init() {
	addCmd.Flags().StringVarP(&addCmdName, "name", "n", "", "name for the snippet (defaults to last path segment)")
	addCmd.Flags().StringVarP(&addCmdGroup, "group", "g", "", "group for the snippet (can also be defined with slashes in the name)")
	addCmd.Flags().BoolP("confirm", "c", true, "confirm content before installing")
	viper.BindPFlag(cnfg.KEY_NEW_CONFIRM_CONTENT, addCmd.Flags().Lookup("confirm"))
	addCmd.Flags().BoolVarP(&addCmdYesToAll, "yes", "y", false, "skip confirmation with \"yes\"")
}
