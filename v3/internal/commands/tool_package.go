package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/packager"
)

// ToolPackage generates a Linux package in the specified format using nfpm
func ToolPackage(options *flags.ToolPackage) error {
	DisableFooter = true

	if options.ConfigPath == "" {
		return fmt.Errorf("please provide a config file using the -config flag")
	}

	if options.ExecutableName == "" {
		return fmt.Errorf("please provide an executable name using the -name flag")
	}

	// Validate format
	var pkgType packager.PackageType
	switch strings.ToLower(options.Format) {
	case "deb":
		pkgType = packager.DEB
	case "rpm":
		pkgType = packager.RPM
	case "archlinux":
		pkgType = packager.ARCH
	default:
		return fmt.Errorf("unsupported package format '%s'. Supported formats: deb, rpm, archlinux", options.Format)
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
