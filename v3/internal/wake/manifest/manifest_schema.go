package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v3/internal/wake/buildinfo"
)

// manifestDocument is intentionally a closed schema. Keep the YAML surface small:
// this is declarative build intent, never a programmable pipeline.
type manifestDocument struct {
	Run *manifestRun `manifest:"run,block"`

	Version      int                   `manifest:"version"`
	Project      *manifestProject      `manifest:"project,block" required:"true"`
	Frontend     *manifestFrontend     `manifest:"frontend,block"`
	Build        *manifestBuild        `manifest:"build,block"`
	Dev          *manifestDev          `manifest:"dev,block"`
	Windows      *manifestPlatform     `manifest:"windows,block"`
	Darwin       *manifestPlatform     `manifest:"darwin,block"`
	Linux        *manifestPlatform     `manifest:"linux,block"`
	IOS          *manifestPlatform     `manifest:"ios,block"`
	Android      *manifestPlatform     `manifest:"android,block"`
	Targets      []manifestTarget      `manifest:"target,block"`
	Packages     []PackageFormat       `manifest:"package,block"`
	Profiles     []manifestProfile     `manifest:"profile,block"`
	Associations []manifestAssociation `manifest:"file_association,block"`
	Protocols    []manifestProtocol    `manifest:"protocol,block"`
	Hooks        []manifestHook        `manifest:"hook,block"`
}

type manifestHook struct {
	Phase     string    `manifest:",label" schema_label:"phase"`
	Script    *string   `manifest:"script,optional" required:"true" nonempty:"true" path:"true"`
	Directory *string   `manifest:"directory,optional" path:"true"`
	Cache     *bool     `manifest:"cache,optional" default:"false"`
	Inputs    *[]string `manifest:"inputs,optional" default:"[]" path:"true"`
	Outputs   *[]string `manifest:"outputs,optional" default:"[]" path:"true"`
}

type manifestProject struct {
	SupportedPlatforms *[]string `manifest:"supported_platforms,optional"`

	Name        *string `manifest:"name,optional" required:"true" nonempty:"true"`
	ProductName *string `manifest:"product_name,optional" required:"true" nonempty:"true"`
	Identifier  *string `manifest:"identifier,optional" required:"true" nonempty:"true"`
	Version     *string `manifest:"version,optional" required:"true" nonempty:"true"`
	Company     *string `manifest:"company,optional"`
	BinaryName  *string `manifest:"binary_name,optional" default:"$project.name.slug"`
	Icon        *string `manifest:"icon,optional"`
	Description *string `manifest:"description,optional"`
	Copyright   *string `manifest:"copyright,optional"`
	Comments    *string `manifest:"comments,optional"`
	BuildNumber *int    `manifest:"build_number,optional" default:"1"`
}

type manifestFrontend struct {
	Disabled *bool `manifest:"disabled,optional" default:"false"`

	Directory   *string            `manifest:"directory,optional" default:"frontend"`
	Install     *[]string          `manifest:"install,optional" default:"[\"npm\",\"install\"]"`
	Build       *[]string          `manifest:"build,optional" default:"[\"npm\",\"run\",\"build\"]"`
	Dev         *[]string          `manifest:"dev,optional" default:"[\"npm\",\"run\",\"dev\"]"`
	Output      *string            `manifest:"output,optional" default:"dist"`
	Environment *map[string]string `manifest:"environment,optional" default:"{}"`
	Bindings    *manifestBindings  `manifest:"bindings,block"`
}

type manifestBindings struct {
	TypeScript     *bool   `manifest:"typescript,optional" default:"true"`
	Interfaces     *bool   `manifest:"interfaces,optional" default:"true"`
	Output         *string `manifest:"output,optional" default:"bindings"`
	ModelsFilename *string `manifest:"models_filename,optional" default:"models"`
	IndexFilename  *string `manifest:"index_filename,optional" default:"index"`
	TimeType       *string `manifest:"time_type,optional" default:"string"`
}

type manifestBuild struct {
	Output        *string            `manifest:"output,optional" default:"bin"`
	Tags          *[]string          `manifest:"tags,optional" default:"[]"`
	TrimPath      *bool              `manifest:"trim_path,optional" default:"true"`
	Strip         *bool              `manifest:"strip,optional" default:"true"`
	Obfuscated    *bool              `manifest:"obfuscated,optional" default:"false"`
	GarbleArgs    *[]string          `manifest:"garble_args,optional" default:"[]"`
	LDFlags       *[]string          `manifest:"ldflags,optional" default:"[]"`
	CompilerFlags *[]string          `manifest:"compiler_flags,optional" default:"[]"`
	VCSInfo       *bool              `manifest:"vcs_info,optional" default:"false"`
	Environment   *map[string]string `manifest:"environment,optional" default:"{}"`
}

