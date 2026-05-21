//go:build !linux && !darwin

package main

type pidManager struct{}

func (p *pidManager) checkSingleInstance() error { return nil }
func (p *pidManager) write()                     {}
func (p *pidManager) release()                   {}

var pid = &pidManager{}

func daemonize() {}
