//go:build linux || darwin

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func runtimeFile(name string) string {
	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
		return filepath.Join(dir, name)
	}
	return fmt.Sprintf("/tmp/pomogoro-%d-%s", os.Getuid(), name)
}

var pidFile = runtimeFile("pomogoro.pid")
var pidLock *os.File

func daemonize() {
	self, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "daemonize:", err)
		os.Exit(1)
	}
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		fmt.Fprintln(os.Stderr, "daemonize:", err)
		os.Exit(1)
	}
	cmd := exec.Command(self)
	cmd.Stdin = devNull
	cmd.Stdout = devNull
	cmd.Stderr = devNull
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "daemonize:", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func checkSingleInstance() error {
	f, err := os.OpenFile(pidFile, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return fmt.Errorf("already running")
	}
	pidLock = f
	return nil
}

func writePID() {
	if pidLock == nil {
		return
	}
	pidLock.Truncate(0)
	pidLock.WriteAt([]byte(strconv.Itoa(os.Getpid())), 0)
}

func releasePID() {
	if pidLock == nil {
		return
	}
	pidLock.Close()
	os.Remove(pidFile)
}
