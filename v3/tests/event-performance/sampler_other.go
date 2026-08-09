//go:build !darwin && !linux && !windows

package main

import "errors"

// Memory metrics are unsupported off darwin, but the timing and ordering half
// of the harness still runs so engines can be compared later.

const samplerSupported = false

func webContentPids() ([]int, error) {
	return nil, errors.New("process enumeration unsupported on this platform")
}

func footprint(pid int) (uint64, bool) { return 0, false }
