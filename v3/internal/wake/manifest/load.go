package manifest

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
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
	current, err := filepath.Abs(start)
	if err != nil {
		return "", "", err
	}
	if info, statErr := os.Stat(current); statErr == nil && !info.IsDir() {
		current = filepath.Dir(current)
	}
	for {
		candidate := filepath.Join(current, Filename)
		if info, statErr := os.Stat(candidate); statErr == nil {
			if info.IsDir() {
				return "", "", fmt.Errorf("%s is a directory", candidate)
			}
			return current, candidate, nil
		} else if !errors.Is(statErr, fs.ErrNotExist) {
			return "", "", statErr
		}
		module := filepath.Join(current, "go.mod")
		if info, statErr := os.Stat(module); statErr == nil && !info.IsDir() {
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
	root, path, err := Discover(start)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("could not find %s from %s", Filename, start)
		}
		return nil, err
	}
	raw, err := os.ReadFile(path)
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
	frontend := Frontend{
		Directory: "frontend", PackageManager: "npm", InstallCommand: "install",
		BuildCommand: "build", DevCommand: "dev", OutputDirectory: "dist",
		Install: []string{"npm", "install"}, Build: []string{"npm", "run", "build"}, Dev: []string{"npm", "run", "dev"},
		Bindings: Bindings{TypeScript: true, Interfaces: true, OutputDirectory: "bindings", ModelsFilename: "models", IndexFilename: "index", TimeType: "string"},
	}
	build := Build{OutputDirectory: "bin", Production: true, TrimPath: true, Strip: true}
	dev := Dev{Port: 9245, DebounceMS: 250, LogLevel: "warn", Watch: []string{"**/*.go", Filename}, Exclude: []string{".git", ".wails", "bin", "node_modules", "frontend/dist"}, UseGitIgnore: true, GracePeriodMS: 1500}
	targets := Targets{
		Windows: defaultPlatform("amd64", "arm64"), Darwin: defaultPlatform("amd64", "arm64"),
		Linux: defaultPlatform("amd64", "arm64"), IOS: defaultPlatform("arm64"), Android: defaultPlatform("arm64"),
	}
	targets.IOS.MinimumVersion = "15.0"
	targets.IOS.ARM64.Variant = "simulator"
	targets.Android.MinimumVersion = "21"
	packages := Packages{
		Windows: PackagePlatform{Formats: []string{"nsis"}}, Darwin: PackagePlatform{Formats: []string{"app"}},
		Linux: PackagePlatform{Formats: []string{"appimage"}}, IOS: PackagePlatform{Formats: []string{"app"}},
		Android: PackagePlatform{Formats: []string{"apk"}},
	}
	return Document{Project: project, Frontend: frontend, Build: build, Dev: dev, Targets: targets, Package: packages, Profiles: map[string]Profile{}, Extensions: map[string]map[string]any{}}
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
	config := Config{Root: root, Profile: profile, Project: doc.Project, Frontend: doc.Frontend, Build: doc.Build, Dev: doc.Dev, Targets: doc.Targets, Package: doc.Package, Signing: doc.Signing, Associations: doc.Associations, Protocols: doc.Protocols, Hooks: doc.Hooks, Wake: doc.Wake, Profiles: doc.Profiles, Extensions: doc.Extensions}
	if profile != "" {
		config.Selected = doc.Profiles[profile]
	}
	return config
}

func validateProject(project Project) error {
	if project.Name == "" || project.ProductName == "" || project.Identifier == "" || project.Version == "" {
		return fmt.Errorf("%s: [project] requires name, product_name, identifier, and version", Filename)
	}
	return nil
}

