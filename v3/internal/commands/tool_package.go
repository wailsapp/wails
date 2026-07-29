package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	// "github.com/wailsapp/wails/v3/internal/commands/dmg" // TODO: Missing package
	"github.com/wailsapp/wails/v3/internal/commands/macpkg"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/packager"
)

// ToolPackage generates a package in the specified format
func ToolPackage(options *flags.ToolPackage) error {
	DisableFooter = true

	// Check if we're creating a DMG or PKG
	isDMG := strings.ToLower(options.Format) == "dmg" || options.CreateDMG
	isPKG := strings.ToLower(options.Format) == "pkg" || options.CreatePKG

	// Config file is required for Linux packages but optional for DMG/PKG
	if options.ConfigPath == "" && !isDMG && !isPKG {
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

		// Create output path for DMG
		dmgPath := filepath.Join(options.Out, fmt.Sprintf("%s.dmg", options.ExecutableName))

		// DMG creation temporarily disabled - missing dmg package
		_ = dmgPath // avoid unused variable warning
		return fmt.Errorf("DMG creation is temporarily disabled due to missing dmg package")
		
		// // Create DMG creator
		// dmgCreator, err := dmg.New(appPath, dmgPath, options.ExecutableName)
		// if err != nil {
		// 	return fmt.Errorf("error creating DMG: %w", err)
		// }

		// // Set background image if provided
		// if options.BackgroundImage != "" {
		// 	if err := dmgCreator.SetBackgroundImage(options.BackgroundImage); err != nil {
		// 		return fmt.Errorf("error setting background image: %w", err)
		// 	}
		// }

		// // Set default icon positions
		// dmgCreator.AddIconPosition(filepath.Base(appPath), 150, 175)
		// dmgCreator.AddIconPosition("Applications", 450, 175)

		// // Create the DMG
		// if err := dmgCreator.Create(); err != nil {
		// 	return fmt.Errorf("error creating DMG: %w", err)
		// }

		// fmt.Printf("DMG created successfully: %s\n", dmgPath)
		// return nil
	}

	// Handle PKG creation for macOS
	if isPKG {
		if runtime.GOOS != "darwin" {
			return fmt.Errorf("PKG creation is only supported on macOS")
		}

		// Convert ToolPackage options to MacPKG options
		macPKGOptions := &macpkg.BuildMacPKGOptions{
			ConfigPath:        options.ConfigPath,
			GenerateTemplate:  options.GenerateTemplate,
			SkipNotarization:  options.SkipNotarization,
			ValidateOnly:      options.ValidateOnly,
			Verbose:           false, // Could be added to flags if needed
		}

		return macpkg.BuildMacPKG(macPKGOptions)
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
		return fmt.Errorf("unsupported package format '%s'. Supported formats: deb, rpm, archlinux, dmg, pkg", options.Format)
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
