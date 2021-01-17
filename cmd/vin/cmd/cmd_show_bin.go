package cmd

import (
	"github.com/skmatz/vin/cli"
	"github.com/spf13/cobra"
)

func runShowBin(cmd *cobra.Command, args []string) error {
	c := cli.New()
	return c.ShowBinDir()
}

var showBinCmd = &cobra.Command{
	Use:   "bin",
	Short: "Show the path to the bin directory",
	Long:  "Show the path to the bin directory.",
	RunE:  runShowBin,
}

func init() {
	showCmd.AddCommand(showBinCmd)
}
