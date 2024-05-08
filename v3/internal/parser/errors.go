package parser

import (
	"fmt"
	"sync"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/parser/config"
)

// ErrorReport accumulates and logs error
// and warning messages, with deduplication.
//
// It implements the error interface; the Error method
// returns a report counting messages emitted so far.
type ErrorReport struct {
	logger config.Logger

	mu       sync.Mutex
	warnings map[string]bool
	errors   map[string]bool
}

// NewError report initialises an ErrorReport instance
// with the provided Logger implementation.
//
// If logger is nil, messages will be accumulated but not logged.
func NewErrorReport(logger config.Logger) *ErrorReport {
	if logger == nil {
		logger = config.NullLogger
	}

	return &ErrorReport{
		logger:   logger,
		warnings: make(map[string]bool),
		errors:   make(map[string]bool),
	}
}

// Error returns a string reporting the number
// of errors and warnings emitted so far.
func (report *ErrorReport) Error() string {
	report.mu.Lock()
	defer report.mu.Unlock()

	if len(report.errors) > 0 && len(report.warnings) == 0 {
		return fmt.Sprintf("%d errors emitted", len(report.errors))
	} else if len(report.errors) == 0 && len(report.warnings) > 0 {
		return fmt.Sprintf("%d warnings emitted", len(report.warnings))
	} else if len(report.errors) > 0 && len(report.warnings) > 0 {
		return fmt.Sprintf("%d errors and %d warnings emitted", len(report.errors), len(report.warnings))
	} else {
		return "no errors or warnings emitted"
	}
}

// HasErrors returns true if at least one error has been added to the report.
func (report *ErrorReport) HasErrors() bool {
	report.mu.Lock()
	result := len(report.errors) > 0
	report.mu.Unlock()
	return result
}

// HasWarnings returns true if at least one warning has been added to the report.
func (report *ErrorReport) HasWarnings() bool {
	report.mu.Lock()
	result := len(report.warnings) > 0
	report.mu.Unlock()
	return result
}

// Errors returns the list of error messages
// that have been added to the report.
// The order is randomised.
func (report *ErrorReport) Errors() []string {
	report.mu.Lock()
	defer report.mu.Unlock()

	return lo.Keys(report.errors)
}

// Warnings returns the list of warning messages
// that have been added to the report.
// The order is randomised.
func (report *ErrorReport) Warnings() []string {
	report.mu.Lock()
	defer report.mu.Unlock()

	return lo.Keys(report.warnings)
}

// Errorf formats an error message and adds it to the report.
// If not already present, the message is forwarded
// to the logger instance provided during initialisation.
func (report *ErrorReport) Errorf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)

	report.mu.Lock()
	defer report.mu.Unlock()

	present := report.errors[msg]
	report.errors[msg] = true

	if !present {
		report.logger.Errorf(format, a...)
	}
}

// Warningf formats an error message and adds it to the report.
// If not already present, the message is forwarded
// to the logger instance provided during initialisation.
func (report *ErrorReport) Warningf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)

	report.mu.Lock()
	defer report.mu.Unlock()

	present := report.warnings[msg]
	report.warnings[msg] = true

	if !present {
		report.logger.Warningf(format, a...)
	}
}

// Infof forwards the given informational message
// to the logger instance provided during initialisation.
//
// This method is here just for convenience and performs no deduplication.
func (report *ErrorReport) Infof(format string, a ...any) {
	report.logger.Infof(format, a...)
}

// Debugf forwards the given informational message
// to the logger instance provided during initialisation.
//
// This method is here just for convenience and performs no deduplication.
func (report *ErrorReport) Debugf(format string, a ...any) {
	report.logger.Debugf(format, a...)
}
