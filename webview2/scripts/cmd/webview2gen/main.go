// webview2gen is the WebView2 IDL → Go binding generator and capability
// table builder. It is the single entry point for refreshing
// pkg/webview2 against a new SDK release.
//
// Subcommands
//
//	download      Fetch a WebView2 SDK IDL (latest by default) into the
//	              local IDL cache. Use --version to pin.
//	generate      Parse a cached IDL and emit pkg/webview2/*.go.
//	capabilities  Fetch the SDK release notes, derive the
//	              interface→minimum-version mapping, and emit
//	              pkg/webview2/capabilities.go.
//	test          Run `go test ./...` against the generator + internal pkgs.
//	verify        Regenerate everything against the on-disk IDL and fail if
//	              the working tree differs — guards against hand-edits.
//	full          download → generate → capabilities → verify, in that order.
//
// Run `webview2gen <command> --help` for per-command flags.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"updater/generator"
	"updater/internal/capabilities"
	"updater/internal/idl"
	"updater/internal/idlversion"
	"updater/internal/notes"
)

const (
	// IDLDir is the on-disk cache for downloaded IDL files. The default
	// matches the natural invocation `go run ./cmd/webview2gen ...` from
	// inside scripts/ where the cached `WebView2.<version>.idl` files live.
	IDLDir = "."

	// OutputDir is where generated bindings live, relative to scripts/.
	OutputDir = "../pkg/webview2"
)

