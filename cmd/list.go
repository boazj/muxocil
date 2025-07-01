package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		defer w.Flush()
		fmt.Fprintln(w, "Path\tName\t")
		fmt.Fprintln(w, "----------------------------\t-------------\t")

		for _, loc := range cfg.LayoutSearchLocations {
			eloc := os.ExpandEnv(loc)
			dir, err := utils.IsDirectory(eloc)
			if err != nil || !dir {
				if os.IsNotExist(err) || !dir {
					fmt.Printf("Configuration location %s does not exist or is not a directory\n", loc)
					continue
				} else {
					panic(err)
				}
			}
			layouts := utils.GetFilesRecursively(eloc, func(path string) bool {
				return strings.HasSuffix(path, ".yaml")
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
