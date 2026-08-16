package manifest

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/BurntSushi/toml"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Loaded struct {
	Path     string
	Raw      []byte
	Document Document
	Config   Config
}

func Exists(root string) bool {
	_, err := os.Stat(filepath.Join(root, Filename))
	return err == nil
}

func Load(root, profile string) (*Loaded, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(absRoot, Filename)
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", Filename, err)
	}

	var identity Document
	if _, err := toml.Decode(string(raw), &identity); err != nil {
		return nil, fmt.Errorf("parse %s: %w", Filename, err)
	}
	if err := validateProject(identity.Project); err != nil {
		return nil, err
	}

	doc := defaults(identity.Project)
	metadata, err := toml.Decode(string(raw), &doc)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", Filename, err)
	}
	if undecoded := metadata.Undecoded(); len(undecoded) > 0 {
		items := make([]string, 0, len(undecoded))
		for _, key := range undecoded {
			name := key.String()
			// Profiles are validated when selected and extension payloads are
			// deliberately opaque to core Wails.
			if strings.HasPrefix(name, "profiles.") || strings.HasPrefix(name, "extensions.") || name == "wake.migration" || strings.HasPrefix(name, "wake.migration.") {
				continue
			}
			items = append(items, name)
		}
		if len(items) > 0 {
			sort.Strings(items)
			return nil, fmt.Errorf("unknown %s field(s): %s", Filename, strings.Join(items, ", "))
		}
	}
	resolveInferred(absRoot, &doc, metadata)
	config := configFromDocument(absRoot, profile, doc)
	if profile != "" {
		if profile == "default" || !slugPattern.MatchString(profile) {
			return nil, fmt.Errorf("profile name must be a lowercase slug and cannot be default")
		}
		rawProfile, ok := doc.Profiles[profile]
		if !ok {
			return nil, fmt.Errorf("profile %q is not defined", profile)
		}
		if err := validateProfileIdentity(profile, rawProfile); err != nil {
			return nil, err
		}
		layer := profileLayerFromConfig(config)
		var encoded bytes.Buffer
		if err := toml.NewEncoder(&encoded).Encode(rawProfile); err != nil {
			return nil, err
		}
		profileMetadata, err := toml.Decode(encoded.String(), &layer)
		if err != nil {
			return nil, fmt.Errorf("resolve profile %q: %w", profile, err)
		}
		if undecoded := profileMetadata.Undecoded(); len(undecoded) > 0 {
			items := make([]string, 0, len(undecoded))
			for _, key := range undecoded {
				items = append(items, "profiles."+profile+"."+key.String())
			}
			sort.Strings(items)
			return nil, fmt.Errorf("unknown %s field(s): %s", Filename, strings.Join(items, ", "))
		}
		applyProfileLayer(&config, layer)
		config.Profile = profile
	}
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	return &Loaded{Path: path, Raw: raw, Document: doc, Config: config}, nil
}

