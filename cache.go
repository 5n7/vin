package vin

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/skmatz/vin/cache"
)

func (v *Vin) CacheDir() (string, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cache, "vin"), nil
}

var repoReplacer = strings.NewReplacer("/", "---")

func (v *Vin) AppAlreadyInstalled(app App) (bool, error) {
	cacheDir, err := v.CacheDir()
	if err != nil {
		return false, err
	}

	c := cache.New(cacheDir)
	key := repoReplacer.Replace(app.Repo)
	value := c.GetString(key)
	return value == *app.release.TagName, nil
}

func (v *Vin) SaveCache(app App) error {
	cacheDir, err := v.CacheDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return err
	}

	c := cache.New(cacheDir)
	key := repoReplacer.Replace(app.Repo)
	value := *app.release.TagName
	if err := c.SetString(key, value); err != nil {
		return err
	}
	return nil
}
