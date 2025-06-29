package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: "List all available multiplexing providers in muxocil configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("providers called")
	},
}

func init() {
	rootCmd.AddCommand(providersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// providersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// providersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
