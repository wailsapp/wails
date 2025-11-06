//go:generate go run codegen/error_functions/main.go
//nolint:depguard // we want to use pkg/errors only here, but nowhere else
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
	type causer interface {
		Cause() error
	}

	if err != nil {
		cause, ok := err.(causer)
		if !ok {
			return err
		}
		causeErr := cause.Cause()
		if causeErr != nil {
			return causeErr
		}
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
