package main

import "os/exec"

func notify(title, body string) {
	exec.Command("notify-send", "--app-name=pomogoro", "-u", "normal", "-t", "8000", "-i", "dialog-information", title, body).Run()
	exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/complete.oga").Start()
}
