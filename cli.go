package vin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
)

// CLI represents a CLI for Vin.
type CLI struct{}

func NewCLI() *CLI {
	return &CLI{}
}

func (c *CLI) defaultConfigPath() (string, error) {
	cfg, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfg, "vin", "vin.toml"), nil
}

func (c *CLI) defaultTokenPath() (string, error) {
	cfg, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfg, "vin", "token.json"), nil
}

// Run runs the CLI.
func (c *CLI) Run() error {
	configPath, err := c.defaultConfigPath()
	if err != nil {
		return err
	}

	tokenPath, err := c.defaultTokenPath()
	if err != nil {
		return err
	}

	v, err := New(configPath, tokenPath)
	if err != nil {
		return err
	}

	for _, app := range v.Apps {
		urls := app.suitableAssetURLs()
		if len(urls) == 0 {
			fmt.Fprintf(os.Stderr, "no suitable assets are found: %s\n", app.Repo)
			continue
		}

		for _, url := range urls {
			if err := v.install(app, url); err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		}

		if err := app.runCommand(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
	}
	return nil
}

const tokenGenerateURL = "https://github.com/settings/tokens/new?description=Vin" //nolint:gosec

// AskGitHubAccessToken prompts for the GitHub access token.
func (c *CLI) AskGitHubAccessToken() error {
	tokenPath, err := c.defaultTokenPath()
	if err != nil {
		return err
	}

	fmt.Println(tokenGenerateURL)
	var token string
	prompt := &survey.Input{
		Message: "input your token:",
	}
	if err := survey.AskOne(prompt, &token); err != nil {
		return err
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
