package buildtags

import (
	"errors"
	"github.com/samber/lo"
	"strings"
)

// Parse parses the given tags string and returns
// a cleaned slice of strings. Both comma and space delimeted
// tags are supported but not mixed. If mixed, an error is returned.
func Parse(tags string) ([]string, error) {

	if tags == "" {
		return nil, nil
	}

	var userTags []string
	separator := ""
	if strings.Contains(tags, ",") {
		separator = ","
	}
	if strings.Contains(tags, " ") {
		if separator != "" {
			return nil, errors.New("cannot use both space and comma separated values with `-tags` flag")
		}
		separator = " "
	}
	for _, tag := range strings.Split(tags, separator) {
		thisTag := strings.TrimSpace(tag)
		if thisTag != "" {
			userTags = append(userTags, thisTag)
		}
	}
	return userTags, nil
}

// Stringify converts the given tags slice to a string compatible
// with the go build -tags flag
func Stringify(tags []string) string {
	tags = lo.Map(tags, func(tag string, _ int) string {
		return strings.TrimSpace(tag)
	})
	return strings.Join(tags, ",")
}
