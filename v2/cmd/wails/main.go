package main

import (
	"fmt"
	"os"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/cmd/wails/internal"
	"github.com/wailsapp/wails/v2/internal/tui"
)

func banner(_ *clir.Cli) string {
	return fmt.Sprintf("%s %s",
		tui.Green("Wails CLI"),
		tui.DarkRed(internal.Version))
}

func fatal(message string) {
	tui.Fatal(message)
}

func printBulletPoint(text string, args ...any) {
	tui.BulletPoint(text, args...)
}

func printFooter() {
	fmt.Println(tui.Green("♥") + " If Wails is useful to you or your company, please consider sponsoring the project:")
	fmt.Println("  https://github.com/sponsors/leaanthony")
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
	app.NewSubCommandFunction("dev", "Runs the application in development mode", devApplication)
	app.NewSubCommandFunction("doctor", "Diagnose your environment", diagnoseEnvironment)
	app.NewSubCommandFunction("init", "Initialises a new Wails project", initProject)
	app.NewSubCommandFunction("update", "Update the Wails CLI", update)

	show := app.NewSubCommand("show", "Shows various information")
	show.NewSubCommandFunction("releasenotes", "Shows the release notes for the current version", showReleaseNotes)

	generate := app.NewSubCommand("generate", "Code Generation Tools")
	generate.NewSubCommandFunction("module", "Generates a new Wails module", generateModule)
	generate.NewSubCommandFunction("template", "Generates a new Wails template", generateTemplate)

	command := app.NewSubCommand("version", "The Wails CLI version")
	command.Action(func() error {
		fmt.Println(internal.Version)
		return nil
	})

	err = app.Run()
	if err != nil {
		fmt.Println()
		tui.Error(err.Error())
		printFooter()
		os.Exit(1)
	}
}
