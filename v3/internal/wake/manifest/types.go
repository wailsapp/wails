// Package manifest owns the single user-facing build configuration for Wails.
// Callers load one resolved Config; defaulting, profile overlays, inference and
// migration provenance stay behind this package's interface.
package manifest

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
	Wake         Wake                      `toml:"wake" json:"wake"`
	Profiles     map[string]Profile        `toml:"-" json:"profiles,omitempty"`
	Selected     Profile                   `toml:"-" json:"selected_profile,omitempty"`
	Extensions   map[string]map[string]any `toml:"extensions,omitempty" json:"extensions,omitempty"`
	Origins      map[string]Origin         `toml:"-" json:"origins,omitempty"`
}

type OriginKind string

const (
	OriginDefault  OriginKind = "default"
	OriginManifest OriginKind = "manifest"
)

type SourceRange struct {
	Filename    string `json:"filename,omitempty"`
	StartLine   int    `json:"start_line,omitempty"`
	StartColumn int    `json:"start_column,omitempty"`
	EndLine     int    `json:"end_line,omitempty"`
	EndColumn   int    `json:"end_column,omitempty"`
}

type Origin struct {
	Kind  OriginKind  `json:"kind"`
	Range SourceRange `json:"range,omitempty"`
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
	Directory       string            `toml:"directory" json:"directory"`
	PackageManager  string            `toml:"package_manager" json:"package_manager"`
	InstallCommand  string            `toml:"install_command" json:"install_command"`
	BuildCommand    string            `toml:"build_command" json:"build_command"`
	DevCommand      string            `toml:"dev_command" json:"dev_command"`
	OutputDirectory string            `toml:"output_directory" json:"output_directory"`
	Bindings        Bindings          `toml:"bindings" json:"bindings"`
	Environment     map[string]string `toml:"environment,omitempty" json:"environment,omitempty"`
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
	OutputDirectory string            `toml:"output_directory" json:"output_directory"`
	Production      bool              `toml:"production" json:"production"`
	Obfuscation     bool              `toml:"obfuscation" json:"obfuscation"`
	TrimPath        bool              `toml:"trim_path" json:"trim_path"`
	Strip           bool              `toml:"strip" json:"strip"`
	VCSInfo         bool              `toml:"vcs_info" json:"vcs_info"`
	Environment     map[string]string `toml:"environment,omitempty" json:"environment,omitempty"`
	Go              GoBuild           `toml:"go" json:"go"`
}

type GoBuild struct {
	Tags          []string `toml:"tags" json:"tags"`
	LinkerFlags   []string `toml:"linker_flags" json:"linker_flags"`
	CompilerFlags []string `toml:"compiler_flags" json:"compiler_flags"`
	GarbleArgs    []string `toml:"garble_args" json:"garble_args"`
}

type Dev struct {
	Port          int      `toml:"port" json:"port"`
	Tags          []string `toml:"tags,omitempty" json:"tags,omitempty"`
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
	ProductName      string   `toml:"product_name,omitempty" json:"product_name,omitempty"`
	Identifier       string   `toml:"identifier,omitempty" json:"identifier,omitempty"`
	MinimumVersion   string   `toml:"minimum_version,omitempty" json:"minimum_version,omitempty"`
	BuildNumber      int      `toml:"build_number,omitempty" json:"build_number,omitempty"`
	Capabilities     []string `toml:"capabilities,omitempty" json:"capabilities,omitempty"`
	Icon             string   `toml:"icon,omitempty" json:"icon,omitempty"`
	Manifest         string   `toml:"manifest,omitempty" json:"manifest,omitempty"`
	AssetsCar        string   `toml:"assets_car,omitempty" json:"assets_car,omitempty"`
	InfoPlist        string   `toml:"info_plist,omitempty" json:"info_plist,omitempty"`
	Publisher        string   `toml:"publisher,omitempty" json:"publisher,omitempty"`
	DesktopEntry     string   `toml:"desktop_entry,omitempty" json:"desktop_entry,omitempty"`
	CompanyName      string   `toml:"company_name,omitempty" json:"company_name,omitempty"`
	Comments         string   `toml:"comments,omitempty" json:"comments,omitempty"`
	CFBundleIconName string   `toml:"cf_bundle_icon_name,omitempty" json:"cf_bundle_icon_name,omitempty"`
	BackgroundModes  []string `toml:"background_modes,omitempty" json:"background_modes,omitempty"`
	VersionName      string   `toml:"version_name,omitempty" json:"version_name,omitempty"`
	VersionCode      int      `toml:"version_code,omitempty" json:"version_code,omitempty"`
	MinimumSDK       int      `toml:"minimum_sdk,omitempty" json:"minimum_sdk,omitempty"`
	TargetSDK        int      `toml:"target_sdk,omitempty" json:"target_sdk,omitempty"`
	AMD64            Target   `toml:"amd64,omitempty" json:"amd64,omitempty"`
	ARM64            Target   `toml:"arm64,omitempty" json:"arm64,omitempty"`
	ARM              Target   `toml:"arm,omitempty" json:"arm,omitempty"`
	X86              Target   `toml:"386,omitempty" json:"386,omitempty"`
	Universal        Target   `toml:"universal,omitempty" json:"universal,omitempty"`
}

