package pipeline

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/internal/buildmeta"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func resolveBuildConfig(config manifest.Config, request Request) (manifest.Config, error) {
	selected, err := selectExecutables(config, request.Binaries)
	if err != nil {
		return config, err
	}
	if request.SkipFrontend {
		selected.Frontend.Disabled = true
	}
	if request.SkipBuild {
		selected.Frontend.Disabled = true
		selected.Hooks = make(map[manifest.HookPhase]manifest.Hook, len(config.Hooks))
		for phase, hook := range config.Hooks {
			if phase != manifest.BeforeBuild && phase != manifest.AfterBuild {
				selected.Hooks[phase] = hook
			}
		}
		selected.Project.CommitSymbol = ""
		selected.Binaries = append([]manifest.Binary(nil), selected.Binaries...)
		for i := range selected.Binaries {
			selected.Binaries[i].CommitSymbol = ""
		}
	}
	return resolveProjectMetadata(selected)
}

func resolveProjectMetadata(config manifest.Config) (manifest.Config, error) {
	if config.Project.VersionFile != "" {
		path, err := manifest.ResolveProjectPath(config.Root, "version source", config.Project.VersionFile, true)
		if err != nil {
			return config, err
		}
		info, err := buildmeta.Source(path, config.Project.VersionSymbol)
		if err != nil {
			return config, manifest.AnnotateValidationError(err, config.Origins, "project.version_file")
		}
		config.Project.Version = info.Version
	}
	needsCommit := config.Project.CommitSymbol != ""
	for _, binary := range config.Binaries {
		needsCommit = needsCommit || binary.CommitSymbol != ""
	}
	if !needsCommit {
		return config, nil
	}
	revision, err := buildmeta.Revision(context.Background(), config.Root)
	if err != nil {
		return config, fmt.Errorf("resolve application commit injection: %w", err)
	}
	config.Project.Revision = revision
	config.Binaries = append([]manifest.Binary(nil), config.Binaries...)
	for i, binary := range config.Binaries {
		if binary.CommitSymbol != "" {
			config.Binaries[i].LinkerFlags = append(append([]string(nil), binary.LinkerFlags...), "-X", binary.CommitSymbol+"="+revision)
		}
	}
	return config, nil
}

