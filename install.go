package vin

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
	"github.com/mholt/archiver/v3"
)

func (v *Vin) download(url string) (string, error) {
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

	// display progressbar
	tmpl := fmt.Sprintf("{{ \"%s\" | green }} ", archiveName) +
		`{{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}} {{ speed . | magenta }} {{ percent . | cyan }}`
	bar := pb.ProgressBarTemplate(tmpl).Start64(resp.ContentLength)
	barReader := bar.NewProxyReader(resp.Body)
	if _, err := io.Copy(out, barReader); err != nil {
		return "", err
	}
	bar.Finish()
	return archivePath, nil
}

func isExecutable(info os.FileInfo) bool {
	return info.Mode()&0111 != 0
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
			return os.Rename(path, filepath.Join(v.binDir(), app.Name))
		}
		return nil
	})
}

func (v *Vin) install(app App, url string) error {
	archivePath, err := v.download(url)
	if err != nil {
		return err
	}

	tmpDir := filepath.Dir(archivePath)
	defer os.RemoveAll(tmpDir)

	if err := archiver.Unarchive(archivePath, tmpDir); err != nil {
		return err
	}

	if err := v.pickExecutable(app, tmpDir); err != nil {
		return err
	}
	return nil
}
