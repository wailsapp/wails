package manifest

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManifestAllowsStructuredPackageCustomizationWithoutTemplate(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := []byte(`version = 3

project {
  name = "package-options"
  product_name = "Package Options"
  identifier = "com.example.package-options"
  version = "1.0.0"
}

package "nsis" {
  install_scope = "user"
}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	loaded, err := LoadFile(root, path, "")
	require.NoError(t, err)
	assert.Equal(t, "user", loaded.Config.Package.Windows.NSIS.InstallScope)
	assert.Empty(t, loaded.Config.Package.Windows.NSIS.Template)
}

func TestPublishedPrototypeManifestsConform(t *testing.T) {
	prototypes, err := filepath.Glob(filepath.Join("testdata", "prototypes", "*.hcl"))
	require.NoError(t, err)
	require.Len(t, prototypes, 5)
	for _, prototype := range prototypes {
		prototype := prototype
		t.Run(filepath.Base(prototype), func(t *testing.T) {
			root := t.TempDir()
			template := filepath.Join(root, "packaging", "windows", "installer.nsi")
			require.NoError(t, os.MkdirAll(filepath.Dir(template), 0o755))
			require.NoError(t, os.WriteFile(template, []byte("; user-owned NSIS template\n"), 0o644))

			loaded, err := LoadFile(root, prototype, "")
			require.NoError(t, err)
			encoded, err := EncodeConfig(loaded.Config)
			require.NoError(t, err)
			roundTripped, err := decodeHCL(root, "round-trip.hcl", encoded, "")
			require.NoError(t, err)
			assert.Equal(t, loaded.Config.Project, roundTripped.Config.Project)
			assert.Equal(t, loaded.Config.Targets, roundTripped.Config.Targets)
			assert.Equal(t, loaded.Config.Package, roundTripped.Config.Package)
			assert.Equal(t, loaded.Config.Profiles, roundTripped.Config.Profiles)
		})
	}
}

func TestManifestRoundTripsLinuxDesktopEntryAsUserOwnedInput(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "assets"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "assets", "application.desktop"), []byte("[Desktop Entry]\n"), 0o644))
	path := filepath.Join(root, Filename)
	source := []byte(`version = 3

project {
  name = "desktop-entry"
  product_name = "Desktop Entry"
  identifier = "com.example.desktop-entry"
  version = "1.0.0"
}

linux {
  desktop_entry = "assets/application.desktop"
}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	loaded, err := LoadFile(root, path, "")
	require.NoError(t, err)
	assert.Equal(t, "assets/application.desktop", loaded.Config.Targets.Linux.DesktopEntry)

	encoded, err := EncodeConfig(loaded.Config)
	require.NoError(t, err)
	roundTripPath := filepath.Join(root, "round-trip.hcl")
	require.NoError(t, os.WriteFile(roundTripPath, encoded, 0o644))
	reloaded, err := LoadFile(root, roundTripPath, "")
	require.NoError(t, err)
	assert.Equal(t, loaded.Config.Targets.Linux.DesktopEntry, reloaded.Config.Targets.Linux.DesktopEntry)
}

func TestManifestValidationErrorsCarryTheOwningSourceRange(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := []byte(`version = 3

project {
  name = "diagnostics"
  product_name = "Diagnostics"
  identifier = "com.example.diagnostics"
  version = "1.0.0"
}

target "windows/amd64" {
  toolchain = "magic"
}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	_, err := LoadFile(root, path, "")
	require.Error(t, err)
	var validation *ValidationError
	require.True(t, errors.As(err, &validation))
	assert.Equal(t, `target["windows/amd64"].toolchain`, validation.Field)
	assert.Equal(t, SourceRange{Filename: path, StartLine: 11, StartColumn: 3, EndLine: 11, EndColumn: 22}, validation.Range)
}

func TestManifestDuplicateLabelsReportTheDuplicateRange(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := []byte(`version = 3

project {
  name = "duplicates"
  product_name = "Duplicates"
  identifier = "com.example.duplicates"
  version = "1.0.0"
}

target "windows/amd64" {}
target "windows/amd64" {}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	_, err := LoadFile(root, path, "")
	require.Error(t, err)
	var validation *ValidationError
	require.True(t, errors.As(err, &validation))
	assert.Equal(t, `target["windows/amd64"]`, validation.Field)
	assert.Equal(t, SourceRange{Filename: path, StartLine: 11, StartColumn: 8, EndLine: 11, EndColumn: 23}, validation.Range)
}

func TestManifestRequiredFieldsComeFromSchemaMetadata(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := []byte(`version = 3

project {
  product_name = "Required"
  identifier = "com.example.required"
  version = "1.0.0"
}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	_, err := LoadFile(root, path, "")
	require.Error(t, err)
	var validation *ValidationError
	require.True(t, errors.As(err, &validation))
	assert.Equal(t, "project.name", validation.Field)
	assert.Equal(t, SourceRange{Filename: path, StartLine: 3, StartColumn: 1, EndLine: 3, EndColumn: 8}, validation.Range)
}

func TestManifestResolvesProductionEnvironmentToolchainAndMobileFields(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := []byte(`version = 3

project {
  name = "complete"
  product_name = "Complete"
  identifier = "com.example.complete"
  version = "2.4.1"
}

frontend {
  environment = { PUBLIC_RELEASE = "true" }
  bindings {
    typescript = false
    interfaces = false
    output = "generated-bindings"
    models_filename = "domain"
    index_filename = "client"
    time_type = "Date"
  }
}

build {
  environment = { CGO_ENABLED = "0" }
  vcs_info = true
}

darwin {
  cf_bundle_icon_name = "AppIcon"
}

ios {
  bundle_id = "com.example.complete.ios"
  display_name = "Complete Mobile"
  company = "Example Ltd"
  comments = "iOS release"
  background_modes = ["fetch", "remote-notification"]
}

android {
  application_id = "com.example.complete.android"
  display_name = "Complete Android"
  company = "Example Ltd"
  comments = "Android release"
  version_name = "2.4.1"
  version_code = 241
  minimum_sdk = 26
  target_sdk = 35
}

target "windows/amd64" {
  toolchain = "zig"
  environment = { CC = "zig cc" }
  tags = ["enterprise"]
  ldflags = ["-X example.com/build.version=2.4.1"]
  compiler_flags = ["all=-l"]
  garble_args = ["-literals"]
  obfuscated = true
}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	loaded, err := LoadFile(root, path, "")
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"PUBLIC_RELEASE": "true"}, loaded.Config.Frontend.Environment)
	assert.Equal(t, Bindings{TypeScript: false, Interfaces: false, OutputDirectory: "generated-bindings", ModelsFilename: "domain", IndexFilename: "client", TimeType: "Date"}, loaded.Config.Frontend.Bindings)
	assert.Equal(t, map[string]string{"CGO_ENABLED": "0"}, loaded.Config.Build.Environment)
	assert.True(t, loaded.Config.Build.VCSInfo)
	assert.Equal(t, "AppIcon", loaded.Config.Targets.Darwin.CFBundleIconName)
	assert.Equal(t, "com.example.complete.ios", loaded.Config.Targets.IOS.Identifier)
	assert.Equal(t, "Complete Mobile", loaded.Config.Targets.IOS.ProductName)
	assert.Equal(t, "Example Ltd", loaded.Config.Targets.IOS.CompanyName)
	assert.Equal(t, "iOS release", loaded.Config.Targets.IOS.Comments)
	assert.Equal(t, []string{"fetch", "remote-notification"}, loaded.Config.Targets.IOS.BackgroundModes)
	assert.Equal(t, "com.example.complete.android", loaded.Config.Targets.Android.Identifier)
	assert.Equal(t, "Complete Android", loaded.Config.Targets.Android.ProductName)
	assert.Equal(t, "Example Ltd", loaded.Config.Targets.Android.CompanyName)
	assert.Equal(t, "Android release", loaded.Config.Targets.Android.Comments)
	target := loaded.Config.Targets.Windows.AMD64
	assert.Equal(t, "zig", target.Toolchain)
	assert.Equal(t, map[string]string{"CC": "zig cc"}, target.Environment)
	assert.Equal(t, []string{"enterprise"}, target.Tags)
	assert.Equal(t, []string{"-X example.com/build.version=2.4.1"}, target.LinkerFlags)
	assert.Equal(t, []string{"all=-l"}, target.CompilerFlags)
	assert.Equal(t, []string{"-literals"}, target.GarbleArgs)
	assert.True(t, target.Obfuscated)

	encoded, err := EncodeConfig(loaded.Config)
	require.NoError(t, err)
	roundTripped, err := decodeHCL(root, "round-trip.hcl", encoded, "")
	require.NoError(t, err)
	assert.Equal(t, loaded.Config.Frontend.Environment, roundTripped.Config.Frontend.Environment)
	assert.Equal(t, loaded.Config.Build.Environment, roundTripped.Config.Build.Environment)
	assert.Equal(t, loaded.Config.Build.VCSInfo, roundTripped.Config.Build.VCSInfo)
	assert.Equal(t, loaded.Config.Targets.Darwin.CFBundleIconName, roundTripped.Config.Targets.Darwin.CFBundleIconName)
	assert.Equal(t, loaded.Config.Targets.IOS, roundTripped.Config.Targets.IOS)
	assert.Equal(t, loaded.Config.Targets.Android, roundTripped.Config.Targets.Android)
	assert.Equal(t, target, roundTripped.Config.Targets.Windows.AMD64)
}

func TestManifestRecordsExactConfigurationOrigins(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := []byte(`version = 3

project {
  name = "origins"
  product_name = "Origins"
  identifier = "com.example.origins"
  version = "1.0.0"
}

build {
  output = "release"
}

target "windows/amd64" {
  toolchain = "zig"
}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	loaded, err := LoadFile(root, path, "")
	require.NoError(t, err)
	assert.Equal(t, OriginManifest, loaded.Config.Origins["build.output"].Kind)
	assert.Equal(t, SourceRange{Filename: path, StartLine: 11, StartColumn: 3, EndLine: 11, EndColumn: 21}, loaded.Config.Origins["build.output"].Range)
	assert.Equal(t, OriginManifest, loaded.Config.Origins[`target["windows/amd64"].toolchain`].Kind)
	assert.Equal(t, 15, loaded.Config.Origins[`target["windows/amd64"].toolchain`].Range.StartLine)
	assert.Equal(t, OriginDefault, loaded.Config.Origins["build.trim_path"].Kind)
}

func TestSchemaReferenceDerivesFieldsTypesAndDefaults(t *testing.T) {
	fields := SchemaReference()
	require.NotEmpty(t, fields)
	byPath := make(map[string]SchemaField, len(fields))
	previous := ""
	for _, field := range fields {
		assert.Greater(t, field.Path, previous)
		_, duplicate := byPath[field.Path]
		assert.False(t, duplicate, field.Path)
		byPath[field.Path] = field
		previous = field.Path
	}
	assert.Equal(t, "string", byPath["build.output"].Type)
	assert.Equal(t, "bin", byPath["build.output"].Default)
	assert.Equal(t, `"bin"`, byPath["build.output"].Example)
	assert.NotEmpty(t, byPath["build.output"].Description)
	assert.Equal(t, "map(string)", byPath["frontend.environment"].Type)
	assert.Equal(t, "string", byPath[`target["target"].toolchain`].Type)
	assert.Equal(t, "list(string)", byPath[`profile["profile"].target["target"].formats`].Type)
	assert.Equal(t, "string", byPath[`package["format"].install_scope`].Type)
	assert.True(t, byPath["version"].Required)
	assert.Contains(t, string(SchemaReferenceMarkdown()), "| `build.output` | string | no | `bin` | `\"bin\"` |")
	assert.Equal(t, "nsis", byPath[`package["format"].install_scope`].Formats)
}

func TestManifestKeepsSigningAndNotarizationCredentialsDistinct(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := []byte(`version = 3

project {
  name = "credentials"
  product_name = "Credentials"
  identifier = "com.example.credentials"
  version = "1.0.0"
}

darwin {
  identifier = "com.example.credentials.macos"
  signing {
    credential = "apple-signing"
    identity = "Developer ID Application: Example"
  }
  notarization {
    credential = "apple-notary"
  }
}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	loaded, err := LoadFile(root, path, "")
	require.NoError(t, err)
	assert.Equal(t, "apple-signing", loaded.Config.Signing.Darwin.Credential)
	assert.Equal(t, "apple-notary", loaded.Config.Signing.Darwin.NotarizationCredential)

	encoded, err := EncodeConfig(loaded.Config)
	require.NoError(t, err)
	assert.Contains(t, string(encoded), "credential = \"apple-signing\"")
	assert.Contains(t, string(encoded), "credential = \"apple-notary\"")
	roundTripped, err := decodeHCL(root, "round-trip.hcl", encoded, "")
	require.NoError(t, err)
	assert.Equal(t, loaded.Config.Signing.Darwin, roundTripped.Config.Signing.Darwin)
}

func TestManifestEjectionPreservesExplicitFalseTargetOverrides(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := []byte(`version = 3

project {
  name = "false-override"
  product_name = "False Override"
  identifier = "com.example.false-override"
  version = "1.0.0"
}

build {
  obfuscated = true
}

target "linux/amd64" {
  obfuscated = false
}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))
	loaded, err := LoadFile(root, path, "")
	require.NoError(t, err)
	require.True(t, loaded.Config.Targets.Linux.AMD64.ObfuscatedSet)
	require.False(t, loaded.Config.Targets.Linux.AMD64.Obfuscated)

	encoded, err := EncodeConfig(loaded.Config)
	require.NoError(t, err)
	assert.Contains(t, string(encoded), "target \"linux/amd64\"")
	assert.Contains(t, string(encoded), "obfuscated = false")
	roundTripped, err := decodeHCL(root, "round-trip.hcl", encoded, "")
	require.NoError(t, err)
	assert.True(t, roundTripped.Config.Targets.Linux.AMD64.ObfuscatedSet)
	assert.False(t, roundTripped.Config.Targets.Linux.AMD64.Obfuscated)
}