func validateConfig(config Config) error {
	if config.Project.BinaryName == "" {
		return fmt.Errorf("project name does not produce a binary_name")
	}
	if config.Project.BinaryName != filepath.Base(config.Project.BinaryName) || strings.ContainsAny(config.Project.BinaryName, `/\\`) || config.Project.BinaryName == "." || config.Project.BinaryName == ".." {
		return fmt.Errorf("binary_name must be a plain file name, got %q", config.Project.BinaryName)
	}
	if !contains([]string{"npm", "pnpm", "yarn", "bun"}, config.Frontend.PackageManager) {
		return fmt.Errorf("unsupported frontend package_manager %q", config.Frontend.PackageManager)
	}
	if config.Frontend.Directory == "" || config.Frontend.OutputDirectory == "" || config.Build.OutputDirectory == "" {
		return fmt.Errorf("frontend and build directories cannot be empty")
	}
	if config.Frontend.InstallCommand == "" || config.Frontend.BuildCommand == "" || config.Frontend.DevCommand == "" {
		return fmt.Errorf("frontend install_command, build_command, and dev_command cannot be empty")
	}
	for index, pattern := range config.Dev.Watch {
		if err := validateDevWatchPattern(pattern); err != nil {
			return fmt.Errorf("dev.watch[%d] %q: %w", index, pattern, err)
		}
	}
	for name, value := range map[string]string{
		"project.icon":                       config.Project.Icon,
		"frontend.directory":                 config.Frontend.Directory,
		"frontend.output_directory":          config.Frontend.OutputDirectory,
		"frontend.bindings.output_directory": config.Frontend.Bindings.OutputDirectory,
		"build.output_directory":             config.Build.OutputDirectory,
	} {
		if value != "" && (manifestPathIsAbsolute(value) || pathEscapes(value)) {
			return fmt.Errorf("%s must be a project-relative path", name)
		}
	}
	for name, hook := range hookMap(config.Hooks) {
		if hook.Script == "" {
			continue
		}
		if manifestPathIsAbsolute(hook.Script) || pathEscapes(hook.Script) {
			return fmt.Errorf("hook %s script must be relative to %s", name, Filename)
		}
		if hook.Directory != "" && (manifestPathIsAbsolute(hook.Directory) || pathEscapes(hook.Directory)) {
			return fmt.Errorf("hook %s directory must be relative to %s", name, Filename)
		}
		for _, value := range append(append([]string(nil), hook.Inputs...), hook.Outputs...) {
			if manifestPathIsAbsolute(value) || pathEscapes(value) {
				return fmt.Errorf("hook %s input and output paths must be project-relative", name)
			}
		}
		if hook.Cache && (len(hook.Inputs) == 0 || len(hook.Outputs) == 0) {
			return fmt.Errorf("cached hook %s requires inputs and outputs", name)
		}
		if !hook.Cache && (len(hook.Inputs) > 0 || len(hook.Outputs) > 0) {
			return fmt.Errorf("hook %s inputs and outputs require cache = true", name)
		}
		if hook.Cache {
			outputRoot, err := HookOutputRoot(hook.Outputs)
			if err != nil {
				return fmt.Errorf("cached hook %s outputs: %w", name, err)
			}
			for _, input := range append([]string{hook.Script}, hook.Inputs...) {
				if pathContains(outputRoot, input) {
					return fmt.Errorf("cached hook %s output root %q contains input %q", name, outputRoot, input)
				}
			}
		}
	}
	for name, signing := range signingMap(config.Signing) {
		for field, value := range map[string]string{"certificate": signing.Certificate, "entitlements": signing.Entitlements, "provisioning_profile": signing.ProvisioningProfile} {
			if value != "" && (manifestPathIsAbsolute(value) || pathEscapes(value)) {
				return fmt.Errorf("signing.%s.%s must be project-relative", name, field)
			}
		}
	}
	for index, association := range config.Associations {
		if len(association.Extensions) == 0 {
			return fmt.Errorf("associations[%d].extensions requires at least one extension", index)
		}
		for _, extension := range association.Extensions {
			if strings.TrimSpace(strings.TrimPrefix(extension, ".")) == "" {
				return fmt.Errorf("associations[%d].extensions contains an empty extension", index)
			}
		}
		if association.Icon != "" && (manifestPathIsAbsolute(association.Icon) || pathEscapes(association.Icon)) {
			return fmt.Errorf("associations[%d].icon must be project-relative", index)
		}
		if err := validateRegistrationPlatforms(fmt.Sprintf("associations[%d]", index), association.Platforms); err != nil {
			return err
		}
	}
	for index, protocol := range config.Protocols {
		if strings.TrimSpace(protocol.Scheme) == "" {
			return fmt.Errorf("protocols[%d].scheme is required", index)
		}
		if err := validateRegistrationPlatforms(fmt.Sprintf("protocols[%d]", index), protocol.Platforms); err != nil {
			return err
		}
	}
	for name, format := range packageFormatMap(config.Package) {
		if format.Template != "" {
			if manifestPathIsAbsolute(format.Template) || pathEscapes(format.Template) {
				return fmt.Errorf("package.%s.template must be project-relative", name)
			}
			if err := validateResolvedProjectPath(config.Root, format.Template); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("package.%s.template %q does not exist", name, format.Template)
				}
				return fmt.Errorf("package.%s.template %q: %w", name, format.Template, err)
			}
		}
		if err := validatePackageOptions(name, format); err != nil {
			return err
		}
	}
	return nil
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

func validateResolvedProjectPath(root, value string) error {
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return err
	}
	resolvedPath, err := filepath.EvalSymlinks(filepath.Join(root, filepath.FromSlash(value)))
	if err != nil {
		return err
	}
	relative, err := filepath.Rel(resolvedRoot, resolvedPath)
	if err != nil {
		return err
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return fmt.Errorf("must resolve inside the project")
	}
	return nil
}

func pathContains(parent, child string) bool {
	parent = filepath.ToSlash(filepath.Clean(filepath.FromSlash(parent)))
	child = filepath.ToSlash(filepath.Clean(filepath.FromSlash(child)))
	return child == parent || strings.HasPrefix(child, parent+"/")
}

