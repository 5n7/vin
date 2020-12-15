package cli

import (
	"bytes"
	"fmt"

	"github.com/naoina/toml"
	"github.com/skmatz/vin"
)

// Example shows a config example.
func (c *CLI) Example() error {
	v := vin.Vin{
		Apps: []vin.App{{
			Repo: "cli/cli",
		}},
	}
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(v); err != nil {
		return err
	}
	fmt.Print(buf.String())
	return nil
}
