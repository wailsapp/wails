package config

import (
	"fmt"

	"github.com/pterm/pterm"
)

// A Logger instance provides methods to format and report messages
// intended for the end user.
//
// All Logger methods may be called concurrently by its consumers.
type Logger interface {
	// Errorf should process its arguments as if they were passed to fmt.Sprintf
	// and report the resulting string to the user as an error message.
	Errorf(format string, a ...any)

	// Warningf should process its arguments as if they were passed to fmt.Sprintf
	// and report the resulting string to the user as a warning message.
	Warningf(format string, a ...any)

	// Infof should process its arguments as if they were passed to fmt.Sprintf
	// and report the resulting string to the user as an informational message.
	Infof(format string, a ...any)

	// Debugf should process its arguments as if they were passed to fmt.Sprintf
	// and report the resulting string to the user as a debug message.
	Debugf(format string, a ...any)

	// Statusf should process its arguments as if they were passed to fmt.Sprintf
	// and report the resulting string to the user as a status message.
	Statusf(format string, a ...any)
}

// NullLogger is a dummy Logger implementation
// that discards all incoming messages.
var NullLogger Logger = nullLogger{}

type nullLogger struct{}

func (nullLogger) Errorf(format string, a ...any)   {}
func (nullLogger) Warningf(format string, a ...any) {}
func (nullLogger) Infof(format string, a ...any)    {}
func (nullLogger) Debugf(format string, a ...any)   {}
func (nullLogger) Statusf(format string, a ...any)  {}

// DefaultPtermLogger returns a Logger implementation that writes
// to the default pterm printers for each logging level.
//
// If spinner is not nil, it is used to log status updates.
// The spinner must have been started already.
func DefaultPtermLogger(spinner *pterm.SpinnerPrinter) Logger {
	return &PtermLogger{
		&pterm.Error,
		&pterm.Warning,
		&pterm.Info,
		&pterm.Debug,
		spinner,
	}
}

// PtermLogger is a Logger implementation that writes to pterm printers.
// If any field is nil, PtermLogger discards all messages of that level.
type PtermLogger struct {
	Error   pterm.TextPrinter
	Warning pterm.TextPrinter
	Info    pterm.TextPrinter
	Debug   pterm.TextPrinter
	Spinner *pterm.SpinnerPrinter
}

func (logger *PtermLogger) Errorf(format string, a ...any) {
	if logger.Error != nil {
		logger.Error.Printfln(format, a...)
	}
}

func (logger *PtermLogger) Warningf(format string, a ...any) {
	if logger.Warning != nil {
		logger.Warning.Printfln(format, a...)
	}
}

func (logger *PtermLogger) Infof(format string, a ...any) {
	if logger.Info != nil {
		logger.Info.Printfln(format, a...)
	}
}

func (logger *PtermLogger) Debugf(format string, a ...any) {
	if logger.Debug != nil {
		logger.Debug.Printfln(format, a...)
	}
}

func (logger *PtermLogger) Statusf(format string, a ...any) {
	if logger.Spinner != nil {
		logger.Spinner.UpdateText(fmt.Sprintf(format, a...))
	}
}
