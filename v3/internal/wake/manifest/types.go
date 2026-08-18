// Package manifest owns the single user-facing build configuration for Wails.
// Callers load one resolved Config; defaulting, profile overlays, inference and
// migration provenance stay behind this package's interface.
package manifest

import "fmt"

// Filename is the single opt-in marker for Wails' native build system. Its
// presence deliberately disables legacy Taskfile routing.
const Filename = "wails.hcl"

const (
	EjectedFilename  = "wails.ejected.hcl"
	MigratedFilename = "wails.migrated.hcl"
)

type Document struct {
	Project      Project                   `toml:"project"`
	Frontend     Frontend                  `toml:"frontend,omitempty"`
	Build        Build                     `toml:"build,omitempty"`
	Dev          Dev                       `toml:"dev,omitempty"`
	Targets      Targets                   `toml:"targets,omitempty"`
	Package      Packages                  `toml:"package,omitempty"`
	Signing      Signing                   `toml:"signing,omitempty"`
	Associations []Association             `toml:"associations,omitempty"`
	Protocols    []Protocol                `toml:"protocols,omitempty"`
	Hooks        Hooks                     `toml:"hooks,omitempty"`
	Wake         Wake                      `toml:"wake,omitempty"`
	Profiles     map[string]Profile        `toml:"profiles,omitempty"`
	Extensions   map[string]map[string]any `toml:"extensions,omitempty"`
}

type Config struct {
	Root         string                    `toml:"-" json:"root"`
	Profile      string                    `toml:"-" json:"profile,omitempty"`
	Project      Project                   `toml:"project" json:"project"`
	Frontend     Frontend                  `toml:"frontend" json:"frontend"`
	Build        Build                     `toml:"build" json:"build"`
	Dev          Dev                       `toml:"dev" json:"dev"`
	Targets      Targets                   `toml:"targets" json:"targets"`
	Package      Packages                  `toml:"package" json:"package"`
	Signing      Signing                   `toml:"signing" json:"signing"`
	Associations []Association             `toml:"associations,omitempty" json:"associations,omitempty"`
	Protocols    []Protocol                `toml:"protocols,omitempty" json:"protocols,omitempty"`
	Hooks        Hooks                     `toml:"hooks" json:"hooks"`
	Wake         Wake                      `toml:"wake" json:"wake"`
	Profiles     map[string]Profile        `toml:"-" json:"profiles,omitempty"`
	Selected     Profile                   `toml:"-" json:"selected_profile,omitempty"`
	Extensions   map[string]map[string]any `toml:"extensions,omitempty" json:"extensions,omitempty"`
}

type Project struct {
	Name        string `toml:"name" json:"name"`
	ProductName string `toml:"product_name" json:"product_name"`
	Identifier  string `toml:"identifier" json:"identifier"`
	Version     string `toml:"version" json:"version"`
	BinaryName  string `toml:"binary_name,omitempty" json:"binary_name"`
	BuildNumber int    `toml:"build_number,omitempty" json:"build_number"`
	CompanyName string `toml:"company_name,omitempty" json:"company_name,omitempty"`
	Description string `toml:"description,omitempty" json:"description,omitempty"`
	Copyright   string `toml:"copyright,omitempty" json:"copyright,omitempty"`
	Comments    string `toml:"comments,omitempty" json:"comments,omitempty"`
	Icon        string `toml:"icon,omitempty" json:"icon,omitempty"`
}

type Frontend struct {
	Directory       string   `toml:"directory" json:"directory"`
	PackageManager  string   `toml:"package_manager" json:"package_manager"`
	InstallCommand  string   `toml:"install_command" json:"install_command"`
	BuildCommand    string   `toml:"build_command" json:"build_command"`
	DevCommand      string   `toml:"dev_command" json:"dev_command"`
	OutputDirectory string   `toml:"output_directory" json:"output_directory"`
	Bindings        Bindings `toml:"bindings" json:"bindings"`
	// Commands are argument vectors, never shell strings. The string fields
	// above are retained internally while the legacy Taskfile adapter exists.
	Install []string `toml:"-" json:"install,omitempty"`
	Build   []string `toml:"-" json:"build,omitempty"`
	Dev     []string `toml:"-" json:"dev,omitempty"`
}

