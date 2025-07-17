package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/boazj/muxocil/common"
	"github.com/boazj/muxocil/utils"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

var layoutsCmd = &cobra.Command{
	Use:   "layouts",
	Short: "List all layouts available",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := common.GetConfig()

		rows := make([][]string, 0)
		for _, loc := range cfg.GetLayoutSearchLocations() {
			layouts := utils.GetFilesRecursively(loc, utils.IsYaml)
			for _, l := range layouts {
				rows = append(rows, []string{l, filepath.Base(l)})
			}
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
			Headers("Path", "Name").
			Rows(rows...)

		fmt.Println(t)
	},
}

func init() {
	rootCmd.AddCommand(layoutsCmd)
}
