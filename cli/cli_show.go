package cli

import (
	"fmt"

	"github.com/skmatz/vin"
)

// ShowBinDir shows the path to the bin directory.
func (c *CLI) ShowBinDir() error {
	vinDir, err := vin.DefaultVinDir()
	if err != nil {
		return err
	}

	v, err := vin.New(vinDir)
	if err != nil {
		return err
	}
	fmt.Println(v.BinDir())
	return nil
}

// ShowTmpDir shows the path to the tmp directory.
func (c *CLI) ShowTmpDir() error {
	vinDir, err := vin.DefaultVinDir()
	if err != nil {
		return err
	}

	v, err := vin.New(vinDir)
	if err != nil {
		return err
	}
	fmt.Println(v.TmpDir())
	return nil
}
