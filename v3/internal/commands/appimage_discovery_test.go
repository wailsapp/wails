package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestFindGTKFilesInRoots(t *testing.T) {
	root := t.TempDir()
	want := []string{filepath.Join(root, "lib", "WebKitWebProcess"), filepath.Join(root, "lib64", "WebKitWebProcess"), filepath.Join(root, "lib64", "libwebkitgtkinjectedbundle.so")}
	for _, path := range want {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(root, "WebKitWebProcess"), 0700); err != nil {
		t.Fatal(err)
	}
	got, err := findGTKFilesInRoots([]string{root}, []string{"WebKitWebProcess", "libwebkitgtkinjectedbundle.so"})
	if err != nil || !slices.Equal(got, want) {
		t.Fatalf("files=%v, error=%v; want %v", got, err, want)
	}
	if _, err := findGTKFilesInRoots([]string{root}, []string{"missing-process"}); err == nil || !strings.Contains(err.Error(), "missing-process") {
		t.Errorf("missing file error=%v", err)
	}
	if _, err := findGTKFilesInRoots([]string{filepath.Join(root, "absent")}, []string{"WebKitWebProcess"}); err == nil || !strings.Contains(err.Error(), "WebKitWebProcess") {
		t.Errorf("missing root error=%v", err)
	}
}

func BenchmarkFindGTKFilesInRoots(b *testing.B) {
	root := b.TempDir()
	for directory := range 32 {
		dir := filepath.Join(root, fmt.Sprint(directory))
		if err := os.Mkdir(dir, 0700); err != nil {
			b.Fatal(err)
		}
		for file := range 64 {
			name := fmt.Sprint(file)
			if file == 0 {
				name = "WebKitWebProcess"
			}
			if err := os.WriteFile(filepath.Join(dir, name), nil, 0600); err != nil {
				b.Fatal(err)
			}
		}
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		files, err := findGTKFilesInRoots([]string{root}, []string{"WebKitWebProcess"})
		if err != nil || len(files) != 32 {
			b.Fatalf("files=%d, error=%v", len(files), err)
		}
	}
}
