//go:build !linux && !darwin

package main

import "errors"

func newDriver() (driver, error) { return nil, errors.New("no e2e input driver for this OS") }
