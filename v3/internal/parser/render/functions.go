package render

import (
	"bufio"
	"bytes"
	"slices"
	"strings"
	"text/template"
)

// tmplFunctions holds a map of utility functions
// that should be available in every template.
var tmplFunctions = template.FuncMap{
	"jsid": func(ident string) string {
		if _, reserved := slices.BinarySearch(reservedWords, ident); reserved {
			return "$" + ident
		}
		return ident
	},

	"jsdoc": func(comment string, indent string) string {
		var builder strings.Builder
		prefix := []byte("\n" + indent + " * ")

		scanner := bufio.NewScanner(bytes.NewReader([]byte(comment)))
		for scanner.Scan() {
			builder.Write(prefix)

			line := scanner.Bytes()

			// Escape comment terminators.
			for t := bytes.Index(line, commentTerminator); t >= 0; t = bytes.Index(line, commentTerminator) {
				builder.Write(line[:t+1])
				builder.WriteRune(' ')
				line = line[t+1:]
			}

			builder.Write(line)
		}

		return builder.String()
	},
}

func init() {
	// Ensure reserved words are sorted in ascending lexicographical order.
	slices.Sort(reservedWords)
}

var commentTerminator = []byte("*/")

// reservedWords is a list of JS + TS reserved or special meaning words.
// Keep in ascending lexicographical order for best startup performance.
var reservedWords = []string{
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
	"var",
	"void",
	"while",
	"with",
	"yield",
}
