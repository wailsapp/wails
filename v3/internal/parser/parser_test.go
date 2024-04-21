package parser

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
)

func TestPackageDir(t *testing.T) {

	tests := []struct {
		project     string
		basePath    string
		useBaseName bool
		want        map[string]bool
	}{
		{
			project:     "function_from_imported_package",
			basePath:    ".",
			useBaseName: false,
			want: map[string]bool{
				".":        true,
				"services": true,
			},
		},
		{
			project:     "function_from_imported_package",
			basePath:    ".",
			useBaseName: true,
			want: map[string]bool{
				"main":     true,
				"services": true,
			},
		},
		{
			project:     "function_from_imported_package",
			basePath:    "github.com/wailsapp/wails/v3/internal/parser/testdata",
			useBaseName: false,
			want: map[string]bool{
				"function_from_imported_package":          true,
				"function_from_imported_package/services": true,
			},
		},
		{
			project:     "function_from_imported_package",
			basePath:    "github.com/wailsapp/wails/v3/internal/parser/testdata",
			useBaseName: true, // will be ignored
			want: map[string]bool{
				"function_from_imported_package":          true,
				"function_from_imported_package/services": true,
			},
		},
		{
			project:     "function_from_imported_package",
			basePath:    "",
			useBaseName: false,
			want: map[string]bool{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_imported_package":          true,
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_imported_package/services": true,
			},
		},
		{
			project:     "app_outside_main/app",
			basePath:    ".",
			useBaseName: false,
			want: map[string]bool{
				".":      true,
				"models": true,
			},
		},
		{
			project:     "app_outside_main/app",
			basePath:    ".",
			useBaseName: true,
			want: map[string]bool{
				"app":    true,
				"models": true,
			},
		},
		{
			project:     "app_outside_main/app",
			basePath:    "github.com/wailsapp/wails/v3/internal/parser/testdata/app_outside_main",
			useBaseName: false,
			want: map[string]bool{
				"app":        true,
				"app/models": true,
			},
		},
		{
			project:     "multiple_packages",
			basePath:    ".",
			useBaseName: true,
			want: map[string]bool{
				"other":                  true,
				"other/other":            true,
				"runtime/debug":          true,
				"github.com/google/uuid": true,
			},
		},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%s__%s_%v", tt.project, tt.basePath, tt.useBaseName)

		t.Run(name, func(t *testing.T) {

			absDir, err := filepath.Abs("testdata/" + tt.project)
			if err != nil {
				t.Errorf("filepath.Abs() error = %v", err)
				return
			}

			options := &flags.GenerateBindingsOptions{
				ProjectDirectory: absDir,
				BasePath:         tt.basePath,
				UseBaseName:      tt.useBaseName,
			}

			project, err := ParseProjectAndPkgs(options)
			if err != nil {
				t.Errorf("ParseProjectAndPkgs() error = %v", err)
				return
			}

			// Generate Models
			allModels, err := project.GenerateModels()
			if err != nil {
				t.Fatalf("GenerateModels() error = %v", err)
			}

			// Check if models are missing
			for pkgDir := range tt.want {
				if _, ok := allModels[pkgDir]; !ok {
					t.Errorf("GenerateModels() missing model = %v", pkgDir)
				}
			}

			// Check for unexpected models
			for pkgDir := range allModels {
				if _, ok := tt.want[pkgDir]; !ok {
					t.Errorf("GenerateModels() unexpected model = %v", pkgDir)
				}
			}

		})
	}

}
