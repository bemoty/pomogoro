//go:build linux || darwin

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func runtimeFile(name string) string {
	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
		return filepath.Join(dir, name)
	}
	return fmt.Sprintf("/tmp/pomogoro-%d-%s", os.Getuid(), name)
}

type pidManager struct {
	file string
	lock *os.File
}

var pid = &pidManager{file: runtimeFile("pomogoro.pid")}

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

func (p *pidManager) checkSingleInstance() error {
	f, err := os.OpenFile(p.file, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = f.Close()
		return fmt.Errorf("already running")
	}
	p.lock = f
	return nil
}

func (p *pidManager) release() {
	if p.lock == nil {
		return
	}
	_ = p.lock.Close()
	_ = os.Remove(p.file)
}
