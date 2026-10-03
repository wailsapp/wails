package capabilities

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"updater/internal/notes"
)

func TestEmitProducesValidGo(t *testing.T) {
	m := Mapping{
		"ICoreWebView2":    {SDKVersion: "", RuntimeVersion: ""},
		"ICoreWebView2_2":  {SDKVersion: "1.0.705.50", RuntimeVersion: "86.0.616.0"},
		"ICoreWebView2_3":  {SDKVersion: "1.0.721.34", RuntimeVersion: "87.0.658.0"},
		"ICoreWebView2_17": {SDKVersion: "1.0.1823.32", RuntimeVersion: "108.0.1462.37"},
		"ICoreWebView2_27": {SDKVersion: "1.0.2903.40", RuntimeVersion: "121.0.2277.83"},
	}
	out, err := Emit(m)
	if err != nil {
		t.Fatalf("Emit: %v", err)
	}

	// The build tag prefix is stripped before parsing — go/parser doesn't
	// process build constraints itself but the //go:build line is valid.
	_, err = parser.ParseFile(token.NewFileSet(), "capabilities.go", out, parser.AllErrors)
	if err != nil {
		t.Fatalf("emitted file is not valid Go:\n%s\nerror: %v", out, err)
	}

	src := string(out)
	for _, want := range []string{
		`"ICoreWebView2": {SDKVersion: "", MinRuntimeVersion: ""}`,
		`"ICoreWebView2_27": {SDKVersion: "1.0.2903.40", MinRuntimeVersion: "121.0.2277.83"}`,
		"func SupportsInterface",
		"var InterfaceSupportTable",
		"DO NOT EDIT",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("emitted file missing %q", want)
		}
	}
}

func TestSortedByVersionThenName(t *testing.T) {
	m := Mapping{
		"ICoreWebView2_27":  {SDKVersion: "1.0.2903.40"},
		"ICoreWebView2":     {SDKVersion: "1.0.1018.46"},
		"ICoreWebView2_2":   {SDKVersion: "1.0.705.50"},
		"ICoreWebView2Base": {}, // baseline sorts before any versioned entry
	}
	keys := m.Sorted()
	want := []string{"ICoreWebView2Base", "ICoreWebView2_2", "ICoreWebView2", "ICoreWebView2_27"}
	if len(keys) != len(want) {
		t.Fatalf("len = %d, want %d", len(keys), len(want))
	}
	for i := range keys {
		if keys[i] != want[i] {
			t.Errorf("Sorted()[%d] = %q, want %q", i, keys[i], want[i])
		}
	}
}

func TestEmitEmpty(t *testing.T) {
	if _, err := Emit(nil); err == nil {
		t.Error("Emit(empty) should error")
	}
}

func TestEmitJSON(t *testing.T) {
	m := Mapping{
		"ICoreWebView2_2": {SDKVersion: "1.0.705.50", RuntimeVersion: "86.0.616.0"},
		"ICoreWebView2":   {SDKVersion: "1.0.1018.46", RuntimeVersion: "92.0.902.49"},
	}
	js := string(EmitJSON(m))
	// Stable order — older version first.
	if i, j := strings.Index(js, "_2"), strings.Index(js, `"ICoreWebView2"`); i == -1 || j == -1 || i > j {
		t.Errorf("EmitJSON ordering wrong:\n%s", js)
	}
	if !strings.Contains(js, `"runtime": "86.0.616.0"`) {
		t.Errorf("EmitJSON missing runtime version:\n%s", js)
	}
}

func TestBuild(t *testing.T) {
	support := map[string]notes.Support{
		"ICoreWebView2_17":      {SDKVersion: "1.0.1518.46", RuntimeVersion: "110.0.1518.46"},
		"ICoreWebView2Interop9": {SDKVersion: "1.0.1700.0", RuntimeVersion: "111.0.0.0"},
	}
	inventory := []string{"ICoreWebView2", "ICoreWebView2_17"}
	oldestInventory := []string{"ICoreWebView2"}

	m, err := Build(support, inventory, oldestInventory)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	// Notes-dated interface.
	if got := m["ICoreWebView2_17"]; got.RuntimeVersion != "110.0.1518.46" {
		t.Errorf("ICoreWebView2_17 runtime = %q, want 110.0.1518.46", got.RuntimeVersion)
	}
	// Baseline interface: predates the archive, always supported.
	if got := m["ICoreWebView2"]; got.SDKVersion != "" || got.RuntimeVersion != "" {
		t.Errorf("ICoreWebView2 should be a baseline entry, got %+v", got)
	}
	// Notes-only interop interface is kept.
	if _, ok := m["ICoreWebView2Interop9"]; !ok {
		t.Error("notes-only interop interface dropped from mapping")
	}
}

func TestBuildFailsOnUndatedInterface(t *testing.T) {
	// In the current IDL, absent from notes AND absent from the oldest IDL:
	// must fail loudly, not silently omit.
	_, err := Build(nil, []string{"ICoreWebView2Mystery"}, nil)
	if err == nil {
		t.Fatal("Build should fail for an interface it cannot date")
	}
	if !strings.Contains(err.Error(), "ICoreWebView2Mystery") {
		t.Errorf("error should name the offending interface, got: %v", err)
	}
}