func validateProfileIdentity(profile string, raw map[string]any) error {
	if _, exists := raw["project"]; exists {
		return fmt.Errorf("profile %q cannot override project identity", profile)
	}
	targets, ok := raw["targets"].(map[string]any)
	if !ok {
		return nil
	}
	forbidden := map[string]bool{"identifier": true, "product_name": true, "version": true, "build_number": true}
	var visit func(string, map[string]any) error
	visit = func(path string, table map[string]any) error {
		for key, value := range table {
			if forbidden[key] {
				return fmt.Errorf("profile %q cannot override target identity field %s.%s", profile, path, key)
			}
			if child, ok := value.(map[string]any); ok {
				if err := visit(path+"."+key, child); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return visit("targets", targets)
}

func defaults(project Project) Document {
	if project.BinaryName == "" {
		project.BinaryName = deriveBinaryName(project.Name)
	}
	if project.BuildNumber == 0 {
		project.BuildNumber = 1
	}
	frontend := Frontend{
		Directory: "frontend", PackageManager: "auto", InstallCommand: "install",
		BuildCommand: "build", DevCommand: "dev", OutputDirectory: "dist",
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
	return Document{Project: project, Frontend: frontend, Build: build, Dev: dev, Targets: targets, Package: packages, Profiles: map[string]map[string]any{}, Extensions: map[string]map[string]any{}}
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

func resolveInferred(root string, doc *Document, metadata toml.MetaData) {
	frontendRoot := filepath.Join(root, doc.Frontend.Directory)
	if doc.Frontend.PackageManager == "auto" {
		for _, candidate := range []struct{ file, manager string }{{"bun.lock", "bun"}, {"bun.lockb", "bun"}, {"pnpm-lock.yaml", "pnpm"}, {"yarn.lock", "yarn"}, {"package-lock.json", "npm"}, {"npm-shrinkwrap.json", "npm"}} {
			if _, err := os.Stat(filepath.Join(frontendRoot, candidate.file)); err == nil {
				doc.Frontend.PackageManager = candidate.manager
				break
			}
		}
		if doc.Frontend.PackageManager == "auto" {
			doc.Frontend.PackageManager = "npm"
		}
	}
	if _, err := os.Stat(filepath.Join(frontendRoot, "tsconfig.json")); errors.Is(err, fs.ErrNotExist) {
		if !metadata.IsDefined("frontend", "bindings", "typescript") {
			doc.Frontend.Bindings.TypeScript = false
		}
		if !metadata.IsDefined("frontend", "bindings", "interfaces") {
			doc.Frontend.Bindings.Interfaces = false
		}
		if !metadata.IsDefined("frontend", "bindings", "time_type") {
			doc.Frontend.Bindings.TimeType = "Date"
		}
	}
	if doc.Project.Icon == "" {
		for _, path := range []string{"assets/appicon.png", "build/appicon.png"} {
			if _, err := os.Stat(filepath.Join(root, path)); err == nil {
				doc.Project.Icon = path
				break
			}
		}
	}
}

func configFromDocument(root, profile string, doc Document) Config {
	return Config{Root: root, Profile: profile, Project: doc.Project, Frontend: doc.Frontend, Build: doc.Build, Dev: doc.Dev, Targets: doc.Targets, Package: doc.Package, Signing: doc.Signing, Associations: doc.Associations, Protocols: doc.Protocols, Hooks: doc.Hooks, Wake: doc.Wake, Extensions: doc.Extensions}
}

type profileLayer struct {
	Frontend     Frontend                  `toml:"frontend"`
	Build        Build                     `toml:"build"`
	Dev          Dev                       `toml:"dev"`
	Targets      Targets                   `toml:"targets"`
	Package      Packages                  `toml:"package"`
	Signing      Signing                   `toml:"signing"`
	Associations []Association             `toml:"associations,omitempty"`
	Protocols    []Protocol                `toml:"protocols,omitempty"`
	Hooks        Hooks                     `toml:"hooks"`
	Extensions   map[string]map[string]any `toml:"extensions,omitempty"`
}

func profileLayerFromConfig(config Config) profileLayer {
	return profileLayer{Frontend: config.Frontend, Build: config.Build, Dev: config.Dev, Targets: config.Targets, Package: config.Package, Signing: config.Signing, Associations: config.Associations, Protocols: config.Protocols, Hooks: config.Hooks, Extensions: config.Extensions}
}

func applyProfileLayer(config *Config, layer profileLayer) {
	config.Frontend, config.Build, config.Dev, config.Targets = layer.Frontend, layer.Build, layer.Dev, layer.Targets
	config.Package, config.Signing, config.Associations, config.Protocols, config.Hooks, config.Extensions = layer.Package, layer.Signing, layer.Associations, layer.Protocols, layer.Hooks, layer.Extensions
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
	}
	for name, format := range packageFormatMap(config.Package) {
		if format.Template != "" && (manifestPathIsAbsolute(format.Template) || pathEscapes(format.Template)) {
			return fmt.Errorf("package.%s.template must be project-relative", name)
		}
	}
	return nil
}

func hookMap(h Hooks) map[string]Hook {
	return map[string]Hook{"before_build": h.BeforeBuild, "after_build": h.AfterBuild, "before_package": h.BeforePackage, "after_package": h.AfterPackage, "before_sign": h.BeforeSign, "after_sign": h.AfterSign}
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
