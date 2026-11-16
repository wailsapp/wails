package userpaths

import (
	"bufio"
	"io"
	"path/filepath"
	"strings"
)

func (p *Locations) readPersonalDirectories(r io.Reader, getenv func(string) string) {
	targets := p.personalSettings()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		for _, field := range targets {
			if key != field.key || filepath.IsAbs(getenv(key)) {
				continue
			}
			if path, ok := personalValue(strings.TrimSpace(value), p.Home); ok {
				*field.path = path
			}
		}
	}
}

// Parse the quoted data format without evaluating shell commands or other variables.
func personalValue(value, home string) (string, bool) {
	if len(value) < 2 || value[0] != '"' {
		return "", false
	}
	value = value[1:]
	var result strings.Builder
	if strings.HasPrefix(value, "$HOME/") || strings.HasPrefix(value, "$HOME\"") {
		result.WriteString(home)
		value = value[5:]
	}
	for i := 0; i < len(value); i++ {
		c := value[i]
		if c == '"' {
			path := result.String()
			if !filepath.IsAbs(path) {
				return "", false
			}
			return path, true
		}
		if c == '\\' && i+1 < len(value) && strings.ContainsRune("\\\"$`", rune(value[i+1])) {
			i++
			c = value[i]
		}
		result.WriteByte(c)
	}
	return "", false
}
