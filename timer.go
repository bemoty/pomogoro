package main

import (
	"fmt"
	"time"
)

type Phase int

const (
	Work Phase = iota
	ShortBreak
	LongBreak
)

type Command int

const (
	CmdTogglePause Command = iota
	CmdSkip
	CmdReset
	CmdQuit
)

const (
	workDuration       = 25 * time.Minute
	shortBreakDuration = 5 * time.Minute
	longBreakDuration  = 15 * time.Minute
	pomodorosPerCycle  = 4
)

type state struct {
	phase              Phase
	remaining          time.Duration
	completedPomodoros int
	deskUp             bool
	paused             bool
}

type UIUpdate struct {
	title      string
	pauseLabel string
	isWork     bool
	deskState  string
}

func deskLabel(standing bool) string {
	if standing {
		return "Desk state: Standing"
	}
	return "Desk state: Sitting"
}

func (s *state) enterWork() {
	s.deskUp = !s.deskUp
	s.phase = Work
	s.remaining = workDuration
}

func (s *state) trayTitle() string {
	total := int(s.remaining.Seconds())
	mins := total / 60
	secs := total % 60

	var prefix string
	switch s.phase {
	case Work:
		prefix = "W"
	case ShortBreak:
		prefix = "B"
	case LongBreak:
		prefix = "LB"
	}

	title := fmt.Sprintf("%s %02d:%02d", prefix, mins, secs)
	if s.paused {
		title += " ❄"
	}
	return title
}

func runTimer(cmds <-chan Command, update func(UIUpdate)) {
	s := state{}
	s.enterWork()
	notify("Work", deskLabel(s.deskUp))
	update(UIUpdate{title: s.trayTitle(), pauseLabel: "Pause", isWork: true, deskState: deskLabel(s.deskUp)})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if s.paused {
				continue
			}
			s.remaining -= time.Second
			if s.remaining > 0 {
				isWork := s.phase == Work
				ds := deskLabel(s.deskUp)
				if !isWork {
					ds = deskLabel(!s.deskUp)
				}
				update(UIUpdate{title: s.trayTitle(), pauseLabel: "Pause", isWork: isWork, deskState: ds})
				continue
			}

			switch s.phase {
			case Work:
				s.completedPomodoros++
				if s.completedPomodoros >= pomodorosPerCycle {
					s.phase = LongBreak
					s.remaining = longBreakDuration
					s.completedPomodoros = 0
					notify("Long break", "15 minutes. Well done. "+deskLabel(!s.deskUp))
				} else {
					s.phase = ShortBreak
					s.remaining = shortBreakDuration
					notify("Short break", "5 minutes. "+deskLabel(!s.deskUp))
				}
			case ShortBreak, LongBreak:
				s.enterWork()
				notify("Work", deskLabel(s.deskUp))
			}
			isWork := s.phase == Work
			ds := deskLabel(s.deskUp)
			if !isWork {
				ds = deskLabel(!s.deskUp)
			}
			update(UIUpdate{title: s.trayTitle(), pauseLabel: "Pause", isWork: isWork, deskState: ds})

		case cmd := <-cmds:
			switch cmd {
			case CmdTogglePause:
				s.paused = !s.paused
				isWork := s.phase == Work
				ds := deskLabel(s.deskUp)
				if !isWork {
					ds = deskLabel(!s.deskUp)
				}
				update(UIUpdate{title: s.trayTitle(), pauseLabel: pauseLabel(s.paused), isWork: isWork, deskState: ds})

			case CmdSkip:
				s.paused = false
				switch s.phase {
				case Work:
					s.phase = ShortBreak
					s.remaining = shortBreakDuration
					notify("Short break", "5 minutes. "+deskLabel(!s.deskUp))
				case ShortBreak, LongBreak:
					s.enterWork()
					notify("Work", deskLabel(s.deskUp))
				}
				isWork := s.phase == Work
				ds := deskLabel(s.deskUp)
				if !isWork {
					ds = deskLabel(!s.deskUp)
				}
				update(UIUpdate{title: s.trayTitle(), pauseLabel: "Pause", isWork: isWork, deskState: ds})

			case CmdReset:
				s.completedPomodoros = 0
				s.paused = false
				s.enterWork()
				notify("Work", deskLabel(s.deskUp))
				update(UIUpdate{title: s.trayTitle(), pauseLabel: "Pause", isWork: true, deskState: deskLabel(s.deskUp)})

			case CmdQuit:
				return
			}
		}
	}
}

func pauseLabel(paused bool) string {
	if paused {
		return "Resume"
	}
	return "Pause"
}
