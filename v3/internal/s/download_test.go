package s

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// DOWNLOAD exits the process on failure, so exercise it in a child test process.
func TestDownloadProcess(t *testing.T) {
	if os.Getenv("WAILS_TEST_DOWNLOAD_PROCESS") != "1" {
		return
	}
	DOWNLOAD(os.Getenv("WAILS_TEST_DOWNLOAD_URL"), os.Getenv("WAILS_TEST_DOWNLOAD_TARGET"))
}

func TestDownloadHTTPStatus(t *testing.T) {
	cases := []struct {
		status        int
		initial, want string
		wantExit      int
	}{
		{200, "", "response body", 0},
		{200, "original", "response body", 0},
		{302, "", "response body", 0},
		{302, "original", "response body", 0},
		{404, "", "", 1},
		{404, "original", "original", 1},
		{500, "", "", 1},
		{500, "original", "original", 1},
	}
	for _, tc := range cases {
		status := tc.status
		t.Run(fmt.Sprintf("%d/existing=%t", status, tc.initial != ""), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if status == http.StatusFound && r.URL.Path != "/file" {
					http.Redirect(w, r, "/file", http.StatusFound)
					return
				}
				if status != http.StatusFound {
					w.WriteHeader(status)
				}
				_, _ = w.Write([]byte("response body"))
			}))
			defer server.Close()
			target := filepath.Join(t.TempDir(), "tool")
			if tc.initial != "" {
				if err := os.WriteFile(target, []byte(tc.initial), 0600); err != nil {
					t.Fatal(err)
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=^TestDownloadProcess$")
			cmd.Env = append(os.Environ(), "WAILS_TEST_DOWNLOAD_PROCESS=1", "WAILS_TEST_DOWNLOAD_URL="+server.URL, "WAILS_TEST_DOWNLOAD_TARGET="+target)
			output, err := cmd.CombinedOutput()
			data, readErr := os.ReadFile(target)
			if cmd.ProcessState == nil {
				t.Fatalf("start child process: %v", err)
			}
			if got := cmd.ProcessState.ExitCode(); got != tc.wantExit {
				t.Errorf("exit code=%d, want %d: %s", got, tc.wantExit, output)
			}
			if tc.wantExit != 0 && !strings.Contains(string(output), http.StatusText(status)) {
				t.Errorf("missing HTTP status in error: %s", output)
			}
			if tc.want == "" {
				if !errors.Is(readErr, os.ErrNotExist) {
					t.Errorf("failed download created destination: data=%q, error=%v", data, readErr)
				}
				return
			}
			if readErr != nil || string(data) != tc.want {
				t.Errorf("destination=%q, want %q; read error=%v", data, tc.want, readErr)
			}
		})
	}
}

func TestDownloadInterruptedPreservesTarget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		_, _ = io.WriteString(w, "partial")
	}))
	t.Cleanup(server.Close)
	for _, existing := range []bool{false, true} {
		t.Run(fmt.Sprintf("existing=%t", existing), func(t *testing.T) {
			dir := t.TempDir()
			target := filepath.Join(dir, "tool")
			if existing {
				if err := os.WriteFile(target, []byte("original"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			cmd := exec.Command(os.Args[0], "-test.run=^TestDownloadProcess$")
			cmd.Env = append(os.Environ(), "WAILS_TEST_DOWNLOAD_PROCESS=1", "WAILS_TEST_DOWNLOAD_URL="+server.URL, "WAILS_TEST_DOWNLOAD_TARGET="+target)
			output, err := cmd.CombinedOutput()
			if err == nil || cmd.ProcessState == nil || cmd.ProcessState.ExitCode() != 1 {
				t.Fatalf("expected failed download: %v: %s", err, output)
			}
			if !strings.Contains(string(output), "unexpected EOF") {
				t.Fatalf("missing transfer failure: %s", output)
			}
			data, err := os.ReadFile(target)
			if existing {
				if err != nil || string(data) != "original" {
					t.Errorf("existing target=%q, error=%v", data, err)
				}
			} else if !errors.Is(err, os.ErrNotExist) {
				t.Errorf("partial target remains: %q, error=%v", data, err)
			}
			leftovers, err := filepath.Glob(filepath.Join(dir, ".wails-download-*"))
			if err != nil || len(leftovers) != 0 {
				t.Errorf("staging leftovers=%v, error=%v", leftovers, err)
			}
		})
	}
}

func TestDownloadPublishFailureCleansStaging(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, "complete") }))
	t.Cleanup(server.Close)
	dir := t.TempDir()
	target := filepath.Join(dir, "directory")
	if err := os.Mkdir(target, 0700); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(target, "existing")
	if err := os.WriteFile(marker, []byte("original"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := downloadFile(server.URL, target); err == nil {
		t.Fatal("expected directory replacement to fail")
	}
	data, err := os.ReadFile(marker)
	if err != nil || string(data) != "original" {
		t.Fatalf("existing contents=%q, error=%v", data, err)
	}
	leftovers, err := filepath.Glob(filepath.Join(dir, ".wails-download-*"))
	if err != nil || len(leftovers) != 0 {
		t.Fatalf("staging leftovers=%v, error=%v", leftovers, err)
	}
}

func TestDownloadReplacementPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix permission bits")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, "complete") }))
	t.Cleanup(server.Close)
	for _, mode := range []os.FileMode{0600, 0755} {
		target := filepath.Join(t.TempDir(), "tool")
		if err := os.WriteFile(target, []byte("original"), mode); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(target, mode); err != nil {
			t.Fatal(err)
		}
		if err := downloadFile(server.URL, target); err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(target)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != mode {
			t.Errorf("mode=%o, want %o", info.Mode().Perm(), mode)
		}
		data, err := os.ReadFile(target)
		if err != nil || string(data) != "complete" {
			t.Fatalf("replacement=%q, error=%v", data, err)
		}
	}
}
