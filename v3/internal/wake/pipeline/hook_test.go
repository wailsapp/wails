package pipeline

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestHookPlannerBuildAndPackageBarriers(t *testing.T) {
	config := testConfig(t)
	config.Hooks = allPlannerHooks()
	plan, err := PlanBuild(config, Request{Verb: "package", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"deb", "rpm"}})
	require.NoError(t, err)

	beforeBuild := plan.Nodes["hook:before_build"]
	assert.Equal(t, ProjectScope, beforeBuild.Scope)
	assert.Empty(t, beforeBuild.Dependencies)
	beforeBuildSpec := beforeBuild.Spec.(HookSpec)
	assert.Equal(t, "package", beforeBuildSpec.Command)
	assert.Equal(t, ProjectScope, beforeBuildSpec.Scope)
	assert.Equal(t, 1, beforeBuildSpec.ContextVersion)
	assert.Contains(t, plan.Nodes["frontend:install"].Dependencies, beforeBuild.Key)
	assert.Contains(t, plan.Nodes["frontend:bindings"].Dependencies, beforeBuild.Key)

	compile := NodeKey("target:linux/amd64:compile")
	afterBuild := plan.Nodes["hook:after_build:linux-amd64"]
	assert.Equal(t, []NodeKey{compile}, afterBuild.Dependencies)
	afterBuildSpec := afterBuild.Spec.(HookSpec)
	assert.Equal(t, plan.Nodes[compile].Output, afterBuildSpec.ScopeOutput)
	assert.Equal(t, TargetScope, afterBuildSpec.Scope)
	assert.Equal(t, "linux", afterBuildSpec.TargetOS)
	assert.Equal(t, "amd64", afterBuildSpec.TargetArch)

	beforePackage := plan.Nodes["hook:before_package:linux-amd64"]
	assert.Contains(t, beforePackage.Dependencies, afterBuild.Key)
	assert.Contains(t, plan.Nodes["package:linux/amd64:deb"].Dependencies, beforePackage.Key)
	assert.Contains(t, plan.Nodes["package:linux/amd64:rpm"].Dependencies, beforePackage.Key)
	afterPackage := plan.Nodes["hook:after_package:linux-amd64"]
	assert.ElementsMatch(t, []NodeKey{"package:linux/amd64:deb", "package:linux/amd64:rpm"}, afterPackage.Dependencies)
	assert.Equal(t, ".wails/build/default/linux-amd64/artifacts", afterPackage.Spec.(HookSpec).ScopeOutput)
	for _, key := range []NodeKey{"publish:linux/amd64:deb", "publish:linux/amd64:rpm"} {
		assert.Contains(t, plan.Nodes[key].Dependencies, afterPackage.Key)
	}
	assert.NotContains(t, plan.Nodes, NodeKey("hook:before_sign:linux-amd64"))
	assert.NotContains(t, plan.Nodes, NodeKey("hook:after_sign:linux-amd64"))
}

func TestHookContextContractAffectsActionIdentity(t *testing.T) {
	base := HookSpec{Phase: manifest.BeforeBuild, Command: "build", Scope: ProjectScope, ContextVersion: 1}
	baseKey, err := cache.ActionKey(string(RunHook), base, nil, nil)
	require.NoError(t, err)

	for name, changed := range map[string]HookSpec{
		"command":         {Phase: manifest.BeforeBuild, Command: "package", Scope: ProjectScope, ContextVersion: 1},
		"scope":           {Phase: manifest.BeforeBuild, Command: "build", Scope: TargetScope, ContextVersion: 1},
		"context version": {Phase: manifest.BeforeBuild, Command: "build", Scope: ProjectScope, ContextVersion: 2},
	} {
		t.Run(name, func(t *testing.T) {
			changedKey, keyErr := cache.ActionKey(string(RunHook), changed, nil, nil)
			require.NoError(t, keyErr)
			assert.NotEqual(t, baseKey, changedKey)
		})
	}
}

