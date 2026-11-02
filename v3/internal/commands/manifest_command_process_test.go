package commands

import (
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestManifestCommandOwnsChildrenWithoutSystemTools(t *testing.T) {
	for _, mode := range []string{"parent", "orphan-parent"} {
		t.Run(mode, func(t *testing.T) {
			t.Setenv("PATH", "")
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			address := listener.Addr().String()
			if err := listener.Close(); err != nil {
				t.Fatal(err)
			}
			root := t.TempDir()
			ready := filepath.Join(root, "ready")
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			result := make(chan error, 1)
			go func() {
				_, err := runManifestCommand(ctx, root, []string{
					"WAILS_PROCESS_TREE_HELPER=" + mode,
					"WAILS_PROCESS_TREE_ADDRESS=" + address,
					"WAILS_PROCESS_TREE_READY=" + ready,
				}, os.Args[0], "-test.run=^TestManifestProcessTreeHelper$")
				result <- err
			}()
			waitForManifestDevTestFile(t, ready)
			if mode == "parent" {
				cancel()
			}
			select {
			case err := <-result:
				if mode == "parent" && !errors.Is(err, context.Canceled) {
					t.Fatalf("cancellation lost: %v", err)
				}
				if mode == "orphan-parent" && err != nil {
					t.Fatalf("wrapper failed: %v", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("command did not finish")
			}
			deadline := time.Now().Add(time.Second)
			for {
				connection, err := net.DialTimeout("tcp", address, 20*time.Millisecond)
				if err != nil {
					break
				}
				_ = connection.Close()
				if time.Now().After(deadline) {
					t.Fatal("command left its child running")
				}
				time.Sleep(10 * time.Millisecond)
			}
		})
	}
}

func TestManifestCommandCancelledBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	output, err := runManifestCommand(ctx, t.TempDir(), nil, "missing-executable")
	if !errors.Is(err, context.Canceled) || output != "" {
		t.Fatalf("cancelled command attempted to start: output=%q error=%v", output, err)
	}
}
