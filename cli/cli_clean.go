package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/skmatz/vin"
)

var cyan = color.New(color.FgCyan)

// Clean cleans cache files.
func (c *CLI) Clean() error {
	var v vin.Vin

	cacheDir, err := v.CacheDir()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(cacheDir); err != nil {
		return err
	}
	fmt.Println(cyan.Sprintf("cache directory removed: %s", cacheDir))
	return nil
}
