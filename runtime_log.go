package wails

// RuntimeLog exposes the logging interface to the runtime
type RuntimeLog struct {
}

func newRuntimeLog() *RuntimeLog {
	return &RuntimeLog{}
}

// New creates a new logger
func (r *RuntimeLog) New(prefix string) *CustomLogger {
	return newCustomLogger(prefix)
}
