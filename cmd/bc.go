package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/boazj/muxocil/common"
	"github.com/boazj/muxocil/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	TEAMOCIL  = "$HOME/.teamocil"
	ITERMOCIL = "$HOME/.itermocil"
)

var bcCmd = &cobra.Command{
	Use:   "bc",
	Short: "Command compatible with iTermocil and teamocil, can be used in an alias to replace both tools",
	Run: func(cmd *cobra.Command, args []string) {
		if path := viper.GetString("bc.layout"); path != "" {
			if viper.GetBool("bc.edit") {
				cmd := exec.Command(os.ExpandEnv("${EDITOR:-vi}"), path)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
				os.Exit(0)
				return
			}
			// TODO:
			// show := viper.GetBool("bc.show")
			// TODO:
			// here := viper.GetBool("bc.here")
			return
		}

		if viper.GetBool("bc.list") {
			cfg := common.GetConfig()

			cfg.LayoutSearchLocations = append(cfg.LayoutSearchLocations, TEAMOCIL, ITERMOCIL)

			for _, loc := range cfg.GetLayoutSearchLocations() {
				layouts := utils.GetFilesRecursively(loc, func(path string) bool {
					return utils.IsYaml(path)
				})
				for _, l := range layouts {
					fmt.Printf("%s\n", utils.GetFileNamePart(l))
				}
			}

			return
		}
	},
}

func init() {
	bcCmd.Flags().String("layout", "", "Takes a custom file path to a YAML layout file instead of [layout-name]")
	bcCmd.Flags().Bool("here", false, "Uses the current window as the layout’s first window")
	bcCmd.Flags().Bool("edit", false, "Edit the layout file in either $EDITOR or your preferred GUI editor")
	bcCmd.Flags().Bool("show", false, "Shows the layout content instead of executing it")
	bcCmd.Flags().Bool("list", false, "Lists all available layouts in ~/.itermocil, ~/.teamocil & locations per configuration file")

	if err := viper.BindPFlag("bc.layout", bcCmd.Flags().Lookup("layout")); err != nil {
		panic(fmt.Sprintf("failed to bind config: %v", err))
	}
	if err := viper.BindPFlag("bc.here", bcCmd.Flags().Lookup("here")); err != nil {
		panic(fmt.Sprintf("failed to bind config: %v", err))
	}
	if err := viper.BindPFlag("bc.edit", bcCmd.Flags().Lookup("edit")); err != nil {
		panic(fmt.Sprintf("failed to bind config: %v", err))
	}
	if err := viper.BindPFlag("bc.show", bcCmd.Flags().Lookup("show")); err != nil {
		panic(fmt.Sprintf("failed to bind config: %v", err))
	}
	if err := viper.BindPFlag("bc.list", bcCmd.Flags().Lookup("list")); err != nil {
		panic(fmt.Sprintf("failed to bind config: %v", err))
	}
}

func init() {
	rootCmd.AddCommand(bcCmd)
}
