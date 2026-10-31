// Command puregocheck proves that the CGO-free (purego) backends can reach
// every native function they call, and that none of them silently stub out
// behaviour the cgo backends implement.
//
// It works from the source of the build being checked: the files selected by
// GOOS and -tags with cgo disabled, exactly as `go build` would select them.
//
//	puregocheck [-tags T] [-goos OS] [-json FILE] [-allow FILE] static
//	    Inventory and stub report only. Runs anywhere.
//	puregocheck [-tags T] [-json FILE] [-allow FILE] resolve
//	    The static checks, then loads the platform libraries and resolves the
//	    whole inventory: every C symbol (dlsym), and on macOS every
//	    Objective-C class (objc_getClass) and selector (implemented by some
//	    class in the runtime). Must run on the target OS, ideally a clean VM.
//
// Exit status 1 means a required symbol, class or selector is missing, a
// purego function is a stub where its cgo twin is not (and it is not listed
// in the allow file), or an Objective-C message passes a narrow argument on
// the darwin/arm64 stack where purego misplaces it (stacknarrow, stackargs.go). Optional entries (registerOptional, a manifest's
// "optional", a direct Dlsym, respondsTo/classExists guards) are reported but
// never fail.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// The packages that hold purego code, relative to the v3 module root.
var packages = []string{
	"pkg/application",
	"pkg/events",
	"pkg/mac",
	"pkg/services/dock",
	"pkg/services/notifications",
	"internal/assetserver/webview",
	"internal/operatingsystem",
}

type entry struct {
	Name     string   `json:"name"`
	Optional bool     `json:"optional,omitempty"`
	Ours     bool     `json:"ours,omitempty"` // defined by the backend itself (registered class / method)
	Where    []string `json:"where"`

	// Linux, from GIR (-gir): the library and version the symbol appeared in,
	// and whether that is after the supported floor (so it must be optional).
	Namespace string `json:"namespace,omitempty"`
	Since     string `json:"since,omitempty"`
	Gated     bool   `json:"gated,omitempty"`
}

type stub struct {
	Name   string `json:"name"`
	Purego string `json:"purego"`
	Cgo    string `json:"cgo"`
}

type report struct {
	GOOS      string   `json:"goos"`
	Tags      string   `json:"tags"`
	Files     int      `json:"files"`
	Symbols   []*entry `json:"symbols"`
	Classes   []*entry `json:"classes,omitempty"`
	Selectors []*entry `json:"selectors,omitempty"`
	Dynamic   []string `json:"dynamic,omitempty"` // registrations with a non-literal name: not verifiable
	Stubs     []stub   `json:"stubs"`
	Allowed   []stub   `json:"allowed_stubs,omitempty"`
	Missing   []string `json:"missing,omitempty"`   // required, not resolved at runtime
	Absent    []string `json:"absent,omitempty"`    // optional, not resolved at runtime
	Versions  []string `json:"versions,omitempty"`  // rule 1 breaches (see versions.go)
	Unguarded []string `json:"unguarded,omitempty"` // rule 2 breaches
	// poolthread (darwin, see poolthread.go): autorelease pools pushed without
	// pinning the goroutine to its thread.
	PoolThread []string `json:"pool_thread,omitempty"`

	// stacknarrow (darwin only, see stackargs.go): narrow arguments that
	// purego places on the stack at the wrong offset.
	StackNarrow    []stackFinding `json:"stack_narrow,omitempty"`
	StackChecked   int            `json:"stack_checked,omitempty"`   // calls classified
	StackBundled   int            `json:"stack_bundled,omitempty"`   // fixed-signature narrow spills purego packs
	StackUnknown   int            `json:"stack_unknown,omitempty"`   // calls with an unmapped type, not classified
	StackForwarded int            `json:"stack_forwarded,omitempty"` // send helpers passing args... on
	StackSkipped   []string       `json:"stack_skipped,omitempty"`   // where those calls are
	AllowedStack   []stackFinding `json:"allowed_stack,omitempty"`
	Resolved       bool           `json:"resolved"`
	Loaded         []string       `json:"loaded_versions,omitempty"` // library versions found at runtime
	Libraries      []string       `json:"libraries,omitempty"`       // loaded for the resolve step

	libCandidates []string // every library name the selected sources mention
}

func main() {
	os.Exit(run(os.Args[1:]))
}

// run is main with its arguments and exit status passed in and returned, so
// the tests can drive the whole command.
func run(args []string) int {
	fs := flag.NewFlagSet("puregocheck", flag.ContinueOnError)
	tags := fs.String("tags", "", "build tags of the purego build (comma or space separated)")
	goos := fs.String("goos", runtime.GOOS, "target OS (static mode only; resolve uses the running OS)")
	root := fs.String("root", ".", "v3 module root")
	jsonOut := fs.String("json", "", "write the full report as JSON to this file")
	girDir := fs.String("gir", "", "GIR directory (e.g. /usr/share/gir-1.0): check symbols against library version floors")
	floorsFile := fs.String("floors", "", "minimum supported library versions (default internal/puregocheck/floors.json)")
	allowFile := fs.String("allow", "", "file of intentional stubs, one name per line (default internal/puregocheck/allow-<goos>.txt)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	mode := fs.Arg(0)
	if mode != "static" && mode != "resolve" {
		fmt.Fprintln(os.Stderr, "usage: puregocheck [flags] static|resolve")
		return 2
	}
	if mode == "resolve" && *goos != runtime.GOOS {
		fmt.Fprintln(os.Stderr, "resolve runs on the target OS; drop -goos")
		return 2
	}
	tagList := strings.FieldsFunc(*tags, func(r rune) bool { return r == ',' || r == ' ' })
	if *allowFile == "" {
		*allowFile = filepath.Join(*root, "internal/puregocheck", "allow-"+*goos+".txt")
	}

	r, err := scan(*root, *goos, tagList)
	if err != nil {
		fmt.Fprintln(os.Stderr, "puregocheck:", err)
		return 2
	}
	r.Tags = strings.Join(tagList, ",")
	if *girDir != "" {
		if *floorsFile == "" {
			*floorsFile = filepath.Join(*root, "internal/puregocheck", "floors.json")
		}
		fl, err := loadFloors(*floorsFile)
		if err == nil {
			var gir map[string]since
			if gir, err = readGIR(*girDir); err == nil {
				checkVersions(r, gir, fl)
			}
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "puregocheck: versions:", err)
			return 2
		}
	}
	applyAllow(r, *allowFile)

	if mode == "resolve" {
		if err := resolve(r); err != nil {
			fmt.Fprintln(os.Stderr, "puregocheck: resolve:", err)
			return 2
		}
		r.Resolved = true
	}

	if *jsonOut != "" {
		b, _ := json.MarshalIndent(r, "", "  ")
		if err := os.WriteFile(*jsonOut, b, 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "puregocheck:", err)
			return 2
		}
	}
	return summary(r)
}

