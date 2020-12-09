package cmd

import (
	"github.com/skmatz/vin/cli"
	"github.com/spf13/cobra"
)

func runToken(cmd *cobra.Command, args []string) error {
	c := cli.New()
	token, err := c.AskGitHubAccessToken()
	if err != nil {
		return err
	}
	return c.StoreAccessToken(token)
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
