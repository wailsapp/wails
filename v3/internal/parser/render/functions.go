package render

import (
	"bufio"
	"bytes"
	"fmt"
	"go/types"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/parser/collect"
)

// tmplFunctions holds a map of utility functions
// that should be available in every template.
var tmplFunctions = template.FuncMap{
	"isclass":  collect.IsClass,
	"jsdoc":    jsdoc,
	"jsid":     jsid,
	"jsimport": jsimport,
	"jsparam":  jsparam,
	"jstype":   RenderType,
	"jstypeq": func(typ types.Type, imports *collect.ImportMap, collector *collect.Collector, quoted bool) string {
		result, _ := renderType(typ, imports, collector, quoted)
		return result
	},
	"jsvalue":   RenderValue,
	"jsdefault": RenderDefault,
	"jscreate":  RenderCreate,
	"jscreated": renderConditionalCreate,
}

// jsimport formats an external import name
// by joining the name with its occurrence index.
// Names are modified even when the index is 0
// to avoid collisions with Go identifiers.
func jsimport(info collect.ImportInfo) string {
	return fmt.Sprintf("%s$%d", info.Name, info.Index)
}

// jsid escapes identifiers that match JS/TS reserved words
// by prepending a dollar sign.
func jsid(ident string) string {
	if _, reserved := slices.BinarySearch(reservedWords, ident); reserved {
		return "$" + ident
	}
	return ident
}

// jsparam renders the JS name of a parameter.
// Blank parameters are replaced with a dollar sign followed by the given index.
// Non-blank parameters are escaped by [jsid].
func jsparam(index int, param *collect.ParamInfo) string {
	if param.Blank {
		return "$" + strconv.Itoa(index)
	} else {
		return jsid(param.Name)
	}
}

// jsdoc splits the given comment into lines and rewrites it as follows:
//   - first, line terminators are stripped;
//   - then a line terminator, the indent string and ' * '
//     are prepended to each line;
//   - occurrences of the comment terminator '*/' are replaced with '* /'
//     to avoid accidentally terminating the surrounding comment.
//
// All lines thus modified are joined back together.
//
// The returned string can be inserted in a multiline JSDoc comment
// with the given indentation.
func jsdoc(comment string, indent string) string {
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
	"undefined",
	"var",
	"void",
	"while",
	"with",
	"yield",
}
