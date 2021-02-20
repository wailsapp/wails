package main

import (
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/internal/colour"

	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/update"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/build"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/debug"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/dev"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/doctor"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/generate"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/initialise"
)

func fatal(message string) {
	println(message)
	os.Exit(1)
}

func banner(_ *clir.Cli) string {
	return fmt.Sprintf("%s %s - Go/HTML Application Framework", colour.Yellow("Wails"), colour.DarkRed(version))
}

func main() {

	var err error

	app := clir.NewCli("Wails", "Go/HTML Application Framework", version)

	app.SetBannerFunction(banner)

	build.AddBuildSubcommand(app, os.Stdout)
	err = initialise.AddSubcommand(app, os.Stdout)
	if err != nil {
		fatal(err.Error())
	}

	err = debug.AddSubcommand(app, os.Stdout)
	if err != nil {
		fatal(err.Error())
	}
	err = doctor.AddSubcommand(app, os.Stdout)
	if err != nil {
		fatal(err.Error())
	}

	err = dev.AddSubcommand(app, os.Stdout)
	if err != nil {
		fatal(err.Error())
	}

	err = generate.AddSubcommand(app, os.Stdout)
	if err != nil {
		fatal(err.Error())
	}

	err = update.AddSubcommand(app, os.Stdout, version)
	if err != nil {
		fatal(err.Error())
	}

	err = app.Run()
	if err != nil {
		println("\n\nERROR: " + err.Error())
	}
}
