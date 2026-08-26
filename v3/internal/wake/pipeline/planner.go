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
	if request.Verb == "" {
		request.Verb = "build"
	}
	outcomes, err := resolveBuildOutcomes(config, request)
	if err != nil {
		return Plan{}, err
	}
	combined := Plan{Name: request.Verb, Intent: BuildIntent{Command: request.Verb, Profile: config.Profile}, Nodes: map[NodeKey]Node{}}
	for _, outcome := range outcomes {
		combined.Intent.Targets = append(combined.Intent.Targets, TargetIntent{Target: outcome.target, Formats: append([]string(nil), outcome.formats...), Sign: outcome.sign, Notarize: outcome.notarize, Destination: outcome.destination})
		childConfig := config
		if outcome.notarize {
			childConfig.Signing.Darwin.Notarize = true
		}
		childRequest := request
		childRequest.TargetOS, childRequest.TargetArch, childRequest.Targets = outcome.target.OS, outcome.target.Arch, nil
		childRequest.Formats = append([]string(nil), outcome.formats...)
		childRequest.resolved, childRequest.sign, childRequest.notarize, childRequest.destination = true, outcome.sign, outcome.notarize, outcome.destination
		child, err := planTarget(childConfig, childRequest, len(outcomes) > 1)
		if err != nil {
			return Plan{}, err
		}
		if err := addPlanNodes(&combined, child); err != nil {
			return Plan{}, err
		}
	}
	terminalRoots := append([]NodeKey(nil), combined.Roots...)
	if len(combined.Artifacts) == 0 {
		combined.Artifacts = append([]NodeKey(nil), terminalRoots...)
	}
	references := make([]ArtifactReference, 0, len(combined.Artifacts))
	for _, key := range combined.Artifacts {
		node := combined.Nodes[key]
		references = append(references, ArtifactReference{Key: key, Path: node.Output, Identity: node.Artifact})
	}
	collectCache := CacheArtifact
	for _, key := range appendUniqueKeys(append([]NodeKey(nil), combined.Artifacts...), terminalRoots...) {
		if combined.Nodes[key].Cache == CacheNever {
			// A dependency that always runs can produce different bytes with the
			// same Plan identity (signing is the important case). Rebuild the
			// receipt so it always describes those non-reproducible outputs.
			collectCache = CacheNever
			break
		}
	}
	const collectKey NodeKey = "collect:artifacts"
	combined.Nodes[collectKey] = Node{
		Key:          collectKey,
		Kind:         CollectArtifacts,
		Label:        "Verify and collect final artifacts",
		Scope:        InvocationScope,
		Dependencies: appendUniqueKeys(append([]NodeKey(nil), combined.Artifacts...), terminalRoots...),
		Spec:         CollectSpec{Artifacts: references, Receipt: ".wails/artifacts/receipt.json"},
		Output:       ".wails/artifacts/receipt.json",
		Cache:        collectCache,
		Claims:       ResourceClaims{CPU: 1, MemoryMB: 64},
		EstimateMS:   10,
	}
	combined.Roots = []NodeKey{collectKey}
	return combined, combined.Validate(config.Root)
}

func addPlanNodes(combined *Plan, child Plan) error {
	for key, node := range child.Nodes {
		if existing, ok := combined.Nodes[key]; ok {
			existingWithoutOrigins, nodeWithoutOrigins := existing, node
			existingWithoutOrigins.Origins, nodeWithoutOrigins.Origins = nil, nil
			if !reflect.DeepEqual(existingWithoutOrigins, nodeWithoutOrigins) {
				return fmt.Errorf("project node %s resolves differently across targets; use equivalent target settings or separate invocations", key)
			}
			existing.Origins = mergeOrigins(existing.Origins, node.Origins)
			combined.Nodes[key] = existing
			continue
		}
		combined.Nodes[key] = node
	}
	combined.Roots = appendUniqueKeys(combined.Roots, child.Roots...)
	combined.Artifacts = appendUniqueKeys(combined.Artifacts, child.Artifacts...)
	if combined.Target == "" {
		combined.Target = child.Target
	} else {
		combined.Target += "," + child.Target
	}
	return nil
}

