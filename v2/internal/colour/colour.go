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

// apply returns plain text when the legacy ColourEnabled flag is false,
// otherwise it delegates to the tui helper fn.
func apply(fn func(string) string, text string) string {
	if !ColourEnabled {
		return text
	}
	return fn(text)
}

func Col(text string) string     { return text }
func Yellow(text string) string  { return apply(tui.Yellow, text) }
func Red(text string) string     { return apply(tui.Red, text) }
func Blue(text string) string    { return apply(tui.Blue, text) }
func Green(text string) string   { return apply(tui.Green, text) }
func Cyan(text string) string    { return apply(tui.Cyan, text) }
func Magenta(text string) string { return apply(tui.Magenta, text) }
func White(text string) string   { return apply(tui.White, text) }
func Black(text string) string   { return apply(tui.Black, text) }

func DarkYellow(text string) string  { return apply(tui.DarkYellow, text) }
func DarkRed(text string) string     { return apply(tui.DarkRed, text) }
func DarkBlue(text string) string    { return apply(tui.DarkBlue, text) }
func DarkGreen(text string) string   { return apply(tui.Green, text) }
func DarkCyan(text string) string    { return apply(tui.DarkCyan, text) }
func DarkMagenta(text string) string { return apply(tui.DarkMagenta, text) }
func DarkWhite(text string) string   { return apply(tui.DarkWhite, text) }
func DarkBlack(text string) string   { return apply(tui.DarkBlack, text) }

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
