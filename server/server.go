package main

import (
	"github.com/feelfreelinux/ChatPlug/core"
)

func main() {
	app := core.NewApp()
	app.Init()
	app.RunHTTPServer()
}
