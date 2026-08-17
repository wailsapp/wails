//go:build darwin

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// benchmarkConfig is deliberately private to the example. These controls make
// the demo repeatable under a process profiler without enlarging the proposed
// NativeWindow API or changing normal application behaviour.
type benchmarkConfig struct {
	hidden        bool
	documentBytes int
	autoQuit      time.Duration
	readyFile     string
	forceGC       bool
}

func loadBenchmarkConfig() (benchmarkConfig, error) {
	result := benchmarkConfig{
		hidden:    os.Getenv("NATIVE_EDITOR_BENCH_HIDDEN") == "1",
		readyFile: os.Getenv("NATIVE_EDITOR_BENCH_READY_FILE"),
		forceGC:   os.Getenv("NATIVE_EDITOR_BENCH_FORCE_GC") == "1",
	}
	if value := os.Getenv("NATIVE_EDITOR_BENCH_DOCUMENT_BYTES"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			return result, fmt.Errorf("NATIVE_EDITOR_BENCH_DOCUMENT_BYTES must be a non-negative integer")
		}
		result.documentBytes = parsed
	}
	if value := os.Getenv("NATIVE_EDITOR_BENCH_AUTO_QUIT_MS"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			return result, fmt.Errorf("NATIVE_EDITOR_BENCH_AUTO_QUIT_MS must be a positive integer")
		}
		result.autoQuit = time.Duration(parsed) * time.Millisecond
	}
	return result, nil
}
