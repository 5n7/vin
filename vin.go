package vin

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Vin represents a Vin client.
type Vin struct {
	Apps []App `toml:"app"`

	vinDir string
}

// newGitHubClient returns a GitHub client.
func newGitHubClient(token string) *github.Client {
	if token == "" {
		// w/o authentication
		return github.NewClient(nil)
	}
	// w/ authentication
	ctx := context.Background()
	sts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	return github.NewClient(oauth2.NewClient(ctx, sts))
}

func vinDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".vin"), nil
}

func (v *Vin) binDir() string {
	return filepath.Join(v.vinDir, "bin")
}

func (v *Vin) tmpDir() string {
	return filepath.Join(v.vinDir, "tmp")
}

// New returns a Vin client.
func New(path string) (*Vin, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var v Vin
	if _, err := toml.Decode(string(b), &v); err != nil {
		return nil, err
	}

	dir, err := vinDir()
	if err != nil {
		return nil, err
	}
	v.vinDir = dir

	if err := os.MkdirAll(v.binDir(), os.ModePerm); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(v.tmpDir(), os.ModePerm); err != nil {
		return nil, err
	}

	gh := newGitHubClient("")
	for i := range v.Apps {
		if err := v.Apps[i].init(gh); err != nil {
			return nil, err
		}
	}
	return &v, nil
}
