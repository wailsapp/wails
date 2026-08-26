package manifest

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/wailsapp/wails/v3/internal/wake/buildinfo"
	"github.com/zclconf/go-cty/cty/ctystrings"
)

// hclDocument is intentionally a closed schema. Keep the HCL surface small:
// this is declarative build intent, never a programmable pipeline.
type hclDocument struct {
	Version      int              `hcl:"version"`
	Project      *hclProject      `hcl:"project,block" required:"true"`
	Frontend     *hclFrontend     `hcl:"frontend,block"`
	Build        *hclBuild        `hcl:"build,block"`
	Dev          *hclDev          `hcl:"dev,block"`
	Windows      *hclPlatform     `hcl:"windows,block"`
	Darwin       *hclPlatform     `hcl:"darwin,block"`
	Linux        *hclPlatform     `hcl:"linux,block"`
	IOS          *hclPlatform     `hcl:"ios,block"`
	Android      *hclPlatform     `hcl:"android,block"`
	Targets      []hclTarget      `hcl:"target,block"`
	Packages     []PackageFormat  `hcl:"package,block"`
	Profiles     []hclProfile     `hcl:"profile,block"`
	Associations []hclAssociation `hcl:"file_association,block"`
	Protocols    []hclProtocol    `hcl:"protocol,block"`
	Hooks        []hclHook        `hcl:"hook,block"`
}

type hclHook struct {
	Phase     string    `hcl:",label" schema_label:"phase"`
	Script    *string   `hcl:"script,optional" required:"true" nonempty:"true" path:"true"`
	Directory *string   `hcl:"directory,optional" path:"true"`
	Cache     *bool     `hcl:"cache,optional" default:"false"`
	Inputs    *[]string `hcl:"inputs,optional" default:"[]" path:"true"`
	Outputs   *[]string `hcl:"outputs,optional" default:"[]" path:"true"`
}

type hclProject struct {
	Name        *string `hcl:"name,optional" required:"true" nonempty:"true"`
	ProductName *string `hcl:"product_name,optional" required:"true" nonempty:"true"`
	Identifier  *string `hcl:"identifier,optional" required:"true" nonempty:"true"`
	Version     *string `hcl:"version,optional" required:"true" nonempty:"true"`
	Company     *string `hcl:"company,optional"`
	BinaryName  *string `hcl:"binary_name,optional" default:"$project.name.slug"`
	Icon        *string `hcl:"icon,optional"`
	Description *string `hcl:"description,optional"`
	Copyright   *string `hcl:"copyright,optional"`
	Comments    *string `hcl:"comments,optional"`
	BuildNumber *int    `hcl:"build_number,optional" default:"1"`
}

type hclFrontend struct {
	Directory   *string            `hcl:"directory,optional" default:"frontend"`
	Install     *[]string          `hcl:"install,optional" default:"[\"npm\",\"install\"]"`
	Build       *[]string          `hcl:"build,optional" default:"[\"npm\",\"run\",\"build\"]"`
	Dev         *[]string          `hcl:"dev,optional" default:"[\"npm\",\"run\",\"dev\"]"`
	Output      *string            `hcl:"output,optional" default:"dist"`
	Environment *map[string]string `hcl:"environment,optional" default:"{}"`
	Bindings    *hclBindings       `hcl:"bindings,block"`
}

type hclBindings struct {
	TypeScript     *bool   `hcl:"typescript,optional" default:"true"`
	Interfaces     *bool   `hcl:"interfaces,optional" default:"true"`
	Output         *string `hcl:"output,optional" default:"bindings"`
	ModelsFilename *string `hcl:"models_filename,optional" default:"models"`
	IndexFilename  *string `hcl:"index_filename,optional" default:"index"`
	TimeType       *string `hcl:"time_type,optional" default:"string"`
}

type hclBuild struct {
	Output        *string            `hcl:"output,optional" default:"bin"`
	Tags          *[]string          `hcl:"tags,optional" default:"[]"`
	TrimPath      *bool              `hcl:"trim_path,optional" default:"true"`
	Strip         *bool              `hcl:"strip,optional" default:"true"`
	Obfuscated    *bool              `hcl:"obfuscated,optional" default:"false"`
	GarbleArgs    *[]string          `hcl:"garble_args,optional" default:"[]"`
	LDFlags       *[]string          `hcl:"ldflags,optional" default:"[]"`
	CompilerFlags *[]string          `hcl:"compiler_flags,optional" default:"[]"`
	VCSInfo       *bool              `hcl:"vcs_info,optional" default:"false"`
	Environment   *map[string]string `hcl:"environment,optional" default:"{}"`
}