func mergeOrigins(left, right []OriginReference) []OriginReference {
	result := append(append([]OriginReference(nil), left...), right...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].Field < result[j].Field })
	write := 0
	for _, origin := range result {
		if write != 0 && result[write-1].Field == origin.Field {
			continue
		}
		result[write] = origin
		write++
	}
	return result[:write]
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
	capability, supported := lookupTarget(request.TargetOS, request.TargetArch)
	if !supported {
		return Plan{}, unsupportedTargetError(request.TargetOS, request.TargetArch)
	}

	profileName := config.Profile
	if profileName == "" {
		profileName = "default"
	}
	target := request.TargetOS + "/" + request.TargetArch
	project, targetSettings := effectiveTarget(config, request.TargetOS, request.TargetArch)
	destination := request.destination
	if request.TargetOS == "ios" && destination == "" {
		destination = "simulator"
	}
	if !capability.SupportsToolchain(targetSettings.Toolchain) {
		return Plan{}, fmt.Errorf("toolchain %q is not supported for target %s", targetSettings.Toolchain, target)
	}
	platformSettings := platformConfig(config.Targets, request.TargetOS)
	associations := associationsForPlatform(config.Associations, request.TargetOS)
	protocols := protocolsForPlatform(config.Protocols, request.TargetOS)
	targetDirectory := strings.ReplaceAll(target, "/", "-")
	generatedRoot := filepath.Join(".wails", "build", profileName, targetDirectory)
	if request.Development {
		generatedRoot = filepath.Join(".wails", "dev", targetDirectory)
	}
	assetsOut := filepath.ToSlash(filepath.Join(generatedRoot, "assets"))
	binaryName := config.Project.BinaryName
	if capability.Runnable == runnableBinary {
		binaryName += capability.RunnableSuffix
	}
	finalBinaryDirectory := config.Build.OutputDirectory
	if request.Development {
		finalBinaryDirectory = generatedRoot
	} else if multiTarget {
		finalBinaryDirectory = filepath.Join(finalBinaryDirectory, strings.ReplaceAll(target, "/", "-"))
	}
	finalBinaryOut := filepath.ToSlash(filepath.Join(finalBinaryDirectory, binaryName))
	binaryOut := filepath.ToSlash(filepath.Join(generatedRoot, "artifacts", binaryName))
	if request.Development {
		binaryOut = finalBinaryOut
	}
	if request.TargetOS == "ios" {
		binaryOut += ".a"
	}
	if request.TargetOS == "android" {
		binaryOut = filepath.ToSlash(filepath.Join(generatedRoot, "libwails.so"))
	}
	plan := Plan{Name: request.Verb, Target: target, Nodes: map[NodeKey]Node{}}
	add := func(node Node) NodeKey {
		taintCacheAfterAlwaysRunHook(plan, &node)
		return addNode(&plan, node, originsForNode(config, node, target))
	}
	hookNode := func(phase manifest.HookPhase, scope Scope, scopeOutput string, dependencies []NodeKey) NodeKey {
		hook, ok := config.Hooks[phase]
		if !ok {
			return ""
		}
		key := "hook:" + string(phase)
		if scope != ProjectScope {
			key += ":" + strings.ReplaceAll(target, "/", "-")
		}
		cachePolicy, cacheOutput := CacheNever, ""
		scriptInput := hook.Script
		if resolved, err := manifest.ResolveProjectPath(config.Root, "hook script", hook.Script, true); err == nil {
			scriptInput = resolved
		}
		inputs := []InputSpec{{Label: "hook-script", Files: []string{scriptInput}}}
		if hook.Cache {
			cachePolicy = CacheArtifact
			cacheOutput, _ = manifest.HookOutputRoot(hook.Outputs)
			declaredInputs := make([]string, len(hook.Inputs))
			for index, input := range hook.Inputs {
				declaredInputs[index] = input
				if resolved, err := manifest.ResolveProjectPath(config.Root, "hook input", input, true); err == nil {
					declaredInputs[index] = resolved
				}
			}
			inputs = append(inputs, InputSpec{Label: "hook-inputs", Files: declaredInputs})
		}
		return add(Node{Key: NodeKey(key), Kind: RunHook, Label: "Run " + strings.ReplaceAll(string(phase), "_", " ") + " hook", Scope: scope, Dependencies: append([]NodeKey(nil), dependencies...),
			Spec:   HookSpec{Phase: phase, Script: hook.Script, Directory: hook.Directory, Profile: profileName, TargetOS: choose(scope == ProjectScope, "", request.TargetOS), TargetArch: choose(scope == ProjectScope, "", request.TargetArch), ScopeOutput: scopeOutput, DeclaredOutputs: append([]string(nil), hook.Outputs...), EnvironmentVersion: 1},
			Inputs: inputs, Output: cacheOutput, Cache: cachePolicy, Claims: ResourceClaims{CPU: 1, MemoryMB: 256}, EstimateMS: 50})
	}
	beforeBuild := hookNode(manifest.BeforeBuild, ProjectScope, "", nil)
	projectHookDeps := []NodeKey(nil)
	if beforeBuild != "" {
		projectHookDeps = []NodeKey{beforeBuild}
	}
	publish := func(source NodeKey, destination string) NodeKey {
		node := plan.Nodes[source]
		cachePolicy := CacheArtifact
		if node.Cache == CacheNever {
			cachePolicy = CacheNever
		}
		return add(Node{Key: NodeKey("publish:" + strings.TrimPrefix(string(source), "package:")), Kind: PublishArtifact, Label: "Publish " + filepath.Base(destination), Scope: node.Scope, Dependencies: []NodeKey{source},
			Spec: PublishSpec{Source: node.Output, Destination: destination}, Output: destination,
			Cache: cachePolicy, Claims: ResourceClaims{CPU: 1, MemoryMB: 128}, EstimateMS: 25, Artifact: node.Artifact})
	}

	install := add(Node{Key: "frontend:install", Kind: InstallFrontendDependencies, Label: "Install frontend dependencies", Scope: ProjectScope, Dependencies: projectHookDeps,
		Spec:   InstallSpec{Manager: config.Frontend.PackageManager, Directory: config.Frontend.Directory, Command: config.Frontend.InstallCommand, Arguments: append([]string(nil), config.Frontend.Install...), Environment: cloneStringMap(config.Frontend.Environment)},
		Inputs: []InputSpec{{Label: "frontend-install", Files: frontendInstallFiles(config)}}, Marker: filepath.ToSlash(filepath.Join(config.Frontend.Directory, "node_modules")), Cache: CacheReceipt, Claims: ResourceClaims{CPU: 1, MemoryMB: 512}, EstimateMS: 1800})
	var tags []string
	if request.Development {
		tags = append(tags, config.Dev.Tags...)
		tags = appendUnique(tags, request.ExtraTags...)
	} else {
		tags = append(tags, config.Build.Go.Tags...)
		tags = appendUnique(tags, targetSettings.Tags...)
		tags = appendUnique(tags, "production")
		tags = appendUnique(tags, request.ExtraTags...)
	}
	obfuscated := false
	if !request.Development {
		obfuscated = config.Build.Obfuscation
		if targetSettings.ObfuscatedSet {
			obfuscated = targetSettings.Obfuscated
		}
		if request.Obfuscated {
			obfuscated = true
		}
	}
	if obfuscated {
		tags = appendUnique(tags, "wails_obfuscated")
	}
	localInputs, err := goLocalSourceInputs(config.Root)
	if err != nil {
		return Plan{}, err
	}
	if request.Development && config.Dev.UseGitIgnore {
		for index := range localInputs {
			localInputs[index].UseGitIgnore = true
		}
	}
	localRoots := make([]string, len(localInputs))
	for index := range localInputs {
		localRoots[index] = localInputs[index].Root
	}
	bindingsOut := filepath.ToSlash(filepath.Join(config.Frontend.Directory, config.Frontend.Bindings.OutputDirectory))
	bindingInputs := []InputSpec{
		{Label: "go-binding-api", Root: ".", SemanticGo: true, IncludeNames: []string{"go.mod", "go.sum", "go.work", "go.work.sum"}, IncludeExtensions: []string{".go", ".c", ".cc", ".cpp", ".h", ".hh", ".hpp", ".s", ".syso"}, ExcludeDirs: []string{".git", ".wails", "bin", "build", "dist", config.Frontend.Directory, "node_modules"}, ExcludeSuffixes: []string{"_test.go"}, UseGitIgnore: request.Development && config.Dev.UseGitIgnore},
		{Label: "go-module", Files: goMetadataFiles(config.Root)},
	}
	for _, input := range localInputs {
		bindingInputs = append(bindingInputs, InputSpec{Label: "go-binding-local-api", Root: input.Root, SemanticGo: true, IncludeNames: append([]string(nil), input.IncludeNames...), IncludeExtensions: append([]string(nil), input.IncludeExtensions...), ExcludeDirs: append([]string(nil), input.ExcludeDirs...), ExcludeSuffixes: append([]string(nil), input.ExcludeSuffixes...), UseGitIgnore: input.UseGitIgnore})
	}
	bindings := add(Node{Key: "frontend:bindings", Kind: GenerateBindings, Label: "Generate bindings", Scope: ProjectScope, Dependencies: projectHookDeps,
		Spec:   BindingsSpec{Config: config.Frontend.Bindings, Tags: tags, Obfuscated: obfuscated},
		Inputs: bindingInputs,
		Output: bindingsOut, Cache: CacheArtifact, Claims: ResourceClaims{CPU: 1, MemoryMB: 512, Exclusive: "legacy-command-adapter"}, EstimateMS: 1000})
	frontendDeps := []NodeKey{install, bindings}
	frontendOut := frontendOutputPath(config.Frontend.Directory, config.Frontend.OutputDirectory)
	frontendExclude := filepath.ToSlash(filepath.Clean(config.Frontend.OutputDirectory))
	if frontendExclude == config.Frontend.Directory || strings.HasPrefix(frontendExclude, config.Frontend.Directory+"/") {
		frontendExclude = strings.TrimPrefix(strings.TrimPrefix(frontendExclude, config.Frontend.Directory), "/")
	}
	frontendInputs := []InputSpec{{Label: "frontend-source", Root: config.Frontend.Directory, IncludeAll: true, ExcludeDirs: []string{".git", ".wails", "node_modules", config.Frontend.Bindings.OutputDirectory, frontendExclude}}}
	if request.Development {
		// The persistent frontend Dev process owns source changes through HMR.
		// The finite frontend Node only bootstraps an embeddable output for Go;
		// keeping source out of its Action Key prevents a later backend edit from
		// rebuilding the frontend merely because HMR already handled a UI edit.
		frontendInputs = []InputSpec{{Label: "frontend-dev-bootstrap", Files: frontendInstallFiles(config)}}
	}
	frontend := add(Node{Key: "frontend:build", Kind: BuildFrontend, Label: "Build frontend", Scope: ProjectScope, Dependencies: frontendDeps,
		Spec:   FrontendSpec{Manager: config.Frontend.PackageManager, Directory: config.Frontend.Directory, Command: choose(request.Development, config.Frontend.BuildCommand+":dev", config.Frontend.BuildCommand), Arguments: append([]string(nil), config.Frontend.Build...), Output: config.Frontend.OutputDirectory, Production: !request.Development, Environment: cloneStringMap(config.Frontend.Environment)},
		Inputs: frontendInputs,
		Output: frontendOut, Cache: CacheArtifact, Claims: ResourceClaims{CPU: 2, MemoryMB: 1024}, EstimateMS: 900})
	var assets NodeKey
	if request.TargetOS == "windows" || request.TargetOS == "ios" || request.TargetOS == "android" {
		assets = add(platformAssetsNode(config, request, target, assetsOut, projectHookDeps))
	}

	compileDeps := []NodeKey{frontend}
	if assets != "" {
		compileDeps = append(compileDeps, assets)
	}
	compileInputs := []InputSpec{
		{Label: "go-sources", Root: ".", IncludeNames: []string{"go.mod", "go.sum", "go.work", "go.work.sum"}, IncludeExtensions: []string{".go", ".c", ".cc", ".cpp", ".h", ".hh", ".hpp", ".s", ".syso"}, ExcludeDirs: []string{".git", ".wails", "bin", "build", "dist", config.Frontend.Directory, "node_modules"}, ExcludeSuffixes: []string{"_test.go"}, UseGitIgnore: request.Development && config.Dev.UseGitIgnore},
		{Label: "go-module", Files: goMetadataFiles(config.Root)},
	}
	compileInputs = append(compileInputs, localInputs...)
	linkerFlags := append([]string(nil), config.Build.Go.LinkerFlags...)
	compilerFlags := append([]string(nil), config.Build.Go.CompilerFlags...)
	garbleArgs := append([]string(nil), config.Build.Go.GarbleArgs...)
	if targetSettings.LinkerFlags != nil {
		linkerFlags = append([]string(nil), targetSettings.LinkerFlags...)
	}
	if targetSettings.CompilerFlags != nil {
		compilerFlags = append([]string(nil), targetSettings.CompilerFlags...)
	}
	if targetSettings.GarbleArgs != nil {
		garbleArgs = append([]string(nil), targetSettings.GarbleArgs...)
	}
	garbleArgs = append(garbleArgs, request.GarbleArgs...)
	environment := cloneStringMap(config.Build.Environment)
	if targetSettings.Environment != nil {
		environment = cloneStringMap(targetSettings.Environment)
	}
	trimPath, strip := config.Build.TrimPath, config.Build.Strip
	if request.Development {
		linkerFlags, compilerFlags, garbleArgs = nil, nil, nil
		environment = nil
		trimPath, strip = false, false
	}
	compileTargets := []Target{{OS: request.TargetOS, Arch: request.TargetArch}}
	if capability.Synthetic() {
		compileTargets = make([]Target, int(capability.ComponentCount))
		for index := range compileTargets {
			compileTargets[index] = capability.Component(index)
		}
	}
	compileRoots := make([]NodeKey, 0, len(compileTargets))
	componentBinaries := make([]ComponentBinary, 0, len(compileTargets))
	for _, compileTarget := range compileTargets {
		compileTargetName := compileTarget.OS + "/" + compileTarget.Arch
		compileOutput := binaryOut
		if capability.Synthetic() {
			componentRoot := filepath.Join(".wails", "build", profileName, strings.ReplaceAll(compileTargetName, "/", "-"))
			compileOutput = filepath.ToSlash(filepath.Join(componentRoot, binaryName))
			if request.TargetOS == "android" {
				compileOutput = filepath.ToSlash(filepath.Join(componentRoot, "libwails.so"))
			}
		}
		compileKey := add(Node{Key: NodeKey("target:" + compileTargetName + ":compile"), Kind: CompileApplication, Label: "Compile " + compileTargetName, Scope: TargetScope, Dependencies: compileDeps,
			Spec:   CompileSpec{TargetOS: compileTarget.OS, TargetArch: compileTarget.Arch, Output: compileOutput, Assets: assetsOut, Destination: destination, MinimumVersion: targetSettings.MinimumVersion, Tags: append([]string(nil), tags...), LinkerFlags: append([]string(nil), linkerFlags...), CompilerFlags: append([]string(nil), compilerFlags...), GarbleArgs: append([]string(nil), garbleArgs...), LocalRoots: append([]string(nil), localRoots...), Production: !request.Development, Obfuscated: obfuscated, TrimPath: trimPath, Strip: strip, VCSInfo: config.Build.VCSInfo && !request.Development, Toolchain: targetSettings.Toolchain, Environment: cloneStringMap(environment)},
			Inputs: compileInputs,
			Output: compileOutput, Cache: CacheArtifact, Claims: ResourceClaims{CPU: max(1, runtime.GOMAXPROCS(0)/2), MemoryMB: 2048}, EstimateMS: 1500, Artifact: ArtifactIdentity{Kind: ArtifactBinary, Target: Target{OS: compileTarget.OS, Arch: compileTarget.Arch}}})
		compileRoots = append(compileRoots, compileKey)
		componentBinaries = append(componentBinaries, ComponentBinary{Arch: compileTarget.Arch, Path: compileOutput})
	}
	buildDependencies := append([]NodeKey(nil), compileRoots...)
	if capability.Synthetic() && request.TargetOS == "darwin" {
		inputs := make([]string, len(componentBinaries))
		for index := range componentBinaries {
			inputs[index] = componentBinaries[index].Path
		}
		merge := add(Node{Key: NodeKey("assemble:" + target + ":binary"), Kind: MergeUniversalBinaries, Label: "Merge " + target + " binaries", Scope: TargetScope, Dependencies: compileRoots,
			Spec: MergeSpec{Inputs: append([]string(nil), inputs...), Output: binaryOut}, Inputs: []InputSpec{{Label: "component-binaries", Files: inputs}}, Output: binaryOut,
			Cache: CacheArtifact, Claims: ResourceClaims{CPU: 1, MemoryMB: 256, Exclusive: "darwin-lipo"}, EstimateMS: 100, Artifact: ArtifactIdentity{Kind: ArtifactBinary, Target: Target{OS: request.TargetOS, Arch: request.TargetArch}}})
		buildDependencies = []NodeKey{merge}
	}
	lastBuild := buildDependencies[len(buildDependencies)-1]
	buildScopeOutput := plan.Nodes[lastBuild].Output
	if len(buildDependencies) > 1 {
		buildScopeOutput = packageScopeOutput(nodeOutputs(plan, buildDependencies))
	}
	afterBuild := hookNode(manifest.AfterBuild, TargetScope, buildScopeOutput, buildDependencies)
	buildBarrier := append([]NodeKey(nil), buildDependencies...)
	if afterBuild != "" {
		buildBarrier = []NodeKey{afterBuild}
	}
	lastRunnable := lastBuild
	if capability.Runnable == runnableApp {
		assemblyDependencies := append([]NodeKey(nil), buildBarrier...)
		if assets == "" {
			assets = add(platformAssetsNode(config, request, target, assetsOut, projectHookDeps))
		}
		assemblyDependencies = append(assemblyDependencies, assets)
		finalOutput := runnableOutput(config, request.TargetOS, request.TargetArch, capability, multiTarget)
		output := filepath.ToSlash(filepath.Join(generatedRoot, "artifacts", filepath.Base(finalOutput)))
		if request.Development {
			output = finalOutput
		}
		packageConfig := registeredPackageFormat(config.Package, request.TargetOS, "app")
		cachePolicy := CacheArtifact
		if request.TargetOS == "ios" {
			cachePolicy = CacheNever
		}
		lastRunnable = add(Node{Key: NodeKey("assemble:" + target), Kind: AssembleApplication, Label: "Assemble " + target + " application", Scope: TargetScope, Dependencies: assemblyDependencies,
			Spec:   PackageSpec{TargetOS: request.TargetOS, TargetArch: request.TargetArch, Format: "app", Binary: binaryOut, Binaries: append([]ComponentBinary(nil), componentBinaries...), Assets: assetsOut, Output: output, Profile: profileName, Destination: destination, MinimumVersion: targetSettings.MinimumVersion, Config: packageConfig, Project: project, Capabilities: platformSettings.Capabilities, Associations: associations, Protocols: protocols},
			Inputs: []InputSpec{{Label: "compiled-application", Files: []string{binaryOut, assetsOut}}}, Output: output, Cache: cachePolicy, Claims: ResourceClaims{CPU: 1, MemoryMB: 1024, Exclusive: packageExclusive(request.TargetOS, "app")}, EstimateMS: 500, Artifact: ArtifactIdentity{Kind: ArtifactBundle, Target: Target{OS: request.TargetOS, Arch: request.TargetArch}, Format: "app"}})
	}

	formats := append([]string(nil), request.Formats...)
	if len(formats) == 0 {
		if !request.sign {
			if !request.Development {
				destination := finalBinaryOut
				if capability.Runnable == runnableApp {
					destination = runnableOutput(config, request.TargetOS, request.TargetArch, capability, multiTarget)
				}
				published := publish(lastRunnable, destination)
				if afterBuild != "" && capability.Runnable != runnableApp {
					node := plan.Nodes[published]
					node.Dependencies = appendUniqueKeys(node.Dependencies, afterBuild)
					taintCacheAfterAlwaysRunHook(plan, &node)
					plan.Nodes[published] = node
				}
				lastRunnable = published
			}
			plan.Artifacts = []NodeKey{lastRunnable}
			plan.Roots = []NodeKey{lastRunnable}
			if request.Development && afterBuild != "" && capability.Runnable != runnableApp {
				plan.Roots = []NodeKey{afterBuild}
				plan.Artifacts = []NodeKey{lastRunnable}
			}
			return plan, plan.Validate(config.Root)
		}
		input := plan.Nodes[lastRunnable].Output
		identity := plan.Nodes[lastRunnable].Artifact
		identity.Signed, identity.Notarized = true, config.Signing.Darwin.Notarize
		beforeSignDeps := append([]NodeKey{lastRunnable}, buildBarrier...)
		beforeSign := hookNode(manifest.BeforeSign, TargetScope, input, beforeSignDeps)
		if beforeSign != "" {
			beforeSignDeps = []NodeKey{lastRunnable, beforeSign}
		}
		signed := add(Node{Key: NodeKey(string(lastRunnable) + ":sign"), Kind: SignArtifact, Label: "Sign " + filepath.Base(input), Scope: TargetScope, Dependencies: appendUniqueKeys(nil, beforeSignDeps...), Spec: SignSpec{TargetOS: request.TargetOS, TargetArch: request.TargetArch, Format: plan.Nodes[lastRunnable].Artifact.DisplayKind(), Input: input, Config: signingPlatformForPlan(config.Signing, request.TargetOS, assetsOut)}, Output: input + ".signed", Cache: CacheNever, Claims: ResourceClaims{CPU: 1, MemoryMB: 256, Exclusive: "sign"}, EstimateMS: 1000, Artifact: identity})
		afterSign := hookNode(manifest.AfterSign, TargetScope, plan.Nodes[signed].Output, []NodeKey{signed})
		if !request.Development {
			destination := finalBinaryOut
			if capability.Runnable == runnableApp {
				destination = runnableOutput(config, request.TargetOS, request.TargetArch, capability, multiTarget)
			}
			published := publish(signed, destination+".signed")
			if afterSign != "" {
				node := plan.Nodes[published]
				node.Dependencies = appendUniqueKeys(node.Dependencies, afterSign)
				taintCacheAfterAlwaysRunHook(plan, &node)
				plan.Nodes[published] = node
			}
			signed = published
		}
		plan.Artifacts = []NodeKey{signed}
		plan.Roots = []NodeKey{signed}
		if request.Development && afterSign != "" {
			plan.Roots = []NodeKey{afterSign}
			plan.Artifacts = []NodeKey{signed}
		}
		return plan, plan.Validate(config.Root)
	}
	packageDependencies := buildDependencies
	if capability.Runnable == runnableApp {
		packageDependencies = []NodeKey{lastRunnable}
	}
	beforePackageDeps := append([]NodeKey(nil), buildBarrier...)
	if capability.Runnable == runnableApp {
		beforePackageDeps = []NodeKey{lastRunnable}
	}
	for _, format := range formats {
		formatCapability, known := lookupFormat(format)
		if !known {
			return Plan{}, fmt.Errorf("unknown package format %q", format)
		}
		if !capability.SupportsFormat(format, request.Development) {
			mode := "production"
			if request.Development {
				mode = "development"
			}
			return Plan{}, fmt.Errorf("package format %q is not supported for %s in %s", format, target, mode)
		}
		if formatCapability.RequiredDestination != "" && destination != formatCapability.RequiredDestination {
			return Plan{}, fmt.Errorf("%s %s packaging requires profile destination = %q", request.TargetOS, strings.ToUpper(format), formatCapability.RequiredDestination)
		}
	}
	if assets == "" {
		assets = add(platformAssetsNode(config, request, target, assetsOut, beforePackageDeps))
	}
	packageOutputs := make([]string, 0, len(formats))
	for _, format := range formats {
		finalOutput := packageOutput(config, request.TargetOS, request.TargetArch, format, multiTarget)
		output := filepath.ToSlash(filepath.Join(generatedRoot, "artifacts", filepath.Base(finalOutput)))
		if request.Development {
			output = finalOutput
		}
		packageOutputs = append(packageOutputs, output)
	}
	packageHookOutput := packageScopeOutput(packageOutputs)
	beforePackage := hookNode(manifest.BeforePackage, TargetScope, packageHookOutput, appendUniqueKeys(append([]NodeKey(nil), beforePackageDeps...), assets))
	var packageRoots []NodeKey
	finalOutputs := make(map[NodeKey]string, len(formats))
	for _, format := range formats {
		pkgConfig := registeredPackageFormat(config.Package, request.TargetOS, format)
		finalOutput := packageOutput(config, request.TargetOS, request.TargetArch, format, multiTarget)
		output := filepath.ToSlash(filepath.Join(generatedRoot, "artifacts", filepath.Base(finalOutput)))
		if request.Development {
			output = finalOutput
		}
		key := NodeKey("package:" + target + ":" + format)
		packageDeps := appendUniqueKeys(append([]NodeKey(nil), packageDependencies...), assets)
		packageDeps = appendUniqueKeys(packageDeps, beforePackageDeps...)
		if beforePackage != "" {
			packageDeps = appendUniqueKeys(packageDeps, beforePackage)
		}
		packageCache := CacheArtifact
		if request.TargetOS == "ios" {
			// iOS bundle assembly invokes codesign, including ad-hoc simulator
			// signing, so its result must never enter the reusable artifact cache.
			packageCache = CacheNever
		}
		packageBinary := binaryOut
		if format == "ipa" {
			packageBinary = plan.Nodes[lastRunnable].Output
		}
		pkg := add(Node{Key: key, Kind: PackageArtifact, Label: "Package " + format, Scope: PackageScope, Dependencies: packageDeps,
			Spec:   PackageSpec{TargetOS: request.TargetOS, TargetArch: request.TargetArch, Format: format, Binary: packageBinary, Binaries: append([]ComponentBinary(nil), componentBinaries...), Assets: assetsOut, Output: output, Profile: profileName, Destination: destination, MinimumVersion: targetSettings.MinimumVersion, Config: pkgConfig, Project: project, Capabilities: platformSettings.Capabilities, Associations: associations, Protocols: protocols},
			Inputs: packageInputs(config.Root, request.TargetOS, format, pkgConfig), Output: output, Cache: packageCache, Claims: ResourceClaims{CPU: 1, MemoryMB: 1024, Exclusive: packageExclusive(request.TargetOS, format)}, EstimateMS: 1000, Artifact: ArtifactIdentity{Kind: ArtifactPackage, Target: Target{OS: request.TargetOS, Arch: request.TargetArch}, Format: format}})
		packageRoots = append(packageRoots, pkg)
		finalOutputs[pkg] = finalOutput
	}
	packageArtifacts := append([]NodeKey(nil), packageRoots...)
	afterPackage := hookNode(manifest.AfterPackage, TargetScope, packageHookOutput, packageRoots)
	terminalHook := afterPackage
	if request.sign {
		beforeSignDeps := append([]NodeKey(nil), packageRoots...)
		if afterPackage != "" {
			beforeSignDeps = []NodeKey{afterPackage}
		}
		beforeSign := hookNode(manifest.BeforeSign, TargetScope, packageHookOutput, beforeSignDeps)
		var signed []NodeKey
		for _, artifact := range packageArtifacts {
			input := plan.Nodes[artifact].Output
			dependencies := appendUniqueKeys([]NodeKey{artifact}, packageRoots...)
			if beforeSign != "" {
				dependencies = appendUniqueKeys(dependencies, beforeSign)
			} else if afterPackage != "" {
				dependencies = appendUniqueKeys(dependencies, afterPackage)
			}
			identity := plan.Nodes[artifact].Artifact
			identity.Signed, identity.Notarized = true, config.Signing.Darwin.Notarize
			signedKey := add(Node{Key: NodeKey(string(artifact) + ":sign"), Kind: SignArtifact, Label: "Sign " + filepath.Base(input), Scope: PackageScope, Dependencies: dependencies, Spec: SignSpec{TargetOS: request.TargetOS, TargetArch: request.TargetArch, Format: identity.Format, Input: input, Config: signingPlatformForPlan(config.Signing, request.TargetOS, assetsOut)}, Output: input + ".signed", Cache: CacheNever, Claims: ResourceClaims{CPU: 1, MemoryMB: 256, Exclusive: "sign"}, EstimateMS: 1000, Artifact: identity})
			signed = append(signed, signedKey)
			finalOutputs[signedKey] = finalOutputs[artifact] + ".signed"
		}
		packageRoots = signed
		afterSign := hookNode(manifest.AfterSign, TargetScope, packageScopeOutput(nodeOutputs(plan, signed)), signed)
		if afterSign != "" {
			terminalHook = afterSign
		}
	}
	if !request.Development {
		published := make([]NodeKey, 0, len(packageRoots))
		for _, artifact := range packageRoots {
			key := publish(artifact, finalOutputs[artifact])
			node := plan.Nodes[key]
			if terminalHook != "" {
				node.Dependencies = appendUniqueKeys(node.Dependencies, terminalHook)
			}
			taintCacheAfterAlwaysRunHook(plan, &node)
			plan.Nodes[key] = node
			published = append(published, key)
		}
		packageRoots = published
	}
	plan.Artifacts = append([]NodeKey(nil), packageRoots...)
	if len(plan.Roots) == 0 {
		plan.Roots = packageRoots
		if request.Development && terminalHook != "" {
			plan.Roots = []NodeKey{terminalHook}
		}
	}
	return plan, plan.Validate(config.Root)
}

