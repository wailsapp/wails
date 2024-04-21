package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/wailsapp/wails/v3/internal/flags"
)

func TestGenerateModels(t *testing.T) {

	options := []struct {
		name          string
		useInterface  bool
		useTypescript bool
		ext           string
	}{
		{
			name:          "Typescript",
			useInterface:  false,
			useTypescript: true,
			ext:           "ts",
		},
		{
			name:          "Javascript",
			useInterface:  false,
			useTypescript: false,
			ext:           "js",
		},
		{
			name:          "Typescript interfaces",
			useInterface:  true,
			useTypescript: true,
			ext:           "interfaces.ts",
		},
	}

	tests := []struct {
		name string
		dir  string
		want map[string]bool
	}{
		{
			name: "complex_json",
			dir:  "testdata/complex_json",
			want: map[string]bool{
				"main": true,
			},
		},
		{
			name: "complex_method",
			dir:  "testdata/complex_method",
			want: map[string]bool{
				"main": true,
			},
		},
		{
			name: "enum",
			dir:  "testdata/enum",
			want: map[string]bool{
				"main": true,
			},
		},
		{
			name: "function_from_imported_package",
			dir:  "testdata/function_from_imported_package",
			want: map[string]bool{
				"main":     true,
				"services": true,
			},
		},
		{
			name: "function_from_nested_imported_package",
			dir:  "testdata/function_from_nested_imported_package",
			want: map[string]bool{
				"main":           true,
				"services/other": true,
			},
		},
		{
			name: "variable_single_from_other_function",
			dir:  "testdata/variable_single_from_other_function",
			want: map[string]bool{
				"main":     true,
				"services": true,
			},
		},
		{
			name: "struct_literal_single",
			dir:  "testdata/struct_literal_single",
			want: map[string]bool{
				"main": true,
			},
		},

		{
			name: "struct_literal_multiple_other",
			dir:  "testdata/struct_literal_multiple_other",
			want: map[string]bool{
				"main":     true,
				"services": true,
			},
		},
		{
			name: "struct_literal_non_pointer_single",
			dir:  "testdata/struct_literal_non_pointer_single",
			want: map[string]bool{
				"main": true,
			},
		},

		{
			name: "enum_from_imported_package",
			dir:  "testdata/enum_from_imported_package",
			want: map[string]bool{
				"services": true,
			},
		},
		{
			name: "nested_types",
			dir:  "testdata/nested_types",
			want: map[string]bool{
				"main": true,
			},
		},
		{
			name: "multiple_packages",
			dir:  "testdata/multiple_packages",
			want: map[string]bool{
				"github.com/google/uuid": true,
				"runtime/debug":          true,
				"other":                  true,
				"other/other":            true,
			},
		},
		{
			name: "type_alias",
			dir:  "testdata/type_alias",
			want: map[string]bool{
				"main": true,
			},
		},
		{
			name: "interfaces",
			dir:  "testdata/interfaces",
			want: map[string]bool{
				"main": true,
			},
		},
		{
			name: "app_outside_main",
			dir:  "testdata/app_outside_main/app",
			want: map[string]bool{
				"app":    true,
				"models": true,
			},
		},
	}

	type Test struct {
		name          string
		dir           string
		want          map[string]string
		useTypescript bool
		useInterface  bool
	}

	allTests := []Test{}
	for _, tt := range tests {
		for _, option := range options {
			want := make(map[string]string)

			for pkgDir := range tt.want {
				want[pkgDir] = getFile(fmt.Sprintf("%s/frontend/bindings/%s/models.%s", tt.dir, pkgDir, option.ext))
			}

			allTests = append(allTests, Test{
				name:          tt.name + " - " + option.name,
				dir:           tt.dir,
				want:          want,
				useTypescript: option.useTypescript,
				useInterface:  option.useInterface,
			})
		}

	}

	for _, tt := range allTests {
		t.Run(tt.name, func(t *testing.T) {
			// Run parser on directory
			absDir, err := filepath.Abs(tt.dir)
			if err != nil {
				t.Errorf("filepath.Abs() error = %v", err)
				return
			}

			options := &flags.GenerateBindingsOptions{
				TS:               tt.useTypescript,
				UseInterfaces:    tt.useInterface,
				ModelsFilename:   "models",
				OutputDirectory:  "frontend/bindings",
				ProjectDirectory: absDir,
				BasePath:         ".",
				UseBaseName:      true,
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

			for pkgDir, got := range allModels {
				want, ok := tt.want[pkgDir]
				if !ok {
					t.Errorf("GenerateModels() unexpected package = %v", pkgDir)
					continue
				}

				// convert all line endings to \n
				got = convertLineEndings(got)
				want = convertLineEndings(want)

				if diff := cmp.Diff(want, got); diff != "" {
					originalFilename := "models"
					if tt.useTypescript && tt.useInterface {
						originalFilename += ".interfaces"
					}
					outFileName := originalFilename + ".got"
					if tt.useTypescript {
						originalFilename += ".ts"
						outFileName += ".ts"
					} else {
						originalFilename += ".js"
						outFileName += ".js"
					}

					outDir := filepath.Join(tt.dir, options.OutputDirectory, pkgDir)
					originalFile := filepath.Join(outDir, originalFilename)
					// Check if file exists
					if _, err := os.Stat(originalFile); err != nil {
						outFileName = originalFilename
					}

					outFile := filepath.Join(outDir, outFileName)
					os.MkdirAll(outDir, 0755)
					if err != nil {
						t.Errorf("os.MkdirAll() error = %v", err)
						continue
					}
					err = os.WriteFile(outFile, []byte(got), 0644)
					if err != nil {
						t.Errorf("os.WriteFile() error = %v", err)
						continue
					}

					t.Errorf("GenerateModels() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func convertLineEndings(str string) string {
	// replace all \r\n with \n
	return strings.ReplaceAll(str, "\r\n", "\n")
}
