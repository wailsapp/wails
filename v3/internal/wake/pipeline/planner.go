package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func PlanBuild(config manifest.Config, request Request) (Plan, error) {
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
	assetsOut := filepath.ToSlash(filepath.Join(".wails", "build", profileName, strings.ReplaceAll(target, "/", "-"), "assets"))
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
		beforeBuild = add(hookNode("before_build", hook, config, request, ProjectScope, nil))
	}
	install := add(Node{Key: "frontend:install", Kind: InstallFrontendDependencies, Label: "Install frontend dependencies", Scope: ProjectScope,
		Spec:   InstallSpec{Manager: config.Frontend.PackageManager, Directory: config.Frontend.Directory, Command: config.Frontend.InstallCommand},
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
	bindingsDeps := keys(beforeBuild)
	bindingsOut := filepath.ToSlash(filepath.Join(config.Frontend.Directory, config.Frontend.Bindings.OutputDirectory))
	bindings := add(Node{Key: "frontend:bindings", Kind: GenerateBindings, Label: "Generate bindings", Scope: ProjectScope, Dependencies: bindingsDeps,
		Spec: BindingsSpec{Config: config.Frontend.Bindings, Tags: tags, Obfuscated: request.Obfuscated || config.Build.Obfuscation},
		Inputs: []InputSpec{
			{Label: "go-binding-api", Root: ".", SemanticGo: true, ExcludeDirs: []string{".git", ".wails", "bin", "build", "dist", config.Frontend.Directory, "node_modules"}},
			{Label: "go-module", Files: goMetadataFiles(config.Root)},
		},
		Output: bindingsOut, Cache: CacheArtifact, Claims: ResourceClaims{CPU: 1, MemoryMB: 1024, Exclusive: "legacy-command-adapter"}, EstimateMS: 1000})
	frontendDeps := []NodeKey{install, bindings}
	if beforeBuild != "" {
		frontendDeps = append(frontendDeps, beforeBuild)
	}
	frontendOut := filepath.ToSlash(filepath.Join(config.Frontend.Directory, config.Frontend.OutputDirectory))
	frontend := add(Node{Key: "frontend:build", Kind: BuildFrontend, Label: "Build frontend", Scope: ProjectScope, Dependencies: frontendDeps,
		Spec:   FrontendSpec{Manager: config.Frontend.PackageManager, Directory: config.Frontend.Directory, Command: choose(request.Development, config.Frontend.BuildCommand+":dev", config.Frontend.BuildCommand), Output: config.Frontend.OutputDirectory, Production: !request.Development},
		Inputs: []InputSpec{{Label: "frontend-source", Root: config.Frontend.Directory, IncludeAll: true, ExcludeDirs: []string{".git", ".wails", "node_modules", config.Frontend.Bindings.OutputDirectory, config.Frontend.OutputDirectory}}},
		Output: frontendOut, Cache: CacheArtifact, Claims: ResourceClaims{CPU: 2, MemoryMB: 1536}, EstimateMS: 900})
	var assets NodeKey
	if request.TargetOS == "windows" || request.TargetOS == "ios" || request.TargetOS == "android" {
		assets = add(platformAssetsNode(config, request, target, assetsOut, keys(beforeBuild)))
	}

	binaryName := config.Project.BinaryName
	if request.TargetOS == "windows" {
		binaryName += ".exe"
	}
	binaryOut := filepath.ToSlash(filepath.Join(config.Build.OutputDirectory, binaryName))
	if request.TargetOS == "ios" {
		binaryOut += ".a"
	}
	if request.TargetOS == "android" {
		binaryOut = filepath.ToSlash(filepath.Join(".wails", "build", profileName, strings.ReplaceAll(target, "/", "-"), "libwails.so"))
	}
	compileDeps := []NodeKey{frontend}
	if beforeBuild != "" {
		compileDeps = append(compileDeps, beforeBuild)
	}
	if assets != "" {
		compileDeps = append(compileDeps, assets)
	}
	compile := add(Node{Key: NodeKey("target:" + target + ":compile"), Kind: CompileApplication, Label: "Compile " + target, Scope: TargetScope, Dependencies: compileDeps,
		Spec: CompileSpec{TargetOS: request.TargetOS, TargetArch: request.TargetArch, Output: binaryOut, Assets: assetsOut, Variant: targetSettings.Variant, MinimumVersion: targetSettings.MinimumVersion, Tags: tags, LinkerFlags: config.Build.Go.LinkerFlags, CompilerFlags: config.Build.Go.CompilerFlags, GarbleArgs: config.Build.Go.GarbleArgs, Production: !request.Development, Obfuscated: request.Obfuscated || config.Build.Obfuscation, TrimPath: config.Build.TrimPath, Strip: config.Build.Strip},
		Inputs: []InputSpec{
			{Label: "go-sources", Root: ".", IncludeNames: []string{"go.mod", "go.sum", "go.work", "go.work.sum"}, IncludeExtensions: []string{".go", ".c", ".cc", ".cpp", ".h", ".hh", ".hpp", ".s", ".syso"}, ExcludeDirs: []string{".git", ".wails", "bin", "build", "dist", config.Frontend.Directory, "node_modules"}},
			{Label: "go-module", Files: goMetadataFiles(config.Root)},
		},
		Output: binaryOut, Cache: CacheArtifact, Claims: ResourceClaims{CPU: max(1, runtime.GOMAXPROCS(0)-1), MemoryMB: 2048, Exclusive: "go-build"}, EstimateMS: 1500, ArtifactKind: "binary"})
	lastBuild := compile
	if hook := config.Hooks.AfterBuild; hook.Script != "" {
		lastBuild = add(hookNode("after_build", hook, config, request, TargetScope, []NodeKey{compile}))
	}

	if request.Verb == "build" {
		plan.Roots = []NodeKey{lastBuild}
		return plan, plan.Validate(config.Root)
	}
	if request.Verb != "package" && request.Verb != "sign" {
		return Plan{}, fmt.Errorf("unsupported pipeline verb %q", request.Verb)
	}
	beforePackageDeps := []NodeKey{lastBuild}
	if hook := config.Hooks.BeforePackage; hook.Script != "" {
		beforePackageDeps = []NodeKey{add(hookNode("before_package", hook, config, request, TargetScope, beforePackageDeps))}
	}
	if assets == "" {
		assets = add(platformAssetsNode(config, request, target, assetsOut, beforePackageDeps))
	}
	formats := request.Formats
	if len(formats) == 0 {
		formats = packagePlatform(config.Package, request.TargetOS).Formats
	}
	if len(formats) == 0 {
		return Plan{}, fmt.Errorf("no package formats configured for %s", request.TargetOS)
	}
	sort.Strings(formats)
	var packageRoots []NodeKey
	for _, format := range formats {
		if !platformSupportsFormat(request.TargetOS, format) {
			return Plan{}, fmt.Errorf("package format %q is not supported for %s", format, request.TargetOS)
		}
		pkgConfig, err := packageFormat(packagePlatform(config.Package, request.TargetOS), format)
		if err != nil {
			return Plan{}, err
		}
		if pkgConfig.Template != "" && !supportsCustomTemplate(format) {
			return Plan{}, fmt.Errorf("custom templates are not supported for package format %q", format)
		}
		if len(pkgConfig.Options) > 0 {
			return Plan{}, fmt.Errorf("package options are not yet supported for format %q; use a supported user-owned template", format)
		}
		output := packageOutput(config, request.TargetOS, request.TargetArch, format)
		key := NodeKey("package:" + target + ":" + format)
		packageDeps := appendUniqueKeys([]NodeKey{compile, assets}, beforePackageDeps...)
		pkg := add(Node{Key: key, Kind: PackageArtifact, Label: "Package " + format, Scope: PackageScope, Dependencies: packageDeps,
			Spec:   PackageSpec{TargetOS: request.TargetOS, TargetArch: request.TargetArch, Format: format, Binary: binaryOut, Assets: assetsOut, Output: output, Variant: targetSettings.Variant, MinimumVersion: targetSettings.MinimumVersion, Config: pkgConfig, Project: project},
			Inputs: templateInput(pkgConfig), Output: output, Cache: CacheArtifact, Claims: ResourceClaims{CPU: 1, MemoryMB: 1024, Exclusive: packageExclusive(format)}, EstimateMS: 1000, ArtifactKind: format})
		packageRoots = append(packageRoots, pkg)
	}
	packageArtifacts := append([]NodeKey(nil), packageRoots...)
	if hook := config.Hooks.AfterPackage; hook.Script != "" {
		packageRoots = []NodeKey{add(hookNode("after_package", hook, config, request, PackageScope, packageRoots))}
	}
	if request.Verb == "sign" {
		if hook := config.Hooks.BeforeSign; hook.Script != "" {
			packageRoots = []NodeKey{add(hookNode("before_sign", hook, config, request, PackageScope, packageRoots))}
		}
		var signed []NodeKey
		for _, artifact := range packageArtifacts {
			input := plan.Nodes[artifact].Output
			dependencies := appendUniqueKeys([]NodeKey{artifact}, packageRoots...)
			signed = append(signed, add(Node{Key: NodeKey(string(artifact) + ":sign"), Kind: SignArtifact, Label: "Sign " + filepath.Base(input), Scope: PackageScope, Dependencies: dependencies, Spec: SignSpec{TargetOS: request.TargetOS, Format: plan.Nodes[artifact].ArtifactKind, Input: input, Config: signingPlatform(config.Signing, request.TargetOS)}, Output: input + ".signed", Cache: CacheNever, Claims: ResourceClaims{CPU: 1, MemoryMB: 512, Exclusive: "sign"}, EstimateMS: 1000, ArtifactKind: "signed"}))
		}
		packageRoots = signed
		if hook := config.Hooks.AfterSign; hook.Script != "" {
			packageRoots = []NodeKey{add(hookNode("after_sign", hook, config, request, PackageScope, packageRoots))}
		}
	}
	plan.Roots = packageRoots
	return plan, plan.Validate(config.Root)
}

func platformAssetsNode(config manifest.Config, request Request, target, output string, dependencies []NodeKey) Node {
	project, targetSettings := effectiveTarget(config, request.TargetOS, request.TargetArch)
	return Node{Key: NodeKey("target:" + target + ":assets"), Kind: GeneratePlatformAssets, Label: "Generate " + request.TargetOS + " assets", Scope: TargetScope, Dependencies: dependencies,
		Spec:   AssetsSpec{TargetOS: request.TargetOS, TargetArch: request.TargetArch, Directory: output, MinimumVersion: targetSettings.MinimumVersion, Project: project, Associations: config.Associations, Protocols: config.Protocols},
		Inputs: assetInputs(config), Output: output, Cache: CacheArtifact, Claims: ResourceClaims{CPU: 1, MemoryMB: 512, Exclusive: "legacy-command-adapter"}, EstimateMS: 250}
}

func hookNode(phase string, hook manifest.Hook, config manifest.Config, request Request, scope Scope, deps []NodeKey) Node {
	policy := CacheNever
	output := ""
	if hook.Cache {
		policy = CacheArtifact
		output = commonOutput(hook.Outputs)
	}
	inputs := []InputSpec{{Label: "hook-script", Files: append([]string{hook.Script}, hook.Inputs...)}}
	return Node{Key: NodeKey("hook:" + phase + ":" + request.TargetOS + "/" + request.TargetArch), Kind: RunHook, Label: strings.ReplaceAll(phase, "_", " "), Scope: scope, Dependencies: deps, Spec: HookSpec{Phase: phase, TargetOS: request.TargetOS, TargetArch: request.TargetArch, Profile: config.Profile, Output: output, Hook: hook}, Inputs: inputs, Output: output, Cache: policy, Claims: ResourceClaims{CPU: 1, MemoryMB: 256, Exclusive: "hook:" + phase}, EstimateMS: 100}
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
	return []InputSpec{{Label: "platform-assets", Files: files}}
}
func templateInput(f manifest.PackageFormat) []InputSpec {
	if f.Template == "" {
		return nil
	}
	return []InputSpec{{Label: "package-template", Files: []string{f.Template}}}
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
func supportsCustomTemplate(format string) bool {
	return format == "nsis" || format == "deb" || format == "rpm" || format == "archlinux" || format == "apk" || format == "aab"
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
func packageOutput(c manifest.Config, os, arch, format string) string {
	base := filepath.Join(c.Build.OutputDirectory, c.Project.BinaryName)
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
func packageExclusive(_ string) string {
	// Existing platform packagers share process-wide working-directory and
	// output controls. Keep branches independent in the graph, but do not
	// overlap these adapters until those legacy globals are removed.
	return "legacy-command-adapter"
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
