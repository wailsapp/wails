package idlreader

import (
	"fmt"
	"strings"
)

type token struct {
	text         string
	kind         byte
	line, column int
}

type lexer struct {
	input                string
	offset, line, column int
}

func (l *lexer) skip() {
	for l.offset < len(l.input) {
		rest := l.input[l.offset:]
		if rest[0] == ' ' || rest[0] == '\t' || rest[0] == '\r' || rest[0] == '\n' {
			l.offset++
			if rest[0] == '\n' {
				l.line++
				l.column = 1
			} else {
				l.column++
			}
			continue
		}
		if rest[0] == '#' || len(rest) > 1 && rest[:2] == "//" {
			n := strings.IndexByte(rest, '\n')
			if n < 0 {
				n = len(rest)
			}
			l.offset += n
			l.column += n
			continue
		}
		return
	}
}

func (l *lexer) next() (token, error) {
	l.skip()
	t := token{line: l.line, column: l.column}
	if l.offset == len(l.input) {
		return t, nil
	}
	rest := l.input[l.offset:]
	n := 1
	switch {
	case len(rest) >= 36 && validUUID(rest[:36]):
		t.kind = 'u'
		n = 36
	case rest[0] == '"':
		t.kind = 's'
		var closed bool
		n, closed = quotedLength(rest)
		if !closed {
			return t, fmt.Errorf("IDL:%d:%d: unterminated string", t.line, t.column)
		}
	case letter(rest[0]):
		t.kind = 'i'
		for n < len(rest) && (letter(rest[n]) || digit(rest[n]) || rest[n] == '_') {
			n++
		}
	case digit(rest[0]):
		n, t.kind = numberLength(rest)
	default:
		t.kind = 'p'
		if rest[0] < 33 || rest[0] > 126 {
			return t, fmt.Errorf("IDL:%d:%d: invalid character", t.line, t.column)
		}
	}
	t.text = rest[:n]
	l.offset += n
	if t.kind == 's' && strings.Contains(t.text, "\n") {
		l.line += strings.Count(t.text, "\n")
		l.column = n - strings.LastIndexByte(t.text, '\n')
	} else {
		l.column += n
	}
	return t, nil
}

func quotedLength(rest string) (int, bool) {
	n := 1
	for n < len(rest) {
		if rest[n] == '\\' && n+1 < len(rest) && rest[n+1] == '"' {
			n += 2
			continue
		}
		if rest[n] == '"' {
			return n + 1, true
		}
		n++
	}
	return n, false
}
func numberLength(rest string) (int, byte) {
	n := 1
	if len(rest) > 2 && rest[:2] == "0x" && hex(rest[2]) {
		n = 3
		for n < len(rest) && hex(rest[n]) {
			n++
		}
		return n, 'h'
	}
	for n < len(rest) && digit(rest[n]) {
		n++
	}
	return n, 'n'
}

func letter(c byte) bool { return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' }
func digit(c byte) bool  { return c >= '0' && c <= '9' }
func hex(c byte) bool    { return digit(c) || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F' }
func validUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if s[i] != '-' {
				return false
			}
		} else if !hex(s[i]) {
			return false
		}
	}
	return true
}
