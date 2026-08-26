package manifest

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"unicode"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Loaded struct {
	Path     string
	Raw      []byte
	Document Document
	Config   Config
}

func Exists(root string) bool {
	_, _, err := Discover(root)
	return err == nil
}

// Discover searches upward for the nearest manifest. The directory containing
// it is the project root; callers never need to infer a root independently.
func Discover(start string) (root, path string, err error) {
	return discoverWithOperations(start, filepath.Abs, os.Stat)
}

func discoverWithOperations(start string, absolute func(string) (string, error), stat func(string) (fs.FileInfo, error)) (root, path string, err error) {
	current, err := absolute(start)
	if err != nil {
		return "", "", err
	}
	if info, statErr := stat(current); statErr == nil && !info.IsDir() {
		current = filepath.Dir(current)
	}
	for {
		candidate := filepath.Join(current, Filename)
		if info, statErr := stat(candidate); statErr == nil {
			if info.IsDir() {
				return "", "", fmt.Errorf("%s is a directory", candidate)
			}
			return current, candidate, nil
		} else if !errors.Is(statErr, fs.ErrNotExist) {
			return "", "", statErr
		}
		module := filepath.Join(current, "go.mod")
		if info, statErr := stat(module); statErr == nil && !info.IsDir() {
			return "", "", fs.ErrNotExist
		} else if statErr != nil && !errors.Is(statErr, fs.ErrNotExist) {
			return "", "", statErr
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", "", fs.ErrNotExist
		}
		current = parent
	}
}

func Load(start, profile string) (*Loaded, error) {
	return loadWithOperations(start, profile, Discover, os.ReadFile)
}

func loadWithOperations(start, profile string, discover func(string) (string, string, error), readFile func(string) ([]byte, error)) (*Loaded, error) {
	root, path, err := discover(start)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("could not find %s from %s", Filename, start)
		}
		return nil, err
	}
	raw, err := readFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", Filename, err)
	}
	return decodeHCL(root, path, raw, profile)
}

func defaults(project Project) Document {
	if project.BinaryName == "" {
		project.BinaryName = deriveBinaryName(project.Name)
	}
	if project.BuildNumber == 0 {
		project.BuildNumber = 1
	}
	frontendRaw := &hclFrontend{}
	applySchemaDefaults(frontendRaw)
	bindingsRaw := &hclBindings{}
	applySchemaDefaults(bindingsRaw)
	frontendRaw.Bindings = bindingsRaw
	frontend := Frontend{InstallCommand: "install", BuildCommand: "build", DevCommand: "dev"}
	applyFrontend(&frontend, frontendRaw)
	buildRaw := &hclBuild{}
	applySchemaDefaults(buildRaw)
	build := Build{Production: true}
	applyBuild(&build, buildRaw)
	devRaw := &hclDev{}
	applySchemaDefaults(devRaw)
	dev := Dev{Port: 9245}
	applyDev(&dev, devRaw)
	targets := Targets{
		Windows: defaultPlatform("amd64", "arm64"), Darwin: defaultPlatform("amd64", "arm64"),
		Linux: defaultPlatform("amd64", "arm64"), IOS: defaultPlatform("arm64"), Android: defaultPlatform("arm64"),
	}
	targets.IOS.MinimumVersion = "15.0"
	targets.Android.MinimumSDK = 21
	packages := Packages{
		Windows: PackagePlatform{Formats: []string{"nsis"}}, Darwin: PackagePlatform{Formats: []string{"dmg"}},
		Linux: PackagePlatform{Formats: []string{"appimage"}}, IOS: PackagePlatform{Formats: []string{"ipa"}},
		Android: PackagePlatform{Formats: []string{"aab"}},
	}
	return Document{Project: project, Frontend: frontend, Build: build, Dev: dev, Targets: targets, Package: packages, Hooks: map[HookPhase]Hook{}, Profiles: map[string]Profile{}}
}

// NewDocument returns a resolved document seeded with the compiled defaults.
// Programmatic writers should start here so setting a default-true field to
// false remains distinguishable from leaving an optional section unset.
func NewDocument(project Project) Document { return defaults(project) }

func defaultPlatform(architectures ...string) Platform {
	result := Platform{}
	for _, arch := range architectures {
		switch arch {
		case "amd64":
			result.AMD64.Enabled = true
		case "arm64":
			result.ARM64.Enabled = true
		case "arm":
			result.ARM.Enabled = true
		case "386":
			result.X86.Enabled = true
		}
	}
	return result
}

