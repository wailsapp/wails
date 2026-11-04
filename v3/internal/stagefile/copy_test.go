package stagefile

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCopyPreservesBytesAndModeWithoutOverwriting(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	target := filepath.Join(root, "stage", "target")
	data := bytes.Repeat([]byte("file content"), 10000)
	if err := os.WriteFile(source, data, 0755); err != nil {
		t.Fatal(err)
	}
	if err := Copy(context.Background(), source, target, 0, make([]byte, 32768)); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(target)
	if err != nil || !bytes.Equal(actual, data) {
		t.Fatal("copy changed data")
	}
	sourceInfo, err := os.Stat(source)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(target)
	if err != nil || info.Mode().Perm() != sourceInfo.Mode().Perm() {
		t.Fatal("copy lost executable permissions")
	}
	if err := Copy(context.Background(), source, target, 0644, nil); !os.IsExist(err) {
		t.Fatalf("overwrote existing destination: %v", err)
	}
	if err := Copy(context.Background(), source, source, 0, nil); !os.IsExist(err) {
		t.Fatalf("accepted source alias: %v", err)
	}
	actual, err = os.ReadFile(source)
	if err != nil || !bytes.Equal(actual, data) {
		t.Fatal("copy damaged source")
	}
}

func TestCopyCancellationAndRegularFileRequirement(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	target := filepath.Join(root, "stage", "target")
	if err := os.WriteFile(source, []byte("content"), 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := Copy(ctx, source, target, 0, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation=%v", err)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("cancelled copy created destination: %v", err)
	}
	if err := Copy(context.Background(), root, target, 0, nil); err == nil {
		t.Fatal("accepted directory source")
	}
	if err := Copy(context.Background(), source, target, 0644, nil); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(target); err != nil || runtime.GOOS != "windows" && info.Mode().Perm() != 0644 {
		t.Fatal("mode override ignored")
	}
}
