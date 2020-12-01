package vin

import (
	"fmt"
	"os"
	"path/filepath"
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

// Run runs the CLI.
func (c *CLI) Run() error {
	configPath, err := c.defaultConfigPath()
	if err != nil {
		return err
	}

	v, err := New(configPath)
	if err != nil {
		return err
	}

	for _, app := range v.Apps {
		urls := app.suitableAssetURLs()
		for _, url := range urls {
			if err := v.install(url); err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		}
	}
	return nil
}
