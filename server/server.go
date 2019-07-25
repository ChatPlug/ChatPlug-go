package main

import (
	"fmt"
	"github.com/feelfreelinux/ChatPlug/core"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "help")
	}

	switch os.Args[1] {
	case "start":
		start()
	case "install":
		install()
	case "help":
		help()
	default:
		help()
	}
}

func help() {
	fmt.Println(`ChatPlug is an extensible chat bridge

Usage:

	chatplug [...options] <command>

The commands are:

	help		shows this message
	start		launches the daemon
	install		installs a new service from git or tarball`)
	os.Exit(0)
}

func start() {
	app := core.NewApp()
	app.Init()
	app.RunHTTPServer()
}

func install() {
	//err := core.InstallService()
	//if err != nil {
	//	fmt.Printf("Error installing service: %s\n", err)
	//}
}
