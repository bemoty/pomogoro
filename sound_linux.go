//go:build linux

package main

import (
	"os/exec"
)

func playDing() {
	path := ensureDing()
	_ = exec.Command("paplay", path).Start()
}
