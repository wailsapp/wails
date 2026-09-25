package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// allIDLFiles returns every bundled WebView2.*.idl filename in scripts/.
func allIDLFiles(t *testing.T) []string {
	t.Helper()
	matches, err := filepath.Glob("../WebView2.*.idl")
	require.NoError(t, err)
	require.NotEmpty(t, matches, "expected bundled IDL files")
	var names []string
	for _, m := range matches {
		names = append(names, filepath.Base(m))
	}
	return names
}

// TestParseIDL_AllFiles drives the full bindings pipeline (Process + Generate)
// over every bundled IDL and asserts a com.go plus one file per declaration is
// emitted with non-empty content.
func TestParseIDL_AllFiles(t *testing.T) {
	for _, name := range allIDLFiles(t) {
		name := name
		t.Run(name, func(t *testing.T) {
			files, err := ParseIDL(idlFile(t, name))
			require.NoError(t, err)
			require.NotEmpty(t, files)

			var sawCom bool
			for _, f := range files {
				assert.NotEmpty(t, f.FileName, "every file needs a name")
				assert.NotNil(t, f.Content)
				assert.Positive(t, f.Content.Len(), "%s should not be empty", f.FileName)
				if f.FileName == "com.go" {
					sawCom = true
				}
			}
			assert.True(t, sawCom, "com.go must be generated")
		})
	}
}

// TestParseIDLWithTests_AllFiles drives the full pipeline including the test
// emitter (per-interface *_gen_test.go plus the package-wide abi_gen_test.go)
// over every bundled IDL.
func TestParseIDLWithTests_AllFiles(t *testing.T) {
	for _, name := range allIDLFiles(t) {
		name := name
		t.Run(name, func(t *testing.T) {
			withTests, err := ParseIDLWithTests(idlFile(t, name))
			require.NoError(t, err)

			plain, err := ParseIDL(idlFile(t, name))
			require.NoError(t, err)

			// The test pipeline must emit strictly more files than the plain one.
			assert.Greater(t, len(withTests), len(plain),
				"ParseIDLWithTests should add generated test files")

			var sawABI bool
			for _, f := range withTests {
				if f.FileName == "abi_gen_test.go" {
					sawABI = true
					assert.Contains(t, f.Content.String(), "TestVtblSlotCounts")
					assert.Contains(t, f.Content.String(), "TestStructSizes")
					assert.Contains(t, f.Content.String(), "TestCapabilityTableCoversAllInterfaces")
				}
				if strings.HasSuffix(f.FileName, "_gen_test.go") {
					assert.Contains(t, f.Content.String(), "//go:build windows")
				}
			}
			assert.True(t, sawABI, "abi_gen_test.go must be generated")
		})
	}
}

// TestInterfaceNames_AllFiles checks every declared interface is reported in
// declaration order and matches the bindings emitted by Generate.
func TestInterfaceNames_AllFiles(t *testing.T) {
	for _, name := range allIDLFiles(t) {
		name := name
		t.Run(name, func(t *testing.T) {
			data := idlFile(t, name)
			names, err := InterfaceNames(data)
			require.NoError(t, err)
			require.NotEmpty(t, names)

			// ICoreWebView2 is present in every shipped SDK.
			assert.Contains(t, names, "ICoreWebView2")

			// Names are unique and non-empty.
			seen := map[string]bool{}
			for _, n := range names {
				assert.NotEmpty(t, n)
				assert.False(t, seen[n], "duplicate interface name %q", n)
				seen[n] = true
			}
		})
	}
}

// TestInterfaceMethods_AllFiles checks the per-interface method inventory.
func TestInterfaceMethods_AllFiles(t *testing.T) {
	for _, name := range allIDLFiles(t) {
		name := name
		t.Run(name, func(t *testing.T) {
			data := idlFile(t, name)
			methods, err := InterfaceMethods(data)
			require.NoError(t, err)
			require.NotEmpty(t, methods)

			// The set of keys must equal the InterfaceNames inventory.
			names, err := InterfaceNames(data)
			require.NoError(t, err)
			assert.Len(t, methods, len(names))
			for _, n := range names {
				_, ok := methods[n]
				assert.True(t, ok, "method map missing interface %q", n)
			}

			// ICoreWebView2 owns exactly 58 methods in every SDK.
			assert.Len(t, methods["ICoreWebView2"], 58)
		})
	}
}

// TestParseIDL_Error feeds malformed IDL and asserts each entry point surfaces
// the parse error instead of panicking.
func TestParseIDL_Error(t *testing.T) {
	bad := []byte("this is not valid idl {{{")
	if _, err := ParseIDL(bad); err == nil {
		t.Error("ParseIDL should reject malformed input")
	}
	if _, err := ParseIDLWithTests(bad); err == nil {
		t.Error("ParseIDLWithTests should reject malformed input")
	}
	if _, err := InterfaceNames(bad); err == nil {
		t.Error("InterfaceNames should reject malformed input")
	}
	if _, err := InterfaceMethods(bad); err == nil {
		t.Error("InterfaceMethods should reject malformed input")
	}
}

// readIDL is a convenience for sub-tests that need raw bytes by path.
func readIDL(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../" + name)
	require.NoError(t, err)
	return data
}

const craftLibHeader = "[uuid(11111111-1111-1111-1111-111111111111), version(1.0)]\nlibrary WebView2 {\n"
const craftIfaceUUID = "[uuid(22222222-2222-2222-2222-222222222222), object, pointer_default(unique)]\n"

// TestParseIDL_ProcessError: valid syntax that fails during Process (an [in]
// parameter of an unknown by-value type). Exercises the Process error arm of
// both ParseIDL and ParseIDLWithTests.
func TestParseIDL_ProcessError(t *testing.T) {
	idl := []byte(craftLibHeader + craftIfaceUUID +
		"interface IFoo : IUnknown {\n  HRESULT Take([in] Mystery m);\n}\n}\n")
	if _, err := ParseIDL(idl); err == nil {
		t.Error("ParseIDL should surface the Process error")
	}
	if _, err := ParseIDLWithTests(idl); err == nil {
		t.Error("ParseIDLWithTests should surface the Process error")
	}
}

// TestParseIDL_GenerateError: valid syntax that fails during Generate (an
// interface derived from an undeclared base).
func TestParseIDL_GenerateError(t *testing.T) {
	idl := []byte(craftLibHeader + craftIfaceUUID +
		"interface IFoo : INotDeclared {\n  HRESULT Ping();\n}\n}\n")
	if _, err := ParseIDL(idl); err == nil {
		t.Error("ParseIDL should surface the Generate error")
	}
	if _, err := ParseIDLWithTests(idl); err == nil {
		t.Error("ParseIDLWithTests should surface the Generate error")
	}
}

// TestParseIDLWithTests_GenerateTestsError: generates bindings fine, but the
// test emitter has no strategy for a by-value struct parameter.
func TestParseIDLWithTests_GenerateTestsError(t *testing.T) {
	idl := []byte(craftLibHeader +
		"typedef struct S { UINT32 x; } S;\n" + craftIfaceUUID +
		"interface IFoo : IUnknown {\n  HRESULT Take([in] S s);\n}\n}\n")
	if _, err := ParseIDL(idl); err != nil {
		t.Fatalf("bindings should generate fine: %v", err)
	}
	if _, err := ParseIDLWithTests(idl); err == nil {
		t.Error("ParseIDLWithTests should surface the GenerateTests error")
	}
}
