package main

import "os/exec"

func notify(title, body string) {
	notifySilent(title, body)
	exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/complete.oga").Start()
}

func notifySilent(title, body string) {
	exec.Command("notify-send", "--app-name=pomogoro", "-u", "normal", "-t", "8000", "-i", "dialog-information", title, body).Run()
}
