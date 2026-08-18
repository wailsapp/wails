package manifest

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validHCLValidationConfig(root string) Config {
	project := Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}
	return configFromDocument(root, "", NewDocument(project))
}

func TestHCLValidationRejectsInvalidCoreConfiguration(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		name   string
		mutate func(*Config)
	}{
		{"missing binary name", func(config *Config) { config.Project.BinaryName = "" }},
		{"binary path", func(config *Config) { config.Project.BinaryName = "../app" }},
		{"unsupported package manager", func(config *Config) { config.Frontend.PackageManager = "deno" }},
		{"empty frontend directory", func(config *Config) { config.Frontend.Directory = "" }},
		{"empty frontend output", func(config *Config) { config.Frontend.OutputDirectory = "" }},
		{"empty build output", func(config *Config) { config.Build.OutputDirectory = "" }},
		{"empty install command", func(config *Config) { config.Frontend.InstallCommand = "" }},
		{"empty build command", func(config *Config) { config.Frontend.BuildCommand = "" }},
		{"empty dev command", func(config *Config) { config.Frontend.DevCommand = "" }},
		{"empty watch pattern", func(config *Config) { config.Dev.Watch = []string{""} }},
		{"invalid watch pattern", func(config *Config) { config.Dev.Watch = []string{"["} }},
		{"absolute project icon", func(config *Config) { config.Project.Icon = "/tmp/icon.png" }},
		{"escaping frontend directory", func(config *Config) { config.Frontend.Directory = "../frontend" }},
		{"escaping frontend output", func(config *Config) { config.Frontend.OutputDirectory = "../dist" }},
		{"escaping bindings output", func(config *Config) { config.Frontend.Bindings.OutputDirectory = "../bindings" }},
		{"escaping build output", func(config *Config) { config.Build.OutputDirectory = "../bin" }},
		{"absolute signing certificate", func(config *Config) { config.Signing.Linux.Certificate = "/tmp/key" }},
		{"escaping signing entitlements", func(config *Config) { config.Signing.Darwin.Entitlements = "../entitlements" }},
		{"association without extensions", func(config *Config) { config.Associations = []Association{{}} }},
		{"empty association extension", func(config *Config) { config.Associations = []Association{{Extensions: []string{"."}}} }},
		{"escaping association icon", func(config *Config) {
			config.Associations = []Association{{Extensions: []string{"txt"}, Icon: "../icon"}}
		}},
		{"unsupported association platform", func(config *Config) {
			config.Associations = []Association{{Extensions: []string{"txt"}, Platforms: []string{"plan9"}}}
		}},
		{"empty protocol scheme", func(config *Config) { config.Protocols = []Protocol{{Scheme: " "}} }},
		{"unsupported protocol platform", func(config *Config) { config.Protocols = []Protocol{{Scheme: "app", Platforms: []string{"plan9"}}} }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := validHCLValidationConfig(root)
			test.mutate(&config)
			require.Error(t, validateConfig(config))
		})
	}
}

func TestHCLValidationCoversPackageTemplatesAndOptions(t *testing.T) {
	root := t.TempDir()
	validTemplate := filepath.Join(root, "package.tmpl")
	require.NoError(t, os.WriteFile(validTemplate, []byte("template"), 0o644))

	tests := []struct {
		name   string
		format func(*Config)
	}{
		{"missing template", func(config *Config) { config.Package.Darwin.App.Template = "missing.tmpl" }},
		{"escaping template", func(config *Config) { config.Package.Darwin.App.Template = "../package.tmpl" }},
		{"unsupported options without template", func(config *Config) { config.Package.Windows.NSIS.Options = map[string]any{"unknown": true} }},
		{"unknown dmg option", func(config *Config) { config.Package.Darwin.DMG.Options = map[string]any{"unknown": true} }},
		{"invalid dmg option type", func(config *Config) { config.Package.Darwin.DMG.Options = map[string]any{"window_width": "wide"} }},
		{"escaping dmg resource", func(config *Config) {
			config.Package.Darwin.DMG.Options = map[string]any{"background": "../background.png"}
		}},
		{"invalid dmg file entry", func(config *Config) { config.Package.Darwin.DMG.Options = map[string]any{"files": "README.md"} }},
		{"escaping dmg file path", func(config *Config) {
			config.Package.Darwin.DMG.Options = map[string]any{"files": "README=../README.md"}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := validHCLValidationConfig(root)
			test.format(&config)
			require.Error(t, validateConfig(config))
		})
	}

	config := validHCLValidationConfig(root)
	config.Package.Darwin.App.Template = filepath.Base(validTemplate)
	assert.NoError(t, validateConfig(config))
}

