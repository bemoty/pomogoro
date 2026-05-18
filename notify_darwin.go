//go:build darwin

package main

import (
	"fmt"
	"os/exec"
)

func notifyText(title, body string) {
	exec.Command("osascript", "-e",
		fmt.Sprintf(`display notification %q with title %q`, body, title),
	).Run()
}
