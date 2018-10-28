package main

import (
	"github.com/brokenbydefault/Nanollet/GUI"
	"github.com/brokenbydefault/Nanollet/Storage"
	"runtime"
)

//go:generate go run GUI/Generator/gen.go

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	if Storage.Configuration.DebugStatus {
		print("The debug is enabled")
	}

	GUI.Start()
}
