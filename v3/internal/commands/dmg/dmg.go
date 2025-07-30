package dmg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Creator handles DMG creation
type Creator struct {
	sourcePath      string
	outputPath      string
	appName         string
	backgroundImage string
	iconPositions   map[string]Position
}

// Position represents icon coordinates in the DMG
type Position struct {
	X, Y int
}

// New creates a new DMG creator
func New(sourcePath, outputPath, appName string) (*Creator, error) {
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("DMG creation is only supported on macOS")
	}

	// Check if source exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("source path does not exist: %s", sourcePath)
	}

	return &Creator{
		sourcePath:    sourcePath,
		outputPath:    outputPath,
		appName:       appName,
		iconPositions: make(map[string]Position),
	}, nil
}

// SetBackgroundImage sets the background image for the DMG
func (c *Creator) SetBackgroundImage(imagePath string) error {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("background image does not exist: %s", imagePath)
	}
	c.backgroundImage = imagePath
	return nil
}

// AddIconPosition adds an icon position for the DMG layout
func (c *Creator) AddIconPosition(filename string, x, y int) {
	c.iconPositions[filename] = Position{X: x, Y: y}
}

// Create creates the DMG file
func (c *Creator) Create() error {
	// Remove existing DMG if it exists
	if _, err := os.Stat(c.outputPath); err == nil {
		if err := os.Remove(c.outputPath); err != nil {
			return fmt.Errorf("failed to remove existing DMG: %w", err)
		}
	}

	// Create a temporary directory for DMG content
	tempDir, err := os.MkdirTemp("", "dmg-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy the app bundle to temp directory
	appName := filepath.Base(c.sourcePath)
	tempAppPath := filepath.Join(tempDir, appName)
	if err := c.copyDir(c.sourcePath, tempAppPath); err != nil {
		return fmt.Errorf("failed to copy app bundle: %w", err)
	}

	// Create Applications symlink
	applicationsLink := filepath.Join(tempDir, "Applications")
	if err := os.Symlink("/Applications", applicationsLink); err != nil {
		return fmt.Errorf("failed to create Applications symlink: %w", err)
	}

	// Copy background image if provided
	if c.backgroundImage != "" {
		bgName := filepath.Base(c.backgroundImage)
		bgPath := filepath.Join(tempDir, bgName)
		if err := c.copyFile(c.backgroundImage, bgPath); err != nil {
			return fmt.Errorf("failed to copy background image: %w", err)
		}
	}

	// Create DMG using hdiutil
	if err := c.createDMGWithHdiutil(tempDir); err != nil {
		return fmt.Errorf("failed to create DMG with hdiutil: %w", err)
	}

	return nil
}

// copyDir recursively copies a directory
func (c *Creator) copyDir(src, dst string) error {
	cmd := exec.Command("cp", "-R", src, dst)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy directory: %w", err)
	}
	return nil
}

// copyFile copies a file
func (c *Creator) copyFile(src, dst string) error {
	cmd := exec.Command("cp", src, dst)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	return nil
}

// createDMGWithHdiutil creates the DMG using macOS hdiutil
func (c *Creator) createDMGWithHdiutil(sourceDir string) error {
	// Calculate size needed for DMG (roughly 2x the source size for safety)
	sizeCmd := exec.Command("du", "-sk", sourceDir)
	output, err := sizeCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to calculate directory size: %w", err)
	}

	// Parse size and add padding
	sizeStr := strings.Fields(string(output))[0]

	// Create DMG with hdiutil
	args := []string{
		"create",
		"-srcfolder", sourceDir,
		"-format", "UDZO",
		"-volname", c.appName,
		c.outputPath,
	}

	// Add size if we could determine it
	if sizeStr != "" {
		// Add 50% padding to the calculated size
		args = append([]string{"create", "-size", sizeStr + "k"}, args[1:]...)
		args[0] = "create"
	}

	cmd := exec.Command("hdiutil", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("hdiutil failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}
