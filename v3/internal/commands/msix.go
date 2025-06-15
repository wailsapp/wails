package commands

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/flags"
)

//go:embed build_assets/windows/msix/*
var msixAssets embed.FS

// MSIXOptions represents the configuration for MSIX packaging
type MSIXOptions struct {
	// Info from project config
	Info struct {
		CompanyName       string `json:"companyName"`
		ProductName       string `json:"productName"`
		ProductVersion    string `json:"version"`
		ProductIdentifier string `json:"productIdentifier"`
		Description       string `json:"description"`
		Copyright         string `json:"copyright"`
		Comments          string `json:"comments"`
	}
	// File associations
	FileAssociations []struct {
		Ext         string `json:"ext"`
		Name        string `json:"name"`
		Description string `json:"description"`
		IconName    string `json:"iconName"`
		Role        string `json:"role"`
		MimeType    string `json:"mimeType,omitempty"`
	} `json:"fileAssociations"`
	// MSIX specific options
	Publisher             string `json:"publisher"`
	CertificatePath       string `json:"certificatePath"`
	CertificatePassword   string `json:"certificatePassword,omitempty"`
	ProcessorArchitecture string `json:"processorArchitecture"`
	ExecutableName        string `json:"executableName"`
	ExecutablePath        string `json:"executablePath"`
	OutputPath            string `json:"outputPath"`
	UseMsixPackagingTool  bool   `json:"useMsixPackagingTool"`
	UseMakeAppx           bool   `json:"useMakeAppx"`
}

// ToolMSIX creates an MSIX package for Windows applications
func ToolMSIX(options *flags.ToolMSIX) error {
	DisableFooter = true

	if runtime.GOOS != "windows" {
		return fmt.Errorf("MSIX packaging is only supported on Windows")
	}

	// Check if required tools are installed
	if err := checkMSIXTools(options); err != nil {
		return err
	}

	// Load project configuration
	configPath := options.ConfigPath
	if configPath == "" {
		configPath = "wails.json"
	}

	// Read the config file
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Parse the config
	var config struct {
		Info             map[string]interface{}   `json:"info"`
		FileAssociations []map[string]interface{} `json:"fileAssociations"`
	}
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	// Create MSIX options
	msixOptions := MSIXOptions{
		Publisher:             options.Publisher,
		CertificatePath:       options.CertificatePath,
		CertificatePassword:   options.CertificatePassword,
		ProcessorArchitecture: options.Arch,
		ExecutableName:        options.ExecutableName,
		ExecutablePath:        options.ExecutablePath,
		OutputPath:            options.OutputPath,
		UseMsixPackagingTool:  options.UseMsixPackagingTool,
		UseMakeAppx:           options.UseMakeAppx,
	}

	// Copy info from config
	infoBytes, err := json.Marshal(config.Info)
	if err != nil {
		return fmt.Errorf("error marshaling info: %w", err)
	}
	if err := json.Unmarshal(infoBytes, &msixOptions.Info); err != nil {
		return fmt.Errorf("error unmarshaling info: %w", err)
	}

	// Copy file associations from config
	if len(config.FileAssociations) > 0 {
		faBytes, err := json.Marshal(config.FileAssociations)
		if err != nil {
			return fmt.Errorf("error marshaling file associations: %w", err)
		}
		if err := json.Unmarshal(faBytes, &msixOptions.FileAssociations); err != nil {
			return fmt.Errorf("error unmarshaling file associations: %w", err)
		}
	}

	// Validate options
	if err := validateMSIXOptions(&msixOptions); err != nil {
		return err
	}

	// Create MSIX package
	if msixOptions.UseMsixPackagingTool {
		return createMSIXWithPackagingTool(&msixOptions)
	} else if msixOptions.UseMakeAppx {
		return createMSIXWithMakeAppx(&msixOptions)
	}

	// Default to MakeAppx if neither is specified
	return createMSIXWithMakeAppx(&msixOptions)
}

// checkMSIXTools checks if the required tools for MSIX packaging are installed
func checkMSIXTools(options *flags.ToolMSIX) error {
	// Check if MsixPackagingTool is installed if requested
	if options.UseMsixPackagingTool {
		cmd := exec.Command("powershell", "-Command", "Get-AppxPackage -Name Microsoft.MsixPackagingTool")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Microsoft MSIX Packaging Tool is not installed. Please install it from the Microsoft Store")
		}
	}

	// Check if MakeAppx is available if requested
	if options.UseMakeAppx {
		cmd := exec.Command("where", "MakeAppx.exe")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("MakeAppx.exe is not found in PATH. Please install the Windows SDK")
		}
	}

	// If neither is specified, check for MakeAppx as the default
	if !options.UseMsixPackagingTool && !options.UseMakeAppx {
		cmd := exec.Command("where", "MakeAppx.exe")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("MakeAppx.exe is not found in PATH. Please install the Windows SDK")
		}
	}

	// Check if signtool is available for signing
	if options.CertificatePath != "" {
		cmd := exec.Command("where", "signtool.exe")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("signtool.exe is not found in PATH. Please install the Windows SDK")
		}
	}

	return nil
}

