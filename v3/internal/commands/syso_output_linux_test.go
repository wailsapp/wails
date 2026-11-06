package commands

import (
	"errors"
	"os"
	"syscall"
	"testing"
)

func TestGenerateSysoReportsWriteFailure(t *testing.T) {
	if _, err := os.Stat("/dev/full"); err != nil {
		t.Skip("/dev/full unavailable")
	}
	previous := DisableFooter
	t.Cleanup(func() { DisableFooter = previous })
	options := sysoTestOptions(t)
	options.Out = "/dev/full"
	if err := GenerateSyso(options); !errors.Is(err, syscall.ENOSPC) {
		t.Fatalf("error=%v, want disk-full write failure", err)
	}
}
