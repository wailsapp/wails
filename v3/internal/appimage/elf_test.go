package appimage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

func gtkLibraries(t testing.TB) string {
	t.Helper()
	if runtime.GOOS != "linux" {
		t.Skip("requires native Linux linker")
	}
	compiler, err := exec.LookPath("cc")
	if err != nil {
		t.Skip("requires C compiler for ELF dependency fixtures")
	}
	root := t.TempDir()
	for _, library := range []struct{ name, source, dependency string }{
		{"libgtk-4.so.1", "int gtk_fixture(){return 4;}", ""},
		{"libexample.so", "extern int gtk_fixture(); int example(){return gtk_fixture();}", "libgtk-4.so.1"},
		{"libwrapper.so", "extern int example(); int wrapper(){return example();}", "libexample.so"},
	} {
		source := filepath.Join(root, library.name+".c")
		if err := os.WriteFile(source, []byte(library.source), 0600); err != nil {
			t.Fatal(err)
		}
		args := []string{"-shared", "-fPIC", "-nostdlib", "-Wl,-rpath,$ORIGIN", "-Wl,-soname," + library.name, source, "-L" + root, "-o", filepath.Join(root, library.name)}
		if library.dependency != "" {
			args = append(args, "-l:"+library.dependency)
		}
		if output, err := exec.Command(compiler, args...).CombinedOutput(); err != nil {
			t.Fatalf("compile ELF fixture: %s, %v", output, err)
		}
	}
	return root
}

func TestELFAnalysisFindsTransitiveGTKViaOrigin(t *testing.T) {
	root := gtkLibraries(t)
	input := filepath.Join(root, "libwrapper.so")
	analysis, err := analyzeELF(context.Background(), []string{input}, nil)
	if err != nil || analysis.gtk != "4" {
		t.Fatalf("analysis=%+v, err=%v", analysis, err)
	}
	if err := os.Remove(filepath.Join(root, "libgtk-4.so.1")); err != nil {
		t.Fatal(err)
	}
	if _, err := analyzeELF(context.Background(), []string{input}, nil); err == nil || !strings.Contains(err.Error(), "libgtk-4.so.1") {
		t.Fatalf("missing dependency: %v", err)
	}
}

func TestGTKPluginIsExplicitToolInput(t *testing.T) {
	root := gtkLibraries(t)
	options := appImageFixture(t)
	options.Binaries = []Binary{{"app", filepath.Join(root, "libwrapper.so")}}
	options.GTKPlugin = []byte("existing plugin artifact")
	options.Runner = func(ctx context.Context, command Command) error {
		if command.Env["DEPLOY_GTK_VERSION"] != "4" || !strings.HasPrefix(command.Env["PATH"], command.Directory+string(os.PathListSeparator)) {
			t.Fatalf("GTK env=%v", command.Env)
		}
		plugin := filepath.Join(command.Directory, "linuxdeploy-plugin-gtk.sh")
		data, err := os.ReadFile(plugin)
		if err != nil || string(data) != "existing plugin artifact" {
			t.Fatalf("plugin not staged: %q, %v", data, err)
		}
		return os.WriteFile(filepath.Join(command.Directory, "result.AppImage"), minimalELF(runtime.GOARCH, true), 0755)
	}
	if err := Create(context.Background(), options); err != nil {
		t.Fatal(err)
	}
}

