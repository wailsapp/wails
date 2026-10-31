package appimage

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestDefaultLauncherRunsFromBundledLibraries(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("AppImage launcher requires a POSIX shell")
	}
	options := appImageFixture(t)
	options.Tools.AppRun = ""
	name := "app's name $literal `literal`"
	options.Binaries[0].Name = name
	for i := range options.Binaries {
		options.Binaries[i].Path = options.path(options.Binaries[i].Path)
	}
	appdir, _, err := options.prepareDirectory(context.Background(), filepath.Join(t.TempDir(), "directory's name"), "", "")
	if err != nil {
		t.Fatal(err)
	}
	probe := `#!/bin/sh
printf '%s\000' "$PWD" "$OWD" "$LD_LIBRARY_PATH" "$XDG_DATA_DIRS" "$GSETTINGS_SCHEMA_DIR" "$WEBKIT_INJECTED_BUNDLE_PATH" "$PATH" "$@"
`
	if err := os.WriteFile(filepath.Join(appdir, "usr", "bin", name), []byte(probe), 0755); err != nil {
		t.Fatal(err)
	}
	bundle := filepath.Join(appdir, "usr", "lib", "webkitgtk-6.0", "injected-bundle")
	if err := os.MkdirAll(bundle, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bundle, "injected.so"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	working := t.TempDir()
	args := []string{"one argument", "it's literal", "$(touch unwanted)", ""}
	cmd := exec.Command(filepath.Join(appdir, "AppRun"), args...)
	cmd.Dir = working
	cmd.Env = []string{"PATH=/usr/bin:/bin", "LD_LIBRARY_PATH=/existing/lib", "XDG_DATA_DIRS=/existing/share"}
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("launcher: %v: %s", err, output)
	}
	usr := filepath.Join(appdir, "usr")
	want := []string{filepath.Join(usr, "lib"), working, usr + "/lib:" + usr + "/lib/x86_64-linux-gnu:" + usr + "/lib/aarch64-linux-gnu:" + usr + "/lib64:/existing/lib", usr + "/share:/existing/share", usr + "/share/glib-2.0/schemas", bundle + "/", usr + "/bin:/usr/bin:/bin"}
	want = append(want, args...)
	got := strings.Split(strings.TrimSuffix(string(output), "\x00"), "\x00")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v; want %#v", got, want)
	}
}

type launcherDownloadTransport func(*http.Request) (*http.Response, error)

func (f launcherDownloadTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestPrepareToolsDownloadsOnlyLinuxdeploy(t *testing.T) {
	calls := 0
	options := Options{Root: t.TempDir(), Arch: "amd64", Client: &http.Client{Transport: launcherDownloadTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.String() != "https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous/linuxdeploy-x86_64.AppImage" {
			t.Errorf("unexpected download: %s", r.URL)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("tool")), Header: make(http.Header)}, nil
	})}}
	for range 2 {
		tools, err := options.PrepareTools(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		if tools.AppRun != "" || tools.LinuxDeploy == "" {
			t.Fatalf("unexpected tools: %+v", tools)
		}
	}
	if calls != 1 {
		t.Fatalf("got %d downloads, want one cached linuxdeploy", calls)
	}
	options.Tools = Tools{LinuxDeploy: "custom-linuxdeploy", AppRun: "custom-AppRun"}
	tools, err := options.PrepareTools(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if tools.AppRun != filepath.Join(options.Root, "custom-AppRun") || tools.LinuxDeploy != filepath.Join(options.Root, "custom-linuxdeploy") || calls != 1 {
		t.Fatalf("custom overrides lost: %+v, downloads=%d", tools, calls)
	}
}
