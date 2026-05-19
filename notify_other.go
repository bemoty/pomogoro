//go:build !linux && !darwin && !windows

package main

func notifyText(title, body string) {}