func configFromDocument(root, profile string, doc Document) Config {
	config := Config{Root: root, Profile: profile, Project: doc.Project, Frontend: doc.Frontend, Build: doc.Build, Dev: doc.Dev, Targets: doc.Targets, Package: doc.Package, Signing: doc.Signing, Associations: doc.Associations, Protocols: doc.Protocols, Hooks: doc.Hooks, Profiles: doc.Profiles, Origins: defaultOrigins()}
	if profile != "" {
		config.Selected = doc.Profiles[profile]
	}
	return config
}

func validateProject(project Project) error {
	if project.Name == "" || project.ProductName == "" || project.Identifier == "" || project.Version == "" {
		return fieldValidationError("project", "requires name, product_name, identifier, and version")
	}
	return nil
}

func validateConfig(config Config) error {
	paths, err := newProjectPathValidator(config.Root)
	if err != nil {
		return fmt.Errorf("resolve project root: %w", err)
	}
	if config.Project.BinaryName == "" {
		return fieldValidationError("project.binary_name", "project name does not produce a value")
	}
	if config.Project.BinaryName != filepath.Base(config.Project.BinaryName) || strings.ContainsAny(config.Project.BinaryName, `/\\`) || config.Project.BinaryName == "." || config.Project.BinaryName == ".." {
		return fieldValidationError("project.binary_name", "must be a plain file name, got %q", config.Project.BinaryName)
	}
	if !contains([]string{"npm", "pnpm", "yarn", "bun"}, config.Frontend.PackageManager) {
		return fieldValidationError("frontend.install", "unsupported package manager %q", config.Frontend.PackageManager)
	}
	for field, value := range map[string]string{
		"frontend.directory": config.Frontend.Directory,
		"frontend.output":    config.Frontend.OutputDirectory,
		"build.output":       config.Build.OutputDirectory,
	} {
		if value == "" {
			return fieldValidationError(field, "cannot be empty")
		}
	}
	for field, value := range map[string]string{
		"frontend.install": config.Frontend.InstallCommand,
		"frontend.build":   config.Frontend.BuildCommand,
		"frontend.dev":     config.Frontend.DevCommand,
	} {
		if value == "" {
			return fieldValidationError(field, "command cannot be empty")
		}
	}
	if !contains([]string{"string", "Date"}, config.Frontend.Bindings.TimeType) {
		return fieldValidationError("frontend.bindings.time_type", "must be either string or Date")
	}
	if config.Frontend.Bindings.Interfaces && config.Frontend.Bindings.TimeType == "Date" {
		return fieldValidationError("frontend.bindings.time_type", "Date is not supported when interfaces is true")
	}
	for index, pattern := range config.Dev.Watch {
		if err := validateDevWatchPattern(pattern); err != nil {
			return fieldValidationError("dev.watch", "item %d %q: %v", index, pattern, err)
		}
	}
	for name, value := range map[string]string{
		"project.icon":             config.Project.Icon,
		"frontend.directory":       config.Frontend.Directory,
		"frontend.output":          config.Frontend.OutputDirectory,
		"frontend.bindings.output": config.Frontend.Bindings.OutputDirectory,
		"build.output":             config.Build.OutputDirectory,
	} {
		if err := paths.validate(name, value, false); err != nil {
			return err
		}
	}
	for _, platformEntry := range []struct {
		name  string
		value Platform
	}{
		{"windows", config.Targets.Windows}, {"darwin", config.Targets.Darwin}, {"linux", config.Targets.Linux},
		{"ios", config.Targets.IOS}, {"android", config.Targets.Android},
	} {
		name, platform := platformEntry.name, platformEntry.value
		for _, targetEntry := range []struct {
			arch  string
			value Target
		}{
			{"amd64", platform.AMD64}, {"arm64", platform.ARM64}, {"arm", platform.ARM},
			{"386", platform.X86}, {"universal", platform.Universal},
		} {
			arch, target := targetEntry.arch, targetEntry.value
			if target.Toolchain != "" && !contains([]string{"auto", "native", "zig", "docker"}, target.Toolchain) {
				return fieldValidationError(fmt.Sprintf(`target[%q].toolchain`, name+"/"+arch), "unsupported toolchain %q", target.Toolchain)
			}
			if err := validateEnvironment(fmt.Sprintf(`target[%q].environment`, name+"/"+arch), target.Environment); err != nil {
				return err
			}
		}
		for field, value := range map[string]string{"icon": platform.Icon, "manifest": platform.Manifest, "assets_car": platform.AssetsCar, "info_plist": platform.InfoPlist, "desktop_entry": platform.DesktopEntry} {
			if err := paths.validate(name+"."+field, value, false); err != nil {
				return err
			}
		}
	}
	if err := validateEnvironment("frontend.environment", config.Frontend.Environment); err != nil {
		return err
	}
	if err := validateEnvironment("build.environment", config.Build.Environment); err != nil {
		return err
	}
	if err := validateHooks(config.Root, paths, config.Hooks); err != nil {
		return err
	}
	for _, mode := range config.Targets.IOS.BackgroundModes {
		if !contains([]string{"audio", "location", "voip", "fetch", "remote-notification", "newsstand-content", "external-accessory", "bluetooth-central", "bluetooth-peripheral", "network-authentication", "processing"}, mode) {
			return fieldValidationError("ios.background_modes", "contains unsupported mode %q", mode)
		}
	}
	for _, signingEntry := range []struct {
		name  string
		value SigningPlatform
	}{
		{"windows", config.Signing.Windows}, {"darwin", config.Signing.Darwin}, {"linux", config.Signing.Linux},
		{"ios", config.Signing.IOS}, {"android", config.Signing.Android},
	} {
		name, signing := signingEntry.name, signingEntry.value
		for field, value := range map[string]string{"certificate": signing.Certificate, "entitlements": signing.Entitlements, "provisioning_profile": signing.ProvisioningProfile} {
			if err := paths.validate(name+".signing."+field, value, false); err != nil {
				return err
			}
		}
	}
	for _, association := range config.Associations {
		field := fmt.Sprintf(`file_association[%q]`, association.Name)
		if len(association.Extensions) == 0 {
			return fieldValidationError(field+".extensions", "requires at least one extension")
		}
		for _, extension := range association.Extensions {
			if strings.TrimSpace(strings.TrimPrefix(extension, ".")) == "" {
				return fieldValidationError(field+".extensions", "contains an empty extension")
			}
		}
		if err := paths.validate(field+".icon", association.Icon, false); err != nil {
			return err
		}
		if err := validateRegistrationPlatforms(field, association.Platforms); err != nil {
			return err
		}
	}
	for _, protocol := range config.Protocols {
		field := fmt.Sprintf(`protocol[%q]`, protocol.Scheme)
		if strings.TrimSpace(protocol.Scheme) == "" {
			return fieldValidationError(field, "scheme is required")
		}
		if err := validateRegistrationPlatforms(field, protocol.Platforms); err != nil {
			return err
		}
	}
	for _, formatEntry := range []struct{ platform, format string }{
		{"windows", "nsis"}, {"windows", "msix"}, {"darwin", "app"}, {"darwin", "dmg"},
		{"linux", "appimage"}, {"linux", "deb"}, {"linux", "rpm"}, {"linux", "archlinux"},
		{"ios", "app"}, {"ios", "ipa"}, {"android", "apk"}, {"android", "aab"},
	} {
		name := formatEntry.platform + "." + formatEntry.format
		format, _ := ResolvePackageFormat(config.Package, formatEntry.platform, formatEntry.format)
		if err := validatePackageOptions(paths, name, format); err != nil {
			return err
		}
	}
	return nil
}

