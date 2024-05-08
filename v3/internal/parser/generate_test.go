package parser

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/wailsapp/wails/v3/internal/flags"
)

// configString computes a subtest name from the given configuration.
func configString(options *flags.GenerateBindingsOptions) string {
	if options.TS {
		return fmt.Sprintf("lang=TS/UseInterfaces=%v/UseNames=%v", options.UseInterfaces, options.UseNames)
	} else {
		return fmt.Sprintf("lang=JS/UseNames=%v", options.UseNames)
	}
}

func TestGenerator(t *testing.T) {
	const (
		useNamesBit = 1 << iota
		tsBit
		useInterfacesBit
	)

	type configParams struct {
		name string
		*flags.GenerateBindingsOptions
	}

	// Generate configuration matrix.
	configs := make([]configParams, (1<<1)+(1<<2))
	for i := range configs {
		options := &flags.GenerateBindingsOptions{
			TS:            i&(tsBit|useInterfacesBit) != 0,
			UseInterfaces: i&useInterfacesBit != 0,
			UseNames:      i&useNamesBit != 0,
		}

		configs[i] = configParams{
			name:                    configString(options),
			GenerateBindingsOptions: options,
		}
	}

	type testParams struct {
		name string
		pkgs []string
		want map[string]map[string]bool
	}

	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	// Gather tests from data directory.
	tests := make([]testParams, 0, len(entries))
	all := -1

	for _, entry := range entries {
		name := entry.Name()

		if !entry.IsDir() {
			continue
		}

		// Remember position of 'all' test.
		if name == "all" {
			all = len(tests)
		}

		test := testParams{
			name: name,
			pkgs: []string{"github.com/wailsapp/wails/v3/internal/parser/testdata/" + name},
			want: make(map[string]map[string]bool),
		}

		// Fix complex_expressions test by hand.
		if name == "complex_expressions" {
			test.pkgs = append(test.pkgs, "github.com/wailsapp/wails/v3/internal/parser/testdata/complex_expressions/config")
		}

		// Fill wanted file maps.
		for _, config := range configs {
			want := make(map[string]bool)
			test.want[config.name] = want

			// Compute output dir and create it.
			outputDir := filepath.Join("testdata", name, "assets/bindings", config.name)
			if err := os.MkdirAll(outputDir, 0777); err != nil {
				t.Fatal(err)
			}

			// Walk output dir.
			err := filepath.WalkDir(outputDir, func(path string, d fs.DirEntry, err error) error {
				// Skip directories.
				if d.IsDir() {
					return nil
				}

				// Skip got files.
				if strings.HasSuffix(d.Name(), ".got.js") || strings.HasSuffix(d.Name(), ".got.ts") {
					return nil
				}

				// Record file.
				want[filepath.Clean(path[len(outputDir)+1:])] = false
				return nil
			})

			if err != nil {
				t.Fatal(err)
			}
		}

		tests = append(tests, test)
	}

	// Setup 'all' test.
	if all >= 0 {
		tests[all].pkgs = nil
		for _, test := range tests {
			tests[all].pkgs = append(tests[all].pkgs, test.pkgs...)
		}
	}

	// Run tests.
	for _, test := range tests {
		t.Run("pkg="+test.name, func(t *testing.T) {
			for _, config := range configs {
				t.Run(config.name, func(t *testing.T) {
					generator := NewGenerator(
						config.GenerateBindingsOptions,
						outputCreator(t, test.name, config.name, test.want[config.name]),
					)

					_, err := generator.Generate(test.pkgs...)
					if err != nil {
						t.Error(err)
					}
				})
			}
		})
	}
}

func outputCreator(t *testing.T, testName, configName string, want map[string]bool) FileCreator {
	var mu sync.Mutex
	outputDir := filepath.Join("testdata", testName, "assets/bindings", configName)
	return FileCreatorFunc(func(path string) (io.WriteCloser, error) {
		path = filepath.Clean(path)
		prefixedPath := filepath.Join(outputDir, path)

		// Protect want map accesses.
		mu.Lock()
		defer mu.Unlock()

		if seen, ok := want[path]; ok {
			// File exists: compare and mark as seen.
			if seen {
				err := fmt.Errorf("Duplicate output file '%s'", path)
				t.Error(err)
				return nil, err
			} else {
				want[path] = true

				// Open want file.
				wf, err := os.Open(prefixedPath)
				if err != nil {
					t.Error(err)
					return nil, err
				}

				// Create or truncate got file.
				ext := filepath.Ext(prefixedPath)
				gf, err := os.Create(fmt.Sprintf("%s.got%s", prefixedPath[:len(prefixedPath)-len(ext)], ext))
				if err != nil {
					t.Error(err)
					return nil, err
				}

				// Initialise comparer.
				return &outputComparer{t, path, wf, gf}, nil
			}
		} else {
			// File does not exist: create it.
			t.Errorf("Unexpected output file '%s'", path)
			want[path] = true

			if err := os.MkdirAll(filepath.Dir(prefixedPath), 0777); err != nil {
				t.Error(err)
				return nil, err
			}

			return os.Create(prefixedPath)
		}
	})
}

type outputComparer struct {
	t    *testing.T
	path string
	want *os.File
	got  *os.File
}

func (comparer *outputComparer) Write(data []byte) (int, error) {
	return comparer.got.Write(data)
}

func (comparer *outputComparer) Close() error {
	defer comparer.want.Close()
	defer comparer.got.Close()

	comparer.got.Seek(0, io.SeekStart)

	// Read want data.
	want, err := io.ReadAll(comparer.want)
	if err != nil {
		comparer.t.Error(err)
		return err
	}

	got, err := io.ReadAll(comparer.got)
	if err != nil {
		comparer.t.Error(err)
		return err
	}

	if diff := cmp.Diff(want, got); diff != "" {
		comparer.t.Errorf("Output file '%s' mismatch (-want +got):\n%s", comparer.path, diff)
	} else {
		// On success, delete got file.
		comparer.got.Close()
		if err := os.Remove(comparer.got.Name()); err != nil {
			comparer.t.Error(err)
			return err
		}
	}

	return nil
}
