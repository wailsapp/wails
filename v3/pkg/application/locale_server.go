//go:build server

package application

import "os"
import "strings"

// SystemLocale returns the system's configured locale as a BCP-47 language tag.
// In server mode, reads from environment variables.
func SystemLocale() string {
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if val := os.Getenv(env); val != "" {
			return parsePosixLocaleBCP47(val)
		}
	}
	return "en"
}

// parsePosixLocaleBCP47 converts a POSIX locale (e.g. "nb_NO.UTF-8@euro") to BCP-47 ("nb-NO").
func parsePosixLocaleBCP47(posix string) string {
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
