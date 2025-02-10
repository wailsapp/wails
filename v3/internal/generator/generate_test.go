package generator

import (
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v3/internal/generator/render"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/generator/config"
)

const testcases = "github.com/wailsapp/wails/v3/internal/generator/testcases/..."

type testParams struct {
	name      string
	options   *flags.GenerateBindingsOptions
	outputDir string
	want      map[string]bool
}

func TestGenerator(t *testing.T) {
	const (
		useNamesBit = 1 << iota
		useInterfacesBit
		tsBit
	)

	// Generate configuration matrix.
	tests := make([]*testParams, 1<<3)
	for i := range tests {
		options := &flags.GenerateBindingsOptions{
			ModelsFilename: "models",
			IndexFilename:  "index",

			UseBundledRuntime: true,

			TS:            i&tsBit != 0,
			UseInterfaces: i&useInterfacesBit != 0,
			UseNames:      i&useNamesBit != 0,
		}

		name := configString(options)

		tests[i] = &testParams{
			name:      name,
			options:   options,
			outputDir: filepath.Join("testdata/output", name),
			want:      make(map[string]bool),
		}
	}

	for _, test := range tests {
		// Create output dir.
		if err := os.MkdirAll(test.outputDir, 0777); err != nil {
			t.Fatal(err)
		}

		// Walk output dir.
		err := filepath.WalkDir(test.outputDir, func(path string, d fs.DirEntry, err error) error {
			// Skip directories.
			if d.IsDir() {
				return nil
			}

			// Skip got files.
			if strings.HasSuffix(d.Name(), ".got.js") || strings.HasSuffix(d.Name(), ".got.ts") {
				return nil
			}

			// Record file.
			test.want[filepath.Clean("."+path[len(test.outputDir):])] = false
			return nil
		})

		if err != nil {
			t.Fatal(err)
		}
	}

	// Run tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			creator := outputCreator(t, test)

			generator := NewGenerator(
				test.options,
				creator,
				config.DefaultPtermLogger(nil),
			)

			_, err := generator.Generate(testcases)
			if report := (*ErrorReport)(nil); errors.As(err, &report) {
				if report.HasErrors() {
					t.Error(report)
				} else if report.HasWarnings() {
					pterm.Warning.Println(report)
				}

				// Log warnings and compare with reference output.
				if log, err := creator.Create("warnings.log"); err != nil {
					t.Error(err)
				} else {
					func() {
						defer log.Close()

						warnings := report.Warnings()
						slices.Sort(warnings)

						for _, msg := range warnings {
							fmt.Fprint(log, msg, render.Newline)
						}
					}()
				}
			} else if err != nil {
				t.Error(err)
			}

			for path, present := range test.want {
				if !present {
					t.Errorf("Missing output file '%s'", path)
				}
			}
		})
	}
}

// configString computes a subtest name from the given configuration.
func configString(options *flags.GenerateBindingsOptions) string {
	lang := "JS"
	if options.TS {
		lang = "TS"
	}
	return fmt.Sprintf("lang=%s/UseInterfaces=%v/UseNames=%v", lang, options.UseInterfaces, options.UseNames)
}

// outputCreator returns a FileCreator that detects want/got pairs
// and schedules them for comparison.
//
// If no corresponding want file exists, it is created and reported.
func outputCreator(t *testing.T, params *testParams) config.FileCreator {
	var mu sync.Mutex
	return config.FileCreatorFunc(func(path string) (io.WriteCloser, error) {
		path = filepath.Clean(path)
		prefixedPath := filepath.Join(params.outputDir, path)

		// Protect want map accesses.
		mu.Lock()
		defer mu.Unlock()

		if seen, ok := params.want[path]; ok {
			// File exists: mark as seen and compare.
			if seen {
				t.Errorf("Duplicate output file '%s'", path)
			}
			params.want[path] = true

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
			params.want[path] = true

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
