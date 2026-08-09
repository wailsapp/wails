//go:build !darwin && !linux

package main

// Memory metrics are unsupported off darwin, but the timing and ordering half
// of the harness still runs so engines can be compared later.

const samplerSupported = false

func webContentPids() []int { return nil }

func footprint(pid int) (uint64, bool) { return 0, false }
