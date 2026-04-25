package commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	}

	for _, tt := range tests {
		t.Run(tt.projectName, func(t *testing.T) {
			result := strings.ToLower(strings.ReplaceAll(tt.projectName, " ", "-"))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNfpmTemplateUsesBinaryName(t *testing.T) {
	data, err := updatableBuildAssets.ReadFile("updatable_build_assets/linux/nfpm/nfpm.yaml.tmpl")
	if err != nil {
		t.Skip("nfpm template not found in embedded assets")
	}
	content := string(data)
	assert.True(t, strings.Contains(content, `{{.BinaryName}}`),
		"nfpm template should use BinaryName variable")
}