type manifestDev struct {
	Args         *[]string `manifest:"args,optional"`
	Tags         *[]string `manifest:"tags,optional" default:"[]"`
	DebounceMS   *int      `manifest:"debounce_ms,optional" default:"250"`
	LogLevel     *string   `manifest:"log_level,optional" default:"warn"`
	Watch        *[]string `manifest:"watch,optional" default:"[\"**/*.go\",\"wails.yaml\"]"`
	Exclude      *[]string `manifest:"exclude,optional" default:"[\".git\",\".wails\",\"bin\",\"node_modules\",\"frontend/dist\"]"`
	UseGitIgnore *bool     `manifest:"use_git_ignore,optional" default:"true"`
	GracePeriod  *int      `manifest:"grace_period_ms,optional" default:"1500"`
}

type manifestPlatform struct {
	ProductName      *string               `manifest:"product_name,optional" platforms:"windows,darwin,linux"`
	Identifier       *string               `manifest:"identifier,optional" platforms:"windows,darwin,linux"`
	MinimumVersion   *string               `manifest:"minimum_version,optional" platforms:"windows,darwin,linux,ios"`
	BuildNumber      *int                  `manifest:"build_number,optional" platforms:"windows,darwin,linux,ios"`
	Capabilities     *[]string             `manifest:"capabilities,optional" platforms:"windows,darwin,linux,ios"`
	Icon             *string               `manifest:"icon,optional" platforms:"windows,darwin,linux,ios,android"`
	Manifest         *string               `manifest:"manifest,optional" platforms:"windows,android"`
	AssetsCar        *string               `manifest:"assets_car,optional" platforms:"darwin,ios"`
	InfoPlist        *string               `manifest:"info_plist,optional" platforms:"darwin,ios"`
	Publisher        *string               `manifest:"publisher,optional" platforms:"windows"`
	DesktopEntry     *string               `manifest:"desktop_entry,optional" platforms:"linux"`
	BundleID         *string               `manifest:"bundle_id,optional" platforms:"ios"`
	DisplayName      *string               `manifest:"display_name,optional" platforms:"ios,android"`
	VersionName      *string               `manifest:"version_name,optional" platforms:"android"`
	VersionCode      *int                  `manifest:"version_code,optional" platforms:"android"`
	MinimumSDK       *int                  `manifest:"minimum_sdk,optional" platforms:"android"`
	TargetSDK        *int                  `manifest:"target_sdk,optional" platforms:"android"`
	ApplicationID    *string               `manifest:"application_id,optional" platforms:"android"`
	Company          *string               `manifest:"company,optional" platforms:"ios,android"`
	Comments         *string               `manifest:"comments,optional" platforms:"ios,android"`
	CFBundleIconName *string               `manifest:"cf_bundle_icon_name,optional" platforms:"darwin,ios"`
	BackgroundModes  *[]string             `manifest:"background_modes,optional" platforms:"ios"`
	Signing          *manifestSigning      `manifest:"signing,block"`
	Notarization     *manifestNotarization `manifest:"notarization,block" platforms:"darwin"`
}

type manifestSigning struct {
	Credential          *string `manifest:"credential,optional"`
	Identity            *string `manifest:"identity,optional"`
	Certificate         *string `manifest:"certificate,optional"`
	Thumbprint          *string `manifest:"thumbprint,optional"`
	TimestampServer     *string `manifest:"timestamp_server,optional"`
	Entitlements        *string `manifest:"entitlements,optional"`
	ProvisioningProfile *string `manifest:"provisioning_profile,optional"`
	KeyAlias            *string `manifest:"key_alias,optional"`
}

type manifestNotarization struct {
	Credential *string `manifest:"credential,optional"`
}

type manifestRun struct {
	Tags        *[]string          `manifest:"tags,optional"`
	Args        *[]string          `manifest:"args,optional"`
	Environment *map[string]string `manifest:"environment,optional"`
}

type manifestTargetDev struct {
	Args *[]string `manifest:"args,optional"`
}

type manifestTarget struct {
	Dev *manifestTargetDev `manifest:"dev,block"`
	Run *manifestRun       `manifest:"run,block"`

	Name           string             `manifest:",label" schema_label:"target"`
	Tags           *[]string          `manifest:"tags,optional"`
	MinimumVersion *string            `manifest:"minimum_version,optional"`
	BuildNumber    *int               `manifest:"build_number,optional"`
	Toolchain      *string            `manifest:"toolchain,optional"`
	Environment    *map[string]string `manifest:"environment,optional"`
	LDFlags        *[]string          `manifest:"ldflags,optional"`
	CompilerFlags  *[]string          `manifest:"compiler_flags,optional"`
	GarbleArgs     *[]string          `manifest:"garble_args,optional"`
	Obfuscated     *bool              `manifest:"obfuscated,optional"`
}

