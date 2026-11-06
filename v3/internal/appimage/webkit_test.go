package appimage

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"testing"
)

func webKitFixture(t *testing.T) (Options, string, string) {
	t.Helper()
	libs := gtkLibraries(t)
	options := appImageFixture(t)
	options.WebKit = true
	options.GTKPlugin = []byte("existing plugin artifact")
	options.LibraryDirs = []string{libs}
	options.Binaries = []Binary{{"app", filepath.Join(libs, "libwrapper.so")}}
	options.RuntimeRoot = t.TempDir()
	selected := filepath.Join(options.RuntimeRoot, "libexec", "webkitgtk-6.0")
	other := filepath.Join(options.RuntimeRoot, "libexec", "webkit2gtk-4.1")
	compiler, err := exec.LookPath("cc")
	if err != nil {
		t.Skip("requires cc for GTK dependency fixtures")
	}
	old := filepath.Join(libs, "libgtk-3.so.0")
	if output, err := exec.Command(compiler, "-shared", "-fPIC", "-nostdlib", "-Wl,-soname,libgtk-3.so.0", filepath.Join(libs, "libgtk-4.so.1.c"), "-o", old).CombinedOutput(); err != nil {
		t.Fatalf("GTK3 fixture: %s, %v", output, err)
	}
	wrong := filepath.Join(libs, "wrong.so")
	if output, err := exec.Command(compiler, "-shared", "-fPIC", "-nostdlib", filepath.Join(libs, "libexample.so.c"), "-L"+libs, "-l:libgtk-3.so.0", "-o", wrong).CombinedOutput(); err != nil {
		t.Fatalf("GTK3 helper: %s, %v", output, err)
	}
	for _, dir := range []string{selected, other} {
		data := minimalELF(runtime.GOARCH, false)
		if dir == other {
			data, err = os.ReadFile(wrong)
			if err != nil {
				t.Fatal(err)
			}
		}
		for _, name := range []string{"WebKitWebProcess", "WebKitNetworkProcess", "WebKitGPUProcess", "injected-bundle/libwebkitgtkinjectedbundle.so"} {
			filename := filepath.Join(dir, name)
			if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filename, data, 0755); err != nil {
				t.Fatal(err)
			}
		}
	}
	return options, selected, other
}

func TestWebKitDeploymentRelocatesBeforeCompression(t *testing.T) {
	options, selected, other := webKitFixture(t)
	calls := 0
	original := []byte("prefix\x00" + selected + "\x00" + selected + "/injected-bundle/\x00suffix")
	options.Runner = func(ctx context.Context, command Command) error {
		calls++
		appdir := filepath.Join(command.Directory, "AppDir")
		library := filepath.Join(appdir, "usr", "lib", "libwebkitgtk-6.0.so.4")
		if command.Env["DEPLOY_GTK_VERSION"] != "4" {
			t.Fatalf("GTK environment: %v", command.Env)
		}
		if calls == 1 {
			if slices.Contains(command.Args, "--output") || !slices.Contains(command.Args, "--plugin") {
				t.Fatalf("deployment args: %q", command.Args)
			}
			for _, name := range []string{"WebKitWebProcess", "WebKitNetworkProcess", "WebKitGPUProcess"} {
				if _, err := os.Stat(filepath.Join(appdir, selected, name)); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := os.Stat(filepath.Join(appdir, other)); !os.IsNotExist(err) {
				t.Fatalf("unselected GTK stack staged: %v", err)
			}
			return os.WriteFile(library, original, 0755)
		}
		if calls != 2 || slices.Contains(command.Args, "--plugin") || !slices.Contains(command.Args, "--output") {
			t.Fatalf("compression args: %q", command.Args)
		}
		patched, err := os.ReadFile(library)
		if err != nil {
			t.Fatal(err)
		}
		replacement := []byte("././webkitgtk-6.0/")
		replacement = append(replacement, make([]byte, len(selected)+1-len(replacement))...)
		want := bytes.ReplaceAll(original, []byte(selected+"\x00"), replacement)
		if !bytes.Equal(patched, want) {
			t.Fatalf("library patch changed width or injected-bundle prefix: %q", patched)
		}
		for _, name := range []string{"WebKitWebProcess", "WebKitNetworkProcess", "WebKitGPUProcess"} {
			if _, err := os.Stat(filepath.Join(appdir, "usr", "lib", "webkitgtk-6.0", name)); err != nil {
				t.Fatal(err)
			}
		}
		if info, err := os.Stat(library); err != nil || info.Mode().Perm() != 0755 {
			t.Fatalf("library mode: %v, %v", info, err)
		}
		return os.WriteFile(filepath.Join(command.Directory, command.Env["OUTPUT"]), minimalELF(runtime.GOARCH, true), 0755)
	}
	if err := Create(t.Context(), options); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("got %d linuxdeploy invocations", calls)
	}
}

func TestWebKitDeploymentFailureStopsCompression(t *testing.T) {
	for _, failure := range []string{"deploy", "relocate"} {
		t.Run(failure, func(t *testing.T) {
			options, _, _ := webKitFixture(t)
			calls := 0
			deploymentError := errors.New("deployment failed")
			options.Runner = func(context.Context, Command) error {
				calls++
				if failure == "deploy" {
					return deploymentError
				}
				return nil // No deployed WebKit library: relocation must fail.
			}
			err := Create(t.Context(), options)
			if err == nil || calls != 1 {
				t.Fatalf("err=%v calls=%d", err, calls)
			}
			if failure == "deploy" && !errors.Is(err, deploymentError) {
				t.Fatalf("lost deployment error: %v", err)
			}
			if _, err := os.Stat(filepath.Join(options.Root, options.Output)); !os.IsNotExist(err) {
				t.Fatalf("published failed image: %v", err)
			}
		})
	}
}
