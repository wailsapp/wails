//go:build linux && !android && !server

package application

import (
	"os"
	"strings"
)

// SystemLocale returns the system's configured locale as a BCP-47 language tag
// (e.g. "nb-NO", "en-US"). Reads from LANG/LC_ALL/LC_MESSAGES environment
// variables and normalizes to BCP-47 format.
func SystemLocale() string {
	// Check in order of specificity
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if val := os.Getenv(env); val != "" {
			return parsePosixLocale(val)
		}
	}
	return "en"
}

// parsePosixLocale converts a POSIX locale (e.g. "nb_NO.UTF-8") to BCP-47 ("nb-NO").
func parsePosixLocale(posix string) string {
	// Strip encoding (.UTF-8) and modifier (@euro)
	if i := strings.IndexByte(posix, '.'); i >= 0 {
		posix = posix[:i]
	}
	if i := strings.IndexByte(posix, '@'); i >= 0 {
		posix = posix[:i]
	}
	// Replace underscore with hyphen (nb_NO → nb-NO)
	posix = strings.ReplaceAll(posix, "_", "-")
	if posix == "" || posix == "C" || posix == "POSIX" {
		return "en"
	}
	return posix
}
