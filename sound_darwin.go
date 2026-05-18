//go:build darwin

package main

import (
	"os/exec"
)

func playDing() {
	path := ensureDing()
	exec.Command("afplay", path).Start()
}
