package generator

import (
	"os"
	"testing"
	"updater/generator/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// idlFile is a helper that reads an IDL file from the scripts/ directory
// (one level above the generator package).
func idlFile(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../" + name)
	require.NoError(t, err, "reading %s", name)
	return data
}

// idlCounts holds the expected declaration counts for one IDL file.
type idlCounts struct {
	interfaces int
	enums      int
	structs    int
	forwards   int
}

var idlCountTable = map[string]idlCounts{
	"WebView2.1.0.1823.32.idl": {interfaces: 197, enums: 39, structs: 2, forwards: 196},
	"WebView2.1.0.1901.177.idl": {interfaces: 198, enums: 40, structs: 2, forwards: 197},
	"WebView2.1.0.2045.28.idl": {interfaces: 198, enums: 40, structs: 2, forwards: 197},
	"WebView2.1.0.2592.51.idl": {interfaces: 231, enums: 46, structs: 2, forwards: 229},
	"WebView2.1.0.2739.15.idl": {interfaces: 244, enums: 51, structs: 2, forwards: 243},
	"WebView2.1.0.2903.40.idl": {interfaces: 252, enums: 51, structs: 2, forwards: 251},
}

func TestParseAllIDLFiles(t *testing.T) {
	for filename, want := range idlCountTable {
		filename, want := filename, want
		t.Run(filename, func(t *testing.T) {
			data := idlFile(t, filename)

			idl, err := Parser.ParseBytes("", data)
			require.NoError(t, err, "parsing %s", filename)
			require.Len(t, idl.Libraries, 1, "expected exactly one library block")

			var gotIfaces, gotEnums, gotStructs, gotFwds int
			for _, decl := range idl.Libraries[0].Declarations {
				if decl.Interface != nil {
					gotIfaces++
				}
				if decl.Enum != nil {
					gotEnums++
				}
				if decl.Struct != nil {
					gotStructs++
				}
				if decl.InterfaceForewardDecl != "" {
					gotFwds++
				}
			}

			assert.Equal(t, want.interfaces, gotIfaces, "interface count")
			assert.Equal(t, want.enums, gotEnums, "enum count")
			assert.Equal(t, want.structs, gotStructs, "struct count")
			assert.Equal(t, want.forwards, gotFwds, "forward-declaration count")
		})
	}
}

// TestICoreWebView2MethodCount verifies the main interface always has 58 methods.
func TestICoreWebView2MethodCount(t *testing.T) {
	for filename := range idlCountTable {
		filename := filename
		t.Run(filename, func(t *testing.T) {
			data := idlFile(t, filename)
			idl, err := Parser.ParseBytes("", data)
			require.NoError(t, err)

			iface := findInterface(idl, "ICoreWebView2")
			require.NotNil(t, iface, "ICoreWebView2 must be present")
			assert.Equal(t, "IUnknown", iface.BaseClass)
			assert.Equal(t, 58, len(iface.Methods))
		})
	}
}

// TestInheritanceChain verifies that the latest IDL has the full 27-level chain.
func TestInheritanceChain(t *testing.T) {
	data := idlFile(t, "WebView2.1.0.2903.40.idl")
	idl, err := Parser.ParseBytes("", data)
	require.NoError(t, err)

	// Build name → baseClass map.
	bases := map[string]string{}
	for _, lib := range idl.Libraries {
		for _, decl := range lib.Declarations {
			if decl.Interface != nil {
				bases[decl.Interface.Name] = decl.Interface.BaseClass
			}
		}
	}

	// Walk from ICoreWebView2_27 back to ICoreWebView2.
	chain := []string{}
	cur := "ICoreWebView2_27"
	for cur != "" && cur != "IUnknown" {
		chain = append(chain, cur)
		cur = bases[cur]
	}

	// Must end at ICoreWebView2 and include every numbered version.
	require.Equal(t, "ICoreWebView2", chain[len(chain)-1])
	assert.Equal(t, 27, len(chain), "chain should have 27 entries (_27 through base)")

	// Spot-check a few links.
	assert.Equal(t, "ICoreWebView2_27", chain[0])
	assert.Equal(t, "ICoreWebView2_26", chain[1])
	assert.Equal(t, "ICoreWebView2", chain[26])
}

// findInterface returns the named InterfaceDeclaration or nil.
func findInterface(idl *types.IDL, name string) *types.InterfaceDeclaration {
	for _, lib := range idl.Libraries {
		for _, decl := range lib.Declarations {
			if decl.Interface != nil && decl.Interface.Name == name {
				return decl.Interface
			}
		}
	}
	return nil
}
