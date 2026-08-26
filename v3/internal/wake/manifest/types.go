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
	Project      Project
	Frontend     Frontend
	Build        Build
	Dev          Dev
	Targets      Targets
	Package      Packages
	Signing      Signing
	Associations []Association
	Protocols    []Protocol
	Hooks        map[HookPhase]Hook
	Profiles     map[string]Profile
}

type Config struct {
	Root         string             `json:"root"`
	Profile      string             `json:"profile,omitempty"`
	Project      Project            `json:"project"`
	Frontend     Frontend           `json:"frontend"`
	Build        Build              `json:"build"`
	Dev          Dev                `json:"dev"`
	Targets      Targets            `json:"targets"`
	Package      Packages           `json:"package"`
	Signing      Signing            `json:"signing"`
	Associations []Association      `json:"associations,omitempty"`
	Protocols    []Protocol         `json:"protocols,omitempty"`
	Hooks        map[HookPhase]Hook `json:"hooks,omitempty"`
	Profiles     map[string]Profile `json:"profiles,omitempty"`
	Selected     Profile            `json:"selected_profile,omitempty"`
	Origins      map[string]Origin  `json:"origins,omitempty"`
}

type HookPhase string

const (
	BeforeBuild   HookPhase = "before_build"
	AfterBuild    HookPhase = "after_build"
	BeforePackage HookPhase = "before_package"
	AfterPackage  HookPhase = "after_package"
	BeforeSign    HookPhase = "before_sign"
	AfterSign     HookPhase = "after_sign"
)

var HookPhases = []HookPhase{BeforeBuild, AfterBuild, BeforePackage, AfterPackage, BeforeSign, AfterSign}

