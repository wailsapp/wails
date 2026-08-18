package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"golang.org/x/mod/modfile"
)

func PlanBuild(config manifest.Config, request Request) (Plan, error) {
	if config.Selected.Name != "" && len(request.Targets) == 0 && request.TargetOS == "" && request.TargetArch == "" {
		return planSelectedProfile(config, request)
	}
	targets := append([]Target(nil), request.Targets...)
	if len(targets) == 0 {
		targets = []Target{{OS: request.TargetOS, Arch: request.TargetArch}}
	}
	seen := map[string]bool{}
	for i := range targets {
		if targets[i].OS == "" {
			targets[i].OS = runtime.GOOS
		}
		if targets[i].Arch == "" {
			targets[i].Arch = runtime.GOARCH
		}
		key := targets[i].OS + "/" + targets[i].Arch
		if seen[key] {
			return Plan{}, fmt.Errorf("duplicate target %s", key)
		}
		seen[key] = true
	}
	sort.Slice(targets, func(i, j int) bool {
		return targets[i].OS+"/"+targets[i].Arch < targets[j].OS+"/"+targets[j].Arch
	})
	if len(targets) == 1 {
		request.TargetOS, request.TargetArch = targets[0].OS, targets[0].Arch
		request.Targets = nil
		return planTarget(config, request, false)
	}
	combined := Plan{Name: request.Verb, Nodes: map[NodeKey]Node{}}
	for _, target := range targets {
		childRequest := request
		childRequest.TargetOS, childRequest.TargetArch, childRequest.Targets = target.OS, target.Arch, nil
		child, err := planTarget(config, childRequest, true)
		if err != nil {
			return Plan{}, err
		}
		for key, node := range child.Nodes {
			if existing, ok := combined.Nodes[key]; ok {
				if !reflect.DeepEqual(existing, node) {
					return Plan{}, fmt.Errorf("project node %s resolves differently across targets; use equivalent target settings or separate invocations", key)
				}
				continue
			}
			combined.Nodes[key] = node
		}
		combined.Roots = appendUniqueKeys(combined.Roots, child.Roots...)
	}
	names := make([]string, len(targets))
	for i, target := range targets {
		names[i] = target.OS + "/" + target.Arch
	}
	combined.Target = strings.Join(names, ",")
	if combined.Name == "" {
		combined.Name = "build"
	}
	return combined, combined.Validate(config.Root)
}

// planSelectedProfile expands the complete profile request into one child plan
// per concrete target. The fixed pipeline remains private to Wails; profiles
// select artifacts, not stages or dependencies.
func planSelectedProfile(config manifest.Config, request Request) (Plan, error) {
	combined := Plan{Name: request.Verb, Nodes: map[NodeKey]Node{}}
	for _, selected := range config.Selected.Targets {
		goos, arch, err := splitProfileTarget(selected.Target)
		if err != nil {
			return Plan{}, err
		}
		childConfig := config
		if selected.Destination != "" {
			if err := setTargetVariant(&childConfig.Targets, goos, arch, selected.Destination); err != nil {
				return Plan{}, err
			}
		}
		if selected.Notarize && goos == "darwin" {
			childConfig.Signing.Darwin.Notarize = true
		}
		childRequest := request
		childRequest.TargetOS, childRequest.TargetArch = goos, arch
		childRequest.Targets = nil
		childRequest.Formats = append([]string(nil), selected.Formats...)
		if childRequest.Verb == "build" {
			switch {
			case selected.Sign:
				childRequest.Verb = "sign"
			case len(selected.Formats) > 0:
				childRequest.Verb = "package"
			}
		}
		child, err := planTarget(childConfig, childRequest, true)
		if err != nil {
			return Plan{}, err
		}
		if err := addPlanNodes(&combined, child); err != nil {
			return Plan{}, err
		}
	}
	if combined.Name == "" {
		combined.Name = "build"
	}
	return combined, combined.Validate(config.Root)
}

func addPlanNodes(combined *Plan, child Plan) error {
	for key, node := range child.Nodes {
		if existing, ok := combined.Nodes[key]; ok {
			if !reflect.DeepEqual(existing, node) {
				return fmt.Errorf("project node %s resolves differently across targets; use equivalent target settings or separate invocations", key)
			}
			continue
		}
		combined.Nodes[key] = node
	}
	combined.Roots = appendUniqueKeys(combined.Roots, child.Roots...)
	if combined.Target == "" {
		combined.Target = child.Target
	} else {
		combined.Target += "," + child.Target
	}
	return nil
}

