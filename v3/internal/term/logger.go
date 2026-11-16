package term

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/generator/config"
)

// DefaultPtermLogger returns a Logger implementation that writes
// to the default pterm printers for each logging level.
//
// If spinner is not nil, it is used to log status updates.
// The spinner must have been started already.
func DefaultPtermLogger(spinner *pterm.SpinnerPrinter) config.Logger {
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
		// Prefix with [warn] to prevent GitHub Actions Go problem matcher
		// from treating warnings as errors (it matches file.go:line:col: patterns)
		logger.Warning.Printfln("[warn] "+format, a...)
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
