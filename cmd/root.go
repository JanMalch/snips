package cmd

import (
	"os"
	"runtime"
	"snips/internal/cnfg"
	"snips/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const VERSION = "0.1.0"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "snips",
	Short:   "Snips is a CLI to manage code snippets.",
	Version: VERSION,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// The suggested approach is for the parent command to use AddCommand
	// to add its most immediate subcommands.
	rootCmd.AddCommand(repositoryCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(grepCmd)
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.snips.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault(cnfg.KEY_APPLY_ALLOW_DIRTY, false)
	viper.SetDefault(cnfg.KEY_NEW_CONFIRM_CONTENT, true)
	viper.SetDefault(cnfg.KEY_SNIPPET_DIRECTORY, utils.DefaultSnippetDirectory())
	viper.SetDefault(cnfg.KEY_EDITOR, getDefaultEditor())

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".snips" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".snips")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func getDefaultEditor() string {
	if utils.CheckIfExecInPath("code") {
		return "code"
	}
	if utils.CheckIfExecInPath("vim") {
		return "vim"
	}
	if utils.CheckIfExecInPath("vi") {
		return "vi"
	}
	if runtime.GOOS == "windows" {
		return "start"
	}
	return "open"
}
