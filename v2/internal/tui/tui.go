// Package tui provides a unified terminal UI for the Wails CLI.
// It replaces the previous mix of pterm, wzshiming/ctc, and labstack/gommon/color
// with a single coherent system: glyph for interactive elements (spinners, progress)
// and direct ANSI codes for plain styled text.
package tui

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kungfusheep/glyph"
)

// ColourEnabled controls whether ANSI colors and TUI effects are used.
// Set to false via -nocolor flag, NO_COLOR env var, or non-TTY stdout.
var ColourEnabled = true

func init() {
	if os.Getenv("NO_COLOR") != "" {
		ColourEnabled = false
		return
	}
	if os.Getenv("TERM") == "dumb" {
		ColourEnabled = false
		return
	}
	// Disable color when stdout is not a terminal (CI, pipes)
	fi, err := os.Stdout.Stat()
	if err == nil && (fi.Mode()&os.ModeCharDevice) == 0 {
		ColourEnabled = false
	}
}

// SetNoColour disables all color and TUI output.
func SetNoColour() {
	ColourEnabled = false
}

// ansi returns ANSI-escaped text when colour is enabled, plain text otherwise.
func ansi(code, text string) string {
	if !ColourEnabled {
		return text
	}
	return "\033[" + code + "m" + text + "\033[0m"
}

// Fatal prints a FATAL message to stderr and exits with code 1.
func Fatal(message string) {
	fmt.Fprintln(os.Stderr, ansi("1;91", "FATAL")+" "+message)
	os.Exit(1)
}

// Error prints an error message.
func Error(message string) {
	fmt.Println(ansi("91", "  ✗")+" "+message)
}

// Success prints a success message.
func Success(message string) {
	fmt.Println(ansi("92", "  ✓")+" "+message)
}

// Info prints an informational message.
func Info(message string) {
	fmt.Println(ansi("94", "  →")+" "+message)
}

// Warning prints a warning message.
func Warning(message string) {
	fmt.Println(ansi("93", "  !")+" "+message)
}

// Println prints text followed by a newline.
func Println(text string) {
	fmt.Println(text)
}

// Printf prints formatted text.
func Printf(format string, args ...any) {
	fmt.Printf(format, args...)
}

// Section prints a styled section header.
func Section(title string) {
	fmt.Println()
	if ColourEnabled {
		fmt.Println(ansi("1;94", title))
		fmt.Println(ansi("2", strings.Repeat("─", 40)))
	} else {
		fmt.Println("=== " + title + " ===")
	}
}

// BulletPoint prints a formatted bullet point.
func BulletPoint(text string, args ...any) {
	msg := fmt.Sprintf(text, args...)
	if ColourEnabled {
		fmt.Println("  " + ansi("94", "•") + " " + msg)
	} else {
		fmt.Println("  - " + msg)
	}
}

// Table prints a two-column key-value table.
func Table(rows [][]string) {
	if len(rows) == 0 {
		return
	}
	maxKey := 0
	for _, row := range rows {
		if len(row) > 0 && len(row[0]) > maxKey {
			maxKey = len(row[0])
		}
	}
	for _, row := range rows {
		if len(row) < 2 {
			continue
		}
		key := row[0]
		val := row[1]
		if ColourEnabled {
			fmt.Printf("  %s%s  %s\n",
				ansi("2", key),
				strings.Repeat(" ", maxKey-len(key)),
				val,
			)
		} else {
			fmt.Printf("  %-*s  %s\n", maxKey, key, val)
		}
	}
}

// HeaderTable prints a table with the first row as a bold header.
func HeaderTable(rows [][]string) {
	if len(rows) == 0 {
		return
	}
	// Calculate column widths
	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	for i, row := range rows {
		line := "  "
		for j, cell := range row {
			w := 0
			if j < len(widths) {
				w = widths[j]
			}
			padded := cell + strings.Repeat(" ", w-len(cell))
			if i == 0 && ColourEnabled {
				line += ansi("1", padded) + "  "
			} else {
				line += padded + "  "
			}
		}
		fmt.Println(line)
		if i == 0 {
			sep := "  "
			for j := range row {
				w := 0
				if j < len(widths) {
					w = widths[j]
				}
				sep += strings.Repeat("-", w) + "  "
			}
			fmt.Println(sep)
		}
	}
}