func TestHookPlannerSigningBarriersAndMultiTargetSharing(t *testing.T) {
	config := testConfig(t)
	config.Hooks = allPlannerHooks()
	config.Signing.Linux.Enabled = true
	config.Signing.Linux.Identity = "test"
	plan, err := PlanBuild(config, Request{Verb: "sign", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"deb", "rpm"}, sign: true})
	require.NoError(t, err)
	beforeSign := plan.Nodes["hook:before_sign:linux-amd64"]
	assert.Equal(t, []NodeKey{"hook:after_package:linux-amd64"}, beforeSign.Dependencies)
	for _, key := range []NodeKey{"package:linux/amd64:deb:sign", "package:linux/amd64:rpm:sign"} {
		assert.Contains(t, plan.Nodes[key].Dependencies, beforeSign.Key)
	}
	afterSign := plan.Nodes["hook:after_sign:linux-amd64"]
	assert.ElementsMatch(t, []NodeKey{"package:linux/amd64:deb:sign", "package:linux/amd64:rpm:sign"}, afterSign.Dependencies)
	for _, key := range []NodeKey{"publish:linux/amd64:deb:sign", "publish:linux/amd64:rpm:sign"} {
		assert.Contains(t, plan.Nodes[key].Dependencies, afterSign.Key)
	}

	multi, err := PlanBuild(config, Request{Verb: "build", Targets: []Target{{OS: "linux", Arch: "amd64"}, {OS: "linux", Arch: "arm64"}}})
	require.NoError(t, err)
	count := 0
	for key := range multi.Nodes {
		if key == "hook:before_build" {
			count++
		}
	}
	assert.Equal(t, 1, count)
	assert.Contains(t, multi.Nodes, NodeKey("hook:after_build:linux-amd64"))
	assert.Contains(t, multi.Nodes, NodeKey("hook:after_build:linux-arm64"))
}

func TestHookPlannerCacheContract(t *testing.T) {
	config := testConfig(t)
	config.Hooks = map[manifest.HookPhase]manifest.Hook{
		manifest.BeforeBuild: {Script: "scripts/hook.sh", Cache: true, Inputs: []string{"version.txt"}, Outputs: []string{"generated/a", "generated/b"}},
	}
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	node := plan.Nodes["hook:before_build"]
	assert.Equal(t, CacheArtifact, node.Cache)
	assert.Equal(t, "generated", node.Output)
	assert.Equal(t, []InputSpec{{Label: "hook-script", Files: []string{"scripts/hook.sh"}}, {Label: "hook-inputs", Files: []string{"version.txt"}}}, node.Inputs)
}

func TestAlwaysRunHooksMakeOnlyTheirDownstreamOperationsUncacheable(t *testing.T) {
	config := testConfig(t)
	config.Hooks = map[manifest.HookPhase]manifest.Hook{
		manifest.AfterBuild: {Script: "scripts/after.sh"},
	}
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	assert.Equal(t, CacheArtifact, plan.Nodes["target:linux/amd64:compile"].Cache)
	assert.Equal(t, CacheNever, plan.Nodes["hook:after_build:linux-amd64"].Cache)
	assert.Equal(t, CacheNever, plan.Nodes["publish:target:linux/amd64:compile"].Cache)
	assert.Equal(t, CacheNever, plan.Nodes["collect:artifacts"].Cache)

	config = testConfig(t)
	config.Hooks = map[manifest.HookPhase]manifest.Hook{
		manifest.BeforeBuild: {Script: "scripts/before.sh"},
	}
	plan, err = PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	assert.Equal(t, CacheNever, plan.Nodes["frontend:install"].Cache)
	assert.Equal(t, CacheNever, plan.Nodes["frontend:build"].Cache)
	assert.Equal(t, CacheNever, plan.Nodes["target:linux/amd64:compile"].Cache)
}

func TestHookPlannerFingerprintsResolvedInProjectSymlinkTargets(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating symlinks requires elevated privileges on some Windows hosts")
	}
	config := testConfig(t)
	require.NoError(t, os.MkdirAll(filepath.Join(config.Root, "scripts"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(config.Root, "scripts", "hook.sh"), []byte("#!/bin/sh\n"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(config.Root, "version.txt"), []byte("one\n"), 0o644))
	require.NoError(t, os.Symlink("version.txt", filepath.Join(config.Root, "version-link.txt")))
	config.Hooks = map[manifest.HookPhase]manifest.Hook{
		manifest.BeforeBuild: {Script: "scripts/hook.sh", Cache: true, Inputs: []string{"version-link.txt"}, Outputs: []string{"generated/out"}},
	}
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	node := plan.Nodes["hook:before_build"]
	require.Len(t, node.Inputs, 2)
	assert.Equal(t, filepath.Join(config.Root, "version.txt"), node.Inputs[1].Files[0])
}

func allPlannerHooks() map[manifest.HookPhase]manifest.Hook {
	result := make(map[manifest.HookPhase]manifest.Hook)
	for _, phase := range manifest.HookPhases {
		result[phase] = manifest.Hook{Script: "scripts/" + string(phase) + ".sh"}
	}
	return result
}
