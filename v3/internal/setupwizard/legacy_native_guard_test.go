package setupwizard

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLegacySetupRejectsNativeConfiguration(t *testing.T) {
	t.Setenv("WAILS_EXP_USE_WAKE", "")
	if err := os.Unsetenv("WAILS_EXP_USE_WAKE"); err != nil {
		t.Fatal(err)
	}
	for _, data := range []string{
		"",
		"# incomplete native configuration\n",
		"version: 3\nproject:\n  name: probe\n  product_name: Probe\n  identifier: org.wails.probe\n  version: 1.0.0\n",
		"project: [broken yaml\n",
	} {
		t.Run(data, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, "wails.yaml")
			if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
				t.Fatal(err)
			}
			nested := filepath.Join(root, "nested")
			if err := os.Mkdir(nested, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(nested, "go.mod"), []byte("module nested\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			t.Chdir(nested)
			wizard := New()
			for _, method := range []string{http.MethodGet, http.MethodPost} {
				result := httptest.NewRecorder()
				wizard.handleWailsConfig(result, httptest.NewRequest(method, "/api/wails-config", strings.NewReader(`{"info":{"productName":"Changed"}}`)))
				if result.Code != http.StatusConflict || !strings.Contains(result.Body.String(), "WAILS_EXP_USE_WAKE") {
					t.Fatalf("%s: status %d, body %s", method, result.Code, result.Body.String())
				}
			}
			after, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(after, []byte(data)) {
				t.Fatal("legacy setup changed native configuration")
			}
		})
	}
}

func TestLegacySetupUsesNearestLegacyFile(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "wails.yaml"), []byte("version: 3\nproject: {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(root, "legacy")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "wails.yaml"), []byte("info:\n  productName: Legacy\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(nested)
	result := httptest.NewRecorder()
	New().handleLegacyWailsConfig(result, httptest.NewRequest(http.MethodPost, "/api/wails-config", strings.NewReader(`{"info":{"productName":"Changed"}}`)))
	if result.Code != http.StatusOK {
		t.Fatalf("nearest legacy metadata rejected: %d %s", result.Code, result.Body.String())
	}
	data, err := os.ReadFile(filepath.Join(root, "wails.yaml"))
	if err != nil || string(data) != "version: 3\nproject: {}\n" {
		t.Fatalf("ancestor manifest changed: %q, %v", data, err)
	}
}

func TestLegacySetupDoesNotReplaceUnreadableManifest(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	if err := os.Mkdir("wails.yaml", 0o755); err != nil {
		t.Fatal(err)
	}
	result := httptest.NewRecorder()
	New().handleLegacyWailsConfig(result, httptest.NewRequest(http.MethodPost, "/api/wails-config", strings.NewReader(`{"info":{"productName":"Changed"}}`)))
	if result.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", result.Code)
	}
	info, err := os.Stat("wails.yaml")
	if err != nil || !info.IsDir() {
		t.Fatalf("invalid manifest path changed: %v", err)
	}
}
