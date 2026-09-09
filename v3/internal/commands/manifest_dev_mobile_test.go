package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/dev"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
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

func TestMobileDevTracksSuccessfullyLaunchedArtifact(t *testing.T) {
	for _, format := range []string{"apk", "app"} {
		t.Run(format, func(t *testing.T) {
			root := t.TempDir()
			output := "app." + format
			file := filepath.Join(root, output)
			if format == "app" {
				require.NoError(t, os.MkdirAll(file, 0755))
				file = filepath.Join(file, "Info.plist")
			}
			require.NoError(t, os.WriteFile(file, []byte("initial"), 0644))
			plan := pipeline.Plan{Artifacts: []pipeline.NodeKey{"installable"}, Nodes: map[pipeline.NodeKey]pipeline.Node{
				"installable": {Output: output, Artifact: pipeline.ArtifactIdentity{Format: format}},
			}}
			var readyErr error
			ops := dev.Operations{
				Build: func(context.Context, *manifest.Loaded, string, string, string, int) (manifestPipelineRun, error) {
					// No reusable digest, as with real iOS assembly/signing.
					return manifestPipelineRun{Plan: plan}, nil
				},
				BinaryPath: func(root string, run manifestPipelineRun, goos, arch string) (string, error) {
					return mobileDevArtifact(root, run, goos, arch)
				},
				WaitReady: func(context.Context, dev.Process, time.Duration) error { return readyErr },
			}
			trackMobileDevArtifact(&ops)
			loaded := &manifest.Loaded{Config: manifest.Config{Root: root}}
			build := func() manifestPipelineRun {
				run, err := ops.Build(context.Background(), loaded, "ios", "arm64", "", 0)
				require.NoError(t, err)
				return run
			}
			activate := func(run manifestPipelineRun) error {
				_, err := ops.BinaryPath(root, run, "ios", "arm64")
				require.NoError(t, err)
				return ops.WaitReady(context.Background(), nil, time.Second)
			}
			first := build()
			require.True(t, ops.BackendChanged(first, "ios", "arm64"))
			require.NoError(t, activate(first))
			// A manifest-only debounce edit must keep the installed process.
			loaded.Config.Dev.DebounceMS++
			require.False(t, ops.BackendChanged(build(), "ios", "arm64"))
			// Packaging-only changes must restart even without a compile change.
			require.NoError(t, os.WriteFile(file, []byte("packaging changed"), 0644))
			changed := build()
			require.True(t, ops.BackendChanged(changed, "ios", "arm64"))
			readyErr = fmt.Errorf("launch failed")
			require.Error(t, activate(changed))
			require.True(t, ops.BackendChanged(build(), "ios", "arm64"), "failed replacement must remain retryable")
			require.False(t, ops.BackendChanged(first, "ios", "arm64"), "rollback keeps the previous identity")
			readyErr = nil
			require.NoError(t, activate(changed))
			require.False(t, ops.BackendChanged(build(), "ios", "arm64"))
			// Restoring an older cache artifact must replace the currently installed one.
			require.True(t, ops.BackendChanged(first, "ios", "arm64"))
			// Cached APK identities also participate, independent of lookup status.
			cached := manifestPipelineRun{Plan: plan, Results: map[pipeline.NodeKey]pipeline.Result{
				"installable": {Artifact: "cached-apk", Status: cache.LookupRestored},
			}}
			require.True(t, ops.BackendChanged(cached, "android", "arm64"))
			require.NoError(t, activate(cached))
			require.False(t, ops.BackendChanged(cached, "android", "arm64"))
			require.True(t, ops.BackendChanged(manifestPipelineRun{}, "ios", "arm64"))
			require.NoError(t, os.RemoveAll(filepath.Join(root, output)))
			_, err := ops.Build(context.Background(), loaded, "ios", "arm64", "", 0)
			require.Error(t, err, "missing installable output must not be accepted as an unchanged app")
		})
	}
}
