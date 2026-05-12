package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/energye/systray"
)

//go:embed icons/work.png
var iconWork []byte

//go:embed icons/break.png
var iconBreak []byte

const pidFile = "/tmp/pomogoro.pid"

func main() {
	if err := checkSingleInstance(); err != nil {
		notify("pomogoro", "already running")
		os.Exit(1)
	}
	writePID()
	defer os.Remove(pidFile)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sig
		systray.Quit()
	}()

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconWork)
	systray.SetTitle("W 25:00")
	systray.SetTooltip("pomogoro")

	statusItem := systray.AddMenuItem("W 25:00", "")
	statusItem.Disable()
	deskItem := systray.AddMenuItem("Desk state: Standing", "")
	deskItem.Disable()
	systray.AddSeparator()
	pauseItem := systray.AddMenuItem("Pause", "")
	skipItem := systray.AddMenuItem("Skip phase", "")
	resetItem := systray.AddMenuItem("Reset", "")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "")

	cmds := make(chan Command, 4)

	pauseItem.Click(func() { cmds <- CmdTogglePause })
	skipItem.Click(func() { cmds <- CmdSkip })
	resetItem.Click(func() { cmds <- CmdReset })
	quitItem.Click(func() {
		cmds <- CmdQuit
		systray.Quit()
	})

	go runTimer(cmds, func(u UIUpdate) {
		systray.SetTitle(u.title)
		statusItem.SetTitle(u.title)
		pauseItem.SetTitle(u.pauseLabel)
		deskItem.SetTitle(u.deskState)
		if u.isWork {
			systray.SetIcon(iconWork)
		} else {
			systray.SetIcon(iconBreak)
		}
	})
}

func onExit() {}

func checkSingleInstance() error {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return nil
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return nil
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil
	}
	if err := proc.Signal(syscall.Signal(0)); err != nil {
		return nil
	}
	return fmt.Errorf("already running (pid %d)", pid)
}

func writePID() {
	os.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644)
}
