// Package packager provides a simplified interface for creating Linux packages using nfpm
package packager

import (
	"fmt"
	"io"
	"os"

	"github.com/goreleaser/nfpm/v2"
	_ "github.com/goreleaser/nfpm/v2/apk"  // Register APK packager
	_ "github.com/goreleaser/nfpm/v2/arch" // Register Arch Linux packager
	_ "github.com/goreleaser/nfpm/v2/deb"  // Register DEB packager
	_ "github.com/goreleaser/nfpm/v2/ipk"  // Register IPK packager
	_ "github.com/goreleaser/nfpm/v2/rpm"  // Register RPM packager
)

// PackageType represents supported package formats
type PackageType string

const (
	// DEB is for Debian/Ubuntu packages
	DEB PackageType = "deb"
	// RPM is for RedHat/CentOS packages
	RPM PackageType = "rpm"
	// APK is for Alpine Linux packages
	APK PackageType = "apk"
	// IPK is for OpenWrt packages
	IPK PackageType = "ipk"
	// ARCH is for Arch Linux packages
	ARCH PackageType = "archlinux"
)

// CreatePackageFromConfig loads a configuration file and creates a package
func CreatePackageFromConfig(pkgType PackageType, configPath string, output string) error {
	// Parse nfpm config
	config, err := nfpm.ParseFile(configPath)
	if err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	// Get info for the specified packager
	info, err := config.Get(string(pkgType))
	if err != nil {
		return fmt.Errorf("error getting packager info: %w", err)
	}

	// Get the packager
	packager, err := nfpm.Get(string(pkgType))
	if err != nil {
		return fmt.Errorf("error getting packager: %w", err)
	}

	// Create output file
	out, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer out.Close()

	// Create the package
	if err := packager.Package(info, out); err != nil {
		return fmt.Errorf("error creating package: %w", err)
	}

	return nil
}

// CreatePackageFromConfigWriter loads a configuration file and writes the package to the provided writer
func CreatePackageFromConfigWriter(pkgType PackageType, configPath string, output io.Writer) error {
	// Parse nfpm config
	config, err := nfpm.ParseFile(configPath)
	if err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	// Get info for the specified packager
	info, err := config.Get(string(pkgType))
	if err != nil {
		return fmt.Errorf("error getting packager info: %w", err)
	}

	// Get the packager
	packager, err := nfpm.Get(string(pkgType))
	if err != nil {
		return fmt.Errorf("error getting packager: %w", err)
	}

	// Create the package
	if err := packager.Package(info, output); err != nil {
		return fmt.Errorf("error creating package: %w", err)
	}

	return nil
}
