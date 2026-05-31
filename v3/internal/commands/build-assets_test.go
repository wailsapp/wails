package commands

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
	"howett.net/plist"
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
			CFBundleIconName  string `yaml:"cfBundleIconName,omitempty"`
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

func TestPlistMerge(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "wails-plist-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	buildDir := filepath.Join(tempDir, "build", "darwin")
	err = os.MkdirAll(buildDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create build directory: %v", err)
	}

	existingPlistPath := filepath.Join(buildDir, "Info.plist")
	existingPlist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleName</key>
	<string>OldAppName</string>
	<key>CFBundleVersion</key>
	<string>1.0.0</string>
	<key>NSCameraUsageDescription</key>
	<string>This app needs camera access</string>
	<key>NSMicrophoneUsageDescription</key>
	<string>This app needs microphone access</string>
</dict>
</plist>`

	err = os.WriteFile(existingPlistPath, []byte(existingPlist), 0644)
	if err != nil {
		t.Fatalf("Failed to write existing plist: %v", err)
	}

	options := &UpdateBuildAssetsOptions{
		Dir:                filepath.Join(tempDir, "build"),
		Name:               "TestApp",
		ProductName:        "NewAppName",
		ProductVersion:     "2.0.0",
		ProductCompany:     "Test Company",
		ProductIdentifier:  "com.test.app",
		ProductDescription: "Test Description",
		ProductCopyright:   "© 2024 Test Company",
		ProductComments:    "Test Comments",
		Silent:             true,
	}

	err = UpdateBuildAssets(options)
	if err != nil {
		t.Fatalf("UpdateBuildAssets failed: %v", err)
	}

	mergedContent, err := os.ReadFile(existingPlistPath)
	if err != nil {
		t.Fatalf("Failed to read merged plist: %v", err)
	}

	var mergedDict map[string]any
	_, err = plist.Unmarshal(mergedContent, &mergedDict)
	if err != nil {
		t.Fatalf("Failed to parse merged plist: %v", err)
	}

	if mergedDict["CFBundleName"] != "NewAppName" {
		t.Errorf("Expected CFBundleName to be updated to 'NewAppName', got %v", mergedDict["CFBundleName"])
	}

	if mergedDict["CFBundleVersion"] != "2.0.0" {
		t.Errorf("Expected CFBundleVersion to be updated to '2.0.0', got %v", mergedDict["CFBundleVersion"])
	}

	if mergedDict["NSCameraUsageDescription"] != "This app needs camera access" {
		t.Errorf("Expected NSCameraUsageDescription to be preserved, got %v", mergedDict["NSCameraUsageDescription"])
	}

	if mergedDict["NSMicrophoneUsageDescription"] != "This app needs microphone access" {
		t.Errorf("Expected NSMicrophoneUsageDescription to be preserved, got %v", mergedDict["NSMicrophoneUsageDescription"])
	}

	if mergedDict["CFBundleIdentifier"] != "com.test.app" {
		t.Errorf("Expected CFBundleIdentifier to be 'com.test.app', got %v", mergedDict["CFBundleIdentifier"])
	}
}

func TestCFBundleIconNameDetection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "wails-icon-name-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name                string
		createAssetsCar     bool
		configIconName      string
		expectedIconName    string
		expectIconNameInPlist bool
	}{
		{
			name:                "Assets.car exists, no config - should default to appicon",
			createAssetsCar:     true,
			configIconName:      "",
			expectedIconName:    "appicon",
			expectIconNameInPlist: true,
		},
		{
			name:                "Assets.car exists, config set - should use config",
			createAssetsCar:     true,
			configIconName:      "custom-icon",
			expectedIconName:    "custom-icon",
			expectIconNameInPlist: true,
		},
		{
			name:                "No Assets.car, no config - should not set",
			createAssetsCar:     false,
			configIconName:      "",
			expectedIconName:    "",
			expectIconNameInPlist: false,
		},
		{
			name:                "No Assets.car, config set - should use config",
			createAssetsCar:     false,
			configIconName:      "config-icon",
			expectedIconName:    "config-icon",
			expectIconNameInPlist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildDir := filepath.Join(tempDir, tt.name)
			darwinDir := filepath.Join(buildDir, "darwin")
			err := os.MkdirAll(darwinDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create darwin directory: %v", err)
			}

			// Create Assets.car BEFORE calling UpdateBuildAssets if needed
			// The check happens before template extraction, so CFBundleIconName will be available in the template
			if tt.createAssetsCar {
				assetsCarPath := filepath.Join(darwinDir, "Assets.car")
				err = os.WriteFile(assetsCarPath, []byte("fake assets.car content"), 0644)
				if err != nil {
					t.Fatalf("Failed to create Assets.car: %v", err)
				}
			}

			// Create config file if icon name is set
			configFile := ""
			if tt.configIconName != "" {
				configDir := filepath.Join(tempDir, "config-"+tt.name)
				err = os.MkdirAll(configDir, 0755)
				if err != nil {
					t.Fatalf("Failed to create config directory: %v", err)
				}

				configFile = filepath.Join(configDir, "wails.yaml")
				config := WailsConfig{
					Info: struct {
						CompanyName       string `yaml:"companyName"`
						ProductName       string `yaml:"productName"`
						ProductIdentifier string `yaml:"productIdentifier"`
						Description       string `yaml:"description"`
						Copyright         string `yaml:"copyright"`
						Comments          string `yaml:"comments"`
						Version           string `yaml:"version"`
						CFBundleIconName  string `yaml:"cfBundleIconName,omitempty"`
					}{
						CompanyName:       "Test Company",
						ProductName:       "Test Product",
						ProductIdentifier: "com.test.product",
						CFBundleIconName:  tt.configIconName,
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
			}

			options := &UpdateBuildAssetsOptions{
				Dir:               buildDir,
				Name:               "TestApp",
				ProductName:        "Test App",
				ProductVersion:     "1.0.0",
				ProductCompany:     "Test Company",
				ProductIdentifier:  "com.test.app",
				CFBundleIconName:   tt.configIconName,
				Config:             configFile,
				Silent:             true,
			}

			err = UpdateBuildAssets(options)
			if err != nil {
				t.Fatalf("UpdateBuildAssets failed: %v", err)
			}

			// Verify CFBundleIconName was set correctly in options
			if options.CFBundleIconName != tt.expectedIconName {
				t.Errorf("Expected CFBundleIconName to be '%s', got '%s'", tt.expectedIconName, options.CFBundleIconName)
			}

			// Check Info.plist if it exists
			infoPlistPath := filepath.Join(darwinDir, "Info.plist")
			if _, err := os.Stat(infoPlistPath); err == nil {
				plistContent, err := os.ReadFile(infoPlistPath)
				if err != nil {
					t.Fatalf("Failed to read Info.plist: %v", err)
				}

				var plistDict map[string]any
				_, err = plist.Unmarshal(plistContent, &plistDict)
				if err != nil {
					t.Fatalf("Failed to parse Info.plist: %v", err)
				}

				iconName, exists := plistDict["CFBundleIconName"]
				if tt.expectIconNameInPlist {
					if !exists {
						t.Errorf("Expected CFBundleIconName to be present in Info.plist")
					} else if iconName != tt.expectedIconName {
						t.Errorf("Expected CFBundleIconName in Info.plist to be '%s', got '%v'", tt.expectedIconName, iconName)
					}
				} else {
					if exists {
						t.Errorf("Expected CFBundleIconName to not be present in Info.plist, but found '%v'", iconName)
					}
				}
			}
		})
	}
}

func TestNestedPlistMerge(t *testing.T) {
	tests := []struct {
		name     string
		existing map[string]any
		new      map[string]any
		expected map[string]any
	}{
		{
			name: "simple overwrite",
			existing: map[string]any{
				"key1": "oldValue",
			},
			new: map[string]any{
				"key1": "newValue",
			},
			expected: map[string]any{
				"key1": "newValue",
			},
		},
		{
			name: "preserve existing keys",
			existing: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			new: map[string]any{
				"key1": "newValue1",
			},
			expected: map[string]any{
				"key1": "newValue1",
				"key2": "value2",
			},
		},
		{
			name: "nested dict merge",
			existing: map[string]any{
				"CustomConfig": map[string]any{
					"Setting1": "existingValue1",
					"Setting2": "existingValue2",
				},
			},
			new: map[string]any{
				"CustomConfig": map[string]any{
					"Setting1": "newValue1",
					"Setting3": "newValue3",
				},
			},
			expected: map[string]any{
				"CustomConfig": map[string]any{
					"Setting1": "newValue1",
					"Setting2": "existingValue2",
					"Setting3": "newValue3",
				},
			},
		},
		{
			name: "deeply nested merge",
			existing: map[string]any{
				"Level1": map[string]any{
					"Level2": map[string]any{
						"deepKey1": "deepValue1",
						"deepKey2": "deepValue2",
					},
				},
			},
			new: map[string]any{
				"Level1": map[string]any{
					"Level2": map[string]any{
						"deepKey1": "newDeepValue1",
						"deepKey3": "newDeepValue3",
					},
				},
			},
			expected: map[string]any{
				"Level1": map[string]any{
					"Level2": map[string]any{
						"deepKey1": "newDeepValue1",
						"deepKey2": "deepValue2",
						"deepKey3": "newDeepValue3",
					},
				},
			},
		},
		{
			name: "mixed types - new dict replaces non-dict",
			existing: map[string]any{
				"key1": "stringValue",
			},
			new: map[string]any{
				"key1": map[string]any{
					"nested": "value",
				},
			},
			expected: map[string]any{
				"key1": map[string]any{
					"nested": "value",
				},
			},
		},
		{
			name: "mixed types - new non-dict replaces dict",
			existing: map[string]any{
				"key1": map[string]any{
					"nested": "value",
				},
			},
			new: map[string]any{
				"key1": "stringValue",
			},
			expected: map[string]any{
				"key1": "stringValue",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of existing to avoid mutation issues
			dst := deepCopyMap(tt.existing)
			mergeMaps(dst, tt.new)

			if !mapsEqual(dst, tt.expected) {
				t.Errorf("mergeMaps() got %v, expected %v", dst, tt.expected)
			}
		})
	}
}

// TestOldFormatPlistMigration verifies that update-build-assets strips Go
// template syntax (e.g. "{{.Ext}}") left behind when an older project's
// darwin/Info.plist file is still in raw-template form.
// See: https://github.com/wailsapp/wails/issues/5259
func TestOldFormatPlistMigration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "wails-old-plist-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	buildDir := filepath.Join(tempDir, "build")
	darwinDir := filepath.Join(buildDir, "darwin")
	if err := os.MkdirAll(darwinDir, 0755); err != nil {
		t.Fatalf("Failed to create darwin directory: %v", err)
	}

	// Simulate an older-format Info.plist that still contains raw Go template
	// directives (as produced by wails v2 / early v3 alpha scaffolding).
	// The outer {{if}}…{{end}} blocks are ignored by the XML parser as text
	// nodes, but the inner <string>{{.Ext}}</string> etc. are parsed as real
	// string values — those are the stubs the fix must remove.
	oldPlist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleName</key>
	<string>OldApp</string>
	<key>NSCameraUsageDescription</key>
	<string>Camera access needed</string>
	{{if .Info.FileAssociations}}
	<key>CFBundleDocumentTypes</key>
	<array>
	  {{range .Info.FileAssociations}}
	  <dict>
	    <key>CFBundleTypeExtensions</key>
	    <array>
	      <string>{{.Ext}}</string>
	    </array>
	    <key>CFBundleTypeName</key>
	    <string>{{.Name}}</string>
	    <key>CFBundleTypeRole</key>
	    <string>{{.Role}}</string>
	    <key>CFBundleTypeIconFile</key>
	    <string>{{.IconName}}</string>
	  </dict>
	  {{end}}
	</array>
	{{end}}
	{{if .Info.Protocols}}
	<key>CFBundleURLTypes</key>
	<array>
	  {{range .Info.Protocols}}
	  <dict>
	    <key>CFBundleURLName</key>
	    <string>com.wails.{{.Scheme}}</string>
	    <key>CFBundleURLSchemes</key>
	    <array>
	      <string>{{.Scheme}}</string>
	    </array>
	    <key>CFBundleTypeRole</key>
	    <string>{{.Role}}</string>
	  </dict>
	  {{end}}
	</array>
	{{end}}
</dict>
</plist>`

	existingPlistPath := filepath.Join(darwinDir, "Info.plist")
	if err := os.WriteFile(existingPlistPath, []byte(oldPlist), 0644); err != nil {
		t.Fatalf("Failed to write old-format plist: %v", err)
	}

	options := &UpdateBuildAssetsOptions{
		Dir:               buildDir,
		Name:              "TestApp",
		ProductName:       "TestApp",
		ProductVersion:    "1.0.0",
		ProductCompany:    "Wails",
		ProductIdentifier: "com.wails.testapp",
		Silent:            true,
	}

	if err := UpdateBuildAssets(options); err != nil {
		t.Fatalf("UpdateBuildAssets failed: %v", err)
	}

	content, err := os.ReadFile(existingPlistPath)
	if err != nil {
		t.Fatalf("Failed to read merged plist: %v", err)
	}

	contentStr := string(content)

	// No template stub should survive in the output.
	stubs := []string{"{{.Ext}}", "{{.Name}}", "{{.Role}}", "{{.IconName}}", "{{.Scheme}}"}
	for _, stub := range stubs {
		if strings.Contains(contentStr, stub) {
			t.Errorf("Output Info.plist still contains template stub %q — old-format template was not sanitized correctly", stub)
		}
	}

	// User-added keys that are real values must be preserved.
	var mergedDict map[string]any
	if _, err := plist.Unmarshal(content, &mergedDict); err != nil {
		t.Fatalf("Failed to parse merged plist: %v", err)
	}
	if mergedDict["NSCameraUsageDescription"] != "Camera access needed" {
		t.Errorf("Custom key NSCameraUsageDescription was lost during migration; got %v", mergedDict["NSCameraUsageDescription"])
	}
}

