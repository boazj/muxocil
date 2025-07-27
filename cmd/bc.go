package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/boazj/muxocil/common"
	"github.com/boazj/muxocil/utils"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tiendc/gofn"
)

const (
	TEAMOCIL  = "$HOME/.teamocil"
	ITERMOCIL = "$HOME/.itermocil"
)

var bcCmd = &cobra.Command{
	Use:   "bc",
	Short: "Command compatible with iTermocil and teamocil, can be used in an alias to replace both tools",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := common.GetConfig()
		if path := viper.GetString("bc.layout"); path != "" {
			if viper.GetBool("bc.edit") {

				editor := gofn.FirstNonEmpty(cfg.OverrideEditor, cfg.Editor, "vi")
				cmd := exec.Command(editor, path) // #nosec G204
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					log.Fatal("Encountered an error while opening $EDITOR", err)
					utils.ExitError(utils.ExitOpenEditorError)
				}
				utils.ExitError(utils.ExitOk)
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
				layouts := utils.GetFilesRecursively(loc, utils.IsYaml)
				for _, l := range layouts {
					// TODO: log
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

	err1 := viper.BindPFlag("bc.layout", bcCmd.Flags().Lookup("layout"))
	err2 := viper.BindPFlag("bc.here", bcCmd.Flags().Lookup("here"))
	err3 := viper.BindPFlag("bc.edit", bcCmd.Flags().Lookup("edit"))
	err4 := viper.BindPFlag("bc.show", bcCmd.Flags().Lookup("show"))
	err5 := viper.BindPFlag("bc.list", bcCmd.Flags().Lookup("list"))

	if err := errors.Join(err1, err2, err3, err4, err5); err != nil {
		log.Fatal("Failed to bind config", "err", err)
		utils.ExitError(utils.ExitConfigBindError)
	}
}

func init() {
	rootCmd.AddCommand(bcCmd)
}
