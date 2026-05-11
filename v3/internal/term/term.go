// Package term provides a unified, beautifully styled terminal UI for the Wails CLI.
// It uses direct ANSI escape codes for text output and glyph for animated spinners.
package term

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kungfusheep/glyph"
	"github.com/wailsapp/wails/v3/internal/generator/config"
	"github.com/wailsapp/wails/v3/internal/version"
	xterm "golang.org/x/term"
)

// ─── Global state ─────────────────────────────────────────────────────

var (
	colourEnabled = true
	outputEnabled = true
	debugEnabled  = false
)

// ─── ANSI codes ────────────────────────────────────────────────────────

const (
	ansiReset   = "\033[0m"
	ansiBold    = "\033[1m"
	ansiDim     = "\033[2m"
	ansiUnder   = "\033[4m"
	ansiCyan    = "\033[96m"
	ansiGreen   = "\033[92m"
	ansiRed     = "\033[91m"
	ansiYellow  = "\033[93m"
	ansiBlue    = "\033[94m"
	ansiMagenta = "\033[95m"
)

func col(code, text string) string {
	if !colourEnabled {
		return text
	}
	return code + text + ansiReset
}

// ─── Terminal detection ────────────────────────────────────────────────

// IsTerminal reports whether stdout is an interactive terminal (not CI/pipe).
func IsTerminal() bool {
	return xterm.IsTerminal(int(os.Stdout.Fd())) && os.Getenv("CI") != "true"
}

// ─── Toggles ───────────────────────────────────────────────────────────

func DisableColor()  { colourEnabled = false }
func EnableOutput()  { outputEnabled = true }
func DisableOutput() { outputEnabled = false }
func EnableDebug()   { debugEnabled = true }
func DisableDebug()  { debugEnabled = false }

// ─── Header ────────────────────────────────────────────────────────────

// Header prints the Wails brand banner for the given command.
//
//	  ◆ wails  ─  Init Project
func Header(command string) {
	if !outputEnabled {
		return
	}
	mark := col(ansiBold+ansiCyan, "◆")
	brand := col(ansiBold, "wails")
	ver := col(ansiDim, "v"+version.String())
	sep := col(ansiDim, "─")
	cmd := col(ansiBold, command)
	fmt.Printf("\n  %s %s %s  %s %s\n\n", mark, brand, ver, sep, cmd)
}

// ─── Sections ──────────────────────────────────────────────────────────

// Section prints a styled subsection heading.
//
//	  ── System
func Section(title string) {
	if !outputEnabled {
		return
	}
	bar := col(ansiDim, "──")
	fmt.Printf("\n  %s %s\n\n", bar, col(ansiBold, title))
}

// ─── Message functions ─────────────────────────────────────────────────

// Success prints a green check-mark message.
func Success(input any) {
	if !outputEnabled {
		return
	}
	fmt.Println("  " + col(ansiGreen, "✓") + "  " + sprint(input))
}

// Error prints a red cross message.
func Error(input any) {
	if !outputEnabled {
		return
	}
	fmt.Println("  " + col(ansiRed, "✗") + "  " + sprint(input))
}

// Warning prints a yellow warning message.
func Warning(input any) {
	if !outputEnabled {
		return
	}
	fmt.Println("  " + col(ansiYellow, "⚠") + "  " + sprint(input))
}

// Info prints a blue informational message.
func Info(input any) {
	if !outputEnabled {
		return
	}
	fmt.Println("  " + col(ansiBlue, "↳") + "  " + sprint(input))
}

// Successf formats and prints a success message.
func Successf(format string, args ...any) {
	Success(fmt.Sprintf(format, args...))
}

// Infof formats and prints an informational message.
func Infof(input any, args ...any) {
	Info(fmt.Sprintf(sprint(input), args...))
}

// Warningf formats and prints a warning message.
func Warningf(input any, args ...any) {
	Warning(fmt.Sprintf(sprint(input), args...))
}

// Errorf formats and prints an error message.
func Errorf(format string, args ...any) {
	Error(fmt.Sprintf(format, args...))
}

// Println prints an indented line.
func Println(s string) {
	if !outputEnabled {
		return
	}
	fmt.Println("  " + s)
}

// Printf prints an indented formatted line.
func Printf(format string, args ...any) {
	if !outputEnabled {
		return
	}
	fmt.Printf("  "+format, args...)
}

// sprint formats any value to a string.
func sprint(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case error:
		return t.Error()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ─── Table ─────────────────────────────────────────────────────────────

// Table prints a clean two-column table with the key column in dim grey.
//
//	  Platform    linux/amd64
//	  Compiler    gcc
func Table(rows [][]string) {
	if !outputEnabled || len(rows) == 0 {
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
		key := col(ansiDim, row[0]) + strings.Repeat(" ", maxKey-len(row[0]))
		fmt.Printf("  %s  %s\n", key, row[1])
	}
}

// HeaderTable prints a table with the first row as a bold column header.
func HeaderTable(rows [][]string) {
	if !outputEnabled || len(rows) == 0 {
		return
	}
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
			if i == 0 {
				line += col(ansiBold, padded) + "  "
			} else {
				line += padded + "  "
			}
		}
		fmt.Println(strings.TrimRight(line, " "))
		if i == 0 {
			sep := "  "
			for j := range row {
				w := 0
				if j < len(widths) {
					w = widths[j]
				}
				sep += col(ansiDim, strings.Repeat("─", w)) + "  "
			}
			fmt.Println(strings.TrimRight(sep, " "))
		}
	}
}