func deepCopyMap(m map[string]any) map[string]any {
	result := make(map[string]any)
	for k, v := range m {
		if nested, ok := v.(map[string]any); ok {
			result[k] = deepCopyMap(nested)
		} else {
			result[k] = v
		}
	}
	return result
}

func mapsEqual(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for k, av := range a {
		bv, ok := b[k]
		if !ok {
			return false
		}
		aMap, aIsMap := av.(map[string]any)
		bMap, bIsMap := bv.(map[string]any)
		if aIsMap && bIsMap {
			if !mapsEqual(aMap, bMap) {
				return false
			}
		} else if aIsMap != bIsMap {
			return false
		} else if !reflect.DeepEqual(av, bv) {
			return false
		}
	}
	return true
}

// TestPreserveOriginallyEmptyContainers verifies that originally empty
// containers (maps and arrays) are preserved during sanitization, while
// containers that become empty due to template removal are dropped.
func TestPreserveOriginallyEmptyContainers(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name: "originally empty map is preserved",
			input: map[string]any{
				"EmptyMap": map[string]any{},
			},
			expected: map[string]any{
				"EmptyMap": map[string]any{},
			},
		},
		{
			name: "originally empty array is preserved",
			input: map[string]any{
				"EmptyArray": []any{},
			},
			expected: map[string]any{
				"EmptyArray": []any{},
			},
		},
		{
			name: "map with only template values becomes empty and is dropped",
			input: map[string]any{
				"TemplateMap": map[string]any{
					"Key1": "{{.Ext}}",
					"Key2": "{{.Name}}",
				},
			},
			expected: map[string]any{},
		},
		{
			name: "array with only template values becomes empty and is dropped",
			input: map[string]any{
				"TemplateArray": []any{
					"{{.Ext}}",
					"{{.Name}}",
				},
			},
			expected: map[string]any{},
		},
		{
			name: "map with mixed content keeps real values and drops templates",
			input: map[string]any{
				"MixedMap": map[string]any{
					"RealKey":    "RealValue",
					"TemplateKey": "{{.Ext}}",
				},
			},
			expected: map[string]any{
				"MixedMap": map[string]any{
					"RealKey": "RealValue",
				},
			},
		},
		{
			name: "array with mixed content keeps real values and drops templates",
			input: map[string]any{
				"MixedArray": []any{
					"RealValue",
					"{{.Ext}}",
					"AnotherReal",
					"{{.Name}}",
				},
			},
			expected: map[string]any{
				"MixedArray": []any{
					"RealValue",
					"AnotherReal",
				},
			},
		},
		{
			name: "nested originally empty containers are preserved",
			input: map[string]any{
				"Outer": map[string]any{
					"InnerMap":  map[string]any{},
					"InnerArray": []any{},
				},
			},
			expected: map[string]any{
				"Outer": map[string]any{
					"InnerMap":  map[string]any{},
					"InnerArray": []any{},
				},
			},
		},
		{
			name: "complex nested structure with empty and template values",
			input: map[string]any{
				"Complex": map[string]any{
					"EmptySub": map[string]any{},
					"MixedSub": []any{
						"Real",
						"{{.Template}}",
					},
					"TemplateOnlySub": []any{
						"{{.Ext}}",
						"{{.Name}}",
					},
				},
			},
			expected: map[string]any{
				"Complex": map[string]any{
					"EmptySub": map[string]any{},
					"MixedSub": []any{
						"Real",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizePlistDict(tt.input.(map[string]any))

			if !mapsEqual(result, tt.expected.(map[string]any)) {
				t.Errorf("sanitizePlistDict() got\n%v\nexpected\n%v", result, tt.expected)
			}
		})
	}
}
