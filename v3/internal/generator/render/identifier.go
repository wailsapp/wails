package render

import (
	"slices"
)

// jsid escapes identifiers that match JS/TS reserved words
// by prepending a dollar sign.
func jsid(ident string) string {
	if _, reserved := slices.BinarySearch(protectedWords, ident); reserved {
		return "$" + ident
	}
	return ident
}

func init() {
	// Ensure reserved words are sorted in ascending lexicographical order.
	slices.Sort(protectedWords)
}

// protectedWords is a list of JS + TS words that are either reserved
// or have special meaning. Keep in ascending lexicographical order
// for best startup performance.
var protectedWords = []string{
	"JSON",
	"Object",
	"any",
	"arguments",
	"as",
	"async",
	"await",
	"boolean",
	"break",
	"case",
	"catch",
	"class",
	"const",
	"constructor",
	"continue",
	"debugger",
	"declare",
	"default",
	"delete",
	"do",
	"else",
	"enum",
	"export",
	"extends",
	"false",
	"finally",
	"for",
	"from",
	"function",
	"get",
	"if",
	"implements",
	"import",
	"in",
	"instanceof",
	"interface",
	"let",
	"module",
	"namespace",
	"new",
	"null",
	"number",
	"of",
	"package",
	"private",
	"protected",
	"public",
	"require",
	"return",
	"set",
	"static",
	"string",
	"super",
	"switch",
	"symbol",
	"this",
	"throw",
	"true",
	"try",
	"type",
	"typeof",
	"undefined",
	"var",
	"void",
	"while",
	"with",
	"yield",
}
