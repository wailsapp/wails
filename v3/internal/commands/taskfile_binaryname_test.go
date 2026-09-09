package commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/templates"
)

func TestBinaryNameNormalization(t *testing.T) {
	tests := []struct {
		projectName string
		expected    string
	}{
		{"Wails", "wails"},
		{"MyApp", "myapp"},
		{"Hello World", "hello-world"},
		{"Test-App", "test-app"},
		{"lowercase", "lowercase"},
		{"My_App!", "my-app"},
		{"  leading  ", "leading"},
	}

	for _, tt := range tests {
		t.Run(tt.projectName, func(t *testing.T) {
			result := templates.NormalizeBinaryName(tt.projectName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNfpmTemplateUsesBinaryName(t *testing.T) {
	data, err := updatableBuildAssets.ReadFile("updatable_build_assets/linux/nfpm/nfpm.yaml.tmpl")
	require.NoError(t, err, "nfpm template should be present in embedded assets")
	content := string(data)
	assert.True(t, strings.Contains(content, `{{.BinaryName}}`),
		"nfpm template should use BinaryName variable")
}
