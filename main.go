package main

import (
	"embed"
	"vvorker/exec"
	kv "vvorker/ext/kv/src"
	"vvorker/services"

	"github.com/sirupsen/logrus"
)

//go:embed all:www/out/*
var fs embed.FS

// go:embed ext/ai/dist/index.js
var ExtAiScript string

func main() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	defer exec.ExecManager.ExitAllCmd()
	defer kv.Close()

	services.Run(fs)
}
