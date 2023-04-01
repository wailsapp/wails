package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v2/cmd/wails/internal"

	"github.com/wailsapp/wails/v2/internal/colour"

	"github.com/leaanthony/clir"
)

func banner(_ *clir.Cli) string {
	return fmt.Sprintf("%s %s",
		colour.Green("Wails CLI"),
		colour.DarkRed(internal.Version))
}

func fatal(message string) {
	printer := pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.FatalMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.FatalPrefixStyle,
			Text:  " FATAL ",
		},
	}
	printer.Println(message)
	os.Exit(1)
}

func printBulletPoint(text string, args ...any) {
	item := pterm.BulletListItem{
		Level: 2,
		Text:  text,
	}
	t, err := pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{item}).Srender()
	if err != nil {
		fatal(err.Error())
	}
	t = strings.Trim(t, "\n\r")
	pterm.Printfln(t, args...)
}

func printFooter() {
	printer := pterm.PrefixPrinter{
		MessageStyle: pterm.NewStyle(pterm.FgLightGreen),
		Prefix: pterm.Prefix{
			Style: pterm.NewStyle(pterm.FgRed, pterm.BgLightWhite),
			Text:  "♥ ",
		},
	}
	printer.Println("If Wails is useful to you or your company, please consider sponsoring the project:")
	pterm.Println("https://github.com/sponsors/leaanthony")
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
		pterm.Println(internal.Version)
		return nil
	})

	err = app.Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(err.Error())
		printFooter()
		os.Exit(1)
	}
}
