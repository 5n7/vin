package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/skmatz/vin"

	"github.com/AlecAivazis/survey/v2"
)

const tokenGenerateURL = "https://github.com/settings/tokens/new?description=Vin" //nolint:gosec

// AskGitHubAccessToken prompts for the GitHub access token.
func (c *CLI) AskGitHubAccessToken() (string, error) {
	tokenPath, err := c.defaultTokenPath()
	if err != nil {
		return "", err
	}

	var token string
	if _, err := os.Stat(tokenPath); !os.IsNotExist(err) {
		b, err := ioutil.ReadFile(tokenPath)
		if err != nil {
			return "", err
		}

		var t vin.Token
		if err := json.Unmarshal(b, &t); err != nil {
			return "", err
		}
		token = t.Token
	}

	fmt.Println(tokenGenerateURL)
	prompt := &survey.Input{
		Message: "Input your token:",
		Default: token,
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

	t := vin.Token{Token: token}
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
