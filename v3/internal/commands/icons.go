package commands

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/jackmordaunt/icns/v2"
	"github.com/leaanthony/winicon"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"howett.net/plist"
)

type IconsOptions struct {
	Example           bool   `description:"Generate example icon file (appicon.png) in the current directory"`
	Input             string `description:"The input image file"`
	Sizes             string `description:"The sizes to generate in .ico file (comma separated)" default:"256,128,64,48,32,16"`
	WindowsFilename   string `description:"The output filename for the Windows icon"`
	MacFilename       string `description:"The output filename for the Mac icon bundle"`
	IconComposerInput string `description:"The input Icon Composer file (.icon)"`
	MacAssetDir       string `description:"The output directory for the Mac assets (Assets.car and icons.icns)"`
}

func GenerateIcons(options *IconsOptions) error {
	DisableFooter = true

	if options.Example {
		return generateExampleIcon()
	}

	if options.Input == "" && options.IconComposerInput == "" {
		return fmt.Errorf("either input or icon composer input is required")
	}

	if options.Input != "" && options.WindowsFilename == "" && options.MacFilename == "" {
		return fmt.Errorf("either windows filename or mac filename is required")
	}

	if options.IconComposerInput != "" && options.MacAssetDir == "" {
		return fmt.Errorf("mac asset directory is required")
	}

	// Parse sizes
	var sizes = []int{256, 128, 64, 48, 32, 16}
	var err error
	if options.Sizes != "" {
		sizes, err = parseSizes(options.Sizes)
		if err != nil {
			return err
		}
	}

	// Generate Icons from Icon Composer input
	macIconsGenerated := false
	if options.IconComposerInput != "" {
		if options.MacAssetDir != "" {
			err := generateMacAsset(options)
			if err != nil {
				//Ignore error if the error is "mac asset generation is only supported on macOS" to allow for non-macOS systems to build
				if err.Error() != "mac asset generation is only supported on macOS" {
					return err
				}
			}
			macIconsGenerated = true
		}
	}

	// Generate Icons from input image
	if options.Input != "" {
		iconData, err := os.ReadFile(options.Input)
		if err != nil {
			return err
		}

		if options.WindowsFilename != "" {
			err := generateWindowsIcon(iconData, sizes, options)
			if err != nil {
				return err
			}
		}

		// Generate Icons from input image if no Mac icons were generated from Icon Composer input
		if options.MacFilename != "" && !macIconsGenerated {
			err := generateMacIcon(iconData, options)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func generateExampleIcon() error {
	appIcon, err := buildAssets.ReadFile("build_assets/appicon.png")
	if err != nil {
		return err
	}
	return os.WriteFile("appicon.png", appIcon, 0644)
}

func parseSizes(sizes string) ([]int, error) {
	// split the input string by comma and confirm that each one is an integer
	parsedSizes := strings.Split(sizes, ",")
	var result []int
	for _, size := range parsedSizes {
		s, err := strconv.Atoi(size)
		if err != nil {
			return nil, err
		}
		if s == 0 {
			continue
		}
		result = append(result, s)
	}

	// put all integers in a slice and return
	return result, nil
}

func generateMacIcon(iconData []byte, options *IconsOptions) error {

	srcImg, _, err := image.Decode(bytes.NewBuffer(iconData))
	if err != nil {
		return err
	}

	dest, err := os.Create(options.MacFilename)
	if err != nil {
		return err

	}
	defer func() {
		err = dest.Close()
		if err == nil {
			return
		}
	}()
	return icns.Encode(dest, srcImg)
}

func generateMacAsset(options *IconsOptions) error {
	//Check if running on darwin (macOS), because this will only run on a mac
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("mac asset generation is only supported on macOS")
	}
	// Get system info, because this will only run on macOS 26 or later
	info, err := operatingsystem.Info()
	if err != nil {
		return fmt.Errorf("failed to get system information: %w", err)
	}
	majorStr, _, found := strings.Cut(info.Version, ".")
	if !found {
		return fmt.Errorf("failed to get major version from system information")
	}
	major, err := strconv.Atoi(majorStr)
	if err != nil {
		return fmt.Errorf("failed to convert major version to integer: %w", err)
	}
	if major < 26 {
		return fmt.Errorf("mac asset generation is only supported on macOS 26 or later")
	}

	cmd := exec.Command("/usr/bin/actool", "--version")
	versionPlist, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("failed to get actool version: %w", err)
	}

	// Parse the plist to extract short-bundle-version
	var plistData map[string]any
	if _, err := plist.Unmarshal(versionPlist, &plistData); err != nil {
		return fmt.Errorf("failed to parse actool version plist: %w", err)
	}

	// Navigate to com.apple.actool.version -> short-bundle-version
	actoolVersion, ok := plistData["com.apple.actool.version"].(map[string]any)
	if !ok {
		return fmt.Errorf("failed to find com.apple.actool.version in plist")
	}

	shortVersion, ok := actoolVersion["short-bundle-version"].(string)
	if !ok {
		return fmt.Errorf("failed to find short-bundle-version in plist")
	}

	// Parse the major version number (e.g., "26.2" -> 26)
	actoolMajorStr, _, _ := strings.Cut(shortVersion, ".")
	actoolMajor, err := strconv.Atoi(actoolMajorStr)
	if err != nil {
		return fmt.Errorf("failed to parse major version from short-bundle-version %q: %w", shortVersion, err)
	}

	if actoolMajor < 26 {
		return fmt.Errorf("actool version %s is not supported, version 26 or later is required", shortVersion)
	}

	// Convert paths to absolute paths (required for actool)
	iconComposerPath, err := filepath.Abs(options.IconComposerInput)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for icon composer input: %w", err)
	}
	macAssetDirPath, err := filepath.Abs(options.MacAssetDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for mac asset directory: %w", err)
	}

	// Get Filename from Icon Composer input without extension
	iconComposerFilename := filepath.Base(iconComposerPath)
	iconComposerFilename = strings.TrimSuffix(iconComposerFilename, filepath.Ext(iconComposerFilename))

	cmd = exec.Command("/usr/bin/actool", iconComposerPath,
		"--compile", macAssetDirPath,
		"--notices", "--warnings", "--errors",
		"--output-partial-info-plist", filepath.Join(macAssetDirPath, "/temp.plist"),
		"--app-icon", iconComposerFilename,
		"--enable-on-demand-resources", "NO",
		"--development-region", "en",
		"--target-device", "mac",
		"--minimum-deployment-target", "26.0",
		"--platform", "macosx")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run actool: %w", err)
	}

	// Parse the plist output to verify compilation results
	var compilationResults map[string]any
	if _, err := plist.Unmarshal(out, &compilationResults); err != nil {
		return fmt.Errorf("failed to parse actool compilation results: %w", err)
	}

	// Navigate to com.apple.actool.compilation-results -> output-files
	compilationData, ok := compilationResults["com.apple.actool.compilation-results"].(map[string]any)
	if !ok {
		return fmt.Errorf("failed to find com.apple.actool.compilation-results in plist")
	}

	outputFiles, ok := compilationData["output-files"].([]any)
	if !ok {
		return fmt.Errorf("failed to find output-files array in compilation results")
	}

	if len(outputFiles) != 3 {
		return fmt.Errorf("expected 3 output files, got %d", len(outputFiles))
	}

	// Check that we have one .car file and one .plist file
	var carFile, plistFile, icnsFile string
	for _, file := range outputFiles {
		filePath, ok := file.(string)
		if !ok {
			return fmt.Errorf("output file is not a string: %v", file)
		}
		ext := filepath.Ext(filePath)
		switch ext {
		case ".car":
			carFile = filePath
		case ".plist":
			plistFile = filePath
		case ".icns":
			icnsFile = filePath
		default:
			return fmt.Errorf("unexpected output file extension: %s", ext)
		}
	}

	if carFile == "" {
		return fmt.Errorf("no .car file found in output files")
	}
	if plistFile == "" {
		return fmt.Errorf("no .plist file found in output files")
	}
	if icnsFile == "" {
		return fmt.Errorf("no .icns file found in output files")
	}

	// Remove the temporary plist file since compilation was successful
	if err := os.Remove(plistFile); err != nil {
		return fmt.Errorf("failed to remove temporary plist file: %w", err)
	}

	// Rename the .icns file to icons.icns
	if err := os.Rename(icnsFile, filepath.Join(macAssetDirPath, "icons.icns")); err != nil {
		return fmt.Errorf("failed to rename .icns file to icons.icns: %w", err)
	}

	return nil
}

