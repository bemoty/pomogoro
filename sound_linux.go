//go:build linux

package main

import (
	"os/exec"
)

func playDing() {
	path := ensureDing()
	exec.Command("paplay", path).Start()
}
