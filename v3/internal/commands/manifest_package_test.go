package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestDarwinAppUsesExactCustomInfoPlist(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".wails", "assets", "darwin"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "bin"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "bin", "example"), []byte("binary"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".wails", "assets", "darwin", "Info.plist"), []byte("com.example.app|darwin|app"), 0o644))

	spec := packageTestSpec("darwin", "arm64", "app")
	spec.Binary = "bin/example"
	spec.Assets = ".wails/assets"
	spec.Output = "dist/Example.app"
	_, err := (&manifestHandler{root: root}).packageApp(spec)
	require.NoError(t, err)
	data, err := os.ReadFile(filepath.Join(root, "dist", "Example.app", "Contents", "Info.plist"))
	require.NoError(t, err)
	assert.Equal(t, "com.example.app|darwin|app", string(data))
}

func TestDMGReplacementProducesResolvedToolOptions(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "templates"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "resources"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "resources", "background.png"), []byte("png"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "README.md"), []byte("read me"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "templates", "dmg.json.tmpl"), []byte(`{"background":"resources/background.png","files":"Read Me=README.md","window_width":600}`), 0o644))

	spec := packageTestSpec("darwin", "arm64", "dmg")
	spec.Output = "dist/example.dmg"
	spec.Config = manifest.PackageFormat{Template: "templates/dmg.json.tmpl"}
	options, err := (&manifestHandler{root: root}).dmgOptions(spec)
	require.NoError(t, err)
	resolvedRoot, err := filepath.EvalSymlinks(root)
	require.NoError(t, err)
	assert.Equal(t, 600, options.DmgWindowWidth)
	workspace := filepath.Join(resolvedRoot, ".wails", "build", "default", "darwin-arm64", "package", "dmg")
	assert.Equal(t, filepath.Join(workspace, "resources", "background.png"), options.BackgroundImage)
	assert.Equal(t, "Read Me="+filepath.Join(workspace, "resources", "file-000.md"), options.DmgFiles)
	require.NoError(t, os.WriteFile(filepath.Join(root, "resources", "background.png"), []byte("mutated"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "README.md"), []byte("mutated"), 0o644))
	assert.Equal(t, "png", string(readTestFile(t, options.BackgroundImage)))
	assert.Equal(t, "read me", string(readTestFile(t, filepath.Join(workspace, "resources", "file-000.md"))))
}

func TestResolveDMGFilesValidatesEveryEntryAndPath(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "README.md"), []byte("read me"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "LICENSE"), []byte("license"), 0o644))

	resolvedRoot, err := filepath.EvalSymlinks(root)
	require.NoError(t, err)
	resolved, err := resolveDMGFiles(root, "Read Me=README.md,License = LICENSE")
	require.NoError(t, err)
	assert.Equal(t, "Read Me="+filepath.Join(resolvedRoot, "README.md")+",License="+filepath.Join(resolvedRoot, "LICENSE"), resolved)

	for _, value := range []string{"entry", "entry=", "entry=../README.md", "entry=.wails/README.md", "entry=missing"} {
		_, err := resolveDMGFiles(root, value)
		require.Error(t, err, value)
	}
	resolved, err = resolveDMGFiles(root, "")
	require.NoError(t, err)
	assert.Empty(t, resolved)
}

func TestDMGOptionsRejectsInvalidReplacementAndResourcePaths(t *testing.T) {
	root := t.TempDir()
	handler := &manifestHandler{root: root}
	spec := packageTestSpec("darwin", "arm64", "dmg")
	spec.Output = "dist/example.dmg"

	spec.Config.Template = "missing.json"
	_, err := handler.dmgOptions(spec)
	require.Error(t, err)

	require.NoError(t, os.WriteFile(filepath.Join(root, "invalid.json"), []byte("{"), 0o644))
	spec.Config.Template = "invalid.json"
	_, err = handler.dmgOptions(spec)
	require.ErrorContains(t, err, "parse rendered DMG template")

	for _, field := range []string{"background", "volume_icon", "file_icon"} {
		spec.Config = manifest.PackageFormat{}
		switch field {
		case "background":
			spec.Config.Background = "missing.png"
		case "volume_icon":
			spec.Config.VolumeIcon = "missing.png"
		case "file_icon":
			spec.Config.FileIcon = "missing.png"
		}
		_, err = handler.dmgOptions(spec)
		require.ErrorContains(t, err, "DMG "+field)
	}

	spec.Config = manifest.PackageFormat{Files: map[string]string{"Read Me": "missing.md"}}
	_, err = handler.dmgOptions(spec)
	require.ErrorContains(t, err, "DMG files path")
}

