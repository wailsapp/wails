package commands

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/leaanthony/gosod"
	"gopkg.in/yaml.v3"
	"howett.net/plist"
)

//go:embed build_assets
var buildAssets embed.FS

//go:embed updatable_build_assets
var updatableBuildAssets embed.FS

// ProtocolConfig defines the structure for a custom protocol in wails.json/wails.yaml
type ProtocolConfig struct {
	Scheme      string `yaml:"scheme"                json:"scheme"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// Future platform-specific fields can be added here if needed by templates.
	// E.g., for macOS: CFBundleURLName string `yaml:"cfBundleURLName,omitempty" json:"cfBundleURLName,omitempty"`
}

// BuildAssetsOptions defines the options for generating build assets.
type BuildAssetsOptions struct {
	Dir                   string `description:"The directory to generate the files into"            default:"."`
	Name                  string `description:"The name of the project"`
	BinaryName            string `description:"The name of the binary"`
	ProductName           string `description:"The name of the product" default:"My Product"`
	ProductDescription    string `description:"The description of the product" default:"My Product Description"`
	ProductVersion        string `description:"The version of the product" default:"0.1.0"`
	ProductCompany        string `description:"The company of the product" default:"My Company"`
	ProductCopyright      string `description:"The copyright notice" default:"\u00a9 now, My Company"`
	ProductComments       string `description:"Comments to add to the generated files" default:"This is a comment"`
	ProductIdentifier     string `description:"The product identifier, e.g com.mycompany.myproduct"`
	Publisher             string `description:"Publisher name for MSIX package (e.g., CN=CompanyName)"`
	ProcessorArchitecture string `description:"Processor architecture for MSIX package" default:"x64"`
	ExecutablePath        string `description:"Path to executable for MSIX package"`
	ExecutableName        string `description:"Name of executable for MSIX package"`
	OutputPath            string `description:"Output path for MSIX package"`
	CertificatePath       string `description:"Certificate path for MSIX package"`
	Silent                bool   `description:"Suppress output to console"`
	Typescript            bool   `description:"Use typescript" default:"false"`
}

// BuildConfig defines the configuration for generating build assets.
type BuildConfig struct {
	BuildAssetsOptions
	FileAssociations []FileAssociation `yaml:"fileAssociations"`
	Protocols        []ProtocolConfig  `yaml:"protocols,omitempty"`
}

// UpdateBuildAssetsOptions defines the options for updating build assets.
type UpdateBuildAssetsOptions struct {
	Dir                string `description:"The directory to generate the files into"            default:"build"`
	Name               string `description:"The name of the project"`
	BinaryName         string `description:"The name of the binary"`
	ProductName        string `description:"The name of the product"                             default:"My Product"`
	ProductDescription string `description:"The description of the product"                      default:"My Product Description"`
	ProductVersion     string `description:"The version of the product"                          default:"0.1.0"`
	ProductCompany     string `description:"The company of the product"                          default:"My Company"`
	ProductCopyright   string `description:"The copyright notice"                                default:"© now, My Company"`
	ProductComments    string `description:"Comments to add to the generated files"              default:"This is a comment"`
	ProductIdentifier  string `description:"The product identifier, e.g com.mycompany.myproduct"`
	Config             string `description:"The path to the config file"`
	Silent             bool   `description:"Suppress output to console"`
}

// GenerateBuildAssets generates the build assets for the project.
func GenerateBuildAssets(options *BuildAssetsOptions) error {
	DisableFooter = true

	var err error
	options.Dir, err = filepath.Abs(options.Dir)
	if err != nil {
		return err
	}

	// If directory doesn't exist, create it
	if _, err := os.Stat(options.Dir); os.IsNotExist(err) {
		err = os.MkdirAll(options.Dir, 0755)
		if err != nil {
			return err
		}
	}

	var config BuildConfig

	if options.ProductComments == "" {
		options.ProductComments = fmt.Sprintf("(c) %d %s", time.Now().Year(), options.ProductCompany)
	}

	if options.ProductIdentifier == "" {
		options.ProductIdentifier = "com.wails." + normaliseName(options.Name)
	}

	if options.BinaryName == "" {
		options.BinaryName = normaliseName(options.Name)
		if runtime.GOOS == "windows" {
			options.BinaryName += ".exe"
		}
	}

	if options.Publisher == "" {
		options.Publisher = fmt.Sprintf("CN=%s", options.ProductCompany)
	}

	if options.ProcessorArchitecture == "" {
		options.ProcessorArchitecture = "x64"
	}

	if options.ExecutableName == "" {
		options.ExecutableName = options.BinaryName
	}

	if options.ExecutablePath == "" {
		options.ExecutablePath = options.BinaryName
	}

	if options.OutputPath == "" {
		options.OutputPath = fmt.Sprintf("%s.msix", normaliseName(options.Name))
	}

	// CertificatePath is optional, no default needed

	config.BuildAssetsOptions = *options

	tfs, err := fs.Sub(buildAssets, "build_assets")
	if err != nil {
		return err
	}

	if !options.Silent {
		println("Generating build assets in " + options.Dir)
	}
	err = gosod.New(tfs).Extract(options.Dir, config)
	if err != nil {
		return err
	}
	tfs, err = fs.Sub(updatableBuildAssets, "updatable_build_assets")
	if err != nil {
		return err
	}
	err = gosod.New(tfs).Extract(options.Dir, config)
	if err != nil {
		return err
	}

	return nil
}

// FileAssociation defines the structure for a file association.
type FileAssociation struct {
	Ext         string `yaml:"ext"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	IconName    string `yaml:"iconName"`
	Role        string `yaml:"role"`
	MimeType    string `yaml:"mimeType"`
}

