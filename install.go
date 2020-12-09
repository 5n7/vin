package vin

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/mholt/archiver/v3"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

var (
	cyan    = color.New(color.FgCyan)
	green   = color.New(color.FgGreen)
	magenta = color.New(color.FgMagenta)
)

func (v *Vin) download(app App, url string, p *mpb.Progress) (string, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http response not OK: %d", resp.StatusCode)
	}

	tmpDir, err := ioutil.TempDir(v.tmpDir(), "")
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(tmpDir)
		}
	}()

	archiveName := filepath.Base(url)
	archivePath := filepath.Join(tmpDir, archiveName)

	out, err := os.Create(archivePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if p == nil {
		// w/o progressbar
		if _, err := io.Copy(out, resp.Body); err != nil {
			return "", err
		}
		return archivePath, nil
	}

	// display progressbar
	bar := p.AddBar(
		resp.ContentLength,
		mpb.BarStyle("[=>-]"),
		mpb.PrependDecorators(
			decor.Name(green.Sprintf("%s@%s", app.Repo, *app.release.TagName)),
		),
		mpb.AppendDecorators(
			decor.EwmaSpeed(decor.UnitKiB, magenta.Sprint("% .2f"), 60),
			decor.Name(" "),
			decor.NewPercentage(cyan.Sprint("% d"), decor.WCSyncSpace),
		),
	)

	proxyReader := bar.ProxyReader(resp.Body)
	defer proxyReader.Close()

	if _, err := io.Copy(out, proxyReader); err != nil {
		return "", err
	}
	return archivePath, nil
}

func isExecutable(info os.FileInfo) bool {
	return info.Mode()&0111 != 0
}

func (v *Vin) place(app App, src string) error {
	if app.Name == "" {
		app.Name = filepath.Base(src)
	}
	dst := filepath.Join(v.binDir(), app.Name)
	if err := os.Rename(src, dst); err != nil {
		return err
	}
	return os.Chmod(dst, 0755)
}

// pickExecutable picks all executable files and moves them to the bin directory.
func (v *Vin) pickExecutable(app App, rootDir string) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if isExecutable(info) {
			if app.Name == "" {
				app.Name = info.Name()
			}
			return v.place(app, path)
		}
		return nil
	})
}

func (v *Vin) Install(app App, url string, p *mpb.Progress) error {
	archivePath, err := v.download(app, url, p)
	if err != nil {
		return err
	}

	tmpDir := filepath.Dir(archivePath)
	defer os.RemoveAll(tmpDir)

	if name := filepath.Base(archivePath); !anyExtRegexp.MatchString(name) {
		// the asset is a binary file
		if app.Name == "" {
			app.Name = name
		}
		return v.place(app, archivePath)
	}

	if err := archiver.Unarchive(archivePath, tmpDir); err != nil {
		return err
	}

	if err := v.pickExecutable(app, tmpDir); err != nil {
		return err
	}
	return nil
}
