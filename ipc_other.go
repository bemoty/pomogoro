//go:build !linux && !darwin

package main

import (
	"net"
	"sync/atomic"
)

var latestState atomic.Value
var ipcListener net.Listener

func listenIPC(_ chan<- command) {}
func clientCmd(_ string)         {}
