package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/skmatz/vin"
	"github.com/vbauerster/mpb/v5"
	"golang.org/x/sync/errgroup"
)

// CLI represents a CLI for Vin.
type CLI struct{}

func New() *CLI {
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

func (c *CLI) selectApps(v vin.Vin) ([]vin.App, error) {
	// allApps is a map for referencing applications by repository name
	allApps := make(map[string]vin.App)
	for _, app := range v.Apps {
		allApps[app.Repo] = app
	}

	repos := make([]string, 0)
	prompt := &survey.MultiSelect{
		Message: "select applications to install",
		Options: v.Repos(),
	}
	if err := survey.AskOne(prompt, &repos); err != nil {
		return nil, err
	}

	apps := make([]vin.App, 0)
	for _, repo := range repos {
		apps = append(apps, allApps[repo])
	}
	return apps, nil
}

// Options represents options for the CIL.
type Options struct {
	SelectApps bool
}

// Run runs the CLI.
func (c *CLI) Run(opt Options) error {
	configPath, err := c.defaultConfigPath()
	if err != nil {
		return err
	}

	tokenPath, err := c.defaultTokenPath()
	if err != nil {
		return err
	}

	v, err := vin.New(configPath, tokenPath)
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
			urls := app.SuitableAssetURLs()
			if len(urls) == 0 {
				return fmt.Errorf("no suitable assets are found: %s", app.Repo)
			}

			if len(urls) > 1 {
				return fmt.Errorf("cannot identify one asset; %d assets are found: %s", len(urls), app.Repo)
			}

			if err := v.Install(app, urls[0], p); err != nil {
				return err
			}

			if err := app.RunCommand(); err != nil {
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