func applyAllow(r *report, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	allow := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			allow[strings.Fields(line)[0]] = true
		}
	}
	var keep []stub
	var keepStack []stackFinding
	for _, f := range r.StackNarrow {
		if allow[f.Name] {
			r.AllowedStack = append(r.AllowedStack, f)
		} else {
			keepStack = append(keepStack, f)
		}
	}
	r.StackNarrow = keepStack
	for _, s := range r.Stubs {
		if allow[s.Name] {
			r.Allowed = append(r.Allowed, s)
		} else {
			keep = append(keep, s)
		}
	}
	r.Stubs = keep
}

// summary prints the human report (the last line is machine-readable) and
// returns the exit status.
func summary(r *report) int {
	count := func(es []*entry) (req, opt, ours int) {
		for _, e := range es {
			switch {
			case e.Ours:
				ours++
			case e.Optional:
				opt++
			default:
				req++
			}
		}
		return
	}
	sr, so, _ := count(r.Symbols)
	fmt.Printf("purego %s [%s]: %d files\n", r.GOOS, r.Tags, r.Files)
	fmt.Printf("  C symbols: %d required, %d optional\n", sr, so)
	if r.GOOS == "darwin" {
		cr, co, cu := count(r.Classes)
		xr, xo, xu := count(r.Selectors)
		fmt.Printf("  ObjC classes: %d required, %d optional, %d registered by the backend\n", cr, co, cu)
		fmt.Printf("  ObjC selectors: %d required, %d optional, %d implemented by the backend\n", xr, xo, xu)
	}
	if len(r.Dynamic) > 0 {
		fmt.Printf("  not verifiable (name is not a literal): %d\n", len(r.Dynamic))
		for _, d := range r.Dynamic {
			fmt.Println("    " + d)
		}
	}
	if gated := countGated(r.Symbols); gated > 0 {
		fmt.Printf("  version-gated (newer than the floor, optional): %d\n", gated)
	}
	for _, v := range r.Versions {
		fmt.Println("  VERSION " + v)
	}
	for _, u := range r.Unguarded {
		fmt.Println("  UNGUARDED " + u)
	}
	if r.GOOS == "darwin" {
		fmt.Print(stackNarrowSummary(r))
	}
	for _, f := range r.StackNarrow {
		fmt.Println("  " + f.Text + " at " + f.Where)
	}
	if len(r.AllowedStack) > 0 {
		fmt.Printf("  stack findings allowed by the allow file: %d\n", len(r.AllowedStack))
	}
	for _, s := range r.Stubs {
		fmt.Printf("  STUB %s (purego %s is empty; cgo %s is not)\n", s.Name, s.Purego, s.Cgo)
	}
	if len(r.Allowed) > 0 {
		fmt.Printf("  stubs allowed by the allow file: %d\n", len(r.Allowed))
	}
	if r.Resolved {
		fmt.Printf("  libraries loaded: %s\n", strings.Join(r.Libraries, ", "))
		if len(r.Loaded) > 0 {
			fmt.Printf("  library versions: %s\n", strings.Join(r.Loaded, ", "))
		}
		for _, m := range r.Missing {
			fmt.Println("  MISSING " + m)
		}
		for _, a := range r.Absent {
			fmt.Println("  absent (optional) " + a)
		}
	}
	for _, p := range r.PoolThread {
		fmt.Println("  POOLTHREAD " + p + ": NSAutoreleasePool without runtime.LockOSThread")
	}
	status := 0
	if len(r.Stubs) > 0 || len(r.Missing) > 0 || len(r.Versions) > 0 || len(r.Unguarded) > 0 || len(r.StackNarrow) > 0 || len(r.PoolThread) > 0 {
		status = 1
	}
	fmt.Printf("puregocheck goos=%s tags=%s resolved=%t missing=%d stubs=%d versions=%d unguarded=%d stacknarrow=%d poolthread=%d absent=%d dynamic=%d\n",
		r.GOOS, r.Tags, r.Resolved, len(r.Missing), len(r.Stubs), len(r.Versions), len(r.Unguarded), len(r.StackNarrow), len(r.PoolThread), len(r.Absent), len(r.Dynamic))
	return status
}

func countGated(es []*entry) int {
	n := 0
	for _, e := range es {
		if e.Gated {
			n++
		}
	}
	return n
}

func sortedEntries(m map[string]*entry) []*entry {
	out := make([]*entry, 0, len(m))
	for _, e := range m {
		sort.Strings(e.Where)
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
