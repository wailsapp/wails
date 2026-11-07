//go:build !linux && !darwin && !windows

package main

import "errors"

func screenshot(string) error { return errors.New("no screenshots on this OS") }
