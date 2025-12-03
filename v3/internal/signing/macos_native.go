//go:build darwin

package signing

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// MacOSNativeSigner uses codesign for signing macOS binaries
type MacOSNativeSigner struct {
	codesignPath string
}

// NewMacOSNativeSigner creates a new macOS native signer
func NewMacOSNativeSigner() *MacOSNativeSigner {
	path, _ := exec.LookPath("codesign")
	return &MacOSNativeSigner{
		codesignPath: path,
	}
}

// Backend returns the signer backend type
func (s *MacOSNativeSigner) Backend() SignerBackend {
	return BackendNative
}

// Available returns true if codesign is available
func (s *MacOSNativeSigner) Available() bool {
	return s.codesignPath != ""
}

// Sign signs a macOS binary or app bundle using codesign
func (s *MacOSNativeSigner) Sign(ctx context.Context, req SignRequest) (*SignResult, error) {
	if !s.Available() {
		return nil, fmt.Errorf("codesign not found")
	}

	// Determine the identity to use
	identity := req.Certificate.Identity
	if identity == "" {
		// Try to auto-detect Developer ID
		if id, err := FindDeveloperIDIdentity(); err == nil {
			identity = id.Hash
		} else {
			// Fall back to ad-hoc signing
			identity = "-"
		}
	}

	// Handle output path
	inputPath := req.InputPath
	outputPath := req.OutputPath
	if outputPath == "" {
		outputPath = inputPath
	} else if outputPath != inputPath {
		// Copy the file/bundle first
		if err := copyPath(inputPath, outputPath); err != nil {
			return nil, fmt.Errorf("failed to copy for signing: %w", err)
		}
		inputPath = outputPath
	}

	// Check if this is an app bundle
	if req.BundleSign || strings.HasSuffix(inputPath, ".app") {
		if err := s.signAppBundle(ctx, inputPath, identity, req); err != nil {
			return nil, err
		}
	} else {
		if err := s.signBinary(ctx, inputPath, identity, req); err != nil {
			return nil, err
		}
	}

	return &SignResult{
		Backend:    BackendNative,
		OutputPath: outputPath,
	}, nil
}

// signBinary signs a single binary file
func (s *MacOSNativeSigner) signBinary(ctx context.Context, path, identity string, req SignRequest) error {
	args := []string{}

	args = append(args, "--force")
	args = append(args, "--sign", identity)

	if req.HardenedRuntime {
		args = append(args, "--options", "runtime")
	}

	if req.Entitlements != "" {
		if _, err := os.Stat(req.Entitlements); os.IsNotExist(err) {
			return fmt.Errorf("entitlements file does not exist: %s", req.Entitlements)
		}
		args = append(args, "--entitlements", req.Entitlements)
	}

	if req.Verbose {
		args = append(args, "--verbose")
	}

	args = append(args, path)

	cmd := exec.CommandContext(ctx, s.codesignPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("codesign failed: %w", err)
	}

	return nil
}

// signAppBundle signs an entire .app bundle with proper handling of nested components
func (s *MacOSNativeSigner) signAppBundle(ctx context.Context, bundlePath, identity string, req SignRequest) error {
	contentsPath := filepath.Join(bundlePath, "Contents")
	macOSPath := filepath.Join(contentsPath, "MacOS")
	frameworksPath := filepath.Join(contentsPath, "Frameworks")

	// Sign frameworks first if they exist
	if info, err := os.Stat(frameworksPath); err == nil && info.IsDir() {
		entries, err := os.ReadDir(frameworksPath)
		if err != nil {
			return fmt.Errorf("failed to read Frameworks directory: %w", err)
		}

		for _, entry := range entries {
			frameworkPath := filepath.Join(frameworksPath, entry.Name())
			if err := s.signBinary(ctx, frameworkPath, identity, SignRequest{
				HardenedRuntime: req.HardenedRuntime,
				Verbose:         req.Verbose,
			}); err != nil {
				return fmt.Errorf("failed to sign framework %s: %w", entry.Name(), err)
			}
		}
	}

	// Sign any dylibs in MacOS directory
	if info, err := os.Stat(macOSPath); err == nil && info.IsDir() {
		entries, err := os.ReadDir(macOSPath)
		if err != nil {
			return fmt.Errorf("failed to read MacOS directory: %w", err)
		}

		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".dylib") {
				dylibPath := filepath.Join(macOSPath, entry.Name())
				if err := s.signBinary(ctx, dylibPath, identity, SignRequest{
					HardenedRuntime: req.HardenedRuntime,
					Verbose:         req.Verbose,
				}); err != nil {
					return fmt.Errorf("failed to sign dylib %s: %w", entry.Name(), err)
				}
			}
		}
	}

	// Sign the main app bundle
	return s.signBinary(ctx, bundlePath, identity, req)
}

// Verify verifies the signature on a macOS binary or bundle
func (s *MacOSNativeSigner) Verify(ctx context.Context, path string) error {
	if !s.Available() {
		return fmt.Errorf("codesign not found")
	}

	cmd := exec.CommandContext(ctx, s.codesignPath, "--verify", "--verbose=2", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("signature verification failed: %s", stderr.String())
	}

	return nil
}

// copyPath copies a file or directory from src to dst
func copyPath(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

// copyDir copies a directory recursively
func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
