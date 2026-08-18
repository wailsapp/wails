package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// hclDocument is intentionally a closed schema. Keep the HCL surface small:
// this is declarative build intent, never a programmable pipeline.
type hclDocument struct {
	Version      int              `hcl:"version"`
	Project      *hclProject      `hcl:"project,block"`
	Frontend     *hclFrontend     `hcl:"frontend,block"`
	Build        *hclBuild        `hcl:"build,block"`
	Dev          *hclDev          `hcl:"dev,block"`
	Windows      *hclPlatform     `hcl:"windows,block"`
	Darwin       *hclPlatform     `hcl:"darwin,block"`
	Linux        *hclPlatform     `hcl:"linux,block"`
	IOS          *hclPlatform     `hcl:"ios,block"`
	Android      *hclPlatform     `hcl:"android,block"`
	Targets      []hclTarget      `hcl:"target,block"`
	Packages     []hclPackage     `hcl:"package,block"`
	Profiles     []hclProfile     `hcl:"profile,block"`
	Associations []hclAssociation `hcl:"file_association,block"`
	Protocols    []hclProtocol    `hcl:"protocol,block"`
}

type hclProject struct {
	Name        *string `hcl:"name,optional"`
	ProductName *string `hcl:"product_name,optional"`
	Identifier  *string `hcl:"identifier,optional"`
	Version     *string `hcl:"version,optional"`
	Company     *string `hcl:"company,optional"`
	BinaryName  *string `hcl:"binary_name,optional"`
	Icon        *string `hcl:"icon,optional"`
	Description *string `hcl:"description,optional"`
	Copyright   *string `hcl:"copyright,optional"`
	Comments    *string `hcl:"comments,optional"`
	BuildNumber *int    `hcl:"build_number,optional"`
}

type hclFrontend struct {
	Directory *string   `hcl:"directory,optional"`
	Install   *[]string `hcl:"install,optional"`
	Build     *[]string `hcl:"build,optional"`
	Dev       *[]string `hcl:"dev,optional"`
	Output    *string   `hcl:"output,optional"`
}

type hclBuild struct {
	Output        *string   `hcl:"output,optional"`
	Tags          *[]string `hcl:"tags,optional"`
	TrimPath      *bool     `hcl:"trim_path,optional"`
	Strip         *bool     `hcl:"strip,optional"`
	Obfuscated    *bool     `hcl:"obfuscated,optional"`
	GarbleArgs    *[]string `hcl:"garble_args,optional"`
	LDFlags       *[]string `hcl:"ldflags,optional"`
	CompilerFlags *[]string `hcl:"compiler_flags,optional"`
}

type hclDev struct {
	DebounceMS   *int      `hcl:"debounce_ms,optional"`
	LogLevel     *string   `hcl:"log_level,optional"`
	Watch        *[]string `hcl:"watch,optional"`
	Exclude      *[]string `hcl:"exclude,optional"`
	UseGitIgnore *bool     `hcl:"use_git_ignore,optional"`
	GracePeriod  *int      `hcl:"grace_period_ms,optional"`
}

type hclPlatform struct {
	ProductName    *string          `hcl:"product_name,optional"`
	Identifier     *string          `hcl:"identifier,optional"`
	MinimumVersion *string          `hcl:"minimum_version,optional"`
	BuildNumber    *int             `hcl:"build_number,optional"`
	Capabilities   *[]string        `hcl:"capabilities,optional"`
	Icon           *string          `hcl:"icon,optional"`
	Manifest       *string          `hcl:"manifest,optional"`
	AssetsCar      *string          `hcl:"assets_car,optional"`
	InfoPlist      *string          `hcl:"info_plist,optional"`
	Publisher      *string          `hcl:"publisher,optional"`
	BundleID       *string          `hcl:"bundle_id,optional"`
	DisplayName    *string          `hcl:"display_name,optional"`
	VersionName    *string          `hcl:"version_name,optional"`
	VersionCode    *int             `hcl:"version_code,optional"`
	MinimumSDK     *int             `hcl:"minimum_sdk,optional"`
	TargetSDK      *int             `hcl:"target_sdk,optional"`
	Signing        *hclSigning      `hcl:"signing,block"`
	Notarization   *hclNotarization `hcl:"notarization,block"`
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
	Name           string    `hcl:",label"`
	Tags           *[]string `hcl:"tags,optional"`
	MinimumVersion *string   `hcl:"minimum_version,optional"`
	Variant        *string   `hcl:"variant,optional"`
	BuildNumber    *int      `hcl:"build_number,optional"`
}

