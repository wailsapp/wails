// Package sourceignore applies project .gitignore files to build inputs and watches.
// It intentionally reads neither the Git index nor machine-specific exclude files.
package sourceignore

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Matcher reports exclusions for paths split into project-relative components.
// Once a parent directory is excluded, child negations cannot include its files.
type Matcher interface {
	Match(parts []string, directory bool) bool
}

type rule struct {
	base                         string
	expression                   *regexp.Regexp
	basename, directory, include bool
}
type matcher struct{ rules []rule }

// New creates a matcher from rules relative to the project root.
func New(lines ...string) Matcher {
	m := &matcher{}
	m.add("", lines)
	return m
}

// Read loads project rules, with deeper .gitignore files taking precedence.
// Excluded directories and symbolic links are not traversed.
func Read(root string) (Matcher, error) {
	m := &matcher{}
	err := filepath.WalkDir(root, func(current string, entry fs.DirEntry, err error) error {
		if err != nil {
			// A sibling build step may delete its staging directory mid-walk.
			if errors.Is(err, fs.ErrNotExist) && current != root {
				return nil
			}
			return err
		}
		if !entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, current)
		if err != nil {
			return err
		}
		if relative != "." && (entry.Name() == ".git" || m.Match(strings.Split(filepath.ToSlash(relative), "/"), true)) {
			return filepath.SkipDir
		}
		name := filepath.Join(current, ".gitignore")
		info, err := os.Lstat(name)
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		file, err := os.Open(name)
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 4096), 1024*1024)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		base := ""
		if relative != "." {
			base = filepath.ToSlash(relative) + "/"
		}
		m.add(base, lines)
		return nil
	})
	return m, err
}

func (m *matcher) add(base string, lines []string) {
	for index, line := range lines {
		if index == 0 {
			line = strings.TrimPrefix(line, "\ufeff")
		}
		line = strings.TrimSuffix(line, "\r")
		// A space is significant only when preceded by an odd run of backslashes.
		for strings.HasSuffix(line, " ") {
			backslashes := 0
			for i := len(line) - 2; i >= 0 && line[i] == '\\'; i-- {
				backslashes++
			}
			if backslashes%2 == 1 {
				break
			}
			line = line[:len(line)-1]
		}
		if line == "" || line[0] == '#' {
			continue
		}
		r := rule{base: wildcardText(base)}
		if line[0] == '!' {
			r.include = true
			line = line[1:]
		}
		r.directory = strings.HasSuffix(line, "/")
		line = strings.TrimSuffix(line, "/")
		r.basename = !strings.Contains(line, "/")
		line = strings.TrimPrefix(line, "/")
		if line == "" {
			continue
		}
		expression, err := regexp.Compile("^(?:" + glob(wildcardText(line)) + ")$")
		if err != nil {
			continue
		}
		r.expression = expression
		m.rules = append(m.rules, r)
	}
}

func (m *matcher) Match(parts []string, directory bool) bool {
	for depth := 1; depth <= len(parts); depth++ {
		current := wildcardText(strings.Join(parts[:depth], "/"))
		basename := wildcardText(parts[depth-1])
		isDirectory := depth < len(parts) || directory
		excluded := false
		for index := len(m.rules) - 1; index >= 0; index-- {
			r := &m.rules[index]
			if r.directory && !isDirectory || !strings.HasPrefix(current, r.base) {
				continue
			}
			candidate := strings.TrimPrefix(current, r.base)
			if candidate == "" {
				continue
			}
			if r.basename {
				candidate = basename
			}
			if r.expression.MatchString(candidate) {
				excluded = !r.include
				break
			}
		}
		if excluded {
			return true
		}
	}
	return false
}

// glob compiles documented Git wildcard syntax to Go's linear-time regexp engine.
func glob(pattern string) string {
	var expression strings.Builder
	segments := strings.Split(pattern, "/")
	for i, segment := range segments {
		if segment == "**" && len(segments) > 1 {
			if i == len(segments)-1 {
				expression.WriteString("(?s:.+)")
			} else {
				expression.WriteString("(?:[^/]+/)*")
			}
			continue
		}
		expression.WriteString(componentGlob(segment))
		if i < len(segments)-1 {
			expression.WriteByte('/')
		}
	}
	return expression.String()
}

func componentGlob(pattern string) string {
	var expression strings.Builder
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '*':
			expression.WriteString("[^/]*")
		case '?':
			expression.WriteString("[^/]")
		case '\\':
			i++
			if i == len(pattern) {
				return "["
			}
			expression.WriteString(regexp.QuoteMeta(pattern[i : i+1]))
		case '[':
			class, end := characterClass(pattern, i)
			if end < 0 {
				return "["
			}
			expression.WriteString(class)
			i = end
		default:
			expression.WriteString(regexp.QuoteMeta(pattern[i : i+1]))
		}
	}
	return expression.String()
}

func characterClass(pattern string, start int) (string, int) {
	var expression strings.Builder
	expression.WriteByte('[')
	i := start + 1
	if i < len(pattern) && (pattern[i] == '!' || pattern[i] == '^') {
		expression.WriteByte('^')
		i++
	}
	if i < len(pattern) && pattern[i] == ']' {
		expression.WriteString(`\]`)
		i++
	}
	for ; i < len(pattern); i++ {
		if pattern[i] == ']' {
			expression.WriteByte(']')
			return expression.String(), i
		}
		if strings.HasPrefix(pattern[i:], "[:") {
			end := strings.Index(pattern[i+2:], ":]")
			if end >= 0 {
				end += i + 4
				expression.WriteString(pattern[i:end])
				i = end - 1
				continue
			}
		}
		if pattern[i] == '\\' {
			i++
			if i == len(pattern) {
				break
			}
			expression.WriteString(`\`)
		}
		expression.WriteByte(pattern[i])
	}
	return "", -1
}

// Git wildcards count filename bytes, including UTF-8 continuation bytes.
func wildcardText(value string) string {
	for i := 0; i < len(value); i++ {
		if value[i] < 128 {
			continue
		}
		var text strings.Builder
		text.Grow(len(value) * 2)
		for j := 0; j < len(value); j++ {
			text.WriteRune(rune(value[j]))
		}
		return text.String()
	}
	return value
}
