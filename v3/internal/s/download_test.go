package s

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
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
