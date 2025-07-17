// Package cmd represents the cli inteface of muxocil
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/boazj/muxocil/common"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const configDescription = `config file path order of precedence:
1. --config
2. ./.muxocil.toml
3. OS based config path
  3.1 Windows  - %AppData%\.muxocil.toml
  3.2 Apple - $HOME/Library/Application Support/muxocil/.muxocil.toml
  3.3 Linux/Unix - Same as #4 & #5
4. $XDG_CONFIG_HOME/muxocil/.muxocil.toml
5. $HOME/.config/muxocil/.muxocil.toml
6. $HOME/.muxocil.toml
If none of the above are used, then Muxocil will use the default config`

var rootCmd = &cobra.Command{
	Use:     "muxocil",
	Short:   "Create windows and panes layouts in your multiplexer of choice",
	Version: AppVersion,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("General error occurred", "error", err)
		os.Exit(common.ExitGeneralError)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("config", "", "", configDescription)
	err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	if err != nil {
		panic(fmt.Sprintf("failed to bind config: %v", err))
	}

	rootCmd.PersistentFlags().BoolP("verbose", "", false, "Show verbose logging")
	err = viper.BindPFlag("log.verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	if err != nil {
		panic(fmt.Sprintf("failed to bind config: %v", err))
	}

	rootCmd.PersistentFlags().BoolP("debug", "", false, "Show debug logging")
	err = viper.BindPFlag("log.debug", rootCmd.PersistentFlags().Lookup("debug"))
	if err != nil {
		panic(fmt.Sprintf("failed to bind config: %v", err))
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.SetLevel(log.WarnLevel)
	if viper.GetBool("log.verbose") {
		log.SetLevel(log.InfoLevel)
	}
	if viper.GetBool("log.debug") {
		log.SetLevel(log.DebugLevel)
	}

	cfgFile, err := rootCmd.Flags().GetString("config")
	if err != nil {
		panic(fmt.Sprintf("failed to get config: %v", err))
	}
	if cfgFile != "" {
		// Use config file from the flag.
		log.Debug("Using override config file", "path", cfgFile)
		viper.SetConfigFile(cfgFile)
	} else {
		log.Debug("Searching config file in default locations")
		configPath, _ := os.UserConfigDir()

		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(configPath, "muxocil"))
		viper.AddConfigPath("$XDG_CONFIG_HOME/muxocil")
		viper.AddConfigPath(filepath.Join(home, ".config", "muxocil"))
		viper.AddConfigPath(home)

		viper.SetConfigName(".muxocil")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read config file", "err", err)
		os.Exit(common.ExitConfigFailure)
	}
	log.Debug("Loaded config file from", "path", viper.ConfigFileUsed())
}
