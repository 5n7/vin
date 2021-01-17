package cmd

import (
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show any information",
	Long:  "Show any information.",
}

func init() {
	rootCmd.AddCommand(showCmd)
}