// Profile is a complete, named production request. It intentionally has no
// generic overlay fields: a profile declares the exact targets and artifacts
// Wails should produce.
type Profile struct {
	Name    string          `toml:"-" json:"name,omitempty"`
	Targets []ProfileTarget `toml:"targets" json:"targets"`
}

type ProfileTarget struct {
	Target      string   `toml:"target" json:"target"`
	Formats     []string `toml:"formats,omitempty" json:"formats,omitempty"`
	Sign        bool     `toml:"sign,omitempty" json:"sign,omitempty"`
	Notarize    bool     `toml:"notarize,omitempty" json:"notarize,omitempty"`
	Destination string   `toml:"destination,omitempty" json:"destination,omitempty"`
}

type Bindings struct {
	TypeScript      bool   `toml:"typescript" json:"typescript"`
	Interfaces      bool   `toml:"interfaces" json:"interfaces"`
	OutputDirectory string `toml:"output_directory" json:"output_directory"`
	ModelsFilename  string `toml:"models_filename" json:"models_filename"`
	IndexFilename   string `toml:"index_filename" json:"index_filename"`
	TimeType        string `toml:"time_type" json:"time_type"`
}

type Build struct {
	OutputDirectory string  `toml:"output_directory" json:"output_directory"`
	Production      bool    `toml:"production" json:"production"`
	Obfuscation     bool    `toml:"obfuscation" json:"obfuscation"`
	TrimPath        bool    `toml:"trim_path" json:"trim_path"`
	Strip           bool    `toml:"strip" json:"strip"`
	Go              GoBuild `toml:"go" json:"go"`
}

type GoBuild struct {
	Tags          []string `toml:"tags" json:"tags"`
	LinkerFlags   []string `toml:"linker_flags" json:"linker_flags"`
	CompilerFlags []string `toml:"compiler_flags" json:"compiler_flags"`
	GarbleArgs    []string `toml:"garble_args" json:"garble_args"`
}

type Dev struct {
	Port          int      `toml:"port" json:"port"`
	DebounceMS    int      `toml:"debounce_ms" json:"debounce_ms"`
	LogLevel      string   `toml:"log_level" json:"log_level"`
	Watch         []string `toml:"watch" json:"watch"`
	Exclude       []string `toml:"exclude" json:"exclude"`
	UseGitIgnore  bool     `toml:"use_gitignore" json:"use_gitignore"`
	GracePeriodMS int      `toml:"grace_period_ms" json:"grace_period_ms"`
}

type Targets struct {
	Windows Platform `toml:"windows" json:"windows"`
	Darwin  Platform `toml:"darwin" json:"darwin"`
	Linux   Platform `toml:"linux" json:"linux"`
	IOS     Platform `toml:"ios" json:"ios"`
	Android Platform `toml:"android" json:"android"`
}

type Platform struct {
	ProductName    string   `toml:"product_name,omitempty" json:"product_name,omitempty"`
	Identifier     string   `toml:"identifier,omitempty" json:"identifier,omitempty"`
	MinimumVersion string   `toml:"minimum_version,omitempty" json:"minimum_version,omitempty"`
	BuildNumber    int      `toml:"build_number,omitempty" json:"build_number,omitempty"`
	Capabilities   []string `toml:"capabilities,omitempty" json:"capabilities,omitempty"`
	Icon           string   `toml:"icon,omitempty" json:"icon,omitempty"`
	Manifest       string   `toml:"manifest,omitempty" json:"manifest,omitempty"`
	AssetsCar      string   `toml:"assets_car,omitempty" json:"assets_car,omitempty"`
	InfoPlist      string   `toml:"info_plist,omitempty" json:"info_plist,omitempty"`
	Publisher      string   `toml:"publisher,omitempty" json:"publisher,omitempty"`
	VersionName    string   `toml:"version_name,omitempty" json:"version_name,omitempty"`
	VersionCode    int      `toml:"version_code,omitempty" json:"version_code,omitempty"`
	MinimumSDK     int      `toml:"minimum_sdk,omitempty" json:"minimum_sdk,omitempty"`
	TargetSDK      int      `toml:"target_sdk,omitempty" json:"target_sdk,omitempty"`
	AMD64          Target   `toml:"amd64,omitempty" json:"amd64,omitempty"`
	ARM64          Target   `toml:"arm64,omitempty" json:"arm64,omitempty"`
	ARM            Target   `toml:"arm,omitempty" json:"arm,omitempty"`
	X86            Target   `toml:"386,omitempty" json:"386,omitempty"`
	Universal      Target   `toml:"universal,omitempty" json:"universal,omitempty"`
}

