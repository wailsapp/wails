package parser

import (
	"go/types"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/packages"
)

func TestFindServices(t *testing.T) {
	tests := []struct {
		pkg  string
		dir  string
		want []string
	}{
		{
			pkg: "complex_json",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_json.GreetService",
			},
		},
		{
			pkg: "complex_method",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/complex_method.GreetService",
			},
		},
		{
			pkg: "enum",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/enum.GreetService",
			},
		},
		{
			pkg: "enum_from_imported_package",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/enum_from_imported_package.GreetService",
			},
		},
		{
			pkg: "function_from_imported_package",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_imported_package.GreetService",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_imported_package/services.OtherService",
			},
		},
		{
			pkg: "function_from_nested_imported_package",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_nested_imported_package.GreetService",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_from_nested_imported_package/services/other.OtherService",
			},
		},
		{
			pkg: "function_multiple_files",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_multiple_files.GreetService",
			},
		},
		{
			pkg: "function_single",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_single.GreetService",
			},
		},
		{
			pkg: "function_single_context",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/function_single_context.GreetService",
			},
		},
		{
			pkg: "slice_expressions",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/slice_expressions.Service1",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/slice_expressions.Service2",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/slice_expressions.Service3",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/slice_expressions.Service4",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/slice_expressions.Service5",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/slice_expressions.Service6",
			},
		},
		{
			pkg: "struct_literal_multiple",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple.GreetService",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple.OtherService",
			},
		},
		{
			pkg: "struct_literal_multiple_files",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_files.GreetService",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_files.OtherService",
			},
		},
		{
			pkg: "struct_literal_multiple_other",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_other.GreetService",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_multiple_other/services.OtherService",
			},
		},
		{
			pkg: "struct_literal_non_pointer_single",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_non_pointer_single.GreetService",
			},
		},
		{
			pkg: "struct_literal_single",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/struct_literal_single.GreetService",
			},
		},
		{
			pkg: "variable_single",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single.GreetService",
			},
		},
		{
			pkg: "variable_single_from_function",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single_from_function.GreetService",
			},
		},
		{
			pkg: "variable_single_from_other_function",
			want: []string{
				"github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single_from_other_function.GreetService",
				"github.com/wailsapp/wails/v3/internal/parser/testdata/variable_single_from_other_function/services.OtherService",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.pkg, func(t *testing.T) {
			pkgs, err := LoadPackages(nil, true,
				"github.com/wailsapp/wails/v3/internal/parser/testdata/"+tt.pkg,
			)
			if err != nil {
				t.Fatal(err)
			}

			packages.PrintErrors(pkgs)

			services, err := FindServices(pkgs)
			if err != nil {
				t.Fatal(err)
			}

			got := make([]string, len(services))
			for i, srv := range services {
				got[i] = types.TypeString(srv.Type(), nil)
			}

			slices.Sort(tt.want)
			slices.Sort(got)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FindServices() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
