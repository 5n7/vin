package vin

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-github/v33/github"
	"github.com/naoina/toml"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
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

func DefaultVinDir() (string, error) {
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

type Token struct {
	Token string `json:"token"`
}

func TokenFromJSON(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	var t Token
	if err := json.Unmarshal(b, &t); err != nil {
		return "", err
	}
	return t.Token, nil
}

func (v *Vin) ReadTOML(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := toml.NewDecoder(f).Decode(&v); err != nil {
		return err
	}
	return nil
}

func (v *Vin) FetchApps(token string) error {
	// fetch apps from GitHub Releases
	gh := newGitHubClient(token)
	var eg errgroup.Group
	for i := range v.Apps {
		i := i
		eg.Go(func() error {
			if err := v.Apps[i].init(gh); err != nil {
				return err
			}
			return nil
		})

		if err := eg.Wait(); err != nil {
			return err
		}
	}
	return nil
}

// New returns a Vin client.
func New(vinDir string) (*Vin, error) {
	v := Vin{
		vinDir: vinDir,
	}

	if err := os.MkdirAll(v.binDir(), os.ModePerm); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(v.tmpDir(), os.ModePerm); err != nil {
		return nil, err
	}
	return &v, nil
}

func (v *Vin) Repos() []string {
	r := make([]string, 0)
	for _, app := range v.Apps {
		r = append(r, app.Repo)
	}
	return r
}