type manifestProfile struct {
	Name    string                  `manifest:",label" schema_label:"profile"`
	Targets []manifestProfileTarget `manifest:"target,block" required:"true"`
}

type manifestProfileTarget struct {
	Name        string    `manifest:",label" schema_label:"target"`
	Formats     *[]string `manifest:"formats,optional"`
	Sign        *bool     `manifest:"sign,optional"`
	Notarize    *bool     `manifest:"notarize,optional"`
	Destination *string   `manifest:"destination,optional"`
}

type manifestAssociation struct {
	Label       string    `manifest:",label" schema_label:"association"`
	Extensions  *[]string `manifest:"extensions,optional" required:"true"`
	Name        *string   `manifest:"name,optional"`
	Description *string   `manifest:"description,optional"`
	Icon        *string   `manifest:"icon,optional"`
	Role        *string   `manifest:"role,optional"`
	MIMEType    *string   `manifest:"mime_type,optional"`
	Platforms   *[]string `manifest:"platforms,optional"`
}

type manifestProtocol struct {
	Scheme      string    `manifest:",label" schema_label:"scheme"`
	Description *string   `manifest:"description,optional"`
	Platforms   *[]string `manifest:"platforms,optional"`
}

func decodeYAML(root, filename string, src []byte, selectedProfile string) (*Loaded, error) {
	raw, origins, err := decodeManifestSchema(src, filename)
	if err != nil {
		return nil, err
	}
	doc, err := documentFromManifest(raw)
	if err != nil {
		return nil, attachValidationRange(err, origins)
	}
	config := configFromDocument(root, selectedProfile, doc)
	for field, origin := range origins {
		config.Origins[field] = origin
	}
	if selectedProfile != "" {
		if selectedProfile == "default" || !slugPattern.MatchString(selectedProfile) {
			return nil, fmt.Errorf("profile name must be a lowercase slug and cannot be default")
		}
		profile, exists := doc.Profiles[selectedProfile]
		if !exists {
			return nil, fmt.Errorf("profile %q is not defined", selectedProfile)
		}
		config.Selected = profile
	}
	if err := validateConfig(config); err != nil {
		return nil, attachValidationRange(err, config.Origins)
	}
	return &Loaded{Path: filename, Raw: src, Document: doc, Config: config}, nil
}

// LoadFile validates an explicitly named manifest-like YAML file. It is used
// by migration to validate an inactive draft before the atomic cutover rename.
func LoadFile(root, filename, profile string) (*Loaded, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	return decodeYAML(absRoot, filename, src, profile)
}

