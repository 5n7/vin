package cmd

import (
	"github.com/skmatz/vin/cli"
	"github.com/spf13/cobra"
)

func runExample(cmd *cobra.Command, args []string) error {
	c := cli.New()
	return c.Example()
}

var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Show a config example",
	Long:  "Show a config example.",
	RunE:  runExample,
}

func init() {
	rootCmd.AddCommand(exampleCmd)
}
