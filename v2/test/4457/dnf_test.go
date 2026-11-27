package test4457

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
)

func TestDnfPackageInstalledNotInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	mockDnf := filepath.Join(tmpDir, "dnf")
	script := `#!/bin/sh
# Mock dnf that simulates "package not installed" exit code
# dnf -q list --installed returns exit code 1 when package is not installed
exit 1
`
	err := os.WriteFile(mockDnf, []byte(script), 0755)
	if err != nil {
		t.Fatal(err)
	}

	origPath := os.Getenv("PATH")
	os.Setenv("PATH", tmpDir+":"+origPath)
	defer os.Setenv("PATH", origPath)

	dnf := packagemanager.NewDnf("fedora")
	pkg := &packagemanager.Package{Name: "webkit2gtk4.0-devel", SystemPackage: true}

	installed, err := dnf.PackageInstalled(pkg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if installed {
		t.Error("expected package to NOT be installed, but got installed=true")
	}
}

func TestDnfPackageInstalledIsInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	mockDnf := filepath.Join(tmpDir, "dnf")
	script := `#!/bin/sh
# Mock dnf that simulates installed package output
echo "Installed packages"
echo "webkit2gtk4.0-devel.x86_64 2.46.5-1.fc41 updates"
exit 0
`
	err := os.WriteFile(mockDnf, []byte(script), 0755)
	if err != nil {
		t.Fatal(err)
	}

	origPath := os.Getenv("PATH")
	os.Setenv("PATH", tmpDir+":"+origPath)
	defer os.Setenv("PATH", origPath)

	dnf := packagemanager.NewDnf("fedora")
	pkg := &packagemanager.Package{Name: "webkit2gtk4.0-devel", SystemPackage: true}

	installed, err := dnf.PackageInstalled(pkg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !installed {
		t.Error("expected package to be installed, but got installed=false")
	}
}

func TestDnfPackageInstalledOldDnfInfoFallbackFails(t *testing.T) {
	tmpDir := t.TempDir()
	mockDnf := filepath.Join(tmpDir, "dnf")
	script := `#!/bin/sh
# This simulates the OLD behavior: "dnf info installed" returns available
# packages with exit code 0 even when not installed.
# With the fix using "dnf -q list --installed", exit code 1 means not installed.
# We verify the mock is called with the correct arguments.
for arg in "$@"; do
	echo "ARG: $arg"
done
# If called with "info installed" (old behavior), exit 0 with version info (bug)
# If called with "-q list --installed" (new behavior), exit 1 (not installed)
if [ "$1" = "-q" ] && [ "$2" = "list" ] && [ "$3" = "--installed" ]; then
	exit 1
fi
# Old path would reach here and incorrectly report installed
echo "Version : 2.46.5"
exit 0
`
	err := os.WriteFile(mockDnf, []byte(script), 0755)
	if err != nil {
		t.Fatal(err)
	}

	origPath := os.Getenv("PATH")
	os.Setenv("PATH", tmpDir+":"+origPath)
	defer os.Setenv("PATH", origPath)

	dnf := packagemanager.NewDnf("fedora")
	pkg := &packagemanager.Package{Name: "webkit2gtk4.0-devel", SystemPackage: true}

	installed, err := dnf.PackageInstalled(pkg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if installed {
		t.Error("expected package to NOT be installed with new dnf -q list --installed approach")
	}
}
