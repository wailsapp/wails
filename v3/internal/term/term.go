package term

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/generator/config"
	"github.com/wailsapp/wails/v3/internal/version"
	"golang.org/x/term"
)

func Header(header string) {
	// Print Wails with the current version in white on red background with the header in white with a green background
	pterm.BgLightRed.Print(pterm.LightWhite(" Wails (" + version.String() + ") "))
	pterm.BgLightGreen.Println(pterm.LightWhite(" " + header + " "))
}

func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd())) && (os.Getenv("CI") != "true")
}

type Spinner struct {
	spinner *pterm.SpinnerPrinter
}

func (s Spinner) Logger() config.Logger {
	return config.DefaultPtermLogger(s.spinner)
}

func StartSpinner(text string) Spinner {
	if !IsTerminal() {
		return Spinner{}
	}
	spinner, err := pterm.DefaultSpinner.Start(text)
	if err != nil {
		return Spinner{}
	}
	spinner.RemoveWhenDone = true
	return Spinner{spinner}
}

func StopSpinner(s Spinner) {
	if s.spinner != nil {
		_ = s.spinner.Stop()
	}
}

func output(input any, printer pterm.PrefixPrinter, args ...any) {
	switch v := input.(type) {
	case string:
		printer.Printfln(input.(string), args...)
	case error:
		printer.Println(v.Error())
	default:
		printer.Printfln("%v", v)
	}
}

func Info(input any) {
	output(input, pterm.Info)
}

func Infof(input any, args ...interface{}) {
	output(input, pterm.Info, args...)
}

func Warning(input any) {
	output(input, pterm.Warning)
}

func Warningf(input any, args ...any) {
	output(input, pterm.Warning, args...)
}

func Error(input any) {
	output(input, pterm.Error)
}

func Success(input any) {
	output(input, pterm.Success)
}

func Section(s string) {
	style := pterm.NewStyle(pterm.BgDefault, pterm.FgLightBlue, pterm.Bold)
	style.Println("\n# " + s + " \n")
}

func DisableColor() {
	pterm.DisableColor()
}

func EnableOutput() {
	pterm.EnableOutput()
}

func DisableOutput() {
	pterm.DisableOutput()
}

func EnableDebug() {
	pterm.EnableDebugMessages()
}

func DisableDebug() {
	pterm.DisableDebugMessages()
}

func Println(s string) {
	pterm.Println(s)
}

func Hyperlink(url, text string) string {
	// OSC 8 sequence to start a clickable link
	linkStart := "\x1b]8;;"

	// OSC 8 sequence to end a clickable link
	linkEnd := "\x1b]8;;\x1b\\"

	// ANSI escape code for underline
	underlineStart := "\x1b[4m"

	// ANSI escape code to reset text formatting
	resetFormat := "\x1b[0m"

	return fmt.Sprintf("%s%s%s%s%s%s%s", linkStart, url, "\x1b\\", underlineStart, text, resetFormat, linkEnd)
}
