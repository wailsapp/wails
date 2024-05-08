package parser

import (
	"errors"
	"go/types"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pterm/pterm"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/parser/analyse"
)

func TestAnalyser(t *testing.T) {
	tests := []struct {
		name string
		pkgs []string
		dir  string
		want []string
	}{
		{
			name: "complex_expressions",
			pkgs: []string{
				"complex_expressions",
				"complex_expressions/config",
			},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_expressions.Service1",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_expressions.Service2",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_expressions.Service3",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_expressions.Service4",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_expressions.Service5",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_expressions.Service6",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_expressions/config.Service7",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_expressions/config.Service8",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_expressions/config.Service9",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_expressions/config.Service10",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_expressions/config.Service11",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_expressions/config.Service12",
			},
		},
		{
			name: "complex_json",
			pkgs: []string{"complex_json"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_json.GreetService",
			},
		},
		{
			name: "complex_method",
			pkgs: []string{"complex_method"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_method.GreetService",
			},
		},
		{
			name: "enum",
			pkgs: []string{"enum"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/enum.GreetService",
			},
		},
		{
			name: "enum_from_imported_package",
			pkgs: []string{"enum_from_imported_package"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/enum_from_imported_package.GreetService",
			},
		},
		{
			name: "function_from_imported_package",
			pkgs: []string{"function_from_imported_package"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_imported_package.GreetService",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_imported_package/services.OtherService",
			},
		},
		{
			name: "function_from_nested_imported_package",
			pkgs: []string{"function_from_nested_imported_package"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_nested_imported_package.GreetService",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_nested_imported_package/services/other.OtherService",
			},
		},
		{
			name: "function_multiple_files",
			pkgs: []string{"function_multiple_files"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_multiple_files.GreetService",
			},
		},
		{
			name: "function_single",
			pkgs: []string{"function_single"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_single.GreetService",
			},
		},
		{
			name: "function_single_context",
			pkgs: []string{"function_single_context"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_single_context.GreetService",
			},
		},
		{
			name: "struct_literal_multiple",
			pkgs: []string{"struct_literal_multiple"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple.GreetService",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple.OtherService",
			},
		},
		{
			name: "struct_literal_multiple_files",
			pkgs: []string{"struct_literal_multiple_files"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_files.GreetService",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_files.OtherService",
			},
		},
		{
			name: "struct_literal_multiple_other",
			pkgs: []string{"struct_literal_multiple_other"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_other.GreetService",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_other/services.OtherService",
			},
		},
		{
			name: "struct_literal_non_pointer_single",
			pkgs: []string{"struct_literal_non_pointer_single"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_non_pointer_single.GreetService",
			},
		},
		{
			name: "struct_literal_single",
			pkgs: []string{"struct_literal_single"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_single.GreetService",
			},
		},
		{
			name: "variable_single",
			pkgs: []string{"variable_single"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single.GreetService",
			},
		},
		{
			name: "variable_single_from_function",
			pkgs: []string{"variable_single_from_function"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single_from_function.GreetService",
			},
		},
		{
			name: "variable_single_from_other_function",
			pkgs: []string{"variable_single_from_other_function"},
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single_from_other_function.GreetService",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single_from_other_function/services.OtherService",
			},
		},
		{
			name: "all",
		},
	}

	// Test loading and analysing all packages together.
	allTest := &tests[len(tests)-1]
	for _, tt := range tests {
		if tt.name == "all" {
			continue
		}

		allTest.pkgs = append(allTest.pkgs, tt.pkgs...)
		allTest.want = append(allTest.want, tt.want...)
	}

	for _, tt := range tests {
		t.Run("pkg="+tt.name, func(t *testing.T) {
			pkgs, err := LoadPackages(nil, true,
				lo.Map(tt.pkgs, func(p string, _ int) string {
					return "github.com/wailsapp/wails/v3/internal/parser/testdata/" + p
				})...,
			)
			if err != nil {
				t.Fatal(err)
			}

			for _, pkg := range pkgs {
				for _, err := range pkg.Errors {
					pterm.Error.Println(err)
				}
			}

			analyser := analyse.NewAnalyser(pkgs)
			if err := analyser.Run(nil); err != nil && !errors.Is(err, analyse.ErrNoApplicationPackage) {
				t.Fatal(err)
			}

			services := analyser.Results()

			got := make([]string, len(services))
			for i, srv := range services {
				got[i] = types.TypeString(srv.Type(), nil)
			}

			slices.Sort(tt.want)
			slices.Sort(got)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Found services mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
