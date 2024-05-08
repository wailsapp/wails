package parser

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser/config"
)

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
			ModelsFilename:   "models",
			InternalFilename: "internal",
			IndexFilename:    "index",

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
		want map[string]map[string]bool
	}

	// Gather tests from cases directory.
	entries, err := os.ReadDir("testcases")
	if err != nil {
		t.Fatal(err)
	}

	// Add global test.
	entries = append(entries, nil)

	tests := make([]testParams, 0, len(entries))

	for _, entry := range entries {
		name := "all"
		if entry != nil {
			name = entry.Name()

			if !entry.IsDir() {
				continue
			}
		}

		test := testParams{
			name: name,
			want: make(map[string]map[string]bool),
		}

		// Fill wanted file maps.
		for _, config := range configs {
			want := make(map[string]bool)
			test.want[config.name] = want

			// Compute output dir and create it.
			outputDir := filepath.Join("testdata", name, config.name)
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

	// Run tests.
	for _, test := range tests {
		pkgPattern := "github.com/wailsapp/wails/v3/internal/parser/testcases/"
		if test.name != "all" {
			pkgPattern += test.name + "/"
		}
		pkgPattern += "..."

		t.Run("pkg="+test.name, func(t *testing.T) {
			for _, conf := range configs {
				t.Run(conf.name, func(t *testing.T) {
					want := test.want[conf.name]

					generator := NewGenerator(
						conf.GenerateBindingsOptions,
						outputCreator(t, test.name, conf.name, want),
						config.DefaultPtermLogger,
					)

					_, err := generator.Generate(pkgPattern)
					if report := (*ErrorReport)(nil); errors.As(err, &report) {
						if report.HasErrors() {
							t.Error(report)
						} else if report.HasWarnings() {
							pterm.Warning.Println(report)
						}
					} else if err != nil {
						t.Error(err)
					}

					for path, present := range want {
						if !present {
							t.Errorf("Missing output file '%s'", path)
						}
					}
				})
			}
		})
	}
}

// configString computes a subtest name from the given configuration.
func configString(options *flags.GenerateBindingsOptions) string {
	if options.TS {
		return fmt.Sprintf("lang=TS/UseInterfaces=%v/UseNames=%v", options.UseInterfaces, options.UseNames)
	} else {
		return fmt.Sprintf("lang=JS/UseNames=%v", options.UseNames)
	}
}

// outputCreator returns a FileCreator that detects want/got pairs
// and schedules them for comparison.
//
// If no corresponding want file exists, it is created and reported.
func outputCreator(t *testing.T, testName, configName string, want map[string]bool) config.FileCreator {
	var mu sync.Mutex
	outputDir := filepath.Join("testdata", testName, configName)
	return config.FileCreatorFunc(func(path string) (io.WriteCloser, error) {
		path = filepath.Clean(path)
		prefixedPath := filepath.Join(outputDir, path)

		// Protect want map accesses.
		mu.Lock()
		defer mu.Unlock()

		if seen, ok := want[path]; ok {
			// File exists: mark as seen and compare.
			if seen {
				t.Errorf("Duplicate output file '%s'", path)
			}
			want[path] = true

			// Open want file.
			wf, err := os.Open(prefixedPath)
			if err != nil {
				return nil, err
			}

			// Create or truncate got file.
			ext := filepath.Ext(prefixedPath)
			gf, err := os.Create(fmt.Sprintf("%s.got%s", prefixedPath[:len(prefixedPath)-len(ext)], ext))
			if err != nil {
				return nil, err
			}

			// Initialise comparer.
			return &outputComparer{t, path, wf, gf}, nil
		} else {
			// File does not exist: create it.
			t.Errorf("Unexpected output file '%s'", path)
			want[path] = true

			if err := os.MkdirAll(filepath.Dir(prefixedPath), 0777); err != nil {
				return nil, err
			}

			return os.Create(prefixedPath)
		}
	})
}

// outputComparer is a io.WriteCloser that writes to got.
//
// When Close is called, it compares want to got; if they are identical,
// it deletes got; otherwise it reports a testing error.
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
		return nil
	}

	got, err := io.ReadAll(comparer.got)
	if err != nil {
		comparer.t.Error(err)
		return nil
	}

	if diff := cmp.Diff(want, got); diff != "" {
		comparer.t.Errorf("Output file '%s' mismatch (-want +got):\n%s", comparer.path, diff)
	} else {
		// On success, delete got file.
		comparer.got.Close()
		if err := os.Remove(comparer.got.Name()); err != nil {
			comparer.t.Error(err)
		}
	}

	return nil
}
