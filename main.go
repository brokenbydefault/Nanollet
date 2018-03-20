package main

import (
	"github.com/brokenbydefault/Nanollet/GUI"
	"github.com/brokenbydefault/Nanollet/RPC"
	"runtime"
	"github.com/brokenbydefault/Nanollet/Config"
	"fmt"
	"github.com/brokenbydefault/Nanollet/GUI/Storage"
)

//go:generate go run gen.go

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if Config.IsDebugEnabled() {
		fmt.Println("The debug is enabled")
	}

	Storage.Start()
	RPCClient.StartWebsocket()
	GUI.Start()
}