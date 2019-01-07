package wails

type RuntimeLog struct {
}

func newRuntimeLog() *RuntimeLog {
	return &RuntimeLog{}
}

func (r *RuntimeLog) New(prefix string) *CustomLogger {
	return newCustomLogger(prefix)
}
