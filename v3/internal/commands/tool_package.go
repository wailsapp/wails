package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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
	if err := applyDMGIconLayout(&opts, options.DmgIconLayout, options.DmgIconPositions); err != nil {
		return dmg.Options{}, err
	}
	return opts, nil
}

// applyDMGIconLayout either delegates placement to the DMG library or parses
// explicit Finder-window pixel centres supplied as name=x,y;name=x,y.
func applyDMGIconLayout(opts *dmg.Options, layout, positions string) error {
	switch strings.ToLower(strings.TrimSpace(layout)) {
	case "", "auto":
		if strings.TrimSpace(positions) != "" {
			return fmt.Errorf("DMG icon positions require manual icon layout")
		}
		opts.IconPositions = nil
		return nil
	case "manual":
		iconPositions, err := parseDMGIconPositions(positions)
		if err != nil {
			return err
		}
		for name := range iconPositions {
			if name != "Applications" {
				if _, ok := opts.Files[name]; !ok {
					return fmt.Errorf("DMG icon position references unknown item %q", name)
				}
			}
		}
		for name := range opts.Files {
			if _, ok := iconPositions[name]; !ok {
				return fmt.Errorf("manual DMG icon layout is missing a position for %q", name)
			}
		}
		if opts.AddApplicationsSymlink {
			if _, ok := iconPositions["Applications"]; !ok {
				return fmt.Errorf("manual DMG icon layout is missing a position for Applications")
			}
		}
		opts.IconPositions = iconPositions
		return nil
	default:
		return fmt.Errorf("invalid DMG icon layout %q: expected auto or manual", layout)
	}
}

func parseDMGIconPositions(value string) (map[string]dmg.IconPosition, error) {
	if strings.TrimSpace(value) == "" {
		return nil, fmt.Errorf("manual DMG icon layout requires icon positions")
	}

	positions := make(map[string]dmg.IconPosition)
	for _, item := range strings.Split(value, ";") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		name, coordinates, ok := strings.Cut(item, "=")
		if !ok || strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("invalid DMG icon position %q: expected name=x,y", item)
		}
		xy := strings.Split(coordinates, ",")
		if len(xy) != 2 {
			return nil, fmt.Errorf("invalid DMG icon position %q: expected name=x,y", item)
		}
		x, err := strconv.Atoi(strings.TrimSpace(xy[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid DMG icon X coordinate in %q: %w", item, err)
		}
		y, err := strconv.Atoi(strings.TrimSpace(xy[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid DMG icon Y coordinate in %q: %w", item, err)
		}
		positions[strings.TrimSpace(name)] = dmg.IconPosition{X: x, Y: y}
	}
	if len(positions) == 0 {
		return nil, fmt.Errorf("manual DMG icon layout requires icon positions")
	}
	return positions, nil
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