type hclDev struct {
	Tags         *[]string `hcl:"tags,optional" default:"[]"`
	DebounceMS   *int      `hcl:"debounce_ms,optional" default:"250"`
	LogLevel     *string   `hcl:"log_level,optional" default:"warn"`
	Watch        *[]string `hcl:"watch,optional" default:"[\"**/*.go\",\"wails.hcl\"]"`
	Exclude      *[]string `hcl:"exclude,optional" default:"[\".git\",\".wails\",\"bin\",\"node_modules\",\"frontend/dist\"]"`
	UseGitIgnore *bool     `hcl:"use_git_ignore,optional" default:"true"`
	GracePeriod  *int      `hcl:"grace_period_ms,optional" default:"1500"`
}

type hclPlatform struct {
	ProductName      *string          `hcl:"product_name,optional" platforms:"windows,darwin,linux"`
	Identifier       *string          `hcl:"identifier,optional" platforms:"windows,darwin,linux"`
	MinimumVersion   *string          `hcl:"minimum_version,optional" platforms:"windows,darwin,linux,ios"`
	BuildNumber      *int             `hcl:"build_number,optional" platforms:"windows,darwin,linux,ios"`
	Capabilities     *[]string        `hcl:"capabilities,optional" platforms:"windows,darwin,linux,ios"`
	Icon             *string          `hcl:"icon,optional" platforms:"windows,darwin,linux,ios,android"`
	Manifest         *string          `hcl:"manifest,optional" platforms:"windows,android"`
	AssetsCar        *string          `hcl:"assets_car,optional" platforms:"darwin,ios"`
	InfoPlist        *string          `hcl:"info_plist,optional" platforms:"darwin,ios"`
	Publisher        *string          `hcl:"publisher,optional" platforms:"windows"`
	DesktopEntry     *string          `hcl:"desktop_entry,optional" platforms:"linux"`
	BundleID         *string          `hcl:"bundle_id,optional" platforms:"ios"`
	DisplayName      *string          `hcl:"display_name,optional" platforms:"ios,android"`
	VersionName      *string          `hcl:"version_name,optional" platforms:"android"`
	VersionCode      *int             `hcl:"version_code,optional" platforms:"android"`
	MinimumSDK       *int             `hcl:"minimum_sdk,optional" platforms:"android"`
	TargetSDK        *int             `hcl:"target_sdk,optional" platforms:"android"`
	ApplicationID    *string          `hcl:"application_id,optional" platforms:"android"`
	Company          *string          `hcl:"company,optional" platforms:"ios,android"`
	Comments         *string          `hcl:"comments,optional" platforms:"ios,android"`
	CFBundleIconName *string          `hcl:"cf_bundle_icon_name,optional" platforms:"darwin,ios"`
	BackgroundModes  *[]string        `hcl:"background_modes,optional" platforms:"ios"`
	Signing          *hclSigning      `hcl:"signing,block"`
	Notarization     *hclNotarization `hcl:"notarization,block" platforms:"darwin"`
}

type hclSigning struct {
	Credential          *string `hcl:"credential,optional"`
	Identity            *string `hcl:"identity,optional"`
	Certificate         *string `hcl:"certificate,optional"`
	Thumbprint          *string `hcl:"thumbprint,optional"`
	TimestampServer     *string `hcl:"timestamp_server,optional"`
	Entitlements        *string `hcl:"entitlements,optional"`
	ProvisioningProfile *string `hcl:"provisioning_profile,optional"`
	KeyAlias            *string `hcl:"key_alias,optional"`
}

type hclNotarization struct {
	Credential *string `hcl:"credential,optional"`
}

type hclTarget struct {
	Name           string             `hcl:",label" schema_label:"target"`
	Tags           *[]string          `hcl:"tags,optional"`
	MinimumVersion *string            `hcl:"minimum_version,optional"`
	BuildNumber    *int               `hcl:"build_number,optional"`
	Toolchain      *string            `hcl:"toolchain,optional"`
	Environment    *map[string]string `hcl:"environment,optional"`
	LDFlags        *[]string          `hcl:"ldflags,optional"`
	CompilerFlags  *[]string          `hcl:"compiler_flags,optional"`
	GarbleArgs     *[]string          `hcl:"garble_args,optional"`
	Obfuscated     *bool              `hcl:"obfuscated,optional"`
}