func main() {
	if len(os.Args) < 2 {
		usage(os.Stderr)
		os.Exit(2)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	cmds := map[string]func([]string) error{
		"download":     runDownload,
		"generate":     runGenerate,
		"capabilities": runCapabilities,
		"test":         runTest,
		"verify":       runVerify,
		"full":         runFull,
		"help":         func(_ []string) error { usage(os.Stdout); return nil },
		"-h":           func(_ []string) error { usage(os.Stdout); return nil },
		"--help":       func(_ []string) error { usage(os.Stdout); return nil },
	}

	fn, ok := cmds[cmd]
	if !ok {
		fmt.Fprintf(os.Stderr, "webview2gen: unknown command %q\n\n", cmd)
		usage(os.Stderr)
		os.Exit(2)
	}
	if err := fn(args); err != nil {
		fmt.Fprintf(os.Stderr, "webview2gen %s: %v\n", cmd, err)
		os.Exit(1)
	}
}

func usage(w io.Writer) {
	fmt.Fprint(w, `webview2gen — WebView2 IDL → Go binding generator

USAGE
  webview2gen <command> [flags]

COMMANDS
  download      Fetch an SDK IDL into the local cache.
  generate      Generate pkg/webview2 from a cached IDL.
  capabilities  Emit pkg/webview2/capabilities.go from SDK release notes.
  test          Run `+"`"+`go test ./...`+"`"+` for the generator + internal pkgs.
  verify        Regenerate and fail if the working tree differs.
  full          download → generate → capabilities → verify.

Use 'webview2gen <command> --help' for command flags.
`)
}

// -----------------------------------------------------------------------
// download
// -----------------------------------------------------------------------

func runDownload(args []string) error {
	fs := flag.NewFlagSet("download", flag.ContinueOnError)
	version := fs.String("version", "", "SDK version to download (e.g. 1.0.2903.40). If empty, the latest known cached version is used.")
	dir := fs.String("dir", IDLDir, "directory to cache IDL files in")
	if err := fs.Parse(args); err != nil {
		return err
	}

	store := idl.NewStore(*dir)
	fetcher := idl.NewFetcher(store)

	v := *version
	if v == "" {
		// No version given — fall back to the latest cached IDL so offline
		// runs work. Use the release-notes scrape only if explicitly asked.
		cached, err := store.List()
		if err != nil {
			return fmt.Errorf("list cache: %w", err)
		}
		if len(cached) == 0 {
			return errors.New("no cached versions and --version not specified")
		}
		sort.Slice(cached, func(i, j int) bool {
			c, _ := idlversion.Compare(cached[i], cached[j])
			return c < 0
		})
		v = cached[len(cached)-1]
		fmt.Fprintf(os.Stderr, "using latest cached version: %s\n", v)
	}

	if store.Has(v) {
		fmt.Fprintf(os.Stderr, "%s already cached at %s\n", v, store.CachePath(v))
		return nil
	}
	data, err := fetcher.Download(v)
	if err != nil {
		return fmt.Errorf("download %s: %w", v, err)
	}
	fmt.Fprintf(os.Stderr, "downloaded %s (%d bytes) → %s\n", v, len(data), store.CachePath(v))
	return nil
}

// -----------------------------------------------------------------------
// generate
// -----------------------------------------------------------------------

func runGenerate(args []string) error {
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	version := fs.String("version", "", "SDK version of the cached IDL to parse (default: latest cached)")
	dir := fs.String("dir", IDLDir, "IDL cache directory")
	out := fs.String("out", OutputDir, "output directory for generated bindings")
	if err := fs.Parse(args); err != nil {
		return err
	}

	store := idl.NewStore(*dir)
	v, err := resolveVersion(store, *version)
	if err != nil {
		return err
	}
	idlBytes, err := store.Read(v)
	if err != nil {
		return fmt.Errorf("read IDL %s: %w", v, err)
	}

	files, err := generator.ParseIDL(idlBytes)
	if err != nil {
		return fmt.Errorf("parse IDL: %w", err)
	}

	if err := os.MkdirAll(*out, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	for _, f := range files {
		path := filepath.Join(*out, f.FileName)
		if err := os.WriteFile(path, f.Content.Bytes(), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
	}
	fmt.Fprintf(os.Stderr, "generated %d files in %s from SDK %s\n", len(files), *out, v)
	return nil
}

// -----------------------------------------------------------------------
// capabilities
// -----------------------------------------------------------------------

func runCapabilities(args []string) error {
	fs := flag.NewFlagSet("capabilities", flag.ContinueOnError)
	source := fs.String("source", "", "release-notes markdown file (default: fetch from MicrosoftDocs)")
	out := fs.String("out", OutputDir, "output directory for capabilities.go")
	jsonOut := fs.String("json", "", "also write the interface→version map as JSON at this path (empty = skip)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	var md []byte
	var err error
	if *source != "" {
		md, err = os.ReadFile(*source)
		if err != nil {
			return fmt.Errorf("read source: %w", err)
		}
	} else {
		md, err = notes.Fetch()
		if err != nil {
			return fmt.Errorf("fetch release notes: %w", err)
		}
	}

	releases, err := notes.Parse(md)
	if err != nil {
		return fmt.Errorf("parse release notes: %w", err)
	}
	mapping := capabilities.Mapping(notes.InterfaceMinimumVersions(releases))
	if len(mapping) == 0 {
		return errors.New("no interfaces extracted from release notes — check parser")
	}

	emitted, err := capabilities.Emit(mapping, nil)
	if err != nil {
		return fmt.Errorf("emit: %w", err)
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	path := filepath.Join(*out, "capabilities.go")
	if err := os.WriteFile(path, emitted, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	fmt.Fprintf(os.Stderr, "wrote %s (%d interfaces)\n", path, len(mapping))

	if *jsonOut != "" {
		if err := os.WriteFile(*jsonOut, capabilities.EmitJSON(mapping), 0o644); err != nil {
			return fmt.Errorf("write JSON: %w", err)
		}
		fmt.Fprintf(os.Stderr, "wrote %s\n", *jsonOut)
	}
	return nil
}

// -----------------------------------------------------------------------
// test / verify / full
// -----------------------------------------------------------------------

func runTest(args []string) error {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	verbose := fs.Bool("v", false, "verbose test output")
	if err := fs.Parse(args); err != nil {
		return err
	}
	testArgs := []string{"test", "./generator/...", "./internal/..."}
	if *verbose {
		testArgs = append(testArgs, "-v")
	}
	cmd := exec.Command("go", testArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runVerify(args []string) error {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	version := fs.String("version", "", "SDK version to verify against (default: latest cached)")
	dir := fs.String("dir", IDLDir, "IDL cache directory")
	out := fs.String("out", OutputDir, "directory whose contents are compared against fresh generation")
	if err := fs.Parse(args); err != nil {
		return err
	}

	store := idl.NewStore(*dir)
	v, err := resolveVersion(store, *version)
	if err != nil {
		return err
	}
	idlBytes, err := store.Read(v)
	if err != nil {
		return fmt.Errorf("read IDL %s: %w", v, err)
	}

	files, err := generator.ParseIDL(idlBytes)
	if err != nil {
		return fmt.Errorf("parse IDL: %w", err)
	}

	// Build expected set; compare byte-for-byte against committed files.
	var diffs []string
	expected := map[string]bool{}
	for _, f := range files {
		expected[f.FileName] = true
		path := filepath.Join(*out, f.FileName)
		got, err := os.ReadFile(path)
		if err != nil {
			diffs = append(diffs, fmt.Sprintf("missing committed file: %s (%v)", path, err))
			continue
		}
		if !bytes.Equal(got, f.Content.Bytes()) {
			diffs = append(diffs, fmt.Sprintf("changed file: %s", path))
		}
	}

	// Look for committed files that the generator no longer produces.
	// capabilities.go is emitted separately and shouldn't be flagged.
	entries, err := os.ReadDir(*out)
	if err != nil {
		return fmt.Errorf("read output dir: %w", err)
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		if name == "capabilities.go" || strings.HasSuffix(name, "_test.go") {
			continue
		}
		if !expected[name] {
			diffs = append(diffs, fmt.Sprintf("unexpected committed file: %s", filepath.Join(*out, name)))
		}
	}

	if len(diffs) > 0 {
		sort.Strings(diffs)
		for _, d := range diffs {
			fmt.Fprintln(os.Stderr, d)
		}
		return fmt.Errorf("%d differences between regenerated and committed output", len(diffs))
	}
	fmt.Fprintf(os.Stderr, "verify ok: %d files match committed output for SDK %s\n", len(files), v)
	return nil
}

func runFull(args []string) error {
	for _, step := range []struct {
		name string
		fn   func([]string) error
	}{
		{"download", runDownload},
		{"generate", runGenerate},
		{"capabilities", runCapabilities},
		{"verify", runVerify},
	} {
		fmt.Fprintf(os.Stderr, "==> %s\n", step.name)
		if err := step.fn(args); err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}
	return nil
}

// -----------------------------------------------------------------------
// helpers
// -----------------------------------------------------------------------

func resolveVersion(store *idl.Store, want string) (string, error) {
	if want != "" {
		if !store.Has(want) {
			return "", fmt.Errorf("version %s not in cache (%s) — run 'webview2gen download --version %s' first",
				want, store.Dir, want)
		}
		return want, nil
	}
	cached, err := store.List()
	if err != nil {
		return "", fmt.Errorf("list cache: %w", err)
	}
	if len(cached) == 0 {
		return "", errors.New("no IDL cached — run 'webview2gen download --version <v>' first")
	}
	sort.Slice(cached, func(i, j int) bool {
		c, _ := idlversion.Compare(cached[i], cached[j])
		return c < 0
	})
	return cached[len(cached)-1], nil
}
