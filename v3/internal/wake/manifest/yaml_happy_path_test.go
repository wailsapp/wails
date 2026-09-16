package manifest

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const fullYAML = `version: 3
project:
  name: inventory
  product_name: Inventory
  identifier: com.example.inventory
  version: 2.4.1
  build_number: 42
  company: Example Inc
  binary_name: inventory
  icon: build/icon.png
  description: Inventory desktop application
  copyright: Copyright Example Inc
frontend:
  directory: ui
  install: [npm, install, --frozen-lockfile]
  build: [npm, run, bundle]
  dev: [npm, run, dev]
  output: ui/dist
build:
  output: artifacts
  tags: [desktop, release]
  trim_path: true
  strip: true
  obfuscated: false
  garble_args: [-tiny]
  ldflags: [-X example/build.version=2.4.1]
  compiler_flags: [all=-l]
dev:
  tags: [debug, devtools]
  debounce_ms: 125
  log_level: debug
  watch: ['**/*.go', wails.yaml]
  exclude: [tmp, ui/node_modules]
  use_git_ignore: false
  grace_period_ms: 900
windows:
  product_name: Inventory for Windows
  identifier: com.example.inventory.windows
  minimum_version: "10.0"
  build_number: 12
  capabilities: [internetClient]
  icon: build/windows.ico
  manifest: build/windows.manifest
  publisher: CN=Example Inc
  signing:
    credential: WINDOWS_CERT_PASSWORD
    identity: Example Inc
    certificate: build/windows.pfx
    thumbprint: ABC123
    timestamp_server: https://timestamp.example.test
    entitlements: build/windows.entitlements
darwin:
  product_name: Inventory for Mac
  identifier: com.example.inventory.mac
  minimum_version: "13.0"
  build_number: 13
  icon: build/mac.icns
  info_plist: build/Info.plist
  signing:
    credential: MACOS_KEYCHAIN_PROFILE
    identity: 'Developer ID Application: Example Inc'
    certificate: build/mac.p12
    entitlements: build/mac.entitlements
  notarization:
    credential: NOTARY_PROFILE
linux:
  product_name: Inventory for Linux
  identifier: com.example.inventory.linux
  minimum_version: "22.04"
  build_number: 14
  capabilities: [network]
  signing:
    credential: GPG_KEY_PASSWORD
    identity: release@example.test
    certificate: release-key
ios:
  display_name: Inventory Mobile
  bundle_id: com.example.inventory.ios
  minimum_version: "16.0"
  build_number: 15
  assets_car: ios/Assets.car
  info_plist: ios/Info.plist
  signing:
    credential: IOS_KEYCHAIN_PROFILE
    identity: 'Apple Distribution: Example Inc'
    certificate: ios/distribution.p12
    provisioning_profile: ios/profile.mobileprovision
    key_alias: distribution
android:
  display_name: Inventory Mobile
  application_id: com.example.inventory.android
  version_name: 2.4.1
  version_code: 241
  minimum_sdk: 26
  target_sdk: 35
  signing:
    credential: ANDROID_KEYSTORE_PASSWORD
    identity: release
    certificate: android/release.keystore
    key_alias: release
targets:
  windows/arm64:
    tags: [enterprise]
    minimum_version: "11.0"
  darwin/arm64:
    tags: [metal]
    minimum_version: "14.0"
  darwin/universal:
    tags: [universal]
  linux/arm64:
    tags: [portable]
    build_number: 99
  ios/arm64:
    minimum_version: "17.0"
  android/arm64:
    tags: [mobile]
packages:
  dmg:
    background: packaging/background.png
    window_width: 900
    window_height: 620
  nsis:
    install_scope: user
  msix: {}
  appimage: {}
  deb: {}
  rpm: {}
  archlinux: {}
file_associations:
  inventory:
    extensions: [inventory, .inv]
    name: Inventory document
    description: Inventory data file
    icon: build/inventory.icns
    role: Editor
    mime_type: application/x-inventory
    platforms: [darwin, windows]
protocols:
  inventory:
    description: Open an Inventory document
    platforms: [darwin, windows, linux]
profiles:
  release:
    targets:
      windows/amd64:
        formats: [nsis, msix]
        sign: true
      darwin/universal:
        formats: [dmg]
        sign: true
        notarize: true
      linux/arm64:
        formats: [appimage, deb]`

