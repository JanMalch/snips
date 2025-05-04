package config

import (
	"errors"
	"slices"
	"snips/internal/cnfg"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrUnknownConfigKey = errors.New("unknown config key")
)

// SetCmd represents the set command
var SetCmd = &cobra.Command{
	Use:   "set key value",
	Args:  cobra.ExactArgs(2),
	Short: "Sets the given key value pair",
	Long: `Sets the given key value pair.

Examples:
  snips config set apply.allow_dirty_git true`,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]
		if !slices.Contains(cnfg.KEYS, key) {
			cobra.CheckErr(ErrUnknownConfigKey)
		}
		if value == "true" {
			viper.Set(key, true)
		} else if value == "false" {
			viper.Set(key, false)
		} else {
			viper.Set(key, value)
		}
		err := viper.WriteConfig()
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			cobra.CheckErr(viper.SafeWriteConfig())
		} else {
			cobra.CheckErr(err)
		}
	},
}
