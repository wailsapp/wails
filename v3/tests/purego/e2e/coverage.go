package main

import (
	"fmt"
	"os"
	"runtime/coverage"
)

// flushCoverage writes the coverage data of a probe built with -cover into
// GOCOVERDIR. Go writes it on a normal exit, but the toolkit's quit path can
// end the process without one, so the probe flushes before it quits. Without
// -cover or GOCOVERDIR it does nothing.
func flushCoverage() {
	dir := os.Getenv("GOCOVERDIR")
	if dir == "" {
		return
	}
	if err := coverage.WriteMetaDir(dir); err != nil {
		fmt.Fprintln(os.Stderr, "coverage:", err)
	}
	if err := coverage.WriteCountersDir(dir); err != nil {
		fmt.Fprintln(os.Stderr, "coverage:", err)
	}
}

// exit flushes the coverage data, then exits.
func exit(code int) {
	flushCoverage()
	os.Exit(code)
}
