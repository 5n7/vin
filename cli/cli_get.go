package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/skmatz/vin"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
	"golang.org/x/sync/errgroup"
)

// GetOptions represents options for the CIL.
type GetOptions struct {
	ConfigPath   string
	IgnoreCache  bool
	IgnoreFilter bool
	Priority     int
	SelectApps   bool
}

func (c *CLI) defaultConfigPath() (string, error) {
	cfg, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfg, "vin", "vin.toml"), nil
}

func (c *CLI) configPath(opt GetOptions) (string, error) {
	if opt.ConfigPath == "" {
		path, err := c.defaultConfigPath()
		if err != nil {
			return "", err
		}
		return path, nil
	}
	return opt.ConfigPath, nil
}

func (c *CLI) defaultTokenPath() (string, error) {
	cfg, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfg, "vin", "token.json"), nil
}

// sanityCheck checks the app is ready to install.
func (c *CLI) sanityCheck(app vin.App) error {
	urls := app.SuitableAssetURLs()
	if len(urls) == 0 {
		return fmt.Errorf("no suitable assets are found: %s", app.Repo)
	}

	if len(urls) > 1 {
		return fmt.Errorf("cannot identify one asset; %d assets are found: %s", len(urls), app.Repo)
	}
	return nil
}

func (c *CLI) applyFilters(v *vin.Vin, opt GetOptions) (*vin.Vin, error) {
	host, err := os.Hostname()
	if err != nil {
		return v, err
	}
	v = v.FilterByHost(host)

	v = v.FilterByPriority(opt.Priority)
	return v, nil
}

func (c *CLI) selectApps(v *vin.Vin) (*vin.Vin, error) {
	repos := make([]string, 0)
	prompt := &survey.MultiSelect{
		Message: "Select applications to install",
		Options: v.Repos(),
	}
	if err := survey.AskOne(prompt, &repos); err != nil {
		return nil, err
	}
	return v.FilterByRepo(repos), nil
}

type defaultReadCloserWrapper struct {
	p *mpb.Progress
}

func (w defaultReadCloserWrapper) Wrap(a vin.App, r io.ReadCloser, l int64) io.ReadCloser {
	tag, err := a.TagName()
	if err != nil {
		return r
	}

	bar := w.p.AddBar(
		l,
		mpb.BarStyle("[=>-]"),
		mpb.PrependDecorators(
			decor.Name(green.Sprintf("%s@%s", a.Repo, tag)),
		),
		mpb.AppendDecorators(
			decor.EwmaSpeed(decor.UnitKiB, magenta.Sprint("% .2f"), 60),
			decor.Name(" "),
			decor.NewPercentage(cyan.Sprint("% d"), decor.WCSyncSpace),
		),
	)
	return bar.ProxyReader(r)
}

// Get gets applications.
func (c *CLI) Get(opt GetOptions) error { //nolint:funlen,gocognit
	vinDir, err := vin.DefaultVinDir()
	if err != nil {
		return err
	}

	v, err := vin.New(vinDir)
	if err != nil {
		return err
	}

	tokenPath, err := c.defaultTokenPath()
	if err != nil {
		return err
	}

	var token string
	if t := os.Getenv("GITHUB_TOKEN"); t != "" {
		// read token from environment variable
		token = t
	} else if _, err := os.Stat(tokenPath); !os.IsNotExist(err) {
		// read token from JSON file
		t, err := vin.TokenFromJSON(tokenPath)
		if err != nil {
			return err
		}
		token = t
	}

	configPath, err := c.configPath(opt)
	if err != nil {
		return err
	}

	if err := v.ReadTOML(configPath); err != nil {
		return err
	}

	if err := v.FetchApps(token); err != nil {
		return err
	}

	repos := make([]string, 0)
	var result *multierror.Error
	for _, app := range v.Apps {
		if err := c.sanityCheck(app); err != nil {
			result = multierror.Append(result, err)
		} else {
			repos = append(repos, app.Repo)
		}
	}

	if err := result.ErrorOrNil(); err != nil {
		result.ErrorFormat = func(errs []error) string {
			points := make([]string, len(errs))
			for i, err := range errs {
				points[i] = fmt.Sprintf("\t- %s", err)
			}
			return fmt.Sprintf(
				"%d app(s) skipped installation:\n%s\n",
				len(errs),
				strings.Join(points, "\n"),
			)
		}
		fmt.Fprintln(os.Stderr, err)
		v = v.FilterByRepo(repos)
	}

	if !opt.IgnoreFilter {
		vin, err := c.applyFilters(v, opt)
		if err != nil {
			return err
		}
		v = vin
	}

	if opt.SelectApps {
		vin, err := c.selectApps(v)
		if err != nil {
			return err
		}
		v = vin
	}

	p := mpb.New(
		mpb.WithRefreshRate(180 * time.Millisecond), //nolint:gomnd
	)
	wrapper := defaultReadCloserWrapper{p: p}

	var eg errgroup.Group
	for _, app := range v.Apps {
		app := app
		eg.Go(func() error {
			if !opt.IgnoreCache {
				exists, err := v.AppAlreadyInstalled(app)
				if err != nil || exists {
					return err
				}
			}

			if err := v.Install(app, app.SuitableAssetURLs()[0], wrapper); err != nil {
				return err
			}

			if err := app.RunCommand(); err != nil {
				return err
			}

			if err := v.SaveCache(app); err != nil {
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
