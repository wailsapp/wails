package common

import (
	"errors"
	"strings"
)

// ParseUserTags parses the user tags from the command line
func ParseUserTags(tags string) ([]string, error) {

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
