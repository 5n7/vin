package vin

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/google/go-github/github"
)

// App represents an application.
type App struct {
	// Repo is the GitHub repository name in "owner/repo" format.
	Repo string `toml:"repo,omitempty"`

	// Tag is the tag on GitHub.
	Tag string `toml:"tag,omitempty"`

	// Keywords is a list of keywords for selecting suitable assets from multiple assets.
	Keywords []string `toml:"keywords,omitempty"`

	// Name is the name of the executable file.
	Name string `toml:"name,omitempty"`

	// Hosts is a list of host names.
	Hosts []string `toml:"hosts,omitempty"`

	// Priority is the priority of the application.
	Priority int `toml:"priority,omitempty"`

	// Command is the command to run after installation.
	Command string `toml:"command,omitempty"`

	release *github.RepositoryRelease
}

func (a *App) defaultKeywords() []string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	return []string{goos, goarch}
}

func (a *App) parseOwnerAndRepo(s string) (string, string, error) {
	x := strings.Split(s, "/")
	if len(x) != 2 { //nolint:gomnd
		return "", "", fmt.Errorf("invalid format: %s", s)
	}
	return x[0], x[1], nil
}

func (a *App) getRepositoryRelease(gh *github.Client, owner, repo, tag string) (*github.RepositoryRelease, error) {
	ctx := context.Background()
	if tag == "" {
		release, _, err := gh.Repositories.GetLatestRelease(ctx, owner, repo)
		if err != nil {
			return nil, err
		}
		return release, err
	}

	release, _, err := gh.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
	if err != nil {
		return nil, err
	}
	return release, err
}

func (a *App) init(gh *github.Client) error {
	if len(a.Keywords) == 0 {
		a.Keywords = a.defaultKeywords()
	}

	owner, repo, err := a.parseOwnerAndRepo(a.Repo)
	if err != nil {
		return err
	}

	release, err := a.getRepositoryRelease(gh, owner, repo, a.Tag)
	if err != nil {
		return err
	}
	a.release = release
	return nil
}

// contains returns whether all substrs are within s (ignore case).
func contains(s string, substrs []string) bool {
	s = strings.ToLower(s)
	for _, substr := range substrs {
		if !strings.Contains(s, strings.ToLower(substr)) {
			return false
		}
	}
	return true
}

var (
	anyExtRegexp  = regexp.MustCompile(`\.[a-zA-Z0-9]+$`)
	archiveRegexp = regexp.MustCompile(`\.(tar\.gz|tgz|zip)$`)
)

func (a *App) suitableURLs(urls []string) []string {
	r := make([]string, 0)
	for _, url := range urls {
		name := filepath.Base(url)
		if contains(name, a.Keywords) && (!anyExtRegexp.MatchString(name) || archiveRegexp.MatchString(name)) {
			r = append(r, url)
		}
	}
	return r
}

func (a *App) SuitableAssetURLs() []string {
	urls := make([]string, 0)
	for _, asset := range a.release.Assets {
		urls = append(urls, *asset.BrowserDownloadURL)
	}
	return a.suitableURLs(urls)
}

func (a *App) TagName() (string, error) {
	if a.release == nil {
		return "", fmt.Errorf("failed to reference tag name; invalid release")
	}
	return *a.release.TagName, nil
}
