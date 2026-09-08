package setupwizard

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetupWithoutExperimentPreservesLegacyConfiguration(t *testing.T) {
	t.Setenv("WAILS_EXP_USE_WAKE", "restore")
	require.NoError(t, os.Unsetenv("WAILS_EXP_USE_WAKE"))
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.WriteFile("Taskfile.yml", []byte("legacy tasks"), 0644))
	require.NoError(t, os.WriteFile("wails.hcl", []byte("existing HCL must stay untouched"), 0644))
	require.NoError(t, os.WriteFile("wails.yaml", []byte("# keep header\ninfo:\n  productName: Before # keep field comment\n  customInfo: keep\ncustomTop: keep\n"), 0644))
	nested := filepath.Join(root, "nested")
	require.NoError(t, os.Mkdir(nested, 0755))
	t.Chdir(nested)
	wizard := New()
	result := httptest.NewRecorder()
	wizard.handleWailsConfig(result, httptest.NewRequest(http.MethodGet, "/api/wails-config", nil))
	require.Equal(t, http.StatusOK, result.Code)
	require.Contains(t, result.Body.String(), "Before")
	result = httptest.NewRecorder()
	wizard.handleWailsConfig(result, httptest.NewRequest(http.MethodPost, "/api/wails-config", bytes.NewBufferString(`{"info":{"productName":"After"}}`)))
	require.Equal(t, http.StatusOK, result.Code)
	data, err := os.ReadFile(filepath.Join(root, "wails.yaml"))
	require.NoError(t, err)
	require.Contains(t, string(data), "After")
	for _, expected := range []string{"# keep header", "# keep field comment", "customInfo: keep", "customTop: keep"} {
		require.Contains(t, string(data), expected)
	}
	for name, expected := range map[string]string{"Taskfile.yml": "legacy tasks", "wails.hcl": "existing HCL must stay untouched"} {
		data, err := os.ReadFile(filepath.Join(root, name))
		require.NoError(t, err)
		require.Equal(t, expected, string(data))
	}
}

func TestLegacySetupInitialisesEmptyAndCommentOnlyYAML(t *testing.T) {
	for _, original := range []string{"", "# keep this header", "# keep this header\n\n"} {
		t.Run(original, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "wails.yaml")
			require.NoError(t, os.WriteFile(path, []byte(original), 0644))
			data, err := mergeLegacyProjectInfo(path, WailsConfigInfo{ProductName: "New App"})
			require.NoError(t, err)
			require.Contains(t, string(data), "productName: New App")
			if original != "" {
				require.Contains(t, string(data), "# keep this header\n")
			}
		})
	}
}
