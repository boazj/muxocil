package cmd

import (
	"fmt"
	"strings"

	"github.com/boazj/muxocil/provider"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"github.com/tiendc/gofn"
)

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: "List all available multiplexing providers in muxocil configuration",
	Run: func(cmd *cobra.Command, args []string) {
		rows := make([][]string, 0)

		for _, p := range provider.ProviderDefs.GetSupportedProviders() {
			rows = append(rows, []string{
				p.Display,
				p.Kind.String(),
				strings.Join(gofn.ToStringSlice[string](p.SupportedOs), ", "),
				string(p.ID),
			})
		}

		headerStyle := lipgloss.NewStyle().Bold(true).Align(lipgloss.Center)
		rowStyle := lipgloss.NewStyle().Padding(0, 1)

		t := table.New().
			Border(lipgloss.DoubleBorder()).
			BorderStyle(lipgloss.NewStyle()).
			StyleFunc(func(row int, col int) lipgloss.Style {
				if row == table.HeaderRow {
					return headerStyle
				}
				return rowStyle
			}).
			Headers("Name", "Kind", "Supported OS", "Config String").
			Rows(rows...)

		fmt.Println(t)
	},
}

func init() {
	rootCmd.AddCommand(providersCmd)
}