func TestRuntimeDiscoveryAndCancellation(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"WebKitWebProcess", "WebKitNetworkProcess", "libwebkitgtkinjectedbundle.so"} {
		path := filepath.Join(root, "libexec", "webkitgtk-6.0", name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	files, err := FindRuntimeFiles(context.Background(), root, []string{"WebKitWebProcess", "WebKitNetworkProcess", "libwebkitgtkinjectedbundle.so"})
	if err != nil || len(files) != 3 {
		t.Fatalf("files=%v, err=%v", files, err)
	}
	if _, err := FindRuntimeFiles(context.Background(), root, []string{"missing"}); err == nil {
		t.Fatal("accepted missing helper")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := FindRuntimeFiles(ctx, root, []string{"WebKitWebProcess"}); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation: %v", err)
	}
}

func TestRuntimeHelpersMatchGTKStack(t *testing.T) {
	root := gtkLibraries(t)
	matching := filepath.Join(root, "libwrapper.so")
	compiler, err := exec.LookPath("cc")
	if err != nil {
		t.Skip("requires cc")
	}
	old := filepath.Join(root, "libgtk-3.so.0")
	if output, err := exec.Command(compiler, "-shared", "-fPIC", "-nostdlib", "-Wl,-soname,libgtk-3.so.0", filepath.Join(root, "libgtk-4.so.1.c"), "-o", old).CombinedOutput(); err != nil {
		t.Fatalf("old stack: %s, %v", output, err)
	}
	wrong := filepath.Join(root, "wrong-stack.so")
	if output, err := exec.Command(compiler, "-shared", "-fPIC", "-nostdlib", "-Wl,-rpath,$ORIGIN", filepath.Join(root, "libexample.so.c"), "-L"+root, "-l:libgtk-3.so.0", "-o", wrong).CombinedOutput(); err != nil {
		t.Fatalf("old helper: %s, %v", output, err)
	}
	selected, err := matchingRuntimeFiles(context.Background(), []string{matching, wrong}, "4", nil)
	if err == nil {
		t.Fatal("different basenames should report a missing GTK4 helper")
	}
	shared := filepath.Join(root, "gtk3", "libwrapper.so")
	if err := os.MkdirAll(filepath.Dir(shared), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(wrong, shared); err != nil {
		t.Fatal(err)
	}
	selected, err = matchingRuntimeFiles(context.Background(), []string{matching, shared}, "4", []string{root})
	if err != nil || len(selected) != 1 || selected[0] != matching {
		t.Fatalf("selected=%v, err=%v", selected, err)
	}
}

func BenchmarkAnalyzeELF(b *testing.B) {
	root := gtkLibraries(b)
	filename := filepath.Join(root, "libwrapper.so")
	b.Run("GoELF", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			analysis, err := analyzeELF(context.Background(), []string{filename}, nil)
			if err != nil || analysis.gtk != "4" {
				b.Fatalf("%+v, %v", analysis, err)
			}
		}
	})
	b.Run("LegacyLdd", func(b *testing.B) {
		program, err := exec.LookPath("ldd")
		if err != nil {
			b.Skip("ldd unavailable for baseline")
		}
		b.ReportAllocs()
		for b.Loop() {
			output, err := exec.Command(program, filename).CombinedOutput()
			if err != nil || !strings.Contains(string(output), "libgtk-4.so.1") {
				b.Fatalf("%s, %v", output, err)
			}
		}
	})
}

func TestRuntimeDiscoverySkipsUnrelatedTreesAndFallsBack(t *testing.T) {
	root := t.TempDir()
	names := []string{"WebKitWebProcess", "WebKitNetworkProcess"}
	write := func(relative string) string {
		t.Helper()
		path := filepath.Join(root, relative)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, nil, 0600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	first := write("libexec/webkitgtk-6.0/" + names[0])
	second := write("lib64/webkitgtk-6.0/" + names[1])
	unrelated := write("share/old-app.AppDir/" + names[0])
	files, err := FindRuntimeFiles(context.Background(), root, names)
	if err != nil || len(files) != 2 {
		t.Fatalf("runtime files=%v err=%v", files, err)
	}
	for _, file := range files {
		if file == unrelated {
			t.Fatal("searched unrelated tree despite complete runtime set")
		}
	}
	if err := os.Remove(second); err != nil {
		t.Fatal(err)
	}
	fallback := write("vendor/custom-webkit/" + names[1])
	files, err = FindRuntimeFiles(context.Background(), root, names)
	if err != nil || !slices.Contains(files, fallback) || !slices.Contains(files, first) {
		t.Fatalf("fallback files=%v err=%v", files, err)
	}
}

func BenchmarkRuntimeDiscovery(b *testing.B) {
	root := b.TempDir()
	for i := range 2000 {
		path := filepath.Join(root, "share", fmt.Sprintf("package-%04d", i), "asset")
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			b.Fatal(err)
		}
		if err := os.WriteFile(path, nil, 0600); err != nil {
			b.Fatal(err)
		}
	}
	names := []string{"WebKitWebProcess", "WebKitNetworkProcess", "libwebkitgtkinjectedbundle.so"}
	for _, name := range names {
		path := filepath.Join(root, "libexec", "webkitgtk-6.0", name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			b.Fatal(err)
		}
		if err := os.WriteFile(path, nil, 0600); err != nil {
			b.Fatal(err)
		}
	}
	b.Run("RuntimeDirectories", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := FindRuntimeFiles(context.Background(), root, names); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("WholeRoot", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, _, err := findRuntimeFiles(context.Background(), []string{root}, names); err != nil {
				b.Fatal(err)
			}
		}
	})
}
