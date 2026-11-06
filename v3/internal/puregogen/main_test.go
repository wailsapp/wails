package main

import (
	"strings"
	"testing"
)

func TestHeaderTypes(t *testing.T) {
	d := declarations{aliases: map[string]string{"gsize": "unsigned long", "gboolean": "int", "gchar": "char", "GCallback": "void (*)(void)"}}
	for _, tc := range []struct {
		c, goType string
		str       bool
	}{
		{"gboolean", "int32", false}, {"gsize", "uint64", false},
		{"GCallback", "uintptr", false}, {"const gchar *", "string", true},
		{"GError **", "uintptr", false}, {"double", "float64", false},
		{"float", "float32", false}, {"unsigned char", "uint8", false},
	} {
		got, err := d.goType(tc.c, tc.str)
		if err != nil || got != tc.goType {
			t.Errorf("%s: got %q, %v; want %s", tc.c, got, err, tc.goType)
		}
	}
	for _, tc := range []struct {
		c   string
		str bool
	}{
		{"struct GValue", false}, {"union data", false}, {"__int128", false},
		{"int *", true}, {"char **", true}, {"gboolean", true},
	} {
		if _, err := d.goType(tc.c, tc.str); err == nil {
			t.Errorf("accepted unsupported type %s", tc.c)
		}
	}
}

func TestGenerateFromClangAST(t *testing.T) {
	input := `{"kind":"TranslationUnitDecl","inner":[
		{"kind":"TypedefDecl","name":"gboolean","type":{"qualType":"int"}},
		{"kind":"FunctionDecl","name":"native_open","type":{"qualType":"void *(const char *, gboolean)"},"inner":[
			{"kind":"ParmVarDecl","type":{"qualType":"const char *"}},
			{"kind":"ParmVarDecl","type":{"qualType":"gboolean"}}
		]}
	]}`
	d, err := readDeclarations(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	config := manifest{Package: "native", Build: "linux", Functions: []function{{Name: "native_open", Library: "lib", Strings: []int{0}}}}
	output, err := generate(config, d)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(output), "native_open func(string, int32) uintptr") {
		t.Fatalf("unexpected signature:\n%s", output)
	}
	if !strings.Contains(string(output), `register(&_native_native_open, lib, "native_open")`) {
		t.Fatalf("missing registration:\n%s", output)
	}
	config.Functions[0].Strings = []int{2}
	if _, err := generate(config, d); err == nil {
		t.Fatal("accepted out-of-range string parameter")
	}
	config.Functions[0].Strings = nil
	n := d.functions["native_open"]
	n.Variadic = true
	d.functions["native_open"] = n
	if _, err := generate(config, d); err == nil {
		t.Fatal("accepted variadic function")
	}
	delete(d.functions, "native_open")
	if _, err := generate(config, d); err == nil {
		t.Fatal("silently omitted missing symbol")
	}
}