type Target struct {
	Enabled        bool              `toml:"enabled" json:"enabled"`
	MinimumVersion string            `toml:"minimum_version,omitempty" json:"minimum_version,omitempty"`
	BuildNumber    int               `toml:"build_number,omitempty" json:"build_number,omitempty"`
	Tags           []string          `toml:"tags,omitempty" json:"tags,omitempty"`
	Toolchain      string            `toml:"toolchain,omitempty" json:"toolchain,omitempty"`
	Environment    map[string]string `toml:"environment,omitempty" json:"environment,omitempty"`
	LinkerFlags    []string          `toml:"linker_flags,omitempty" json:"linker_flags,omitempty"`
	CompilerFlags  []string          `toml:"compiler_flags,omitempty" json:"compiler_flags,omitempty"`
	GarbleArgs     []string          `toml:"garble_args,omitempty" json:"garble_args,omitempty"`
	Obfuscated     bool              `toml:"obfuscated,omitempty" json:"obfuscated,omitempty"`
	ObfuscatedSet  bool              `toml:"-" json:"-"`
}

type Packages struct {
	Windows PackagePlatform `toml:"windows" json:"windows"`
	Darwin  PackagePlatform `toml:"darwin" json:"darwin"`
	Linux   PackagePlatform `toml:"linux" json:"linux"`
	IOS     PackagePlatform `toml:"ios" json:"ios"`
	Android PackagePlatform `toml:"android" json:"android"`
}

type PackagePlatform struct {
	Formats   []string        `toml:"formats" json:"formats"`
	NSIS      NSISPackage     `toml:"nsis,omitempty" json:"nsis,omitempty"`
	MSIX      MSIXPackage     `toml:"msix,omitempty" json:"msix,omitempty"`
	App       RunnablePackage `toml:"app,omitempty" json:"app,omitempty"`
	DMG       DMGPackage      `toml:"dmg,omitempty" json:"dmg,omitempty"`
	AppImage  AppImagePackage `toml:"appimage,omitempty" json:"appimage,omitempty"`
	Deb       LinuxPackage    `toml:"deb,omitempty" json:"deb,omitempty"`
	RPM       LinuxPackage    `toml:"rpm,omitempty" json:"rpm,omitempty"`
	ArchLinux LinuxPackage    `toml:"archlinux,omitempty" json:"archlinux,omitempty"`
	IPA       RunnablePackage `toml:"ipa,omitempty" json:"ipa,omitempty"`
	APK       RunnablePackage `toml:"apk,omitempty" json:"apk,omitempty"`
	AAB       RunnablePackage `toml:"aab,omitempty" json:"aab,omitempty"`
}

type RunnablePackage struct{}

type NSISPackage struct {
	Template     string `toml:"template,omitempty" json:"template,omitempty"`
	InstallScope string `toml:"install_scope,omitempty" json:"install_scope,omitempty"`
}

type MSIXPackage struct {
	Publisher string `toml:"publisher,omitempty" json:"publisher,omitempty"`
	Manifest  string `toml:"manifest,omitempty" json:"manifest,omitempty"`
}

type DMGPackage struct {
	Template     string            `toml:"template,omitempty" json:"template,omitempty"`
	Background   string            `toml:"background,omitempty" json:"background,omitempty"`
	VolumeIcon   string            `toml:"volume_icon,omitempty" json:"volume_icon,omitempty"`
	FileIcon     string            `toml:"file_icon,omitempty" json:"file_icon,omitempty"`
	Files        map[string]string `toml:"files,omitempty" json:"files,omitempty"`
	WindowWidth  int               `toml:"window_width,omitempty" json:"window_width,omitempty"`
	WindowHeight int               `toml:"window_height,omitempty" json:"window_height,omitempty"`
}

type AppImagePackage struct {
	Icon         string   `toml:"icon,omitempty" json:"icon,omitempty"`
	DesktopEntry string   `toml:"desktop_entry,omitempty" json:"desktop_entry,omitempty"`
	Categories   []string `toml:"categories,omitempty" json:"categories,omitempty"`
}

type LinuxPackage struct {
	Template     string   `toml:"template,omitempty" json:"template,omitempty"`
	Maintainer   string   `toml:"maintainer,omitempty" json:"maintainer,omitempty"`
	Section      string   `toml:"section,omitempty" json:"section,omitempty"`
	Dependencies []string `toml:"dependencies,omitempty" json:"dependencies,omitempty"`
	PreInstall   string   `toml:"pre_install,omitempty" json:"pre_install,omitempty"`
	PostInstall  string   `toml:"post_install,omitempty" json:"post_install,omitempty"`
	PreRemove    string   `toml:"pre_remove,omitempty" json:"pre_remove,omitempty"`
	PostRemove   string   `toml:"post_remove,omitempty" json:"post_remove,omitempty"`
}

