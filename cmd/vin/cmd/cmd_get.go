package cmd

import (
	"github.com/skmatz/vin"
	"github.com/spf13/cobra"
)

func runGet(cmd *cobra.Command, args []string) error {
	cli := vin.NewCLI()
	return cli.Run()
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get applications",
	Long:  "Get applications.",
	RunE:  runGet,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(getCmd)
}
