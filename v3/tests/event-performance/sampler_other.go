//go:build !darwin && !linux && !windows

package main

import "errors"

// Fallback for platforms with no sampler of their own — darwin, linux and
// windows each have one. Memory metrics are unavailable here, but the timing
// and ordering half of the harness still runs, so engines can be compared.

const samplerSupported = false

func webContentPids() ([]int, error) {
	return nil, errors.New("process enumeration unsupported on this platform")
}

func footprint(pid int) (uint64, bool) { return 0, false }
