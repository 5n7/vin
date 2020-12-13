package cmd

import (
	"github.com/skmatz/vin/cli"
	"github.com/spf13/cobra"
)

func runClean(cmd *cobra.Command, args []string) error {
	c := cli.New()
	return c.Clean()
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean cache files",
	Long:  "Clean cache files.",
	RunE:  runClean,
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
