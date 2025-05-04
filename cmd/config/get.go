package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// SetCmd represents the set command
var GetCmd = &cobra.Command{
	Use:   "get key",
	Args:  cobra.ExactArgs(1),
	Short: "Gets the configured value for the given key",
	Long: `Gets the configured value for the given key.

Examples:
  snips config get snippets.directory`,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := viper.Get(key)
		if value == nil {
			fmt.Print()
		} else {
			fmt.Printf("%v\n", value)
		}
	},
}
