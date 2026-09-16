package templates

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/flags"
)

func TestExperimentalTemplateInstallsWithLegacyTaskfilesWhenDisabled(t *testing.T) {
	t.Setenv("WAILS_EXP_USE_WAKE", "anything")
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, GenerateTemplate(&BaseTemplate{Name: "portable", Author: "Test", Version: "1.0.0", Dir: root}))
	template := filepath.Join(root, "portable")
	require.NoFileExists(t, filepath.Join(template, "Taskfile.tmpl.yml"))
	require.NoError(t, os.Unsetenv("WAILS_EXP_USE_WAKE"))
	for _, custom := range []bool{false, true} {
		t.Run(map[bool]string{false: "missing-taskfile", true: "custom-taskfile"}[custom], func(t *testing.T) {
			t.Chdir(root)
			if custom {
				require.NoError(t, os.WriteFile(filepath.Join(template, "Taskfile.yaml"), []byte("version: '3'\n# keep custom tasks\n"), 0644))
			}
			options := &flags.Init{TemplateName: template, ProjectName: "portable", ProjectDir: t.TempDir(), ModulePath: "example.com/portable", SkipGoModTidy: true}
			require.NoError(t, Install(options))
			dir := options.ProjectDir
			if custom {
				require.NoFileExists(t, filepath.Join(dir, "Taskfile.yml"))
				data, err := os.ReadFile(filepath.Join(dir, "Taskfile.yaml"))
				require.NoError(t, err)
				require.Contains(t, string(data), "# keep custom tasks")
			} else {
				require.FileExists(t, filepath.Join(dir, "Taskfile.yml"))
			}
		})
	}
}