func TestDMGOptionsReturnsReplacementReadFailure(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "dmg.json"), []byte(`{}`), 0o644))
	spec := packageTestSpec("darwin", "arm64", "dmg")
	spec.Config.Template = "dmg.json"
	want := assert.AnError
	_, err := (&manifestHandler{root: root}).dmgOptionsWithRead(spec, func(string) ([]byte, error) {
		return nil, want
	})
	require.ErrorIs(t, err, want)
}

func TestDMGWorkspaceReplacementIsAtomicOnPreparationFailure(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "dmg.json"), []byte(`{"background":"missing.png"}`), 0o644))
	spec := packageTestSpec("darwin", "arm64", "dmg")
	spec.Config.Template = "dmg.json"
	handler := &manifestHandler{root: root}
	workspace := handler.packageWorkspace(spec)
	require.NoError(t, os.MkdirAll(workspace, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(workspace, "last-complete"), []byte("preserve"), 0o644))

	_, err := handler.dmgOptions(spec)
	require.ErrorContains(t, err, "DMG background")
	assert.Equal(t, "preserve", string(readTestFile(t, filepath.Join(workspace, "last-complete"))))
	staging, globErr := filepath.Glob(filepath.Join(filepath.Dir(workspace), ".dmg-stage-*"))
	require.NoError(t, globErr)
	assert.Empty(t, staging)
}

func TestLinuxPackageReplacementCopiesIntoOwnedWorkspace(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "templates"), 0o755))
	contents := "name: example\nsource: /project/bin/example\n"
	require.NoError(t, os.WriteFile(filepath.Join(root, "templates", "nfpm.yaml.tmpl"), []byte(contents), 0o644))
	spec := packageTestSpec("linux", "amd64", "deb")
	spec.Binary = "bin/example"
	spec.Config.Template = "templates/nfpm.yaml.tmpl"

	path, err := (&manifestHandler{root: root}).prepareLinuxPackageConfig(spec)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, ".wails", "build", "default", "linux-amd64", "package", "deb", "nfpm.yaml"), path)
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, contents, string(data))
}

func TestPackageReplacementRevalidatesSymlinksAtExecution(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "templates"), 0o755))
	source := filepath.Join(root, "templates", "package.conf")
	require.NoError(t, os.WriteFile(source, []byte("safe"), 0o644))
	spec := packageTestSpec("linux", "amd64", "deb")
	spec.Config.Template = "templates/package.conf"

	outside := filepath.Join(t.TempDir(), "outside.conf")
	require.NoError(t, os.WriteFile(outside, []byte("outside"), 0o644))
	require.NoError(t, os.Remove(source))
	require.NoError(t, os.Symlink(outside, source))
	destination := filepath.Join(root, ".wails", "workspace", "package.conf")
	require.NoError(t, os.MkdirAll(filepath.Dir(destination), 0o755))
	require.NoError(t, os.WriteFile(destination, []byte("previous"), 0o644))

	err := (&manifestHandler{root: root}).copyPackageReplacement(spec, destination)
	require.ErrorContains(t, err, "resolves outside the project")
	assert.Equal(t, "previous", readTestFile(t, destination))
}