func validateHooks(root string, paths *projectPathValidator, hooks map[HookPhase]Hook) error {
	phases := make([]string, 0, len(hooks))
	for phase := range hooks {
		phases = append(phases, string(phase))
	}
	sort.Strings(phases)
	for _, phaseName := range phases {
		phase := HookPhase(phaseName)
		hook := hooks[phase]
		field := fmt.Sprintf(`hook[%q]`, phase)
		if !containsHookPhase(phase) {
			return fieldValidationError(field, "unsupported hook phase")
		}
		if hook.Script == "" {
			return fieldValidationError(field+".script", "cannot be empty")
		}
		if wailsManagedPath(hook.Script) {
			return fieldValidationError(field+".script", "must not be inside Wails-owned .wails state")
		}
		if err := paths.validate(field+".script", hook.Script, true); err != nil {
			return err
		}
		scriptPath, err := ResolveProjectPath(root, field+".script", hook.Script, true)
		if err != nil {
			return err
		}
		info, err := os.Stat(scriptPath)
		if err != nil {
			return fieldValidationError(field+".script", "%v", err)
		}
		if !info.Mode().IsRegular() {
			return fieldValidationError(field+".script", "must be a regular file")
		}
		if runtime.GOOS == "windows" {
			if !validWindowsHookScript(hook.Script) {
				return fieldValidationError(field+".script", "must use a .cmd, .bat, or .ps1 script on Windows")
			}
		} else {
			if info.Mode().Perm()&0o111 == 0 {
				return fieldValidationError(field+".script", "must be executable")
			}
			file, openErr := os.Open(scriptPath)
			if openErr != nil {
				return fieldValidationError(field+".script", "%v", openErr)
			}
			var prefix [2]byte
			_, readErr := io.ReadFull(file, prefix[:])
			closeErr := file.Close()
			if readErr != nil {
				return fieldValidationError(field+".script", "%v", readErr)
			}
			if closeErr != nil {
				return fieldValidationError(field+".script", "%v", closeErr)
			}
			if string(prefix[:]) != "#!" {
				return fieldValidationError(field+".script", "must start with a shebang")
			}
		}
		if hook.Directory != "" {
			if wailsManagedPath(hook.Directory) {
				return fieldValidationError(field+".directory", "must not be inside Wails-owned .wails state")
			}
			if err := paths.validate(field+".directory", hook.Directory, true); err != nil {
				return err
			}
			directory, err := ResolveProjectPath(root, field+".directory", hook.Directory, true)
			if err != nil {
				return err
			}
			info, err := os.Stat(directory)
			if err != nil || !info.IsDir() {
				return fieldValidationError(field+".directory", "must be an existing directory")
			}
		}
		if !hook.Cache && (len(hook.Inputs) != 0 || len(hook.Outputs) != 0) {
			return fieldValidationError(field+".cache", "must be true when inputs or outputs are declared")
		}
		if hook.Cache && (len(hook.Inputs) == 0 || len(hook.Outputs) == 0) {
			return fieldValidationError(field+".cache", "requires complete inputs and outputs")
		}
		for index, input := range hook.Inputs {
			if input == "" || wailsManagedPath(input) {
				return fieldValidationError(fmt.Sprintf("%s.inputs[%d]", field, index), "must be a non-empty user-owned project path")
			}
			if err := paths.validate(fmt.Sprintf("%s.inputs[%d]", field, index), input, false); err != nil {
				return err
			}
		}
		for index, output := range hook.Outputs {
			if output == "" || wailsManagedPath(output) {
				return fieldValidationError(fmt.Sprintf("%s.outputs[%d]", field, index), "must be a non-empty user-owned project path")
			}
			if err := paths.validate(fmt.Sprintf("%s.outputs[%d]", field, index), output, false); err != nil {
				return err
			}
		}
		if hook.Cache {
			outputRoot, err := HookOutputRoot(hook.Outputs)
			if err != nil {
				return fieldValidationError(field+".outputs", "%v", err)
			}
			for _, input := range append([]string{hook.Script}, hook.Inputs...) {
				if pathContains(outputRoot, input) {
					return fieldValidationError(field+".outputs", "output root %q contains script or input %q", outputRoot, input)
				}
			}
		}
	}
	return nil
}

