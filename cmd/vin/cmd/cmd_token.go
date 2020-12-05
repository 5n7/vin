package cmd

import (
	"github.com/skmatz/vin"
	"github.com/spf13/cobra"
)

func runToken(cmd *cobra.Command, args []string) error {
	cli := vin.NewCLI()
	return cli.AskGitHubAccessToken()
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Prompt for the GitHub access token",
	Long:  "Prompt for the GitHub access token.",
	RunE:  runToken,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(tokenCmd)
}
