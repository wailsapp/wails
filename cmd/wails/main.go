package main

import (
	"github.com/wailsapp/wails/cmd"
)

// Create Logger
var logger = cmd.NewLogger()

// Create main app
var app = cmd.NewCli("wails", "A cli tool for building Wails applications.")

// Prints the cli banner
func printBanner(app *cmd.Cli) error {
	logger.PrintBanner()
	return nil
}

// Main!
func main() {
	app.PreRun(printBanner)
	err := app.Run()
	if err != nil {
		logger.Error(err.Error())
	}
}
