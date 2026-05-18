package main

import (
	_ "embed"
	"os"
	"sync"
)

//go:embed assets/ding.ogg
var dingSound []byte

var (
	dingOnce sync.Once
	dingPath string
)

func ensureDing() string {
	dingOnce.Do(func() {
		f, err := os.CreateTemp("", "pomogoro-ding.ogg")
		if err != nil {
			return
		}
		defer f.Close()
		f.Write(dingSound)
		dingPath = f.Name()
	})
	return dingPath
}
