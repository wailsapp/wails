package webview2

// This test parses the package sources (so it runs on any GOOS, even though
// the COM wrappers themselves are windows-only) and verifies that every
// derived ICoreWebView2* vtable struct embeds its parent's vtable, so method
// slots line up with the real COM vtable layout. See wailsapp/wails#6129.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
	"testing"
)

type vtblInfo struct {
	embeds []string // embedded struct types, in order
	fields int      // directly declared (named) vtable slots
}

func parseVtbls(t *testing.T) map[string]*vtblInfo {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]*vtblInfo{}
	fset := token.NewFileSet()
	for _, e := range entries {
		n := e.Name()
		if !strings.HasSuffix(n, ".go") || strings.HasSuffix(n, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, n, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", n, err)
		}
		for _, d := range f.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, sp := range gd.Specs {
				ts, ok := sp.(*ast.TypeSpec)
				if !ok || !strings.HasSuffix(ts.Name.Name, "Vtbl") {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				info := &vtblInfo{}
				for _, fl := range st.Fields.List {
					if len(fl.Names) == 0 {
						if id, ok := fl.Type.(*ast.Ident); ok {
							info.embeds = append(info.embeds, id.Name)
						}
						continue
					}
					info.fields += len(fl.Names)
				}
				out[ts.Name.Name] = info
			}
		}
	}
	return out
}

// slots returns the total number of pointer-sized vtable slots of the named
// vtable struct, following embedded parents.
func slots(t *testing.T, all map[string]*vtblInfo, name string) int {
	t.Helper()
	if name == "IUnknownVtbl" {
		return 3
	}
	v, ok := all[name]
	if !ok {
		t.Fatalf("vtable %s not found", name)
	}
	n := v.fields
	for _, e := range v.embeds {
		n += slots(t, all, e)
	}
	return n
}

var numberedName = regexp.MustCompile(`^(.*?)(_?)(\d+)Vtbl$`)

// expectedParent returns the vtable the given struct must embed, per the
// WebView2 IDL inheritance (X -> X2 -> X3 ..., ICoreWebView2 -> _2 -> _3 ...).
func expectedParent(name string) string {
	m := numberedName.FindStringSubmatch(name)
	if m == nil || name == "ICoreWebView2Vtbl" {
		return "IUnknownVtbl"
	}
	base, us, k := m[1], m[2], m[3]
	// ICoreWebView2EnvironmentOptions2..8 derive directly from IUnknown.
	if strings.HasPrefix(base, "ICoreWebView2EnvironmentOptions") {
		return "IUnknownVtbl"
	}
	if k == "1" {
		return "IUnknownVtbl"
	}
	if k == "2" {
		return base + "Vtbl"
	}
	n := 0
	for _, c := range k {
		n = n*10 + int(c-'0')
	}
	return base + us + itoa(n-1) + "Vtbl"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for ; n > 0; n /= 10 {
		s = string(rune('0'+n%10)) + s
	}
	return s
}

func TestVtableInheritanceChain(t *testing.T) {
	all := parseVtbls(t)
	if len(all) < 200 {
		t.Fatalf("only found %d vtables; parsing is broken", len(all))
	}
	for name, v := range all {
		if name == "IUnknownVtbl" {
			continue
		}
		want := expectedParent(name)
		if len(v.embeds) != 1 || v.embeds[0] != want {
			t.Errorf("%s embeds %v, want exactly [%s]", name, v.embeds, want)
		}
	}
}

func TestICoreWebView2VtableSlots(t *testing.T) {
	all := parseVtbls(t)

	// Methods added by each ICoreWebView2_N, N = 2..27, per the WebView2 IDL.
	// ICoreWebView2 itself adds 58 on top of IUnknown's 3.
	added := map[int]int{
		2: 7, 3: 5, 4: 4, 5: 2, 6: 1, 7: 1, 8: 7, 9: 9, 10: 2, 11: 3, 12: 3,
		13: 1, 14: 3, 15: 4, 16: 3, 17: 1, 18: 2, 19: 2, 20: 1, 21: 1, 22: 2,
		23: 1, 24: 2, 25: 3, 26: 2, 27: 2,
	}
	total := 3 + 58
	if got := slots(t, all, "ICoreWebView2Vtbl"); got != total {
		t.Fatalf("ICoreWebView2Vtbl has %d slots, want %d", got, total)
	}
	for n := 2; n <= 27; n++ {
		total += added[n]
		name := "ICoreWebView2_" + itoa(n) + "Vtbl"
		if got := slots(t, all, name); got != total {
			t.Errorf("%s has %d slots, want %d", name, got, total)
		}
	}

	// Anchors from the issue: _2 ends at slot count 68, _3 at 73, and _22 at
	// 125 total with its first own method at index 123.
	for name, want := range map[string]int{
		"ICoreWebView2_2Vtbl":  68,
		"ICoreWebView2_3Vtbl":  73,
		"ICoreWebView2_22Vtbl": 125,
	} {
		if got := slots(t, all, name); got != want {
			t.Errorf("%s has %d slots, want %d", name, got, want)
		}
	}
	if got := slots(t, all, "ICoreWebView2_21Vtbl"); got != 123 {
		t.Errorf("AddWebResourceRequestedFilterWithRequestSourceKinds would be at slot %d, want 123", got)
	}
}