type hclProfile struct {
	Name    string             `hcl:",label" schema_label:"profile"`
	Targets []hclProfileTarget `hcl:"target,block"`
}

type hclProfileTarget struct {
	Name        string    `hcl:",label" schema_label:"target"`
	Formats     *[]string `hcl:"formats,optional"`
	Sign        *bool     `hcl:"sign,optional"`
	Notarize    *bool     `hcl:"notarize,optional"`
	Destination *string   `hcl:"destination,optional"`
}

type hclAssociation struct {
	Label       string    `hcl:",label" schema_label:"association"`
	Extensions  *[]string `hcl:"extensions,optional"`
	Name        *string   `hcl:"name,optional"`
	Description *string   `hcl:"description,optional"`
	Icon        *string   `hcl:"icon,optional"`
	Role        *string   `hcl:"role,optional"`
	MIMEType    *string   `hcl:"mime_type,optional"`
	Platforms   *[]string `hcl:"platforms,optional"`
}

type hclProtocol struct {
	Scheme      string    `hcl:",label" schema_label:"scheme"`
	Description *string   `hcl:"description,optional"`
	Platforms   *[]string `hcl:"platforms,optional"`
}

func decodeHCL(root, filename string, src []byte, selectedProfile string) (*Loaded, error) {
	file, diagnostics := hclsyntax.ParseConfig(src, filename, hcl.InitialPos)
	if diagnostics.HasErrors() {
		return nil, validationFromDiagnostics(diagnostics)
	}
	// hclsyntax.ParseConfig always returns a syntax body; the concrete assertion
	// keeps the literal validator independent from the generic hcl.Body API.
	body := file.Body.(*hclsyntax.Body)
	raw, err := decodeManifestSchema(body)
	if err != nil {
		return nil, err
	}
	origins := manifestOrigins(body)
	doc, err := documentFromHCL(raw)
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

// LoadFile validates an explicitly named manifest-like HCL file. It is used
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
	return decodeHCL(absRoot, filename, src, profile)
}

func validateLiteralExpression(expression hcl.Expression) error {
	switch expression := expression.(type) {
	case *hclsyntax.LiteralValueExpr:
		return nil
	case *hclsyntax.TemplateExpr:
		for _, part := range expression.Parts {
			if err := validateLiteralExpression(part); err != nil {
				return err
			}
		}
		return nil
	case *hclsyntax.TupleConsExpr:
		for _, item := range expression.Exprs {
			if err := validateLiteralExpression(item); err != nil {
				return err
			}
		}
		return nil
	case *hclsyntax.ObjectConsExpr:
		for _, item := range expression.Items {
			if err := validateLiteralExpression(item.KeyExpr); err != nil {
				return err
			}
			if err := validateLiteralExpression(item.ValueExpr); err != nil {
				return err
			}
		}
		return nil
	case *hclsyntax.ObjectConsKeyExpr:
		if !expression.ForceNonLiteral && hcl.ExprAsKeyword(expression.Wrapped) != "" {
			return nil
		}
		return validateLiteralExpression(expression.Wrapped)
	default:
		return &ValidationError{Detail: "only literal values are allowed; expressions, references, and calls are not supported", Range: sourceRange(expression.Range())}
	}
}