func splitProfileTarget(value string) (string, string, error) {
	parts := strings.Split(value, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid profile target %q", value)
	}
	return parts[0], parts[1], nil
}

func setTargetVariant(targets *manifest.Targets, platform, arch, variant string) error {
	var configured *manifest.Platform
	switch platform {
	case "windows":
		configured = &targets.Windows
	case "darwin":
		configured = &targets.Darwin
	case "linux":
		configured = &targets.Linux
	case "ios":
		configured = &targets.IOS
	case "android":
		configured = &targets.Android
	default:
		return fmt.Errorf("unsupported profile target platform %q", platform)
	}
	switch arch {
	case "amd64":
		configured.AMD64.Variant = variant
	case "arm64":
		configured.ARM64.Variant = variant
	case "arm":
		configured.ARM.Variant = variant
	case "386":
		configured.X86.Variant = variant
	case "universal":
		configured.Universal.Variant = variant
	default:
		return fmt.Errorf("unsupported profile target architecture %q", arch)
	}
	return nil
}

func planTarget(config manifest.Config, request Request, multiTarget bool) (Plan, error) {
	if request.Verb == "" {
		request.Verb = "build"
	}
	if request.TargetOS == "" {
		request.TargetOS = runtime.GOOS
	}
	if request.TargetArch == "" {
		request.TargetArch = runtime.GOARCH
	}
	if !containsPlatform(request.TargetOS) {
		return Plan{}, fmt.Errorf("unsupported target platform %q", request.TargetOS)
	}
	if !containsArchitecture(request.TargetArch) {
		return Plan{}, fmt.Errorf("unsupported target architecture %q", request.TargetArch)
	}
	if request.TargetArch == "universal" && request.TargetOS != "darwin" {
		return Plan{}, fmt.Errorf("universal target is only valid for darwin")
	}

	profileName := config.Profile
	if profileName == "" {
		profileName = "default"
	}
	target := request.TargetOS + "/" + request.TargetArch
	project, targetSettings := effectiveTarget(config, request.TargetOS, request.TargetArch)
	platformSettings := platformConfig(config.Targets, request.TargetOS)
	associations := associationsForPlatform(config.Associations, request.TargetOS)
	protocols := protocolsForPlatform(config.Protocols, request.TargetOS)
	assetsOut := filepath.ToSlash(filepath.Join(".wails", "build", profileName, strings.ReplaceAll(target, "/", "-"), "assets"))
	binaryName := config.Project.BinaryName
	if request.TargetOS == "windows" {
		binaryName += ".exe"
	}
	binaryDirectory := config.Build.OutputDirectory
	if multiTarget {
		binaryDirectory = filepath.Join(binaryDirectory, strings.ReplaceAll(target, "/", "-"))
	}
	binaryOut := filepath.ToSlash(filepath.Join(binaryDirectory, binaryName))
	if request.TargetOS == "ios" {
		binaryOut += ".a"
	}
	if request.TargetOS == "android" {
		binaryOut = filepath.ToSlash(filepath.Join(".wails", "build", profileName, strings.ReplaceAll(target, "/", "-"), "libwails.so"))
	}
	plan := Plan{Name: request.Verb, Target: target, Nodes: map[NodeKey]Node{}}
	add := func(node Node) NodeKey {
		if _, exists := plan.Nodes[node.Key]; exists {
			panic("duplicate planner key: " + string(node.Key))
		}
		plan.Nodes[node.Key] = node
		return node.Key
	}

	var beforeBuild NodeKey
	if hook := config.Hooks.BeforeBuild; hook.Script != "" {
		hookRequest := request
		// before_build attaches to shared project preparation. It has no Target
		// or phase Artifact, regardless of whether one or many Targets were
		// requested, so the Node remains structurally shareable.
		hookRequest.TargetOS, hookRequest.TargetArch = "", ""
		node, hookErr := hookNode("before_build", hook, config, hookRequest, ProjectScope, "", nil)
		if hookErr != nil {
			return Plan{}, hookErr
		}
		beforeBuild = add(node)
	}
	install := add(Node{Key: "frontend:install", Kind: InstallFrontendDependencies, Label: "Install frontend dependencies", Scope: ProjectScope,
		Spec:   InstallSpec{Manager: config.Frontend.PackageManager, Directory: config.Frontend.Directory, Command: config.Frontend.InstallCommand, Arguments: append([]string(nil), config.Frontend.Install...)},
		Inputs: []InputSpec{{Label: "frontend-install", Files: frontendInstallFiles(config)}}, Marker: filepath.ToSlash(filepath.Join(config.Frontend.Directory, "node_modules")), Cache: CacheReceipt, Claims: ResourceClaims{CPU: 1, MemoryMB: 512}, EstimateMS: 1800})
	tags := append([]string(nil), config.Build.Go.Tags...)
	tags = appendUnique(tags, targetSettings.Tags...)
	if request.Development {
		tags = appendUnique(tags, request.ExtraTags...)
	} else {
		tags = appendUnique(tags, "production")
		tags = appendUnique(tags, request.ExtraTags...)
	}
	if request.Obfuscated || config.Build.Obfuscation {
		tags = appendUnique(tags, "wails_obfuscated")
	}
	localInputs, err := goLocalSourceInputs(config.Root)
	if err != nil {
		return Plan{}, err
	}
	bindingsDeps := keys(beforeBuild)
	bindingsOut := filepath.ToSlash(filepath.Join(config.Frontend.Directory, config.Frontend.Bindings.OutputDirectory))
	bindingInputs := []InputSpec{
		{Label: "go-binding-api", Root: ".", SemanticGo: true, ExcludeDirs: []string{".git", ".wails", "bin", "build", "dist", config.Frontend.Directory, "node_modules"}},
		{Label: "go-module", Files: goMetadataFiles(config.Root)},
	}
	for _, input := range localInputs {
		bindingInputs = append(bindingInputs, InputSpec{Label: "go-binding-local-api", Root: input.Root, SemanticGo: true, ExcludeDirs: input.ExcludeDirs})
	}
	bindings := add(Node{Key: "frontend:bindings", Kind: GenerateBindings, Label: "Generate bindings", Scope: ProjectScope, Dependencies: bindingsDeps,
		Spec:   BindingsSpec{Config: config.Frontend.Bindings, Tags: tags, Obfuscated: request.Obfuscated || config.Build.Obfuscation},
		Inputs: bindingInputs,
		Output: bindingsOut, Cache: CacheArtifact, Claims: ResourceClaims{CPU: 1, MemoryMB: 512, Exclusive: "legacy-command-adapter"}, EstimateMS: 1000})
	frontendDeps := []NodeKey{install, bindings}
	if beforeBuild != "" {
		frontendDeps = append(frontendDeps, beforeBuild)
	}
	frontendOut := frontendOutputPath(config.Frontend.Directory, config.Frontend.OutputDirectory)
	frontendExclude := filepath.ToSlash(filepath.Clean(config.Frontend.OutputDirectory))
	if frontendExclude == config.Frontend.Directory || strings.HasPrefix(frontendExclude, config.Frontend.Directory+"/") {
		frontendExclude = strings.TrimPrefix(strings.TrimPrefix(frontendExclude, config.Frontend.Directory), "/")
	}
	frontend := add(Node{Key: "frontend:build", Kind: BuildFrontend, Label: "Build frontend", Scope: ProjectScope, Dependencies: frontendDeps,
		Spec:   FrontendSpec{Manager: config.Frontend.PackageManager, Directory: config.Frontend.Directory, Command: choose(request.Development, config.Frontend.BuildCommand+":dev", config.Frontend.BuildCommand), Arguments: append([]string(nil), config.Frontend.Build...), Output: config.Frontend.OutputDirectory, Production: !request.Development},
		Inputs: []InputSpec{{Label: "frontend-source", Root: config.Frontend.Directory, IncludeAll: true, ExcludeDirs: []string{".git", ".wails", "node_modules", config.Frontend.Bindings.OutputDirectory, frontendExclude}}},
		Output: frontendOut, Cache: CacheArtifact, Claims: ResourceClaims{CPU: 2, MemoryMB: 1024}, EstimateMS: 900})
	var assets NodeKey
	if request.TargetOS == "windows" || request.TargetOS == "ios" || request.TargetOS == "android" {
		assets = add(platformAssetsNode(config, request, target, assetsOut, keys(beforeBuild)))
	}

	compileDeps := []NodeKey{frontend}
	if beforeBuild != "" {
		compileDeps = append(compileDeps, beforeBuild)
	}
	if assets != "" {
		compileDeps = append(compileDeps, assets)
	}
	compileInputs := []InputSpec{
		{Label: "go-sources", Root: ".", IncludeNames: []string{"go.mod", "go.sum", "go.work", "go.work.sum"}, IncludeExtensions: []string{".go", ".c", ".cc", ".cpp", ".h", ".hh", ".hpp", ".s", ".syso"}, ExcludeDirs: []string{".git", ".wails", "bin", "build", "dist", config.Frontend.Directory, "node_modules"}},
		{Label: "go-module", Files: goMetadataFiles(config.Root)},
	}
	compileInputs = append(compileInputs, localInputs...)
	compile := add(Node{Key: NodeKey("target:" + target + ":compile"), Kind: CompileApplication, Label: "Compile " + target, Scope: TargetScope, Dependencies: compileDeps,
		Spec:   CompileSpec{TargetOS: request.TargetOS, TargetArch: request.TargetArch, Output: binaryOut, Assets: assetsOut, Variant: targetSettings.Variant, MinimumVersion: targetSettings.MinimumVersion, Tags: tags, LinkerFlags: config.Build.Go.LinkerFlags, CompilerFlags: config.Build.Go.CompilerFlags, GarbleArgs: config.Build.Go.GarbleArgs, Production: !request.Development, Obfuscated: request.Obfuscated || config.Build.Obfuscation, TrimPath: config.Build.TrimPath, Strip: config.Build.Strip},
		Inputs: compileInputs,
		Output: binaryOut, Cache: CacheArtifact, Claims: ResourceClaims{CPU: max(1, runtime.GOMAXPROCS(0)/2), MemoryMB: 2048}, EstimateMS: 1500, ArtifactKind: "binary"})
	lastBuild := compile
	if hook := config.Hooks.AfterBuild; hook.Script != "" {
		node, hookErr := hookNode("after_build", hook, config, request, TargetScope, binaryOut, []NodeKey{compile})
		if hookErr != nil {
			return Plan{}, hookErr
		}
		lastBuild = add(node)
	}

	if request.Verb == "build" {
		plan.Roots = []NodeKey{lastBuild}
		return plan, plan.Validate(config.Root)
	}
	if request.Verb != "package" && request.Verb != "sign" {
		return Plan{}, fmt.Errorf("unsupported pipeline verb %q", request.Verb)
	}
	beforePackageDeps := []NodeKey{lastBuild}
	formats := append([]string(nil), request.Formats...)
	if len(formats) == 0 {
		formats = append([]string(nil), packagePlatform(config.Package, request.TargetOS).Formats...)
	}
	if len(formats) == 0 {
		return Plan{}, fmt.Errorf("no package formats configured for %s", request.TargetOS)
	}
	sort.Strings(formats)
	packageOutputs := make([]string, 0, len(formats))
	for _, format := range formats {
		if !platformSupportsFormat(request.TargetOS, format) {
			return Plan{}, fmt.Errorf("package format %q is not supported for %s", format, request.TargetOS)
		}
		if request.TargetOS == "ios" && format == "ipa" && targetSettings.Variant != "device" {
			return Plan{}, fmt.Errorf("IPA packaging requires [targets.ios.%s] variant = \"device\"", request.TargetArch)
		}
		packageOutputs = append(packageOutputs, packageOutput(config, request.TargetOS, request.TargetArch, format, multiTarget))
	}
	packageArtifact := commonOutput(packageOutputs)
	if hook := config.Hooks.BeforePackage; hook.Script != "" {
		node, hookErr := hookNode("before_package", hook, config, request, TargetScope, packageArtifact, beforePackageDeps)
		if hookErr != nil {
			return Plan{}, hookErr
		}
		beforePackageDeps = []NodeKey{add(node)}
	}
	if assets == "" {
		assets = add(platformAssetsNode(config, request, target, assetsOut, beforePackageDeps))
	}
	var packageRoots []NodeKey
	for _, format := range formats {
		pkgConfig, err := packageFormat(packagePlatform(config.Package, request.TargetOS), format)
		if err != nil {
			return Plan{}, err
		}
		output := packageOutput(config, request.TargetOS, request.TargetArch, format, multiTarget)
		key := NodeKey("package:" + target + ":" + format)
		packageDeps := appendUniqueKeys([]NodeKey{compile, assets}, beforePackageDeps...)
		packageCache := CacheArtifact
		if request.TargetOS == "ios" {
			// iOS bundle assembly invokes codesign, including ad-hoc simulator
			// signing, so its result must never enter the reusable artifact cache.
			packageCache = CacheNever
		}
		templateRoot := ""
		if pkgConfig.Template != "" {
			templateRoot = config.Root
		}
		pkg := add(Node{Key: key, Kind: PackageArtifact, Label: "Package " + format, Scope: PackageScope, Dependencies: packageDeps,
			Spec:   PackageSpec{TargetOS: request.TargetOS, TargetArch: request.TargetArch, Format: format, Binary: binaryOut, Assets: assetsOut, Output: output, Profile: profileName, Variant: targetSettings.Variant, MinimumVersion: targetSettings.MinimumVersion, TemplateRoot: templateRoot, Config: pkgConfig, Project: project, Capabilities: platformSettings.Capabilities, Associations: associations, Protocols: protocols},
			Inputs: packageInputs(config.Root, request.TargetOS, format, pkgConfig), Output: output, Cache: packageCache, Claims: ResourceClaims{CPU: 1, MemoryMB: 1024, Exclusive: packageExclusive(request.TargetOS, format)}, EstimateMS: 1000, ArtifactKind: format})
		packageRoots = append(packageRoots, pkg)
	}
	packageArtifacts := append([]NodeKey(nil), packageRoots...)
	if hook := config.Hooks.AfterPackage; hook.Script != "" {
		node, hookErr := hookNode("after_package", hook, config, request, PackageScope, packageArtifact, packageRoots)
		if hookErr != nil {
			return Plan{}, hookErr
		}
		packageRoots = []NodeKey{add(node)}
	}
	if request.Verb == "sign" {
		if hook := config.Hooks.BeforeSign; hook.Script != "" {
			node, hookErr := hookNode("before_sign", hook, config, request, PackageScope, packageArtifact, packageRoots)
			if hookErr != nil {
				return Plan{}, hookErr
			}
			packageRoots = []NodeKey{add(node)}
		}
		var signed []NodeKey
		var signedOutputs []string
		for _, artifact := range packageArtifacts {
			input := plan.Nodes[artifact].Output
			dependencies := appendUniqueKeys([]NodeKey{artifact}, packageRoots...)
			signed = append(signed, add(Node{Key: NodeKey(string(artifact) + ":sign"), Kind: SignArtifact, Label: "Sign " + filepath.Base(input), Scope: PackageScope, Dependencies: dependencies, Spec: SignSpec{TargetOS: request.TargetOS, Format: plan.Nodes[artifact].ArtifactKind, Input: input, Config: signingPlatform(config.Signing, request.TargetOS)}, Output: input + ".signed", Cache: CacheNever, Claims: ResourceClaims{CPU: 1, MemoryMB: 256, Exclusive: "sign"}, EstimateMS: 1000, ArtifactKind: "signed"}))
			signedOutputs = append(signedOutputs, input+".signed")
		}
		packageRoots = signed
		if hook := config.Hooks.AfterSign; hook.Script != "" {
			node, hookErr := hookNode("after_sign", hook, config, request, PackageScope, commonOutput(signedOutputs), packageRoots)
			if hookErr != nil {
				return Plan{}, hookErr
			}
			packageRoots = []NodeKey{add(node)}
		}
	}
	plan.Roots = packageRoots
	return plan, plan.Validate(config.Root)
}

