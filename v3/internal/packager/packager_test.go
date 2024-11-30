package packager

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCreatePackageFromConfig(t *testing.T) {
	// Create a temporary file for testing
	content := []byte("test content")
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Create a temporary config file
	configContent := []byte(`
name: test-package
version: v1.0.0
arch: amd64
description: Test package
maintainer: Test User <test@example.com>
license: MIT
contents:
- src: ` + tmpfile.Name() + `
  dst: /usr/local/bin/test-file
`)

	configFile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(configFile.Name())

	if _, err := configFile.Write(configContent); err != nil {
		t.Fatal(err)
	}
	if err := configFile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test creating packages for each format
	formats := []struct {
		pkgType PackageType
		ext     string
	}{
		{DEB, "deb"},
		{RPM, "rpm"},
		{APK, "apk"},
		{IPK, "ipk"},
		{ARCH, "pkg.tar.zst"},
	}

	for _, format := range formats {
		t.Run(string(format.pkgType), func(t *testing.T) {
			// Test file-based package creation
			outputPath := filepath.Join(os.TempDir(), "test-package."+format.ext)
			err := CreatePackageFromConfig(format.pkgType, configFile.Name(), outputPath)
			if err != nil {
				t.Errorf("CreatePackageFromConfig failed for %s: %v", format.pkgType, err)
			}
			defer os.Remove(outputPath)

			// Verify the file was created
			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				t.Errorf("Package file was not created for %s", format.pkgType)
			}

			// Test writer-based package creation
			var buf bytes.Buffer
			err = CreatePackageFromConfigWriter(format.pkgType, configFile.Name(), &buf)
			if err != nil {
				t.Errorf("CreatePackageFromConfigWriter failed for %s: %v", format.pkgType, err)
			}

			// Verify some content was written
			if buf.Len() == 0 {
				t.Errorf("No content was written for %s", format.pkgType)
			}
		})
	}

	// Test with invalid config file
	t.Run("InvalidConfig", func(t *testing.T) {
		err := CreatePackageFromConfig(DEB, "nonexistent.yaml", "output.deb")
		if err == nil {
			t.Error("Expected error for invalid config, got nil")
		}
	})
}
