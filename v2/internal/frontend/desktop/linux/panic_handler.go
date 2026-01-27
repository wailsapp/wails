//go:build linux
// +build linux

package linux

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/options"
)

func getStackTrace(skipStart int, skipEnd int) string {
	// Get all program counters first
	pc := make([]uintptr, 32)
	n := runtime.Callers(skipStart+1, pc)
	if n == 0 {
		return ""
	}

	pc = pc[:n]
	frames := runtime.CallersFrames(pc)

	// Collect all frames first
	var allFrames []runtime.Frame
	for {
		frame, more := frames.Next()
		allFrames = append(allFrames, frame)
		if !more {
			break
		}
	}

	// Remove frames from the end
	if len(allFrames) > skipEnd {
		allFrames = allFrames[:len(allFrames)-skipEnd]
	}

	// Build the output string
	var builder strings.Builder
	for _, frame := range allFrames {
		fmt.Fprintf(&builder, "%s\n\tat %s:%d\n",
			frame.Function, frame.File, frame.Line)
	}
	return builder.String()
}

type handlePanicOptions struct {
	skipEnd int
}

func newPanicDetails(err error, trace string) *options.PanicDetails {
	return &options.PanicDetails{
		Error:          err,
		Time:           time.Now(),
		StackTrace:     trace,
		FullStackTrace: string(debug.Stack()),
	}
}

// handlePanic recovers from panics and processes them through the configured handler.
// Returns true if a panic was recovered.
func handlePanic(handler options.PanicHandler, logger interface{ Error(string, ...interface{}) }, opts ...handlePanicOptions) bool {
	// Try to recover
	e := recover()
	if e == nil {
		return false
	}

	// Get the error
	err, ok := e.(error)
	if !ok {
		err = fmt.Errorf("%v", e)
	}

	// Get the stack trace
	var stackTrace string
	skipEnd := 0
	if len(opts) > 0 {
		skipEnd = opts[0].skipEnd
	}
	stackTrace = getStackTrace(3, skipEnd)

	panicDetails := newPanicDetails(err, stackTrace)

	// Use custom handler if provided
	if handler != nil {
		handler(panicDetails)
		return true
	}

	// Default behavior: log the panic
	if logger != nil {
		logger.Error("panic error: %v\n%s", panicDetails.Error, panicDetails.StackTrace)
	}
	return true
}
