package main

import (
	"embed"
	"fmt"
	"vvorker/conf"
	"vvorker/exec"
	kv "vvorker/ext/kv/src"
	"vvorker/services"

	"github.com/sirupsen/logrus"
)

//go:embed all:admin/dist/*
var fs embed.FS

//go:embed VERSION.txt
var Version string

func printBanner() {
	banner := `
╔═════════════════════════════════════════════════════════════════════════════════════════╗
║                                                                                         ║
║     ___      ___ ___      ___ ________  ________  ___  __    _______   ________         ║
║    |\  \    /  /|\  \    /  /|\   __  \|\   __  \|\  \|\  \ |\  ___ \ |\   __  \        ║
║    \ \  \  /  / | \  \  /  / | \  \|\  \ \  \|\  \ \  \/  /|\ \   __/|\ \  \|\  \       ║
║     \ \  \/  / / \ \  \/  / / \ \  \\\  \ \   _  _\ \   ___  \ \  \_|/_\ \   _  _\      ║
║      \ \    / /   \ \    / /   \ \  \\\  \ \  \\  \\ \  \\ \  \ \  \_|\ \ \  \\  \|     ║
║       \ \__/ /     \ \__/ /     \ \_______\ \__\\ _\\ \__\\ \__\ \_______\ \__\\ _\     ║
║        \|__|/       \|__|/       \|_______|\|__|\|__|\|__| \|__|\|_______|\|__|\|__|    ║
║                                                                                         ║
║                                                                                         ║
                                      VERSION: %s                                     
║                                       BY: Dicarne                                       ║
║                                                                                         ║
╚═════════════════════════════════════════════════════════════════════════════════════════╝
`
	fmt.Printf(banner, Version)
	fmt.Println()
}

func main() {

	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	printBanner()
	defer exec.ExecManager.ExitAllCmd()
	defer kv.Close()
	conf.Version = Version
	services.Run(fs)
}