func documentFromManifest(raw manifestDocument) (Document, error) {
	if raw.Project == nil {
		return Document{}, fieldValidationError("project", "block is required")
	}
	project := Project{}
	setStrings(&project.SupportedPlatforms, raw.Project.SupportedPlatforms)
	setString(&project.Name, raw.Project.Name)
	setString(&project.ProductName, raw.Project.ProductName)
	setString(&project.Identifier, raw.Project.Identifier)
	setString(&project.Version, raw.Project.Version)
	setString(&project.CompanyName, raw.Project.Company)
	setString(&project.BinaryName, raw.Project.BinaryName)
	setString(&project.Icon, raw.Project.Icon)
	setString(&project.Description, raw.Project.Description)
	setString(&project.Copyright, raw.Project.Copyright)
	setString(&project.Comments, raw.Project.Comments)
	setInt(&project.BuildNumber, raw.Project.BuildNumber)
	if err := validateProject(project); err != nil {
		return Document{}, err
	}
	doc := defaults(project)
	applyRun(&doc.Run, raw.Run)
	if raw.Frontend != nil {
		applyFrontend(&doc.Frontend, raw.Frontend)
	}
	if raw.Build != nil {
		applyBuild(&doc.Build, raw.Build)
	}
	if raw.Dev != nil {
		applyDev(&doc.Dev, raw.Dev)
	}
	for _, item := range []struct {
		name string
		body *manifestPlatform
	}{{"windows", raw.Windows}, {"darwin", raw.Darwin}, {"linux", raw.Linux}, {"ios", raw.IOS}, {"android", raw.Android}} {
		if item.body != nil {
			applyPlatform(&doc, item.name, item.body)
		}
	}
	// OS-wide defaults are applied before architecture overrides regardless of file order.
	sort.SliceStable(raw.Targets, func(i, j int) bool {
		return !strings.Contains(raw.Targets[i].Name, "/") && strings.Contains(raw.Targets[j].Name, "/")
	})
	seenTargets := map[string]bool{}
	for _, target := range raw.Targets {
		field := `targets[` + strconv.Quote(target.Name) + `]`
		if seenTargets[target.Name] {
			return Document{}, fieldValidationError(field, "duplicate target")
		}
		seenTargets[target.Name] = true
		if err := applyTarget(&doc.Targets, target); err != nil {
			return Document{}, fieldValidationCause(field, err, "%v", err)
		}
	}
	seenPackages := map[string]bool{}
	for _, pkg := range raw.Packages {
		field := `packages[` + strconv.Quote(pkg.Format) + `]`
		if seenPackages[pkg.Format] {
			return Document{}, fieldValidationError(field, "duplicate package block")
		}
		seenPackages[pkg.Format] = true
		if err := applyPackage(&doc.Package, pkg); err != nil {
			return Document{}, fieldValidationCause(field, err, "%v", err)
		}
	}
	for _, association := range raw.Associations {
		field := `file_associations[` + strconv.Quote(association.Label) + `]`
		if association.Extensions == nil || len(*association.Extensions) == 0 {
			return Document{}, fieldValidationError(field+".extensions", "requires at least one extension")
		}
		entry := Association{Extensions: append([]string(nil), (*association.Extensions)...)}
		setString(&entry.Name, association.Name)
		if entry.Name == "" {
			entry.Name = association.Label
		}
		setString(&entry.Description, association.Description)
		setString(&entry.Icon, association.Icon)
		setString(&entry.Role, association.Role)
		setString(&entry.MIMEType, association.MIMEType)
		setStrings(&entry.Platforms, association.Platforms)
		doc.Associations = append(doc.Associations, entry)
	}
	seenProfiles := map[string]bool{}
	for _, rawProfile := range raw.Profiles {
		profileField := `profiles[` + strconv.Quote(rawProfile.Name) + `]`
		if rawProfile.Name == "" || rawProfile.Name == "default" || !slugPattern.MatchString(rawProfile.Name) {
			return Document{}, fieldValidationError(profileField, "name must be a lowercase slug and cannot be default")
		}
		if seenProfiles[rawProfile.Name] {
			return Document{}, fieldValidationError(profileField, "duplicate profile")
		}
		seenProfiles[rawProfile.Name] = true
		if len(rawProfile.Targets) == 0 {
			return Document{}, fieldValidationError(profileField, "requires at least one target")
		}
		profile := Profile{Name: rawProfile.Name}
		seen := map[string]bool{}
		for _, rawTarget := range rawProfile.Targets {
			targetField := profileField + `.targets[` + strconv.Quote(rawTarget.Name) + `]`
			if seen[rawTarget.Name] {
				return Document{}, fieldValidationError(targetField, "duplicate target")
			}
			seen[rawTarget.Name] = true
			platform, arch, err := parseTargetName(rawTarget.Name)
			if err != nil {
				return Document{}, fieldValidationCause(targetField, err, "%v", err)
			}
			capability, _ := buildinfo.LookupTarget(platform, arch)
			entry := ProfileTarget{Target: rawTarget.Name}
			setStrings(&entry.Formats, rawTarget.Formats)
			setBool(&entry.Sign, rawTarget.Sign)
			setBool(&entry.Notarize, rawTarget.Notarize)
			setString(&entry.Destination, rawTarget.Destination)
			if err := validateProfileTarget(profile.Name, entry, capability); err != nil {
				return Document{}, err
			}
			profile.Targets = append(profile.Targets, entry)
		}
		doc.Profiles[profile.Name] = profile
	}
	seenProtocols := map[string]bool{}
	for _, protocol := range raw.Protocols {
		field := `protocols[` + strconv.Quote(protocol.Scheme) + `]`
		if protocol.Scheme == "" {
			return Document{}, fieldValidationError(field, "label cannot be empty")
		}
		if seenProtocols[protocol.Scheme] {
			return Document{}, fieldValidationError(field, "duplicate protocol block")
		}
		seenProtocols[protocol.Scheme] = true
		entry := Protocol{Scheme: protocol.Scheme}
		setString(&entry.Description, protocol.Description)
		setStrings(&entry.Platforms, protocol.Platforms)
		doc.Protocols = append(doc.Protocols, entry)
	}
	seenHooks := map[HookPhase]bool{}
	for _, rawHook := range raw.Hooks {
		phase := HookPhase(rawHook.Phase)
		field := `hooks[` + strconv.Quote(rawHook.Phase) + `]`
		if !containsHookPhase(phase) {
			return Document{}, fieldValidationError(field, "phase is not supported")
		}
		if seenHooks[phase] {
			return Document{}, fieldValidationError(field, "duplicate hook block")
		}
		seenHooks[phase] = true
		hook := Hook{}
		setString(&hook.Script, rawHook.Script)
		setString(&hook.Directory, rawHook.Directory)
		setBool(&hook.Cache, rawHook.Cache)
		setStrings(&hook.Inputs, rawHook.Inputs)
		setStrings(&hook.Outputs, rawHook.Outputs)
		if doc.Hooks == nil {
			doc.Hooks = make(map[HookPhase]Hook)
		}
		doc.Hooks[phase] = hook
	}
	return doc, nil
}