type Target struct {
	Enabled        bool     `toml:"enabled" json:"enabled"`
	Variant        string   `toml:"variant,omitempty" json:"variant,omitempty"`
	MinimumVersion string   `toml:"minimum_version,omitempty" json:"minimum_version,omitempty"`
	BuildNumber    int      `toml:"build_number,omitempty" json:"build_number,omitempty"`
	Tags           []string `toml:"tags,omitempty" json:"tags,omitempty"`
}

type Packages struct {
	Windows PackagePlatform `toml:"windows" json:"windows"`
	Darwin  PackagePlatform `toml:"darwin" json:"darwin"`
	Linux   PackagePlatform `toml:"linux" json:"linux"`
	IOS     PackagePlatform `toml:"ios" json:"ios"`
	Android PackagePlatform `toml:"android" json:"android"`
}

type PackagePlatform struct {
	Formats   []string      `toml:"formats" json:"formats"`
	NSIS      PackageFormat `toml:"nsis,omitempty" json:"nsis,omitempty"`
	MSIX      PackageFormat `toml:"msix,omitempty" json:"msix,omitempty"`
	App       PackageFormat `toml:"app,omitempty" json:"app,omitempty"`
	DMG       PackageFormat `toml:"dmg,omitempty" json:"dmg,omitempty"`
	AppImage  PackageFormat `toml:"appimage,omitempty" json:"appimage,omitempty"`
	Deb       PackageFormat `toml:"deb,omitempty" json:"deb,omitempty"`
	RPM       PackageFormat `toml:"rpm,omitempty" json:"rpm,omitempty"`
	ArchLinux PackageFormat `toml:"archlinux,omitempty" json:"archlinux,omitempty"`
	IPA       PackageFormat `toml:"ipa,omitempty" json:"ipa,omitempty"`
	APK       PackageFormat `toml:"apk,omitempty" json:"apk,omitempty"`
	AAB       PackageFormat `toml:"aab,omitempty" json:"aab,omitempty"`
}

type PackageFormat struct {
	Template string         `toml:"template,omitempty" json:"template,omitempty"`
	Options  map[string]any `toml:"options,omitempty" json:"options,omitempty"`
}

type Signing struct {
	Windows SigningPlatform `toml:"windows" json:"windows"`
	Darwin  SigningPlatform `toml:"darwin" json:"darwin"`
	Linux   SigningPlatform `toml:"linux" json:"linux"`
	IOS     SigningPlatform `toml:"ios" json:"ios"`
	Android SigningPlatform `toml:"android" json:"android"`
}

type SigningPlatform struct {
	Enabled             bool   `toml:"enabled" json:"enabled"`
	Identity            string `toml:"identity,omitempty" json:"identity,omitempty"`
	Certificate         string `toml:"certificate,omitempty" json:"certificate,omitempty"`
	Thumbprint          string `toml:"thumbprint,omitempty" json:"thumbprint,omitempty"`
	TimestampServer     string `toml:"timestamp_server,omitempty" json:"timestamp_server,omitempty"`
	Entitlements        string `toml:"entitlements,omitempty" json:"entitlements,omitempty"`
	ProvisioningProfile string `toml:"provisioning_profile,omitempty" json:"provisioning_profile,omitempty"`
	KeyAlias            string `toml:"key_alias,omitempty" json:"key_alias,omitempty"`
	Notarize            bool   `toml:"notarize" json:"notarize"`
	Credential          string `toml:"credential,omitempty" json:"credential,omitempty"`
}

