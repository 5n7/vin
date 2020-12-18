package vin

import (
	"os/exec"
)

func (a *App) RunCommand() error {
	return exec.Command("/bin/sh", "-c", a.Command).Run() //nolint:gosec
}