func platformAssetsNode(config manifest.Config, request Request, target, output string, dependencies []NodeKey) Node {
	project, targetSettings := effectiveTarget(config, request.TargetOS, request.TargetArch)
	platformSettings := platformConfig(config.Targets, request.TargetOS)
	return Node{Key: NodeKey("target:" + target + ":assets"), Kind: GeneratePlatformAssets, Label: "Generate " + request.TargetOS + " assets", Scope: TargetScope, Dependencies: dependencies,
		Spec:   AssetsSpec{TargetOS: request.TargetOS, TargetArch: request.TargetArch, Directory: output, MinimumVersion: targetSettings.MinimumVersion, Capabilities: platformSettings.Capabilities, Project: project, Associations: associationsForPlatform(config.Associations, request.TargetOS), Protocols: protocolsForPlatform(config.Protocols, request.TargetOS)},
		Inputs: assetInputs(config), Output: output, Cache: CacheArtifact, Claims: ResourceClaims{CPU: 1, MemoryMB: 256, Exclusive: "legacy-command-adapter"}, EstimateMS: 250}
}

func associationsForPlatform(values []manifest.Association, platform string) []manifest.Association {
	result := make([]manifest.Association, 0, len(values))
	for _, value := range values {
		if len(value.Platforms) == 0 || contains(value.Platforms, platform) {
			result = append(result, value)
		}
	}
	return result
}