func nodeOutputs(plan Plan, keys []NodeKey) []string {
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, plan.Nodes[key].Output)
	}
	return result
}

func dependsOnAlwaysRunHook(plan Plan, dependencies []NodeKey) bool {
	visited := make(map[NodeKey]bool, len(dependencies))
	var visit func(NodeKey) bool
	visit = func(key NodeKey) bool {
		if visited[key] {
			return false
		}
		visited[key] = true
		node, ok := plan.Nodes[key]
		if !ok {
			return false
		}
		if node.Kind == RunHook && node.Cache == CacheNever {
			return true
		}
		for _, dependency := range node.Dependencies {
			if visit(dependency) {
				return true
			}
		}
		return false
	}
	for _, dependency := range dependencies {
		if visit(dependency) {
			return true
		}
	}
	return false
}

func taintCacheAfterAlwaysRunHook(plan Plan, node *Node) {
	if node.Cache != CacheNever && dependsOnAlwaysRunHook(plan, node.Dependencies) {
		node.Cache = CacheNever
	}
}

func packageScopeOutput(outputs []string) string {
	if len(outputs) == 1 {
		return outputs[0]
	}
	if len(outputs) == 0 {
		return ""
	}
	root := filepath.ToSlash(filepath.Dir(outputs[0]))
	for _, output := range outputs[1:] {
		for root != "." && output != root && !strings.HasPrefix(filepath.ToSlash(filepath.Clean(output)), root+"/") {
			root = filepath.ToSlash(filepath.Dir(root))
		}
	}
	return root
}

