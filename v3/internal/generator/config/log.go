package config

// A Logger instance provides methods to format and report messages
// intended for the end user.
//
// All Logger methods may be called concurrently by its consumers.
type Logger interface {
	// Errorf reports an error message.
	Errorf(format string, a ...any)
	// Warningf reports a warning message.
	Warningf(format string, a ...any)
	// Infof reports an informational message.
	Infof(format string, a ...any)
	// Debugf reports a debug message.
	Debugf(format string, a ...any)
	// Statusf reports a status update (e.g. for updating a spinner label).
	Statusf(format string, a ...any)
}

// NullLogger discards all incoming messages.
var NullLogger Logger = nullLogger{}

type nullLogger struct{}

func (nullLogger) Errorf(format string, a ...any)   {}
func (nullLogger) Warningf(format string, a ...any) {}
func (nullLogger) Infof(format string, a ...any)    {}
func (nullLogger) Debugf(format string, a ...any)   {}
func (nullLogger) Statusf(format string, a ...any)  {}
