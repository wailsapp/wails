package commands

import (
	"debug/pe"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tc-hib/winres"
	"github.com/tc-hib/winres/version"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestManifestWindowsResourcesLinkFromGeneratedPackage(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/windows-resources\ngo 1.25.0\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nfunc main() {}\n"), 0644))
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "resources", ProductName: "Resources", Identifier: "com.example.resources", Version: "3.7.2"}))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	handler := &manifestHandler{root: root, config: loaded.Config}
	for _, arch := range []string{"amd64", "arm64"} {
		t.Run(arch, func(t *testing.T) {
			plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build", TargetOS: "windows", TargetArch: arch})
			require.NoError(t, err)
			_, err = handler.Run(t.Context(), plan.Nodes[pipeline.NodeKey("target:windows/"+arch+":assets")])
			require.NoError(t, err)
			tools, err := os.ReadFile(filepath.Join(root, plan.Nodes[pipeline.NodeKey("target:windows/"+arch+":assets")].Output, "windows", "nsis", "wails_tools.nsh"))
			require.NoError(t, err)
			require.True(t, strings.Contains(string(tools), `!define INFO_PRODUCTVERSION "3.7.2"`))
			node := plan.Nodes[pipeline.NodeKey("target:windows/"+arch+":compile")]
			result, err := handler.Run(t.Context(), node)
			require.NoError(t, err, result.Detail)
			binary, err := pe.Open(filepath.Join(root, node.Output))
			require.NoError(t, err)
			defer binary.Close()
			require.NotNil(t, binary.Section(".rsrc"), "resource object must reach the final executable")
			file, err := os.Open(filepath.Join(root, node.Output))
			require.NoError(t, err)
			defer file.Close()
			resources, err := winres.LoadFromEXE(file)
			require.NoError(t, err)
			info, err := version.FromBytes(resources.Get(winres.RT_VERSION, winres.ID(1), winres.LCIDNeutral))
			require.NoError(t, err)
			require.Equal(t, [4]uint16{3, 7, 2, 0}, info.FileVersion)
			require.Equal(t, info.FileVersion, info.ProductVersion)
			table := *info.Table()[version.LangNeutral]
			require.Equal(t, "3.7.2", table[version.FileVersion])
			require.Equal(t, "Resources", table[version.ProductName])
		})
	}
	require.NoFileExists(t, filepath.Join(root, "wails_windows_resources.go"))
	matches, err := filepath.Glob(filepath.Join(root, "*.syso"))
	require.NoError(t, err)
	require.Empty(t, matches)
}
