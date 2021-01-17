package cmd

import (
	"github.com/skmatz/vin/cli"
	"github.com/spf13/cobra"
)

var opt cli.GetOptions

func runGet(cmd *cobra.Command, args []string) error {
	c := cli.New()
	return c.Get(opt)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get applications",
	Long:  "Get applications.",
	RunE:  runGet,
}

func init() {
	getCmd.Flags().StringVarP(&opt.ConfigPath, "config", "c", "", "specify the path to the config file")
	getCmd.Flags().BoolVar(&opt.IgnoreCache, "ignore-cache", false, "ignore cache and install all applications")
	getCmd.Flags().BoolVar(&opt.IgnoreFilter, "ignore-filter", false, "ignore all filters")
	getCmd.Flags().IntVarP(&opt.Priority, "priority", "p", 0, "minimum priority for applications to install")
	getCmd.Flags().BoolVarP(&opt.SelectApps, "select", "s", false, "select applications to install")

	rootCmd.AddCommand(getCmd)
}
