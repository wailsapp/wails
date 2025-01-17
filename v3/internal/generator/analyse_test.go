package generator

import (
	"encoding/json"
	"errors"
	"go/types"
	"os"
	"path"
	"path/filepath"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/generator/config"
)

func TestAnalyser(t *testing.T) {
	type testParams struct {
		name string
		want []string
	}

	// Gather tests from cases directory.
	entries, err := os.ReadDir("testcases")
	if err != nil {
		t.Fatal(err)
	}

	tests := make([]testParams, 0, len(entries)+1)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		test := testParams{
			name: entry.Name(),
			want: make([]string, 0),
		}

		want, err := os.Open(filepath.Join("testcases", entry.Name(), "bound_types.json"))
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				t.Fatal(err)
			}
		} else {
			err = json.NewDecoder(want).Decode(&test.want)
			want.Close()
			if err != nil {
				t.Fatal(err)
			}
		}

		for i := range test.want {
			test.want[i] = path.Clean("github.com/wailsapp/wails/v3/internal/generator/testcases/" + test.name + test.want[i])
		}
		slices.Sort(test.want)

		tests = append(tests, test)
	}

	// Add global test.
	{
		all := testParams{
			name: "...",
		}

		for _, test := range tests {
			all.want = append(all.want, test.want...)
		}
		slices.Sort(all.want)

		tests = append(tests, all)
	}

	// Resolve system package paths.
	systemPaths, err := ResolveSystemPaths(nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		pkgPattern := "github.com/wailsapp/wails/v3/internal/generator/testcases/" + test.name

		t.Run("pkg="+test.name, func(t *testing.T) {
			pkgs, err := LoadPackages(nil, pkgPattern)
			if err != nil {
				t.Fatal(err)
			}

			for _, pkg := range pkgs {
				for _, err := range pkg.Errors {
					pterm.Warning.Println(err)
				}
			}

			got := make([]string, 0)

			services, err := FindServices(pkgs, systemPaths, config.DefaultPtermLogger(nil))
			if err != nil {
				t.Error(err)
			}

			for obj := range services {
				got = append(got, types.TypeString(obj.Type(), nil))
			}

			slices.Sort(got)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Found services mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
