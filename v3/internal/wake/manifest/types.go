// Package manifest owns the single user-facing build configuration for Wails.
// Callers load one resolved Config; defaulting, profile overlays, inference and
// migration provenance stay behind this package's interface.
package manifest

// Filename is the single opt-in marker for Wails' native build system. Its
// presence deliberately disables legacy Taskfile routing.
const Filename = "wails.yaml"

const (
	EjectedFilename  = "wails.ejected.yaml"
	MigratedFilename = "wails.migrated.yaml"
)

type Run struct {
	Tags        []string          `json:"tags,omitempty"`
	Args        []string          `json:"args,omitempty"`
	ArgsSet     bool              `json:"-"`
	Environment map[string]string `json:"environment,omitempty"`
}

type Document struct {
	Binaries []Binary
	Run      Run

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
	PrimaryExecutable *Binary  `json:"-"`
	Binaries          []Binary `json:"binaries,omitempty"`
	Run               Run      `json:"run"`

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
	Revision           string   `json:"revision,omitempty"`
	Package            string   `json:"package,omitempty"`
	Kind               string   `json:"kind,omitempty"`
	VersionFile        string   `json:"version_file,omitempty"`
	VersionSymbol      string   `json:"version_symbol,omitempty"`
	CommitSymbol       string   `json:"commit_symbol,omitempty"`
	UPXLevel           int      `json:"upx_level,omitempty"`
	SupportedPlatforms []string `json:"supported_platforms,omitempty"`

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

// Binary declares an additional executable shipped with the primary application.
// Its package is local to the project; arguments are passed directly to Go.
type Binary struct {
	Name          string            `manifest:",label" schema_label:"binary" json:"name"`
	Package       string            `manifest:"package,optional" path:"true" json:"package"`
	Kind          string            `manifest:"kind,optional" json:"kind,omitempty"`
	Tags          []string          `manifest:"tags,optional" json:"tags,omitempty"`
	LinkerFlags   []string          `manifest:"ldflags,optional" json:"ldflags,omitempty"`
	CompilerFlags []string          `manifest:"compiler_flags,optional" json:"compiler_flags,omitempty"`
	Environment   map[string]string `manifest:"environment,optional" json:"environment,omitempty"`
	CommitSymbol  string            `manifest:"commit_symbol,optional" json:"commit_symbol,omitempty"`
	UPXLevel      int               `manifest:"upx_level,optional" json:"upx_level,omitempty"`
}

type Frontend struct {
	Disabled bool `json:"disabled"`

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

// DevLaunch holds development-only application arguments.
type DevLaunch struct {
	Args    []string `json:"args,omitempty"`
	ArgsSet bool     `json:"-"`
}

type Dev struct {
	Args          []string `json:"args,omitempty"`
	ArgsSet       bool     `json:"-"`
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
	MicrophoneUsageDescription string   `json:"microphone_usage_description,omitempty"`
	ProductName                string   `json:"product_name,omitempty"`
	Identifier                 string   `json:"identifier,omitempty"`
	MinimumVersion             string   `json:"minimum_version,omitempty"`
	BuildNumber                int      `json:"build_number,omitempty"`
	Capabilities               []string `json:"capabilities,omitempty"`
	Icon                       string   `json:"icon,omitempty"`
	Manifest                   string   `json:"manifest,omitempty"`
	AssetsCar                  string   `json:"assets_car,omitempty"`
	InfoPlist                  string   `json:"info_plist,omitempty"`
	Publisher                  string   `json:"publisher,omitempty"`
	DesktopEntry               string   `json:"desktop_entry,omitempty"`
	CompanyName                string   `json:"company_name,omitempty"`
	Comments                   string   `json:"comments,omitempty"`
	CFBundleIconName           string   `json:"cf_bundle_icon_name,omitempty"`
	BackgroundModes            []string `json:"background_modes,omitempty"`
	VersionName                string   `json:"version_name,omitempty"`
	VersionCode                int      `json:"version_code,omitempty"`
	MinimumSDK                 int      `json:"minimum_sdk,omitempty"`
	TargetSDK                  int      `json:"target_sdk,omitempty"`
	AMD64                      Target   `json:"amd64,omitempty"`
	ARM64                      Target   `json:"arm64,omitempty"`
	ARM                        Target   `json:"arm,omitempty"`
	X86                        Target   `json:"386,omitempty"`
	Universal                  Target   `json:"universal,omitempty"`
}

type Target struct {
	Dev DevLaunch `json:"dev"`
	Run Run       `json:"run"`

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
	Void      LinuxPackage    `json:"void,omitempty"`
	ArchLinux LinuxPackage    `json:"archlinux,omitempty"`
	IPA       RunnablePackage `json:"ipa,omitempty"`
	APK       RunnablePackage `json:"apk,omitempty"`
	AAB       RunnablePackage `json:"aab,omitempty"`
}

type RunnablePackage struct{}

type NSISPackage struct {
	Defines      map[string]string `json:"defines,omitempty"`
	Template     string            `json:"template,omitempty"`
	InstallScope string            `json:"install_scope,omitempty"`
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
	RuntimeRoot  string   `json:"runtime_root,omitempty"`
	RuntimePath  string   `json:"runtime_path,omitempty"`
	GTKVersion   int      `json:"gtk_version,omitempty"`
	UpdateInfo   string   `json:"update_info,omitempty"`
	Checksum     bool     `json:"checksum,omitempty"`
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
	RuntimeRoot  string            `manifest:"runtime_root,optional" formats:"appimage" path:"true" json:"runtime_root,omitempty"`
	RuntimePath  string            `manifest:"runtime_path,optional" formats:"appimage" path:"true" json:"runtime_path,omitempty"`
	GTKVersion   int               `manifest:"gtk_version,optional" formats:"appimage" json:"gtk_version,omitempty"`
	Defines      map[string]string `manifest:"defines,optional" formats:"nsis" json:"defines,omitempty"`
	UpdateInfo   string            `manifest:"update_info,optional" formats:"appimage" json:"update_info,omitempty"`
	Checksum     bool              `manifest:"checksum,optional" formats:"appimage" json:"checksum,omitempty"`
	Format       string            `manifest:",label" schema_label:"format" json:"format,omitempty"`
	Template     string            `manifest:"template,optional" formats:"nsis,dmg,deb,rpm,archlinux" path:"true" json:"template,omitempty"`
	InstallScope string            `manifest:"install_scope,optional" formats:"nsis" json:"install_scope,omitempty"`
	Publisher    string            `manifest:"publisher,optional" formats:"msix" json:"publisher,omitempty"`
	Manifest     string            `manifest:"manifest,optional" formats:"msix" path:"true" json:"manifest,omitempty"`
	Background   string            `manifest:"background,optional" formats:"dmg" path:"true" json:"background,omitempty"`
	VolumeIcon   string            `manifest:"volume_icon,optional" formats:"dmg" path:"true" json:"volume_icon,omitempty"`
	FileIcon     string            `manifest:"file_icon,optional" formats:"dmg" path:"true" json:"file_icon,omitempty"`
	Files        map[string]string `manifest:"files,optional" formats:"dmg" path:"true" json:"files,omitempty"`
	WindowWidth  int               `manifest:"window_width,optional" formats:"dmg" json:"window_width,omitempty"`
	WindowHeight int               `manifest:"window_height,optional" formats:"dmg" json:"window_height,omitempty"`
	Icon         string            `manifest:"icon,optional" formats:"appimage" path:"true" json:"icon,omitempty"`
	DesktopEntry string            `manifest:"desktop_entry,optional" formats:"appimage" path:"true" json:"desktop_entry,omitempty"`
	Categories   []string          `manifest:"categories,optional" formats:"appimage" json:"categories,omitempty"`
	Maintainer   string            `manifest:"maintainer,optional" formats:"deb,rpm,archlinux,void" json:"maintainer,omitempty"`
	Section      string            `manifest:"section,optional" formats:"deb,rpm,archlinux,void" json:"section,omitempty"`
	Dependencies []string          `manifest:"dependencies,optional" formats:"deb,rpm,archlinux,void" json:"dependencies"`
	PreInstall   string            `manifest:"pre_install,optional" formats:"deb,rpm,archlinux,void" path:"true" json:"pre_install,omitempty"`
	PostInstall  string            `manifest:"post_install,optional" formats:"deb,rpm,archlinux,void" path:"true" json:"post_install,omitempty"`
	PreRemove    string            `manifest:"pre_remove,optional" formats:"deb,rpm,archlinux,void" path:"true" json:"pre_remove,omitempty"`
	PostRemove   string            `manifest:"post_remove,optional" formats:"deb,rpm,archlinux,void" path:"true" json:"post_remove,omitempty"`
}

type Signing struct {
	Windows SigningPlatform `json:"windows"`
	Darwin  SigningPlatform `json:"darwin"`
	Linux   SigningPlatform `json:"linux"`
	IOS     SigningPlatform `json:"ios"`
	Android SigningPlatform `json:"android"`
}

// GlobalSigningPaths records exact paths resolved from user-level defaults.
// It is runtime provenance, never manifest input or serialized configuration.
type GlobalSigningPaths struct {
	Certificate, Entitlements, NotarizationKey string
}

type SigningPlatform struct {
	GlobalPaths            GlobalSigningPaths `json:"-" yaml:"-"`
	Enabled                bool               `json:"enabled"`
	Identity               string             `json:"identity,omitempty"`
	Certificate            string             `json:"certificate,omitempty"`
	Thumbprint             string             `json:"thumbprint,omitempty"`
	TimestampServer        string             `json:"timestamp_server,omitempty"`
	Entitlements           string             `json:"entitlements,omitempty"`
	ProvisioningProfile    string             `json:"provisioning_profile,omitempty"`
	KeyAlias               string             `json:"key_alias,omitempty"`
	Notarize               bool               `json:"notarize"`
	Credential             string             `json:"credential,omitempty"`
	NotarizationCredential string             `json:"notarization_credential,omitempty"`
	NotarizationKey        string             `json:"notarization_key,omitempty"`
	NotarizationKeyID      string             `json:"notarization_key_id,omitempty"`
	NotarizationIssuerID   string             `json:"notarization_issuer_id,omitempty"`
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