func validateRegistrationPlatforms(field string, platforms []string) error {
	for _, platform := range platforms {
		if !contains([]string{"windows", "darwin", "linux", "ios", "android"}, platform) {
			return fmt.Errorf("%s.platforms contains unsupported platform %q", field, platform)
		}
	}
	return nil
}

// HookOutputRoot returns the single file/directory that owns a cacheable
// hook's declared outputs. Multiple outputs must share a non-project-root
// ancestor so recording the Artifact cannot capture unrelated project files.
func HookOutputRoot(outputs []string) (string, error) {
	if len(outputs) == 0 {
		return "", fmt.Errorf("at least one output is required")
	}
	cleaned := make([]string, len(outputs))
	for index, output := range outputs {
		cleaned[index] = filepath.ToSlash(filepath.Clean(filepath.FromSlash(output)))
		if cleaned[index] == "." {
			return "", fmt.Errorf("project root cannot be a cached output")
		}
	}
	if len(cleaned) == 1 {
		return cleaned[0], nil
	}
	common := strings.Split(cleaned[0], "/")
	for _, output := range cleaned[1:] {
		parts := strings.Split(output, "/")
		limit := len(common)
		if len(parts) < limit {
			limit = len(parts)
		}
		matched := 0
		for matched < limit && common[matched] == parts[matched] {
			matched++
		}
		common = common[:matched]
	}
	if len(common) == 0 {
		return "", fmt.Errorf("multiple outputs must share a top-level directory")
	}
	return strings.Join(common, "/"), nil
}

func validatePackageOptions(name string, format PackageFormat) error {
	if len(format.Options) == 0 || format.Template != "" {
		return nil
	}
	allowed := map[string]map[string]string{
		"darwin.dmg": {
			"background": "string", "volume_icon": "string", "file_icon": "string", "files": "string",
			"window_width": "integer", "window_height": "integer",
		},
		"linux.appimage": {"categories": "string"},
	}[name]
	if allowed == nil {
		return fmt.Errorf("package.%s.options requires a custom template", name)
	}
	for key, value := range format.Options {
		typeName, ok := allowed[key]
		if !ok {
			return fmt.Errorf("unknown package.%s.options field %q", name, key)
		}
		valid := false
		switch typeName {
		case "string":
			_, valid = value.(string)
		case "integer":
			switch value.(type) {
			case int, int64:
				valid = true
			}
		}
		if !valid {
			return fmt.Errorf("package.%s.options.%s must be a %s", name, key, typeName)
		}
		if name == "darwin.dmg" && contains([]string{"background", "volume_icon", "file_icon"}, key) {
			path := value.(string)
			if path != "" && (manifestPathIsAbsolute(path) || pathEscapes(path)) {
				return fmt.Errorf("package.%s.options.%s must be project-relative", name, key)
			}
		}
		if name == "darwin.dmg" && key == "files" {
			for _, item := range strings.Split(value.(string), ",") {
				item = strings.TrimSpace(item)
				if item == "" {
					continue
				}
				entry, path, ok := strings.Cut(item, "=")
				path = strings.TrimSpace(path)
				if !ok || strings.TrimSpace(entry) == "" || path == "" {
					return fmt.Errorf("package.%s.options.files entry %q must be name=path", name, item)
				}
				if manifestPathIsAbsolute(path) || pathEscapes(path) {
					return fmt.Errorf("package.%s.options.files path %q must be project-relative", name, path)
				}
			}
		}
	}
	return nil
}

func hookMap(h Hooks) map[string]Hook {
	return map[string]Hook{"before_build": h.BeforeBuild, "after_build": h.AfterBuild, "before_package": h.BeforePackage, "after_package": h.AfterPackage, "before_sign": h.BeforeSign, "after_sign": h.AfterSign}
}

func signingMap(signing Signing) map[string]SigningPlatform {
	return map[string]SigningPlatform{"windows": signing.Windows, "darwin": signing.Darwin, "linux": signing.Linux, "ios": signing.IOS, "android": signing.Android}
}

func packageFormatMap(packages Packages) map[string]PackageFormat {
	return map[string]PackageFormat{
		"windows.nsis": packages.Windows.NSIS, "windows.msix": packages.Windows.MSIX,
		"darwin.app": packages.Darwin.App, "darwin.dmg": packages.Darwin.DMG,
		"linux.appimage": packages.Linux.AppImage, "linux.deb": packages.Linux.Deb,
		"linux.rpm": packages.Linux.RPM, "linux.archlinux": packages.Linux.ArchLinux,
		"ios.app": packages.IOS.App, "ios.ipa": packages.IOS.IPA,
		"android.apk": packages.Android.APK, "android.aab": packages.Android.AAB,
	}
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