func validWindowsHookScript(value string) bool {
	switch strings.ToLower(filepath.Ext(value)) {
	case ".cmd", ".bat", ".ps1":
		return true
	default:
		return false
	}
}

// HookOutputRoot returns the bounded artifact restored for a cacheable hook.
func HookOutputRoot(outputs []string) (string, error) {
	if len(outputs) == 0 {
		return "", fmt.Errorf("at least one output is required")
	}
	clean := filepath.ToSlash(filepath.Clean(outputs[0]))
	if clean == "." || clean == "" {
		return "", fmt.Errorf("output root must not be the project root")
	}
	if len(outputs) == 1 {
		return clean, nil
	}
	common := strings.Split(clean, "/")
	for _, output := range outputs[1:] {
		parts := strings.Split(filepath.ToSlash(filepath.Clean(output)), "/")
		limit := min(len(common), len(parts))
		index := 0
		for index < limit && common[index] == parts[index] {
			index++
		}
		common = common[:index]
	}
	if len(common) == 0 {
		return "", fmt.Errorf("multiple outputs must share one non-root directory")
	}
	return strings.Join(common, "/"), nil
}

func pathContains(root, pathValue string) bool {
	root = filepath.ToSlash(filepath.Clean(root))
	pathValue = filepath.ToSlash(filepath.Clean(pathValue))
	return pathValue == root || strings.HasPrefix(pathValue, root+"/")
}

