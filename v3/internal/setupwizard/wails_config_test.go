package setupwizard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestWailsConfigEndpointWritesValidHCLFromProjectState(t *testing.T) {
	originalDirectory, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, os.Chdir(originalDirectory)) })
	root := t.TempDir()
	require.NoError(t, os.Chdir(root))

	payload := `{
  "info": {
    "companyName": "Acme Limited",
    "productName": "Wizard Product",
    "productIdentifier": "com.example.wizard",
    "description": "Created through the setup wizard",
    "copyright": "Copyright 2026 Acme Limited",
    "comments": "Internal preview",
    "version": "2.3.4"
  },
  "bindings": {
    "typescript": false,
    "interfaces": true
  }
}`
	mux := http.NewServeMux()
	New().setupRoutes(mux)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/wails-config", strings.NewReader(payload)))

	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	assert.NoFileExists(t, filepath.Join(root, "wails.yaml"))
	assert.FileExists(t, filepath.Join(root, manifest.Filename))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	assert.Equal(t, "Wizard_Product", loaded.Config.Project.Name)
	assert.Equal(t, "Wizard Product", loaded.Config.Project.ProductName)
	assert.Equal(t, "Acme Limited", loaded.Config.Project.CompanyName)
	assert.Equal(t, "com.example.wizard", loaded.Config.Project.Identifier)
	assert.Equal(t, "Created through the setup wizard", loaded.Config.Project.Description)
	assert.Equal(t, "2.3.4", loaded.Config.Project.Version)
	assert.Equal(t, "Copyright 2026 Acme Limited", loaded.Config.Project.Copyright)
	assert.Equal(t, "Internal preview", loaded.Config.Project.Comments)
	assert.False(t, loaded.Config.Frontend.Bindings.TypeScript)
	assert.False(t, loaded.Config.Frontend.Bindings.Interfaces)

	getResponse := httptest.NewRecorder()
	mux.ServeHTTP(getResponse, httptest.NewRequest(http.MethodGet, "/api/wails-config", nil))
	require.Equal(t, http.StatusOK, getResponse.Code, getResponse.Body.String())
	var state ProjectConfigState
	require.NoError(t, json.NewDecoder(getResponse.Body).Decode(&state))
	assert.Equal(t, "Wizard_Product", state.Info.ProjectName)
	assert.Equal(t, "Wizard Product", state.Info.ProductName)
	assert.Equal(t, "Acme Limited", state.Info.CompanyName)
	assert.Equal(t, "com.example.wizard", state.Info.ProductIdentifier)
	assert.Equal(t, "Created through the setup wizard", state.Info.Description)
	assert.Equal(t, "2.3.4", state.Info.Version)
	assert.Equal(t, "Copyright 2026 Acme Limited", state.Info.Copyright)
	assert.Equal(t, "Internal preview", state.Info.Comments)
	require.NotNil(t, state.Bindings)
	assert.False(t, state.Bindings.TypeScript)
	assert.False(t, state.Bindings.Interfaces)
}

func TestWailsConfigEndpointPreservesBuildIntentWhenProjectStateChanges(t *testing.T) {
	originalDirectory, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, os.Chdir(originalDirectory)) })
	root := t.TempDir()
	require.NoError(t, os.Chdir(root))

	original := `version = 3

# This user-owned comment must survive wizard edits.
project {
  name         = "existing-project"
  product_name = "Old Product"
  identifier   = "com.example.old"
  version      = "1.0.0"
  binary_name  = "custom-binary"
  build_number = 42
}

frontend {
  directory = "web"
  install   = ["pnpm", "install", "--frozen-lockfile"]
  build     = ["pnpm", "run", "bundle"]
  dev       = ["pnpm", "run", "serve"]
  output    = "web/public"
}

build {
  output    = "release-bin"
  tags      = ["enterprise"]
  trim_path = false
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(original), 0o640))

	payload := `{
  "info": {
    "companyName": "New Company",
    "productName": "New Product",
    "productIdentifier": "com.example.new",
    "description": "Updated description",
    "copyright": "Copyright 2026 New Company",
    "comments": "Updated comments",
    "version": "2.0.0"
  }
}`
	mux := http.NewServeMux()
	New().setupRoutes(mux)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/wails-config", strings.NewReader(payload)))

	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	assert.Equal(t, "existing-project", loaded.Config.Project.Name)
	assert.Equal(t, "custom-binary", loaded.Config.Project.BinaryName)
	assert.Equal(t, 42, loaded.Config.Project.BuildNumber)
	assert.Equal(t, "New Product", loaded.Config.Project.ProductName)
	assert.Equal(t, "New Company", loaded.Config.Project.CompanyName)
	assert.Equal(t, "com.example.new", loaded.Config.Project.Identifier)
	assert.Equal(t, "2.0.0", loaded.Config.Project.Version)
	assert.Equal(t, "web", loaded.Config.Frontend.Directory)
	assert.Equal(t, []string{"pnpm", "install", "--frozen-lockfile"}, loaded.Config.Frontend.Install)
	assert.Equal(t, "web/public", loaded.Config.Frontend.OutputDirectory)
	assert.Equal(t, "release-bin", loaded.Config.Build.OutputDirectory)
	assert.Equal(t, []string{"enterprise"}, loaded.Config.Build.Go.Tags)
	assert.False(t, loaded.Config.Build.TrimPath)
	updated, err := os.ReadFile(filepath.Join(root, manifest.Filename))
	require.NoError(t, err)
	assert.Contains(t, string(updated), "# This user-owned comment must survive wizard edits.")
	if runtime.GOOS != "windows" {
		info, statErr := os.Stat(filepath.Join(root, manifest.Filename))
		require.NoError(t, statErr)
		assert.Equal(t, os.FileMode(0o640), info.Mode().Perm())
	}
}

func TestWailsConfigEndpointRejectsIncompleteProjectStateWithoutWritingHCL(t *testing.T) {
	originalDirectory, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, os.Chdir(originalDirectory)) })
	root := t.TempDir()
	require.NoError(t, os.Chdir(root))

	mux := http.NewServeMux()
	New().setupRoutes(mux)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/wails-config", strings.NewReader(`{
  "info": {
    "productName": "Incomplete Product",
    "version": "1.0.0"
  }
}`)))

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "productIdentifier")
	assert.NoFileExists(t, filepath.Join(root, manifest.Filename))
	assert.NoFileExists(t, filepath.Join(root, "wails.yaml"))
}
