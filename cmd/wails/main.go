package main

import (
	"os"
	"os/exec"

	"github.com/wailsapp/wails/cmd"
)

// Create Logger
var logger = cmd.NewLogger()

// Create main app
var app = cmd.NewCli("wails", "A cli tool for building Wails applications.")

// Main!
func main() {
	err := app.Run()
	if err != nil {
		logger.Error(err.Error())
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}
