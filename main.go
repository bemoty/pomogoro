package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"fyne.io/systray"
)

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if len(arg) > 0 && arg[0] != '-' {
			switch arg {
			case "pause", "skip", "reset", "stop", "status":
				clientCmd(arg)
				return
			default:
				fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", arg)
				os.Exit(1)
			}
		}
	}

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
	signal.Notify(sig, shutdownSignals...)
	go func() {
		<-sig
		signal.Reset()
		systray.Quit()
	}()

	systray.Run(onReady, onExit)
}
