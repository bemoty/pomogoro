//go:build !linux && !darwin && !windows

package main

import (
	"os"
	"syscall"
)

var shutdownSignals = []os.Signal{syscall.SIGTERM, syscall.SIGINT}
