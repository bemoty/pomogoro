//go:build windows

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func notifyText(title, body string) {
	t := strings.ReplaceAll(title, "'", "''")
	b := strings.ReplaceAll(body, "'", "''")
	script := fmt.Sprintf(
		`Add-Type -AssemblyName System.Windows.Forms; `+
			`$n = [System.Windows.Forms.NotifyIcon]::new(); `+
			`$n.Icon = [System.Drawing.SystemIcons]::Information; `+
			`$n.Visible = $true; `+
			`$n.ShowBalloonTip(8000, '%s', '%s', [System.Windows.Forms.ToolTipIcon]::Info); `+
			`Start-Sleep -Milliseconds 500; `+
			`$n.Dispose()`,
		t, b,
	)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", script).Run()
}