func TestLoadFullYAMLHappyPath(t *testing.T) {
	root := t.TempDir()
	writeManifest(t, root, fullYAML)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "packaging"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "packaging", "android.tmpl"), []byte("android template"), 0o644))
	loaded, err := Load(root, "release")
	require.NoError(t, err)

	assert.Equal(t, "inventory", loaded.Config.Project.Name)
	assert.Equal(t, 42, loaded.Config.Project.BuildNumber)
	assert.Equal(t, "Example Inc", loaded.Config.Project.CompanyName)
	assert.Equal(t, []string{"npm", "run", "bundle"}, loaded.Config.Frontend.Build)
	assert.Equal(t, "npm", loaded.Config.Frontend.PackageManager)
	assert.Equal(t, []string{"desktop", "release"}, loaded.Config.Build.Go.Tags)
	assert.Equal(t, []string{"-X example/build.version=2.4.1"}, loaded.Config.Build.Go.LinkerFlags)
	assert.Equal(t, []string{"all=-l"}, loaded.Config.Build.Go.CompilerFlags)
	assert.Equal(t, 125, loaded.Config.Dev.DebounceMS)
	assert.Equal(t, []string{"debug", "devtools"}, loaded.Config.Dev.Tags)
	assert.False(t, loaded.Config.Dev.UseGitIgnore)

	assert.Equal(t, "11.0", loaded.Config.Targets.Windows.ARM64.MinimumVersion)
	assert.Equal(t, []string{"enterprise"}, loaded.Config.Targets.Windows.ARM64.Tags)
	assert.Equal(t, "14.0", loaded.Config.Targets.Darwin.ARM64.MinimumVersion)
	assert.Equal(t, []string{"universal"}, loaded.Config.Targets.Darwin.Universal.Tags)
	assert.Equal(t, "17.0", loaded.Config.Targets.IOS.ARM64.MinimumVersion)
	assert.Equal(t, "build/windows.pfx", loaded.Config.Signing.Windows.Certificate)
	assert.True(t, loaded.Config.Signing.Windows.Enabled)
	assert.Equal(t, "MACOS_KEYCHAIN_PROFILE", loaded.Config.Signing.Darwin.Credential)
	assert.Equal(t, "NOTARY_PROFILE", loaded.Config.Signing.Darwin.NotarizationCredential)
	assert.True(t, loaded.Config.Signing.Darwin.Enabled)
	assert.True(t, loaded.Config.Signing.Darwin.Notarize)
	assert.Equal(t, "ios/profile.mobileprovision", loaded.Config.Signing.IOS.ProvisioningProfile)
	assert.Equal(t, "distribution", loaded.Config.Signing.IOS.KeyAlias)
	assert.Equal(t, "release", loaded.Config.Signing.Android.KeyAlias)
	assert.Equal(t, "com.example.inventory.android", loaded.Config.Targets.Android.Identifier)
	assert.Equal(t, 241, loaded.Config.Targets.Android.VersionCode)

	require.Len(t, loaded.Config.Associations, 1)
	assert.Equal(t, []string{"darwin", "windows"}, loaded.Config.Associations[0].Platforms)
	require.Len(t, loaded.Config.Protocols, 1)
	assert.Equal(t, "Open an Inventory document", loaded.Config.Protocols[0].Description)
	assert.Equal(t, "packaging/background.png", loaded.Config.Package.Darwin.DMG.Background)
	assert.Equal(t, 900, loaded.Config.Package.Darwin.DMG.WindowWidth)
	assert.Equal(t, "user", loaded.Config.Package.Windows.NSIS.InstallScope)

	require.Equal(t, "release", loaded.Config.Selected.Name)
	require.Len(t, loaded.Config.Selected.Targets, 3)
	assert.Empty(t, loaded.Config.Selected.Targets[1].Destination)
	assert.True(t, loaded.Config.Selected.Targets[1].Notarize)
}

