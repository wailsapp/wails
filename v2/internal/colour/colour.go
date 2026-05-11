// Package colour provides backward-compatible color functions.
// All functions delegate to the tui package which is the single source of truth.
package colour

import (
	"strings"

	"github.com/wailsapp/wails/v2/internal/tui"
)

// ColourEnabled is kept for backward compatibility.
// Use tui.SetNoColour() or tui.ColourEnabled directly instead.
var ColourEnabled = true

func Col(text string) string     { return text }
func Yellow(text string) string  { return tui.Yellow(text) }
func Red(text string) string     { return tui.Red(text) }
func Blue(text string) string    { return tui.Blue(text) }
func Green(text string) string   { return tui.Green(text) }
func Cyan(text string) string    { return tui.Cyan(text) }
func Magenta(text string) string { return tui.Magenta(text) }
func White(text string) string   { return tui.White(text) }
func Black(text string) string   { return tui.Black(text) }

func DarkYellow(text string) string  { return tui.DarkYellow(text) }
func DarkRed(text string) string     { return tui.DarkRed(text) }
func DarkBlue(text string) string    { return tui.DarkBlue(text) }
func DarkGreen(text string) string   { return tui.Green(text) }
func DarkCyan(text string) string    { return tui.DarkCyan(text) }
func DarkMagenta(text string) string { return tui.DarkMagenta(text) }
func DarkWhite(text string) string   { return tui.DarkWhite(text) }
func DarkBlack(text string) string   { return tui.DarkBlack(text) }

var rainbowCols = []func(string) string{
	tui.Red, tui.Yellow, tui.Green, tui.Cyan, tui.Blue, tui.Magenta,
}

func Rainbow(text string) string {
	if !ColourEnabled {
		return text
	}
	var builder strings.Builder
	for i := 0; i < len(text); i++ {
		fn := rainbowCols[i%len(rainbowCols)]
		builder.WriteString(fn(text[i : i+1]))
	}
	return builder.String()
}
