package parser

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateBindings(t *testing.T) {

	tests := []string{
		"struct_literal_single",
		"struct_literal_multiple",
		"struct_literal_multiple_other",
		"struct_literal_multiple_files",
		"function_single",
		"function_from_imported_package",
		"variable_single",
		"variable_single_from_function",
		"variable_single_from_other_function",
	}
	for _, projectDir := range tests {
		t.Run(projectDir, func(t *testing.T) {
			projectDir = "testdata/" + projectDir
			// Run parser on directory
			project, err := ParseProject(projectDir)
			if err != nil {
				t.Errorf("ParseProject() error = %v", err)
				return
			}

			// Generate Bindings
			got := GenerateBindings(project.BoundMethods)

			// Load bindings.js from project directory
			expected, err := os.ReadFile(projectDir + "/bindings.js")
			if err != nil {
				// Write file to project directory
				err = os.WriteFile(projectDir+"/bindings.got.js", []byte(got), 0644)
				if err != nil {
					t.Errorf("os.WriteFile() error = %v", err)
					return
				}
				t.Errorf("os.ReadFile() error = %v", err)
				return
			}

			// Compare
			if diff := cmp.Diff(string(expected), got); diff != "" {
				// Write file to project directory
				err = os.WriteFile(projectDir+"/bindings.got.js", []byte(got), 0644)
				if err != nil {
					t.Errorf("os.WriteFile() error = %v", err)
					return
				}
				t.Fatalf("GenerateService() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