func TestYAMLDevTagsRoundTripThroughEjection(t *testing.T) {
	root := t.TempDir()
	writeManifest(t, root, `version: 3
project:
  name: dev-tags
  product_name: Dev Tags
  identifier: com.example.devtags
  version: 1.0.0
dev:
  tags: [debug, sqlite]`)
	loaded, err := Load(root, "")
	require.NoError(t, err)
	assert.Equal(t, []string{"debug", "sqlite"}, loaded.Config.Dev.Tags)
	encoded, err := EncodeEjectedYAML(loaded.Config, "test")
	require.NoError(t, err)
	assert.Contains(t, string(encoded), `tags: [debug, sqlite]`)
}

func TestYAMLTypedPackageOptionsRoundTripThroughTheClosedSchema(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "packaging"), 0o755))
	for _, name := range []string{"background.png", "volume.icns", "file.icns", "README.md"} {
		require.NoError(t, os.WriteFile(filepath.Join(root, "packaging", name), []byte(name), 0o644))
	}
	writeManifest(t, root, `version: 3
project:
  name: options
  product_name: Options
  identifier: com.example.options
  version: 1.0.0
packages:
  dmg:
    background: packaging/background.png
    volume_icon: packaging/volume.icns
    file_icon: packaging/file.icns
    files:
      Read Me: packaging/README.md
    window_width: 900
    window_height: 620
  appimage:
    categories: [Development, IDE]`)
	loaded, err := Load(root, "")
	require.NoError(t, err)
	dmg := loaded.Config.Package.Darwin.DMG
	assert.Equal(t, "packaging/volume.icns", dmg.VolumeIcon)
	assert.Equal(t, "packaging/file.icns", dmg.FileIcon)
	assert.Equal(t, map[string]string{"Read Me": "packaging/README.md"}, dmg.Files)
	assert.Equal(t, []string{"Development", "IDE"}, loaded.Config.Package.Linux.AppImage.Categories)

	encoded, err := EncodeConfig(loaded.Config)
	require.NoError(t, err)
	assert.Contains(t, string(encoded), `volume_icon: packaging/volume.icns`)
	assert.Contains(t, string(encoded), `categories: [Development, IDE]`)
}

func TestEjectedYAMLRoundTripsCompleteBuildIntent(t *testing.T) {
	root := t.TempDir()
	writeManifest(t, root, fullYAML)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "packaging"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "packaging", "android.tmpl"), []byte("android template"), 0o644))
	loaded, err := Load(root, "release")
	require.NoError(t, err)

	data, err := EncodeEjectedYAML(loaded.Config, "3.0.0")
	require.NoError(t, err)
	ejected := filepath.Join(root, EjectedFilename)
	require.NoError(t, os.WriteFile(ejected, data, 0o644))
	reloaded, err := LoadFile(root, ejected, "")
	require.NoError(t, err)
	assert.Equal(t, loaded.Config.Project, reloaded.Config.Project)
	assert.Equal(t, loaded.Config.Frontend, reloaded.Config.Frontend)
	assert.Equal(t, loaded.Config.Build, reloaded.Config.Build)
	assert.Equal(t, loaded.Config.Dev, reloaded.Config.Dev)
	assert.Equal(t, loaded.Config.Targets.Windows.BuildNumber, reloaded.Config.Targets.Windows.BuildNumber)
	assert.Equal(t, loaded.Config.Targets.Linux.ARM64.BuildNumber, reloaded.Config.Targets.Linux.ARM64.BuildNumber)
	assert.Equal(t, loaded.Config.Associations, reloaded.Config.Associations)
	assert.Equal(t, loaded.Config.Protocols, reloaded.Config.Protocols)
	assert.Equal(t, loaded.Config.Package.Darwin.DMG, reloaded.Config.Package.Darwin.DMG)
	assert.Equal(t, loaded.Config.Package.Android.APK, reloaded.Config.Package.Android.APK)
	assert.Contains(t, string(data), "Generated by Wails CLI 3.0.0")
}