func addNode(plan *Plan, node Node, origins []OriginReference) NodeKey {
	if _, exists := plan.Nodes[node.Key]; exists {
		panic("duplicate planner key: " + string(node.Key))
	}
	node.Origins = origins
	plan.Nodes[node.Key] = node
	return node.Key
}

func registeredPackageFormat(packages manifest.Packages, platform, format string) manifest.PackageFormat {
	result, err := manifest.ResolvePackageFormat(packages, platform, format)
	if err != nil {
		// Selection has already accepted this pair from the closed capability
		// registry. A mismatch here is an internal registry defect, not a user
		// configuration error.
		panic(fmt.Sprintf("pipeline package registry mismatch for %s/%s: %v", platform, format, err))
	}
	return result
}

func originsForNode(config manifest.Config, node Node, target string) []OriginReference {
	platform, _, _ := strings.Cut(target, "/")
	profile := config.Profile
	if profile == "" {
		profile = config.Selected.Name
	}
	prefixes := []string{"project.binary_name"}
	switch node.Kind {
	case InstallFrontendDependencies:
		prefixes = []string{"frontend.directory", "frontend.install", "frontend.environment"}
	case GenerateBindings:
		prefixes = []string{"frontend.bindings", "build.tags", "build.obfuscated", `target["` + target + `"].tags`, `target["` + target + `"].obfuscated`}
	case BuildFrontend:
		prefixes = []string{"frontend.directory", "frontend.build", "frontend.output", "frontend.environment"}
	case CompileApplication:
		prefixes = append(prefixes, "build", platform, `target["`+target+`"]`)
	case MergeUniversalBinaries:
		prefixes = append(prefixes, "build.output", `target["`+target+`"]`)
	case GeneratePlatformAssets:
		prefixes = []string{"project", platform, "file_association", "protocol"}
	case AssembleApplication, PackageArtifact:
		prefixes = []string{"project", "build.output", platform, `target["` + target + `"]`, "file_association", "protocol"}
		if spec, ok := node.Spec.(PackageSpec); ok {
			prefixes = append(prefixes, `package["`+spec.Format+`"]`)
		}
	case SignArtifact:
		prefixes = []string{platform + ".signing", platform + ".notarization"}
	case CollectArtifacts:
		return nil
	case RunHook:
		if spec, ok := node.Spec.(HookSpec); ok {
			prefixes = []string{`hook["` + string(spec.Phase) + `"]`}
		}
	}
	if profile != "" {
		prefixes = append(prefixes, `profile["`+profile+`"].target["`+target+`"]`)
	}
	result := make([]OriginReference, 0, len(prefixes))
	for field, origin := range config.Origins {
		for _, prefix := range prefixes {
			if field == prefix || strings.HasPrefix(field, prefix+".") || strings.HasPrefix(field, prefix+"[") {
				result = append(result, OriginReference{Field: field, Origin: origin})
				break
			}
		}
	}
	sort.Slice(result, func(left, right int) bool { return result[left].Field < result[right].Field })
	return result
}

