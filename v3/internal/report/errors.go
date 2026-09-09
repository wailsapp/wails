package report

import "errors"

// reportedError preserves the cause while marking a diagnostic already shown
// by a Wails reporter. Build and dev failures share the same CLI exit contract.
type reportedError struct{ err error }

func (e reportedError) Error() string { return e.err.Error() }
func (e reportedError) Unwrap() error { return e.err }
func IsReported(err error) bool       { var rendered reportedError; return errors.As(err, &rendered) }
func MarkReported(err error) error {
	if err == nil || IsReported(err) {
		return err
	}
	return reportedError{err}
}
