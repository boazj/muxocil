package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/boazj/muxocil/provider"
	"github.com/spf13/cobra"
	"github.com/tiendc/gofn"
)

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: "List all available multiplexing providers in muxocil configuration",
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintln(w, "Name\tKind\tSupported OS\tConfig String\t")
		fmt.Fprintln(w, "----\t------------\t------------------\t-------------\t")
		defer w.Flush()

		for _, p := range provider.Providers {
			fmt.Fprintf(
				w,
				"%s\t%s\t%s\t%s\t\n",
				p.Display,
				p.Kind,
				strings.Join(gofn.ToStringSlice[string](p.SupportedOs), ", "),
				p.ID,
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(providersCmd)
}