// validateMSIXOptions validates the MSIX options
func validateMSIXOptions(options *MSIXOptions) error {
	// Check required fields
	if options.Info.ProductName == "" {
		return fmt.Errorf("product name is required")
	}
	if options.Info.ProductIdentifier == "" {
		return fmt.Errorf("product identifier is required")
	}
	if options.Info.CompanyName == "" {
		return fmt.Errorf("company name is required")
	}
	if options.ExecutableName == "" {
		return fmt.Errorf("executable name is required")
	}
	if options.ExecutablePath == "" {
		return fmt.Errorf("executable path is required")
	}

	// Validate executable path
	if _, err := os.Stat(options.ExecutablePath); os.IsNotExist(err) {
		return fmt.Errorf("executable file not found: %s", options.ExecutablePath)
	}

	// Validate certificate path if provided
	if options.CertificatePath != "" {
		if _, err := os.Stat(options.CertificatePath); os.IsNotExist(err) {
			return fmt.Errorf("certificate file not found: %s", options.CertificatePath)
		}
	}

	// Set default processor architecture if not provided
	if options.ProcessorArchitecture == "" {
		options.ProcessorArchitecture = "x64"
	}

	// Set default publisher if not provided
	if options.Publisher == "" {
		options.Publisher = fmt.Sprintf("CN=%s", options.Info.CompanyName)
	}

	// Set default output path if not provided
	if options.OutputPath == "" {
		options.OutputPath = filepath.Join(".", fmt.Sprintf("%s.msix", options.Info.ProductName))
	}

	return nil
}

