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

	"github.com/stretchr/testify/assert"
)

func TestAnalyser(t *testing.T) {
	type testParams struct {
		name   string
		want   []string
		events bool
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

		events, err := os.Open(filepath.Join("testcases", entry.Name(), "events.json"))
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				t.Fatal(err)
			}
		} else {
			err = json.NewDecoder(events).Decode(&test.events)
			events.Close()
			if err != nil {
				t.Fatal(err)
			}
		}

		tests = append(tests, test)
	}

	// Add global test.
	{
		all := testParams{
			name: "...",
		}

		for _, test := range tests {
			all.want = append(all.want, test.want...)
			all.events = all.events || test.events
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
					t.Logf("package warning: %v", err)
				}
			}

			got := make([]string, 0)

			services, registerService, err := FindServices(pkgs, systemPaths, analyserTestLogger{t})
			if err != nil {
				t.Error(err)
			}

			for obj := range services {
				got = append(got, types.TypeString(obj.Type(), nil))
			}

			slices.Sort(got)

			assert.Equal(t, test.want, got, "Found services mismatch")

			if test.events != (registerService != nil) {
				t.Errorf("Found events mismatch: wanted=%t, got=%t", test.events, registerService != nil)
			}
		})
	}
}

type analyserTestLogger struct{ t *testing.T }

func (logger analyserTestLogger) Errorf(format string, args ...any) {
	logger.t.Helper()
	logger.t.Logf("error: "+format, args...)
}

func (logger analyserTestLogger) Warningf(format string, args ...any) {
	logger.t.Helper()
	logger.t.Logf("warning: "+format, args...)
}

func (logger analyserTestLogger) Infof(format string, args ...any) {
	logger.t.Helper()
	logger.t.Logf("info: "+format, args...)
}

func (logger analyserTestLogger) Debugf(format string, args ...any) {
	logger.t.Helper()
	logger.t.Logf("debug: "+format, args...)
}

func (logger analyserTestLogger) Statusf(format string, args ...any) {
	logger.t.Helper()
	logger.t.Logf("status: "+format, args...)
}
