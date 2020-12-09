package cmd

import (
	"github.com/skmatz/vin/cli"
	"github.com/spf13/cobra"
)

var opt cli.Options

func runGet(cmd *cobra.Command, args []string) error {
	c := cli.New()
	return c.Run(opt)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get applications",
	Long:  "Get applications.",
	RunE:  runGet,
}

func init() {
	getCmd.Flags().BoolVarP(&opt.IgnoreFilter, "ignore-filter", "i", false, "ignore all filters")
	getCmd.Flags().IntVarP(&opt.Priority, "priority", "p", 0, "minimum priority for applications to install")
	getCmd.Flags().BoolVarP(&opt.SelectApps, "select", "s", false, "select applications to install")

	rootCmd.AddCommand(getCmd)
}
