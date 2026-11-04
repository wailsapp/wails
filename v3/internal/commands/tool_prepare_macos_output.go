package commands

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"
)

type PrepareMacOSOutputOptions struct {
	Dir string `name:"dir" description:"Project build output directory" default:"bin"`
}

// ToolPrepareMacOSOutput keeps generated app bundles on a local filesystem when
// a file provider adds Finder metadata to .app directories. A symlink preserves
// the project's usual bin path while the actual bundle stays signable.
func ToolPrepareMacOSOutput(options *PrepareMacOSOutputOptions) error {
	DisableFooter = true
	if runtime.GOOS != "darwin" {
		return nil
	}
	if options == nil || options.Dir == "" {
		return fmt.Errorf("build output directory is required")
	}
	projectRoot, err := os.Getwd()
	if err != nil {
		return err
	}
	binPath, err := filepath.Abs(options.Dir)
	if err != nil {
		return err
	}
	insideProject, err := isProjectOutputDirectory(projectRoot, binPath)
	if err != nil {
		return err
	}
	// An explicitly selected directory outside the project remains the user's
	// chosen destination. Only generated output inside the project is relocated.
	if !insideProject {
		return os.MkdirAll(binPath, 0o755)
	}
	unlock, err := lockMacOSOutput(binPath)
	if err != nil {
		return err
	}
	defer unlock()
	linked, err := isOutputLink(binPath)
	if err != nil {
		return err
	}
	if linked {
		return ensureMacOSOutputIsLocal(binPath)
	}
	if err := os.MkdirAll(binPath, 0o755); err != nil {
		return err
	}
	if macOSOutputIsLocal(binPath) {
		return nil
	}
	provider, err := isFileProviderDirectory(binPath)
	if err != nil {
		return err
	}
	if !provider {
		return markMacOSOutputLocal(binPath)
	}
	return relocateMacOSOutput(projectRoot, binPath)
}

func isProjectOutputDirectory(projectRoot, binPath string) (bool, error) {
	relative, err := filepath.Rel(projectRoot, binPath)
	if err != nil {
		return false, err
	}
	if relative == "." {
		return false, fmt.Errorf("build output directory must not be the project root")
	}
	return relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)), nil
}

func isFileProviderDirectory(binPath string) (bool, error) {
	if provider, err := hasFileProviderAncestor(binPath); provider || err != nil {
		return provider, err
	}
	probe, err := os.MkdirTemp(binPath, "wails-probe-*.app")
	if err != nil {
		return false, err
	}
	defer os.Remove(probe)
	// Probe as a fallback for providers without domain metadata. A cached result
	// still checks ancestors so asynchronously attached metadata revokes trust.
	time.Sleep(2500 * time.Millisecond)
	output, err := exec.Command("xattr", probe).Output()
	if err != nil {
		return false, fmt.Errorf("read macOS signing metadata: %w", err)
	}
	attributes := strings.Split(strings.TrimSpace(string(output)), "\n")
	return slices.Contains(attributes, "com.apple.fileprovider.fpfs#P") ||
		slices.Contains(attributes, "com.apple.FinderInfo"), nil
}

func hasFileProviderAncestor(path string) (bool, error) {
	resolved, err := resolveMacOSOutput(path)
	if err != nil {
		return false, err
	}
	for {
		if provider, err := macOSDirectoryHasFileProviderMetadata(resolved); provider || err != nil {
			return provider, err
		}
		parent := filepath.Dir(resolved)
		if parent == resolved {
			return false, nil
		}
		resolved = parent
	}
}

func localMacOSOutputRoot() (string, error) {
	outputRoot := os.Getenv("WAILS_MACOS_OUTPUT_ROOT")
	if outputRoot == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		outputRoot = filepath.Join(configDir, "Wails", "build-output")
	}
	outputRoot, err := filepath.Abs(outputRoot)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(outputRoot, 0o700); err != nil {
		return "", err
	}
	if err := ensureMacOSOutputIsLocal(outputRoot); err != nil {
		return "", fmt.Errorf("WAILS_MACOS_OUTPUT_ROOT: %w", err)
	}
	return outputRoot, nil
}

func ensureMacOSOutputIsLocal(path string) error {
	if macOSOutputIsLocal(path) {
		return nil
	}
	provider, err := isFileProviderDirectory(path)
	if err != nil {
		return err
	}
	if provider {
		return fmt.Errorf("macOS output %s reattaches Finder metadata; choose an unsynced output directory", path)
	}
	return markMacOSOutputLocal(path)
}

func relocateMacOSOutput(projectRoot, binPath string) error {
	outputRoot, err := localMacOSOutputRoot()
	if err != nil {
		return err
	}
	sameFilesystem, err := sameMacOSFilesystem(binPath, outputRoot)
	if err != nil {
		return fmt.Errorf("check macOS output volume: %w", err)
	}
	if !sameFilesystem {
		return fmt.Errorf("cannot relocate build output across volumes; set WAILS_MACOS_OUTPUT_ROOT to an unsynced directory on the project's volume, or configure BIN_DIR directly outside the project")
	}
	digest := sha256.Sum256([]byte(binPath))
	target := filepath.Join(outputRoot, fmt.Sprintf("%s-%x", filepath.Base(projectRoot), digest[:8]))
	if _, err := os.Lstat(target); err == nil {
		return fmt.Errorf("local build output already exists at %s; resolve the existing bin directory before retrying", target)
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.Rename(binPath, target); err != nil {
		if linked, linkErr := isOutputLink(binPath); linked && linkErr == nil {
			return ensureMacOSOutputIsLocal(binPath)
		}
		return fmt.Errorf("move build output outside file provider: %w", err)
	}
	if err := os.Symlink(target, binPath); err != nil {
		if restoreErr := os.Rename(target, binPath); restoreErr != nil {
			return fmt.Errorf("link local build output: %w; could not restore the original directory: %v; output remains at %s", err, restoreErr, target)
		}
		return fmt.Errorf("link local build output into project: %w", err)
	}
	if err := markMacOSOutputLocal(target); err != nil {
		return fmt.Errorf("mark local build output at %s: %w", target, err)
	}
	fmt.Fprintf(os.Stderr, "Wails: macOS file provider adds signing metadata to app bundles; build output is stored at %s (linked from %s)\n", target, binPath)
	return nil
}

func macOSOutputIsLocal(path string) bool {
	resolved, err := resolveMacOSOutput(path)
	if err != nil {
		return false
	}
	data, err := os.ReadFile(filepath.Join(resolved, ".wails-provider-checked"))
	if err != nil || string(data) != "v2\n"+resolved {
		return false
	}
	provider, err := hasFileProviderAncestor(resolved)
	return err == nil && !provider
}

func markMacOSOutputLocal(path string) error {
	resolved, err := resolveMacOSOutput(path)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(resolved, ".wails-provider-checked"), []byte("v2\n"+resolved), 0o644)
}

func resolveMacOSOutput(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(absolute)
}

func isOutputLink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Stat(path)
		if err != nil {
			return false, fmt.Errorf("build output link %s is broken: %w", path, err)
		}
		if !target.IsDir() {
			return false, fmt.Errorf("build output link %s does not point to a directory", path)
		}
		return true, nil
	}
	if !info.IsDir() {
		return false, fmt.Errorf("build output %s is not a directory", path)
	}
	return false, nil
}
