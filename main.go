package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"fyne.io/systray"
)

const pidFile = "/tmp/pomogoro.pid"

func main() {
	daemon := flag.Bool("d", false, "run in background")
	flag.Parse()

	if *daemon {
		daemonize()
	}

	if err := checkSingleInstance(); err != nil {
		notifyText("pomogoro", "already running")
		os.Exit(1)
	}
	writePID()
	defer os.Remove(pidFile)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		<-sig
		signal.Reset()
		systray.Quit()
	}()

	systray.Run(onReady, onExit)
}

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

func onReady() {
	systray.SetTemplateIcon(renderTemplateIcon(0), renderIcon(0, true))
	systray.SetTooltip("pomogoro")

	statusItem := systray.AddMenuItem("W 25:00", "")
	statusItem.Disable()
	pomodoroItem := systray.AddMenuItem("Pomodoro 1 / 4", "")
	pomodoroItem.Disable()
	deskItem := systray.AddMenuItem("Desk state: Standing", "")
	deskItem.Disable()
	systray.AddSeparator()
	pauseItem := systray.AddMenuItem("Pause", "")
	skipItem := systray.AddMenuItem("Skip to break", "")
	resetItem := systray.AddMenuItem("Reset", "")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "")

	cmds := make(chan Command, 4)

	go func() {
		for {
			select {
			case <-pauseItem.ClickedCh:
				cmds <- CmdTogglePause
			case <-skipItem.ClickedCh:
				cmds <- CmdSkip
			case <-resetItem.ClickedCh:
				cmds <- CmdReset
			case <-quitItem.ClickedCh:
				cmds <- CmdQuit
				systray.Quit()
			}
		}
	}()

	go runTimer(cmds, func(u UIUpdate) {
		systray.SetTooltip("pomogoro: " + u.title)
		systray.SetTemplateIcon(renderTemplateIcon(u.progress), renderIcon(u.progress, u.isWork))
		statusItem.SetTitle(u.title)
		pomodoroItem.SetTitle(u.pomodoroCount)
		pauseItem.SetTitle(u.pauseLabel)
		skipItem.SetTitle(u.skipLabel)
		deskItem.SetTitle(u.deskState)
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
	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "writePID:", err)
	}
}