// UpdateConfig defines the configuration for updating build assets.
type UpdateConfig struct {
	UpdateBuildAssetsOptions
	FileAssociations []FileAssociation `yaml:"fileAssociations"`
	Protocols        []ProtocolConfig  `yaml:"protocols,omitempty"`
}

// WailsConfig defines the structure for a Wails configuration.
type WailsConfig struct {
	Info struct {
		CompanyName       string `yaml:"companyName"`
		ProductName       string `yaml:"productName"`
		ProductIdentifier string `yaml:"productIdentifier"`
		Description       string `yaml:"description"`
		Copyright         string `yaml:"copyright"`
		Comments          string `yaml:"comments"`
		Version           string `yaml:"version"`
	} `yaml:"info"`
	FileAssociations []FileAssociation `yaml:"fileAssociations,omitempty"`
	Protocols        []ProtocolConfig  `yaml:"protocols,omitempty"`
}

// UpdateBuildAssets updates the build assets for the project.
func UpdateBuildAssets(options *UpdateBuildAssetsOptions) error {
	DisableFooter = true

	var err error
	options.Dir, err = filepath.Abs(options.Dir)
	if err != nil {
		return err
	}

	var config UpdateConfig
	if options.Config != "" {
		var wailsConfig WailsConfig
		bytes, err := os.ReadFile(options.Config)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(bytes, &wailsConfig)
		if err != nil {
			return err
		}

		if options.ProductCompany == "My Company" && wailsConfig.Info.CompanyName != "" {
			options.ProductCompany = wailsConfig.Info.CompanyName
		}
		if options.ProductName == "My Product" && wailsConfig.Info.ProductName != "" {
			options.ProductName = wailsConfig.Info.ProductName
		}
		if options.ProductIdentifier == "" {
			options.ProductIdentifier = wailsConfig.Info.ProductIdentifier
		}
		if options.ProductDescription == "My Product Description" && wailsConfig.Info.Description != "" {
			options.ProductDescription = wailsConfig.Info.Description
		}
		if options.ProductCopyright == "© now, My Company" && wailsConfig.Info.Copyright != "" {
			options.ProductCopyright = wailsConfig.Info.Copyright
		}
		if options.ProductComments == "This is a comment" && wailsConfig.Info.Comments != "" {
			options.ProductComments = wailsConfig.Info.Comments
		}
		if options.ProductVersion == "0.1.0" && wailsConfig.Info.Version != "" {
			options.ProductVersion = wailsConfig.Info.Version
		}
		config.FileAssociations = wailsConfig.FileAssociations
		config.Protocols = wailsConfig.Protocols
	}

	config.UpdateBuildAssetsOptions = *options

	// If directory doesn't exist, create it
	if _, err := os.Stat(options.Dir); os.IsNotExist(err) {
		err = os.MkdirAll(options.Dir, 0755)
		if err != nil {
			return err
		}
	}

	tfs, err := fs.Sub(updatableBuildAssets, "updatable_build_assets")
	if err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", "wails-build-assets-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	err = gosod.New(tfs).Extract(tempDir, config)
	if err != nil {
		return err
	}

	err = mergePlistFiles(tempDir, options.Dir)
	if err != nil {
		return err
	}

	err = copyNonPlistFiles(tempDir, options.Dir)
	if err != nil {
		return err
	}

	if !options.Silent {
		println("Successfully updated build assets in " + options.Dir)
	}

	return nil
}

func normaliseName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}

// mergeMaps recursively merges src into dst.
// For nested maps, it merges recursively. For other types, src overwrites dst.
func mergeMaps(dst, src map[string]any) {
	for key, srcValue := range src {
		if dstValue, exists := dst[key]; exists {
			// If both are maps, merge recursively
			srcMap, srcIsMap := srcValue.(map[string]any)
			dstMap, dstIsMap := dstValue.(map[string]any)
			if srcIsMap && dstIsMap {
				mergeMaps(dstMap, srcMap)
				continue
			}
		}
		// Otherwise, src overwrites dst
		dst[key] = srcValue
	}
}

func mergePlistFile(filePath string, newContent []byte) error {
	var newDict map[string]any
	_, err := plist.Unmarshal(newContent, &newDict)
	if err != nil {
		return fmt.Errorf("failed to parse new plist content: %w", err)
	}

	existingDict := make(map[string]any)
	if _, err := os.Stat(filePath); err == nil {
		existingContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read existing plist file: %w", err)
		}
		_, err = plist.Unmarshal(existingContent, &existingDict)
		if err != nil {
			return fmt.Errorf("failed to parse existing plist file: %w", err)
		}
	}

	mergeMaps(existingDict, newDict)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create plist file: %w", err)
	}
	defer file.Close()

	encoder := plist.NewEncoder(file)
	encoder.Indent("\t")
	err = encoder.Encode(existingDict)
	if err != nil {
		return fmt.Errorf("failed to encode merged plist: %w", err)
	}

	return nil
}

func mergePlistFiles(tempDir, targetDir string) error {
	return filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".plist") {
			return nil
		}

		relPath, err := filepath.Rel(tempDir, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(targetDir, relPath)

		newContent, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read new plist file %s: %w", path, err)
		}

		targetDirPath := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create target directory %s: %w", targetDirPath, err)
		}

		if err := mergePlistFile(targetPath, newContent); err != nil {
			return fmt.Errorf("failed to merge plist file %s: %w", targetPath, err)
		}

		return nil
	})
}

func copyNonPlistFiles(tempDir, targetDir string) error {
	return filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".plist") {
			return nil
		}

		relPath, err := filepath.Rel(tempDir, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(targetDir, relPath)

		targetDirPath := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create target directory %s: %w", targetDirPath, err)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		if err := os.WriteFile(targetPath, content, info.Mode()); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}

		return nil
	})
}
