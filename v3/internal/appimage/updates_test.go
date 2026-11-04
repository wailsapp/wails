package appimage

import (
	"context"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestUpdateControlPublication(t *testing.T) {
	for _, failure := range []string{"", "missing", "wrong image", "wrong name", "cancel", "directory"} {
		t.Run(failure, func(t *testing.T) {
			options := appImageFixture(t)
			options.UpdateInfo = "zsync|https://example.test/My.AppImage.zsync"
			output := filepath.Join(options.Root, options.Output)
			if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(output, []byte("old image"), 0755); err != nil {
				t.Fatal(err)
			}
			if failure == "directory" {
				if err := os.Mkdir(output+".zsync", 0755); err != nil {
					t.Fatal(err)
				}
			} else if err := os.WriteFile(output+".zsync", []byte("old control"), 0644); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			options.Runner = func(_ context.Context, command Command) error {
				if command.Env["LDAI_UPDATE_INFORMATION"] != options.UpdateInfo || command.Env["OUTPUT"] != filepath.Base(output) {
					t.Fatalf("update settings lost: %v", command.Env)
				}
				image := minimalELF(runtime.GOARCH, true)
				if err := os.WriteFile(filepath.Join(command.Directory, command.Env["OUTPUT"]), image, 0755); err != nil {
					return err
				}
				if failure == "missing" {
					return nil
				}
				filename := command.Env["OUTPUT"]
				if failure == "wrong name" {
					filename = "result.AppImage"
				}
				digest := fmt.Sprintf("%x", sha1.Sum(image))
				if failure == "wrong image" {
					digest = strings.Repeat("0", 40)
				}
				header := fmt.Sprintf("zsync: 0.6.2\nFilename: %s\nLength: %d\nSHA-1: %s\n\n", filename, len(image), digest)
				if failure == "cancel" {
					cancel()
				}
				return os.WriteFile(filepath.Join(command.Directory, command.Env["OUTPUT"]+".zsync"), []byte(header), 0644)
			}
			err := Create(ctx, options)
			if failure == "" {
				if err != nil {
					t.Fatal(err)
				}
				if err := verifyZsync(context.Background(), output); err != nil {
					t.Fatal(err)
				}
				if options.Env["LDAI_UPDATE_INFORMATION"] != "" {
					t.Fatal("mutated caller environment")
				}
			} else {
				if err == nil {
					t.Fatal("expected update failure")
				}
				old, err := os.ReadFile(output)
				if err != nil || string(old) != "old image" {
					t.Fatalf("previous image changed: %q, %v", old, err)
				}
				if failure != "directory" {
					old, err = os.ReadFile(output + ".zsync")
					if err != nil || string(old) != "old control" {
						t.Fatalf("previous control changed: %q, %v", old, err)
					}
				}
			}
		})
	}
}
