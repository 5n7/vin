package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version bool

func runRoot(cmd *cobra.Command, args []string) error {
	if version {
		return runVersion(cmd, args)
	}
	return fmt.Errorf("vin requires a sub-command")
}

var rootCmd = &cobra.Command{
	Use:   "vin",
	Short: "Vin is a next-generation GitHub Releases installer",
	Long:  "Vin is a next-generation GitHub Releases installer.",
	RunE:  runRoot,
}

func init() { //nolint:gochecknoinits
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "show version")
}

func Execute() {
	rootCmd.Execute() //nolint:errcheck
}
