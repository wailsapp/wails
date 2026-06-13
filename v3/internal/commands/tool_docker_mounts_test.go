package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// captureStdout runs fn with os.Stdout swapped for a pipe and returns whatever fn wrote.
func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fnErr := fn()
	if err := w.Close(); err != nil {
		t.Fatalf("closing pipe: %v", err)
	}
	os.Stdout = oldStdout
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("reading stdout: %v", err)
	}
	return buf.String(), fnErr
}

// setupDockerMountsProject prepares a temp project layout and chdirs into it.
// Layout:
//
//	<parent>/
//	  project/         (returned as cwd)
//	    go.mod         (caller writes contents)
//	    internal/foo/  (target for relative-inside replace)
//	  sibling/         (target for relative-outside replace)
//	  absolute/        (target for absolute replace)
//	  fake-gopath/     (GOPATH for deterministic /go/pkg/mod mount)
func setupDockerMountsProject(t *testing.T) (project, sibling, absolute, fakeGopath string) {
	t.Helper()
	parent := t.TempDir()
	project = filepath.Join(parent, "project")
	sibling = filepath.Join(parent, "sibling")
	absolute = filepath.Join(parent, "absolute")
	fakeGopath = filepath.Join(parent, "fake-gopath")
	for _, d := range []string{
		filepath.Join(project, "internal", "foo"),
		sibling,
		absolute,
		fakeGopath,
	} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}
	t.Setenv("GOPATH", fakeGopath)
	t.Chdir(project)
	return project, sibling, absolute, fakeGopath
}

func writeGoMod(t *testing.T, dir, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(body), 0o644); err != nil {
		t.Fatalf("writing go.mod: %v", err)
	}
}

func TestToolDockerMounts_RelativeReplaces(t *testing.T) {
	project, _, _, fakeGopath := setupDockerMountsProject(t)
	writeGoMod(t, project, `module example.com/app

go 1.25

replace foo => ../sibling
replace bar => ./internal/foo
`)

	out, err := captureStdout(t, func() error { return ToolDockerMounts(&DockerMountsOptions{}) })
	if err != nil {
		t.Fatalf("ToolDockerMounts returned error: %v", err)
	}

	gopathMount := fmt.Sprintf(`-v "%s/pkg/mod:/go/pkg/mod"`, filepath.ToSlash(fakeGopath))
	siblingHost := filepath.ToSlash(filepath.Join(filepath.Dir(project), "sibling"))
	internalFooHost := filepath.ToSlash(filepath.Join(project, "internal", "foo"))
	wantSibling := fmt.Sprintf(`-v "%s:/sibling:ro"`, siblingHost)
	wantInternal := fmt.Sprintf(`-v "%s:/app/internal/foo:ro"`, internalFooHost)

	for _, want := range []string{gopathMount, wantSibling, wantInternal} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\nfull output: %s", want, out)
		}
	}
}

func TestToolDockerMounts_AbsoluteReplace(t *testing.T) {
	// Skip on Windows because we'd need a non-drive-letter absolute path,
	// which isn't a thing on Windows hosts.
	if runtime.GOOS == "windows" {
		t.Skip("Unix absolute path semantics; covered separately by drive-letter test")
	}
	project, _, absolute, _ := setupDockerMountsProject(t)
	writeGoMod(t, project, fmt.Sprintf(`module example.com/app

go 1.25

replace foo => %s
`, absolute))

	out, err := captureStdout(t, func() error { return ToolDockerMounts(&DockerMountsOptions{}) })
	if err != nil {
		t.Fatalf("ToolDockerMounts returned error: %v", err)
	}

	absSlash := filepath.ToSlash(absolute)
	want := fmt.Sprintf(`-v "%s:%s:ro"`, absSlash, absSlash)
	if !strings.Contains(out, want) {
		t.Errorf("output missing literal-path mount %q\nfull output: %s", want, out)
	}
}

func TestToolDockerMounts_NonExistentTargetSkipped(t *testing.T) {
	project, _, _, _ := setupDockerMountsProject(t)
	writeGoMod(t, project, `module example.com/app

go 1.25

replace foo => ../does-not-exist
`)

	out, err := captureStdout(t, func() error { return ToolDockerMounts(&DockerMountsOptions{}) })
	if err != nil {
		t.Fatalf("ToolDockerMounts returned error: %v", err)
	}
	if strings.Contains(out, "does-not-exist") {
		t.Errorf("non-existent target should be skipped, got: %s", out)
	}
}

