package vin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/vbauerster/mpb/v5"
	"golang.org/x/sync/errgroup"
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

func (c *CLI) selectApps(v Vin) ([]App, error) {
	// allApps is a map for referencing applications by repository name
	allApps := make(map[string]App)
	for _, app := range v.Apps {
		allApps[app.Repo] = app
	}

	repos := make([]string, 0)
	prompt := &survey.MultiSelect{
		Message: "select applications to install",
		Options: v.repos(),
	}
	if err := survey.AskOne(prompt, &repos); err != nil {
		return nil, err
	}

	apps := make([]App, 0)
	for _, repo := range repos {
		apps = append(apps, allApps[repo])
	}
	return apps, nil
}

// CLIOptions represents options for the CIL.
type CLIOptions struct {
	SelectApps bool
}

// Run runs the CLI.
func (c *CLI) Run(opt CLIOptions) error {
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

	if opt.SelectApps {
		apps, err := c.selectApps(*v)
		if err != nil {
			return err
		}
		v.Apps = apps
	}

	p := mpb.New(
		mpb.WithRefreshRate(180 * time.Millisecond), //nolint:gomnd
	)

	var eg errgroup.Group
	for _, app := range v.Apps {
		app := app
		eg.Go(func() error {
			urls := app.suitableAssetURLs()
			if len(urls) == 0 {
				return fmt.Errorf("no suitable assets are found: %s", app.Repo)
			}

			if len(urls) > 1 {
				return fmt.Errorf("cannot identify one asset; %d assets are found: %s", len(urls), app.Repo)
			}

			if err := v.install(app, urls[0], p); err != nil {
				return err
			}

			if err := app.runCommand(); err != nil {
				return err
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
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
