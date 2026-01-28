// Package defaults provides functionality for loading and saving global default settings
// for Wails projects. Settings are stored in ~/.config/wails/defaults.yaml
package defaults

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// GlobalDefaults represents the user's default project settings
// These are stored in ~/.config/wails/defaults.yaml and used when creating new projects
type GlobalDefaults struct {
	// Author information
	Author AuthorDefaults `json:"author" yaml:"author"`

	// Default project settings
	Project ProjectDefaults `json:"project" yaml:"project"`

	// Code signing configuration (optional)
	Signing SigningDefaults `json:"signing,omitempty" yaml:"signing,omitempty"`
}

// AuthorDefaults contains the author's information
type AuthorDefaults struct {
	Name    string `json:"name" yaml:"name"`
	Company string `json:"company" yaml:"company"`
}

// ProjectDefaults contains default project settings
type ProjectDefaults struct {
	ProductIdentifierPrefix string `json:"productIdentifierPrefix" yaml:"productIdentifierPrefix"`
	DefaultTemplate         string `json:"defaultTemplate" yaml:"defaultTemplate"`
	Framework               string `json:"framework" yaml:"framework"`
	Language                string `json:"language" yaml:"language"`
	CopyrightTemplate       string `json:"copyrightTemplate" yaml:"copyrightTemplate"`
	DescriptionTemplate     string `json:"descriptionTemplate" yaml:"descriptionTemplate"`
	DefaultVersion          string `json:"defaultVersion" yaml:"defaultVersion"`
	UseInterfaces           bool   `json:"useInterfaces" yaml:"useInterfaces"`
}

// SigningDefaults contains code signing configuration for all platforms
type SigningDefaults struct {
	Darwin  DarwinSigningDefaults  `json:"darwin,omitempty" yaml:"darwin,omitempty"`
	Windows WindowsSigningDefaults `json:"windows,omitempty" yaml:"windows,omitempty"`
	Linux   LinuxSigningDefaults   `json:"linux,omitempty" yaml:"linux,omitempty"`
}

// DarwinSigningDefaults contains macOS code signing configuration
type DarwinSigningDefaults struct {
	Identity        string `json:"identity,omitempty" yaml:"identity,omitempty"`
	TeamID          string `json:"teamID,omitempty" yaml:"teamID,omitempty"`
	KeychainProfile string `json:"keychainProfile,omitempty" yaml:"keychainProfile,omitempty"`
	Entitlements    string `json:"entitlements,omitempty" yaml:"entitlements,omitempty"`
	P12Path         string `json:"p12Path,omitempty" yaml:"p12Path,omitempty"`
	APIKeyPath      string `json:"apiKeyPath,omitempty" yaml:"apiKeyPath,omitempty"`
	APIKeyID        string `json:"apiKeyID,omitempty" yaml:"apiKeyID,omitempty"`
	APIIssuerID     string `json:"apiIssuerID,omitempty" yaml:"apiIssuerID,omitempty"`
}

// WindowsSigningDefaults contains Windows code signing configuration
type WindowsSigningDefaults struct {
	CertificatePath string `json:"certificatePath,omitempty" yaml:"certificatePath,omitempty"`
	Thumbprint      string `json:"thumbprint,omitempty" yaml:"thumbprint,omitempty"`
	TimestampServer string `json:"timestampServer,omitempty" yaml:"timestampServer,omitempty"`
	CloudProvider   string `json:"cloudProvider,omitempty" yaml:"cloudProvider,omitempty"`
	CloudKeyID      string `json:"cloudKeyID,omitempty" yaml:"cloudKeyID,omitempty"`
}

// LinuxSigningDefaults contains Linux package signing configuration
type LinuxSigningDefaults struct {
	GPGKeyPath string `json:"gpgKeyPath,omitempty" yaml:"gpgKeyPath,omitempty"`
	GPGKeyID   string `json:"gpgKeyID,omitempty" yaml:"gpgKeyID,omitempty"`
	SignRole   string `json:"signRole,omitempty" yaml:"signRole,omitempty"`
}

