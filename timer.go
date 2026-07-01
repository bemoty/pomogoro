package main

import (
	"fmt"
	"time"
)

type phase int

const (
	work phase = iota
	shortBreak
	longBreak
)

type command int

const (
	cmdTogglePause command = iota
	cmdSkip
	cmdReset
	cmdQuit
)

const (
	workDuration       = 25 * time.Minute
	shortBreakDuration = 5 * time.Minute
	longBreakDuration  = 15 * time.Minute
	pomodorosPerCycle  = 4
)

type state struct {
	phase              phase
	remaining          time.Duration
	completedPomodoros int
	deskUp             bool
	paused             bool
}

func newState() state {
	s := state{}
	s.enterWork()
	return s
}

func (p phase) letter() string {
	switch p {
	case work:
		return "W"
	case shortBreak:
		return "B"
	default:
		return "LB"
	}
}

func (s *state) standing() bool {
	if s.phase == work {
		return s.deskUp
	}
	return !s.deskUp
}

func deskLabel(standing bool) string {
	if standing {
		return "Desk state: Standing"
	}
	return "Desk state: Sitting"
}

func (s *state) enterWork() {
	s.deskUp = !s.deskUp
	s.phase = work
	s.remaining = workDuration
}

func (s *state) phaseDuration() time.Duration {
	switch s.phase {
	case work:
		return workDuration
	case shortBreak:
		return shortBreakDuration
	default:
		return longBreakDuration
	}
}

func (s *state) progress() float64 {
	total := s.phaseDuration()
	return (total - s.remaining).Seconds() / total.Seconds()
}

func (s *state) trayTitle() string {
	total := int(s.remaining.Seconds())
	mins := total / 60
	secs := total % 60

	title := fmt.Sprintf("%s %02d:%02d", s.phase.letter(), mins, secs)
	if s.paused {
		title += " ❄"
	}
	return title
}

func (s *state) advance() {
	switch s.phase {
	case work:
		s.completedPomodoros++
		if s.completedPomodoros >= pomodorosPerCycle {
			s.phase = longBreak
			s.remaining = longBreakDuration
			notify("Long break", "15 minutes. Well done. "+deskLabel(!s.deskUp))
		} else {
			s.phase = shortBreak
			s.remaining = shortBreakDuration
			notify("Short break", "5 minutes. "+deskLabel(!s.deskUp))
		}
	case shortBreak, longBreak:
		if s.phase == longBreak {
			s.completedPomodoros = 0
		}
		s.enterWork()
		notify("Work", deskLabel(s.deskUp))
	}
}

func (s *state) skip() {
	switch s.phase {
	case work:
		s.phase = shortBreak
		s.remaining = shortBreakDuration
		notifyText("Short break", "5 minutes. "+deskLabel(!s.deskUp))
	case shortBreak, longBreak:
		if s.phase == longBreak {
			s.completedPomodoros = 0
		}
		s.enterWork()
		notifyText("Work", deskLabel(s.deskUp))
	}
}

func runTimer(cmds <-chan command, update func(uiUpdate)) {
	s := newState()
	notifyText("Work", deskLabel(s.deskUp))

	emit := func(pl string) {
		latestState.Store(s)
		update(s.toUpdate(pl))
	}
	emit("Pause")

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
				emit("Pause")
				continue
			}
			s.advance()
			emit("Pause")

		case cmd := <-cmds:
			switch cmd {
			case cmdTogglePause:
				s.paused = !s.paused
				emit(pauseLabel(s.paused))

			case cmdSkip:
				s.paused = false
				s.skip()
				emit("Pause")

			case cmdReset:
				s.completedPomodoros = 0
				s.paused = false
				s.enterWork()
				notify("Work", deskLabel(s.deskUp))
				emit("Pause")

			case cmdQuit:
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
