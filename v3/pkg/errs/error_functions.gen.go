package errs

import (
	"fmt"
)

type wailsError struct {
	cause     error
	msg       string
	errorType ErrorType
}

func (w *wailsError) Cause() error { return w.cause }
func (w *wailsError) Error() string {
	errMsg := fmt.Sprintf("%s: %s", w.errorType, w.msg)
	if w.cause != nil {
		return fmt.Sprintf("%s: %s", errMsg, w.cause.Error())
	}
	return errMsg
}
func (w *wailsError) Msg() string          { return w.msg }
func (w *wailsError) ErrorType() ErrorType { return w.errorType }

func NewInvalidWindowCallErrorf(message string, args ...any) error {
	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     nil,
		msg:       msg,
		errorType: InvalidWindowCallError,
	}
}

func WrapInvalidWindowCallErrorf(err error, message string, args ...any) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     err,
		msg:       msg,
		errorType: InvalidWindowCallError,
	}
}

func IsInvalidWindowCallError(err error) bool {
	return Is(err, InvalidWindowCallError)
}

func HasInvalidWindowCallError(err error) bool {
	return Has(err, InvalidWindowCallError)
}

func NewInvalidApplicationCallErrorf(message string, args ...any) error {
	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     nil,
		msg:       msg,
		errorType: InvalidApplicationCallError,
	}
}

func WrapInvalidApplicationCallErrorf(err error, message string, args ...any) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     err,
		msg:       msg,
		errorType: InvalidApplicationCallError,
	}
}

func IsInvalidApplicationCallError(err error) bool {
	return Is(err, InvalidApplicationCallError)
}

func HasInvalidApplicationCallError(err error) bool {
	return Has(err, InvalidApplicationCallError)
}

func NewInvalidBrowserCallErrorf(message string, args ...any) error {
	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     nil,
		msg:       msg,
		errorType: InvalidBrowserCallError,
	}
}

func WrapInvalidBrowserCallErrorf(err error, message string, args ...any) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     err,
		msg:       msg,
		errorType: InvalidBrowserCallError,
	}
}

func IsInvalidBrowserCallError(err error) bool {
	return Is(err, InvalidBrowserCallError)
}

func HasInvalidBrowserCallError(err error) bool {
	return Has(err, InvalidBrowserCallError)
}

func NewInvalidSystemCallErrorf(message string, args ...any) error {
	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     nil,
		msg:       msg,
		errorType: InvalidSystemCallError,
	}
}

func WrapInvalidSystemCallErrorf(err error, message string, args ...any) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     err,
		msg:       msg,
		errorType: InvalidSystemCallError,
	}
}

func IsInvalidSystemCallError(err error) bool {
	return Is(err, InvalidSystemCallError)
}

func HasInvalidSystemCallError(err error) bool {
	return Has(err, InvalidSystemCallError)
}

func NewInvalidScreensCallErrorf(message string, args ...any) error {
	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     nil,
		msg:       msg,
		errorType: InvalidScreensCallError,
	}
}

func WrapInvalidScreensCallErrorf(err error, message string, args ...any) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     err,
		msg:       msg,
		errorType: InvalidScreensCallError,
	}
}

func IsInvalidScreensCallError(err error) bool {
	return Is(err, InvalidScreensCallError)
}

func HasInvalidScreensCallError(err error) bool {
	return Has(err, InvalidScreensCallError)
}

func NewInvalidDialogCallErrorf(message string, args ...any) error {
	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     nil,
		msg:       msg,
		errorType: InvalidDialogCallError,
	}
}

func WrapInvalidDialogCallErrorf(err error, message string, args ...any) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     err,
		msg:       msg,
		errorType: InvalidDialogCallError,
	}
}

func IsInvalidDialogCallError(err error) bool {
	return Is(err, InvalidDialogCallError)
}

func HasInvalidDialogCallError(err error) bool {
	return Has(err, InvalidDialogCallError)
}

func NewInvalidContextMenuCallErrorf(message string, args ...any) error {
	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     nil,
		msg:       msg,
		errorType: InvalidContextMenuCallError,
	}
}

func WrapInvalidContextMenuCallErrorf(err error, message string, args ...any) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     err,
		msg:       msg,
		errorType: InvalidContextMenuCallError,
	}
}

func IsInvalidContextMenuCallError(err error) bool {
	return Is(err, InvalidContextMenuCallError)
}

func HasInvalidContextMenuCallError(err error) bool {
	return Has(err, InvalidContextMenuCallError)
}

func NewInvalidClipboardCallErrorf(message string, args ...any) error {
	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     nil,
		msg:       msg,
		errorType: InvalidClipboardCallError,
	}
}

func WrapInvalidClipboardCallErrorf(err error, message string, args ...any) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     err,
		msg:       msg,
		errorType: InvalidClipboardCallError,
	}
}

func IsInvalidClipboardCallError(err error) bool {
	return Is(err, InvalidClipboardCallError)
}

func HasInvalidClipboardCallError(err error) bool {
	return Has(err, InvalidClipboardCallError)
}

func NewInvalidBindingCallErrorf(message string, args ...any) error {
	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     nil,
		msg:       msg,
		errorType: InvalidBindingCallError,
	}
}

func WrapInvalidBindingCallErrorf(err error, message string, args ...any) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     err,
		msg:       msg,
		errorType: InvalidBindingCallError,
	}
}

func IsInvalidBindingCallError(err error) bool {
	return Is(err, InvalidBindingCallError)
}

func HasInvalidBindingCallError(err error) bool {
	return Has(err, InvalidBindingCallError)
}

func NewBindingCallFailedErrorf(message string, args ...any) error {
	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     nil,
		msg:       msg,
		errorType: BindingCallFailedError,
	}
}

func WrapBindingCallFailedErrorf(err error, message string, args ...any) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     err,
		msg:       msg,
		errorType: BindingCallFailedError,
	}
}

func IsBindingCallFailedError(err error) bool {
	return Is(err, BindingCallFailedError)
}

func HasBindingCallFailedError(err error) bool {
	return Has(err, BindingCallFailedError)
}

func NewInvalidEventsCallErrorf(message string, args ...any) error {
	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     nil,
		msg:       msg,
		errorType: InvalidEventsCallError,
	}
}

func WrapInvalidEventsCallErrorf(err error, message string, args ...any) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     err,
		msg:       msg,
		errorType: InvalidEventsCallError,
	}
}

func IsInvalidEventsCallError(err error) bool {
	return Is(err, InvalidEventsCallError)
}

func HasInvalidEventsCallError(err error) bool {
	return Has(err, InvalidEventsCallError)
}

func NewInvalidRuntimeCallErrorf(message string, args ...any) error {
	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     nil,
		msg:       msg,
		errorType: InvalidRuntimeCallError,
	}
}

func WrapInvalidRuntimeCallErrorf(err error, message string, args ...any) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(message, args...)
	return &wailsError{
		cause:     err,
		msg:       msg,
		errorType: InvalidRuntimeCallError,
	}
}

func IsInvalidRuntimeCallError(err error) bool {
	return Is(err, InvalidRuntimeCallError)
}

func HasInvalidRuntimeCallError(err error) bool {
	return Has(err, InvalidRuntimeCallError)
}
