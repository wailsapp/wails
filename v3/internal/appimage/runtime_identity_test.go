package appimage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRuntimeIdentityTracksSelectedFilesAndSameMetadataChanges(t *testing.T) {
	root := t.TempDir()
	names := []string{"usr/lib64/libglib-2.0.so.0", "usr/lib64/libgobject-2.0.so.0", "usr/lib64/libgio-2.0.so.0", "usr/lib64/libsoup-3.0.so.0", "usr/lib64/libgtk-4.so.1", "usr/lib64/libwebkitgtk-6.0.so.4", "usr/lib64/libjavascriptcoregtk-6.0.so.1", "usr/libexec/webkitgtk-6.0/WebKitWebProcess", "usr/libexec/webkitgtk-6.0/WebKitNetworkProcess", "usr/lib64/webkitgtk-6.0/injected-bundle/libwebkitgtkinjectedbundle.so", "app"}
	for _, name := range names {
		file := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, minimalELF("amd64", false), 0755); err != nil {
			t.Fatal(err)
		}
	}
	schema := filepath.Join(root, "usr/share/glib-2.0/schemas/gschemas.compiled")
	if err := os.MkdirAll(filepath.Dir(schema), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(schema, []byte("compiled schema fixture"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "runtime"), minimalELF("amd64", true), 0755); err != nil {
		t.Fatal(err)
	}
	options := Options{Root: root, RuntimeRoot: root, RuntimePath: "runtime", Arch: "amd64", GTKVersion: 4, WebKit: true, Binaries: []Binary{{Name: "app", Path: "app"}}}
	identity := func() string {
		t.Helper()
		value, err := options.RuntimeIdentity(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		return value
	}
	first := identity()
	if identity() != first {
		t.Fatal("unstable identity")
	}
	info, err := os.Stat(schema)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(schema, []byte("modified schema fixture"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(schema, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	second := identity()
	if second == first {
		t.Fatal("same-size same-mtime content change was missed")
	}
	module := filepath.Join(root, "usr/lib64/girepository-1.0/extra.typelib")
	if err := os.MkdirAll(filepath.Dir(module), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(module, []byte("added selected data"), 0644); err != nil {
		t.Fatal(err)
	}
	third := identity()
	if third == second {
		t.Fatal("new runtime data was missed")
	}
	if err := os.WriteFile(filepath.Join(root, "unrelated"), []byte("unrelated file"), 0644); err != nil {
		t.Fatal(err)
	}
	if identity() != third {
		t.Fatal("unrelated file invalidates runtime identity")
	}
}
