package commands

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
	"howett.net/plist"
)

func TestMobileDevPlansRemainDevelopmentBuilds(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "wails.hcl"), []byte("version = 3\nproject {\n name = \"devtest\"\n product_name = \"Dev Test\"\n version = \"1.0.0\"\n identifier = \"org.wails.devtest\"\n}\n"), 0644))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	for _, platform := range []string{"ios", "android"} {
		for _, destination := range []string{"simulator", "device"} {
			if platform == "android" && destination == "device" {
				continue
			}
			options := &DevOptions{}
			if platform == "ios" {
				options.Destination = destination
			}
			run := mobileDevBuildOptions(context.Background(), options, loaded, platform, "arm64", "http://127.0.0.1:9353", 9353)
			plan, err := pipeline.PlanBuild(run.Loaded.Config, pipeline.Request{Verb: run.Verb, TargetOS: run.TargetOS, TargetArch: run.TargetArch, Development: run.Development, Formats: run.Formats})
			require.NoError(t, err)
			for _, node := range plan.Nodes {
				switch spec := node.Spec.(type) {
				case pipeline.CompileSpec:
					require.False(t, spec.Production)
					require.True(t, strings.HasPrefix(node.Output, ".wails/dev/"))
				case pipeline.PackageSpec:
					require.True(t, spec.Development)
					require.True(t, strings.HasPrefix(node.Output, ".wails/dev/"))
				}
			}
			require.Equal(t, "", loaded.Config.Selected.Name, "dev selection must not mutate loaded configuration")
		}
	}
}
func TestIOSDevNetworkOnlyChangesStagedPlist(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source.plist")
	stage := filepath.Join(root, "stage.plist")
	data, err := plist.Marshal(map[string]any{"CFBundleIdentifier": "org.wails.test", "NSLocalNetworkUsageDescription": "Custom explanation", "NSAppTransportSecurity": map[string]any{"NSAllowsArbitraryLoads": false}}, plist.XMLFormat)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(source, data, 0644))
	require.NoError(t, copyManifestPath(source, stage))
	require.NoError(t, prepareIOSDevNetwork(stage))
	original, err := os.ReadFile(source)
	require.NoError(t, err)
	require.Equal(t, data, original)
	result, err := os.ReadFile(stage)
	require.NoError(t, err)
	var fields map[string]any
	_, err = plist.Unmarshal(result, &fields)
	require.NoError(t, err)
	require.Equal(t, "Custom explanation", fields["NSLocalNetworkUsageDescription"])
	require.Equal(t, false, fields["NSAppTransportSecurity"].(map[string]any)["NSAllowsArbitraryLoads"])
}
func TestAndroidDevNetworkPolicyIsDebugOnlyAndLoopbackOnly(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, prepareAndroidDevNetwork(root))
	require.NoFileExists(t, filepath.Join(root, "app/src/main/AndroidManifest.xml"))
	data, err := os.ReadFile(filepath.Join(root, "app/src/debug/res/xml/wails_dev_network.xml"))
	require.NoError(t, err)
	require.Contains(t, string(data), `<base-config cleartextTrafficPermitted="false"/>`)
	require.Contains(t, string(data), `<domain>127.0.0.1</domain>`)
}
func TestDevelopmentPackageWorkspacesDoNotOverlapProduction(t *testing.T) {
	h := manifestHandler{root: t.TempDir()}
	s := pipeline.PackageSpec{TargetOS: "android", TargetArch: "arm64", Format: "apk"}
	production := h.packageWorkspace(s)
	s.Development = true
	development := h.packageWorkspace(s)
	require.NotEqual(t, production, development)
	require.Contains(t, filepath.ToSlash(development), "/.wails/dev/android-arm64/")
}
