package parser

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/wailsapp/wails/v3/internal/flags"
)

//go:embed testdata
var testdata embed.FS

func getFile(filename string) string {
	// get the file from the testdata FS
	file, err := fs.ReadFile(testdata, filename)
	if err != nil {
		return ""
	}
	return string(file)
}

func ParseProjectAndPkgs(options *flags.GenerateBindingsOptions) (*Project, error) {
	project, err := ParseProject(options)
	if err != nil {
		return project, err
	}

	project.pkgs, err = ParsePackages(project)
	return project, err
}

func TestGenerateBindings(t *testing.T) {

	testOptions := []struct {
		name          string
		useNames      bool
		useTypescript bool
		ext           string
	}{
		{
			name:          "Javascript - CallById",
			useNames:      false,
			useTypescript: false,
			ext:           "js",
		},
		{
			name:          "Javascript - CallByName",
			useNames:      true,
			useTypescript: false,
			ext:           "name.js",
		},
		{
			name:          "Typescript - CallById",
			useNames:      false,
			useTypescript: true,
			ext:           "ts",
		},
		{
			name:          "Typescript - CallByName",
			useNames:      true,
			useTypescript: true,
			ext:           "name.ts",
		},
	}

	tests := []struct {
		name              string
		dir               string
		want              map[string]map[string]bool
		useBundledRuntime bool
	}{
		{
			name: "complex_json",
			dir:  "testdata/complex_json",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
			},
		},
		{
			name: "complex_method",
			dir:  "testdata/complex_method",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
			},
		},
		{
			name: "enum",
			dir:  "testdata/enum",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
			},
		},
		{
			name: "enum_from_imported_package",
			dir:  "testdata/enum_from_imported_package",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
			},
		},
		{
			name: "function_single",
			dir:  "testdata/function_single",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
			},
		},
		{
			name: "function_multiple_files",
			dir:  "testdata/function_multiple_files",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
			},
		},
		{
			name: "function_from_imported_package",
			dir:  "testdata/function_from_imported_package",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
				"services": {
					"OtherService": true,
				},
			},
		},
		{
			name: "function_from_nested_imported_package",
			dir:  "testdata/function_from_nested_imported_package",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
				"services/other": {
					"OtherService": true,
				},
			},
		},
		{
			name: "struct_literal_multiple",
			dir:  "testdata/struct_literal_multiple",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
					"OtherService": true,
				},
			},
		},
		{
			name: "function_from_imported_package",
			dir:  "testdata/function_from_imported_package",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
				"services": {
					"OtherService": true,
				},
			},
		},
		{
			name: "variable_single",
			dir:  "testdata/variable_single",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
			},
		},
		{
			name: "variable_single_from_function",
			dir:  "testdata/variable_single_from_function",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
			},
		},
		{
			name: "variable_single_from_other_function",
			dir:  "testdata/variable_single_from_other_function",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
				"services": {
					"OtherService": true,
				},
			},
		},
		{
			name: "struct_literal_single",
			dir:  "testdata/struct_literal_single",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
			},
		},
		{
			name: "struct_literal_multiple_other",
			dir:  "testdata/struct_literal_multiple_other",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
				"services": {
					"OtherService": true,
				},
			},
		},
		{
			name: "struct_literal_non_pointer_single",
			dir:  "testdata/struct_literal_non_pointer_single",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
			},
		},
		{
			name: "struct_literal_multiple_files",
			dir:  "testdata/struct_literal_multiple_files",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
					"OtherService": true,
				},
			},
		},
		{
			name: "function_single_context",
			dir:  "testdata/function_single_context",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
			},
		},
		{
			name: "nested_types",
			dir:  "testdata/nested_types",
			want: map[string]map[string]bool{
				"main": {
					"OtherService": true,
				},
			},
		},
		{
			name: "multiple_packages",
			dir:  "testdata/multiple_packages",
			want: map[string]map[string]bool{
				"log": {
					"Logger": true,
				},
				"main": {
					"GreetService": true,
					"Greeter":      true,
				},
				"other": {
					"OtherService": true,
				},
				"other/other": {
					"OtherService": true,
				},
				// "github.com-samber-lo": {
				// 	"Tupel2": true,
				// },
			},
			useBundledRuntime: true,
		},
		{
			name: "renamed_import",
			dir:  "testdata/renamed_import",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
			},
		},
		{
			name: "type_alias",
			dir:  "testdata/type_alias",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
			},
		},
		{
			name: "interfaces",
			dir:  "testdata/interfaces",
			want: map[string]map[string]bool{
				"main": {
					"GreetService": true,
				},
			},
		},
		{
			name: "app_outside_main",
			dir:  "testdata/app_outside_main/app",
			want: map[string]map[string]bool{
				"app": {
					"GreetService": true,
				},
				"services": {
					"OtherService": true,
				},
				"github.com/wailsapp/wails/v3/internal/parser/testdata/app_outside_main/other": {
					"OtherService": true,
				},
			},
		},
	}

	type Test struct {
		name              string
		dir               string
		want              map[string]map[string]string
		useNames          bool
		useTypescript     bool
		useBundledRuntime bool
	}

	allTests := []Test{}
	for _, tt := range tests {
		for _, option := range testOptions {

			want := make(map[string]map[string]string)

			for pkgDir, services := range tt.want {
				files := make(map[string]string)

				for serviceName := range services {
					filePath := fmt.Sprintf("%s/frontend/bindings/%s/%s.%s", tt.dir, pkgDir, serviceName, option.ext)
					files[serviceName] = getFile(filePath)
				}
				want[pkgDir] = files
			}

			allTests = append(allTests, Test{
				name:              tt.name + " - " + option.name,
				dir:               tt.dir,
				want:              want,
				useNames:          option.useNames,
				useTypescript:     option.useTypescript,
				useBundledRuntime: tt.useBundledRuntime,
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
				ModelsFilename:    "models",
				UseNames:          tt.useNames,
				TS:                tt.useTypescript,
				OutputDirectory:   "frontend/bindings",
				ProjectDirectory:  absDir,
				UseBundledRuntime: tt.useBundledRuntime,
				BasePath:          ".",
				UseBaseName:       true,
			}

			project, err := ParseProjectAndPkgs(options)
			if err != nil {
				t.Errorf("ParseProjectAndPkgs() error = %v", err)
				return
			}

			// Generate Bindings
			got, err := project.GenerateBindings()
			if err != nil {
				t.Fatalf("GenerateBindings() error = %v", err)
			}

			// check if bindings are missing
			for dirName, structDetails := range tt.want {
				for name := range structDetails {
					_, ok := got[dirName][name]
					if !ok {
						t.Errorf("GenerateBindings() missing binding = %v/%v", dirName, name)
						continue
					}
				}
			}

			for dirName, structDetails := range got {
				// iterate the struct names in structDetails
				for name, binding := range structDetails {
					expected, ok := tt.want[dirName][name]
					if !ok {
						t.Errorf("GenerateBindings() unexpected binding = %v/%v", dirName, name)
						continue
					}

					// compare the binding

					// convert all line endings to \n
					binding = convertLineEndings(binding)
					expected = convertLineEndings(expected)

					if diff := cmp.Diff(expected, binding); diff != "" {
						originalFilename := name
						if tt.useNames {
							originalFilename += ".name"
						}
						outFileName := originalFilename + ".got"
						if tt.useTypescript {
							originalFilename += ".ts"
							outFileName += ".ts"
						} else {
							originalFilename += ".js"
							outFileName += ".js"
						}

						outDir := filepath.Join(tt.dir, options.OutputDirectory, dirName)
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
						err = os.WriteFile(outFile, []byte(binding), 0644)
						if err != nil {
							t.Errorf("os.WriteFile() error = %v", err)
							continue
						}

						t.Errorf("GenerateBindings() mismatch (-want +got):\n%s", diff)
					}
				}
			}
		})
	}
}