func containsHookPhase(want HookPhase) bool {
	for _, phase := range HookPhases {
		if phase == want {
			return true
		}
	}
	return false
}

func applyFrontend(target *Frontend, raw *manifestFrontend) {
	setBool(&target.Disabled, raw.Disabled)
	setString(&target.Directory, raw.Directory)
	setString(&target.OutputDirectory, raw.Output)
	setStrings(&target.Install, raw.Install)
	setStrings(&target.Build, raw.Build)
	setStrings(&target.Dev, raw.Dev)
	setStringMap(&target.Environment, raw.Environment)
	if raw.Bindings != nil {
		applyBindings(&target.Bindings, raw.Bindings)
	}
	if len(target.Install) > 0 {
		target.PackageManager = target.Install[0]
	}
	if len(target.Build) > 0 && target.PackageManager == "" {
		target.PackageManager = target.Build[0]
	}
	if len(target.Dev) > 0 && target.PackageManager == "" {
		target.PackageManager = target.Dev[0]
	}
}

func applyBindings(target *Bindings, raw *manifestBindings) {
	setBool(&target.TypeScript, raw.TypeScript)
	setBool(&target.Interfaces, raw.Interfaces)
	setString(&target.OutputDirectory, raw.Output)
	setString(&target.ModelsFilename, raw.ModelsFilename)
	setString(&target.IndexFilename, raw.IndexFilename)
	setString(&target.TimeType, raw.TimeType)
}

func applyBuild(target *Build, raw *manifestBuild) {
	setString(&target.OutputDirectory, raw.Output)
	setBool(&target.TrimPath, raw.TrimPath)
	setBool(&target.Strip, raw.Strip)
	setBool(&target.Obfuscation, raw.Obfuscated)
	setStrings(&target.Go.Tags, raw.Tags)
	setStrings(&target.Go.GarbleArgs, raw.GarbleArgs)
	setStrings(&target.Go.LinkerFlags, raw.LDFlags)
	setStrings(&target.Go.CompilerFlags, raw.CompilerFlags)
	setBool(&target.VCSInfo, raw.VCSInfo)
	setStringMap(&target.Environment, raw.Environment)
}

func applyDev(target *Dev, raw *manifestDev) {
	if raw.Args != nil {
		target.Args = append([]string{}, (*raw.Args)...)
		target.ArgsSet = true
	}
	setStrings(&target.Tags, raw.Tags)
	setInt(&target.DebounceMS, raw.DebounceMS)
	setInt(&target.GracePeriodMS, raw.GracePeriod)
	setString(&target.LogLevel, raw.LogLevel)
	setStrings(&target.Watch, raw.Watch)
	setStrings(&target.Exclude, raw.Exclude)
	setBool(&target.UseGitIgnore, raw.UseGitIgnore)
}

func applyPlatform(doc *Document, name string, raw *manifestPlatform) {
	platform := platformByName(&doc.Targets, name)
	setString(&platform.ProductName, raw.ProductName)
	setString(&platform.Identifier, raw.Identifier)
	setString(&platform.MinimumVersion, raw.MinimumVersion)
	setInt(&platform.BuildNumber, raw.BuildNumber)
	setStrings(&platform.Capabilities, raw.Capabilities)
	setString(&platform.Icon, raw.Icon)
	setString(&platform.Manifest, raw.Manifest)
	setString(&platform.AssetsCar, raw.AssetsCar)
	setString(&platform.InfoPlist, raw.InfoPlist)
	setString(&platform.Publisher, raw.Publisher)
	setString(&platform.DesktopEntry, raw.DesktopEntry)
	setString(&platform.Identifier, raw.BundleID)
	setString(&platform.Identifier, raw.ApplicationID)
	setString(&platform.ProductName, raw.DisplayName)
	setString(&platform.CompanyName, raw.Company)
	setString(&platform.Comments, raw.Comments)
	setString(&platform.CFBundleIconName, raw.CFBundleIconName)
	setStrings(&platform.BackgroundModes, raw.BackgroundModes)
	setString(&platform.VersionName, raw.VersionName)
	setInt(&platform.VersionCode, raw.VersionCode)
	setInt(&platform.MinimumSDK, raw.MinimumSDK)
	setInt(&platform.TargetSDK, raw.TargetSDK)
	if raw.Signing != nil {
		signing := signingByName(&doc.Signing, name)
		signing.Enabled = true
		setString(&signing.Credential, raw.Signing.Credential)
		setString(&signing.Identity, raw.Signing.Identity)
		setString(&signing.Certificate, raw.Signing.Certificate)
		setString(&signing.Thumbprint, raw.Signing.Thumbprint)
		setString(&signing.TimestampServer, raw.Signing.TimestampServer)
		setString(&signing.Entitlements, raw.Signing.Entitlements)
		setString(&signing.ProvisioningProfile, raw.Signing.ProvisioningProfile)
		setString(&signing.KeyAlias, raw.Signing.KeyAlias)
	}
	if raw.Notarization != nil {
		signing := signingByName(&doc.Signing, name)
		signing.Enabled = true
		signing.Notarize = true
		setString(&signing.NotarizationCredential, raw.Notarization.Credential)
	}
}

