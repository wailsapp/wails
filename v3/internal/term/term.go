package term

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/generator/config"
	"github.com/wailsapp/wails/v3/internal/version"
	"golang.org/x/term"
	"os"
)

func Header(header string) {
	// Print Wails with the current version in white on red background with the header in white with a green background
	pterm.BgLightRed.Print(pterm.LightWhite(" Wails (" + version.String() + ") "))
	pterm.BgLightGreen.Println(pterm.LightWhite(" " + header + " "))
}

func Infof(format string, args ...interface{}) {
	pterm.Info.Printf(format, args...)
}
func Infofln(format string, args ...interface{}) {
	pterm.Info.Printfln(format, args...)
}

func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd())) && (os.Getenv("CI") != "true")
}

type Spinner struct {
	spinner *pterm.SpinnerPrinter
}

func (s *Spinner) Logger() config.Logger {
	if s == nil {
		return nil
	}
	return config.DefaultPtermLogger(s.spinner)
}

func StartSpinner(text string) *Spinner {
	if !IsTerminal() {
		return nil
	}
	spin, err := pterm.DefaultSpinner.Start(text)
	if err != nil {
		return nil
	}
	return &Spinner{
		spinner: spin,
	}
}

func StopSpinner(s *Spinner) {
	if s == nil {
		return
	}
	_ = s.spinner.Stop()
}

func output(input any, printer pterm.PrefixPrinter, args ...any) {
	switch v := input.(type) {
	case string:
		printer.Println(fmt.Sprintf(input.(string), args...))
	case error:
		printer.Println(v.Error())
	default:
		printer.Printfln("%v", v)
	}
}

func Warning(input any) {
	output(input, pterm.Warning)
}

func Warningf(input any, args ...any) {
	output(input, pterm.Warning, args)
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

func Println(s string) {
	pterm.Println(s)
}