func TestYAMLWritersProduceLoadableManifests(t *testing.T) {
	root := t.TempDir()
	project := Project{Name: "writer", ProductName: "Writer", Identifier: "com.example.writer", Version: "1.0.0"}
	require.NoError(t, WriteMinimal(root, project))
	initial, err := os.ReadFile(filepath.Join(root, Filename))
	require.NoError(t, err)
	assert.Contains(t, string(initial), "# yaml-language-server: $schema="+JSONSchemaURL)
	assert.True(t, Exists(root))
	loaded, err := Load(root, "")
	require.NoError(t, err)

	encoded, err := EncodeConfig(loaded.Config)
	require.NoError(t, err)
	configPath := filepath.Join(root, "encoded.yaml")
	require.NoError(t, os.WriteFile(configPath, encoded, 0o644))
	reloaded, err := LoadFile(root, configPath, "")
	require.NoError(t, err)
	assert.Equal(t, loaded.Config.Project, reloaded.Config.Project)
	assert.Equal(t, loaded.Config.Frontend, reloaded.Config.Frontend)

	secondRoot := t.TempDir()
	doc := NewDocument(project)
	doc.Frontend.Directory = "web"
	doc.Frontend.OutputDirectory = "web/dist"
	require.NoError(t, WriteDocument(secondRoot, doc))
	second, err := Load(secondRoot, "")
	require.NoError(t, err)
	assert.Equal(t, "web", second.Config.Frontend.Directory)
	assert.Equal(t, "web/dist", second.Config.Frontend.OutputDirectory)
}

func TestYAMLWriterPreservesExplicitFalseBuildAndDevSettings(t *testing.T) {
	root := t.TempDir()
	project := Project{Name: "flags", ProductName: "Flags", Identifier: "com.example.flags", Version: "1.0.0"}
	doc := NewDocument(project)
	doc.Build.TrimPath = false
	doc.Build.Strip = false
	doc.Dev.UseGitIgnore = false
	doc.Frontend.Directory = "ui"
	require.NoError(t, WriteDocument(root, doc))
	loaded, err := Load(root, "")
	require.NoError(t, err)
	assert.False(t, loaded.Config.Build.TrimPath)
	assert.False(t, loaded.Config.Build.Strip)
	assert.False(t, loaded.Config.Dev.UseGitIgnore)
	assert.Equal(t, "ui", loaded.Config.Frontend.Directory)
}

func TestProjectOnlyYAMLUsesCompiledBuildDefaults(t *testing.T) {
	root := t.TempDir()
	writeManifest(t, root, `version: 3
project:
  name: defaults
  product_name: Defaults
  identifier: com.example.defaults
  version: 1.0.0`)
	loaded, err := Load(root, "")
	require.NoError(t, err)
	assert.Equal(t, "npm", loaded.Config.Frontend.PackageManager)
	assert.Equal(t, []string{"npm", "install"}, loaded.Config.Frontend.Install)
	assert.Equal(t, "dist", loaded.Config.Frontend.OutputDirectory)
	assert.Equal(t, "bin", loaded.Config.Build.OutputDirectory)
	assert.Equal(t, 9245, loaded.Config.Dev.Port)
}

