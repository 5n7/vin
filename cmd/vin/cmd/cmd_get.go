package cmd

import (
	"github.com/skmatz/vin"
	"github.com/spf13/cobra"
)

var selectApps bool

func runGet(cmd *cobra.Command, args []string) error {
	cli := vin.NewCLI()
	opt := vin.CLIOptions{
		SelectApps: selectApps,
	}
	return cli.Run(opt)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get applications",
	Long:  "Get applications.",
	RunE:  runGet,
}

func init() {
	getCmd.Flags().BoolVarP(&selectApps, "select", "s", false, "select applications to install")

	rootCmd.AddCommand(getCmd)
}
