package main

import (
	"fmt"
	"github.com/wailsapp/wails/v2/cmd/wails/internal"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/show"
	"os"

	"github.com/wailsapp/wails/v2/internal/colour"

	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/update"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/dev"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/doctor"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/generate"
)

func fatal(message string) {
	println(message)
	os.Exit(1)
}

func banner(_ *clir.Cli) string {
	return fmt.Sprintf("%s %s",
		colour.Green("Wails CLI"),
		colour.DarkRed(internal.Version))
}

func printFooter() {
	println(colour.Green("\nIf Wails is useful to you or your company, please consider sponsoring the project:\nhttps://github.com/sponsors/leaanthony\n"))
}

func bool2Str(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

var app *clir.Cli

func main() {

	var err error

	app = clir.NewCli("Wails", "Go/HTML Appkit", internal.Version)

	app.SetBannerFunction(banner)
	defer printFooter()

	app.NewSubCommandFunction("build", "Builds the application", buildApplication)

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

	show.AddSubcommand(app, os.Stdout)

	err = update.AddSubcommand(app, os.Stdout, internal.Version)
	if err != nil {
		fatal(err.Error())
	}

	command := app.NewSubCommand("version", "The Wails CLI version")
	command.Action(func() error {
		println(internal.Version)
		return nil
	})

	err = app.Run()
	if err != nil {
		println("\n\nERROR: " + err.Error())
		printFooter()
		os.Exit(1)
	}
}
