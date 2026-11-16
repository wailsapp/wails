package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func planPrebuilt(plan Plan, config manifest.Config, request Request, multiTarget bool) (Plan, error) {
	if request.Development || (request.TargetOS != "linux" && request.TargetOS != "windows" && request.TargetOS != "darwin") {
		return Plan{}, fmt.Errorf("--skip-build supports production desktop binaries only")
	}
	directory := config.Build.OutputDirectory
	if multiTarget {
		directory = filepath.Join(directory, request.TargetOS+"-"+request.TargetArch)
	}
	for key, node := range plan.Nodes {
		if node.Kind != CompileApplication && node.Kind != MergeUniversalBinaries {
			continue
		}
		name := filepath.Base(node.Output)
		source := filepath.Join(directory, name)
		if request.TargetOS == "darwin" && name == config.Project.BinaryName && config.Project.Kind != "cli" && config.Project.Kind != "server" {
			if _, err := os.Lstat(filepath.Join(config.Root, source)); os.IsNotExist(err) {
				source = filepath.Join(directory, name+".app", "Contents", "MacOS", name)
			}
		}
		source = filepath.ToSlash(source)
		node.Kind, node.Label = ImportBinary, "Import prebuilt "+name
		node.Dependencies = nil
		node.Spec = PrebuiltSpec{Source: source, Output: node.Output, Target: node.Artifact.Target}
		node.Inputs = []InputSpec{{Label: "prebuilt-binary", Files: []string{source}}}
		node.Cache, node.Claims, node.EstimateMS = CacheArtifact, ResourceClaims{CPU: 1, MemoryMB: 64}, 25
		plan.Nodes[key] = node
	}
	// Imported executables are inputs, not publication destinations. Keep them
	// untouched when running or assembling an application from a prebuilt binary.
	for key, node := range plan.Nodes {
		if node.Kind != PublishArtifact || node.Artifact.Signed || (node.Artifact.Kind != ArtifactBinary && node.Artifact.Format != "app") {
			continue
		}
		source := node.Dependencies[0]
		for i, root := range plan.Roots {
			if root == key {
				plan.Roots[i] = source
			}
		}
		for i, artifact := range plan.Artifacts {
			if artifact == key {
				plan.Artifacts[i] = source
			}
		}
	}
	// A bare build/run still needs its companions beside the primary executable.
	// Universal component nodes are replaced by the imported merged executable.
	for key, node := range plan.Nodes {
		if node.Kind == ImportBinary && node.Artifact.Target == (Target{OS: request.TargetOS, Arch: request.TargetArch}) {
			plan.Roots = appendUniqueKeys(plan.Roots, key)
		}
	}
	slices.Sort(plan.Roots)
	reachable := map[NodeKey]bool{}
	var visit func(NodeKey)
	visit = func(key NodeKey) {
		if reachable[key] {
			return
		}
		reachable[key] = true
		for _, dependency := range plan.Nodes[key].Dependencies {
			visit(dependency)
		}
	}
	for _, key := range plan.Roots {
		visit(key)
	}
	for key := range plan.Nodes {
		if !reachable[key] {
			delete(plan.Nodes, key)
		}
	}
	if err := publishAppImageSidecars(&plan); err != nil {
		return Plan{}, err
	}
	return plan, plan.Validate(config.Root)
}
