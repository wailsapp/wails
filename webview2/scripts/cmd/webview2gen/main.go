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

	// Incremental write: only touch files whose content actually changed, so
	// repeated runs are cheap, file mtimes stay stable for build caching, and
	// the diff after an SDK bump contains exactly the real changes.
	written, unchanged := 0, 0
	expected := map[string]bool{}
	for _, f := range files {
		expected[f.FileName] = true
		path := filepath.Join(*out, f.FileName)
		content := normalizeNewlines(f.Content.Bytes())
		if old, err := os.ReadFile(path); err == nil && bytes.Equal(normalizeNewlines(old), content) {
			unchanged++
			continue
		}
		if err := os.WriteFile(path, content, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
		written++
	}

	// Remove generated files the current IDL no longer produces, applying the
	// same exclusions as verify (separately-emitted and hand-written files).
	stale, err := staleGeneratedFiles(*out, expected)
	if err != nil {
		return err
	}
	for _, name := range stale {
		path := filepath.Join(*out, name)
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("remove stale %s: %w", path, err)
		}
		fmt.Fprintf(os.Stderr, "removed stale file: %s\n", path)
	}

	fmt.Fprintf(os.Stderr, "generated %d files in %s from SDK %s (%d written, %d unchanged, %d removed)\n",
		len(files), *out, v, written, unchanged, len(stale))
	return nil
}

// generatedFileExclusions are .go files in the output directory that the
// `generate` subcommand does not produce and must never delete or flag:
// capabilities.go comes from the `capabilities` subcommand, doc.go is the
// only hand-written file in the package, and *_test.go files are managed by
// the `tests` emission. Shared by generate (stale removal) and verify.
func isExcludedGeneratedFile(name string) bool {
	return name == "capabilities.go" || name == "doc.go" || strings.HasSuffix(name, "_test.go")
}

func staleGeneratedFiles(dir string, expected map[string]bool) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read output dir: %w", err)
	}
	var stale []string
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || isExcludedGeneratedFile(name) {
			continue
		}
		if !expected[name] {
			stale = append(stale, name)
		}
	}
	sort.Strings(stale)
	return stale, nil
}

// normalizeNewlines collapses CRLF to LF. text/template preserves whatever
// line endings the template source carried, so on Windows checkouts with
// core.autocrlf=true the templates become CRLF and the emitted Go files
// diverge from the LF golden files committed to the repo. Stripping CR
// in the generator keeps the output deterministic across platforms.
func normalizeNewlines(b []byte) []byte {
	return bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
}

// -----------------------------------------------------------------------
// capabilities
// -----------------------------------------------------------------------

func runCapabilities(args []string) error {
	fs := flag.NewFlagSet("capabilities", flag.ContinueOnError)
	source := fs.String("source", "", "release-notes markdown file (default: fetch from MicrosoftDocs)")
	version := fs.String("version", "", "SDK version of the cached IDL providing the interface inventory (default: latest cached)")
	dir := fs.String("dir", IDLDir, "IDL cache directory")
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
	support := notes.InterfaceSupport(releases)
	if len(support) == 0 {
		return errors.New("no interfaces extracted from release notes — check parser")
	}

	// The capability table must cover every interface the generator emits.
	// The current IDL provides the inventory; the oldest cached IDL vouches
	// for interfaces that predate the release-notes archive.
	store := idl.NewStore(*dir)
	v, err := resolveVersion(store, *version)
	if err != nil {
		return err
	}
	idlBytes, err := store.Read(v)
	if err != nil {
		return fmt.Errorf("read IDL %s: %w", v, err)
	}
	inventory, err := generator.InterfaceNames(idlBytes)
	if err != nil {
		return fmt.Errorf("parse IDL %s: %w", v, err)
	}
	oldestVersion, err := oldestCachedVersion(store)
	if err != nil {
		return err
	}
	oldestBytes, err := store.Read(oldestVersion)
	if err != nil {
		return fmt.Errorf("read IDL %s: %w", oldestVersion, err)
	}
	oldestInventory, err := generator.InterfaceNames(oldestBytes)
	if err != nil {
		return fmt.Errorf("parse IDL %s: %w", oldestVersion, err)
	}

	mapping, err := capabilities.Build(support, inventory, oldestInventory)
	if err != nil {
		return err
	}

	emitted, err := capabilities.Emit(mapping)
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
	testArgs := []string{"test", "./..."}
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
		if !bytes.Equal(normalizeNewlines(got), normalizeNewlines(f.Content.Bytes())) {
			diffs = append(diffs, fmt.Sprintf("changed file: %s", path))
		}
	}

	// Look for committed files that the generator no longer produces.
	stale, err := staleGeneratedFiles(*out, expected)
	if err != nil {
		return err
	}
	for _, name := range stale {
		diffs = append(diffs, fmt.Sprintf("unexpected committed file: %s", filepath.Join(*out, name)))
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

// runFull chains download → generate → capabilities → verify. Each step has
// its own flag set, so full parses a shared superset once and forwards only
// the flags each step understands — passing raw args through verbatim breaks
// any step that doesn't define one of them.
func runFull(args []string) error {
	fs := flag.NewFlagSet("full", flag.ContinueOnError)
	version := fs.String("version", "", "SDK version to download/generate (default: latest cached)")
	dir := fs.String("dir", IDLDir, "IDL cache directory")
	out := fs.String("out", OutputDir, "output directory for generated bindings")
	source := fs.String("source", "", "release-notes markdown file for capabilities (default: fetch from MicrosoftDocs)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	common := []string{"-version", *version, "-dir", *dir}
	capArgs := append([]string{"-out", *out}, common...)
	if *source != "" {
		capArgs = append(capArgs, "-source", *source)
	}
	for _, step := range []struct {
		name string
		fn   func([]string) error
		args []string
	}{
		{"download", runDownload, common},
		{"generate", runGenerate, append([]string{"-out", *out}, common...)},
		{"capabilities", runCapabilities, capArgs},
		{"verify", runVerify, append([]string{"-out", *out}, common...)},
	} {
		fmt.Fprintf(os.Stderr, "==> %s\n", step.name)
		if err := step.fn(step.args); err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}
	return nil
}

// -----------------------------------------------------------------------
// helpers
// -----------------------------------------------------------------------

func oldestCachedVersion(store *idl.Store) (string, error) {
	cached, err := store.List()
	if err != nil {
		return "", fmt.Errorf("list cache: %w", err)
	}
	if len(cached) == 0 {
		return "", errors.New("no IDL cached")
	}
	sort.Slice(cached, func(i, j int) bool {
		c, _ := idlversion.Compare(cached[i], cached[j])
		return c < 0
	})
	return cached[0], nil
}

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
