package render

import (
	"bufio"
	"bytes"
	"go/ast"
	"strings"
	"unicode"

	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

// hasdoc checks whether the given comment group contains actual doc comments.
func hasdoc(group *ast.CommentGroup) bool {
	if group == nil {
		return false
	}

	// TODO: this is horrible, make it more efficient?
	return strings.ContainsFunc(group.Text(), func(r rune) bool { return !unicode.IsSpace(r) })
}

var commentTerminator = []byte("*/")

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
	prefix := []byte(Newline + indent + " * ")

	scanner := bufio.NewScanner(bytes.NewReader([]byte(comment)))
	for scanner.Scan() {
		line := scanner.Bytes()

		// Prepend prefix.
		builder.Write(prefix)

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

// jsdocline removes all newlines in the given comment
// and escapes comment terminators using the same strategy as jsdoc.
func jsdocline(comment string) string {
	var builder strings.Builder

	scanner := bufio.NewScanner(bytes.NewReader([]byte(comment)))
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			// Skip empty lines.
			continue
		}

		// Prepend space to separate lines.
		builder.WriteRune(' ')

		// Escape comment terminators.
		for t := bytes.Index(line, commentTerminator); t >= 0; t = bytes.Index(line, commentTerminator) {
			builder.Write(line[:t+1])
			builder.WriteRune(' ')
			line = line[t+1:]
		}

		builder.Write(line)
	}

	// Return resulting string, but skip initial space.
	return builder.String()[1:]
}

// isjsdocid returns true if the given string is a valid ECMAScript identifier,
// excluding unicode escape sequences. This is the property name format supported by JSDoc.
func isjsdocid(name string) bool {
	for i, r := range name {
		if i == 0 && !id_start(r) && r != '$' && r != '_' {
			return false
		} else if i > 0 && !id_continue(r) && r != '$' {
			return false
		}
	}
	return true
}

// isjsdocobj returns true if all field names in the given model
// are valid jsdoc property names.
func isjsdocobj(model *collect.ModelInfo) bool {
	if len(model.Fields) == 0 {
		return false
	}
	for _, decl := range model.Fields {
		for _, field := range decl {
			if !isjsdocid(field.JsonName) {
				return false
			}
		}
	}
	return true
}

// id_start returns true if the given rune is in the ID_Start category
// according to UAX#31 (https://unicode.org/reports/tr31/).
func id_start(r rune) bool {
	return (unicode.IsLetter(r) ||
		unicode.Is(unicode.Nl, r) ||
		unicode.Is(unicode.Other_ID_Start, r)) && !unicode.Is(unicode.Pattern_Syntax, r) && !unicode.Is(unicode.Pattern_White_Space, r)
}

// id_continue returns true if the given rune is in the ID_Continue category
// according to UAX#31 (https://unicode.org/reports/tr31/).
func id_continue(r rune) bool {
	return (id_start(r) ||
		unicode.Is(unicode.Mn, r) ||
		unicode.Is(unicode.Mc, r) ||
		unicode.Is(unicode.Nd, r) ||
		unicode.Is(unicode.Pc, r) ||
		unicode.Is(unicode.Other_ID_Continue, r)) && !unicode.Is(unicode.Pattern_Syntax, r) && !unicode.Is(unicode.Pattern_White_Space, r)
}
