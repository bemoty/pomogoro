//go:build linux || darwin

package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync/atomic"

	"fyne.io/systray"
)

var sockFile = runtimeFile("pomogoro.sock")
var ipcListener net.Listener
var latestState atomic.Value // stores state

func formatStatus(s state) string {
	mins := int(s.remaining.Minutes())
	secs := int(s.remaining.Seconds()) % 60

	pomNum := s.completedPomodoros
	if s.phase == work {
		pomNum++
	}

	desk := "sitting"
	if s.standing() {
		desk = "standing"
	}

	line := fmt.Sprintf("%s %02d:%02d %d/%d %s", s.phase.letter(), mins, secs, pomNum, pomodorosPerCycle, desk)
	if s.paused {
		line += " paused"
	}
	return line
}

func listenIPC(cmds chan<- command) {
	_ = os.Remove(sockFile)
	ln, err := net.Listen("unix", sockFile)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "ipc listen:", err)
		return
	}
	ipcListener = ln
	defer func() { _ = os.Remove(sockFile) }()
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		go handleConn(conn, cmds)
	}
}

func handleConn(conn net.Conn, cmds chan<- command) {
	defer func() { _ = conn.Close() }()
	buf := make([]byte, 64)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	switch strings.TrimSpace(string(buf[:n])) {
	case "pause":
		cmds <- cmdTogglePause
		_, _ = fmt.Fprintln(conn, "ok")
	case "skip":
		cmds <- cmdSkip
		_, _ = fmt.Fprintln(conn, "ok")
	case "reset":
		cmds <- cmdReset
		_, _ = fmt.Fprintln(conn, "ok")
	case "stop":
		_, _ = fmt.Fprintln(conn, "ok")
		_ = conn.Close()
		cmds <- cmdQuit
		systray.Quit()
	case "status":
		v := latestState.Load()
		if v == nil {
			_, _ = fmt.Fprintln(conn, "starting")
			return
		}
		_, _ = fmt.Fprintln(conn, formatStatus(v.(state)))
	default:
		_, _ = fmt.Fprintln(conn, "unknown command")
	}
}

func clientCmd(cmd string) {
	conn, err := net.Dial("unix", sockFile)
	if err != nil {
		if cmd == "status" {
			fmt.Println("stopped")
		} else {
			_, _ = fmt.Fprintln(os.Stderr, "pomogoro not running")
		}
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()
	_, _ = fmt.Fprintln(conn, cmd)
	var buf strings.Builder
	_, _ = io.Copy(&buf, conn)
	if out := strings.TrimSpace(buf.String()); out != "ok" {
		fmt.Println(out)
	}
}
