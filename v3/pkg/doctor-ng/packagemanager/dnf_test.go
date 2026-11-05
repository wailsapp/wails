//go:build linux

package packagemanager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDnfInstallCommandUsesRpmOstreeOnAtomicSystems(t *testing.T) {
	pkg := &Package{Name: "gtk4-devel", SystemPackage: true}

	if got := newDnf("fedora", true).InstallCommand(pkg); got != "sudo rpm-ostree install gtk4-devel" {
		t.Fatalf("atomic install command = %q", got)
	}
	if got := newDnf("fedora", false).InstallCommand(pkg); got != "sudo dnf install gtk4-devel" {
		t.Fatalf("regular install command = %q", got)
	}
}

func TestIsAtomicSystemAt(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, "ostree-booted")

	if isAtomicSystemAt(marker) {
		t.Fatal("missing ostree marker detected as atomic")
	}
	if err := os.WriteFile(marker, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if !isAtomicSystemAt(marker) {
		t.Fatal("ostree marker was not detected")
	}
}
