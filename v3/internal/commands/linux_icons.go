package commands

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/image/draw"
)

const (
	defaultLinuxIconSizes = "16,32,48,64,128,256,512"
	maxLinuxIconSize      = 4096
)

func parseLinuxIconSizes(value string) ([]int, error) {
	if strings.TrimSpace(value) == "" {
		value = defaultLinuxIconSizes
	}

	seen := make(map[int]struct{})
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, fmt.Errorf("invalid Linux icon sizes %q: sizes must be positive integers", value)
		}

		size, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid Linux icon size %q: %w", part, err)
		}
		if size <= 0 || size > maxLinuxIconSize {
			return nil, fmt.Errorf("invalid Linux icon size %d: must be between 1 and %d", size, maxLinuxIconSize)
		}
		seen[size] = struct{}{}
	}

	sizes := make([]int, 0, len(seen))
	for size := range seen {
		sizes = append(sizes, size)
	}
	sort.Ints(sizes)
	return sizes, nil
}

func generateLinuxIcons(iconData []byte, sizes []int, outputDir string) error {
	if err := validateLinuxIconOutputDir(outputDir); err != nil {
		return err
	}
	outputDir = filepath.Clean(outputDir)
	if len(sizes) == 0 {
		return fmt.Errorf("at least one Linux icon size is required")
	}
	for _, size := range sizes {
		if size <= 0 || size > maxLinuxIconSize {
			return fmt.Errorf("invalid Linux icon size %d: must be between 1 and %d", size, maxLinuxIconSize)
		}
	}

	parent := filepath.Dir(outputDir)
	base := filepath.Base(outputDir)
	// Keep the previous complete set in a predictable sibling so the next run
	// can restore it if the process stops between the two directory renames.
	backupDir := filepath.Join(parent, "."+base+".wails-backup")
	if err := recoverLinuxIconOutput(outputDir, backupDir); err != nil {
		return err
	}
	if err := validateExistingLinuxIconOutput(outputDir); err != nil {
		return err
	}

	source, _, err := image.Decode(bytes.NewReader(iconData))
	if err != nil {
		return fmt.Errorf("decode Linux icon source: %w", err)
	}
	if source.Bounds().Dx() <= 0 || source.Bounds().Dy() <= 0 {
		return fmt.Errorf("decode Linux icon source: image has empty bounds")
	}

	if err := os.MkdirAll(parent, 0o755); err != nil {
		return fmt.Errorf("create Linux icon parent directory: %w", err)
	}

	// A sibling staging directory keeps publication on one filesystem. Its
	// unique name also prevents concurrent invocations from deleting each
	// other's in-progress encodes.
	stagingDir, err := os.MkdirTemp(parent, "."+base+".wails-staging-")
	if err != nil {
		return fmt.Errorf("create Linux icon staging directory: %w", err)
	}
	defer os.RemoveAll(stagingDir)
	if err := os.Chmod(stagingDir, 0o755); err != nil {
		return fmt.Errorf("set Linux icon staging directory permissions: %w", err)
	}

	for _, size := range sizes {
		filename := filepath.Join(stagingDir, linuxIconFilename(size))
		if err := encodeLinuxIcon(filename, resizeAndPadIcon(source, size)); err != nil {
			return err
		}
	}

	if err := publishLinuxIconOutput(outputDir, stagingDir, backupDir); err != nil {
		return err
	}
	return nil
}

func validateLinuxIconOutputDir(outputDir string) error {
	if strings.TrimSpace(outputDir) == "" {
		return fmt.Errorf("Linux icon output directory is required")
	}

	clean := filepath.Clean(outputDir)
	volumeRoot := filepath.VolumeName(clean) + string(filepath.Separator)
	if clean == "." || clean == ".." || clean == string(filepath.Separator) || clean == volumeRoot || filepath.Base(clean) == "." || filepath.Base(clean) == ".." {
		return fmt.Errorf("Linux icon output directory must be a dedicated subdirectory")
	}
	return nil
}

