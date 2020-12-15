package cli

import "github.com/fatih/color"

// CLI represents a CLI for Vin.
type CLI struct{}

func New() *CLI {
	return &CLI{}
}

var (
	cyan    = color.New(color.FgCyan)
	green   = color.New(color.FgGreen)
	magenta = color.New(color.FgMagenta)
)