func TestYAMLDiscoveryAcceptsAFileAndReturnsTheNearestManifest(t *testing.T) {
	root := t.TempDir()
	writeManifest(t, root, minimalYAML)
	file := filepath.Join(root, "frontend", "src", "main.ts")
	require.NoError(t, os.MkdirAll(filepath.Dir(file), 0o755))
	require.NoError(t, os.WriteFile(file, []byte("export {}\n"), 0o644))
	loaded, err := Load(file, "")
	require.NoError(t, err)
	assert.Equal(t, root, loaded.Config.Root)
	assert.True(t, Exists(root))
}

func TestYAMLDiscoveryStopsAtTheNearestGoModuleBoundary(t *testing.T) {
	parent := t.TempDir()
	writeManifest(t, parent, minimalYAML)
	module := filepath.Join(parent, "nested-module")
	require.NoError(t, os.MkdirAll(filepath.Join(module, "subdir"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(module, "go.mod"), []byte("module example.com/nested\n"), 0o644))
	_, _, err := Discover(filepath.Join(module, "subdir"))
	assert.ErrorIs(t, err, fs.ErrNotExist)
}

func TestYAMLValidationRejectsMalformedAndAmbiguousBuildRequests(t *testing.T) {
	projectOnly := `version: 3
project:
  name: app
  product_name: App
  identifier: com.example.app
  version: 1.0.0`
	validProfile := projectOnly + "\nprofiles:\n  release:\n    targets:\n      linux/amd64: {}\n"
	tests := []struct {
		name    string
		source  string
		profile string
	}{
		{name: "syntax", source: "version:"},
		{name: "missing project", source: "version: 3\n"},
		{name: "wrong version", source: strings.Replace(projectOnly, "version: 3", "version: 2", 1)},
		{name: "missing version", source: strings.TrimPrefix(projectOnly, "version: 3\n")},
		{name: "version must be first", source: "other: value\n" + projectOnly},
		{name: "default profile is reserved", source: validProfile, profile: "default"},
		{name: "profile must be lowercase slug", source: validProfile, profile: "Release"},
		{name: "profile must exist", source: validProfile, profile: "missing"},
		{name: "duplicate target", source: projectOnly + "\ntargets:\n  linux/amd64: {}\n  linux/amd64: {}\n"},
		{name: "unsupported target architecture", source: projectOnly + "\ntargets:\n  linux/riscv64: {}\n"},
		{name: "invalid universal target", source: projectOnly + "\ntargets:\n  windows/universal: {}\n"},
		{name: "duplicate package", source: projectOnly + "\npackages:\n  dmg: {}\n  dmg: {}\n"},
		{name: "unsupported package", source: projectOnly + "\npackages:\n  unknown: {}\n"},
		{name: "association requires extension", source: projectOnly + "\nfile_associations:\n  app: {}\n"},
		{name: "invalid profile name", source: projectOnly + "\nprofiles:\n  Release:\n    targets:\n      linux/amd64: {}\n"},
		{name: "duplicate profile", source: projectOnly + "\nprofiles:\n  release:\n    targets:\n      linux/amd64: {}\n  release:\n    targets:\n      darwin/amd64: {}\n"},
		{name: "profile duplicate target", source: projectOnly + "\nprofiles:\n  release:\n    targets:\n      linux/amd64: {}\n      linux/amd64: {}\n"},
		{name: "profile invalid target", source: projectOnly + "\nprofiles:\n  release:\n    targets:\n      linux/riscv64: {}\n"},
		{name: "duplicate protocol", source: projectOnly + "\nprotocols:\n  app: {}\n  app: {}\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			writeManifest(t, root, test.source)
			_, err := Load(root, test.profile)
			require.Error(t, err)
		})
	}
	missingRoot := t.TempDir()
	_, err := LoadFile(missingRoot, filepath.Join(missingRoot, "missing.yaml"), "")
	assert.Error(t, err)
}

