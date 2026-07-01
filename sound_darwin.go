//go:build darwin

package main

import (
	"os/exec"
)

func playDing() {
	path := ensureDing()
	_ = exec.Command("afplay", path).Start()
}
