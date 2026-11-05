package changelog

import "strings"

type commentBlock struct {
	active bool
}

func (c *commentBlock) skip(line string) bool {
	if !c.active && !strings.HasPrefix(line, "<!--") {
		return false
	}
	c.active = !strings.Contains(line, "-->")
	return true
}
