package application

import (
	"os"
	"path/filepath"
	"testing"
)

func TestThumbnailOptionsWithDefaults(t *testing.T) {
	tests := []struct {
		name string
		in   ThumbnailOptions
		want ThumbnailOptions
	}{
		{name: "zero", in: ThumbnailOptions{}, want: ThumbnailOptions{Width: 256, Height: 256, Scale: 1}},
		{name: "width only", in: ThumbnailOptions{Width: 100}, want: ThumbnailOptions{Width: 100, Height: 100, Scale: 1}},
		{name: "height only", in: ThumbnailOptions{Height: 80}, want: ThumbnailOptions{Width: 80, Height: 80, Scale: 1}},
		{name: "negative scale", in: ThumbnailOptions{Width: 10, Height: 20, Scale: -2}, want: ThumbnailOptions{Width: 10, Height: 20, Scale: 1}},
		{name: "explicit", in: ThumbnailOptions{Width: 64, Height: 32, Scale: 2, IconMode: true}, want: ThumbnailOptions{Width: 64, Height: 32, Scale: 2, IconMode: true}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.in.withDefaults(); got != test.want {
				t.Fatalf("withDefaults() = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestValidateQuickLookPath(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(existing, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{name: "existing file", path: existing},
		{name: "existing directory", path: dir},
		{name: "empty", path: "", wantErr: true},
		{name: "relative", path: "file.txt", wantErr: true},
		{name: "missing", path: filepath.Join(dir, "missing.txt"), wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateQuickLookPath(test.path)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateQuickLookPath(%q) error = %v, wantErr %v", test.path, err, test.wantErr)
			}
		})
	}
	if err := validateQuickLookPaths(nil); err == nil {
		t.Fatal("expected an error for no paths")
	}
	if err := validateQuickLookPaths([]string{existing, "relative"}); err == nil {
		t.Fatal("expected an error when any path is invalid")
	}
	if err := validateQuickLookPaths([]string{existing, dir}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
