package vin

import (
	"fmt"
	"os/exec"
	"strings"
)

func (a *App) RunCommand() error {
	cmds := strings.Split(a.Command, "\n")
	for _, cmd := range cmds {
		if cmd == "" {
			continue
		}

		out, err := exec.Command("/bin/sh", "-c", cmd).Output() //nolint:gosec
		if err != nil {
			return err
		}

		if o := string(out); o != "" {
			fmt.Println(o)
		}
	}
	return nil
}
