package render

import (
	"io"
	"path"
	"strconv"
	"strings"
	"unicode"

	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

// EventNamesDir is the directory, relative to the bindings root,
// that holds the module exporting typed event names.
const EventNamesDir = "events"

// EventNamesFile returns the path, relative to the bindings root,
// of the module that exports one named constant for each registered event.
func (renderer *Renderer) EventNamesFile() string {
	return path.Join(EventNamesDir, renderer.options.IndexFilename+renderer.ext)
}

type eventName struct {
	Ident string
	Name  string
}

// EventNames renders the given event map to w as a module
// exporting one constant per registered event,
// so that events can be referred to by (typed) identifier
// instead of by string literal.
func (renderer *Renderer) EventNames(w io.Writer, events *collect.EventMap) error {
	self := events.Imports.Self
	dir := "../" + self

	used := make(map[string]bool, len(events.Defs))
	names := make([]eventName, 0, len(events.Defs))
	for _, def := range events.Defs {
		ident := eventIdentifier(def.Name)
		for base, i := ident, 2; used[ident]; i++ {
			ident = base + strconv.Itoa(i)
		}
		used[ident] = true
		names = append(names, eventName{ident, def.Name})
	}

	return tmplEventNames.Execute(w, &struct {
		TS         bool
		HasEvents  bool
		DataPath   string
		DataFile   string
		CreatePath string
		Names      []eventName
	}{
		TS:         renderer.options.TS,
		HasEvents:  len(events.Defs) > 0,
		DataPath:   dir + "/eventdata.js",
		DataFile:   dir + "/" + renderer.EventDataFile(),
		CreatePath: dir + "/eventcreate.js",
		Names:      names,
	})
}

// reserved holds identifiers that cannot be used as export names.
var reserved = map[string]bool{
	"await": true, "break": true, "case": true, "catch": true, "class": true, "const": true,
	"continue": true, "debugger": true, "default": true, "delete": true, "do": true, "else": true,
	"enum": true, "export": true, "extends": true, "false": true, "finally": true, "for": true,
	"function": true, "if": true, "import": true, "in": true, "instanceof": true, "new": true,
	"null": true, "return": true, "super": true, "switch": true, "this": true, "throw": true,
	"true": true, "try": true, "typeof": true, "var": true, "void": true, "while": true,
	"with": true, "yield": true, "let": true, "static": true, "implements": true,
	"interface": true, "package": true, "private": true, "protected": true, "public": true,
	"arguments": true, "eval": true, "undefined": true,
}

// eventIdentifier derives a PascalCase JavaScript identifier from an event name,
// splitting at any character that is not a letter or digit
// (e.g. "user:login" becomes "UserLogin").
func eventIdentifier(name string) string {
	var b strings.Builder
	upper := true
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			upper = true
			continue
		}
		if upper {
			r = unicode.ToUpper(r)
			upper = false
		}
		b.WriteRune(r)
	}

	ident := b.String()
	if ident == "" {
		return "Event"
	}
	if first := []rune(ident)[0]; unicode.IsDigit(first) {
		ident = "Event" + ident
	}
	if reserved[ident] {
		ident = "Event" + ident
	}
	return ident
}