func applyRun(target *Run, raw *manifestRun) {
	if raw == nil {
		return
	}
	if raw.Tags != nil {
		target.Tags = appendUniqueRunTags(target.Tags, (*raw.Tags)...)
	}
	if raw.Args != nil {
		target.Args = append([]string{}, (*raw.Args)...)
		target.ArgsSet = true
	}
	if raw.Environment != nil {
		if target.Environment == nil {
			target.Environment = map[string]string{}
		}
		for k, v := range *raw.Environment {
			target.Environment[k] = v
		}
	}
}

func appendUniqueRunTags(tags []string, extra ...string) []string {
	result := append([]string(nil), tags...)
	for _, tag := range extra {
		found := false
		for _, existing := range result {
			if existing == tag {
				found = true
				break
			}
		}
		if !found {
			result = append(result, tag)
		}
	}
	return result
}

func applyTarget(targets *Targets, raw manifestTarget) error {
	if !strings.Contains(raw.Name, "/") {
		matched := false
		for _, name := range buildinfo.SupportedTargetNames() {
			if strings.HasPrefix(name, raw.Name+"/") {
				matched = true
				item := raw
				item.Name = name
				if err := applyTarget(targets, item); err != nil {
					return err
				}
			}
		}
		if !matched {
			return fmt.Errorf("unsupported platform %q", raw.Name)
		}
		return nil
	}

	platform, arch, err := parseTargetName(raw.Name)
	if err != nil {
		return err
	}
	target := targetByName(platformByName(targets, platform), arch)
	// parseTargetName accepts only registry targets, and every registry
	// architecture has a concrete slot in Platform.
	if raw.Tags != nil {
		target.Tags = appendUniqueRunTags(target.Tags, (*raw.Tags)...)
	}
	applyRun(&target.Run, raw.Run)
	if raw.Dev != nil && raw.Dev.Args != nil {
		target.Dev.Args = append([]string{}, (*raw.Dev.Args)...)
		target.Dev.ArgsSet = true
	}
	setString(&target.MinimumVersion, raw.MinimumVersion)
	setInt(&target.BuildNumber, raw.BuildNumber)
	setString(&target.Toolchain, raw.Toolchain)
	if raw.Environment != nil {
		if target.Environment == nil {
			target.Environment = map[string]string{}
		}
		for key, value := range *raw.Environment {
			target.Environment[key] = value
		}
	}
	setStrings(&target.LinkerFlags, raw.LDFlags)
	setStrings(&target.CompilerFlags, raw.CompilerFlags)
	setStrings(&target.GarbleArgs, raw.GarbleArgs)
	setBool(&target.Obfuscated, raw.Obfuscated)
	target.ObfuscatedSet = target.ObfuscatedSet || raw.Obfuscated != nil
	return nil
}

func applyPackage(packages *Packages, raw PackageFormat) error {
	switch raw.Format {
	case "nsis":
		packages.Windows.NSIS = NSISPackage{Template: raw.Template, InstallScope: raw.InstallScope}
	case "msix":
		packages.Windows.MSIX = MSIXPackage{Publisher: raw.Publisher, Manifest: raw.Manifest}
	case "dmg":
		packages.Darwin.DMG = DMGPackage{Template: raw.Template, Background: raw.Background, VolumeIcon: raw.VolumeIcon, FileIcon: raw.FileIcon, Files: cloneStringMapValue(raw.Files), WindowWidth: raw.WindowWidth, WindowHeight: raw.WindowHeight}
	case "appimage":
		packages.Linux.AppImage = AppImagePackage{Icon: raw.Icon, DesktopEntry: raw.DesktopEntry, Categories: append([]string(nil), raw.Categories...)}
	case "deb", "rpm", "archlinux":
		value := LinuxPackage{Template: raw.Template, Maintainer: raw.Maintainer, Section: raw.Section, Dependencies: append([]string(nil), raw.Dependencies...), PreInstall: raw.PreInstall, PostInstall: raw.PostInstall, PreRemove: raw.PreRemove, PostRemove: raw.PostRemove}
		switch raw.Format {
		case "deb":
			packages.Linux.Deb = value
		case "rpm":
			packages.Linux.RPM = value
		case "archlinux":
			packages.Linux.ArchLinux = value
		}
	default:
		return fmt.Errorf("package format %q has no configurable package block", raw.Format)
	}
	return nil
}