func wailsManagedPath(value string) bool {
	clean := strings.TrimPrefix(filepath.ToSlash(filepath.Clean(value)), "./")
	if runtime.GOOS == "windows" {
		clean = strings.ToLower(clean)
	}
	return clean == ".wails" || strings.HasPrefix(clean, ".wails/")
}

func validateDevWatchPattern(pattern string) error {
	pattern = strings.Trim(strings.TrimPrefix(filepath.ToSlash(pattern), "./"), "/")
	if pattern == "" {
		return fmt.Errorf("pattern cannot be empty")
	}
	for _, segment := range strings.Split(pattern, "/") {
		if segment == "**" {
			continue
		}
		if _, err := path.Match(segment, ""); err != nil {
			return fmt.Errorf("invalid pattern: %w", err)
		}
	}
	return nil
}

func validateProjectPath(root, field, value string, requireExisting bool) error {
	_, err := ResolveProjectPath(root, field, value, requireExisting)
	return err
}

// ResolveProjectPath applies the manifest's single path-safety contract and
// returns a path suitable for immediate use. Existing inputs are symlink-
// resolved; prospective outputs retain their project-relative spelling after
// every existing ancestor has been checked.
func ResolveProjectPath(root, field, value string, requireExisting bool) (string, error) {
	validator, err := newProjectPathValidator(root)
	if err != nil {
		return "", fmt.Errorf("%s %q: %w", field, value, err)
	}
	if err := validator.validate(field, value, requireExisting); err != nil {
		return "", err
	}
	candidate := filepath.Clean(filepath.Join(root, filepath.FromSlash(strings.ReplaceAll(value, `\`, "/"))))
	if !requireExisting || value == "" {
		return candidate, nil
	}
	// validateResolved populated this exact key while proving the existing path
	// stays within the project. Reuse that result instead of resolving twice.
	return validator.resolved[candidate].path, nil
}

type projectPathValidator struct {
	root, resolvedRoot string
	resolved           map[string]pathResolution
	ops                projectPathOperations
}

type pathResolution struct {
	path   string
	err    error
	exists bool
}

type projectPathOperations struct {
	eval  func(string) (string, error)
	lstat func(string) (fs.FileInfo, error)
	rel   func(string, string) (string, error)
}

func newProjectPathValidator(root string) (*projectPathValidator, error) {
	return newProjectPathValidatorWithOperations(root, projectPathOperations{eval: filepath.EvalSymlinks, lstat: os.Lstat, rel: filepath.Rel})
}

func newProjectPathValidatorWithOperations(root string, ops projectPathOperations) (*projectPathValidator, error) {
	resolvedRoot, err := ops.eval(root)
	if err != nil {
		return nil, err
	}
	return &projectPathValidator{
		root: root, resolvedRoot: resolvedRoot,
		resolved: map[string]pathResolution{filepath.Clean(root): {path: resolvedRoot, exists: true}},
		ops:      ops,
	}, nil
}

func (v *projectPathValidator) validate(field, value string, requireExisting bool) error {
	if value == "" {
		return nil
	}
	if manifestPathIsAbsolute(value) || pathEscapes(value) {
		return fieldValidationError(field, "must be project-relative, got %q", value)
	}
	if pathEntersWails(value) {
		return fieldValidationError(field, "must not reference .wails, got %q", value)
	}
	if err := v.validateResolved(value, requireExisting); err != nil {
		return fieldValidationCause(field, err, "%q: %v", value, err)
	}
	return nil
}

func validateResolvedProjectPath(root, value string) error {
	validator, err := newProjectPathValidator(root)
	if err != nil {
		return err
	}
	return validator.validateResolved(value, true)
}

func (v *projectPathValidator) validateResolved(value string, requireExisting bool) error {
	candidate := filepath.Join(v.root, filepath.FromSlash(strings.ReplaceAll(value, `\`, "/")))
	resolvedPath, err := v.resolve(candidate, requireExisting)
	if err != nil {
		return err
	}
	relative, err := v.ops.rel(v.resolvedRoot, resolvedPath)
	if err != nil {
		return err
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return fmt.Errorf("resolves outside the project")
	}
	return nil
}

func (v *projectPathValidator) resolve(candidate string, requireExisting bool) (string, error) {
	candidate = filepath.Clean(candidate)
	if cached, ok := v.resolved[candidate]; ok {
		if requireExisting && !cached.exists && cached.err == nil {
			return "", fs.ErrNotExist
		}
		return cached.path, cached.err
	}
	if _, err := v.ops.lstat(candidate); err == nil {
		resolved, resolveErr := v.ops.eval(candidate)
		v.resolved[candidate] = pathResolution{path: resolved, err: resolveErr, exists: true}
		return resolved, resolveErr
	} else if !errors.Is(err, fs.ErrNotExist) || requireExisting {
		v.resolved[candidate] = pathResolution{err: err}
		return "", err
	}
	parent := filepath.Dir(candidate)
	if parent == candidate {
		return "", fs.ErrNotExist
	}
	resolved, err := v.resolve(parent, false)
	v.resolved[candidate] = pathResolution{path: resolved, err: err}
	return resolved, err
}

func pathEntersWails(value string) bool {
	normalized := filepath.ToSlash(filepath.Clean(strings.ReplaceAll(value, `\`, "/")))
	for _, segment := range strings.Split(normalized, "/") {
		if strings.EqualFold(segment, ".wails") {
			return true
		}
	}
	return false
}

func validateRegistrationPlatforms(field string, platforms []string) error {
	for _, platform := range platforms {
		if !contains([]string{"windows", "darwin", "linux", "ios", "android"}, platform) {
			return fieldValidationError(field+".platforms", "contains unsupported platform %q", platform)
		}
	}
	return nil
}

func validatePackageOptions(paths *projectPathValidator, name string, format PackageFormat) error {
	manifestField := packageManifestField(name)
	if format.Format != "" {
		_, expected, _ := strings.Cut(name, ".")
		if format.Format != expected {
			return fieldValidationError(manifestField, "format identity is %q, expected %q", format.Format, expected)
		}
	}
	if format.Template != "" && packageHasStructuredConfiguration(format) {
		return fieldValidationError(manifestField, "cannot combine a complete template replacement with structured options")
	}
	node := schemaNodesByType[reflect.TypeOf(PackageFormat{})]
	value := reflect.ValueOf(format)
	for _, attributeName := range node.attributeOrder {
		descriptor := node.attributes[attributeName]
		attribute := value.Field(descriptor.fieldIndex)
		if attribute.IsZero() {
			continue
		}
		if !schemaFormatAllowed(manifestField, descriptor.formatMask) {
			return fieldValidationError(manifestField+"."+attributeName, "field is not supported for this package format")
		}
		if !descriptor.path {
			continue
		}
		field := manifestField + "." + attributeName
		switch attribute.Kind() {
		case reflect.String:
			if err := paths.validate(field, attribute.String(), attributeName == "template"); err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					return fieldValidationError(field, "%q does not exist", attribute.String())
				}
				return err
			}
		case reflect.Map:
			keys := attribute.MapKeys()
			sort.Slice(keys, func(left, right int) bool { return keys[left].String() < keys[right].String() })
			for _, key := range keys {
				if strings.TrimSpace(key.String()) == "" {
					return fieldValidationError(field, "file name must not be empty")
				}
				if err := paths.validate(field, attribute.MapIndex(key).String(), false); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func packageHasStructuredConfiguration(format PackageFormat) bool {
	copy := format
	copy.Format = ""
	copy.Template = ""
	return !reflect.ValueOf(copy).IsZero()
}

func validateEnvironment(field string, environment map[string]string) error {
	names := make([]string, 0, len(environment))
	for name := range environment {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if strings.TrimSpace(name) == "" || strings.Contains(name, "=") {
			return fieldValidationError(field, "contains invalid variable name %q", name)
		}
	}
	return nil
}

func packageManifestField(name string) string {
	_, format, _ := strings.Cut(name, ".")
	return fmt.Sprintf(`package[%q]`, format)
}

func manifestPathIsAbsolute(value string) bool {
	normalized := strings.ReplaceAll(value, `\`, "/")
	return filepath.IsAbs(normalized) || strings.HasPrefix(normalized, "//") || len(normalized) >= 3 && normalized[1] == ':' && normalized[2] == '/'
}

func pathEscapes(path string) bool {
	normalized := strings.ReplaceAll(path, `\`, "/")
	clean := filepath.ToSlash(filepath.Clean(normalized))
	return clean == ".." || strings.HasPrefix(clean, "../")
}
func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func deriveBinaryName(name string) string {
	var out strings.Builder
	dash := false
	for _, r := range strings.ToLower(strings.TrimSpace(name)) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			if dash && out.Len() > 0 {
				out.WriteByte('-')
			}
			dash = false
			out.WriteRune(r)
		default:
			dash = true
		}
	}
	return strings.Trim(out.String(), "-")
}
