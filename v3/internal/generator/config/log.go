package config

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