func cloneStringMapValue(source map[string]string) map[string]string {
	if source == nil {
		return nil
	}
	result := make(map[string]string, len(source))
	for name, value := range source {
		result[name] = value
	}
	return result
}

// ResolvePackageFormat returns the closed planner representation for one
// compatible platform/format pair. Resolved Config remains compact by storing
// only each format's own fields; callers never need to know that layout.
func ResolvePackageFormat(packages Packages, platform, format string) (PackageFormat, error) {
	switch platform + "/" + format {
	case "windows/nsis":
		value := packages.Windows.NSIS
		return PackageFormat{Template: value.Template, InstallScope: value.InstallScope}, nil
	case "windows/msix":
		value := packages.Windows.MSIX
		return PackageFormat{Publisher: value.Publisher, Manifest: value.Manifest}, nil
	case "darwin/app", "ios/app", "ios/ipa", "android/apk", "android/aab":
		return PackageFormat{}, nil
	case "darwin/dmg":
		value := packages.Darwin.DMG
		return PackageFormat{Template: value.Template, Background: value.Background, VolumeIcon: value.VolumeIcon, FileIcon: value.FileIcon, Files: cloneStringMapValue(value.Files), WindowWidth: value.WindowWidth, WindowHeight: value.WindowHeight}, nil
	case "linux/appimage":
		value := packages.Linux.AppImage
		return PackageFormat{Icon: value.Icon, DesktopEntry: value.DesktopEntry, Categories: append([]string(nil), value.Categories...)}, nil
	case "linux/deb", "linux/rpm", "linux/archlinux":
		var value LinuxPackage
		switch format {
		case "deb":
			value = packages.Linux.Deb
		case "rpm":
			value = packages.Linux.RPM
		case "archlinux":
			value = packages.Linux.ArchLinux
		}
		return PackageFormat{Template: value.Template, Maintainer: value.Maintainer, Section: value.Section, Dependencies: append([]string(nil), value.Dependencies...), PreInstall: value.PreInstall, PostInstall: value.PostInstall, PreRemove: value.PreRemove, PostRemove: value.PostRemove}, nil
	default:
		return PackageFormat{}, fmt.Errorf("package format %q is not supported for %s", format, platform)
	}
}

func parseTargetName(value string) (string, string, error) {
	parts := strings.Split(value, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("target must be platform/architecture, got %q", value)
	}
	if _, ok := buildinfo.LookupTarget(parts[0], parts[1]); !ok {
		return "", "", fmt.Errorf("unsupported target %q; supported targets: %v", value, buildinfo.SupportedTargetNames())
	}
	return parts[0], parts[1], nil
}

func validateProfileTarget(profile string, target ProfileTarget, capability buildinfo.TargetCapability) error {
	field := `profiles[` + strconv.Quote(profile) + `].targets[` + strconv.Quote(target.Target) + `]`
	seenFormats := make(map[string]bool, len(target.Formats))
	for _, format := range target.Formats {
		if seenFormats[format] {
			return fieldValidationError(field+".formats", "contains duplicate format %q", format)
		}
		seenFormats[format] = true
		formatCapability, ok := buildinfo.LookupFormat(format)
		if !ok || !formatCapability.Production || !capability.SupportsFormat(format, false) {
			return fieldValidationError(field+".formats", "format %q is not a production format for %s", format, target.Target)
		}
	}
	if target.Destination != "" {
		if capability.Target.OS != "ios" {
			return fieldValidationError(field+".destination", "is only valid for iOS targets")
		}
		if target.Destination != "simulator" && target.Destination != "device" {
			return fieldValidationError(field+".destination", "must be simulator or device")
		}
	}
	if capability.Target.OS == "ios" && target.Destination == "" {
		return fieldValidationError(field+".destination", "requires destination = %q or %q", "simulator", "device")
	}
	if seenFormats["ipa"] && target.Destination != "device" {
		return fieldValidationError(field+".destination", "IPA requires %q", "device")
	}
	if target.Notarize && capability.Target.OS != "darwin" {
		return fieldValidationError(field+".notarize", "is only valid for darwin targets")
	}
	if target.Notarize && !target.Sign {
		return fieldValidationError(field+".sign", "must be signed before notarization")
	}
	return nil
}

func platformByName(targets *Targets, name string) *Platform {
	switch name {
	case "windows":
		return &targets.Windows
	case "darwin":
		return &targets.Darwin
	case "linux":
		return &targets.Linux
	case "ios":
		return &targets.IOS
	default:
		return &targets.Android
	}
}

func signingByName(signing *Signing, name string) *SigningPlatform {
	switch name {
	case "windows":
		return &signing.Windows
	case "darwin":
		return &signing.Darwin
	case "linux":
		return &signing.Linux
	case "ios":
		return &signing.IOS
	default:
		return &signing.Android
	}
}