type hclPackage struct {
	Format       string  `hcl:",label"`
	Template     *string `hcl:"template,optional"`
	InstallScope *string `hcl:"install_scope,optional"`
	Background   *string `hcl:"background,optional"`
	VolumeIcon   *string `hcl:"volume_icon,optional"`
	FileIcon     *string `hcl:"file_icon,optional"`
	Files        *string `hcl:"files,optional"`
	Categories   *string `hcl:"categories,optional"`
	WindowWidth  *int    `hcl:"window_width,optional"`
	WindowHeight *int    `hcl:"window_height,optional"`
	PreInstall   *string `hcl:"pre_install,optional"`
	PostInstall  *string `hcl:"post_install,optional"`
	PreRemove    *string `hcl:"pre_remove,optional"`
	PostRemove   *string `hcl:"post_remove,optional"`
}

type hclProfile struct {
	Name    string             `hcl:",label"`
	Targets []hclProfileTarget `hcl:"target,block"`
}

type hclProfileTarget struct {
	Name        string    `hcl:",label"`
	Formats     *[]string `hcl:"formats,optional"`
	Sign        *bool     `hcl:"sign,optional"`
	Notarize    *bool     `hcl:"notarize,optional"`
	Destination *string   `hcl:"destination,optional"`
}

type hclAssociation struct {
	Label       string    `hcl:",label"`
	Extensions  *[]string `hcl:"extensions,optional"`
	Name        *string   `hcl:"name,optional"`
	Description *string   `hcl:"description,optional"`
	Icon        *string   `hcl:"icon,optional"`
	Role        *string   `hcl:"role,optional"`
	MIMEType    *string   `hcl:"mime_type,optional"`
	Platforms   *[]string `hcl:"platforms,optional"`
}

type hclProtocol struct {
	Scheme      string    `hcl:",label"`
	Description *string   `hcl:"description,optional"`
	Platforms   *[]string `hcl:"platforms,optional"`
}

func decodeHCL(root, filename string, src []byte, selectedProfile string) (*Loaded, error) {
	file, diagnostics := hclsyntax.ParseConfig(src, filename, hcl.InitialPos)
	if diagnostics.HasErrors() {
		return nil, fmt.Errorf("parse %s: %s", Filename, diagnostics.Error())
	}
	// hclsyntax.ParseConfig always returns a syntax body; the concrete assertion
	// keeps the literal validator independent from the generic hcl.Body API.
	body := file.Body.(*hclsyntax.Body)
	if err := validateLiteralOnlyBody(body, true); err != nil {
		return nil, err
	}
	var raw hclDocument
	diagnostics = gohcl.DecodeBody(file.Body, nil, &raw)
	if diagnostics.HasErrors() {
		return nil, fmt.Errorf("parse %s: %s", Filename, diagnostics.Error())
	}
	if raw.Version != 3 {
		return nil, fmt.Errorf("%s: version must be 3", Filename)
	}
	doc, err := documentFromHCL(raw)
	if err != nil {
		return nil, err
	}
	config := configFromDocument(root, selectedProfile, doc)
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
		return nil, err
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

func validateLiteralOnlyBody(body *hclsyntax.Body, topLevel bool) error {
	if topLevel {
		version, exists := body.Attributes["version"]
		if !exists {
			return fmt.Errorf("%s: first attribute must be version = 3", Filename)
		}
		first := version.Range().Start.Byte
		for _, attribute := range body.Attributes {
			if attribute.Range().Start.Byte < first {
				return fmt.Errorf("%s:%d: version must be the first attribute", Filename, attribute.Range().Start.Line)
			}
		}
		for _, block := range body.Blocks {
			if block.TypeRange.Start.Byte < first {
				return fmt.Errorf("%s:%d: version must be the first attribute", Filename, block.TypeRange.Start.Line)
			}
		}
	}
	for _, attribute := range body.Attributes {
		if err := validateLiteralExpression(attribute.Expr); err != nil {
			return fmt.Errorf("%s:%d: %s", Filename, attribute.Range().Start.Line, err)
		}
	}
	for _, block := range body.Blocks {
		if err := validateLiteralOnlyBody(block.Body, false); err != nil {
			return err
		}
	}
	return nil
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
		return validateLiteralExpression(expression.Wrapped)
	default:
		return fmt.Errorf("only literal values are allowed; expressions, references, and calls are not supported")
	}
}