func TestPackageWorkspacesDoNotMutateGeneratedPlatformAssets(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, ".wails", "assets")
	require.NoError(t, os.MkdirAll(filepath.Join(assets, "windows", "nsis"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(assets, "appicon.png"), []byte("icon"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(assets, "windows", "nsis", "project.nsi"), []byte("generated-default"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "templates"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "templates", "project.nsi.tmpl"), []byte("custom amd64"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "templates", "app.desktop.tmpl"), []byte("Name=Example"), 0o644))

	nsis := packageTestSpec("windows", "amd64", "nsis")
	nsis.Assets = ".wails/assets"
	nsis.Config.Template = "templates/project.nsi.tmpl"
	dir, _, err := (&manifestHandler{root: root}).prepareNSISWorkspace(nsis)
	require.NoError(t, err)
	rendered, err := os.ReadFile(filepath.Join(dir, "project.nsi"))
	require.NoError(t, err)
	assert.Equal(t, "custom amd64", string(rendered))

	appImage := packageTestSpec("linux", "amd64", "appimage")
	appImage.Assets = ".wails/assets"
	appImage.Config.DesktopEntry = "templates/app.desktop.tmpl"
	desktop, icon, _, err := (&manifestHandler{root: root}).prepareAppImageInputs(appImage)
	require.NoError(t, err)
	assert.FileExists(t, desktop)
	assert.Equal(t, filepath.Join(root, ".wails", "build", "default", "linux-amd64", "package", "appimage", "example.png"), icon)
	assert.FileExists(t, icon)

	original, err := os.ReadFile(filepath.Join(assets, "windows", "nsis", "project.nsi"))
	require.NoError(t, err)
	assert.Equal(t, "generated-default", string(original))
	assert.NoFileExists(t, filepath.Join(assets, "example.desktop"))
}

func TestGeneratedPackageWorkspacePreparationPreservesLastCompleteState(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		prepare func(*manifestHandler, pipeline.PackageSpec) error
		spec    pipeline.PackageSpec
	}{
		{
			name: "appimage", pattern: ".appimage-stage-*",
			spec: func() pipeline.PackageSpec {
				spec := packageTestSpec("linux", "amd64", "appimage")
				spec.Assets = ".wails/assets"
				spec.Config.DesktopEntry = "missing.desktop"
				return spec
			}(),
			prepare: func(handler *manifestHandler, spec pipeline.PackageSpec) error {
				_, _, _, err := handler.prepareAppImageInputs(spec)
				return err
			},
		},
		{
			name: "linux", pattern: ".linux-package-stage-*",
			spec: func() pipeline.PackageSpec {
				spec := packageTestSpec("linux", "amd64", "deb")
				spec.Config.Template = "missing.yaml"
				return spec
			}(),
			prepare: func(handler *manifestHandler, spec pipeline.PackageSpec) error {
				_, err := handler.prepareLinuxPackageConfig(spec)
				return err
			},
		},
		{
			name: "nsis", pattern: ".nsis-stage-*",
			spec: func() pipeline.PackageSpec {
				spec := packageTestSpec("windows", "amd64", "nsis")
				spec.Assets = ".wails/assets"
				spec.Config.Template = "missing.nsi"
				return spec
			}(),
			prepare: func(handler *manifestHandler, spec pipeline.PackageSpec) error {
				_, _, err := handler.prepareNSISWorkspace(spec)
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			handler := &manifestHandler{root: root}
			workspace := handler.packageWorkspace(test.spec)
			require.NoError(t, os.MkdirAll(workspace, 0o755))
			require.NoError(t, os.WriteFile(filepath.Join(workspace, "last-complete"), []byte("preserve"), 0o644))
			if test.name == "appimage" {
				require.NoError(t, os.MkdirAll(filepath.Join(root, ".wails", "assets"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(root, ".wails", "assets", "appicon.png"), []byte("icon"), 0o644))
			}
			if test.name == "nsis" {
				require.NoError(t, os.MkdirAll(filepath.Join(root, ".wails", "assets", "windows", "nsis"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(root, ".wails", "assets", "windows", "nsis", "project.nsi"), []byte("generated"), 0o644))
			}

			err := test.prepare(handler, test.spec)
			require.Error(t, err)
			assert.Equal(t, "preserve", string(readTestFile(t, filepath.Join(workspace, "last-complete"))))
			staging, globErr := filepath.Glob(filepath.Join(filepath.Dir(workspace), test.pattern))
			require.NoError(t, globErr)
			assert.Empty(t, staging)
		})
	}
}

func TestMSIXWorkspacePreparationPreservesLastCompleteState(t *testing.T) {
	root := t.TempDir()
	spec := packageTestSpec("windows", "amd64", "msix")
	spec.Config.Manifest = "missing.xml"
	handler := &manifestHandler{root: root}
	workspace := handler.packageWorkspace(spec)
	require.NoError(t, os.MkdirAll(workspace, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(workspace, "last-complete"), []byte("preserve"), 0o644))
	previousHost := manifestHostOS
	manifestHostOS = "windows"
	t.Cleanup(func() { manifestHostOS = previousHost })

	_, err := handler.packageMSIX(spec)
	require.ErrorContains(t, err, "package msix manifest")
	assert.Equal(t, "preserve", string(readTestFile(t, filepath.Join(workspace, "last-complete"))))
	staging, globErr := filepath.Glob(filepath.Join(filepath.Dir(workspace), ".msix-stage-*"))
	require.NoError(t, globErr)
	assert.Empty(t, staging)
}

func TestAndroidPackageFailurePreservesLastCompleteWorkspaceAndArtifact(t *testing.T) {
	root := t.TempDir()
	spec := packageTestSpec("android", "arm64", "aab")
	spec.Assets = ".wails/assets"
	spec.Binary = ".wails/binaries/libwails.so"
	spec.Output = ".wails/artifacts/example.aab"
	handler := &manifestHandler{root: root}
	assets := filepath.Join(root, ".wails", "assets", "android")
	require.NoError(t, os.MkdirAll(assets, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(assets, "gradlew"), []byte("#!/bin/sh\nexit 17\n"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".wails", "binaries"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, spec.Binary), []byte("binary"), 0o755))
	workspace := handler.packageWorkspace(spec)
	require.NoError(t, os.MkdirAll(workspace, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(workspace, "last-complete"), []byte("preserve"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(root, spec.Output)), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, spec.Output), []byte("previous artifact"), 0o644))

	_, err := handler.packageAndroid(t.Context(), spec)
	require.Error(t, err)
	assert.Equal(t, "preserve", string(readTestFile(t, filepath.Join(workspace, "last-complete"))))
	assert.Equal(t, "previous artifact", string(readTestFile(t, filepath.Join(root, spec.Output))))
	staging, globErr := filepath.Glob(filepath.Join(filepath.Dir(workspace), ".android-package-stage-*"))
	require.NoError(t, globErr)
	assert.Empty(t, staging)
}

func TestMSIXStructureUsesCustomAppxManifest(t *testing.T) {
	root := t.TempDir()
	executable := filepath.Join(root, "example.exe")
	custom := filepath.Join(root, "AppxManifest.xml")
	require.NoError(t, os.WriteFile(executable, []byte("binary"), 0o755))
	require.NoError(t, os.WriteFile(custom, []byte("<Package Custom=\"true\"/>"), 0o644))
	out := filepath.Join(root, "structure")

	err := createMSIXPackageStructure(&MSIXOptions{ExecutablePath: executable, AppxManifest: custom}, out)
	require.NoError(t, err)
	data, err := os.ReadFile(filepath.Join(out, "AppxManifest.xml"))
	require.NoError(t, err)
	assert.Equal(t, "<Package Custom=\"true\"/>", string(data))
}

func TestDefaultMSIXManifestIncludesProtocols(t *testing.T) {
	path := filepath.Join(t.TempDir(), "AppxManifest.xml")
	options := &MSIXOptions{}
	options.Protocols = append(options.Protocols, struct {
		Scheme      string `json:"scheme"`
		Description string `json:"description"`
	}{Scheme: "example", Description: "Example links"})

	require.NoError(t, generateAppxManifest(options, path))
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), `<uap:Protocol Name="example">`)
	assert.Contains(t, string(data), `<uap:DisplayName>Example links</uap:DisplayName>`)
}

func packageTestSpec(platform, arch, format string) pipeline.PackageSpec {
	return pipeline.PackageSpec{
		TargetOS: platform, TargetArch: arch, Format: format,
		Binary: "bin/example", Assets: ".wails/assets", Output: "dist/example." + format,
		Profile: "default",
		Project: manifest.Project{Name: "example", ProductName: "Example", BinaryName: "example", Identifier: "com.example.app", Version: "1.0.0"},
	}
}
