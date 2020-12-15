package vin

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mholt/archiver/v3"
)

type ReadCloserWrapper interface {
	Wrap(app App, reader io.ReadCloser, contentLength int64) io.ReadCloser
}

func (v *Vin) download(app App, url string, wrapper ReadCloserWrapper) (string, error) {
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

	reader := wrapper.Wrap(app, resp.Body, resp.ContentLength)
	defer reader.Close()

	if _, err := io.Copy(out, reader); err != nil {
		return "", err
	}
	return archivePath, nil
}

func isExecutable(info os.FileInfo) bool {
	if runtime.GOOS == "windows" {
		return filepath.Ext(info.Name()) == ".exe"
	}
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

func (v *Vin) Install(app App, url string, wrapper ReadCloserWrapper) error {
	archivePath, err := v.download(app, url, wrapper)
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
