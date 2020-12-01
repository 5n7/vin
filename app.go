package vin

import (
	"context"
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"github.com/google/go-github/github"
)

// App represents an application.
type App struct {
	// Repo is the GitHub repository name in "owner/repo" format.
	Repo string `toml:"repo"`

	// Tag is the tag on GitHub.
	Tag string `toml:"tag"`

	// Keywords is a list of keywords for selecting suitable assets from multiple assets.
	Keywords []string `toml:"keywords"`

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

// contains returns whether all substrs are within s.
func contains(s string, substrs []string) bool {
	for _, substr := range substrs {
		if !strings.Contains(s, substr) {
			return false
		}
	}
	return true
}

var archiveRegexp = regexp.MustCompile(`\.(tar\.gz|tgz|zip)$`)

func (a *App) suitableURLs(urls []string) []string {
	r := make([]string, 0)
	for _, url := range urls {
		if contains(url, a.Keywords) && archiveRegexp.MatchString(url) {
			r = append(r, url)
		}
	}
	return r
}

func (a *App) suitableAssetURLs() []string {
	urls := make([]string, 0)
	for _, asset := range a.release.Assets {
		urls = append(urls, *asset.BrowserDownloadURL)
	}
	return a.suitableURLs(urls)
}