func documentFromHCL(raw hclDocument) (Document, error) {
	if raw.Project == nil {
		return Document{}, fmt.Errorf("%s: project block is required", Filename)
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
		if seenTargets[target.Name] {
			return Document{}, fmt.Errorf("duplicate target %q", target.Name)
		}
		seenTargets[target.Name] = true
		if err := applyTarget(&doc.Targets, target); err != nil {
			return Document{}, err
		}
	}
	seenPackages := map[string]bool{}
	for _, pkg := range raw.Packages {
		if seenPackages[pkg.Format] {
			return Document{}, fmt.Errorf("duplicate package %q", pkg.Format)
		}
		seenPackages[pkg.Format] = true
		if err := applyPackage(&doc.Package, pkg); err != nil {
			return Document{}, err
		}
	}
	for _, association := range raw.Associations {
		if association.Extensions == nil || len(*association.Extensions) == 0 {
			return Document{}, fmt.Errorf("file_association %q requires extensions", association.Label)
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
		if rawProfile.Name == "" || rawProfile.Name == "default" || !slugPattern.MatchString(rawProfile.Name) {
			return Document{}, fmt.Errorf("profile name must be a lowercase slug and cannot be default")
		}
		if seenProfiles[rawProfile.Name] {
			return Document{}, fmt.Errorf("duplicate profile %q", rawProfile.Name)
		}
		seenProfiles[rawProfile.Name] = true
		if len(rawProfile.Targets) == 0 {
			return Document{}, fmt.Errorf("profile %q requires at least one target", rawProfile.Name)
		}
		profile := Profile{Name: rawProfile.Name}
		seen := map[string]bool{}
		for _, rawTarget := range rawProfile.Targets {
			if seen[rawTarget.Name] {
				return Document{}, fmt.Errorf("profile %q contains duplicate target %q", rawProfile.Name, rawTarget.Name)
			}
			seen[rawTarget.Name] = true
			if _, _, err := parseTargetName(rawTarget.Name); err != nil {
				return Document{}, err
			}
			entry := ProfileTarget{Target: rawTarget.Name}
			setStrings(&entry.Formats, rawTarget.Formats)
			setBool(&entry.Sign, rawTarget.Sign)
			setBool(&entry.Notarize, rawTarget.Notarize)
			setString(&entry.Destination, rawTarget.Destination)
			profile.Targets = append(profile.Targets, entry)
		}
		doc.Profiles[profile.Name] = profile
	}
	seenProtocols := map[string]bool{}
	for _, protocol := range raw.Protocols {
		if protocol.Scheme == "" {
			return Document{}, fmt.Errorf("protocol label cannot be empty")
		}
		if seenProtocols[protocol.Scheme] {
			return Document{}, fmt.Errorf("duplicate protocol %q", protocol.Scheme)
		}
		seenProtocols[protocol.Scheme] = true
		entry := Protocol{Scheme: protocol.Scheme}
		setString(&entry.Description, protocol.Description)
		setStrings(&entry.Platforms, protocol.Platforms)
		doc.Protocols = append(doc.Protocols, entry)
	}
	return doc, nil
}

func applyFrontend(target *Frontend, raw *hclFrontend) {
	setString(&target.Directory, raw.Directory)
	setString(&target.OutputDirectory, raw.Output)
	setStrings(&target.Install, raw.Install)
	setStrings(&target.Build, raw.Build)
	setStrings(&target.Dev, raw.Dev)
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

func applyBuild(target *Build, raw *hclBuild) {
	setString(&target.OutputDirectory, raw.Output)
	setBool(&target.TrimPath, raw.TrimPath)
	setBool(&target.Strip, raw.Strip)
	setBool(&target.Obfuscation, raw.Obfuscated)
	setStrings(&target.Go.Tags, raw.Tags)
	setStrings(&target.Go.GarbleArgs, raw.GarbleArgs)
	setStrings(&target.Go.LinkerFlags, raw.LDFlags)
	setStrings(&target.Go.CompilerFlags, raw.CompilerFlags)
}

func applyDev(target *Dev, raw *hclDev) {
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
	setString(&platform.Identifier, raw.BundleID)
	setString(&platform.ProductName, raw.DisplayName)
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
		setString(&signing.Credential, raw.Notarization.Credential)
	}
}

func applyTarget(targets *Targets, raw hclTarget) error {
	platform, arch, err := parseTargetName(raw.Name)
	if err != nil {
		return err
	}
	target := targetByName(platformByName(targets, platform), arch)
	if target == nil {
		return fmt.Errorf("unsupported target %q", raw.Name)
	}
	setStrings(&target.Tags, raw.Tags)
	setString(&target.MinimumVersion, raw.MinimumVersion)
	setString(&target.Variant, raw.Variant)
	setInt(&target.BuildNumber, raw.BuildNumber)
	return nil
}

func applyPackage(packages *Packages, raw hclPackage) error {
	platforms := packagePlatformsForFormat(packages, raw.Format)
	if len(platforms) == 0 {
		return fmt.Errorf("unsupported package format %q", raw.Format)
	}
	for _, platform := range platforms {
		// packagePlatformsForFormat has already rejected unknown formats and
		// returns platforms whose field matches raw.Format.
		format, _ := packageFormatPointer(platform, raw.Format)
		setString(&format.Template, raw.Template)
		if format.Options == nil {
			format.Options = map[string]any{}
		}
		setOption(format.Options, "install_scope", raw.InstallScope)
		setOption(format.Options, "background", raw.Background)
		setOption(format.Options, "volume_icon", raw.VolumeIcon)
		setOption(format.Options, "file_icon", raw.FileIcon)
		setOption(format.Options, "files", raw.Files)
		setOption(format.Options, "categories", raw.Categories)
		setOption(format.Options, "window_width", raw.WindowWidth)
		setOption(format.Options, "window_height", raw.WindowHeight)
		setOption(format.Options, "pre_install", raw.PreInstall)
		setOption(format.Options, "post_install", raw.PostInstall)
		setOption(format.Options, "pre_remove", raw.PreRemove)
		setOption(format.Options, "post_remove", raw.PostRemove)
	}
	return nil
}

func parseTargetName(value string) (string, string, error) {
	parts := strings.Split(value, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("target must be platform/architecture, got %q", value)
	}
	if parts[1] == "universal" && parts[0] != "darwin" && parts[0] != "android" {
		return "", "", fmt.Errorf("universal target is only valid for darwin or android")
	}
	switch parts[0] {
	case "windows", "darwin", "linux", "ios", "android":
	default:
		return "", "", fmt.Errorf("unsupported target platform %q", parts[0])
	}
	return parts[0], parts[1], nil
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

func packagePlatformsForFormat(packages *Packages, format string) []*PackagePlatform {
	switch format {
	case "nsis", "msix":
		return []*PackagePlatform{&packages.Windows}
	case "dmg":
		return []*PackagePlatform{&packages.Darwin}
	case "appimage", "deb", "rpm", "archlinux":
		return []*PackagePlatform{&packages.Linux}
	case "ipa":
		return []*PackagePlatform{&packages.IOS}
	case "apk", "aab":
		return []*PackagePlatform{&packages.Android}
	case "app":
		return []*PackagePlatform{&packages.Darwin, &packages.IOS}
	default:
		return nil
	}
}

func packageFormatPointer(platform *PackagePlatform, format string) (*PackageFormat, error) {
	switch format {
	case "nsis":
		return &platform.NSIS, nil
	case "msix":
		return &platform.MSIX, nil
	case "app":
		return &platform.App, nil
	case "dmg":
		return &platform.DMG, nil
	case "appimage":
		return &platform.AppImage, nil
	case "deb":
		return &platform.Deb, nil
	case "rpm":
		return &platform.RPM, nil
	case "archlinux":
		return &platform.ArchLinux, nil
	case "ipa":
		return &platform.IPA, nil
	case "apk":
		return &platform.APK, nil
	case "aab":
		return &platform.AAB, nil
	default:
		return nil, fmt.Errorf("unsupported package format %q", format)
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
func boolPointer(value bool) *bool { return &value }
func setOption(options map[string]any, key string, value any) {
	switch value := value.(type) {
	case *string:
		if value != nil {
			options[key] = *value
		}
	case *int:
		if value != nil {
			options[key] = *value
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
	var output strings.Builder
	output.WriteString(header)
	fmt.Fprintln(&output, "version = 3")
	fmt.Fprintln(&output)
	fmt.Fprintln(&output, "project {")
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
		fmt.Fprintf(&output, "  build_number = %d\n", config.Project.BuildNumber)
	}
	fmt.Fprintln(&output, "}")
	fmt.Fprintln(&output)
	fmt.Fprintln(&output, "frontend {")
	hclString(&output, "directory", config.Frontend.Directory)
	hclStrings(&output, "install", config.Frontend.Install)
	hclStrings(&output, "build", config.Frontend.Build)
	hclStrings(&output, "dev", config.Frontend.Dev)
	hclString(&output, "output", config.Frontend.OutputDirectory)
	fmt.Fprintln(&output, "}")
	fmt.Fprintln(&output)
	fmt.Fprintln(&output, "build {")
	hclString(&output, "output", config.Build.OutputDirectory)
	hclStrings(&output, "tags", config.Build.Go.Tags)
	fmt.Fprintf(&output, "  trim_path = %t\n  strip = %t\n  obfuscated = %t\n", config.Build.TrimPath, config.Build.Strip, config.Build.Obfuscation)
	hclStrings(&output, "garble_args", config.Build.Go.GarbleArgs)
	hclStrings(&output, "ldflags", config.Build.Go.LinkerFlags)
	hclStrings(&output, "compiler_flags", config.Build.Go.CompilerFlags)
	fmt.Fprintln(&output, "}")
	fmt.Fprintln(&output)
	fmt.Fprintln(&output, "dev {")
	fmt.Fprintf(&output, "  debounce_ms = %d\n", config.Dev.DebounceMS)
	hclString(&output, "log_level", config.Dev.LogLevel)
	hclStrings(&output, "watch", config.Dev.Watch)
	hclStrings(&output, "exclude", config.Dev.Exclude)
	fmt.Fprintf(&output, "  use_git_ignore = %t\n  grace_period_ms = %d\n", config.Dev.UseGitIgnore, config.Dev.GracePeriodMS)
	fmt.Fprintln(&output, "}")
	fmt.Fprintln(&output)
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
		if !target.value.Enabled && len(target.value.Tags) == 0 && target.value.Variant == "" && target.value.MinimumVersion == "" && target.value.BuildNumber == 0 {
			continue
		}
		fmt.Fprintf(&output, "target %s {\n", strconv.Quote(target.name))
		hclStringsIndented(&output, "tags", target.value.Tags, "  ")
		hclStringIndented(&output, "minimum_version", target.value.MinimumVersion, "  ")
		hclStringIndented(&output, "variant", target.value.Variant, "  ")
		if target.value.BuildNumber != 0 {
			fmt.Fprintf(&output, "  build_number = %d\n", target.value.BuildNumber)
		}
		fmt.Fprintln(&output, "}")
		fmt.Fprintln(&output)
	}
	writeEjectedPackages(&output, config.Package)
	for _, association := range config.Associations {
		fmt.Fprintf(&output, "file_association %s {\n", strconv.Quote(association.Name))
		hclStringsIndented(&output, "extensions", association.Extensions, "  ")
		hclStringIndented(&output, "name", association.Name, "  ")
		hclStringIndented(&output, "description", association.Description, "  ")
		hclStringIndented(&output, "icon", association.Icon, "  ")
		hclStringIndented(&output, "role", association.Role, "  ")
		hclStringIndented(&output, "mime_type", association.MIMEType, "  ")
		hclStringsIndented(&output, "platforms", association.Platforms, "  ")
		fmt.Fprintln(&output, "}")
		fmt.Fprintln(&output)
	}
	for _, protocol := range config.Protocols {
		fmt.Fprintf(&output, "protocol %s {\n", strconv.Quote(protocol.Scheme))
		hclStringIndented(&output, "description", protocol.Description, "  ")
		hclStringsIndented(&output, "platforms", protocol.Platforms, "  ")
		fmt.Fprintln(&output, "}")
		fmt.Fprintln(&output)
	}
	for _, profile := range sortedProfiles(config.Profiles) {
		fmt.Fprintf(&output, "profile %s {\n", strconv.Quote(profile.Name))
		for _, target := range profile.Targets {
			fmt.Fprintf(&output, "  target %s {\n", strconv.Quote(target.Target))
			hclStringsIndented(&output, "formats", target.Formats, "    ")
			fmt.Fprintf(&output, "    sign = %t\n    notarize = %t\n", target.Sign, target.Notarize)
			hclStringIndented(&output, "destination", target.Destination, "    ")
			fmt.Fprintln(&output, "  }")
		}
		fmt.Fprintln(&output, "}")
		fmt.Fprintln(&output)
	}
	return []byte(output.String()), nil
}

func writeEjectedPlatform(output *strings.Builder, name string, platform Platform, signing SigningPlatform) {
	if !platformConfigured(platform) && signing == (SigningPlatform{}) {
		return
	}
	fmt.Fprintf(output, "%s {\n", name)
	hclStringIndented(output, "product_name", platform.ProductName, "  ")
	hclStringIndented(output, "identifier", platform.Identifier, "  ")
	hclStringIndented(output, "minimum_version", platform.MinimumVersion, "  ")
	if platform.BuildNumber != 0 {
		fmt.Fprintf(output, "  build_number = %d\n", platform.BuildNumber)
	}
	hclStringIndented(output, "icon", platform.Icon, "  ")
	hclStringIndented(output, "manifest", platform.Manifest, "  ")
	hclStringIndented(output, "assets_car", platform.AssetsCar, "  ")
	hclStringIndented(output, "info_plist", platform.InfoPlist, "  ")
	hclStringIndented(output, "publisher", platform.Publisher, "  ")
	hclStringIndented(output, "version_name", platform.VersionName, "  ")
	if platform.VersionCode != 0 {
		fmt.Fprintf(output, "  version_code = %d\n", platform.VersionCode)
	}
	if platform.MinimumSDK != 0 {
		fmt.Fprintf(output, "  minimum_sdk = %d\n", platform.MinimumSDK)
	}
	if platform.TargetSDK != 0 {
		fmt.Fprintf(output, "  target_sdk = %d\n", platform.TargetSDK)
	}
	hclStringsIndented(output, "capabilities", platform.Capabilities, "  ")
	if signing != (SigningPlatform{}) {
		fmt.Fprintln(output, "  signing {")
		hclStringIndented(output, "credential", signing.Credential, "    ")
		hclStringIndented(output, "identity", signing.Identity, "    ")
		hclStringIndented(output, "certificate", signing.Certificate, "    ")
		hclStringIndented(output, "thumbprint", signing.Thumbprint, "    ")
		hclStringIndented(output, "timestamp_server", signing.TimestampServer, "    ")
		hclStringIndented(output, "entitlements", signing.Entitlements, "    ")
		hclStringIndented(output, "provisioning_profile", signing.ProvisioningProfile, "    ")
		hclStringIndented(output, "key_alias", signing.KeyAlias, "    ")
		fmt.Fprintln(output, "  }")
	}
	if signing.Notarize {
		fmt.Fprintln(output, "  notarization {")
		fmt.Fprintln(output, "  }")
	}
	fmt.Fprintln(output, "}")
	fmt.Fprintln(output)
}

func writeEjectedPackages(output *strings.Builder, packages Packages) {
	formats := []struct {
		name  string
		value PackageFormat
	}{
		{"nsis", packages.Windows.NSIS}, {"msix", packages.Windows.MSIX}, {"app", packages.Darwin.App}, {"dmg", packages.Darwin.DMG}, {"appimage", packages.Linux.AppImage}, {"deb", packages.Linux.Deb}, {"rpm", packages.Linux.RPM}, {"archlinux", packages.Linux.ArchLinux}, {"ipa", packages.IOS.IPA}, {"apk", packages.Android.APK}, {"aab", packages.Android.AAB},
	}
	for _, item := range formats {
		if item.value.Template == "" && len(item.value.Options) == 0 {
			continue
		}
		fmt.Fprintf(output, "package %s {\n", strconv.Quote(item.name))
		hclStringIndented(output, "template", item.value.Template, "  ")
		for _, key := range []string{"install_scope", "background", "volume_icon", "file_icon", "files", "categories", "pre_install", "post_install", "pre_remove", "post_remove"} {
			if value, ok := item.value.Options[key].(string); ok {
				hclStringIndented(output, key, value, "  ")
			}
		}
		for _, key := range []string{"window_width", "window_height"} {
			if value, ok := item.value.Options[key].(int); ok {
				fmt.Fprintf(output, "  %s = %d\n", key, value)
			}
		}
		fmt.Fprintln(output, "}")
		fmt.Fprintln(output)
	}
}

func platformConfigured(platform Platform) bool {
	return platform.ProductName != "" || platform.Identifier != "" || platform.MinimumVersion != "" || platform.BuildNumber != 0 || len(platform.Capabilities) > 0 || platform.Icon != "" || platform.Manifest != "" || platform.AssetsCar != "" || platform.InfoPlist != "" || platform.Publisher != "" || platform.VersionName != "" || platform.VersionCode != 0 || platform.MinimumSDK != 0 || platform.TargetSDK != 0
}

func hclString(output *strings.Builder, name, value string) {
	hclStringIndented(output, name, value, "  ")
}
func hclStringIndented(output *strings.Builder, name, value, indent string) {
	if value != "" {
		fmt.Fprintf(output, "%s%s = %s\n", indent, name, strconv.Quote(value))
	}
}
func hclStrings(output *strings.Builder, name string, values []string) {
	hclStringsIndented(output, name, values, "  ")
}
func hclStringsIndented(output *strings.Builder, name string, values []string, indent string) {
	if len(values) == 0 {
		return
	}
	quoted := make([]string, len(values))
	for index, value := range values {
		quoted[index] = strconv.Quote(value)
	}
	fmt.Fprintf(output, "%s%s = [%s]\n", indent, name, strings.Join(quoted, ", "))
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
