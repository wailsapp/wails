package cfd

import "errors"

var (
	ErrCancelled = errors.New("cancelled by user")
	ErrInvalidGUID = errors.New("guid cannot be nil")
	ErrEmptyFilters = errors.New("must specify at least one filter")
)