func TestManifestRoundTripPreservesExplicitEmptyWatcherCollections(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := []byte(`version = 3

project {
  name = "empty-watch"
  product_name = "Empty Watch"
  identifier = "com.example.empty-watch"
  version = "1.0.0"
}

dev {
  watch = []
  exclude = []
}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))
	loaded, err := LoadFile(root, path, "")
	require.NoError(t, err)
	require.Empty(t, loaded.Config.Dev.Watch)
	require.Empty(t, loaded.Config.Dev.Exclude)

	encoded, err := EncodeConfig(loaded.Config)
	require.NoError(t, err)
	assert.Contains(t, string(encoded), "  watch = []\n")
	assert.Contains(t, string(encoded), "  exclude = []\n")
	roundTripped, err := decodeHCL(root, "round-trip.hcl", encoded, "")
	require.NoError(t, err)
	assert.Empty(t, roundTripped.Config.Dev.Watch)
	assert.Empty(t, roundTripped.Config.Dev.Exclude)
}

func TestManifestRoundTripsTypedPackageConfiguration(t *testing.T) {
	root := t.TempDir()
	for name, contents := range map[string]string{
		"AppxManifest.xml":    "<Package/>",
		"background.png":      "background",
		"volume.icns":         "volume",
		"file.icns":           "file",
		"LICENSE.txt":         "license",
		"application.png":     "icon",
		"application.desktop": "[Desktop Entry]\n",
		"preinstall.sh":       "#!/bin/sh\n",
	} {
		require.NoError(t, os.WriteFile(filepath.Join(root, name), []byte(contents), 0o644))
	}
	path := filepath.Join(root, Filename)
	source := []byte(`version = 3

project {
  name = "typed-packages"
  product_name = "Typed Packages"
  identifier = "com.example.typed-packages"
  version = "1.0.0"
}

package "msix" {
  publisher = "CN=Example"
  manifest = "AppxManifest.xml"
}

package "dmg" {
  background = "background.png"
  volume_icon = "volume.icns"
  file_icon = "file.icns"
  files = { "License" = "LICENSE.txt" }
  window_width = 640
  window_height = 480
}

package "appimage" {
  icon = "application.png"
  desktop_entry = "application.desktop"
  categories = ["Development", "IDE"]
}

package "deb" {
  maintainer = "Example <release@example.com>"
  section = "utils"
  dependencies = ["libgtk-4-1", "libwebkitgtk-6.0-4"]
  pre_install = "preinstall.sh"
}
`)
	require.NoError(t, os.WriteFile(path, source, 0o644))

	loaded, err := LoadFile(root, path, "")
	require.NoError(t, err)
	assert.Equal(t, "CN=Example", loaded.Config.Package.Windows.MSIX.Publisher)
	assert.Equal(t, "AppxManifest.xml", loaded.Config.Package.Windows.MSIX.Manifest)
	assert.Equal(t, map[string]string{"License": "LICENSE.txt"}, loaded.Config.Package.Darwin.DMG.Files)
	assert.Equal(t, []string{"Development", "IDE"}, loaded.Config.Package.Linux.AppImage.Categories)
	assert.Equal(t, []string{"libgtk-4-1", "libwebkitgtk-6.0-4"}, loaded.Config.Package.Linux.Deb.Dependencies)

	encoded, err := EncodeConfig(loaded.Config)
	require.NoError(t, err)
	roundTripped, err := decodeHCL(root, "round-trip.hcl", encoded, "")
	require.NoError(t, err)
	assert.Equal(t, loaded.Config.Package, roundTripped.Config.Package)
}

func TestEjectionCoversEveryPlatformSchemaField(t *testing.T) {
	project := Project{Name: "schema-ejection", ProductName: "Schema Ejection", Identifier: "com.example.schema-ejection", Version: "1.0.0"}
	doc := defaults(project)
	all := Platform{
		ProductName: "Product", Identifier: "com.example.platform", MinimumVersion: "1.0", BuildNumber: 7,
		Capabilities: []string{"capability"}, Icon: "icon", Manifest: "manifest", AssetsCar: "Assets.car",
		InfoPlist: "Info.plist", Publisher: "CN=Publisher", DesktopEntry: "application.desktop",
		CompanyName: "Company", Comments: "Comments", CFBundleIconName: "AppIcon",
		BackgroundModes: []string{"audio"}, VersionName: "1.0", VersionCode: 7, MinimumSDK: 24, TargetSDK: 35,
	}
	doc.Targets.Windows, doc.Targets.Darwin, doc.Targets.Linux, doc.Targets.IOS, doc.Targets.Android = all, all, all, all, all
	signing := SigningPlatform{Enabled: true, Identity: "identity", Certificate: "certificate", Thumbprint: "thumbprint", TimestampServer: "timestamp", Entitlements: "entitlements", ProvisioningProfile: "profile", KeyAlias: "alias", Credential: "credential", Notarize: true, NotarizationCredential: "notary"}
	doc.Signing.Windows, doc.Signing.Darwin, doc.Signing.Linux, doc.Signing.IOS, doc.Signing.Android = signing, signing, signing, signing, signing

	encoded, err := EncodeConfig(configFromDocument(t.TempDir(), "", doc))
	require.NoError(t, err)
	body := parseEjectedBody(t, encoded)
	for _, platform := range []string{"windows", "darwin", "linux", "ios", "android"} {
		block := requireEjectedBlock(t, body, platform, "")
		descriptor := manifestSchema.blocks[platform].node
		var expectedAttributes []string
		for _, name := range descriptor.attributeOrder {
			if schemaFieldAllowed(platform, descriptor.attributes[name].platformMask) {
				expectedAttributes = append(expectedAttributes, name)
			}
		}
		sort.Strings(expectedAttributes)
		assert.Equal(t, expectedAttributes, sortedEjectedAttributeNames(block.Body), platform)

		signingBlock := requireEjectedBlock(t, block.Body, "signing", "")
		expectedSigning := append([]string(nil), manifestSchema.blocks[platform].node.blocks["signing"].node.attributeOrder...)
		sort.Strings(expectedSigning)
		assert.Equal(t, expectedSigning, sortedEjectedAttributeNames(signingBlock.Body), platform+" signing")
		if platform == "darwin" {
			notarization := requireEjectedBlock(t, block.Body, "notarization", "")
			assert.Equal(t, []string{"credential"}, sortedEjectedAttributeNames(notarization.Body))
		} else {
			assert.Nil(t, findEjectedBlock(block.Body, "notarization", ""), platform)
		}
	}
}

func TestEjectionCoversEveryPackageSchemaField(t *testing.T) {
	project := Project{Name: "package-ejection", ProductName: "Package Ejection", Identifier: "com.example.package-ejection", Version: "1.0.0"}
	packageNode := schemaNodesByType[reflect.TypeOf(PackageFormat{})]
	for _, format := range []string{"nsis", "msix", "dmg", "appimage", "deb", "rpm", "archlinux"} {
		t.Run(format, func(t *testing.T) {
			doc := defaults(project)
			raw := PackageFormat{
				Format: format, InstallScope: "user", Publisher: "CN=Publisher", Manifest: "manifest",
				Background: "background", VolumeIcon: "volume", FileIcon: "file", Files: map[string]string{"License": "LICENSE"},
				WindowWidth: 640, WindowHeight: 480, Icon: "icon", DesktopEntry: "application.desktop",
				Categories: []string{"Utility"}, Maintainer: "Maintainer", Section: "utils", Dependencies: []string{"dependency"},
				PreInstall: "preinstall", PostInstall: "postinstall", PreRemove: "preremove", PostRemove: "postremove",
			}
			require.NoError(t, applyPackage(&doc.Package, raw))
			encoded, err := EncodeConfig(configFromDocument(t.TempDir(), "", doc))
			require.NoError(t, err)
			block := requireEjectedBlock(t, parseEjectedBody(t, encoded), "package", format)
			var expected []string
			for _, name := range packageNode.attributeOrder {
				descriptor := packageNode.attributes[name]
				if name != "template" && schemaFormatNameAllowed(format, descriptor.formatMask) {
					expected = append(expected, name)
				}
			}
			sort.Strings(expected)
			assert.Equal(t, expected, sortedEjectedAttributeNames(block.Body))
		})
	}

	for _, format := range []string{"nsis", "dmg", "deb", "rpm", "archlinux"} {
		t.Run(format+" template", func(t *testing.T) {
			doc := defaults(project)
			require.NoError(t, applyPackage(&doc.Package, PackageFormat{Format: format, Template: "package.template"}))
			encoded, err := EncodeConfig(configFromDocument(t.TempDir(), "", doc))
			require.NoError(t, err)
			block := requireEjectedBlock(t, parseEjectedBody(t, encoded), "package", format)
			assert.Equal(t, []string{"template"}, sortedEjectedAttributeNames(block.Body))
		})
	}
}

func TestSemanticProfileErrorsCarryTheOwningAttributeRange(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := `version = 3

project {
  name = "diagnostics"
  product_name = "Diagnostics"
  identifier = "com.example.diagnostics"
  version = "1.0.0"
}

profile "release" {
  target "linux/amd64" {
    formats = ["aab"]
  }
}
`
	require.NoError(t, os.WriteFile(path, []byte(source), 0o644))

	_, err := LoadFile(root, path, "")
	require.Error(t, err)
	var validation *ValidationError
	require.ErrorAs(t, err, &validation)
	assert.Equal(t, `profile["release"].target["linux/amd64"].formats`, validation.Field)
	assert.Equal(t, path, validation.Range.Filename)
	assert.Equal(t, 12, validation.Range.StartLine)
}

func TestSemanticTargetErrorsCarryTheOwningLabelRange(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, Filename)
	source := `version = 3

project {
  name = "diagnostics"
  product_name = "Diagnostics"
  identifier = "com.example.diagnostics"
  version = "1.0.0"
}

profile "release" {
  target "plan9/amd64" {}
}
`
	require.NoError(t, os.WriteFile(path, []byte(source), 0o644))

	_, err := LoadFile(root, path, "")
	require.Error(t, err)
	var validation *ValidationError
	require.ErrorAs(t, err, &validation)
	assert.Equal(t, `profile["release"].target["plan9/amd64"]`, validation.Field)
	assert.Equal(t, path, validation.Range.Filename)
	assert.Equal(t, 11, validation.Range.StartLine)
}

func parseEjectedBody(t *testing.T, source []byte) *hclsyntax.Body {
	t.Helper()
	file, diagnostics := hclsyntax.ParseConfig(source, "ejected.hcl", hcl.InitialPos)
	require.False(t, diagnostics.HasErrors(), diagnostics.Error())
	return file.Body.(*hclsyntax.Body)
}

func requireEjectedBlock(t *testing.T, body *hclsyntax.Body, blockType, label string) *hclsyntax.Block {
	t.Helper()
	block := findEjectedBlock(body, blockType, label)
	require.NotNil(t, block, "%s %q", blockType, label)
	return block
}

func findEjectedBlock(body *hclsyntax.Body, blockType, label string) *hclsyntax.Block {
	for _, block := range body.Blocks {
		if block.Type == blockType && (label == "" && len(block.Labels) == 0 || len(block.Labels) == 1 && block.Labels[0] == label) {
			return block
		}
	}
	return nil
}

func sortedEjectedAttributeNames(body *hclsyntax.Body) []string {
	result := make([]string, 0, len(body.Attributes))
	for name := range body.Attributes {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}
