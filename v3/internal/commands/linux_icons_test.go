package commands

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseLinuxIconSizes(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    []int
		wantErr bool
	}{
		{name: "defaults", want: []int{16, 32, 48, 64, 128, 256, 512}},
		{name: "sorts and deduplicates", value: "64, 16,32,16", want: []int{16, 32, 64}},
		{name: "empty element", value: "16,,32", wantErr: true},
		{name: "not a number", value: "16,large", wantErr: true},
		{name: "zero", value: "0", wantErr: true},
		{name: "negative", value: "-16", wantErr: true},
		{name: "too large", value: "4097", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseLinuxIconSizes(test.value)
			if (err != nil) != test.wantErr {
				t.Fatalf("parseLinuxIconSizes(%q) error = %v, wantErr %v", test.value, err, test.wantErr)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("parseLinuxIconSizes(%q) = %v, want %v", test.value, got, test.want)
			}
		})
	}
}

func TestIsLinuxIconFilename(t *testing.T) {
	tests := map[string]bool{
		"16x16.png":     true,
		"4096x4096.png": true,
		"0x0.png":       false,
		"016x016.png":   false,
		"16x32.png":     false,
		"4097x4097.png": false,
		"16x16.jpg":     false,
		"README.txt":    false,
	}
	for filename, want := range tests {
		if got := isLinuxIconFilename(filename); got != want {
			t.Errorf("isLinuxIconFilename(%q) = %v, want %v", filename, got, want)
		}
	}
}

func TestGenerateLinuxIconsDimensionsPaddingAndDeterminism(t *testing.T) {
	iconData := encodeTestPNG(t, 40, 20, color.NRGBA{R: 220, G: 40, B: 10, A: 128})
	outputDir := filepath.Join(t.TempDir(), "icons")
	sizes := []int{16, 32}

	if err := generateLinuxIcons(iconData, sizes, outputDir); err != nil {
		t.Fatalf("generateLinuxIcons() error = %v", err)
	}
	firstRun := readGeneratedIcons(t, outputDir, sizes)

	for _, size := range sizes {
		filename := linuxIconFilename(size)
		generated, err := png.Decode(bytes.NewReader(firstRun[filename]))
		if err != nil {
			t.Fatalf("decode %s: %v", filename, err)
		}
		if got, want := generated.Bounds(), image.Rect(0, 0, size, size); got != want {
			t.Errorf("%s bounds = %v, want %v", filename, got, want)
		}

		nrgba, ok := generated.(*image.NRGBA)
		if !ok {
			t.Fatalf("%s decoded as %T, want *image.NRGBA", filename, generated)
		}
		if alpha := nrgba.NRGBAAt(0, 0).A; alpha != 0 {
			t.Errorf("%s corner alpha = %d, want transparent padding", filename, alpha)
		}
		if alpha := nrgba.NRGBAAt(size/2, size/2).A; alpha == 0 || alpha == 255 {
			t.Errorf("%s center alpha = %d, want preserved source transparency", filename, alpha)
		}
	}

	if err := generateLinuxIcons(iconData, sizes, outputDir); err != nil {
		t.Fatalf("second generateLinuxIcons() error = %v", err)
	}
	secondRun := readGeneratedIcons(t, outputDir, sizes)
	if !reflect.DeepEqual(secondRun, firstRun) {
		t.Error("identical input produced different Linux PNG bytes")
	}
}

func TestGenerateLinuxIconsRecoversInterruptedReplacement(t *testing.T) {
	iconData := encodeTestPNG(t, 20, 40, color.NRGBA{R: 10, G: 80, B: 230, A: 255})
	outputDir := filepath.Join(t.TempDir(), "icons")
	backupDir := filepath.Join(filepath.Dir(outputDir), ".icons.wails-backup")

	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(backupDir, "32x32.png"), []byte("previous complete set"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := generateLinuxIcons(iconData, []int{24, 48}, outputDir); err != nil {
		t.Fatalf("generateLinuxIcons() error = %v", err)
	}
	readGeneratedIcons(t, outputDir, []int{24, 48})
	if _, err := os.Stat(backupDir); !os.IsNotExist(err) {
		t.Errorf("recovery artifact %q still exists; stat error = %v", backupDir, err)
	}
}

func TestGenerateLinuxIconsRestoresBackupBeforeDecoding(t *testing.T) {
	outputDir := filepath.Join(t.TempDir(), "icons")
	backupDir := filepath.Join(filepath.Dir(outputDir), ".icons.wails-backup")
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		t.Fatal(err)
	}
	previous := []byte("previous complete set")
	if err := os.WriteFile(filepath.Join(backupDir, "32x32.png"), previous, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := generateLinuxIcons([]byte("not an image"), []int{32}, outputDir); err == nil {
		t.Fatal("generateLinuxIcons() error = nil, want decode error")
	}
	got, err := os.ReadFile(filepath.Join(outputDir, "32x32.png"))
	if err != nil {
		t.Fatalf("read restored icon: %v", err)
	}
	if !bytes.Equal(got, previous) {
		t.Errorf("restored icon = %q, want %q", got, previous)
	}
	if _, err := os.Stat(backupDir); !os.IsNotExist(err) {
		t.Errorf("backup directory still exists; stat error = %v", err)
	}
}

func TestGenerateLinuxIconsFailureLeavesExistingSetUntouched(t *testing.T) {
	outputDir := filepath.Join(t.TempDir(), "icons")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}
	existingPath := filepath.Join(outputDir, "32x32.png")
	existing := []byte("existing complete output")
	if err := os.WriteFile(existingPath, existing, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := generateLinuxIcons([]byte("not an image"), []int{32}, outputDir); err == nil {
		t.Fatal("generateLinuxIcons() error = nil, want decode error")
	}
	got, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatalf("read existing icon after failure: %v", err)
	}
	if !bytes.Equal(got, existing) {
		t.Errorf("existing icon changed after failed generation: got %q, want %q", got, existing)
	}
}

