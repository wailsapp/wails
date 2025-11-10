//go:generate go run codegen/error_functions/main.go
package errs

import (
	"errors"
)

func Is(err error, errorType ErrorType) bool {
	if err == nil {
		return false
	}
	var general WailsError
	ok := errors.As(err, &general)
	if !ok {
		return false
	}
	return general.ErrorType() == errorType
}

func Cause(err error) error {
	if err == nil {
		return nil
	}

	type causer interface {
		Cause() error
	}

	if cause, ok := err.(causer); ok {
		if causeErr := cause.Cause(); causeErr != nil {
			return causeErr
		}
	}

	if unwrapErr := errors.Unwrap(err); unwrapErr != nil {
		return unwrapErr
	}

	return err
}

func Has(err error, errorType ErrorType) bool {
	if err == nil {
		return false
	}

	for {
		if Is(err, errorType) {
			return true
		}

		cause := Cause(err)
		if errors.Is(cause, err) || cause == nil {
			break
		}
		err = cause
	}

	return false
}
