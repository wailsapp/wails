package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/leaanthony/dmg/dmg"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/packager"
)

// ToolPackage generates a package in the specified format
func ToolPackage(options *flags.ToolPackage) error {
	DisableFooter = true

	// Check if we're creating a DMG
	isDMG := strings.ToLower(options.Format) == "dmg" || options.CreateDMG

	// Config file is required for Linux packages but optional for DMG
	if options.ConfigPath == "" && !isDMG {
		return fmt.Errorf("please provide a config file using the -config flag")
	}

	if options.ExecutableName == "" {
		return fmt.Errorf("please provide an executable name using the -name flag")
	}

	// Handle DMG creation for macOS
	if isDMG {
		if runtime.GOOS != "darwin" {
			return fmt.Errorf("DMG creation is only supported on macOS")
		}

		// For DMG, we expect the .app bundle to already exist
		appPath := filepath.Join(options.Out, fmt.Sprintf("%s.app", options.ExecutableName))
		if _, err := os.Stat(appPath); os.IsNotExist(err) {
			return fmt.Errorf("application bundle not found: %s", appPath)
		}

		// Create output path for DMG.
		dmgPath := filepath.Join(options.Out, fmt.Sprintf("%s.dmg", options.ExecutableName))

		opts, err := buildDMGOptions(options, appPath, dmgPath)
		if err != nil {
			return err
		}
		if err := dmg.Build(opts); err != nil {
			return fmt.Errorf("error creating DMG: %w", err)
		}
		return nil
	}

	// For Linux packages, continue with existing logic
	var pkgType packager.PackageType
	switch strings.ToLower(options.Format) {
	case "deb":
		pkgType = packager.DEB
	case "rpm":
		pkgType = packager.RPM
	case "archlinux":
		pkgType = packager.ARCH
	default:
		return fmt.Errorf("unsupported package format '%s'. Supported formats: deb, rpm, archlinux, dmg", options.Format)
	}

	// Get absolute path of config file
	configPath, err := filepath.Abs(options.ConfigPath)
	if err != nil {
		return fmt.Errorf("error getting absolute path of config file: %w", err)
	}

	// Check if config file exists
	if info, err := os.Stat(configPath); err != nil {
		return fmt.Errorf("config file not found: %s", configPath)
	} else if info.Mode().Perm()&0444 == 0 {
		return fmt.Errorf("config file is not readable: %s", configPath)
	}

	// Generate output filename based on format
	if options.Format == "archlinux" {
		// Arch linux packages are not .archlinux files, they are .pkg.tar.zst
		options.Format = "pkg.tar.zst"
	}
	outputFile := filepath.Join(options.Out, fmt.Sprintf("%s.%s", options.ExecutableName, options.Format))

	// Create the package
	err = packager.CreatePackageFromConfig(pkgType, configPath, outputFile)
	if err != nil {
		return fmt.Errorf("error creating package: %w", err)
	}

	return nil
}

// buildDMGOptions converts package flags into a tested DMG configuration without
// invoking macOS tooling. The application bundle is always added to the image;
// extra files, if any, are configured with the same name=path syntax exposed by
// the Taskfile.
func buildDMGOptions(options *flags.ToolPackage, appPath, dmgPath string) (dmg.Options, error) {
	opts := dmg.DefaultOptions(appPath, dmgPath)
	opts.VolumeName = options.ExecutableName
	windowWidth, windowHeight := options.DmgWindowWidth, options.DmgWindowHeight
	if windowWidth <= 0 {
		windowWidth = 540
	}
	if windowHeight <= 0 {
		windowHeight = 380
	}
	opts.Window = dmg.WindowConfig{X: 100, Y: 100, Width: windowWidth, Height: windowHeight}
	opts.Icon = dmg.IconConfig{Size: 96, TextSize: 12, GridSpace: 100}
	opts.Files[filepath.Base(appPath)] = appPath
	// Finder renders the stored Y from the top of its content area, below a
	// small title-bar offset. Subtract it so the default row is visibly centred.
	iconY := windowHeight/2 - 26
	if iconY < 0 {
		iconY = windowHeight / 2
	}
	opts.IconPositions = map[string]dmg.IconPosition{
		filepath.Base(appPath): {X: windowWidth * 28 / 100, Y: iconY},
		"Applications":         {X: windowWidth * 72 / 100, Y: iconY},
	}
	if options.BackgroundImage != "" {
		opts.Background = &dmg.BackgroundConfig{File: options.BackgroundImage}
	}
	if options.DmgVolumeIcon != "" {
		opts.VolumeIcon = options.DmgVolumeIcon
	}
	if options.DmgFileIcon != "" {
		opts.FileIcon = options.DmgFileIcon
	}
	if err := addDMGFiles(&opts, options.DmgFiles); err != nil {
		return dmg.Options{}, err
	}
	return opts, nil
}

// addDMGFiles adds optional installer resources configured in a Taskfile. The
// deliberately small name=path syntax keeps this low-level packaging command
// useful without introducing another project configuration format.
func addDMGFiles(opts *dmg.Options, value string) error {
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		name, path, ok := strings.Cut(item, "=")
		if !ok || strings.TrimSpace(name) == "" || strings.TrimSpace(path) == "" {
			return fmt.Errorf("invalid DMG file %q: expected name=path", item)
		}
		name = strings.TrimSpace(name)
		path = strings.TrimSpace(path)
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("DMG file %q: %w", name, err)
		}
		if _, exists := opts.Files[name]; exists {
			return fmt.Errorf("DMG file %q conflicts with an existing entry", name)
		}
		opts.Files[name] = path
	}
	return nil
}
