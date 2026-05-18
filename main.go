package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/systray"
)

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
	defer releasePID()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		<-sig
		signal.Reset()
		systray.Quit()
	}()

	systray.Run(onReady, onExit)
}
