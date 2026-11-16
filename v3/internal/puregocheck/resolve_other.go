//go:build !darwin && !linux

package main

import "errors"

func resolve(*report) error {
	return errors.New("resolve is only implemented for darwin and linux")
}
