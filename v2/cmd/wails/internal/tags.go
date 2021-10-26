package internal

import "strings"

// ParseUserTags takes the string form of tags and converts to a slice of strings
func ParseUserTags(tagString string) []string {
	userTags := make([]string, 0)
	for _, tag := range strings.Split(tagString, " ") {
		thisTag := strings.TrimSpace(tag)
		if thisTag != "" {
			userTags = append(userTags, thisTag)
		}
	}
	return userTags
}
