//go:build unix

package voidpackage

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"

	"golang.org/x/sys/unix"
)

func TestNamedPipeInputIsRejectedWithoutOpening(t *testing.T) {
	options := fixture(t)
	pipe := filepath.Join(t.TempDir(), "input.pipe")
	if err := unix.Mkfifo(pipe, 0600); err != nil {
		t.Fatal(err)
	}
	options.Contents[0].Source = pipe
	var output bytes.Buffer
	if err := Write(context.Background(), &output, options); err == nil || output.Len() != 0 {
		t.Fatalf("named pipe accepted: %v", err)
	}
}
