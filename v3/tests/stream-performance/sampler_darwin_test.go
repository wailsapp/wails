//go:build darwin

package main

import (
	"os"
	"testing"
)

func TestSamplerReadsLibproc(t *testing.T) {
	if _, err := webContentPids(); err != nil {
		t.Fatal(err)
	}
	name, ok := processName(os.Getpid())
	if !ok || name == "" {
		t.Fatalf("processName(self) = %q, %v", name, ok)
	}
	if size, ok := footprint(os.Getpid()); !ok || size == 0 {
		t.Fatalf("footprint(self) = %d, %v", size, ok)
	}
}