// BoxedTable prints a section header followed by a table.
func BoxedTable(title string, rows [][]string) {
	if title != "" {
		Section(title)
	}
	Table(rows)
}

// spinnerFrames are the braille spinner animation frames.
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// WithSpinner runs fn while displaying an animated spinner.
// On completion, a success or failure indicator replaces the spinner.
// In non-TTY mode it prints a plain start/done message instead.
func WithSpinner(message string, fn func() error) error {
	if !ColourEnabled {
		fmt.Println(message + "...")
		err := fn()
		if err != nil {
			fmt.Println("  Failed: " + err.Error())
		} else {
			fmt.Println("  Done.")
		}
		return err
	}

	var (
		done   bool
		runErr error
		mu     sync.Mutex
		frame  int
	)

	displayMsg := message
	app := glyph.NewInlineApp()
	app.Height(1)
	app.SetView(
		glyph.HBox(
			glyph.Spinner(&frame).Frames(glyph.SpinnerBraille).FG(glyph.Cyan),
			glyph.SpaceW(1),
			glyph.Text(&displayMsg),
		),
	)

	go func() {
		for {
			mu.Lock()
			d := done
			mu.Unlock()
			if d {
				return
			}
			frame++
			app.RequestRender()
			time.Sleep(80 * time.Millisecond)
		}
	}()

	go func() {
		runErr = fn()
		mu.Lock()
		done = true
		mu.Unlock()
		app.Stop()
	}()

	_ = app.RunNonInteractive()

	if runErr != nil {
		fmt.Println(ansi("91", "  ✗") + " " + message)
	} else {
		fmt.Println(ansi("92", "  ✓") + " " + message)
	}
	return runErr
}

// Green returns green-colored text when colour is enabled.
func Green(text string) string { return ansi("92", text) }

// Red returns bright red text when colour is enabled.
func Red(text string) string { return ansi("91", text) }

// DarkRed returns dark red text when colour is enabled.
func DarkRed(text string) string { return ansi("31", text) }

// Yellow returns yellow text when colour is enabled.
func Yellow(text string) string { return ansi("93", text) }

// DarkYellow returns dark yellow text when colour is enabled.
func DarkYellow(text string) string { return ansi("33", text) }

// Blue returns blue text when colour is enabled.
func Blue(text string) string { return ansi("94", text) }

// DarkBlue returns dark blue text when colour is enabled.
func DarkBlue(text string) string { return ansi("34", text) }

// Cyan returns cyan text when colour is enabled.
func Cyan(text string) string { return ansi("96", text) }

// DarkCyan returns dark cyan text when colour is enabled.
func DarkCyan(text string) string { return ansi("36", text) }

// Magenta returns magenta text when colour is enabled.
func Magenta(text string) string { return ansi("95", text) }

// DarkMagenta returns dark magenta text when colour is enabled.
func DarkMagenta(text string) string { return ansi("35", text) }

// White returns bright white text when colour is enabled.
func White(text string) string { return ansi("97", text) }

// DarkWhite returns regular white text when colour is enabled.
func DarkWhite(text string) string { return ansi("37", text) }

// Black returns bright black (dark grey) text when colour is enabled.
func Black(text string) string { return ansi("90", text) }

// DarkBlack returns black text when colour is enabled.
func DarkBlack(text string) string { return ansi("30", text) }

// Gray returns grey text when colour is enabled.
func Gray(text string) string { return ansi("2", text) }

// Bold returns bold text when colour is enabled.
func Bold(text string) string { return ansi("1", text) }

// Rainbow returns text with cycling colors when colour is enabled.
func Rainbow(text string) string {
	if !ColourEnabled {
		return text
	}
	codes := []string{"91", "93", "92", "96", "94", "95"}
	var sb strings.Builder
	for i, ch := range text {
		sb.WriteString(ansi(codes[i%len(codes)], string(ch)))
	}
	return sb.String()
}
