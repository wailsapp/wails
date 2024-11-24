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
	if _, err := os.Stat(configPath); err != nil {
		return fmt.Errorf("config file not found: %s", configPath)
	}

	// Generate output filename based on format
	outputFile := fmt.Sprintf("package.%s", options.Format)

	// Create the package
	err = packager.CreatePackageFromConfig(pkgType, configPath, outputFile)
	if err != nil {
		return fmt.Errorf("error creating package: %w", err)
	}

	return nil
}
