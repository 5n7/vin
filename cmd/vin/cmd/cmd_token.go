package cmd

import (
	"github.com/skmatz/vin/cli"
	"github.com/spf13/cobra"
)

func runToken(cmd *cobra.Command, args []string) error {
	c := cli.New()
	return c.AskGitHubAccessToken()
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Prompt for the GitHub access token",
	Long:  "Prompt for the GitHub access token.",
	RunE:  runToken,
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
