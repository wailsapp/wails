package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goreleaser/nfpm/v2"
	"github.com/goreleaser/nfpm/v2/files"
)

func TestFreshBuildAssetsEnableLinuxHicolorIcons(t *testing.T) {
	buildDir := filepath.Join(t.TempDir(), "build")
	options := &BuildAssetsOptions{
		Dir:                buildDir,
		Name:               "Test App",
		BinaryName:         "test-app",
		ProductName:        "Test App",
		ProductDescription: "Test application",
		ProductVersion:     "1.0.0",
		ProductCompany:     "Test Company",
		Silent:             true,
	}
	if err := GenerateBuildAssets(options); err != nil {
		t.Fatalf("GenerateBuildAssets() error = %v", err)
	}

	taskfile := readTextFile(t, filepath.Join(buildDir, "Taskfile.yml"))
	if !strings.Contains(taskfile, "-linuxoutputdir linux/icons") {
		t.Error("generated Taskfile does not enable Linux icon generation")
	}
	if !strings.Contains(taskfile, "-linuxsizes "+defaultLinuxIconSizes) {
		t.Error("generated Taskfile Linux sizes do not match the command default")
	}

	sizes, err := parseLinuxIconSizes("")
	if err != nil {
		t.Fatal(err)
	}
	for _, size := range sizes {
		output := fmt.Sprintf(`"linux/icons/%dx%d.png"`, size, size)
		if !strings.Contains(taskfile, output) {
			t.Errorf("generated Taskfile does not declare output %s", output)
		}
	}

	nfpm := readTextFile(t, filepath.Join(buildDir, "linux", "nfpm", "nfpm.yaml"))
	assertNfpmHicolorEntries(t, nfpm, "test-app", sizes)
	if strings.Contains(nfpm, `src: "./build/appicon.png"`) {
		t.Error("fresh nfpm config contains the legacy mismatched 128x128 icon mapping")
	}
	iconData, err := os.ReadFile(filepath.Join(buildDir, "appicon.png"))
	if err != nil {
		t.Fatal(err)
	}
	if err := generateLinuxIcons(iconData, sizes, filepath.Join(buildDir, "linux", "icons")); err != nil {
		t.Fatalf("generate package icons: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(filepath.Dir(buildDir), "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(filepath.Dir(buildDir), "bin", "test-app"), []byte("test binary"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(buildDir, "linux", "test-app.desktop"), []byte("[Desktop Entry]\nName=Test App\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(filepath.Dir(buildDir))
	assertPreparedPackageContentsContainHicolorIcons(t, filepath.Join(buildDir, "linux", "nfpm", "nfpm.yaml"), "test-app", sizes)
}

func TestUpdateBuildAssetsKeepsLegacyLinuxIconMappingWithoutTaskSupport(t *testing.T) {
	buildDir := filepath.Join(t.TempDir(), "build")
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(buildDir, "Taskfile.yml"), []byte("tasks:\n  generate:icons:\n    cmds: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	updateBuildAssetsForIconTest(t, buildDir)
	nfpm := readTextFile(t, filepath.Join(buildDir, "linux", "nfpm", "nfpm.yaml"))
	if !strings.Contains(nfpm, `src: "./build/appicon.png"`) {
		t.Error("updated legacy project did not retain its compatible appicon.png mapping")
	}
	if strings.Contains(nfpm, "./build/linux/icons/") {
		t.Error("updated legacy project references Linux icons its Taskfile cannot generate")
	}
}

func TestUpdateBuildAssetsEnablesLinuxIconMappingsForCapableTaskfile(t *testing.T) {
	buildDir := filepath.Join(t.TempDir(), "build")
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		t.Fatal(err)
	}
	sizes, err := parseLinuxIconSizes("")
	if err != nil {
		t.Fatal(err)
	}
	var capableTaskfile strings.Builder
	capableTaskfile.WriteString("tasks:\n  generate:icons:\n    generates:\n")
	for _, size := range sizes {
		fmt.Fprintf(&capableTaskfile, "      - \"linux/icons/%dx%d.png\"\n", size, size)
	}
	fmt.Fprintf(&capableTaskfile, "    cmds:\n      - wails3 generate icons -linuxoutputdir linux/icons -linuxsizes %s\n", defaultLinuxIconSizes)
	if err := os.WriteFile(filepath.Join(buildDir, "Taskfile.yml"), []byte(capableTaskfile.String()), 0o644); err != nil {
		t.Fatal(err)
	}

	updateBuildAssetsForIconTest(t, buildDir)
	nfpm := readTextFile(t, filepath.Join(buildDir, "linux", "nfpm", "nfpm.yaml"))
	assertNfpmHicolorEntries(t, nfpm, "test-app", sizes)
	if strings.Contains(nfpm, `src: "./build/appicon.png"`) {
		t.Error("capable project retained the legacy icon mapping")
	}
}

func TestTaskfileSupportsLinuxHicolorIconsIsConservative(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		if taskfileSupportsLinuxHicolorIcons(t.TempDir()) {
			t.Error("missing Taskfile reported Linux hicolor support")
		}
	})

	t.Run("unreadable as file", func(t *testing.T) {
		buildDir := t.TempDir()
		if err := os.Mkdir(filepath.Join(buildDir, "Taskfile.yml"), 0o755); err != nil {
			t.Fatal(err)
		}
		if taskfileSupportsLinuxHicolorIcons(buildDir) {
			t.Error("unreadable Taskfile reported Linux hicolor support")
		}
	})

	t.Run("requires flag and output", func(t *testing.T) {
		buildDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(buildDir, "Taskfile.yml"), []byte("-linuxoutputdir linux/icons\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if taskfileSupportsLinuxHicolorIcons(buildDir) {
			t.Error("Taskfile without declared Linux output reported support")
		}
	})

	t.Run("requires every default output", func(t *testing.T) {
		buildDir := t.TempDir()
		content := "-linuxoutputdir linux/icons -linuxsizes 512\nlinux/icons/512x512.png\n"
		if err := os.WriteFile(filepath.Join(buildDir, "Taskfile.yml"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		if taskfileSupportsLinuxHicolorIcons(buildDir) {
			t.Error("Taskfile with only one Linux size reported complete hicolor support")
		}
	})
}

func updateBuildAssetsForIconTest(t *testing.T, buildDir string) {
	t.Helper()
	options := &UpdateBuildAssetsOptions{
		Dir:                buildDir,
		Name:               "Test App",
		BinaryName:         "test-app",
		ProductName:        "Test App",
		ProductDescription: "Test application",
		ProductVersion:     "1.0.0",
		ProductCompany:     "Test Company",
		Silent:             true,
	}
	if err := UpdateBuildAssets(options); err != nil {
		t.Fatalf("UpdateBuildAssets() error = %v", err)
	}
}

func assertNfpmHicolorEntries(t *testing.T, content, binaryName string, sizes []int) {
	t.Helper()
	for _, size := range sizes {
		source := fmt.Sprintf(`src: "./build/linux/icons/%dx%d.png"`, size, size)
		destination := fmt.Sprintf(`dst: "/usr/share/icons/hicolor/%dx%d/apps/%s.png"`, size, size, binaryName)
		if !strings.Contains(content, source) {
			t.Errorf("nfpm config missing %s", source)
		}
		if !strings.Contains(content, destination) {
			t.Errorf("nfpm config missing %s", destination)
		}
	}
}

func assertPreparedPackageContentsContainHicolorIcons(t *testing.T, configPath, binaryName string, sizes []int) {
	t.Helper()
	t.Setenv("GOARCH", "amd64")
	t.Setenv("GIT_COMMITTER_NAME", "Wails Test")
	t.Setenv("GIT_COMMITTER_EMAIL", "test@example.com")
	config, err := nfpm.ParseFile(configPath)
	if err != nil {
		t.Fatalf("parse generated nfpm config: %v", err)
	}
	for _, format := range []string{"deb", "rpm", "archlinux"} {
		info, err := config.Get(format)
		if err != nil {
			t.Fatalf("resolve %s package config: %v", format, err)
		}
		prepared, err := files.PrepareForPackager(info.Contents, info.Umask, format, info.DisableGlobbing, info.MTime)
		if err != nil {
			t.Fatalf("prepare %s package contents: %v", format, err)
		}
		for _, size := range sizes {
			destination := fmt.Sprintf("/usr/share/icons/hicolor/%dx%d/apps/%s.png", size, size, binaryName)
			if !prepared.ContainsDestination(destination) {
				t.Errorf("%s package config does not contain %s", format, destination)
			}
		}
	}
}

func readTextFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
