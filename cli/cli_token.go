package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
)

const tokenGenerateURL = "https://github.com/settings/tokens/new?description=Vin" //nolint:gosec

// AskGitHubAccessToken prompts for the GitHub access token.
func (c *CLI) AskGitHubAccessToken() (string, error) {
	fmt.Println(tokenGenerateURL)
	var token string
	prompt := &survey.Input{
		Message: "Input your token:",
	}
	if err := survey.AskOne(prompt, &token); err != nil {
		return "", err
	}
	return token, nil
}

// StoreAccessToken stores the GitHub access token.
func (c *CLI) StoreAccessToken(token string) error {
	tokenPath, err := c.defaultTokenPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(tokenPath); !os.IsNotExist(err) {
		var overwrite bool
		prompt := &survey.Confirm{
			Message: "Token file already exists. Overwrite?",
		}
		if err := survey.AskOne(prompt, &overwrite); err != nil {
			return err
		}

		if !overwrite {
			return nil
		}
	}

	var t = struct {
		Token string `json:"token"`
	}{Token: token}
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(tokenPath), os.ModePerm); err != nil {
		return err
	}

	if err := ioutil.WriteFile(tokenPath, b, os.ModePerm); err != nil {
		return err
	}
	fmt.Printf("your token is stored in %s\n", tokenPath)
	return nil
}
