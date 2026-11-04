//go:build !windows

package commands

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestManifestDevInputCancellationDoesNotWaitForNewline(t *testing.T) {
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer read.Close()
	defer write.Close()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- readManifestDevInput(ctx, read, func(string) bool { t.Error("unexpected command"); return true })
	}()
	if _, err := write.WriteString("incomplete"); err != nil {
		t.Fatal(err)
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("stdin reader blocked after cancellation")
	}
	if _, err := write.WriteString("still-owned"); err != nil {
		t.Fatal("reader closed caller stdin", err)
	}
}

func TestManifestDevInputLinesAndEOF(t *testing.T) {
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer read.Close()
	done := make(chan error, 1)
	var lines []string
	go func() {
		done <- readManifestDevInput(context.Background(), read, func(line string) bool { lines = append(lines, line); return true })
	}()
	if _, err := write.WriteString("r\nbindings\r\nf\nc\nquit\n"); err != nil {
		t.Fatal(err)
	}
	write.Close()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("EOF did not stop reader")
	}
	if !reflect.DeepEqual(lines, []string{"r", "bindings", "f", "c", "quit"}) {
		t.Fatal(lines)
	}
}
