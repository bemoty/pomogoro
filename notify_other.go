//go:build !linux && !darwin

package main

func notifyText(title, body string) {}
