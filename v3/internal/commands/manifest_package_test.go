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

func TestDarwinAppUsesRenderedCustomInfoPlist(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "templates"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".wails", "assets", "darwin"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "bin"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "bin", "example"), []byte("binary"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "templates", "Info.plist.tmpl"), []byte("{{.Project.Identifier}}|{{.Target.OS}}|{{.Package.Format}}"), 0o644))

	spec := packageTestSpec("darwin", "arm64", "app")
	spec.Binary = "bin/example"
	spec.Assets = ".wails/assets"
	spec.Output = "dist/Example.app"
	spec.Config.Template = "templates/Info.plist.tmpl"
	_, err := (&manifestHandler{root: root}).packageApp(spec)
	require.NoError(t, err)
	data, err := os.ReadFile(filepath.Join(root, "dist", "Example.app", "Contents", "Info.plist"))
	require.NoError(t, err)
	assert.Equal(t, "com.example.app|darwin|app", string(data))
}

func TestDMGTemplateAndStructuredOptionsProduceResolvedToolOptions(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "templates"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "resources"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "resources", "background.png"), []byte("png"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "README.md"), []byte("read me"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "templates", "dmg.json.tmpl"), []byte(`{"background":"resources/background.png","files":"Read Me=README.md","window_width":600}`), 0o644))

	spec := packageTestSpec("darwin", "arm64", "dmg")
	spec.Output = "dist/example.dmg"
	spec.Config = manifest.PackageFormat{Template: "templates/dmg.json.tmpl", Options: map[string]any{"window_width": 720}}
	options, err := (&manifestHandler{root: root}).dmgOptions(spec)
	require.NoError(t, err)
	assert.Equal(t, 720, options.DmgWindowWidth, "structured manifest values override rendered defaults")
	assert.Equal(t, filepath.Join(root, "resources", "background.png"), options.BackgroundImage)
	assert.Equal(t, "Read Me="+filepath.Join(root, "README.md"), options.DmgFiles)
}

func TestLinuxPackageTemplateRendersIntoOwnedWorkspace(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "templates"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "templates", "nfpm.yaml.tmpl"), []byte("name: {{.Project.BinaryName}}\nsource: {{.Paths.Binary}}\n"), 0o644))
	spec := packageTestSpec("linux", "amd64", "deb")
	spec.Binary = "bin/example"
	spec.Config.Template = "templates/nfpm.yaml.tmpl"

	path, err := (&manifestHandler{root: root}).prepareLinuxPackageConfig(spec)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, ".wails", "build", "default", "linux-amd64", "package", "deb", "nfpm.yaml"), path)
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), "source: "+filepath.Join(root, "bin", "example"))
}

func TestPackageWorkspacesDoNotMutateGeneratedPlatformAssets(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, ".wails", "assets")
	require.NoError(t, os.MkdirAll(filepath.Join(assets, "windows", "nsis"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(assets, "appicon.png"), []byte("icon"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(assets, "windows", "nsis", "project.nsi"), []byte("generated-default"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "templates"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "templates", "project.nsi.tmpl"), []byte("custom {{.Target.Arch}}"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "templates", "app.desktop.tmpl"), []byte("Name={{.Project.ProductName}}"), 0o644))

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
	appImage.Config.Template = "templates/app.desktop.tmpl"
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

func TestHookCommandUsesOnlyPlatformScriptInterpreters(t *testing.T) {
	command, args := hookCommand("windows", `C:\project with spaces\before.cmd`)
	assert.Equal(t, "cmd.exe", command)
	assert.Equal(t, []string{"/d", "/s", "/c", `C:\project with spaces\before.cmd`}, args)

	command, args = hookCommand("windows", `C:\project\before.ps1`)
	assert.Equal(t, "powershell.exe", command)
	assert.Equal(t, "-File", args[len(args)-2])
	assert.Equal(t, `C:\project\before.ps1`, args[len(args)-1])

	command, args = hookCommand("linux", "/project/before.sh")
	assert.Equal(t, "/project/before.sh", command)
	assert.Empty(t, args)
}

func packageTestSpec(platform, arch, format string) pipeline.PackageSpec {
	return pipeline.PackageSpec{
		TargetOS: platform, TargetArch: arch, Format: format,
		Binary: "bin/example", Assets: ".wails/assets", Output: "dist/example." + format,
		Profile: "default",
		Project: manifest.Project{Name: "example", ProductName: "Example", BinaryName: "example", Identifier: "com.example.app", Version: "1.0.0"},
	}
}