func TestGenerateLinuxIconsRefusesDirectoryWithUnrelatedFiles(t *testing.T) {
	outputDir := filepath.Join(t.TempDir(), "icons")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}
	unrelatedPath := filepath.Join(outputDir, "README.txt")
	if err := os.WriteFile(unrelatedPath, []byte("keep me"), 0o644); err != nil {
		t.Fatal(err)
	}

	iconData := encodeTestPNG(t, 32, 32, color.NRGBA{R: 50, G: 180, B: 90, A: 255})
	if err := generateLinuxIcons(iconData, []int{16}, outputDir); err == nil {
		t.Fatal("generateLinuxIcons() error = nil, want unsafe replacement error")
	}
	if got, err := os.ReadFile(unrelatedPath); err != nil || string(got) != "keep me" {
		t.Errorf("unrelated file was changed: content = %q, error = %v", got, err)
	}
}

func TestGenerateIconsLinuxOnlyRegeneratesMissingSize(t *testing.T) {
	tempDir := t.TempDir()
	input := filepath.Join(tempDir, "source.png")
	if err := os.WriteFile(input, encodeTestPNG(t, 32, 32, color.NRGBA{R: 50, G: 180, B: 90, A: 255}), 0o644); err != nil {
		t.Fatal(err)
	}
	outputDir := filepath.Join(tempDir, "icons")
	options := &IconsOptions{
		Input:          input,
		LinuxOutputDir: outputDir,
		LinuxSizes:     "16,32",
	}

	if err := GenerateIcons(options); err != nil {
		t.Fatalf("GenerateIcons() error = %v", err)
	}
	missing := filepath.Join(outputDir, "16x16.png")
	if err := os.Remove(missing); err != nil {
		t.Fatal(err)
	}
	if err := GenerateIcons(options); err != nil {
		t.Fatalf("GenerateIcons() after missing output error = %v", err)
	}
	readGeneratedIcons(t, outputDir, []int{16, 32})
}

func TestGenerateIconsRequiresRasterInputForLinux(t *testing.T) {
	options := &IconsOptions{
		IconComposerInput: "appicon.icon",
		MacAssetDir:       t.TempDir(),
		LinuxOutputDir:    filepath.Join(t.TempDir(), "icons"),
	}

	err := GenerateIcons(options)
	if err == nil || !strings.Contains(err.Error(), "input is required for Linux icon generation") {
		t.Fatalf("GenerateIcons() error = %v, want Linux raster input error", err)
	}
}

func TestGenerateLinuxIconsCleansOutputDirectoryPath(t *testing.T) {
	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "parent", "..", "icons") + string(filepath.Separator)
	iconData := encodeTestPNG(t, 32, 32, color.NRGBA{R: 50, G: 180, B: 90, A: 255})

	if err := generateLinuxIcons(iconData, []int{16}, outputDir); err != nil {
		t.Fatalf("generateLinuxIcons() error = %v", err)
	}
	readGeneratedIcons(t, filepath.Join(tempDir, "icons"), []int{16})
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".icons.wails-staging-") {
			t.Errorf("staging directory %q still exists", entry.Name())
		}
	}
}

func encodeTestPNG(t *testing.T, width, height int, fill color.NRGBA) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetNRGBA(x, y, fill)
		}
	}

	var data bytes.Buffer
	if err := png.Encode(&data, img); err != nil {
		t.Fatalf("encode source PNG: %v", err)
	}
	return data.Bytes()
}

func readGeneratedIcons(t *testing.T, outputDir string, sizes []int) map[string][]byte {
	t.Helper()
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("read output directory: %v", err)
	}
	if got, want := len(entries), len(sizes); got != want {
		t.Fatalf("output entry count = %d, want %d", got, want)
	}

	result := make(map[string][]byte, len(sizes))
	for _, size := range sizes {
		filename := linuxIconFilename(size)
		data, err := os.ReadFile(filepath.Join(outputDir, filename))
		if err != nil {
			t.Fatalf("read %s: %v", filename, err)
		}
		result[filename] = data
	}
	return result
}
