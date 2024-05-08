package templates

import (
	"bufio"
	"bytes"
	"embed"
	"slices"
	"strings"
	"text/template"
)

//go:embed *.tmpl
var templates embed.FS

var functions = template.FuncMap{
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

var BindingsJS = template.Must(template.New("bindings.js.tmpl").Funcs(functions).ParseFS(templates, "bindings.js.tmpl"))
var BindingsTS = template.Must(template.New("bindings.ts.tmpl").Funcs(functions).ParseFS(templates, "bindings.ts.tmpl"))

var ModelsJS = template.Must(template.New("models.js.tmpl").Funcs(functions).ParseFS(templates, "models.js.tmpl"))
var ModelsTS = template.Must(template.New("models.ts.tmpl").Funcs(functions).ParseFS(templates, "models.ts.tmpl"))
var InterfacesTS = template.Must(template.New("interfaces.ts.tmpl").Funcs(functions).ParseFS(templates, "interfaces.ts.tmpl"))

var IndexJS = template.Must(template.New("index.js.tmpl").Funcs(functions).ParseFS(templates, "index.js.tmpl"))
var IndexTS = template.Must(template.New("index.ts.tmpl").Funcs(functions).ParseFS(templates, "index.ts.tmpl"))

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
