package main

import (
	"fmt"
	"github.com/brokenbydefault/Nanollet/Config"
	"github.com/brokenbydefault/Nanollet/GUI"
	"github.com/brokenbydefault/Nanollet/RPC"
	"runtime"
)

//go:generate go run GUI/Generator/gen.go

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	if Config.Configuration().DebugStatus {
		fmt.Println("The debug is enabled")
	}

	RPCClient.StartWebsocket()

	GUI.Start()
}
