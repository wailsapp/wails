package parser

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
	"github.com/wailsapp/wails/v3/internal/parser/config"
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
			test.want[i] = path.Clean("github.com/wailsapp/wails/v3/internal/parser/testcases/" + test.name + test.want[i])
		}
		slices.Sort(test.want)

		tests = append(tests, test)
	}

	// Add global test.
	{
		all := testParams{
			name: "all",
		}

		for _, test := range tests {
			all.want = append(all.want, test.want...)
		}
		slices.Sort(all.want)

		tests = append(tests, all)
	}

	// Resolve wails app pkg path.
	wailsAppPkgPaths, err := ResolvePatterns(nil, WailsAppPkgPath)
	if err != nil {
		return
	}

	if len(wailsAppPkgPaths) < 1 {
		t.Fatal(ErrNoApplicationPackage)
	} else if len(wailsAppPkgPaths) > 1 {
		// This should never happen...
		t.Fatal("wails application package path matched multiple packages")
	}

	for _, test := range tests {
		pkgPattern := "github.com/wailsapp/wails/v3/internal/parser/testcases/"
		if test.name == "all" {
			pkgPattern += "..."
		} else {
			pkgPattern += test.name
		}

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

			err = FindServices(pkgs, wailsAppPkgPaths[0], config.DefaultPtermLogger(nil), func(tn *types.TypeName) bool {
				got = append(got, types.TypeString(tn.Type(), nil))
				return true
			})
			if err != nil && !errors.Is(err, ErrNoServices) {
				t.Error(err)
			}

			slices.Sort(got)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Found services mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
