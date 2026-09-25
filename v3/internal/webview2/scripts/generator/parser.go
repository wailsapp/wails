package generator

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"updater/generator/types"
)

var (
	idlLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"Comment", `(?:#|//)[^\n]*\n?`},
		{"String", `"(\\"|[^"])*"`},
		{"UUID", `[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}`},
		{"Hex", `0x[a-fA-F0-9]+`},
		{"Int", `[0-9]+`},
		{"Ident", `[a-zA-Z]\w*`},
		{"Number", `(?:@Int\.)?@Int`},
		{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
		{"Whitespace", `[ \t\n\r]+`},
	})
	Parser = participle.MustBuild[types.IDL](
		participle.UseLookahead(4),
		participle.Elide("Comment", "Whitespace"),
		participle.Lexer(idlLexer),
	)
)

func ParseIDL(idlData []byte) ([]*types.GeneratedFile, error) {

	idl, err := Parser.ParseBytes("", idlData)
	if err != nil {
		return nil, err
	}

	err = idl.Process()
	if err != nil {
		return nil, err
	}

	generatedFiles, err := idl.Generate()
	if err != nil {
		return nil, err
	}
	return generatedFiles, nil
}

// ParseIDLWithTests is ParseIDL plus the generated test files (per-interface
// marshalling tests and the package-wide ABI test file).
func ParseIDLWithTests(idlData []byte) ([]*types.GeneratedFile, error) {
	idl, err := Parser.ParseBytes("", idlData)
	if err != nil {
		return nil, err
	}
	if err := idl.Process(); err != nil {
		return nil, err
	}
	files, err := idl.Generate()
	if err != nil {
		return nil, err
	}
	tests, err := idl.GenerateTests()
	if err != nil {
		return nil, err
	}
	return append(files, tests...), nil
}

// InterfaceNames parses the IDL and returns the names of all fully declared
// interfaces (forward references excluded), in declaration order. This is
// the inventory the capability table must cover.
func InterfaceNames(idlData []byte) ([]string, error) {
	idl, err := Parser.ParseBytes("", idlData)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, lib := range idl.Libraries {
		for _, d := range lib.Declarations {
			if d.Interface != nil {
				names = append(names, d.Interface.Name)
			}
		}
	}
	return names, nil
}

// InterfaceMethods parses the IDL and returns each declared interface's own
// method names (inherited methods excluded), used to diff SDK versions for
// changelog entries.
func InterfaceMethods(idlData []byte) (map[string][]string, error) {
	idl, err := Parser.ParseBytes("", idlData)
	if err != nil {
		return nil, err
	}
	out := map[string][]string{}
	for _, lib := range idl.Libraries {
		for _, d := range lib.Declarations {
			if d.Interface == nil {
				continue
			}
			var methods []string
			for _, m := range d.Interface.Methods {
				methods = append(methods, string(m.Name))
			}
			out[d.Interface.Name] = methods
		}
	}
	return out, nil
}