func protocolsForPlatform(values []manifest.Protocol, platform string) []manifest.Protocol {
	result := make([]manifest.Protocol, 0, len(values))
	for _, value := range values {
		if len(value.Platforms) == 0 || contains(value.Platforms, platform) {
			result = append(result, value)
		}
	}
	return result
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func hookNode(phase string, hook manifest.Hook, config manifest.Config, request Request, scope Scope, artifact string, deps []NodeKey) (Node, error) {
	policy := CacheNever
	output := ""
	if hook.Cache {
		root, err := manifest.HookOutputRoot(hook.Outputs)
		if err != nil {
			return Node{}, fmt.Errorf("hook %s outputs: %w", phase, err)
		}
		policy, output = CacheArtifact, root
	}
	inputs := []InputSpec{{Label: "hook-script", Files: append([]string{hook.Script}, hook.Inputs...)}}
	scopeKey := request.TargetOS + "/" + request.TargetArch
	if request.TargetOS == "" && request.TargetArch == "" {
		scopeKey = "project"
	}
	profile := config.Profile
	if profile == "" {
		profile = "default"
	}
	return Node{Key: NodeKey("hook:" + phase + ":" + scopeKey), Kind: RunHook, Label: strings.ReplaceAll(phase, "_", " "), Scope: scope, Dependencies: deps, Spec: HookSpec{Phase: phase, TargetOS: request.TargetOS, TargetArch: request.TargetArch, Profile: profile, Artifact: artifact, Hook: hook}, Inputs: inputs, Output: output, Cache: policy, Claims: ResourceClaims{CPU: 1, MemoryMB: 512, Exclusive: "hook:" + phase}, EstimateMS: 100}, nil
}

func goLocalSourceInputs(root string) ([]InputSpec, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	local := map[string]bool{}
	add := func(base, value string) {
		if value == "" {
			return
		}
		if !filepath.IsAbs(value) {
			value = filepath.Join(base, value)
		}
		value, err = filepath.Abs(value)
		if err != nil || pathWithin(value, root) {
			return
		}
		if info, statErr := os.Stat(value); statErr == nil && info.IsDir() {
			local[filepath.Clean(value)] = true
		}
	}
	for _, metadata := range goMetadataFiles(root) {
		data, readErr := os.ReadFile(metadata)
		if readErr != nil {
			if os.IsNotExist(readErr) {
				continue
			}
			return nil, readErr
		}
		base := filepath.Dir(metadata)
		switch filepath.Base(metadata) {
		case "go.mod":
			file, parseErr := modfile.Parse(metadata, data, nil)
			if parseErr != nil {
				return nil, parseErr
			}
			for _, replacement := range file.Replace {
				if replacement.New.Version == "" {
					add(base, replacement.New.Path)
				}
			}
		case "go.work":
			file, parseErr := modfile.ParseWork(metadata, data, nil)
			if parseErr != nil {
				return nil, parseErr
			}
			for _, use := range file.Use {
				add(base, use.Path)
			}
			for _, replacement := range file.Replace {
				if replacement.New.Version == "" {
					add(base, replacement.New.Path)
				}
			}
		}
	}
	paths := make([]string, 0, len(local))
	for path := range local {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	inputs := make([]InputSpec, 0, len(paths))
	for _, path := range paths {
		inputs = append(inputs, InputSpec{Label: "go-local-source", Root: path, IncludeNames: []string{"go.mod", "go.sum"}, IncludeExtensions: []string{".go", ".c", ".cc", ".cpp", ".h", ".hh", ".hpp", ".s", ".syso"}, ExcludeDirs: []string{".git", ".wails", "bin", "dist", "node_modules"}})
	}
	return inputs, nil
}

func pathWithin(path, root string) bool {
	relative, err := filepath.Rel(root, path)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

// frontendOutputPath accepts both forms users naturally write in a single
// project manifest: a path relative to the frontend directory ("dist") and a
// project-relative path ("frontend/dist"). The planner stores one canonical
// project-relative output so cache snapshots and artifact reporting agree.
func frontendOutputPath(directory, output string) string {
	directory = filepath.ToSlash(filepath.Clean(directory))
	output = filepath.ToSlash(filepath.Clean(output))
	if output == directory || strings.HasPrefix(output, directory+"/") {
		return output
	}
	return filepath.ToSlash(filepath.Join(directory, output))
}

func frontendInstallFiles(c manifest.Config) []string {
	root := c.Frontend.Directory
	result := []string{filepath.Join(root, "package.json"), filepath.Join(root, "package-lock.json"), filepath.Join(root, "npm-shrinkwrap.json"), filepath.Join(root, "pnpm-lock.yaml"), filepath.Join(root, "yarn.lock"), filepath.Join(root, "bun.lock"), filepath.Join(root, "bun.lockb"), filepath.Join(root, ".npmrc")}
	return result
}
func goMetadataFiles(root string) []string {
	root, err := filepath.Abs(root)
	if err != nil {
		return []string{"go.mod", "go.sum", "go.work", "go.work.sum"}
	}
	var result []string
	foundModule := false
	foundWorkspace := false
	for dir := root; ; dir = filepath.Dir(dir) {
		if !foundModule {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				result = append(result, filepath.Join(dir, "go.mod"), filepath.Join(dir, "go.sum"))
				foundModule = true
			}
		}
		if !foundWorkspace {
			if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
				result = append(result, filepath.Join(dir, "go.work"), filepath.Join(dir, "go.work.sum"))
				foundWorkspace = true
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}
	if len(result) == 0 {
		return []string{"go.mod", "go.sum", "go.work", "go.work.sum"}
	}
	return result
}
func assetInputs(c manifest.Config) []InputSpec {
	files := []string{manifest.Filename}
	if c.Project.Icon != "" {
		files = append(files, c.Project.Icon)
	}
	for _, platform := range []manifest.Platform{c.Targets.Windows, c.Targets.Darwin, c.Targets.Linux, c.Targets.IOS, c.Targets.Android} {
		for _, path := range []string{platform.Icon, platform.Manifest, platform.AssetsCar, platform.InfoPlist} {
			if path != "" {
				files = append(files, path)
			}
		}
	}
	for _, signing := range []manifest.SigningPlatform{c.Signing.Windows, c.Signing.Darwin, c.Signing.Linux, c.Signing.IOS, c.Signing.Android} {
		for _, path := range []string{signing.Certificate, signing.Entitlements, signing.ProvisioningProfile} {
			if path != "" {
				files = append(files, path)
			}
		}
	}
	return []InputSpec{{Label: "platform-assets", Files: files}}
}
func packageInputs(root, platform, format string, f manifest.PackageFormat) []InputSpec {
	var result []InputSpec
	if f.Template != "" {
		path := f.Template
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, filepath.FromSlash(path))
		}
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			result = append(result, InputSpec{Label: "package-template", Root: f.Template, IncludeAll: true})
		} else {
			result = append(result, InputSpec{Label: "package-template", Files: []string{f.Template}})
		}
	}
	if platform == "darwin" && format == "dmg" && f.Template == "" {
		var resources []string
		for _, key := range []string{"background", "volume_icon", "file_icon"} {
			if path, ok := f.Options[key].(string); ok && path != "" {
				resources = append(resources, path)
			}
		}
		if files, ok := f.Options["files"].(string); ok {
			for _, item := range strings.Split(files, ",") {
				if _, path, found := strings.Cut(item, "="); found && strings.TrimSpace(path) != "" {
					resources = append(resources, strings.TrimSpace(path))
				}
			}
		}
		if len(resources) > 0 {
			result = append(result, InputSpec{Label: "package-resources", Files: resources})
		}
	}
	return result
}
func commonOutput(paths []string) string {
	if len(paths) == 1 {
		return filepath.ToSlash(paths[0])
	}
	if len(paths) > 1 {
		return filepath.ToSlash(filepath.Dir(paths[0]))
	}
	return ""
}
func keys(values ...NodeKey) []NodeKey {
	var result []NodeKey
	for _, v := range values {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}
func appendUniqueKeys(values []NodeKey, extra ...NodeKey) []NodeKey {
	seen := map[NodeKey]bool{}
	for _, value := range values {
		if value != "" {
			seen[value] = true
		}
	}
	result := values[:0]
	for _, value := range values {
		if value != "" {
			result = append(result, value)
		}
	}
	for _, value := range extra {
		if value != "" && !seen[value] {
			result = append(result, value)
			seen[value] = true
		}
	}
	return result
}
func choose(test bool, a, b string) string {
	if test {
		return a
	}
	return b
}
func appendUnique(values []string, extra ...string) []string {
	seen := map[string]bool{}
	for _, v := range values {
		seen[v] = true
	}
	for _, v := range extra {
		if v != "" && !seen[v] {
			values = append(values, v)
			seen[v] = true
		}
	}
	return values
}
func containsPlatform(v string) bool {
	for _, p := range []string{"windows", "darwin", "linux", "ios", "android"} {
		if v == p {
			return true
		}
	}
	return false
}
func containsArchitecture(value string) bool {
	return value == "amd64" || value == "arm64" || value == "arm" || value == "386" || value == "universal"
}
func platformSupportsFormat(platform, format string) bool {
	allowed := map[string][]string{
		"windows": {"nsis", "msix"},
		"darwin":  {"app", "dmg"},
		"linux":   {"appimage", "deb", "rpm", "archlinux"},
		"ios":     {"app", "ipa"},
		"android": {"apk", "aab"},
	}
	for _, candidate := range allowed[platform] {
		if candidate == format {
			return true
		}
	}
	return false
}
func packagePlatform(p manifest.Packages, os string) manifest.PackagePlatform {
	switch os {
	case "windows":
		return p.Windows
	case "darwin":
		return p.Darwin
	case "linux":
		return p.Linux
	case "ios":
		return p.IOS
	default:
		return p.Android
	}
}
func targetConfig(targets manifest.Targets, platform, arch string) manifest.Target {
	var value manifest.Platform
	switch platform {
	case "windows":
		value = targets.Windows
	case "darwin":
		value = targets.Darwin
	case "linux":
		value = targets.Linux
	case "ios":
		value = targets.IOS
	case "android":
		value = targets.Android
	}
	switch arch {
	case "amd64":
		return value.AMD64
	case "arm64":
		return value.ARM64
	case "arm":
		return value.ARM
	case "386":
		return value.X86
	case "universal":
		return value.Universal
	}
	return manifest.Target{}
}
func platformConfig(targets manifest.Targets, platform string) manifest.Platform {
	switch platform {
	case "windows":
		return targets.Windows
	case "darwin":
		return targets.Darwin
	case "linux":
		return targets.Linux
	case "ios":
		return targets.IOS
	default:
		return targets.Android
	}
}
func effectiveTarget(config manifest.Config, platform, arch string) (manifest.Project, manifest.Target) {
	project := config.Project
	common := platformConfig(config.Targets, platform)
	target := targetConfig(config.Targets, platform, arch)
	if common.ProductName != "" {
		project.ProductName = common.ProductName
	}
	if common.Identifier != "" {
		project.Identifier = common.Identifier
	}
	if common.BuildNumber != 0 {
		project.BuildNumber = common.BuildNumber
	}
	if target.BuildNumber != 0 {
		project.BuildNumber = target.BuildNumber
	}
	if target.MinimumVersion == "" {
		target.MinimumVersion = common.MinimumVersion
	}
	return project, target
}
func signingPlatform(s manifest.Signing, os string) manifest.SigningPlatform {
	switch os {
	case "windows":
		return s.Windows
	case "darwin":
		return s.Darwin
	case "linux":
		return s.Linux
	case "ios":
		return s.IOS
	default:
		return s.Android
	}
}
func packageFormat(p manifest.PackagePlatform, f string) (manifest.PackageFormat, error) {
	switch f {
	case "nsis":
		return p.NSIS, nil
	case "msix":
		return p.MSIX, nil
	case "app":
		return p.App, nil
	case "dmg":
		return p.DMG, nil
	case "appimage":
		return p.AppImage, nil
	case "deb":
		return p.Deb, nil
	case "rpm":
		return p.RPM, nil
	case "archlinux":
		return p.ArchLinux, nil
	case "ipa":
		return p.IPA, nil
	case "apk":
		return p.APK, nil
	case "aab":
		return p.AAB, nil
	default:
		return manifest.PackageFormat{}, fmt.Errorf("unsupported package format %q", f)
	}
}
func packageOutput(c manifest.Config, os, arch, format string, multiTarget bool) string {
	directory := c.Build.OutputDirectory
	if multiTarget {
		directory = filepath.Join(directory, os+"-"+arch)
	}
	base := filepath.Join(directory, c.Project.BinaryName)
	switch format {
	case "app":
		return filepath.ToSlash(base + ".app")
	case "dmg":
		return filepath.ToSlash(base + ".dmg")
	case "nsis":
		return filepath.ToSlash(base + "-installer.exe")
	case "msix":
		return filepath.ToSlash(base + ".msix")
	case "appimage":
		return filepath.ToSlash(base + "-" + arch + ".AppImage")
	case "deb":
		return filepath.ToSlash(base + "_" + c.Project.Version + "_" + arch + ".deb")
	case "rpm":
		return filepath.ToSlash(base + "-" + c.Project.Version + "." + arch + ".rpm")
	case "archlinux":
		return filepath.ToSlash(base + "-" + c.Project.Version + "-" + arch + ".pkg.tar.zst")
	case "ipa":
		return filepath.ToSlash(base + ".ipa")
	case "apk":
		return filepath.ToSlash(base + ".apk")
	case "aab":
		return filepath.ToSlash(base + ".aab")
	}
	return filepath.ToSlash(base + "." + format)
}
func packageExclusive(platform, format string) string {
	if platform == "darwin" && (format == "app" || format == "dmg") {
		// DMG assembly materialises the same .app bundle as the app format.
		return "darwin-app-bundle"
	}
	if platform == "ios" {
		return "codesign"
	}
	return ""
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
