package macpkg

import (
	"fmt"
	"os"
	"runtime"
)

// BuildMacPKGOptions holds options for the Mac PKG building command
type BuildMacPKGOptions struct {
	ConfigPath        string `name:"config" description:"Path to PKG configuration file" default:""`
	GenerateTemplate  bool   `name:"generate-template" description:"Generate a sample configuration file" default:"false"`
	SkipNotarization  bool   `name:"skip-notarization" description:"Skip notarization step" default:"false"`
	ValidateOnly      bool   `name:"validate-only" description:"Only validate configuration without building" default:"false"`
	Verbose           bool   `name:"verbose" description:"Enable verbose output" default:"false"`
}

// BuildMacPKG is the main entry point for Mac PKG building
func BuildMacPKG(options *BuildMacPKGOptions) error {
	// Check if running on macOS
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("Mac PKG building is only supported on macOS")
	}
	
	// Handle template generation
	if options.GenerateTemplate {
		return handleTemplateGeneration(options.ConfigPath)
	}
	
	// Validate config path is provided
	if options.ConfigPath == "" {
		return fmt.Errorf("config file is required. Use --config to specify a configuration file, or --generate-template to create one")
	}
	
	// Check dependencies first
	if err := CheckDependencies(); err != nil {
		return fmt.Errorf("dependency check failed: %w", err)
	}
	
	// Check notarization dependencies if not skipping
	if !options.SkipNotarization {
		if err := CheckNotarizationDependencies(); err != nil {
			fmt.Printf("Warning: Notarization dependencies not available: %v\n", err)
			fmt.Println("Continuing without notarization...")
			options.SkipNotarization = true
		}
	}
	
	// Load configuration
	config, err := LoadPKGConfig(options.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	
	// Validate configuration
	if err := ValidateConfig(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}
	
	// If validate-only, stop here
	if options.ValidateOnly {
		fmt.Println("Configuration is valid!")
		return nil
	}
	
	// Create PKG builder
	builder, err := NewPKGBuilder(config)
	if err != nil {
		return fmt.Errorf("failed to create PKG builder: %w", err)
	}
	defer builder.Cleanup()
	
	// Build the package
	fmt.Printf("Building Mac PKG installer for %s...\n", config.AppName)
	if err := builder.Build(); err != nil {
		return fmt.Errorf("PKG build failed: %w", err)
	}
	
	fmt.Printf("PKG installer created successfully: %s\n", config.OutputPath)
	
	// Notarization (if enabled and credentials provided)
	if !options.SkipNotarization && shouldNotarize(config) {
		fmt.Println("Starting notarization process...")
		
		notarizer := NewNotarizationService(config.AppleID, config.AppPassword, config.TeamID)
		if err := notarizer.NotarizePackage(config.OutputPath); err != nil {
			return fmt.Errorf("notarization failed: %w", err)
		}
	} else if !options.SkipNotarization {
		fmt.Println("Skipping notarization: credentials not provided in configuration")
	}
	
	// Validate the final package
	fmt.Println("Validating package signature...")
	if err := ValidatePackageSignature(config.OutputPath); err != nil {
		fmt.Printf("Warning: Package signature validation failed: %v\n", err)
	} else {
		fmt.Println("Package signature validation passed!")
	}
	
	fmt.Printf("\nMac PKG installer ready: %s\n", config.OutputPath)
	return nil
}

// handleTemplateGeneration creates a sample configuration file
func handleTemplateGeneration(configPath string) error {
	templatePath := "wails-pkg.yaml"
	if configPath != "" {
		templatePath = configPath
	}
	
	// Check if file already exists
	if _, err := os.Stat(templatePath); err == nil {
		return fmt.Errorf("file already exists: %s", templatePath)
	}
	
	if err := GenerateTemplate(templatePath); err != nil {
		return fmt.Errorf("failed to generate template: %w", err)
	}
	
	fmt.Printf("Configuration template created: %s\n", templatePath)
	fmt.Println("Edit this file to customize your PKG installer settings.")
	return nil
}

// shouldNotarize determines if notarization should be attempted
func shouldNotarize(config *PKGConfig) bool {
	return config.AppleID != "" && config.AppPassword != "" && config.TeamID != ""
}