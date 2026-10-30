//go:build linux

package fileexplorer

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExplorerQueryCancelledContext(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "xdg-mime"), []byte("#!/bin/sh\nexit 1\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", root)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := explorerBinArgs(ctx, root, false); !errors.Is(err, context.Canceled) {
		t.Errorf("query error = %v, want context.Canceled", err)
	}
}

func TestExplorerQueryDeadline(t *testing.T) {
	if _, err := os.Stat("/bin/sleep"); err != nil {
		t.Skip("/bin/sleep is needed for the blocking query fixture")
	}
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "xdg-mime"), []byte("#!/bin/sh\nexec /bin/sleep 30\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", root)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	result := make(chan error, 1)
	go func() {
		_, _, err := explorerBinArgs(ctx, root, false)
		result <- err
	}()
	select {
	case err := <-result:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("query error = %v, want context.DeadlineExceeded", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("query ignored its context deadline")
	}
}
