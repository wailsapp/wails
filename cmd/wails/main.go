package main

import (
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
	}
}