func documentFromHCL(raw hclDocument) (Document, error) {
	if raw.Project == nil {
		return Document{}, fieldValidationError("project", "block is required")
	}
	project := Project{}
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
		body *hclPlatform
	}{{"windows", raw.Windows}, {"darwin", raw.Darwin}, {"linux", raw.Linux}, {"ios", raw.IOS}, {"android", raw.Android}} {
		if item.body != nil {
			applyPlatform(&doc, item.name, item.body)
		}
	}
	seenTargets := map[string]bool{}
	for _, target := range raw.Targets {
		field := `target[` + strconv.Quote(target.Name) + `]`
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
		field := `package[` + strconv.Quote(pkg.Format) + `]`
		if seenPackages[pkg.Format] {
			return Document{}, fieldValidationError(field, "duplicate package block")
		}
		seenPackages[pkg.Format] = true
		if err := applyPackage(&doc.Package, pkg); err != nil {
			return Document{}, fieldValidationCause(field, err, "%v", err)
		}
	}
	for _, association := range raw.Associations {
		field := `file_association[` + strconv.Quote(association.Label) + `]`
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
		profileField := `profile[` + strconv.Quote(rawProfile.Name) + `]`
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
			targetField := profileField + `.target[` + strconv.Quote(rawTarget.Name) + `]`
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
		field := `protocol[` + strconv.Quote(protocol.Scheme) + `]`
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
		field := `hook[` + strconv.Quote(rawHook.Phase) + `]`
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

func applyFrontend(target *Frontend, raw *hclFrontend) {
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

func applyBindings(target *Bindings, raw *hclBindings) {
	setBool(&target.TypeScript, raw.TypeScript)
	setBool(&target.Interfaces, raw.Interfaces)
	setString(&target.OutputDirectory, raw.Output)
	setString(&target.ModelsFilename, raw.ModelsFilename)
	setString(&target.IndexFilename, raw.IndexFilename)
	setString(&target.TimeType, raw.TimeType)
}

func applyBuild(target *Build, raw *hclBuild) {
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

func applyDev(target *Dev, raw *hclDev) {
	setStrings(&target.Tags, raw.Tags)
	setInt(&target.DebounceMS, raw.DebounceMS)
	setInt(&target.GracePeriodMS, raw.GracePeriod)
	setString(&target.LogLevel, raw.LogLevel)
	setStrings(&target.Watch, raw.Watch)
	setStrings(&target.Exclude, raw.Exclude)
	setBool(&target.UseGitIgnore, raw.UseGitIgnore)
}

func applyPlatform(doc *Document, name string, raw *hclPlatform) {
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

func applyTarget(targets *Targets, raw hclTarget) error {
	platform, arch, err := parseTargetName(raw.Name)
	if err != nil {
		return err
	}
	target := targetByName(platformByName(targets, platform), arch)
	// parseTargetName accepts only registry targets, and every registry
	// architecture has a concrete slot in Platform.
	setStrings(&target.Tags, raw.Tags)
	setString(&target.MinimumVersion, raw.MinimumVersion)
	setInt(&target.BuildNumber, raw.BuildNumber)
	setString(&target.Toolchain, raw.Toolchain)
	setStringMap(&target.Environment, raw.Environment)
	setStrings(&target.LinkerFlags, raw.LDFlags)
	setStrings(&target.CompilerFlags, raw.CompilerFlags)
	setStrings(&target.GarbleArgs, raw.GarbleArgs)
	setBool(&target.Obfuscated, raw.Obfuscated)
	target.ObfuscatedSet = raw.Obfuscated != nil
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
	field := `profile[` + strconv.Quote(profile) + `].target[` + strconv.Quote(target.Target) + `]`
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

// EncodeEjectedHCL produces an inactive, complete reference manifest. It is
// intentionally a sibling file: ejection must never rewrite active project
// configuration or turn a routine CLI upgrade into a merge conflict.
func EncodeEjectedHCL(config Config, cliVersion string) ([]byte, error) {
	return encodeConfigHCL(config, fmt.Sprintf("# Generated by Wails CLI %s. This file is inactive; review and copy values into %s deliberately.\n\n", cliVersion, Filename))
}

func encodeConfigHCL(config Config, header string) ([]byte, error) {
	var output bytes.Buffer
	output.Grow(len(header) + 2048 + len(config.Associations)*256 + len(config.Protocols)*128 + len(config.Profiles)*256)
	output.WriteString(header)
	output.WriteString("version = 3\n\nproject {\n")
	hclString(&output, "name", config.Project.Name)
	hclString(&output, "product_name", config.Project.ProductName)
	hclString(&output, "identifier", config.Project.Identifier)
	hclString(&output, "version", config.Project.Version)
	hclString(&output, "company", config.Project.CompanyName)
	hclString(&output, "binary_name", config.Project.BinaryName)
	hclString(&output, "icon", config.Project.Icon)
	hclString(&output, "description", config.Project.Description)
	hclString(&output, "copyright", config.Project.Copyright)
	hclString(&output, "comments", config.Project.Comments)
	if config.Project.BuildNumber != 0 {
		hclIntIndented(&output, "build_number", config.Project.BuildNumber, "  ")
	}
	output.WriteString("}\n\nfrontend {\n")
	hclString(&output, "directory", config.Frontend.Directory)
	hclStrings(&output, "install", config.Frontend.Install)
	hclStrings(&output, "build", config.Frontend.Build)
	hclStrings(&output, "dev", config.Frontend.Dev)
	hclString(&output, "output", config.Frontend.OutputDirectory)
	hclStringMapIndented(&output, "environment", config.Frontend.Environment, "  ")
	output.WriteString("  bindings {\n")
	hclBoolIndented(&output, "typescript", config.Frontend.Bindings.TypeScript, "    ")
	hclBoolIndented(&output, "interfaces", config.Frontend.Bindings.Interfaces, "    ")
	hclStringIndented(&output, "output", config.Frontend.Bindings.OutputDirectory, "    ")
	hclStringIndented(&output, "models_filename", config.Frontend.Bindings.ModelsFilename, "    ")
	hclStringIndented(&output, "index_filename", config.Frontend.Bindings.IndexFilename, "    ")
	hclStringIndented(&output, "time_type", config.Frontend.Bindings.TimeType, "    ")
	output.WriteString("  }\n}\n\nbuild {\n")
	hclString(&output, "output", config.Build.OutputDirectory)
	hclStrings(&output, "tags", config.Build.Go.Tags)
	hclBoolIndented(&output, "trim_path", config.Build.TrimPath, "  ")
	hclBoolIndented(&output, "strip", config.Build.Strip, "  ")
	hclBoolIndented(&output, "vcs_info", config.Build.VCSInfo, "  ")
	hclBoolIndented(&output, "obfuscated", config.Build.Obfuscation, "  ")
	hclStringMapIndented(&output, "environment", config.Build.Environment, "  ")
	hclStrings(&output, "garble_args", config.Build.Go.GarbleArgs)
	hclStrings(&output, "ldflags", config.Build.Go.LinkerFlags)
	hclStrings(&output, "compiler_flags", config.Build.Go.CompilerFlags)
	output.WriteString("}\n\ndev {\n")
	hclStrings(&output, "tags", config.Dev.Tags)
	hclIntIndented(&output, "debounce_ms", config.Dev.DebounceMS, "  ")
	hclString(&output, "log_level", config.Dev.LogLevel)
	hclStringsIndentedPresent(&output, "watch", config.Dev.Watch, "  ", manifestValueWasExplicit(config.Origins, "dev.watch"))
	hclStringsIndentedPresent(&output, "exclude", config.Dev.Exclude, "  ", manifestValueWasExplicit(config.Origins, "dev.exclude"))
	hclBoolIndented(&output, "use_git_ignore", config.Dev.UseGitIgnore, "  ")
	hclIntIndented(&output, "grace_period_ms", config.Dev.GracePeriodMS, "  ")
	output.WriteString("}\n\n")
	for _, platform := range []struct {
		name    string
		value   Platform
		signing SigningPlatform
	}{{"windows", config.Targets.Windows, config.Signing.Windows}, {"darwin", config.Targets.Darwin, config.Signing.Darwin}, {"linux", config.Targets.Linux, config.Signing.Linux}, {"ios", config.Targets.IOS, config.Signing.IOS}, {"android", config.Targets.Android, config.Signing.Android}} {
		writeEjectedPlatform(&output, platform.name, platform.value, platform.signing)
	}
	for _, target := range []struct {
		name  string
		value Target
	}{
		{"windows/amd64", config.Targets.Windows.AMD64}, {"windows/arm64", config.Targets.Windows.ARM64}, {"darwin/amd64", config.Targets.Darwin.AMD64}, {"darwin/arm64", config.Targets.Darwin.ARM64}, {"darwin/universal", config.Targets.Darwin.Universal}, {"linux/amd64", config.Targets.Linux.AMD64}, {"linux/arm64", config.Targets.Linux.ARM64}, {"ios/arm64", config.Targets.IOS.ARM64}, {"android/amd64", config.Targets.Android.AMD64}, {"android/arm64", config.Targets.Android.ARM64}, {"android/universal", config.Targets.Android.Universal},
	} {
		if !target.value.Enabled && len(target.value.Tags) == 0 && target.value.MinimumVersion == "" && target.value.BuildNumber == 0 && target.value.Toolchain == "" && len(target.value.Environment) == 0 && len(target.value.LinkerFlags) == 0 && len(target.value.CompilerFlags) == 0 && len(target.value.GarbleArgs) == 0 && !target.value.ObfuscatedSet {
			continue
		}
		hclLabeledBlockStart(&output, "", "target", target.name)
		hclStringsIndented(&output, "tags", target.value.Tags, "  ")
		hclStringIndented(&output, "minimum_version", target.value.MinimumVersion, "  ")
		hclStringIndented(&output, "toolchain", target.value.Toolchain, "  ")
		hclStringMapIndented(&output, "environment", target.value.Environment, "  ")
		hclStringsIndented(&output, "ldflags", target.value.LinkerFlags, "  ")
		hclStringsIndented(&output, "compiler_flags", target.value.CompilerFlags, "  ")
		hclStringsIndented(&output, "garble_args", target.value.GarbleArgs, "  ")
		if target.value.ObfuscatedSet {
			hclBoolIndented(&output, "obfuscated", target.value.Obfuscated, "  ")
		}
		if target.value.BuildNumber != 0 {
			hclIntIndented(&output, "build_number", target.value.BuildNumber, "  ")
		}
		output.WriteString("}\n\n")
	}
	writeEjectedPackages(&output, config.Package)
	for _, phase := range HookPhases {
		hook, ok := config.Hooks[phase]
		if !ok {
			continue
		}
		hclLabeledBlockStart(&output, "", "hook", string(phase))
		hclStringIndented(&output, "script", hook.Script, "  ")
		hclStringIndented(&output, "directory", hook.Directory, "  ")
		hclBoolIndented(&output, "cache", hook.Cache, "  ")
		hclStringsIndented(&output, "inputs", hook.Inputs, "  ")
		hclStringsIndented(&output, "outputs", hook.Outputs, "  ")
		output.WriteString("}\n\n")
	}
	for _, association := range config.Associations {
		hclLabeledBlockStart(&output, "", "file_association", association.Name)
		hclStringsIndented(&output, "extensions", association.Extensions, "  ")
		hclStringIndented(&output, "name", association.Name, "  ")
		hclStringIndented(&output, "description", association.Description, "  ")
		hclStringIndented(&output, "icon", association.Icon, "  ")
		hclStringIndented(&output, "role", association.Role, "  ")
		hclStringIndented(&output, "mime_type", association.MIMEType, "  ")
		hclStringsIndented(&output, "platforms", association.Platforms, "  ")
		output.WriteString("}\n\n")
	}
	for _, protocol := range config.Protocols {
		hclLabeledBlockStart(&output, "", "protocol", protocol.Scheme)
		hclStringIndented(&output, "description", protocol.Description, "  ")
		hclStringsIndented(&output, "platforms", protocol.Platforms, "  ")
		output.WriteString("}\n\n")
	}
	for _, profile := range sortedProfiles(config.Profiles) {
		hclLabeledBlockStart(&output, "", "profile", profile.Name)
		for _, target := range profile.Targets {
			hclLabeledBlockStart(&output, "  ", "target", target.Target)
			hclStringsIndented(&output, "formats", target.Formats, "    ")
			hclBoolIndented(&output, "sign", target.Sign, "    ")
			hclBoolIndented(&output, "notarize", target.Notarize, "    ")
			hclStringIndented(&output, "destination", target.Destination, "    ")
			output.WriteString("  }\n")
		}
		output.WriteString("}\n\n")
	}
	return output.Bytes(), nil
}

func writeEjectedPlatform(output *bytes.Buffer, name string, platform Platform, signing SigningPlatform) {
	if !platformConfigured(platform) && signing == (SigningPlatform{}) {
		return
	}
	output.WriteString(name)
	output.WriteString(" {\n")
	switch name {
	case "ios":
		hclStringIndented(output, "display_name", platform.ProductName, "  ")
		hclStringIndented(output, "bundle_id", platform.Identifier, "  ")
	case "android":
		hclStringIndented(output, "display_name", platform.ProductName, "  ")
		hclStringIndented(output, "application_id", platform.Identifier, "  ")
	default:
		hclStringIndented(output, "product_name", platform.ProductName, "  ")
		hclStringIndented(output, "identifier", platform.Identifier, "  ")
	}
	if platformAttributeAllowed(name, "minimum_version") {
		hclStringIndented(output, "minimum_version", platform.MinimumVersion, "  ")
	}
	if platformAttributeAllowed(name, "build_number") && platform.BuildNumber != 0 {
		hclIntIndented(output, "build_number", platform.BuildNumber, "  ")
	}
	writePlatformString(output, name, "icon", platform.Icon)
	writePlatformString(output, name, "manifest", platform.Manifest)
	writePlatformString(output, name, "assets_car", platform.AssetsCar)
	writePlatformString(output, name, "info_plist", platform.InfoPlist)
	writePlatformString(output, name, "publisher", platform.Publisher)
	writePlatformString(output, name, "desktop_entry", platform.DesktopEntry)
	writePlatformString(output, name, "company", platform.CompanyName)
	writePlatformString(output, name, "comments", platform.Comments)
	writePlatformString(output, name, "cf_bundle_icon_name", platform.CFBundleIconName)
	if platformAttributeAllowed(name, "background_modes") {
		hclStringsIndented(output, "background_modes", platform.BackgroundModes, "  ")
	}
	writePlatformString(output, name, "version_name", platform.VersionName)
	if platformAttributeAllowed(name, "version_code") && platform.VersionCode != 0 {
		hclIntIndented(output, "version_code", platform.VersionCode, "  ")
	}
	if platformAttributeAllowed(name, "minimum_sdk") && platform.MinimumSDK != 0 {
		hclIntIndented(output, "minimum_sdk", platform.MinimumSDK, "  ")
	}
	if platformAttributeAllowed(name, "target_sdk") && platform.TargetSDK != 0 {
		hclIntIndented(output, "target_sdk", platform.TargetSDK, "  ")
	}
	if platformAttributeAllowed(name, "capabilities") {
		hclStringsIndented(output, "capabilities", platform.Capabilities, "  ")
	}
	if signing != (SigningPlatform{}) {
		output.WriteString("  signing {\n")
		hclStringIndented(output, "credential", signing.Credential, "    ")
		hclStringIndented(output, "identity", signing.Identity, "    ")
		hclStringIndented(output, "certificate", signing.Certificate, "    ")
		hclStringIndented(output, "thumbprint", signing.Thumbprint, "    ")
		hclStringIndented(output, "timestamp_server", signing.TimestampServer, "    ")
		hclStringIndented(output, "entitlements", signing.Entitlements, "    ")
		hclStringIndented(output, "provisioning_profile", signing.ProvisioningProfile, "    ")
		hclStringIndented(output, "key_alias", signing.KeyAlias, "    ")
		output.WriteString("  }\n")
	}
	if signing.Notarize && platformBlockAllowed(name, "notarization") {
		output.WriteString("  notarization {\n")
		hclStringIndented(output, "credential", signing.NotarizationCredential, "    ")
		output.WriteString("  }\n")
	}
	output.WriteString("}\n\n")
}

func writePlatformString(output *bytes.Buffer, platform, name, value string) {
	if platformAttributeAllowed(platform, name) {
		hclStringIndented(output, name, value, "  ")
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

func writeEjectedPackages(output *bytes.Buffer, packages Packages) {
	formats := []struct{ platform, format string }{
		{"windows", "nsis"}, {"windows", "msix"}, {"darwin", "dmg"},
		{"linux", "appimage"}, {"linux", "deb"}, {"linux", "rpm"}, {"linux", "archlinux"},
	}
	for _, item := range formats {
		value, _ := ResolvePackageFormat(packages, item.platform, item.format)
		if !packageConfigured(value) {
			continue
		}
		hclLabeledBlockStart(output, "", "package", item.format)
		writePackageString(output, item.format, "template", value.Template)
		writePackageString(output, item.format, "install_scope", value.InstallScope)
		writePackageString(output, item.format, "publisher", value.Publisher)
		writePackageString(output, item.format, "manifest", value.Manifest)
		writePackageString(output, item.format, "background", value.Background)
		writePackageString(output, item.format, "volume_icon", value.VolumeIcon)
		writePackageString(output, item.format, "file_icon", value.FileIcon)
		if packageAttributeAllowed(item.format, "files") {
			hclStringMapIndented(output, "files", value.Files, "  ")
		}
		if packageAttributeAllowed(item.format, "window_width") && value.WindowWidth != 0 {
			hclIntIndented(output, "window_width", value.WindowWidth, "  ")
		}
		if packageAttributeAllowed(item.format, "window_height") && value.WindowHeight != 0 {
			hclIntIndented(output, "window_height", value.WindowHeight, "  ")
		}
		writePackageString(output, item.format, "icon", value.Icon)
		writePackageString(output, item.format, "desktop_entry", value.DesktopEntry)
		if packageAttributeAllowed(item.format, "categories") {
			hclStringsIndented(output, "categories", value.Categories, "  ")
		}
		writePackageString(output, item.format, "maintainer", value.Maintainer)
		writePackageString(output, item.format, "section", value.Section)
		if packageAttributeAllowed(item.format, "dependencies") {
			hclStringsIndented(output, "dependencies", value.Dependencies, "  ")
		}
		writePackageString(output, item.format, "pre_install", value.PreInstall)
		writePackageString(output, item.format, "post_install", value.PostInstall)
		writePackageString(output, item.format, "pre_remove", value.PreRemove)
		writePackageString(output, item.format, "post_remove", value.PostRemove)
		output.WriteString("}\n\n")
	}
}

func writePackageString(output *bytes.Buffer, format, name, value string) {
	if packageAttributeAllowed(format, name) {
		hclStringIndented(output, name, value, "  ")
	}
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

func hclString(output *bytes.Buffer, name, value string) {
	hclStringIndented(output, name, value, "  ")
}

func hclLabeledBlockStart(output *bytes.Buffer, indent, blockType, label string) {
	output.WriteString(indent)
	output.WriteString(blockType)
	output.WriteByte(' ')
	writeHCLStringLiteral(output, label)
	output.WriteString(" {\n")
}

func hclBoolIndented(output *bytes.Buffer, name string, value bool, indent string) {
	output.WriteString(indent)
	output.WriteString(name)
	if value {
		output.WriteString(" = true\n")
	} else {
		output.WriteString(" = false\n")
	}
}

func hclIntIndented(output *bytes.Buffer, name string, value int, indent string) {
	output.WriteString(indent)
	output.WriteString(name)
	output.WriteString(" = ")
	var buffer [24]byte
	output.Write(strconv.AppendInt(buffer[:0], int64(value), 10))
	output.WriteByte('\n')
}

func hclStringIndented(output *bytes.Buffer, name, value, indent string) {
	if value != "" {
		output.WriteString(indent)
		output.WriteString(name)
		output.WriteString(" = ")
		writeHCLStringLiteral(output, value)
		output.WriteByte('\n')
	}
}
func hclStrings(output *bytes.Buffer, name string, values []string) {
	hclStringsIndented(output, name, values, "  ")
}
func hclStringsIndented(output *bytes.Buffer, name string, values []string, indent string) {
	hclStringsIndentedPresent(output, name, values, indent, false)
}

func hclStringsIndentedPresent(output *bytes.Buffer, name string, values []string, indent string, present bool) {
	if len(values) == 0 && !present {
		return
	}
	output.WriteString(indent)
	output.WriteString(name)
	output.WriteString(" = [")
	for index, value := range values {
		if index != 0 {
			output.WriteString(", ")
		}
		writeHCLStringLiteral(output, value)
	}
	output.WriteString("]\n")
}

func manifestValueWasExplicit(origins map[string]Origin, path string) bool {
	return origins[path].Kind == OriginManifest
}
func hclStringMapIndented(output *bytes.Buffer, name string, values map[string]string, indent string) {
	if len(values) == 0 {
		return
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	output.WriteString(indent)
	output.WriteString(name)
	output.WriteString(" = {\n")
	for _, key := range keys {
		output.WriteString(indent)
		output.WriteString("  ")
		writeHCLStringLiteral(output, key)
		output.WriteString(" = ")
		writeHCLStringLiteral(output, values[key])
		output.WriteByte('\n')
	}
	output.WriteString(indent)
	output.WriteString("}\n")
}

func writeHCLStringLiteral(output *bytes.Buffer, value string) {
	for index := 0; index < len(value); index++ {
		if value[index] >= 0x80 {
			value = ctystrings.Normalize(value)
			break
		}
	}
	output.WriteByte('"')
	const hexadecimal = "0123456789abcdef"
	for index := 0; index < len(value); index++ {
		character := value[index]
		switch character {
		case '\n':
			output.WriteString(`\n`)
		case '\r':
			output.WriteString(`\r`)
		case '\t':
			output.WriteString(`\t`)
		case '"':
			output.WriteString(`\"`)
		case '\\':
			output.WriteString(`\\`)
		case '$', '%':
			output.WriteByte(character)
			if index+1 < len(value) && value[index+1] == '{' {
				output.WriteByte(character)
			}
		default:
			if character < 0x20 {
				output.WriteString(`\u00`)
				output.WriteByte(hexadecimal[character>>4])
				output.WriteByte(hexadecimal[character&0x0f])
			} else {
				output.WriteByte(character)
			}
		}
	}
	output.WriteByte('"')
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
