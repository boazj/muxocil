package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/boazj/muxocil/common"
	"github.com/boazj/muxocil/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all layouts available",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := common.GetConfig()
		// TODO: tea
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		defer w.Flush()
		fmt.Fprintln(w, "Path\tName\t")
		fmt.Fprintln(w, "----------------------------\t-------------\t")

		for _, loc := range cfg.GetLayoutSearchLocations() {
			layouts := utils.GetFilesRecursively(loc, func(path string) bool {
				return utils.HasSuffix(path, ".yaml", ".yml")
			})
			for _, l := range layouts {
				fmt.Fprintf(
					w,
					"%s\t%s\t\n",
					l,
					filepath.Base(l),
				)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
