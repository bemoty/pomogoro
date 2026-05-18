//go:build linux

package main

import "os/exec"

func notifyText(title, body string) {
	exec.Command("notify-send", "--app-name=pomogoro", "-u", "normal", "-t", "8000", "-i", "dialog-information", title, body).Run()
}
