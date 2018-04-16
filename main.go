package main

import (
	"fmt"
	"github.com/brokenbydefault/Nanollet/Config"
	"github.com/brokenbydefault/Nanollet/GUI"
	"github.com/brokenbydefault/Nanollet/GUI/Storage"
	"github.com/brokenbydefault/Nanollet/RPC"
	"runtime"
)

//go:generate go run GUI/Generator/gen.go

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	if Config.IsDebugEnabled() {
		fmt.Println("The debug is enabled")
	}

	Storage.Start()
	RPCClient.StartWebsocket()

	GUI.Unpack()
	GUI.Start()
}
