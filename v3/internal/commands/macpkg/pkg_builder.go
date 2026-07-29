package macpkg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

// PKGBuilder handles Mac .pkg creation using pkgbuild and productbuild
type PKGBuilder struct {
	Config   *PKGConfig
	tempDir  string
	buildDir string
}

// PKGConfig holds configuration for Mac pkg building
type PKGConfig struct {
	// App information
	AppName        string `yaml:"app_name"`
	AppPath        string `yaml:"app_path"`
	BundleID       string `yaml:"bundle_id"`
	Version        string `yaml:"version"`
	
	// Signing
	SigningIdentity string `yaml:"signing_identity"`
	
	// Installation
	InstallLocation string `yaml:"install_location"`
	
	// Distribution
	Title           string `yaml:"title"`
	Background      string `yaml:"background"`
	WelcomeFile     string `yaml:"welcome_file"`
	ReadmeFile      string `yaml:"readme_file"`
	LicenseFile     string `yaml:"license_file"`
	
	// Notarization
	AppleID         string `yaml:"apple_id"`
	AppPassword     string `yaml:"app_password"`
	TeamID          string `yaml:"team_id"`
	
	// Output
	OutputPath      string `yaml:"output_path"`
}

// NewPKGBuilder creates a new PKG builder instance
func NewPKGBuilder(config *PKGConfig) (*PKGBuilder, error) {
	// Create temporary directories
	tempDir, err := os.MkdirTemp("", "wails-pkg-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	
	buildDir := filepath.Join(tempDir, "build")
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create build directory: %w", err)
	}
	
	return &PKGBuilder{
		Config:   config,
		tempDir:  tempDir,
		buildDir: buildDir,
	}, nil
}

// Build creates the Mac pkg installer
func (b *PKGBuilder) Build() error {
	if err := b.validateConfig(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	
	// Step 1: Create component package with pkgbuild
	componentPkgPath := filepath.Join(b.buildDir, "component.pkg")
	if err := b.createComponentPackage(componentPkgPath); err != nil {
		return fmt.Errorf("failed to create component package: %w", err)
	}
	
	// Step 2: Create distribution XML
	distributionPath := filepath.Join(b.buildDir, "distribution.xml")
	if err := b.createDistributionXML(distributionPath); err != nil {
		return fmt.Errorf("failed to create distribution XML: %w", err)
	}
	
	// Step 3: Create product package with productbuild
	if err := b.createProductPackage(distributionPath, componentPkgPath); err != nil {
		return fmt.Errorf("failed to create product package: %w", err)
	}
	
	return nil
}

// createComponentPackage uses pkgbuild to create the component package
func (b *PKGBuilder) createComponentPackage(outputPath string) error {
	args := []string{
		"pkgbuild",
		"--root", b.Config.AppPath,
		"--identifier", b.Config.BundleID,
		"--version", b.Config.Version,
		"--install-location", b.Config.InstallLocation,
	}
	
	// Add signing if configured
	if b.Config.SigningIdentity != "" {
		args = append(args, "--sign", b.Config.SigningIdentity)
	}
	
	args = append(args, outputPath)
	
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pkgbuild failed: %w", err)
	}
	
	return nil
}

// createDistributionXML generates the distribution XML file
func (b *PKGBuilder) createDistributionXML(outputPath string) error {
	tmpl := `<?xml version="1.0" encoding="utf-8"?>
<installer-gui-script minSpecVersion="1">
    <title>{{.Title}}</title>
    <pkg-ref id="{{.BundleID}}"/>
    
    <!-- Domain configuration -->
    <domains enable_localSystem="true"/>
    
    <!-- Options -->
    <options customize="never" require-scripts="false" rootVolumeOnly="true" hostArchitectures="x86_64,arm64"/>
    
    <!-- Choices outline -->
    <choices-outline>
        <line choice="default">
            <line choice="{{.BundleID}}"/>
        </line>
    </choices-outline>
    
    <!-- Choice definitions -->
    <choice id="default" title="{{.Title}}"/>
    <choice id="{{.BundleID}}" visible="false" title="{{.AppName}}">
        <pkg-ref id="{{.BundleID}}"/>
    </choice>
    
    <!-- Package reference -->
    <pkg-ref id="{{.BundleID}}" version="{{.Version}}" onConclusion="none">component.pkg</pkg-ref>
    
    {{if .Background}}<background file="{{.Background}}" mime-type="image/png" alignment="topleft" scaling="none"/>{{end}}
    {{if .WelcomeFile}}<welcome file="{{.WelcomeFile}}" mime-type="text/rtf"/>{{end}}
    {{if .ReadmeFile}}<readme file="{{.ReadmeFile}}" mime-type="text/rtf"/>{{end}}
    {{if .LicenseFile}}<license file="{{.LicenseFile}}" mime-type="text/rtf"/>{{end}}
</installer-gui-script>`
	
	t, err := template.New("distribution").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse distribution template: %w", err)
	}
	
	var buf bytes.Buffer
	if err := t.Execute(&buf, b.Config); err != nil {
		return fmt.Errorf("failed to execute distribution template: %w", err)
	}
	
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write distribution XML: %w", err)
	}
	
	return nil
}

// createProductPackage uses productbuild to create the final installer
func (b *PKGBuilder) createProductPackage(distributionPath, componentPkgPath string) error {
	args := []string{
		"productbuild",
		"--distribution", distributionPath,
		"--package-path", b.buildDir,
		"--version", b.Config.Version,
	}
	
	// Add signing if configured
	if b.Config.SigningIdentity != "" {
		args = append(args, "--sign", b.Config.SigningIdentity)
	}
	
	args = append(args, b.Config.OutputPath)
	
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("productbuild failed: %w", err)
	}
	
	return nil
}

// validateConfig ensures all required configuration is present
func (b *PKGBuilder) validateConfig() error {
	if b.Config.AppName == "" {
		return fmt.Errorf("app_name is required")
	}
	if b.Config.AppPath == "" {
		return fmt.Errorf("app_path is required")
	}
	if b.Config.BundleID == "" {
		return fmt.Errorf("bundle_id is required")
	}
	if b.Config.Version == "" {
		return fmt.Errorf("version is required")
	}
	if b.Config.OutputPath == "" {
		return fmt.Errorf("output_path is required")
	}
	
	// Verify app path exists
	if _, err := os.Stat(b.Config.AppPath); os.IsNotExist(err) {
		return fmt.Errorf("app_path does not exist: %s", b.Config.AppPath)
	}
	
	// Set default install location if not specified
	if b.Config.InstallLocation == "" {
		b.Config.InstallLocation = fmt.Sprintf("/Applications/%s.app", b.Config.AppName)
	}
	
	// Set default title if not specified
	if b.Config.Title == "" {
		b.Config.Title = b.Config.AppName
	}
	
	return nil
}

// Cleanup removes temporary directories
func (b *PKGBuilder) Cleanup() error {
	if b.tempDir != "" {
		return os.RemoveAll(b.tempDir)
	}
	return nil
}

// CheckDependencies verifies that required tools are available
func CheckDependencies() error {
	tools := []string{"pkgbuild", "productbuild"}
	
	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err != nil {
			return fmt.Errorf("required tool '%s' not found in PATH", tool)
		}
	}
	
	return nil
}