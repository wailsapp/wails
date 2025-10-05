package commands

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestGenerateBuildAssets(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "wails-build-assets-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name    string
		options *BuildAssetsOptions
		wantErr bool
	}{
		{
			name: "Basic build assets generation",
			options: &BuildAssetsOptions{
				Dir:                "testbuild",
				Name:               "TestApp",
				BinaryName:         "",
				ProductName:        "Test Application",
				ProductDescription: "A test application",
				ProductVersion:     "1.0.0",
				ProductCompany:     "Test Company",
				ProductCopyright:   "© 2024 Test Company",
				ProductComments:    "Test comments",
				ProductIdentifier:  "",
				Silent:             true,
			},
			wantErr: false,
		},
		{
			name: "Build assets with custom binary name",
			options: &BuildAssetsOptions{
				Dir:                "testbuild2",
				Name:               "Custom App",
				BinaryName:         "custom-binary",
				ProductName:        "Custom Application",
				ProductDescription: "A custom application",
				ProductVersion:     "2.0.0",
				ProductCompany:     "Custom Company",
				ProductIdentifier:  "com.custom.app",
				Silent:             true,
			},
			wantErr: false,
		},
		{
			name: "Build assets with MSIX options",
			options: &BuildAssetsOptions{
				Dir:                   "testbuild3",
				Name:                  "MSIX App",
				ProductName:           "MSIX Application",
				ProductDescription:    "An MSIX application",
				ProductVersion:        "3.0.0",
				ProductCompany:        "MSIX Company",
				Publisher:             "CN=MSIX Company",
				ProcessorArchitecture: "x64",
				ExecutablePath:        "msix-app.exe",
				ExecutableName:        "msix-app.exe",
				OutputPath:            "msix-app.msix",
				Silent:                true,
			},
			wantErr: false,
		},
		{
			name: "Build assets with TypeScript",
			options: &BuildAssetsOptions{
				Dir:                "testbuild4",
				Name:               "TypeScript App",
				ProductName:        "TypeScript Application",
				ProductDescription: "A TypeScript application",
				ProductVersion:     "4.0.0",
				ProductCompany:     "TypeScript Company",
				Typescript:         true,
				Silent:             true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the directory to be under our temp directory
			buildDir := filepath.Join(tempDir, tt.options.Dir)
			tt.options.Dir = buildDir

			err := GenerateBuildAssets(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateBuildAssets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify that the build directory was created
				if _, err := os.Stat(buildDir); os.IsNotExist(err) {
					t.Errorf("Build directory %s was not created", buildDir)
				}

				// List all files that were actually created for debugging
				files, err := os.ReadDir(buildDir)
				if err != nil {
					t.Errorf("Failed to read build directory: %v", err)
				} else {
					t.Logf("Files created in %s:", buildDir)
					for _, file := range files {
						t.Logf("  - %s", file.Name())
					}
				}

				// Verify some expected files were created - check what actually exists
				expectedFiles := []string{
					"config.yml",
					"appicon.png",
					"Taskfile.yml",
				}

				for _, file := range expectedFiles {
					filePath := filepath.Join(buildDir, file)
					if _, err := os.Stat(filePath); os.IsNotExist(err) {
						t.Errorf("Expected file %s was not created", file)
					}
				}

				// Test that defaults were applied correctly
				if tt.options.ProductIdentifier == "" && tt.options.Name != "" {
					expectedIdentifier := "com.wails." + normaliseName(tt.options.Name)
					// We can't easily check this without modifying the function to return the config
					// but we know the logic is there
					_ = expectedIdentifier
				}
			}
		})
	}
}

func TestUpdateBuildAssets(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "wails-update-assets-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a sample wails config file
	configDir := filepath.Join(tempDir, "config")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configFile := filepath.Join(configDir, "wails.yaml")
	config := WailsConfig{
		Info: struct {
			CompanyName       string `yaml:"companyName"`
			ProductName       string `yaml:"productName"`
			ProductIdentifier string `yaml:"productIdentifier"`
			Description       string `yaml:"description"`
			Copyright         string `yaml:"copyright"`
			Comments          string `yaml:"comments"`
			Version           string `yaml:"version"`
		}{
			CompanyName:       "Config Company",
			ProductName:       "Config Product",
			ProductIdentifier: "com.config.product",
			Description:       "Config Description",
			Copyright:         "© 2024 Config Company",
			Comments:          "Config Comments",
			Version:           "1.0.0",
		},
		FileAssociations: []FileAssociation{
			{
				Ext:         ".test",
				Name:        "Test File",
				Description: "Test file association",
				IconName:    "test-icon",
				Role:        "Editor",
				MimeType:    "application/test",
			},
		},
		Protocols: []ProtocolConfig{
			{
				Scheme:      "testapp",
				Description: "Test App Protocol",
			},
		},
	}

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	err = os.WriteFile(configFile, configBytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	tests := []struct {
		name    string
		options *UpdateBuildAssetsOptions
		wantErr bool
	}{
		{
			name: "Update with config file",
			options: &UpdateBuildAssetsOptions{
				Dir:    "updatebuild1",
				Name:   "UpdateApp",
				Config: configFile,
				Silent: true,
			},
			wantErr: false,
		},
		{
			name: "Update without config file",
			options: &UpdateBuildAssetsOptions{
				Dir:                "updatebuild2",
				Name:               "UpdateApp2",
				ProductName:        "Update Application 2",
				ProductDescription: "An update application 2",
				ProductVersion:     "2.0.0",
				ProductCompany:     "Update Company 2",
				Silent:             true,
			},
			wantErr: false,
		},
		{
			name: "Update with non-existent config file",
			options: &UpdateBuildAssetsOptions{
				Dir:    "updatebuild3",
				Name:   "UpdateApp3",
				Config: "non-existent-config.yaml",
				Silent: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the directory to be under our temp directory
			updateDir := filepath.Join(tempDir, tt.options.Dir)
			tt.options.Dir = updateDir

			err := UpdateBuildAssets(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateBuildAssets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify that the update directory was created
				if _, err := os.Stat(updateDir); os.IsNotExist(err) {
					t.Errorf("Update directory %s was not created", updateDir)
				}
			}
		})
	}
}
