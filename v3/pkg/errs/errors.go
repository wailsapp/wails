//go:generate go run codegen/error_functions/main.go
package errs

type ErrorType string

type WailsError interface {
	Cause() error
	Error() string
	Msg() string
	ErrorType() ErrorType
}

const (
	InvalidWindowCallError      ErrorType = "Invalid window call"
	InvalidApplicationCallError ErrorType = "Invalid application call"
	InvalidBrowserCallError     ErrorType = "Invalid browser call"
	InvalidSystemCallError      ErrorType = "Invalid system call"
	InvalidScreensCallError     ErrorType = "Invalid screens call"
	InvalidDialogCallError      ErrorType = "Invalid dialog call"
	InvalidContextMenuCallError ErrorType = "Invalid context menu call"
	InvalidClipboardCallError   ErrorType = "Invalid clipboard call"
	InvalidBindingCallError     ErrorType = "Invalid binding call"
	BindingCallFailedError      ErrorType = "Binding call failed"
	InvalidEventsCallError      ErrorType = "Invalid events call"
	InvalidRuntimeCallError     ErrorType = "Invalid runtime call"
)