func TestHCLCorePathHelpersCoverResolvedPathsAndHookOutputRoots(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "templates"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "templates", "package.tmpl"), []byte("template"), 0o644))
	assert.NoError(t, validateResolvedProjectPath(root, "templates/package.tmpl"))
	assert.Error(t, validateResolvedProjectPath(root, "templates/missing.tmpl"))
	assert.Error(t, validateResolvedProjectPath(filepath.Join(root, "missing"), "package.tmpl"))
	outside := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(outside, "package.tmpl"), []byte("outside"), 0o644))
	link := filepath.Join(root, "outside")
	if err := os.Symlink(outside, link); err == nil {
		assert.Error(t, validateResolvedProjectPath(root, "outside/package.tmpl"))
	}

	for _, test := range []struct {
		outputs []string
		want    string
	}{
		{[]string{"dist/app"}, "dist/app"},
		{[]string{"dist/app", "dist/assets/icon"}, "dist"},
		{[]string{"dist/app", "dist"}, "dist"},
	} {
		got, err := HookOutputRoot(test.outputs)
		require.NoError(t, err)
		assert.Equal(t, test.want, got)
	}
	for _, outputs := range [][]string{nil, {"."}, {"dist/app", "build/app"}} {
		_, err := HookOutputRoot(outputs)
		assert.Error(t, err)
	}
	assert.True(t, pathContains("dist", "dist"))
	assert.True(t, pathContains("dist", "dist/app"))
	assert.False(t, pathContains("dist", "build/app"))
}

func TestHCLDefaultsAndBinaryNameDerivationCoverAllArchitectures(t *testing.T) {
	platform := defaultPlatform("amd64", "arm64", "arm", "386", "unknown")
	assert.True(t, platform.AMD64.Enabled)
	assert.True(t, platform.ARM64.Enabled)
	assert.True(t, platform.ARM.Enabled)
	assert.True(t, platform.X86.Enabled)
	assert.Equal(t, "hello-world-3", deriveBinaryName(" Hello, World! 3 "))
	assert.Equal(t, "", deriveBinaryName("---"))
}

func TestHCLWritersAndDiscoveryReportInvalidRequests(t *testing.T) {
	root := t.TempDir()
	project := Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}
	invalid := Project{Name: "app"}
	assert.Error(t, WriteMinimal(root, invalid))
	_, encodeErr := EncodeDocument(Document{Project: invalid})
	assert.Error(t, encodeErr)
	assert.Error(t, WriteDocument(root, Document{Project: invalid}))

	_, _, err := Discover(filepath.Join(root, "missing"))
	assert.ErrorIs(t, err, fs.ErrNotExist)
	_, err = Load(filepath.Join(root, "missing"), "")
	assert.ErrorContains(t, err, "could not find")

	manifestDir := filepath.Join(root, Filename)
	require.NoError(t, os.Mkdir(manifestDir, 0o755))
	_, _, err = Discover(root)
	assert.ErrorContains(t, err, "is a directory")

	validRoot := t.TempDir()
	require.NoError(t, WriteMinimal(validRoot, project))
	assert.Error(t, Eject(validRoot, "release", "3.0.0", false))
	assert.Error(t, Eject(validRoot, "", "3.0.0", true))
	require.NoError(t, Eject(validRoot, "", "3.0.0", false))
	assert.Error(t, Eject(validRoot, "", "3.0.0", false))
}
