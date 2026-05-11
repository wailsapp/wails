package clilogger

import (
	"fmt"
	"io"
	"os"

	"github.com/wailsapp/wails/v2/internal/tui"
)

// CLILogger is used by the cli
type CLILogger struct {
	Writer io.Writer
	mute   bool
}

// New cli logger
func New(writer io.Writer) *CLILogger {
	return &CLILogger{
		Writer: writer,
	}
}

// Mute sets whether the logger should be muted
func (c *CLILogger) Mute(value bool) {
	c.mute = value
}

// Print works like Printf
func (c *CLILogger) Print(message string, args ...interface{}) {
	if c.mute {
		return
	}

	_, err := fmt.Fprintf(c.Writer, message, args...)
	if err != nil {
		c.Fatal("%s", "FATAL: "+err.Error())
	}
}

// Println works like Printf but with a line ending
func (c *CLILogger) Println(message string, args ...interface{}) {
	if c.mute {
		return
	}
	temp := fmt.Sprintf(message, args...)
	_, err := fmt.Fprintln(c.Writer, temp)
	if err != nil {
		c.Fatal("%s", "FATAL: "+err.Error())
	}
}

// Fatal prints the given message then aborts
func (c *CLILogger) Fatal(message string, args ...interface{}) {
	temp := fmt.Sprintf(message, args...)
	_, err := fmt.Fprintln(c.Writer, tui.Red("FATAL: "+temp))
	if err != nil {
		println(tui.Red("FATAL: " + err.Error()))
	}
	os.Exit(1)
}