// createMSIXWithPackagingTool creates an MSIX package using the Microsoft MSIX Packaging Tool
func createMSIXWithPackagingTool(options *MSIXOptions) error {
	// Create a temporary directory for the template
	tempDir, err := os.MkdirTemp("", "wails-msix-")
	if err != nil {
		return fmt.Errorf("error creating temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Generate the template file
	templatePath := filepath.Join(tempDir, "template.xml")
	if err := generateMSIXTemplate(options, templatePath); err != nil {
		return fmt.Errorf("error generating MSIX template: %w", err)
	}

	// Create the MSIX package
	fmt.Println("Creating MSIX package using Microsoft MSIX Packaging Tool...")
	args := []string{"create-package", "--template", templatePath}

	// Add certificate password if provided
	if options.CertificatePassword != "" {
		args = append(args, "--certPassword", options.CertificatePassword)
	}

	cmd := exec.Command("MsixPackagingTool.exe", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error creating MSIX package: %w", err)
	}

	fmt.Printf("MSIX package created successfully: %s\n", options.OutputPath)
	return nil
}

// createMSIXWithMakeAppx creates an MSIX package using MakeAppx.exe
func createMSIXWithMakeAppx(options *MSIXOptions) error {
	// Create a temporary directory for the package structure
	tempDir, err := os.MkdirTemp("", "wails-msix-")
	if err != nil {
		return fmt.Errorf("error creating temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the package structure
	if err := createMSIXPackageStructure(options, tempDir); err != nil {
		return fmt.Errorf("error creating MSIX package structure: %w", err)
	}

	// Create the MSIX package
	fmt.Println("Creating MSIX package using MakeAppx.exe...")
	cmd := exec.Command("MakeAppx.exe", "pack", "/d", tempDir, "/p", options.OutputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error creating MSIX package: %w", err)
	}

	// Sign the package if certificate is provided
	if options.CertificatePath != "" {
		fmt.Println("Signing MSIX package...")
		signArgs := []string{"sign", "/fd", "SHA256", "/a", "/f", options.CertificatePath}

		// Add certificate password if provided
		if options.CertificatePassword != "" {
			signArgs = append(signArgs, "/p", options.CertificatePassword)
		}

		signArgs = append(signArgs, options.OutputPath)

		cmd = exec.Command("signtool.exe", signArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error signing MSIX package: %w", err)
		}
	}

	fmt.Printf("MSIX package created successfully: %s\n", options.OutputPath)
	return nil
}

// generateMSIXTemplate generates the MSIX template file for the Microsoft MSIX Packaging Tool
func generateMSIXTemplate(options *MSIXOptions, outputPath string) error {
	// Read the template file
	templateData, err := msixAssets.ReadFile("build_assets/windows/msix/template.xml.tmpl")
	if err != nil {
		return fmt.Errorf("error reading template file: %w", err)
	}

	// Parse the template
	tmpl, err := template.New("msix-template").Parse(string(templateData))
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer file.Close()

	// Execute the template
	if err := tmpl.Execute(file, options); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

// createMSIXPackageStructure creates the MSIX package structure for MakeAppx.exe
func createMSIXPackageStructure(options *MSIXOptions, outputDir string) error {
	// Create the Assets directory
	assetsDir := filepath.Join(outputDir, "Assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return fmt.Errorf("error creating Assets directory: %w", err)
	}

	// Generate the AppxManifest.xml file
	manifestPath := filepath.Join(outputDir, "AppxManifest.xml")
	if err := generateAppxManifest(options, manifestPath); err != nil {
		return fmt.Errorf("error generating AppxManifest.xml: %w", err)
	}

	// Copy the executable
	executableDest := filepath.Join(outputDir, filepath.Base(options.ExecutablePath))
	if err := copyFile(options.ExecutablePath, executableDest); err != nil {
		return fmt.Errorf("error copying executable: %w", err)
	}

	// Copy any additional files needed for the application
	// This would include DLLs, resources, etc.
	// For now, we'll just copy the executable

	// Generate placeholder assets
	assets := []string{
		"Square150x150Logo.png",
		"Square44x44Logo.png",
		"Wide310x150Logo.png",
		"SplashScreen.png",
		"StoreLogo.png",
	}

	// Add FileIcon.png if there are file associations
	if len(options.FileAssociations) > 0 {
		assets = append(assets, "FileIcon.png")
	}

	// Generate placeholder assets
	for _, asset := range assets {
		assetPath := filepath.Join(assetsDir, asset)
		if err := generatePlaceholderImage(assetPath); err != nil {
			return fmt.Errorf("error generating placeholder image %s: %w", asset, err)
		}
	}

	return nil
}

// generateAppxManifest generates the AppxManifest.xml file
func generateAppxManifest(options *MSIXOptions, outputPath string) error {
	// Read the template file
	templateData, err := msixAssets.ReadFile("build_assets/windows/msix/app_manifest.xml.tmpl")
	if err != nil {
		return fmt.Errorf("error reading template file: %w", err)
	}

	// Parse the template
	tmpl, err := template.New("appx-manifest").Parse(string(templateData))
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer file.Close()

	// Execute the template
	if err := tmpl.Execute(file, options); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

// generatePlaceholderImage generates a placeholder image file
func generatePlaceholderImage(outputPath string) error {
	// For now, we'll create a simple 1x1 transparent PNG
	// In a real implementation, we would generate proper icons based on the application icon

	// Create a minimal valid PNG file (1x1 transparent pixel)
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
		0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	return os.WriteFile(outputPath, pngData, 0644)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	// Read the source file
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Write the destination file
	return os.WriteFile(dst, data, 0644)
}

// InstallMSIXTools installs the required tools for MSIX packaging
func InstallMSIXTools() error {
	// Check if running on Windows
	if runtime.GOOS != "windows" {
		return fmt.Errorf("MSIX packaging is only supported on Windows")
	}

	fmt.Println("Installing MSIX packaging tools...")

	// Install MSIX Packaging Tool from Microsoft Store
	fmt.Println("Installing Microsoft MSIX Packaging Tool from Microsoft Store...")
	cmd := exec.Command("powershell", "-Command", "Start-Process ms-windows-store://pdp/?ProductId=9N5R1TQPJVBP")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error launching Microsoft Store: %w", err)
	}

	// Check if Windows SDK is installed
	fmt.Println("Checking for Windows SDK...")
	sdkInstalled := false
	cmd = exec.Command("where", "MakeAppx.exe")
	if err := cmd.Run(); err == nil {
		sdkInstalled = true
		fmt.Println("Windows SDK is already installed.")
	}

	// Install Windows SDK if not installed
	if !sdkInstalled {
		fmt.Println("Windows SDK is not installed. Please download and install from:")
		fmt.Println("https://developer.microsoft.com/en-us/windows/downloads/windows-sdk/")

		// Open the download page
		cmd = exec.Command("powershell", "-Command", "Start-Process https://developer.microsoft.com/en-us/windows/downloads/windows-sdk/")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error opening Windows SDK download page: %w", err)
		}
	}

	fmt.Println("MSIX packaging tools installation initiated. Please complete the installation process in the opened windows.")
	return nil
}

// init registers the MSIX command
func init() {
	// Register the MSIX command in the CLI
	// This will be called by the CLI framework
}
