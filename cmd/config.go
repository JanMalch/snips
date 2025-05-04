package cmd

import (
	"fmt"

	"snips/cmd/config"
	"snips/internal/cnfg"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configList bool

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Parent command for snips configuration",
	Long: `Parent command for snips configuration.

Examples:
  snips config --list`,
	Run: func(cmd *cobra.Command, args []string) {
		if configList {
			for _, k := range cnfg.KEYS {
				fmt.Printf("%s=%v\n", k, viper.Get(k))
			}
		} else {
			cobra.CheckErr("no action specified")
		}
	},
}

func init() {
	configCmd.AddCommand(config.GetCmd)
	configCmd.AddCommand(config.SetCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	configCmd.Flags().BoolVar(&configList, "list", false, "List all configured values")
}