func TestYAMLSchemaAcceptsAllDeclarativeValueShapes(t *testing.T) {
	for _, test := range []struct {
		source      string
		destination any
	}{
		{`text`, (*string)(nil)},
		{`TRUE`, (*bool)(nil)},
		{`[one, two]`, (*[]string)(nil)},
		{"key: value\n", (*map[string]string)(nil)},
	} {
		var node yaml.Node
		require.NoError(t, yaml.Unmarshal([]byte(test.source), &node))
		value := node.Content[0]
		_, err := decodeYAMLAttribute(value, reflect.TypeOf(test.destination))
		require.NoError(t, err, test.source)
	}
	target := Frontend{}
	applyFrontend(&target, &manifestFrontend{Build: &[]string{"pnpm", "run", "build"}})
	assert.Equal(t, "pnpm", target.PackageManager)
	target = Frontend{}
	applyFrontend(&target, &manifestFrontend{Dev: &[]string{"bun", "run", "dev"}})
	assert.Equal(t, "bun", target.PackageManager)
	_, err := documentFromManifest(manifestDocument{})
	assert.Error(t, err)
}

func TestYAMLEncoderUsesValidUniqueFileAssociationKeys(t *testing.T) {
	project := Project{Name: "writer", ProductName: "Writer", Identifier: "com.example.writer", Version: "1.0.0"}
	document := NewDocument(project)
	document.Associations = []Association{{Extensions: []string{".note"}}}
	encoded, err := EncodeDocument(document)
	require.NoError(t, err)
	assert.Contains(t, string(encoded), "  note:\n")

	root := t.TempDir()
	path := filepath.Join(root, "associations.yaml")
	require.NoError(t, os.WriteFile(path, encoded, 0o644))
	loaded, err := LoadFile(root, path, "")
	require.NoError(t, err)
	require.Len(t, loaded.Config.Associations, 1)
	assert.Equal(t, "note", loaded.Config.Associations[0].Name)

	document.Associations = []Association{
		{Name: "document", Extensions: []string{"doc"}},
		{Name: "document", Extensions: []string{"docx"}},
	}
	_, err = EncodeDocument(document)
	require.ErrorContains(t, err, "duplicate mapping key")
}

func TestYAMLDecoderRejectsNonLiteralCompositionFeatures(t *testing.T) {
	project := `version: 3
project:
  name: app
  product_name: App
  identifier: com.example.app
  version: 1.0.0
`
	for _, suffix := range []string{
		"build: &shared\n  output: bin\n",
		"build: *shared\n",
		"build:\n  <<: {output: bin}\n",
		"build: !custom {output: bin}\n",
		"build: null\n",
	} {
		_, err := decodeYAML(t.TempDir(), "wails.yaml", []byte(project+suffix), "")
		require.Error(t, err, suffix)
	}
}

func TestYAMLDocumentUsesAssociationLabelAndRejectsInvalidProjects(t *testing.T) {
	root := t.TempDir()
	writeManifest(t, root, `version: 3
project:
  name: labels
  product_name: Labels
  identifier: com.example.labels
  version: 1.0.0
file_associations:
  notes:
    extensions: [note]`)
	loaded, err := Load(root, "")
	require.NoError(t, err)
	require.Len(t, loaded.Config.Associations, 1)
	assert.Equal(t, "notes", loaded.Config.Associations[0].Name)

	writeManifest(t, root, `version: 3
project:
  name: incomplete`)
	_, err = Load(root, "")
	assert.ErrorContains(t, err, "project.product_name: required field is missing")
}

func TestYAMLTargetAndPackageNameParsersRejectMalformedNames(t *testing.T) {
	for _, value := range []string{"linux", "/amd64", "linux/", "linux/amd64/extra", "windows/universal", "linux/arm", "linux/386", "ios/amd64", "android/386", "plan9/amd64"} {
		_, _, err := parseTargetName(value)
		assert.Error(t, err, value)
	}

	for _, value := range []string{"unknown"} {
		_, err := ResolvePackageFormat(Packages{}, "linux", value)
		assert.Error(t, err, value)
	}
}

