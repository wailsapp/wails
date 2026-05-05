package application

import (
	"fmt"
	"os"
	"strings"
)

// FatalError instances are passed to the registered error handler
// in case of catastrophic, unrecoverable failures that require immediate termination.
// FatalError wraps the original error value in an informative message.
// The underlying error may be retrieved through the [FatalError.Unwrap] method.
type FatalError struct {
	err      error
	internal bool
}

// Internal returns true when the error was triggered from wails' internal code.
func (e *FatalError) Internal() bool {
	return e.internal
}

// Unwrap returns the original cause of the fatal error,
// for easy inspection using the [errors.As] API.
func (e *FatalError) Unwrap() error {
	return e.err
}

func (e *FatalError) Error() string {
	var buffer strings.Builder
	buffer.WriteString("\n\n******************************** FATAL *********************************\n")
	buffer.WriteString("*      There has been a catastrophic failure in your application.      *\n")
	if e.internal {
		buffer.WriteString("* Please report this error at https://github.com/wailsapp/wails/issues *\n")
	}
	buffer.WriteString("**************************** Error Details *****************************\n")
	buffer.WriteString(e.err.Error())
	buffer.WriteString("************************************************************************\n")
	return buffer.String()
}

func Fatal(message string, args ...any) {
	err := &FatalError{
		err:      fmt.Errorf(message, args...),
		internal: true,
	}

	if globalApplication != nil {
		globalApplication.handleError(err)
	} else {
		fmt.Println(err)
	}

	os.Exit(1)
}
