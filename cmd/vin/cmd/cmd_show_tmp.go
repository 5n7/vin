package cmd

import (
	"github.com/skmatz/vin/cli"
	"github.com/spf13/cobra"
)

func runShowTmp(cmd *cobra.Command, args []string) error {
	c := cli.New()
	return c.ShowTmpDir()
}

var showTmpCmd = &cobra.Command{
	Use:   "tmp",
	Short: "Show the path to the tmp directory",
	Long:  "Show the path to the tmp directory.",
	RunE:  runShowTmp,
}

func init() {
	showCmd.AddCommand(showTmpCmd)
}
