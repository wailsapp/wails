package templates

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
)

func TestTemplateOptionsBinaryName(t *testing.T) {
	tests := []struct {
		projectName     string
		expectedBinName string
	}{
		{"Wails", "wails"},
		{"MyApp", "myapp"},
		{"My App", "my-app"},
		{"My Cool App", "my-cool-app"},
		{"wails", "wails"},
		{"already-lower", "already-lower"},
	}

	for _, tt := range tests {
		t.Run(tt.projectName, func(t *testing.T) {
			options := &flags.Init{
				ProjectName: tt.projectName,
			}
			templateData := TemplateOptions{
				Init:       options,
				BinaryName: strings.ToLower(strings.ReplaceAll(options.ProjectName, " ", "-")),
			}
			if templateData.BinaryName != tt.expectedBinName {
				t.Errorf("BinaryName = %q, want %q", templateData.BinaryName, tt.expectedBinName)
			}
		})
	}
}

func TestTaskfileTemplateUsesBinaryName(t *testing.T) {
	subFS, err := fs.Sub(templates, "_common")
	if err != nil {
		t.Fatalf("Failed to get _common sub FS: %v", err)
	}

	data, err := fs.ReadFile(subFS, "Taskfile.tmpl.yml")
	if err != nil {
		t.Fatalf("Failed to read Taskfile template: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "{{.BinaryName}}") {
		t.Error("Taskfile.tmpl.yml should use {{.BinaryName}} for APP_NAME, not {{.ProjectName}}")
	}

	if strings.Contains(content, `APP_NAME: "{{.ProjectName}}"`) {
		t.Error("Taskfile.tmpl.yml should NOT use {{.ProjectName}} for APP_NAME")
	}
}