func TestYAMLDocumentValidationBranches(t *testing.T) {
	stringPointer := func(value string) *string { return &value }
	project := &manifestProject{Name: stringPointer("app"), ProductName: stringPointer("App"), Identifier: stringPointer("com.example.app"), Version: stringPointer("1.0.0")}
	target := manifestTarget{Name: "linux/amd64"}
	profile := func(name string, targets ...manifestProfileTarget) manifestProfile {
		return manifestProfile{Name: name, Targets: targets}
	}
	profileTarget := func(name string) manifestProfileTarget { return manifestProfileTarget{Name: name} }

	tests := []struct {
		name string
		raw  manifestDocument
	}{
		{name: "duplicate target", raw: manifestDocument{Project: project, Targets: []manifestTarget{target, target}}},
		{name: "duplicate package", raw: manifestDocument{Project: project, Packages: []PackageFormat{{Format: "deb"}, {Format: "deb"}}}},
		{name: "association needs extensions", raw: manifestDocument{Project: project, Associations: []manifestAssociation{{Label: "app"}}}},
		{name: "invalid profile name", raw: manifestDocument{Project: project, Profiles: []manifestProfile{profile("Release", profileTarget("linux/amd64"))}}},
		{name: "duplicate profile", raw: manifestDocument{Project: project, Profiles: []manifestProfile{profile("release", profileTarget("linux/amd64")), profile("release", profileTarget("darwin/amd64"))}}},
		{name: "profile needs targets", raw: manifestDocument{Project: project, Profiles: []manifestProfile{profile("release")}}},
		{name: "profile duplicate target", raw: manifestDocument{Project: project, Profiles: []manifestProfile{profile("release", profileTarget("linux/amd64"), profileTarget("linux/amd64"))}}},
		{name: "profile invalid target", raw: manifestDocument{Project: project, Profiles: []manifestProfile{profile("release", profileTarget("plan9/amd64"))}}},
		{name: "empty protocol", raw: manifestDocument{Project: project, Protocols: []manifestProtocol{{}}}},
		{name: "duplicate protocol", raw: manifestDocument{Project: project, Protocols: []manifestProtocol{{Scheme: "app"}, {Scheme: "app"}}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := documentFromManifest(test.raw)
			require.Error(t, err)
		})
	}
}

func TestYAMLProfileSelectionErrorsAreReported(t *testing.T) {
	root := t.TempDir()
	writeManifest(t, root, `version: 3
project:
  name: profiles
  product_name: Profiles
  identifier: com.example.profiles
  version: 1.0.0
profiles:
  release:
    targets:
      linux/amd64: {}`)
	for _, profile := range []string{"default", "Release", "missing"} {
		_, err := Load(root, profile)
		require.Error(t, err, profile)
	}
}

func TestLoadFileReportsWorkingDirectoryResolutionFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not allow removing the active working directory")
	}
	parent := t.TempDir()
	working := filepath.Join(parent, "working")
	draft := filepath.Join(parent, "draft.yaml")
	require.NoError(t, os.MkdirAll(working, 0o755))
	require.NoError(t, os.WriteFile(draft, []byte(minimalYAML), 0o644))
	t.Chdir(working)
	require.NoError(t, os.RemoveAll(working))
	_, err := LoadFile(".", draft, "")
	assert.Error(t, err)
}

func TestYAMLAcceptsExplicitCustomFrontendCommands(t *testing.T) {
	root := t.TempDir()
	source := `version: 3
project:
  name: custom
  product_name: Custom
  identifier: com.example.custom
  version: 1.0.0
frontend:
  install: [python3, install.py]
  build: [corepack, pnpm, build]
  dev: [./serve, --dev]`
	require.NoError(t, os.WriteFile(filepath.Join(root, Filename), []byte(source), 0644))
	loaded, err := Load(root, "")
	require.NoError(t, err)
	require.Equal(t, []string{"python3", "install.py"}, loaded.Config.Frontend.Install)
}
