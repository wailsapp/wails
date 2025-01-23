package application

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
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

type PanicDetails struct {
	StackTrace     string
	Error          error
	Time           time.Time
	FullStackTrace string
}

func newPanicDetails(err error, trace string) *PanicDetails {
	return &PanicDetails{
		Error:          err,
		Time:           time.Now(),
		StackTrace:     trace,
		FullStackTrace: string(debug.Stack()),
	}
}

// handlePanic handles any panics
// Returns the error if there was one
func handlePanic(options ...handlePanicOptions) bool {
	// Try to recover
	e := recover()
	if e == nil {
		return false
	}

	// Get the error
	var err error
	if errPanic, ok := e.(error); ok {
		err = errPanic
	} else {
		err = fmt.Errorf("%v", e)
	}

	// Get the stack trace
	var stackTrace string
	skipEnd := 0
	if len(options) > 0 {
		skipEnd = options[0].skipEnd
	}
	stackTrace = getStackTrace(3, skipEnd)

	processPanic(newPanicDetails(err, stackTrace))
	return false
}

func processPanic(panicDetails *PanicDetails) {
	h := globalApplication.options.PanicHandler
	if h != nil {
		h(panicDetails)
		return
	}
	defaultPanicHandler(panicDetails)
}

func defaultPanicHandler(panicDetails *PanicDetails) {
	errorMessage := fmt.Sprintf("panic error: %s\n%s", panicDetails.Error.Error(), panicDetails.StackTrace)
	globalApplication.fatal(errorMessage)
}
