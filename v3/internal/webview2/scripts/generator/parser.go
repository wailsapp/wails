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
