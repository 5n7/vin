package main

import (
	"fmt"
	"os"

	"github.com/skmatz/vin"
)

func run() error {
	cli := vin.NewCLI()
	return cli.Run()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