// Hook is a bounded invocation of one project-owned executable script. It is
// intentionally not a shell command or a programmable pipeline definition.
type Hook struct {
	Script    string   `json:"script"`
	Directory string   `json:"directory,omitempty"`
	Cache     bool     `json:"cache,omitempty"`
	Inputs    []string `json:"inputs,omitempty"`
	Outputs   []string `json:"outputs,omitempty"`
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
	Name        string `json:"name"`
	ProductName string `json:"product_name"`
	Identifier  string `json:"identifier"`
	Version     string `json:"version"`
	BinaryName  string `json:"binary_name"`
	BuildNumber int    `json:"build_number"`
	CompanyName string `json:"company_name,omitempty"`
	Description string `json:"description,omitempty"`
	Copyright   string `json:"copyright,omitempty"`
	Comments    string `json:"comments,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

type Frontend struct {
	Directory       string            `json:"directory"`
	PackageManager  string            `json:"package_manager"`
	InstallCommand  string            `json:"install_command"`
	BuildCommand    string            `json:"build_command"`
	DevCommand      string            `json:"dev_command"`
	OutputDirectory string            `json:"output_directory"`
	Bindings        Bindings          `json:"bindings"`
	Environment     map[string]string `json:"environment,omitempty"`
	// Commands are argument vectors, never shell strings. The string fields
	// above are retained internally while the legacy Taskfile adapter exists.
	Install []string `json:"install,omitempty"`
	Build   []string `json:"build,omitempty"`
	Dev     []string `json:"dev,omitempty"`
}

// Profile is a complete, named production request. It intentionally has no
// generic overlay fields: a profile declares the exact targets and artifacts
// Wails should produce.
type Profile struct {
	Name    string          `json:"name,omitempty"`
	Targets []ProfileTarget `json:"targets"`
}

type ProfileTarget struct {
	Target      string   `json:"target"`
	Formats     []string `json:"formats,omitempty"`
	Sign        bool     `json:"sign,omitempty"`
	Notarize    bool     `json:"notarize,omitempty"`
	Destination string   `json:"destination,omitempty"`
}

type Bindings struct {
	TypeScript      bool   `json:"typescript"`
	Interfaces      bool   `json:"interfaces"`
	OutputDirectory string `json:"output_directory"`
	ModelsFilename  string `json:"models_filename"`
	IndexFilename   string `json:"index_filename"`
	TimeType        string `json:"time_type"`
}

type Build struct {
	OutputDirectory string            `json:"output_directory"`
	Production      bool              `json:"production"`
	Obfuscation     bool              `json:"obfuscation"`
	TrimPath        bool              `json:"trim_path"`
	Strip           bool              `json:"strip"`
	VCSInfo         bool              `json:"vcs_info"`
	Environment     map[string]string `json:"environment,omitempty"`
	Go              GoBuild           `json:"go"`
}

type GoBuild struct {
	Tags          []string `json:"tags"`
	LinkerFlags   []string `json:"linker_flags"`
	CompilerFlags []string `json:"compiler_flags"`
	GarbleArgs    []string `json:"garble_args"`
}

type Dev struct {
	Port          int      `json:"port"`
	Tags          []string `json:"tags,omitempty"`
	DebounceMS    int      `json:"debounce_ms"`
	LogLevel      string   `json:"log_level"`
	Watch         []string `json:"watch"`
	Exclude       []string `json:"exclude"`
	UseGitIgnore  bool     `json:"use_gitignore"`
	GracePeriodMS int      `json:"grace_period_ms"`
}

type Targets struct {
	Windows Platform `json:"windows"`
	Darwin  Platform `json:"darwin"`
	Linux   Platform `json:"linux"`
	IOS     Platform `json:"ios"`
	Android Platform `json:"android"`
}

type Platform struct {
	ProductName      string   `json:"product_name,omitempty"`
	Identifier       string   `json:"identifier,omitempty"`
	MinimumVersion   string   `json:"minimum_version,omitempty"`
	BuildNumber      int      `json:"build_number,omitempty"`
	Capabilities     []string `json:"capabilities,omitempty"`
	Icon             string   `json:"icon,omitempty"`
	Manifest         string   `json:"manifest,omitempty"`
	AssetsCar        string   `json:"assets_car,omitempty"`
	InfoPlist        string   `json:"info_plist,omitempty"`
	Publisher        string   `json:"publisher,omitempty"`
	DesktopEntry     string   `json:"desktop_entry,omitempty"`
	CompanyName      string   `json:"company_name,omitempty"`
	Comments         string   `json:"comments,omitempty"`
	CFBundleIconName string   `json:"cf_bundle_icon_name,omitempty"`
	BackgroundModes  []string `json:"background_modes,omitempty"`
	VersionName      string   `json:"version_name,omitempty"`
	VersionCode      int      `json:"version_code,omitempty"`
	MinimumSDK       int      `json:"minimum_sdk,omitempty"`
	TargetSDK        int      `json:"target_sdk,omitempty"`
	AMD64            Target   `json:"amd64,omitempty"`
	ARM64            Target   `json:"arm64,omitempty"`
	ARM              Target   `json:"arm,omitempty"`
	X86              Target   `json:"386,omitempty"`
	Universal        Target   `json:"universal,omitempty"`
}

type Target struct {
	Enabled        bool              `json:"enabled"`
	MinimumVersion string            `json:"minimum_version,omitempty"`
	BuildNumber    int               `json:"build_number,omitempty"`
	Tags           []string          `json:"tags,omitempty"`
	Toolchain      string            `json:"toolchain,omitempty"`
	Environment    map[string]string `json:"environment,omitempty"`
	LinkerFlags    []string          `json:"linker_flags,omitempty"`
	CompilerFlags  []string          `json:"compiler_flags,omitempty"`
	GarbleArgs     []string          `json:"garble_args,omitempty"`
	Obfuscated     bool              `json:"obfuscated,omitempty"`
	ObfuscatedSet  bool              `json:"-"`
}

type Packages struct {
	Windows PackagePlatform `json:"windows"`
	Darwin  PackagePlatform `json:"darwin"`
	Linux   PackagePlatform `json:"linux"`
	IOS     PackagePlatform `json:"ios"`
	Android PackagePlatform `json:"android"`
}

type PackagePlatform struct {
	Formats   []string        `json:"formats"`
	NSIS      NSISPackage     `json:"nsis,omitempty"`
	MSIX      MSIXPackage     `json:"msix,omitempty"`
	App       RunnablePackage `json:"app,omitempty"`
	DMG       DMGPackage      `json:"dmg,omitempty"`
	AppImage  AppImagePackage `json:"appimage,omitempty"`
	Deb       LinuxPackage    `json:"deb,omitempty"`
	RPM       LinuxPackage    `json:"rpm,omitempty"`
	ArchLinux LinuxPackage    `json:"archlinux,omitempty"`
	IPA       RunnablePackage `json:"ipa,omitempty"`
	APK       RunnablePackage `json:"apk,omitempty"`
	AAB       RunnablePackage `json:"aab,omitempty"`
}

type RunnablePackage struct{}

type NSISPackage struct {
	Template     string `json:"template,omitempty"`
	InstallScope string `json:"install_scope,omitempty"`
}

type MSIXPackage struct {
	Publisher string `json:"publisher,omitempty"`
	Manifest  string `json:"manifest,omitempty"`
}

type DMGPackage struct {
	Template     string            `json:"template,omitempty"`
	Background   string            `json:"background,omitempty"`
	VolumeIcon   string            `json:"volume_icon,omitempty"`
	FileIcon     string            `json:"file_icon,omitempty"`
	Files        map[string]string `json:"files,omitempty"`
	WindowWidth  int               `json:"window_width,omitempty"`
	WindowHeight int               `json:"window_height,omitempty"`
}

type AppImagePackage struct {
	Icon         string   `json:"icon,omitempty"`
	DesktopEntry string   `json:"desktop_entry,omitempty"`
	Categories   []string `json:"categories,omitempty"`
}

type LinuxPackage struct {
	Template     string   `json:"template,omitempty"`
	Maintainer   string   `json:"maintainer,omitempty"`
	Section      string   `json:"section,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
	PreInstall   string   `json:"pre_install,omitempty"`
	PostInstall  string   `json:"post_install,omitempty"`
	PreRemove    string   `json:"pre_remove,omitempty"`
	PostRemove   string   `json:"post_remove,omitempty"`
}