// ─── Interactive ────────────────────────────────────────────────────────

// Confirm asks a yes/no question and returns true if the user answered yes.
func Confirm(prompt string) bool {
	fmt.Printf("  %s %s", col(ansiYellow, "?"), prompt+" [y/N] ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
		return answer == "y" || answer == "yes"
	}
	return false
}

// ─── Spinner ────────────────────────────────────────────────────────────

type spinnerState struct {
	app    *glyph.App
	msgPtr *string
	mu     sync.Mutex
	done   chan struct{}
	wg     sync.WaitGroup
}

// Spinner wraps a running glyph spinner. The zero value is safe (no-op).
type Spinner struct {
	state *spinnerState
}

// UpdateText updates the spinner's displayed label while it's running.
func (s Spinner) UpdateText(text string) {
	if s.state == nil || s.state.app == nil {
		return
	}
	s.state.mu.Lock()
	*s.state.msgPtr = text
	s.state.mu.Unlock()
	s.state.app.RequestRender()
}

// Success stops the spinner and prints a success message.
func (s Spinner) Success(msg string) {
	StopSpinner(s)
	Success(msg)
}

// Fail stops the spinner and prints an error message.
func (s Spinner) Fail(msg string) {
	StopSpinner(s)
	Error(msg)
}

// Logger returns a config.Logger that routes messages through this spinner.
func (s Spinner) Logger() config.Logger {
	return &spinnerLogger{s}
}

// StartSpinner starts an animated spinner with the given label.
// On non-TTY outputs (CI, pipes) it prints a plain info line instead.
func StartSpinner(text string) Spinner {
	if !IsTerminal() {
		if outputEnabled {
			fmt.Println("  " + col(ansiBlue, "↳") + "  " + text + "…")
		}
		return Spinner{}
	}

	msg := text
	frame := 0

	app := glyph.NewInlineApp()
	app.Height(1)
	app.SetView(
		glyph.HBox(
			glyph.SpaceW(2),
			glyph.Spinner(&frame).Frames(glyph.SpinnerBraille).FG(glyph.Cyan),
			glyph.SpaceW(1),
			glyph.Text(&msg),
		),
	)

	st := &spinnerState{
		app:    app,
		msgPtr: &msg,
		done:   make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				frame++
				app.RequestRender()
			case <-st.done:
				return
			}
		}
	}()

	st.wg.Add(1)
	go func() {
		defer st.wg.Done()
		_ = app.RunNonInteractive()
	}()

	return Spinner{state: st}
}

// StopSpinner stops the spinner and waits for its goroutines to exit.
func StopSpinner(s Spinner) {
	if s.state == nil || s.state.app == nil {
		return
	}
	select {
	case <-s.state.done:
	default:
		close(s.state.done)
	}
	s.state.app.Stop()
	s.state.wg.Wait()
}

// spinnerLogger implements config.Logger by routing to the Spinner.
type spinnerLogger struct{ s Spinner }

func (l *spinnerLogger) Errorf(format string, a ...any) {
	Error(fmt.Sprintf(format, a...))
}
func (l *spinnerLogger) Warningf(format string, a ...any) {
	Warning(fmt.Sprintf(format, a...))
}
func (l *spinnerLogger) Infof(format string, a ...any) {
	Info(fmt.Sprintf(format, a...))
}
func (l *spinnerLogger) Debugf(format string, a ...any) {
	if debugEnabled {
		Info(fmt.Sprintf(format, a...))
	}
}
func (l *spinnerLogger) Statusf(format string, a ...any) {
	l.s.UpdateText(fmt.Sprintf(format, a...))
}

// ─── Hyperlink ──────────────────────────────────────────────────────────

// Hyperlink returns an OSC 8 terminal hyperlink (clickable in supporting terminals).
func Hyperlink(url, text string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s%s%s\x1b]8;;\x1b\\",
		url, ansiUnder, text, ansiReset)
}

// ─── Inline color helpers ────────────────────────────────────────────────

// These are provided for callers that need to compose colored strings directly.
func Cyan(s string) string    { return col(ansiCyan, s) }
func Green(s string) string   { return col(ansiGreen, s) }
func Red(s string) string     { return col(ansiRed, s) }
func Yellow(s string) string  { return col(ansiYellow, s) }
func Blue(s string) string    { return col(ansiBlue, s) }
func Magenta(s string) string { return col(ansiMagenta, s) }
func Bold(s string) string    { return col(ansiBold, s) }
func Dim(s string) string     { return col(ansiDim, s) }
