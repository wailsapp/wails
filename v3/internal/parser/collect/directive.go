package collect

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// IsDirective returns true if the given comment
// is a directive of the form //wails: + directive.
func IsDirective(comment string, directive string) bool {
	if strings.HasPrefix(comment, "//wails:"+directive) {
		length := len("//wails:") + len(directive)
		if len(comment) == length {
			return true
		}

		next, _ := utf8.DecodeRuneInString(comment[length:])
		return unicode.IsSpace(next)
	}

	return false
}

// ParseDirective extracts the argument portion of a //wails: + directive comment.
func ParseDirective(comment string, directive string) string {
	rawArg := comment[len("//wails:")+len(directive):]

	if directive != "inject" {
		return strings.TrimSpace(rawArg)
	}

	// wails:inject requires special parsing:
	// do not trim all surrounding space, just the one space
	// immediately after the directive.
	_, wsize := utf8.DecodeRuneInString(rawArg)
	return rawArg[wsize:]
}