func generateWindowsIcon(iconData []byte, sizes []int, options *IconsOptions) error {

	var output bytes.Buffer

	err := winicon.GenerateIcon(bytes.NewBuffer(iconData), &output, sizes)
	if err != nil {
		return err
	}

	err = os.WriteFile(options.WindowsFilename, output.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

func GenerateTemplateIcon(data []byte, outputFilename string) (err error) {
	// Decode the input file as a PNG
	buffer := bytes.NewBuffer(data)
	var img image.Image
	img, err = png.Decode(buffer)
	if err != nil {
		return fmt.Errorf("failed to decode input file as PNG: %w", err)
	}

	// Create a new image with the same dimensions and RGBA color model
	bounds := img.Bounds()
	iconImg := image.NewRGBA(bounds)

	// Iterate over each pixel of the input image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the alpha of the pixel
			_, _, _, a := img.At(x, y).RGBA()
			iconImg.SetRGBA(x, y, color.RGBA{R: 0, G: 0, B: 0, A: uint8(a)})
		}
	}

	// Create the output file
	var outFile *os.File
	outFile, err = os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		err = outFile.Close()
	}()

	// Encode the template icon image as a PNG and write it to the output file
	if err = png.Encode(outFile, iconImg); err != nil {
		return fmt.Errorf("failed to encode output image as PNG: %w", err)
	}

	return nil
}
