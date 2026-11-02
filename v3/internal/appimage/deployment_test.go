package appimage

import (
	"context"
	"debug/elf"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRuntimeSymlinksRemainWithinSuppliedTree(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink fixture needs Unix symlinks")
	}
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "usr/lib"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "usr/lib/library"), []byte("target"), 0644); err != nil {
		t.Fatal(err)
	}
	for name, target := range map[string]string{"lib": "/usr/lib", "relative": "usr/lib", "escape": "../../etc/passwd", "cycle": "cycle"} {
		if err := os.Symlink(target, filepath.Join(root, name)); err != nil {
			t.Fatal(err)
		}
	}
	d, err := newRuntimeDeployment(root, "amd64", 4, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"/lib/library", "/relative/library"} {
		got, err := d.resolve(name)
		if err != nil || got != filepath.Join(root, "usr/lib/library") {
			t.Fatalf("resolve %s: %s, %v", name, got, err)
		}
	}
	for _, name := range []string{"/escape", "/cycle"} {
		if _, err := d.resolve(name); err == nil {
			t.Fatalf("accepted %s", name)
		}
	}
	if _, err := d.library("libgtk-4.so.1"); err == nil {
		t.Fatal("resolved absent target library from host")
	}
}

func TestRuntimeDependencyClosureAndConflicts(t *testing.T) {
	libs := gtkLibraries(t)
	d, err := newRuntimeDeployment(libs, runtime.GOARCH, 4, []string{"/"})
	if err != nil {
		t.Fatal(err)
	}
	if err := d.dependencies(t.Context(), []string{filepath.Join(libs, "libwrapper.so")}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"libexample.so", "libgtk-4.so.1"} {
		if d.files["usr/lib/"+name] == "" {
			t.Fatalf("missing transitive library %s", name)
		}
	}
	different := filepath.Join(libs, "different")
	if err := os.WriteFile(different, []byte("different library"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := d.add("usr/lib/libexample.so", different); err == nil {
		t.Fatal("accepted conflicting same-name libraries")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := d.dependencies(ctx, []string{filepath.Join(libs, "libwrapper.so")}); err == nil {
		t.Fatal("ignored cancellation")
	}
}

func TestBundledLibraryPathOverridesLegacyRpath(t *testing.T) {
	root := gtkLibraries(t)
	compiler, err := exec.LookPath("cc")
	if err != nil {
		t.Skip("needs ELF fixture compiler")
	}
	file := filepath.Join(root, "rpath.so")
	if out, err := exec.Command(compiler, "-shared", "-fPIC", "-nostdlib", "-Wl,--disable-new-dtags,-rpath,/old/build/machine", filepath.Join(root, "libgtk-4.so.1.c"), "-o", file).CombinedOutput(); err != nil {
		t.Fatalf("compile: %s %v", out, err)
	}
	before, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if err := preferBundledLibraries(file); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(before) != len(after) {
		t.Fatal("loader patch resized ELF")
	}
	image, err := elf.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer image.Close()
	old, _ := image.DynString(elf.DT_RPATH)
	current, _ := image.DynString(elf.DT_RUNPATH)
	if len(old) != 0 || len(current) != 1 || current[0] != "/old/build/machine" {
		t.Fatalf("RPATH=%q RUNPATH=%q", old, current)
	}
	if err := preferBundledLibraries(file); err != nil {
		t.Fatal(err)
	}
}

func TestRuntimeLoaderCachePathsBecomeRelative(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "cache")
	if err := os.WriteFile(source, []byte("# loader cache\n\"/usr/lib64/gdk-pixbuf/loaders/libexample.so\"\n\"image/png\" \"\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	d, err := newRuntimeDeployment(root, "amd64", 4, nil)
	if err != nil {
		t.Fatal(err)
	}
	d.files["usr/lib/gdk-pixbuf/loaders.cache"] = source
	stage := t.TempDir()
	if err := d.stage(t.Context(), stage); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(stage, "usr/lib/gdk-pixbuf/loaders.cache"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "\"libexample.so\"") || !strings.Contains(string(data), "\"image/png\"") {
		t.Fatalf("bad relocated cache: %q", data)
	}
}

func BenchmarkRuntimeDependencyCollection(b *testing.B) {
	root := gtkLibraries(b)
	inputs := []string{filepath.Join(root, "libwrapper.so")}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		d, err := newRuntimeDeployment(root, runtime.GOARCH, 4, []string{"/"})
		if err != nil {
			b.Fatal(err)
		}
		if err := d.dependencies(b.Context(), inputs); err != nil {
			b.Fatal(err)
		}
	}
}

func TestRuntimeRequiresCacheForExternalLoaderModules(t *testing.T) {
	d := runtimeDeployment{files: map[string]string{"usr/lib/gdk-pixbuf-2.0/2.10.0/loaders/libtest.so": "module"}}
	if err := d.checkLoaderCaches(); err == nil {
		t.Fatal("accepted pixbuf modules without cache")
	}
	d.files["usr/lib/gdk-pixbuf-2.0/2.10.0/loaders.cache"] = "cache"
	if err := d.checkLoaderCaches(); err != nil {
		t.Fatal(err)
	}
	d.files["usr/lib/gtk-3.0/3.0.0/immodules/im-test.so"] = "module"
	if err := d.checkLoaderCaches(); err == nil {
		t.Fatal("accepted GTK input modules without cache")
	}
	d.files["usr/lib/gtk-3.0/3.0.0/immodules.cache"] = "cache"
	if err := d.checkLoaderCaches(); err != nil {
		t.Fatal(err)
	}
}
