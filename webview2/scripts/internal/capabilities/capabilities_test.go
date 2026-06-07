package capabilities

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestEmitProducesValidGo(t *testing.T) {
	m := Mapping{
		"ICoreWebView2":    "1.0.1018.46",
		"ICoreWebView2_2":  "1.0.705.50",
		"ICoreWebView2_3":  "1.0.721.34",
		"ICoreWebView2_17": "1.0.1823.32",
		"ICoreWebView2_27": "1.0.2903.40",
	}
	out, err := Emit(m, nil)
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
		`"ICoreWebView2": "1.0.1018.46"`,
		`"ICoreWebView2_27": "1.0.2903.40"`,
		"func SupportsInterface",
		"func HasCapability",
		"var AllCapabilities",
		"DO NOT EDIT",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("emitted file missing %q", want)
		}
	}
}

func TestSortedByVersionThenName(t *testing.T) {
	m := Mapping{
		"ICoreWebView2_27": "1.0.2903.40",
		"ICoreWebView2":    "1.0.1018.46",
		"ICoreWebView2_2":  "1.0.705.50",
	}
	keys := m.Sorted()
	// 705.50 is oldest, then 1018.46, then 2903.40.
	want := []string{"ICoreWebView2_2", "ICoreWebView2", "ICoreWebView2_27"}
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
	if _, err := Emit(nil, nil); err == nil {
		t.Error("Emit(empty) should error")
	}
}

func TestEmitJSON(t *testing.T) {
	m := Mapping{
		"ICoreWebView2_2": "1.0.705.50",
		"ICoreWebView2":   "1.0.1018.46",
	}
	js := string(EmitJSON(m))
	// Stable order — older version first.
	if i, j := strings.Index(js, "_2"), strings.Index(js, `"ICoreWebView2"`); i == -1 || j == -1 || i > j {
		t.Errorf("EmitJSON ordering wrong:\n%s", js)
	}
}