func recoverLinuxIconOutput(outputDir, backupDir string) error {
	outputExists, err := pathExists(outputDir)
	if err != nil {
		return fmt.Errorf("inspect Linux icon output directory: %w", err)
	}
	backupExists, err := pathExists(backupDir)
	if err != nil {
		return fmt.Errorf("inspect Linux icon backup directory: %w", err)
	}

	if backupExists {
		if err := validateExistingLinuxIconOutput(backupDir); err != nil {
			return fmt.Errorf("inspect Linux icon backup directory: %w", err)
		}
		if outputExists {
			if err := validateExistingLinuxIconOutput(outputDir); err != nil {
				return err
			}
			if err := os.RemoveAll(backupDir); err != nil {
				return fmt.Errorf("remove committed Linux icon backup: %w", err)
			}
		} else if err := os.Rename(backupDir, outputDir); err != nil {
			return fmt.Errorf("restore Linux icon backup: %w", err)
		}
	}
	return nil
}

func validateExistingLinuxIconOutput(outputDir string) error {
	info, err := os.Lstat(outputDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect Linux icon output directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("Linux icon output path %q is not a directory", outputDir)
	}

	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return fmt.Errorf("inspect Linux icon output directory: %w", err)
	}
	for _, entry := range entries {
		if !entry.Type().IsRegular() || !isLinuxIconFilename(entry.Name()) {
			return fmt.Errorf("refusing to replace Linux icon output directory %q: unexpected entry %q", outputDir, entry.Name())
		}
	}
	return nil
}

func isLinuxIconFilename(filename string) bool {
	stem, found := strings.CutSuffix(filename, ".png")
	if !found {
		return false
	}
	width, height, found := strings.Cut(stem, "x")
	if !found || width != height {
		return false
	}
	size, err := strconv.Atoi(width)
	return err == nil && size > 0 && size <= maxLinuxIconSize && strconv.Itoa(size) == width
}

func publishLinuxIconOutput(outputDir, stagingDir, backupDir string) error {
	// Portable filesystems do not provide an atomic exchange for non-empty
	// directories. The two renames may briefly leave outputDir absent, but both
	// the staging and backup directories always contain complete sets, and the
	// next invocation restores the backup after an interrupted replacement.
	outputExists, err := pathExists(outputDir)
	if err != nil {
		return fmt.Errorf("inspect Linux icon output directory: %w", err)
	}
	if !outputExists {
		if err := os.Rename(stagingDir, outputDir); err != nil {
			return fmt.Errorf("publish Linux icons: %w", err)
		}
		return nil
	}

	if err := validateExistingLinuxIconOutput(outputDir); err != nil {
		return err
	}

	if err := os.Rename(outputDir, backupDir); err != nil {
		return fmt.Errorf("prepare Linux icon output replacement: %w", err)
	}
	if err := os.Rename(stagingDir, outputDir); err != nil {
		if rollbackErr := os.Rename(backupDir, outputDir); rollbackErr != nil {
			return fmt.Errorf("publish Linux icons: %w (restore previous icons: %v)", err, rollbackErr)
		}
		return fmt.Errorf("publish Linux icons: %w", err)
	}
	if err := os.RemoveAll(backupDir); err != nil {
		return fmt.Errorf("remove previous Linux icons: %w", err)
	}
	return nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func resizeAndPadIcon(source image.Image, size int) *image.NRGBA {
	sourceBounds := source.Bounds()
	sourceWidth := sourceBounds.Dx()
	sourceHeight := sourceBounds.Dy()

	targetWidth := size
	targetHeight := size
	if sourceWidth > sourceHeight {
		targetHeight = scaledIconDimension(sourceHeight, sourceWidth, size)
	} else if sourceHeight > sourceWidth {
		targetWidth = scaledIconDimension(sourceWidth, sourceHeight, size)
	}

	x := (size - targetWidth) / 2
	y := (size - targetHeight) / 2
	target := image.NewNRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(
		target,
		image.Rect(x, y, x+targetWidth, y+targetHeight),
		source,
		sourceBounds,
		draw.Src,
		nil,
	)
	return target
}

func scaledIconDimension(sourceDimension, sourceMaximum, targetMaximum int) int {
	scaled := (int64(sourceDimension)*int64(targetMaximum) + int64(sourceMaximum)/2) / int64(sourceMaximum)
	if scaled < 1 {
		return 1
	}
	return int(scaled)
}

func encodeLinuxIcon(filename string, icon image.Image) (err error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("create Linux icon %q: %w", filename, err)
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close Linux icon %q: %w", filename, closeErr)
		}
	}()

	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	if err := encoder.Encode(file, icon); err != nil {
		return fmt.Errorf("encode Linux icon %q: %w", filename, err)
	}
	return nil
}

func linuxIconFilename(size int) string {
	return fmt.Sprintf("%dx%d.png", size, size)
}
