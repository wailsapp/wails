package signing

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// WindowsNativeSigner uses signtool.exe for signing Windows binaries
type WindowsNativeSigner struct {
	signtoolPath string
}

// NewWindowsNativeSigner creates a new Windows native signer
func NewWindowsNativeSigner() *WindowsNativeSigner {
	return &WindowsNativeSigner{
		signtoolPath: findSigntool(),
	}
}

// findSigntool attempts to locate signtool.exe on the system
func findSigntool() string {
	if runtime.GOOS != "windows" {
		return ""
	}

	// Check if signtool is in PATH
	if path, err := exec.LookPath("signtool.exe"); err == nil {
		return path
	}
	if path, err := exec.LookPath("signtool"); err == nil {
		return path
	}

	// Common Windows SDK locations
	programFiles := os.Getenv("ProgramFiles(x86)")
	if programFiles == "" {
		programFiles = `C:\Program Files (x86)`
	}

	sdkRoot := filepath.Join(programFiles, "Windows Kits", "10", "bin")

	// Find the latest SDK version
	entries, err := os.ReadDir(sdkRoot)
	if err != nil {
		return ""
	}

	var latestVersion string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "10.") {
			if entry.Name() > latestVersion {
				latestVersion = entry.Name()
			}
		}
	}

	if latestVersion != "" {
		// Try x64 first, then x86
		for _, arch := range []string{"x64", "x86"} {
			path := filepath.Join(sdkRoot, latestVersion, arch, "signtool.exe")
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	return ""
}

// Backend returns the signer backend type
func (s *WindowsNativeSigner) Backend() SignerBackend {
	return BackendNative
}

// Available returns true if signtool.exe is available
func (s *WindowsNativeSigner) Available() bool {
	return s.signtoolPath != ""
}

// Sign signs a Windows binary using signtool.exe
func (s *WindowsNativeSigner) Sign(ctx context.Context, req SignRequest) (*SignResult, error) {
	if !s.Available() {
		return nil, fmt.Errorf("signtool.exe not found")
	}

	args := []string{"sign"}

	// Certificate configuration
	if req.Certificate.PKCS12Path != "" {
		args = append(args, "/f", req.Certificate.PKCS12Path)
		if req.Certificate.PKCS12Password != "" {
			args = append(args, "/p", req.Certificate.PKCS12Password)
		}
	} else if req.Certificate.Thumbprint != "" {
		args = append(args, "/sha1", req.Certificate.Thumbprint)
	} else {
		return nil, fmt.Errorf("either PKCS12Path or Thumbprint must be specified")
	}

	// Signature algorithm (default to SHA256)
	args = append(args, "/fd", "SHA256")

	// Timestamp server
	if req.TimestampServer != "" {
		args = append(args, "/t", req.TimestampServer)
	} else {
		// Default timestamp server
		args = append(args, "/t", "http://timestamp.digicert.com")
	}

	// Description
	if req.Description != "" {
		args = append(args, "/d", req.Description)
	}

	// URL
	if req.URL != "" {
		args = append(args, "/du", req.URL)
	}

	// Verbose
	if req.Verbose {
		args = append(args, "/v")
	}

	// Output path handling
	inputPath := req.InputPath
	outputPath := req.OutputPath
	if outputPath == "" {
		outputPath = inputPath
	} else if outputPath != inputPath {
		// Copy file first since signtool modifies in place
		if err := copyFile(inputPath, outputPath); err != nil {
			return nil, fmt.Errorf("failed to copy file for signing: %w", err)
		}
		inputPath = outputPath
	}

	args = append(args, inputPath)

	cmd := exec.CommandContext(ctx, s.signtoolPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("signtool failed: %w", err)
	}

	return &SignResult{
		Backend:    BackendNative,
		OutputPath: outputPath,
	}, nil
}

// Verify verifies the signature on a Windows binary
func (s *WindowsNativeSigner) Verify(ctx context.Context, path string) error {
	if !s.Available() {
		return fmt.Errorf("signtool.exe not found")
	}

	cmd := exec.CommandContext(ctx, s.signtoolPath, "verify", "/pa", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("signature verification failed: %s", string(output))
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, info.Mode())
}