func targetByName(platform *Platform, arch string) *Target {
	switch arch {
	case "amd64":
		return &platform.AMD64
	case "arm64":
		return &platform.ARM64
	case "arm":
		return &platform.ARM
	case "386":
		return &platform.X86
	case "universal":
		return &platform.Universal
	default:
		return nil
	}
}

func setString(target *string, value *string) {
	if value != nil {
		*target = *value
	}
}
func setInt(target *int, value *int) {
	if value != nil {
		*target = *value
	}
}
func setBool(target *bool, value *bool) {
	if value != nil {
		*target = *value
	}
}
func setStrings(target *[]string, value *[]string) {
	if value != nil {
		*target = append([]string(nil), (*value)...)
	}
}
func setStringMap(target *map[string]string, value *map[string]string) {
	if value != nil {
		*target = make(map[string]string, len(*value))
		for key, item := range *value {
			(*target)[key] = item
		}
	}
}

func platformAttributeAllowed(platform, name string) bool {
	attribute := manifestSchema.blocks[platform].node.attributes[name]
	return schemaFieldAllowed(platform, attribute.platformMask)
}

func platformBlockAllowed(platform, name string) bool {
	child := manifestSchema.blocks[platform].node.blocks[name]
	return schemaFieldAllowed(platform, child.platformMask)
}

func packageAttributeAllowed(format, name string) bool {
	descriptor := schemaNodesByType[reflect.TypeOf(PackageFormat{})].attributes[name]
	return schemaFormatNameAllowed(format, descriptor.formatMask)
}

func packageConfigured(format PackageFormat) bool {
	return format.Template != "" || format.InstallScope != "" || format.Publisher != "" || format.Manifest != "" || format.Background != "" || format.VolumeIcon != "" || format.FileIcon != "" || len(format.Files) > 0 || format.WindowWidth != 0 || format.WindowHeight != 0 || format.Icon != "" || format.DesktopEntry != "" || len(format.Categories) > 0 || format.Maintainer != "" || format.Section != "" || len(format.Dependencies) > 0 || format.PreInstall != "" || format.PostInstall != "" || format.PreRemove != "" || format.PostRemove != ""
}

func platformConfigured(platform Platform) bool {
	return platform.ProductName != "" || platform.Identifier != "" || platform.MinimumVersion != "" || platform.BuildNumber != 0 || len(platform.Capabilities) > 0 || platform.Icon != "" || platform.Manifest != "" || platform.AssetsCar != "" || platform.InfoPlist != "" || platform.Publisher != "" || platform.DesktopEntry != "" || platform.CompanyName != "" || platform.Comments != "" || platform.CFBundleIconName != "" || len(platform.BackgroundModes) > 0 || platform.VersionName != "" || platform.VersionCode != 0 || platform.MinimumSDK != 0 || platform.TargetSDK != 0
}

func manifestValueWasExplicit(origins map[string]Origin, path string) bool {
	return origins[path].Kind == OriginManifest
}

func sortedProfiles(profiles map[string]Profile) []Profile {
	keys := make([]string, 0, len(profiles))
	for key := range profiles {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]Profile, 0, len(keys))
	for _, key := range keys {
		result = append(result, profiles[key])
	}
	return result
}

func runConfigured(run Run) bool { return len(run.Tags) > 0 || run.ArgsSet || len(run.Environment) > 0 }

// RunForTarget returns independent resolved launch defaults for one target.
func (config Config) RunForTarget(goos, goarch string) (Run, error) {
	if _, _, err := parseTargetName(goos + "/" + goarch); err != nil {
		return Run{}, err
	}
	target := targetByName(platformByName(&config.Targets, goos), goarch).Run
	result := config.Run
	result.Tags = appendUniqueRunTags(config.Run.Tags, target.Tags...)
	result.Args = append([]string(nil), config.Run.Args...)
	if target.ArgsSet {
		result.Args = append([]string{}, target.Args...)
		result.ArgsSet = true
	}
	result.Environment = cloneStringMapValue(config.Run.Environment)
	if result.Environment == nil {
		result.Environment = map[string]string{}
	}
	for k, v := range target.Environment {
		result.Environment[k] = v
	}
	return result, nil
}

// DevArgsForTarget resolves application arguments independently of run/build policy.
func (config Config) DevArgsForTarget(goos, goarch string) ([]string, error) {
	if _, _, err := parseTargetName(goos + "/" + goarch); err != nil {
		return nil, err
	}
	target := targetByName(platformByName(&config.Targets, goos), goarch).Dev
	args := config.Dev.Args
	if target.ArgsSet {
		args = target.Args
	}
	return append([]string{}, args...), nil
}