func TestToolDockerMounts_MissingGoModError(t *testing.T) {
	parent := t.TempDir()
	t.Setenv("GOPATH", parent)
	t.Chdir(parent)

	_, err := captureStdout(t, func() error { return ToolDockerMounts(&DockerMountsOptions{}) })
	if err == nil {
		t.Fatal("expected error when go.mod is missing, got nil")
	}
	if !strings.Contains(err.Error(), "go.mod") {
		t.Errorf("expected error to mention go.mod, got: %v", err)
	}
}

func TestToolDockerMounts_DriveLetterSkipped(t *testing.T) {
	if runtime.GOOS != "windows" {
		// filepath.IsAbs("C:\\vendor\\lib") is false on non-Windows, so the
		// drive-letter skip branch can't be exercised on a Linux/macOS runner.
		t.Skip("drive-letter handling is Windows-specific")
	}
	project, _, _, _ := setupDockerMountsProject(t)
	writeGoMod(t, project, `module example.com/app

go 1.25

replace foo => C:\vendor\lib
`)

	out, err := captureStdout(t, func() error { return ToolDockerMounts(&DockerMountsOptions{}) })
	if err != nil {
		t.Fatalf("ToolDockerMounts returned error: %v", err)
	}
	if strings.Contains(out, `vendor/lib`) {
		t.Errorf("drive-letter path should be skipped, got: %s", out)
	}
}

func TestToolHas(t *testing.T) {
	// Create a temp dir with a stub executable so we control exactly what is
	// and isn't in PATH — no reliance on the host environment.
	dir := t.TempDir()
	stubName := createStubExecutable(t, dir)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	tests := []struct {
		tool    string
		want    string
		wantErr bool
	}{
		{tool: stubName, want: "true"},                            // found: stub is in PATH
		{tool: "nonexistent-wails-tool-xyz", want: "false"},      // not found
		{tool: "nonexistent-wails-tool-xyz|" + stubName, want: "true"}, // second alternative found
		{tool: stubName + "|nonexistent-wails-tool-xyz", want: "true"}, // first alternative found
		{tool: "", wantErr: true},                                 // missing arg: error
	}

	for _, tc := range tests {
		out, err := captureStdout(t, func() error { return ToolHas(&HasOptions{Tool: tc.tool}) })
		if tc.wantErr {
			if err == nil {
				t.Errorf("ToolHas(%q): expected error, got nil (output %q)", tc.tool, out)
			}
			continue
		}
		if err != nil {
			t.Errorf("ToolHas(%q): unexpected error: %v", tc.tool, err)
			continue
		}
		if out != tc.want {
			t.Errorf("ToolHas(%q): got %q, want %q", tc.tool, out, tc.want)
		}
	}

	// Deprecated alias: put a gcc stub in the same temp dir so the result is
	// deterministic regardless of what the host has installed.
	createNamedStub(t, dir, "gcc")
	out1, err1 := captureStdout(t, func() error { return ToolHasCC(&HasCCOptions{}) })
	out2, err2 := captureStdout(t, func() error { return ToolHas(&HasOptions{Tool: "gcc|clang"}) })
	if err1 != nil || err2 != nil {
		t.Errorf("ToolHasCC/ToolHas errors: %v / %v", err1, err2)
	}
	if out1 != "true" || out2 != "true" {
		t.Errorf("expected \"true\" with gcc stub in PATH, got ToolHasCC=%q ToolHas=%q", out1, out2)
	}
}

// createStubExecutable writes a minimal executable named "wails-test-stub" in
// dir and returns its base name. LookPath only checks existence and the
// executable bit, so content is irrelevant.
func createStubExecutable(t *testing.T, dir string) string {
	t.Helper()
	return createNamedStub(t, dir, "wails-test-stub")
}

// createNamedStub writes a minimal executable with the given name in dir and
// returns the base name (without any platform extension).
func createNamedStub(t *testing.T, dir, name string) string {
	t.Helper()
	filename := filepath.Join(dir, name)
	if runtime.GOOS == "windows" {
		filename += ".exe"
	}
	if err := os.WriteFile(filename, []byte{}, 0755); err != nil {
		t.Fatalf("creating stub executable %q: %v", name, err)
	}
	return name
}