// Default returns sensible defaults for first-time users
func Default() GlobalDefaults {
	return GlobalDefaults{
		Author: AuthorDefaults{
			Name:    "",
			Company: "",
		},
		Project: ProjectDefaults{
			ProductIdentifierPrefix: "com.example",
			DefaultTemplate:         "vanilla",
			CopyrightTemplate:       "© {year}, {company}",
			DescriptionTemplate:     "A {name} application",
			DefaultVersion:          "0.1.0",
			UseInterfaces:           true,
		},
	}
}

// GetConfigDir returns the path to the Wails config directory
func GetConfigDir() (string, error) {
	// Use XDG_CONFIG_HOME if set, otherwise use ~/.config
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configHome = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configHome, "wails"), nil
}

// GetDefaultsPath returns the path to the defaults.yaml file
func GetDefaultsPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "defaults.yaml"), nil
}

// Load loads the global defaults from the config file
// Returns default values if the file doesn't exist
func Load() (GlobalDefaults, error) {
	defaults := Default()

	path, err := GetDefaultsPath()
	if err != nil {
		return defaults, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaults, nil
		}
		return defaults, err
	}

	if err := yaml.Unmarshal(data, &defaults); err != nil {
		return Default(), err
	}

	return defaults, nil
}

// Save saves the global defaults to the config file
func Save(defaults GlobalDefaults) error {
	path, err := GetDefaultsPath()
	if err != nil {
		return err
	}

	// Ensure the config directory exists
	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(&defaults)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GenerateCopyright generates a copyright string from the template
func (d *GlobalDefaults) GenerateCopyright() string {
	template := d.Project.CopyrightTemplate
	if template == "" {
		template = "© {year}, {company}"
	}

	year := time.Now().Format("2006")
	company := d.Author.Company
	if company == "" {
		company = "My Company"
	}

	result := template
	result = replaceAll(result, "{year}", year)
	result = replaceAll(result, "{company}", company)
	return result
}

// GenerateProductIdentifier generates a product identifier from prefix and project name
func (d *GlobalDefaults) GenerateProductIdentifier(projectName string) string {
	prefix := d.Project.ProductIdentifierPrefix
	if prefix == "" {
		prefix = "com.example"
	}
	return prefix + "." + sanitizeIdentifier(projectName)
}

// GenerateDescription generates a description string from the template
func (d *GlobalDefaults) GenerateDescription(projectName string) string {
	template := d.Project.DescriptionTemplate
	if template == "" {
		template = "A {name} application"
	}
	return replaceAll(template, "{name}", projectName)
}

// GetDefaultVersion returns the default version or the fallback
func (d *GlobalDefaults) GetDefaultVersion() string {
	if d.Project.DefaultVersion != "" {
		return d.Project.DefaultVersion
	}
	return "0.1.0"
}

// replaceAll replaces all occurrences of old with new in s
func replaceAll(s, old, new string) string {
	result := s
	for {
		newResult := replaceOnce(result, old, new)
		if newResult == result {
			break
		}
		result = newResult
	}
	return result
}

func replaceOnce(s, old, new string) string {
	for i := 0; i <= len(s)-len(old); i++ {
		if s[i:i+len(old)] == old {
			return s[:i] + new + s[i+len(old):]
		}
	}
	return s
}

var (
	Frameworks = []string{
		"vanilla", "vue", "react", "react-swc", "svelte", "sveltekit",
		"preact", "lit", "solid", "qwik", "ios",
	}

	Languages = []string{"JavaScript", "TypeScript"}
)

func IsValidFramework(f string) bool {
	for _, fw := range Frameworks {
		if fw == f {
			return true
		}
	}
	return false
}

func (d *GlobalDefaults) GetTemplateName() string {
	framework := d.Project.Framework
	lang := d.Project.Language

	if framework != "" && IsValidFramework(framework) {
		if lang == "TypeScript" {
			return framework + "-ts"
		}
		return framework
	}

	if d.Project.DefaultTemplate != "" {
		return d.Project.DefaultTemplate
	}

	return "vanilla"
}

func sanitizeIdentifier(name string) string {
	var result []byte
	for i := 0; i < len(name); i++ {
		c := name[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			result = append(result, c)
		}
	}
	if len(result) == 0 {
		return "app"
	}
	// Lowercase the result
	for i := range result {
		if result[i] >= 'A' && result[i] <= 'Z' {
			result[i] = result[i] + 32
		}
	}
	return string(result)
}
