package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Use a specific layout file, instead of `~/.teamocil/<layout>.yml`",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("use called")
	},
}

func init() {
	rootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// layoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// layoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
