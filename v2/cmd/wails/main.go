package main

import (
	"os"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/build"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/dev"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/doctor"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/initialise"
)

func fatal(message string) {
	println(message)
	os.Exit(1)
}

func main() {

	var err error
	version := "v2.0.0-alpha"

	app := clir.NewCli("Wails", "Go/HTML Application Framework", version)

	build.AddBuildSubcommand(app)
	err = initialise.AddSubcommand(app)
	if err != nil {
		fatal(err.Error())
	}
	err = doctor.AddSubcommand(app)
	if err != nil {
		fatal(err.Error())
	}

	err = dev.AddSubcommand(app)
	if err != nil {
		fatal(err.Error())
	}

	err = app.Run()
	if err != nil {
		println("\n\nERROR: " + err.Error())
	}
}
