//go:build !linux && !darwin

package main

func daemonize()                  {}
func checkSingleInstance() error  { return nil }
func writePID()                   {}
func releasePID()                 {}
