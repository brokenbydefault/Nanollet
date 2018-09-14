package main

import (
	"fmt"
	"github.com/brokenbydefault/Nanollet/GUI"
	"runtime"
	"github.com/brokenbydefault/Nanollet/Storage"
)

//go:generate go run GUI/Generator/gen.go

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	if Storage.Configuration.DebugStatus {
		fmt.Println("The debug is enabled")
	}

	GUI.Start()
}