func planCompanions(plan *Plan, config manifest.Config, request Request, targets []Target, generated string) ([]NodeKey, []NamedBinary, error) {
	if len(config.Binaries) == 0 {
		return nil, nil, nil
	}
	if request.TargetOS == "ios" || request.TargetOS == "android" {
		return nil, nil, fmt.Errorf("additional executables are supported for desktop targets only")
	}
	var roots []NodeKey
	var binaries []NamedBinary
	for _, binary := range config.Binaries {
		kind := binary.Kind
		if kind == "" {
			kind = "cli"
		}
		name := binary.Name
		if request.TargetOS == "windows" {
			name += ".exe"
		}
		output := filepath.ToSlash(filepath.Join(generated, "artifacts", name))
		var components []NodeKey
		var paths []string
		for _, target := range targets {
			base := plan.Nodes[NodeKey("target:"+target.OS+"/"+target.Arch+":compile")]
			spec := base.Spec.(CompileSpec)
			spec.Package, spec.Console = binary.Package, kind != "gui"
			spec.Output = output
			if len(targets) > 1 {
				spec.Output = filepath.ToSlash(filepath.Join(generated, "companions", target.Arch, name))
			}
			spec.Tags = appendUnique(append([]string(nil), spec.Tags...), binary.Tags...)
			// Remove only the primary's injected metadata pair before adding this
			// executable's own flags; both programs may use the same Go symbol name.
			spec.LinkerFlags = append([]string(nil), spec.LinkerFlags...)
			if config.Project.CommitSymbol != "" && len(spec.LinkerFlags) >= 2 {
				n := len(spec.LinkerFlags)
				if spec.LinkerFlags[n-2] == "-X" && spec.LinkerFlags[n-1] == config.Project.CommitSymbol+"="+config.Project.Revision {
					spec.LinkerFlags = spec.LinkerFlags[:n-2]
				}
			}
			spec.LinkerFlags = append(spec.LinkerFlags, binary.LinkerFlags...)
			spec.CompilerFlags = append(append([]string(nil), spec.CompilerFlags...), binary.CompilerFlags...)
			spec.Environment = cloneStringMap(base.Spec.(CompileSpec).Environment)
			if spec.Environment == nil {
				spec.Environment = map[string]string{}
			}
			for key, value := range binary.Environment {
				spec.Environment[key] = value
			}
			spec.UPXLevel = binary.UPXLevel
			if request.Development || target.OS == "darwin" {
				spec.UPXLevel = 0
			}
			key := NodeKey("binary:" + request.TargetOS + "/" + request.TargetArch + ":" + binary.Name + ":" + target.Arch)
			node := Node{Key: key, Kind: CompileApplication, Label: "Compile " + binary.Name + " (" + target.OS + "/" + target.Arch + ")", Scope: TargetScope, Dependencies: append([]NodeKey(nil), base.Dependencies...), Spec: spec, Inputs: append([]InputSpec(nil), base.Inputs...), Output: spec.Output, Cache: CacheArtifact, Claims: base.Claims, EstimateMS: base.EstimateMS, Artifact: ArtifactIdentity{Kind: ArtifactBinary, Target: target}}
			taintCacheAfterAlwaysRunHook(*plan, &node)
			addNode(plan, node, nil)
			components = append(components, key)
			paths = append(paths, spec.Output)
		}
		root := components[0]
		if len(components) > 1 {
			root = NodeKey("binary:" + request.TargetOS + "/" + request.TargetArch + ":" + binary.Name + ":merge")
			addNode(plan, Node{Key: root, Kind: MergeUniversalBinaries, Label: "Merge " + binary.Name, Scope: TargetScope, Dependencies: components, Spec: MergeSpec{Inputs: paths, Output: output}, Inputs: []InputSpec{{Label: "companion-components", Files: paths}}, Output: output, Cache: CacheArtifact, Claims: ResourceClaims{CPU: 1, MemoryMB: 256}, Artifact: ArtifactIdentity{Kind: ArtifactBinary, Target: Target{OS: request.TargetOS, Arch: request.TargetArch}}}, nil)
		}
		roots = append(roots, root)
		binaries = append(binaries, NamedBinary{Name: name, Path: output, Kind: kind})
	}
	return roots, binaries, nil
}

func publishCompanions(plan *Plan, config manifest.Config, request Request, multiTarget bool) error {
	if len(config.Binaries) == 0 || request.Development {
		return nil
	}
	for _, binary := range config.Binaries {
		prefix := "binary:" + request.TargetOS + "/" + request.TargetArch + ":" + binary.Name + ":"
		var source Node
		for key, node := range plan.Nodes {
			if strings.HasPrefix(string(key), prefix) && (strings.HasSuffix(string(key), ":merge") || source.Key == "") {
				source = node
			}
		}
		if source.Key == "" {
			return fmt.Errorf("missing planned executable %s", binary.Name)
		}
		name := binary.Name
		if request.TargetOS == "windows" {
			name += ".exe"
		}
		directory := config.Build.OutputDirectory
		if multiTarget {
			directory = filepath.Join(directory, request.TargetOS+"-"+request.TargetArch)
		}
		output := filepath.ToSlash(filepath.Join(directory, name))
		key := NodeKey("publish:" + prefix + "executable")
		node := Node{Key: key, Kind: PublishArtifact, Label: "Publish " + name, Scope: TargetScope, Dependencies: appendUniqueKeys([]NodeKey{source.Key}, plan.Roots...), Spec: PublishSpec{Source: source.Output, Destination: output}, Output: output, Cache: source.Cache, Claims: ResourceClaims{CPU: 1, MemoryMB: 64}, Artifact: source.Artifact}
		taintCacheAfterAlwaysRunHook(*plan, &node)
		addNode(plan, node, nil)
		plan.Roots = append(plan.Roots, key)
		plan.Artifacts = append(plan.Artifacts, key)
	}
	return nil
}

