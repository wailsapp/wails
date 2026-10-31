package artifactsum

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestWriteSidecar(t *testing.T) {
	payload := bytes.Repeat([]byte("payload\x00\xff"), 20000)
	type sidecarCase struct {
		name, rendered, prefix string
		data                   []byte
	}
	cases := []sidecarCase{
		{"My App ü.dmg", "My App ü.dmg", "", payload},
		{"empty.deb", "empty.deb", "", nil},
	}
	if runtime.GOOS != "windows" {
		cases = append(cases,
			sidecarCase{"line\nbreak\\.AppImage", "line\\nbreak\\\\.AppImage", "\\", payload},
			sidecarCase{"carriage\rreturn.deb", "carriage\\rreturn.deb", "\\", payload},
		)
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, tc.name)
			data := tc.data
			for run := 0; run < 2; run++ {
				if err := os.WriteFile(path, data, 0600); err != nil {
					t.Fatal(err)
				}
				if err := Write(context.Background(), path); err != nil {
					t.Fatal(err)
				}
				want := fmt.Sprintf("%s%x  %s\n", tc.prefix, sha256.Sum256(data), tc.rendered)
				got, err := os.ReadFile(path + ".sha256")
				if err != nil || string(got) != want {
					t.Fatalf("sidecar=%q want=%q err=%v", got, want, err)
				}
				data = []byte("replaced artifact")
			}
			entries, err := os.ReadDir(root)
			if err != nil || len(entries) != 2 {
				t.Fatalf("temporary files leaked: %v, %v", entries, err)
			}
		})
	}
}

func TestWriteFailurePreservesSidecar(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "app.AppImage")
	if err := os.WriteFile(path+".sha256", []byte("previous"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Write(context.Background(), path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing artifact: %v", err)
	}
	if err := os.Mkdir(path, 0700); err != nil {
		t.Fatal(err)
	}
	if err := Write(context.Background(), path); err == nil {
		t.Fatal("accepted a directory")
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("changed"), 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := Write(ctx, path); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation: %v", err)
	}
	got, err := os.ReadFile(path + ".sha256")
	if err != nil || string(got) != "previous" {
		t.Fatalf("previous sidecar changed: %q %v", got, err)
	}
}

func BenchmarkWrite(b *testing.B) {
	root := b.TempDir()
	path := filepath.Join(root, "app.AppImage")
	data := bytes.Repeat([]byte(strings.Repeat("payload", 1024)), 2048)
	if err := os.WriteFile(path, data, 0600); err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Write(context.Background(), path); err != nil {
			b.Fatal(err)
		}
	}
}
