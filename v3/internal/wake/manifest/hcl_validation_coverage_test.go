package manifest

import (
	"fmt"
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
		{"unsupported binding time", func(config *Config) { config.Frontend.Bindings.TimeType = "Temporal" }},
		{"invalid frontend environment", func(config *Config) { config.Frontend.Environment = map[string]string{"BAD=NAME": "x"} }},
		{"invalid build environment", func(config *Config) { config.Build.Environment = map[string]string{"": "x"} }},
		{"invalid target environment", func(config *Config) { config.Targets.Windows.AMD64.Environment = map[string]string{"BAD=NAME": "x"} }},
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

func TestHCLValidationRejectsAnUnresolvableProjectRoot(t *testing.T) {
	config := validHCLValidationConfig(filepath.Join(t.TempDir(), "missing"))
	require.Error(t, validateConfig(config))
}

func TestHCLValidationCoversPackageTemplatesAndOptions(t *testing.T) {
	root := t.TempDir()
	validTemplate := filepath.Join(root, "package.tmpl")
	require.NoError(t, os.WriteFile(validTemplate, []byte("template"), 0o644))

	tests := []struct {
		name   string
		format func(*Config)
	}{
		{"missing template", func(config *Config) { config.Package.Windows.NSIS.Template = "missing.tmpl" }},
		{"escaping template", func(config *Config) { config.Package.Windows.NSIS.Template = "../package.tmpl" }},
		{"replacement conflicts with structured options", func(config *Config) {
			config.Package.Windows.NSIS.Template = filepath.Base(validTemplate)
			config.Package.Windows.NSIS.InstallScope = "user"
		}},
		{"escaping dmg resource", func(config *Config) {
			config.Package.Darwin.DMG.Background = "../background.png"
		}},
		{"invalid dmg file entry", func(config *Config) { config.Package.Darwin.DMG.Files = map[string]string{"": "README.md"} }},
		{"escaping dmg file path", func(config *Config) {
			config.Package.Darwin.DMG.Files = map[string]string{"README": "../README.md"}
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
	config.Package.Windows.NSIS.Template = filepath.Base(validTemplate)
	assert.NoError(t, validateConfig(config))
	paths, err := newProjectPathValidator(root)
	require.NoError(t, err)
	require.Error(t, validatePackageOptions(paths, "windows.nsis", PackageFormat{Publisher: "CN=Example"}))
	require.Error(t, validatePackageOptions(paths, "windows.msix", PackageFormat{Format: "dmg"}))
}

func TestHCLValidationRejectsEveryConfigurablePathIntoWails(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		field  string
		mutate func(*Config)
	}{
		{"project.icon", func(c *Config) { c.Project.Icon = ".wails/icon.png" }},
		{"frontend.directory", func(c *Config) { c.Frontend.Directory = ".wails/frontend" }},
		{"frontend.output", func(c *Config) { c.Frontend.OutputDirectory = ".wails/dist" }},
		{"frontend.bindings.output", func(c *Config) { c.Frontend.Bindings.OutputDirectory = ".wails/bindings" }},
		{"build.output", func(c *Config) { c.Build.OutputDirectory = ".wails/bin" }},
		{"windows.icon", func(c *Config) { c.Targets.Windows.Icon = ".wails/icon.ico" }},
		{"windows.manifest", func(c *Config) { c.Targets.Windows.Manifest = ".wails/app.manifest" }},
		{"windows.assets_car", func(c *Config) { c.Targets.Windows.AssetsCar = ".wails/assets.car" }},
		{"windows.info_plist", func(c *Config) { c.Targets.Windows.InfoPlist = ".wails/Info.plist" }},
		{"windows.signing.certificate", func(c *Config) { c.Signing.Windows.Certificate = ".wails/cert.pfx" }},
		{"darwin.signing.entitlements", func(c *Config) { c.Signing.Darwin.Entitlements = ".wails/app.entitlements" }},
		{"ios.signing.provisioning_profile", func(c *Config) { c.Signing.IOS.ProvisioningProfile = ".wails/app.mobileprovision" }},
		{`file_association[""].icon`, func(c *Config) {
			c.Associations = []Association{{Extensions: []string{"txt"}, Icon: ".wails/file.ico"}}
		}},
		{`package["msix"].manifest`, func(c *Config) { c.Package.Windows.MSIX.Manifest = ".wails/AppxManifest.xml" }},
		{`package["dmg"].background`, func(c *Config) { c.Package.Darwin.DMG.Background = ".wails/background.png" }},
		{`package["dmg"].files`, func(c *Config) { c.Package.Darwin.DMG.Files = map[string]string{"Read Me": ".wails/README.md"} }},
	}
	for _, test := range tests {
		t.Run(test.field, func(t *testing.T) {
			config := validHCLValidationConfig(root)
			test.mutate(&config)
			err := validateConfig(config)
			require.Error(t, err)
			assert.Contains(t, err.Error(), test.field)
			assert.Contains(t, err.Error(), ".wails")
		})
	}
}

func TestProjectPathValidationIsCrossPlatformAndChecksSymlinkedParents(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "existing"), []byte("asset"), 0o644))
	resolvedRoot, err := filepath.EvalSymlinks(root)
	require.NoError(t, err)
	resolved, err := ResolveProjectPath(root, "source.field", "existing", true)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(resolvedRoot, "existing"), resolved)
	resolved, err = ResolveProjectPath(root, "source.field", "future/output", false)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "future", "output"), resolved)

	for _, value := range []string{"/tmp/file", `C:\temp\file`, `\\server\share\file`, "../file", `..\file`, ".wails/file", `folder\.WAILS\file`} {
		err := validateProjectPath(root, "source.field", value, false)
		require.Error(t, err, value)
		assert.Contains(t, err.Error(), "source.field")
	}
	assert.NoError(t, validateProjectPath(root, "source.field", "assets/.wails-icons/icon.png", false))

	outside := t.TempDir()
	require.NoError(t, os.Symlink(outside, filepath.Join(root, "linked")))
	err = validateProjectPath(root, "source.field", "linked/not-created-yet", false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source.field")
}

