package application

import (
	"os"
	"strings"
	"testing"
)

func TestAutostartSlug(t *testing.T) {
	cases := map[string]string{
		"Simple":             "simple",
		"My App":             "my-app",
		"My  Cool   App":     "my--cool---app",
		"Foo/Bar":            "foobar",
		"  trim me  ":        "trim-me",
		"All Punc!@#$":       "all-punc",
		"":                   "wails-app",
		"héllo":              "hllo",
		"v1.2.3":             "v1.2.3",
		"___leading_under":   "leading_under",
		"...":                "wails-app",
	}
	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			if got := autostartSlug(in); got != want {
				t.Errorf("autostartSlug(%q) = %q, want %q", in, got, want)
			}
		})
	}
}

func TestValidateAutostartIdentifier(t *testing.T) {
	good := []string{
		"",
		"com.example.app",
		"my-app",
		"my_app",
		"App123",
		strings.Repeat("a", 200),
	}
	for _, s := range good {
		if err := validateAutostartIdentifier(s); err != nil {
			t.Errorf("validateAutostartIdentifier(%q) unexpectedly failed: %v", s, err)
		}
	}
	bad := []string{
		"with space",
		"slash/in/it",
		"back\\slash",
		"line\nbreak",
		"quote\"in",
		strings.Repeat("a", 201),
	}
	for _, s := range bad {
		if err := validateAutostartIdentifier(s); err == nil {
			t.Errorf("validateAutostartIdentifier(%q) should have failed", s)
		}
	}
}

func TestWriteFileAtomic(t *testing.T) {
	dir := t.TempDir()
	target := dir + "/out.txt"
	if err := writeFileAtomic(target, []byte("hello"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	// Subsequent writes overwrite, not append, and leave no leftover tempfiles.
	if err := writeFileAtomic(target, []byte("world"), 0644); err != nil {
		t.Fatalf("rewrite: %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "out.txt" {
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("expected exactly out.txt in dir, got %v", names)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "world" {
		t.Errorf("contents=%q want %q", string(data), "world")
	}
}