type Association struct {
	Extensions  []string `toml:"extensions" json:"extensions"`
	Name        string   `toml:"name,omitempty" json:"name,omitempty"`
	Description string   `toml:"description,omitempty" json:"description,omitempty"`
	Icon        string   `toml:"icon,omitempty" json:"icon,omitempty"`
	Role        string   `toml:"role,omitempty" json:"role,omitempty"`
	MIMEType    string   `toml:"mime_type,omitempty" json:"mime_type,omitempty"`
	Platforms   []string `toml:"platforms,omitempty" json:"platforms,omitempty"`
}

type Protocol struct {
	Scheme      string   `toml:"scheme" json:"scheme"`
	Description string   `toml:"description,omitempty" json:"description,omitempty"`
	Platforms   []string `toml:"platforms,omitempty" json:"platforms,omitempty"`
}

type Hooks struct {
	BeforeBuild   Hook `toml:"before_build,omitempty" json:"before_build,omitempty"`
	AfterBuild    Hook `toml:"after_build,omitempty" json:"after_build,omitempty"`
	BeforePackage Hook `toml:"before_package,omitempty" json:"before_package,omitempty"`
	AfterPackage  Hook `toml:"after_package,omitempty" json:"after_package,omitempty"`
	BeforeSign    Hook `toml:"before_sign,omitempty" json:"before_sign,omitempty"`
	AfterSign     Hook `toml:"after_sign,omitempty" json:"after_sign,omitempty"`
}

type Hook struct {
	Script    string   `toml:"script" json:"script"`
	Directory string   `toml:"directory,omitempty" json:"directory,omitempty"`
	Cache     bool     `toml:"cache" json:"cache"`
	Inputs    []string `toml:"inputs,omitempty" json:"inputs,omitempty"`
	Outputs   []string `toml:"outputs,omitempty" json:"outputs,omitempty"`
}

func (h *Hook) UnmarshalTOML(value any) error {
	switch value := value.(type) {
	case string:
		h.Script = value
		return nil
	case map[string]any:
		for key, raw := range value {
			switch key {
			case "script":
				text, ok := raw.(string)
				if !ok {
					return fmt.Errorf("hook field %q must be a string", key)
				}
				h.Script = text
			case "directory":
				text, ok := raw.(string)
				if !ok {
					return fmt.Errorf("hook field %q must be a string", key)
				}
				h.Directory = text
			case "cache":
				flag, ok := raw.(bool)
				if !ok {
					return fmt.Errorf("hook field %q must be a boolean", key)
				}
				h.Cache = flag
			case "inputs":
				values, ok := stringsFromAny(raw)
				if !ok {
					return fmt.Errorf("hook field %q must be an array of strings", key)
				}
				h.Inputs = values
			case "outputs":
				values, ok := stringsFromAny(raw)
				if !ok {
					return fmt.Errorf("hook field %q must be an array of strings", key)
				}
				h.Outputs = values
			default:
				return fmt.Errorf("unknown hook field %q", key)
			}
		}
		return nil
	default:
		return fmt.Errorf("hook must be a script path or table")
	}
}

func stringsFromAny(value any) ([]string, bool) {
	raw, ok := value.([]any)
	if !ok {
		return nil, false
	}
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		s, ok := item.(string)
		if !ok {
			return nil, false
		}
		result = append(result, s)
	}
	return result, true
}

type Wake struct {
	EjectedBy       string            `toml:"ejected_by,omitempty" json:"ejected_by,omitempty"`
	EjectedProfiles map[string]string `toml:"ejected_profiles,omitempty" json:"ejected_profiles,omitempty"`
}
