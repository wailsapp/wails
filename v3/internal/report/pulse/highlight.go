package pulse

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// keywordRE matches FAIL, PASS, ERROR (and lowercase variants) as whole
// words. These keywords flood `go test` output and similar logs; bolding
// them lets the eye land on the failure context immediately rather than
// having to scan plain paragraphs of text.
var keywordRE = regexp.MustCompile(`\b(FAIL|PASS|ERROR|fail|pass|error)\b`)

// fileRefRE matches "file:line[:col]" file references — including bare-name
// references like "analyser_test.go:142" that `go test` emits, as well as
// path-prefixed ones like "internal/foo/bar.go:88:6". The recognised
// extension list ([a-zA-Z]{1,5}) prevents most false positives from prose
// (e.g. "42:99" or "word:1").
//
// Examples it matches:
//
//	internal/foo/bar.go:142
//	pkg/baz.go:88:6
//	./cmd/main.go:1:1
//	analyser_test.go:142
//	src\windows\winapi.cc:50  — also matches backslash separators
//
// Examples it leaves alone:
//
//	exit status 1
//	Service13:142   — no extension
//	2026-06-03 10:49:12  — no extension
var fileRefRE = regexp.MustCompile(
	`(?:^|[\s\[(])((?:\./)?(?:[\w.\-]+[/\\])*[\w\-]+\.[a-zA-Z]{1,5})(:\d+(?::\d+)?)`,
)

// highlightBody decorates one body line of a failure panel: bolds the
// FAIL/PASS/ERROR keywords and wraps file:line references in OSC 8
// hyperlinks pointing at the resolved absolute path. cwd is the working
// directory used to absolutise relative paths.
//
// The output's *visible* width is unchanged — we only insert SGR and OSC 8
// escape sequences, both of which `visibleWidth` skips. That means the
// failure panel's width math (top/bottom border rule, body padding) keeps
// working without changes.
func (s *styler) highlightBody(line, cwd string) string {
	// File references first — they often contain alphabetic identifiers that
	// could overlap with the keyword regex (e.g. "error.go"), and we don't
	// want to bold "error" inside a path.
	line = fileRefRE.ReplaceAllStringFunc(line, func(match string) string {
		// Re-split into the leading whitespace (if any), the file, and the
		// location tail — the outer regex captured these but ReplaceAllString
		// loses the submatches when called through the Func variant.
		sub := fileRefRE.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		leadEnd := strings.Index(match, sub[1])
		lead := match[:leadEnd]
		file, loc := sub[1], sub[2]
		abs := file
		if cwd != "" && !filepath.IsAbs(file) {
			abs = filepath.Join(cwd, file)
		}
		uri := "file://" + abs
		return lead + s.link(uri, file+loc)
	})

	// Then keywords. Wrap each match in bold + the appropriate accent.
	line = keywordRE.ReplaceAllStringFunc(line, func(match string) string {
		upper := strings.ToUpper(match)
		switch upper {
		case "PASS":
			return s.bold(s.fg(Success, match))
		case "FAIL", "ERROR":
			return s.bold(s.fg(Failure, match))
		}
		return match
	})

	return line
}

// formatExitCode renders an exit-code line with the number itself bolded
// so the panel's "status exited 1" reads "status exited **1**" — the
// number is what the eye is looking for.
func (s *styler) formatExitCode(code int) string {
	return fmt.Sprintf("exited %s", s.bold(fmt.Sprintf("%d", code)))
}