func TestProjectPathValidatorFaultsAndCacheAreCovered(t *testing.T) {
	root := t.TempDir()
	_, err := newProjectPathValidatorWithOperations(root, projectPathOperations{eval: func(string) (string, error) { return "", fs.ErrPermission }})
	require.ErrorIs(t, err, fs.ErrPermission)
	require.ErrorIs(t, validateProjectPath(filepath.Join(root, "missing"), "source.field", "asset", false), fs.ErrNotExist)

	validator, err := newProjectPathValidator(root)
	require.NoError(t, err)
	require.NoError(t, validator.validate("source.field", "missing/file", false))
	require.ErrorIs(t, validator.validate("source.field", "missing/file", true), fs.ErrNotExist)

	defaultOps := projectPathOperations{eval: filepath.EvalSymlinks, lstat: os.Lstat, rel: filepath.Rel}
	t.Run("lstat", func(t *testing.T) {
		ops := defaultOps
		ops.lstat = func(string) (fs.FileInfo, error) { return nil, fs.ErrPermission }
		validator, err := newProjectPathValidatorWithOperations(root, ops)
		require.NoError(t, err)
		require.ErrorIs(t, validator.validate("source.field", "asset", false), fs.ErrPermission)
	})
	t.Run("eval", func(t *testing.T) {
		asset := filepath.Join(root, "asset")
		require.NoError(t, os.WriteFile(asset, []byte("asset"), 0o644))
		ops := defaultOps
		ops.eval = func(path string) (string, error) {
			if path == asset {
				return "", fs.ErrPermission
			}
			return filepath.EvalSymlinks(path)
		}
		validator, err := newProjectPathValidatorWithOperations(root, ops)
		require.NoError(t, err)
		require.ErrorIs(t, validator.validate("source.field", "asset", true), fs.ErrPermission)
	})
	t.Run("relative", func(t *testing.T) {
		ops := defaultOps
		ops.rel = func(string, string) (string, error) { return "", fs.ErrInvalid }
		validator, err := newProjectPathValidatorWithOperations(root, ops)
		require.NoError(t, err)
		require.ErrorIs(t, validator.validate("source.field", "asset", false), fs.ErrInvalid)
	})
	validator = &projectPathValidator{resolved: map[string]pathResolution{}, ops: projectPathOperations{lstat: func(string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }}}
	_, err = validator.resolve(string(filepath.Separator), false)
	require.ErrorIs(t, err, fs.ErrNotExist)

	config := validHCLValidationConfig(root)
	config.Package.Darwin.DMG.WindowWidth = 700
	for _, name := range []string{"first.txt", "second.txt"} {
		require.NoError(t, os.WriteFile(filepath.Join(root, name), []byte(name), 0o644))
	}
	config.Package.Darwin.DMG.Files = map[string]string{"Second": "second.txt", "First": "first.txt"}
	require.NoError(t, validateConfig(config))
}

func BenchmarkValidateLargePathManifest(b *testing.B) {
	root := b.TempDir()
	config := validHCLValidationConfig(root)
	config.Associations = make([]Association, 1000)
	for index := range config.Associations {
		config.Associations[index] = Association{Extensions: []string{"ext"}, Icon: filepath.ToSlash(filepath.Join("assets", "icons", fmt.Sprintf("%04d.png", index)))}
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if err := validateConfig(config); err != nil {
			b.Fatal(err)
		}
	}
}

func TestHCLCorePathHelpersCoverResolvedPaths(t *testing.T) {
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
	require.NoError(t, Eject(validRoot, "", "3.0.0", false))
	assert.Error(t, Eject(validRoot, "", "3.0.0", false))
	require.NoError(t, Eject(validRoot, "", "3.0.1", true))
}
