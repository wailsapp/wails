package cache

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// snapshotGoEmbed follows directives independently of the source tree's
// extension, directory and gitignore filters. Generated frontend output and
// user-owned resources are both compiler inputs. Pattern parsing is cached by
// source content, but matches are rediscovered on every snapshot so additions
// and deletions cannot reuse an old binary.
func (c *Cache) snapshotGoEmbed(source, sourceDigest string) (string, error) {
	if c.index.GoEmbeds == nil {
		c.index.GoEmbeds = make(map[string][]string)
	}
	patterns, ok := c.index.GoEmbeds[sourceDigest]
	if !ok {
		data, err := os.ReadFile(source)
		if err != nil {
			return "", err
		}
		if bytes.Contains(data, []byte("//go:embed")) {
			file, err := parser.ParseFile(token.NewFileSet(), source, data, parser.ParseComments)
			if err != nil {
				return "", err
			}
			for _, group := range file.Comments {
				for _, comment := range group.List {
					text, found := strings.CutPrefix(comment.Text, "//go:embed")
					if !found || text == "" || (text[0] != ' ' && text[0] != '\t') {
						continue
					}
					values, err := parseEmbedPatterns(text)
					if err != nil {
						return "", fmt.Errorf("%s: go:embed: %w", source, err)
					}
					patterns = append(patterns, values...)
				}
			}
		}
		c.index.GoEmbeds[sourceDigest] = patterns
		c.dirty = true
	}
	if len(patterns) == 0 {
		return "", nil
	}
	files := make(map[string]bool)
	for _, pattern := range patterns {
		// Including hidden files conservatively for both forms is safe: it can
		// cause an extra rebuild, but cannot omit anything Go embeds with all:.
		pattern = strings.TrimPrefix(pattern, "all:")
		if pattern == "." || !fs.ValidPath(pattern) || strings.Contains(pattern, "\\") {
			return "", fmt.Errorf("%s: invalid go:embed pattern %q", source, pattern)
		}
		matches, err := filepath.Glob(filepath.Join(filepath.Dir(source), filepath.FromSlash(pattern)))
		if err != nil {
			return "", err
		}
		for _, match := range matches {
			err := filepath.WalkDir(match, func(path string, entry fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !entry.IsDir() {
					files[path] = true
				}
				return nil
			})
			if err != nil {
				return "", err
			}
		}
	}
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return c.SnapshotFiles("go-embed", paths...)
}

func parseEmbedPatterns(text string) ([]string, error) {
	var patterns []string
	for text = strings.TrimSpace(text); text != ""; text = strings.TrimSpace(text) {
		var pattern string
		if text[0] == '"' || text[0] == '`' {
			quote := text[0]
			end := 1
			for end < len(text) && text[end] != quote {
				if quote == '"' && text[end] == '\\' {
					end++
				}
				end++
			}
			if end >= len(text) {
				return nil, fmt.Errorf("unterminated quoted pattern")
			}
			var err error
			pattern, err = strconv.Unquote(text[:end+1])
			if err != nil {
				return nil, err
			}
			text = text[end+1:]
			if text != "" && text[0] != ' ' && text[0] != '\t' {
				return nil, fmt.Errorf("patterns must be separated by whitespace")
			}
		} else {
			end := strings.IndexFunc(text, unicode.IsSpace)
			if end < 0 {
				end = len(text)
			}
			pattern, text = text[:end], text[end:]
		}
		patterns = append(patterns, pattern)
	}
	return patterns, nil
}