type PackageFormat struct {
	Format       string            `hcl:",label" schema_label:"format" json:"format,omitempty"`
	Template     string            `hcl:"template,optional" formats:"nsis,dmg,deb,rpm,archlinux" path:"true" json:"template,omitempty"`
	InstallScope string            `hcl:"install_scope,optional" formats:"nsis" json:"install_scope,omitempty"`
	Publisher    string            `hcl:"publisher,optional" formats:"msix" json:"publisher,omitempty"`
	Manifest     string            `hcl:"manifest,optional" formats:"msix" path:"true" json:"manifest,omitempty"`
	Background   string            `hcl:"background,optional" formats:"dmg" path:"true" json:"background,omitempty"`
	VolumeIcon   string            `hcl:"volume_icon,optional" formats:"dmg" path:"true" json:"volume_icon,omitempty"`
	FileIcon     string            `hcl:"file_icon,optional" formats:"dmg" path:"true" json:"file_icon,omitempty"`
	Files        map[string]string `hcl:"files,optional" formats:"dmg" path:"true" json:"files,omitempty"`
	WindowWidth  int               `hcl:"window_width,optional" formats:"dmg" json:"window_width,omitempty"`
	WindowHeight int               `hcl:"window_height,optional" formats:"dmg" json:"window_height,omitempty"`
	Icon         string            `hcl:"icon,optional" formats:"appimage" path:"true" json:"icon,omitempty"`
	DesktopEntry string            `hcl:"desktop_entry,optional" formats:"appimage" path:"true" json:"desktop_entry,omitempty"`
	Categories   []string          `hcl:"categories,optional" formats:"appimage" json:"categories,omitempty"`
	Maintainer   string            `hcl:"maintainer,optional" formats:"deb,rpm,archlinux" json:"maintainer,omitempty"`
	Section      string            `hcl:"section,optional" formats:"deb,rpm,archlinux" json:"section,omitempty"`
	Dependencies []string          `hcl:"dependencies,optional" formats:"deb,rpm,archlinux" json:"dependencies,omitempty"`
	PreInstall   string            `hcl:"pre_install,optional" formats:"deb,rpm,archlinux" path:"true" json:"pre_install,omitempty"`
	PostInstall  string            `hcl:"post_install,optional" formats:"deb,rpm,archlinux" path:"true" json:"post_install,omitempty"`
	PreRemove    string            `hcl:"pre_remove,optional" formats:"deb,rpm,archlinux" path:"true" json:"pre_remove,omitempty"`
	PostRemove   string            `hcl:"post_remove,optional" formats:"deb,rpm,archlinux" path:"true" json:"post_remove,omitempty"`
}

type Signing struct {
	Windows SigningPlatform `json:"windows"`
	Darwin  SigningPlatform `json:"darwin"`
	Linux   SigningPlatform `json:"linux"`
	IOS     SigningPlatform `json:"ios"`
	Android SigningPlatform `json:"android"`
}

type SigningPlatform struct {
	Enabled                bool   `json:"enabled"`
	Identity               string `json:"identity,omitempty"`
	Certificate            string `json:"certificate,omitempty"`
	Thumbprint             string `json:"thumbprint,omitempty"`
	TimestampServer        string `json:"timestamp_server,omitempty"`
	Entitlements           string `json:"entitlements,omitempty"`
	ProvisioningProfile    string `json:"provisioning_profile,omitempty"`
	KeyAlias               string `json:"key_alias,omitempty"`
	Notarize               bool   `json:"notarize"`
	Credential             string `json:"credential,omitempty"`
	NotarizationCredential string `json:"notarization_credential,omitempty"`
}

type Association struct {
	Extensions  []string `json:"extensions"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Icon        string   `json:"icon,omitempty"`
	Role        string   `json:"role,omitempty"`
	MIMEType    string   `json:"mime_type,omitempty"`
	Platforms   []string `json:"platforms,omitempty"`
}

type Protocol struct {
	Scheme      string   `json:"scheme"`
	Description string   `json:"description,omitempty"`
	Platforms   []string `json:"platforms,omitempty"`
}
