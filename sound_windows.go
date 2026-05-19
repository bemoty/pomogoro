//go:build windows

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func playDing() {
	path := ensureDing()
	uri := strings.ReplaceAll(path, `\`, `/`)
	script := fmt.Sprintf(
		`Add-Type -AssemblyName PresentationCore; `+
			`$mp = [System.Windows.Media.MediaPlayer]::new(); `+
			`$mp.Open([uri]'file:///%s'); `+
			`$mp.Play(); `+
			`[System.Threading.Thread]::Sleep(4000)`,
		uri,
	)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", script).Start()
}
