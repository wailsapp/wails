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

func TestWizardDoesNotExposeProjectConfigurationWrites(t *testing.T) {
	for _, wake := range []bool{false, true} {
		for _, init := range []bool{false, true} {
			for name, config := range map[string]string{
				"absent":  "",
				"native":  "# preserve this comment\nversion: 3\nproject:\n  name: existing\n  identifier: com.example.existing\n  version: 1.0.0\n",
				"legacy":  "# preserve legacy settings\ninfo:\n  productName: Existing\ncustom: keep\n",
				"invalid": "project: [broken yaml\n",
			} {
				mode := "global"
				if init {
					mode = "init"
				}
				experiment := "legacy"
				if wake {
					experiment = "wake"
				}
				t.Run(mode+"/"+experiment+"/"+name, func(t *testing.T) {
					t.Setenv("WAILS_EXP_USE_WAKE", "1")
					if !wake {
						if err := os.Unsetenv("WAILS_EXP_USE_WAKE"); err != nil {
							t.Fatal(err)
						}
					}
					root := t.TempDir()
					path := filepath.Join(root, "wails.yaml")
					if config != "" {
						if err := os.WriteFile(path, []byte(config), 0600); err != nil {
							t.Fatal(err)
						}
					}
					nested := filepath.Join(root, "nested")
					if err := os.Mkdir(nested, 0700); err != nil {
						t.Fatal(err)
					}
					t.Chdir(nested)
					wizard := New()
					if init {
						wizard = NewInitWizard(InitData{ProjectName: "new-app", BaseDir: nested})
					}
					mux := http.NewServeMux()
					wizard.setupRoutes(mux)
					for _, method := range []string{http.MethodGet, http.MethodPost} {
						response := httptest.NewRecorder()
						mux.ServeHTTP(response, httptest.NewRequest(method, "/api/wails-config", strings.NewReader(`{"info":{"productName":"Changed"}}`)))
						if response.Code != http.StatusNotFound {
							t.Fatalf("%s obsolete project endpoint: status %d, body %s", method, response.Code, response.Body.String())
						}
					}
					data, err := os.ReadFile(path)
					if config == "" {
						if !os.IsNotExist(err) {
							t.Fatalf("created parent configuration: %v", err)
						}
					} else if err != nil || !bytes.Equal(data, []byte(config)) {
						t.Fatalf("changed existing configuration: %q, %v", data, err)
					}
					entries, err := os.ReadDir(nested)
					if err != nil || len(entries) != 0 {
						t.Fatalf("created project files through obsolete endpoint: %v, %v", entries, err)
					}
				})
			}
		}
	}
}
