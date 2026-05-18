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

	var phase string
	switch s.phase {
	case work:
		phase = "W"
	case shortBreak:
		phase = "B"
	case longBreak:
		phase = "LB"
	}

	pomNum := s.completedPomodoros
	if s.phase == work {
		pomNum++
	}

	standing := s.deskUp
	if s.phase != work {
		standing = !s.deskUp
	}
	desk := "sitting"
	if standing {
		desk = "standing"
	}

	line := fmt.Sprintf("%s %02d:%02d %d/%d %s", phase, mins, secs, pomNum, pomodorosPerCycle, desk)
	if s.paused {
		line += " paused"
	}
	return line
}

func listenIPC(cmds chan<- command) {
	os.Remove(sockFile)
	ln, err := net.Listen("unix", sockFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ipc listen:", err)
		return
	}
	ipcListener = ln
	defer os.Remove(sockFile)
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		go handleConn(conn, cmds)
	}
}

func handleConn(conn net.Conn, cmds chan<- command) {
	defer conn.Close()
	buf := make([]byte, 64)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	switch strings.TrimSpace(string(buf[:n])) {
	case "pause":
		cmds <- cmdTogglePause
		fmt.Fprintln(conn, "ok")
	case "skip":
		cmds <- cmdSkip
		fmt.Fprintln(conn, "ok")
	case "reset":
		cmds <- cmdReset
		fmt.Fprintln(conn, "ok")
	case "stop":
		fmt.Fprintln(conn, "ok")
		conn.Close()
		cmds <- cmdQuit
		systray.Quit()
	case "status":
		v := latestState.Load()
		if v == nil {
			fmt.Fprintln(conn, "starting")
			return
		}
		fmt.Fprintln(conn, formatStatus(v.(state)))
	default:
		fmt.Fprintln(conn, "unknown command")
	}
}

func clientCmd(cmd string) {
	conn, err := net.Dial("unix", sockFile)
	if err != nil {
		if cmd == "status" {
			fmt.Println("stopped")
		} else {
			fmt.Fprintln(os.Stderr, "pomogoro not running")
		}
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Fprintln(conn, cmd)
	var buf strings.Builder
	io.Copy(&buf, conn)
	if out := strings.TrimSpace(buf.String()); out != "ok" {
		fmt.Println(out)
	}
}