func cloneStringMap(source map[string]string) map[string]string {
	if source == nil {
		return nil
	}
	result := make(map[string]string, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
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

func goLocalSourceInputs(root string) ([]InputSpec, error) {
	return goLocalSourceInputsWithAbs(root, filepath.Abs)
}

func goLocalSourceInputsWithAbs(root string, abs func(string) (string, error)) ([]InputSpec, error) {
	root, err := abs(root)
	if err != nil {
		return nil, err
	}
	local := map[string]bool{}
	add := func(base, value string) {
		if !filepath.IsAbs(value) {
			value = filepath.Join(base, value)
		}
		value, err = abs(value)
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
		inputs = append(inputs, InputSpec{Label: "go-local-source", Root: path, IncludeNames: []string{"go.mod", "go.sum"}, IncludeExtensions: []string{".go", ".c", ".cc", ".cpp", ".h", ".hh", ".hpp", ".s", ".syso"}, ExcludeDirs: []string{".git", ".wails", "bin", "dist", "node_modules"}, ExcludeSuffixes: []string{"_test.go"}})
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
	return goMetadataFilesWithAbs(root, filepath.Abs)
}

func goMetadataFilesWithAbs(root string, abs func(string) (string, error)) []string {
	root, err := abs(root)
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
		for _, path := range []string{platform.Icon, platform.Manifest, platform.AssetsCar, platform.InfoPlist, platform.DesktopEntry} {
			if path != "" {
				files = append(files, path)
			}
		}
	}
	for _, signing := range []manifest.SigningPlatform{c.Signing.Windows, c.Signing.Darwin, c.Signing.Linux, c.Signing.IOS, c.Signing.Android} {
		for _, path := range []string{signing.Entitlements, signing.ProvisioningProfile} {
			if path != "" {
				files = append(files, path)
			}
		}
	}
	for _, association := range c.Associations {
		if association.Icon != "" {
			files = append(files, association.Icon)
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
	var resources []string
	switch platform {
	case "windows":
		if format == "msix" {
			resources = appendNonempty(resources, f.Manifest)
		}
	case "darwin":
		if format == "dmg" {
			resources = appendNonempty(resources, f.Background, f.VolumeIcon, f.FileIcon)
			fileNames := make([]string, 0, len(f.Files))
			for name := range f.Files {
				fileNames = append(fileNames, name)
			}
			sort.Strings(fileNames)
			for _, name := range fileNames {
				resources = appendNonempty(resources, f.Files[name])
			}
		}
	case "linux":
		switch format {
		case "appimage":
			resources = appendNonempty(resources, f.Icon, f.DesktopEntry)
		case "deb", "rpm", "archlinux":
			resources = appendNonempty(resources, f.PreInstall, f.PostInstall, f.PreRemove, f.PostRemove)
		}
	}
	if len(resources) > 0 {
		result = append(result, InputSpec{Label: "package-resources", Files: resources})
	}
	return result
}

func appendNonempty(destination []string, values ...string) []string {
	for _, value := range values {
		if value != "" {
			destination = append(destination, value)
		}
	}
	return destination
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

func signingPlatformForPlan(signing manifest.Signing, platform, assets string) manifest.SigningPlatform {
	result := signingPlatform(signing, platform)
	if result.Entitlements != "" {
		result.Entitlements = filepath.ToSlash(filepath.Join(assets, "signing", "entitlements"+filepath.Ext(result.Entitlements)))
	}
	if result.ProvisioningProfile != "" {
		result.ProvisioningProfile = filepath.ToSlash(filepath.Join(assets, "signing", "provisioning-profile"+filepath.Ext(result.ProvisioningProfile)))
	}
	return result
}
func packageOutput(c manifest.Config, os, arch, format string, multiTarget bool) string {
	directory := c.Build.OutputDirectory
	if multiTarget {
		directory = filepath.Join(directory, os+"-"+arch)
	}
	capability, ok := lookupFormat(format)
	if !ok {
		return filepath.ToSlash(filepath.Join(directory, c.Project.BinaryName) + "." + format)
	}
	return capability.OutputPath(directory, c.Project.BinaryName, c.Project.Version, arch)
}

func runnableOutput(c manifest.Config, os, arch string, capability targetCapability, multiTarget bool) string {
	directory := c.Build.OutputDirectory
	if multiTarget {
		directory = filepath.Join(directory, os+"-"+arch)
	}
	return filepath.ToSlash(filepath.Join(directory, c.Project.BinaryName) + capability.RunnableSuffix)
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