type PackageFormat struct {
	Format       string            `hcl:",label" schema_label:"format" toml:"format,omitempty" json:"format,omitempty"`
	Template     string            `hcl:"template,optional" formats:"nsis,dmg,deb,rpm,archlinux" path:"true" toml:"template,omitempty" json:"template,omitempty"`
	InstallScope string            `hcl:"install_scope,optional" formats:"nsis" toml:"install_scope,omitempty" json:"install_scope,omitempty"`
	Publisher    string            `hcl:"publisher,optional" formats:"msix" toml:"publisher,omitempty" json:"publisher,omitempty"`
	Manifest     string            `hcl:"manifest,optional" formats:"msix" path:"true" toml:"manifest,omitempty" json:"manifest,omitempty"`
	Background   string            `hcl:"background,optional" formats:"dmg" path:"true" toml:"background,omitempty" json:"background,omitempty"`
	VolumeIcon   string            `hcl:"volume_icon,optional" formats:"dmg" path:"true" toml:"volume_icon,omitempty" json:"volume_icon,omitempty"`
	FileIcon     string            `hcl:"file_icon,optional" formats:"dmg" path:"true" toml:"file_icon,omitempty" json:"file_icon,omitempty"`
	Files        map[string]string `hcl:"files,optional" formats:"dmg" path:"true" toml:"files,omitempty" json:"files,omitempty"`
	WindowWidth  int               `hcl:"window_width,optional" formats:"dmg" toml:"window_width,omitempty" json:"window_width,omitempty"`
	WindowHeight int               `hcl:"window_height,optional" formats:"dmg" toml:"window_height,omitempty" json:"window_height,omitempty"`
	Icon         string            `hcl:"icon,optional" formats:"appimage" path:"true" toml:"icon,omitempty" json:"icon,omitempty"`
	DesktopEntry string            `hcl:"desktop_entry,optional" formats:"appimage" path:"true" toml:"desktop_entry,omitempty" json:"desktop_entry,omitempty"`
	Categories   []string          `hcl:"categories,optional" formats:"appimage" toml:"categories,omitempty" json:"categories,omitempty"`
	Maintainer   string            `hcl:"maintainer,optional" formats:"deb,rpm,archlinux" toml:"maintainer,omitempty" json:"maintainer,omitempty"`
	Section      string            `hcl:"section,optional" formats:"deb,rpm,archlinux" toml:"section,omitempty" json:"section,omitempty"`
	Dependencies []string          `hcl:"dependencies,optional" formats:"deb,rpm,archlinux" toml:"dependencies,omitempty" json:"dependencies,omitempty"`
	PreInstall   string            `hcl:"pre_install,optional" formats:"deb,rpm,archlinux" path:"true" toml:"pre_install,omitempty" json:"pre_install,omitempty"`
	PostInstall  string            `hcl:"post_install,optional" formats:"deb,rpm,archlinux" path:"true" toml:"post_install,omitempty" json:"post_install,omitempty"`
	PreRemove    string            `hcl:"pre_remove,optional" formats:"deb,rpm,archlinux" path:"true" toml:"pre_remove,omitempty" json:"pre_remove,omitempty"`
	PostRemove   string            `hcl:"post_remove,optional" formats:"deb,rpm,archlinux" path:"true" toml:"post_remove,omitempty" json:"post_remove,omitempty"`
}

type Signing struct {
	Windows SigningPlatform `toml:"windows" json:"windows"`
	Darwin  SigningPlatform `toml:"darwin" json:"darwin"`
	Linux   SigningPlatform `toml:"linux" json:"linux"`
	IOS     SigningPlatform `toml:"ios" json:"ios"`
	Android SigningPlatform `toml:"android" json:"android"`
}

type SigningPlatform struct {
	Enabled                bool   `toml:"enabled" json:"enabled"`
	Identity               string `toml:"identity,omitempty" json:"identity,omitempty"`
	Certificate            string `toml:"certificate,omitempty" json:"certificate,omitempty"`
	Thumbprint             string `toml:"thumbprint,omitempty" json:"thumbprint,omitempty"`
	TimestampServer        string `toml:"timestamp_server,omitempty" json:"timestamp_server,omitempty"`
	Entitlements           string `toml:"entitlements,omitempty" json:"entitlements,omitempty"`
	ProvisioningProfile    string `toml:"provisioning_profile,omitempty" json:"provisioning_profile,omitempty"`
	KeyAlias               string `toml:"key_alias,omitempty" json:"key_alias,omitempty"`
	Notarize               bool   `toml:"notarize" json:"notarize"`
	Credential             string `toml:"credential,omitempty" json:"credential,omitempty"`
	NotarizationCredential string `toml:"notarization_credential,omitempty" json:"notarization_credential,omitempty"`
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

type Wake struct {
	EjectedBy       string            `toml:"ejected_by,omitempty" json:"ejected_by,omitempty"`
	EjectedProfiles map[string]string `toml:"ejected_profiles,omitempty" json:"ejected_profiles,omitempty"`
}
