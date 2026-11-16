package installerbuild

import (
	"slices"
	"strings"
)

// CompilerEnvironment gives only the macOS NSIS child a UTF-8 locale. Its Unicode
// compiler can abort under the C locale (upstream NSIS issue 1165).
func CompilerEnvironment(host string, environment []string) []string {
	result := slices.Clone(environment)
	if host != "darwin" {
		return result
	}
	values := map[string]string{}
	for _, entry := range result {
		key, value, found := strings.Cut(entry, "=")
		if found {
			values[key] = value
		}
	}
	effective := values["LC_ALL"]
	if effective == "" {
		effective = values["LC_CTYPE"]
	}
	if effective == "" {
		effective = values["LANG"]
	}
	encoding := effective
	if _, suffix, found := strings.Cut(encoding, "."); found {
		encoding = suffix
	}
	encoding, _, _ = strings.Cut(encoding, "@")
	if strings.EqualFold(strings.ReplaceAll(encoding, "-", ""), "UTF8") {
		return result
	}
	for _, key := range []string{"LC_ALL", "LANG"} {
		found := false
		for index, entry := range result {
			if strings.HasPrefix(entry, key+"=") {
				result[index] = key + "=en_US.UTF-8"
				found = true
			}
		}
		if !found {
			result = append(result, key+"=en_US.UTF-8")
		}
	}
	return result
}
