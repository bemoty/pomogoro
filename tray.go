package main

import (
	"fmt"

	"fyne.io/systray"
)

type uiUpdate struct {
	title         string
	pauseLabel    string
	skipLabel     string
	pomodoroCount string
	isWork        bool
	deskState     string
	progress      float64
}

func (s *state) toUpdate(pl string) uiUpdate {
	isWork := s.phase == work
	ds := deskLabel(s.deskUp)
	if !isWork {
		ds = deskLabel(!s.deskUp)
	}

	skipLabel := "Skip to break"
	if !isWork {
		skipLabel = "Skip to work"
	}

	var pomCount string
	switch s.phase {
	case work:
		pomCount = fmt.Sprintf("Pomodoro %d / %d", s.completedPomodoros+1, pomodorosPerCycle)
	case shortBreak:
		pomCount = fmt.Sprintf("Break (%d / %d done)", s.completedPomodoros, pomodorosPerCycle)
	case longBreak:
		pomCount = fmt.Sprintf("Long break (%d / %d done)", s.completedPomodoros, pomodorosPerCycle)
	}

	return uiUpdate{
		title:         s.trayTitle(),
		pauseLabel:    pl,
		skipLabel:     skipLabel,
		pomodoroCount: pomCount,
		isWork:        isWork,
		deskState:     ds,
		progress:      s.progress(),
	}
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

	cmds := make(chan command, 4)

	go runMenuLoop(pauseItem, skipItem, resetItem, quitItem, cmds)

	go runTimer(cmds, func(u uiUpdate) {
		systray.SetTooltip("pomogoro: " + u.title)
		systray.SetTemplateIcon(renderTemplateIcon(u.progress), renderIcon(u.progress, u.isWork))
		statusItem.SetTitle(u.title)
		pomodoroItem.SetTitle(u.pomodoroCount)
		pauseItem.SetTitle(u.pauseLabel)
		skipItem.SetTitle(u.skipLabel)
		deskItem.SetTitle(u.deskState)
	})
}

func runMenuLoop(pauseItem, skipItem, resetItem, quitItem *systray.MenuItem, cmds chan<- command) {
	for {
		select {
		case <-pauseItem.ClickedCh:
			cmds <- cmdTogglePause
		case <-skipItem.ClickedCh:
			cmds <- cmdSkip
		case <-resetItem.ClickedCh:
			cmds <- cmdReset
		case <-quitItem.ClickedCh:
			cmds <- cmdQuit
			systray.Quit()
		}
	}
}

func onExit() {}