func executableUPXLevel(level int, development bool, targetOS string) int {
	if development || (targetOS != "linux" && targetOS != "windows") {
		return 0
	}
	return level
}

func finalizeTargetPlan(plan Plan, config manifest.Config, request Request, multiTarget bool) (Plan, error) {
	if request.SkipBuild {
		return planPrebuilt(plan, config, request, multiTarget)
	}
	if err := resolveCompileModes(&plan, request); err != nil {
		return Plan{}, err
	}
	refineCompileInputs(&plan, config.Root)
	if err := publishAppImageSidecars(&plan); err != nil {
		return Plan{}, err
	}
	if err := publishCompanions(&plan, config, request, multiTarget); err != nil {
		return Plan{}, err
	}
	return plan, plan.Validate(config.Root)
}

// Selection promotes a requested executable to the primary slot and retains
// only requested companions. Frontend work is shared whenever a GUI is selected.
func selectExecutables(config manifest.Config, names []string) (manifest.Config, error) {
	if len(names) == 0 {
		return config, nil
	}
	primary := manifest.Binary{Name: config.Project.BinaryName, Package: config.Project.Package, Kind: config.Project.Kind, CommitSymbol: config.Project.CommitSymbol, UPXLevel: config.Project.UPXLevel}
	if primary.Kind == "" {
		primary.Kind = "gui"
	}
	available := map[string]manifest.Binary{primary.Name: primary}
	for _, binary := range config.Binaries {
		if binary.Kind == "" {
			binary.Kind = "cli"
		}
		available[binary.Name] = binary
	}
	var selected []manifest.Binary
	seen := map[string]bool{}
	for _, name := range names {
		binary, ok := available[name]
		if !ok {
			return config, fmt.Errorf("unknown executable %q", name)
		}
		if seen[name] {
			return config, fmt.Errorf("duplicate requested executable %q", name)
		}
		seen[name] = true
		selected = append(selected, binary)
	}
	// Package a GUI as the application's primary executable, independent of flag order.
	for i, binary := range selected {
		if binary.Kind == "gui" {
			selected[0], selected[i] = selected[i], selected[0]
			break
		}
	}
	first := selected[0]
	config.Project.BinaryName, config.Project.Package, config.Project.Kind = first.Name, first.Package, first.Kind
	config.Project.CommitSymbol, config.Project.UPXLevel = first.CommitSymbol, first.UPXLevel
	config.PrimaryExecutable = &first
	config.Binaries = append([]manifest.Binary(nil), selected[1:]...)
	if first.Kind != "gui" {
		config.Frontend.Disabled = true
	}
	return config, nil
}

func primaryBuildTags(config manifest.Config, tags []string) []string {
	if config.PrimaryExecutable == nil {
		return tags
	}
	return appendUnique(append([]string(nil), tags...), config.PrimaryExecutable.Tags...)
}

func applyPrimaryExecutable(plan *Plan, config manifest.Config, targets []Target) {
	selected := config.PrimaryExecutable
	if selected == nil {
		return
	}
	for _, target := range targets {
		key := NodeKey("target:" + target.OS + "/" + target.Arch + ":compile")
		node := plan.Nodes[key]
		spec := node.Spec.(CompileSpec)
		spec.Tags = primaryBuildTags(config, spec.Tags)
		spec.LinkerFlags = append(append([]string(nil), spec.LinkerFlags...), selected.LinkerFlags...)
		spec.CompilerFlags = append(append([]string(nil), spec.CompilerFlags...), selected.CompilerFlags...)
		spec.Environment = cloneStringMap(spec.Environment)
		if spec.Environment == nil {
			spec.Environment = map[string]string{}
		}
		for key, value := range selected.Environment {
			spec.Environment[key] = value
		}
		node.Spec = spec
		plan.Nodes[key] = node
	}
}
