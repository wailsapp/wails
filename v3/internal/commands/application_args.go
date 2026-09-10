package commands

import (
	"fmt"
	"github.com/google/shlex"
	"strings"
)

// resolveApplicationArguments parses the compatibility string exactly once.
// nil means inherit defaults; an explicitly empty vector clears them.
func resolveApplicationArguments(compatibility, passthrough []string) ([]string, error) {
	if len(compatibility) > 1 {
		return nil, fmt.Errorf("--appargs may be specified only once")
	}
	if len(compatibility) != 0 && passthrough != nil {
		return nil, fmt.Errorf("use either --appargs or arguments after --, not both")
	}
	args := passthrough
	if len(compatibility) != 0 {
		parsed, err := shlex.Split(compatibility[0])
		if err != nil {
			return nil, fmt.Errorf("invalid --appargs quoting: %w", err)
		}
		args = append([]string{}, parsed...)
	}
	for _, arg := range args {
		if strings.ContainsRune(arg, 0) {
			return nil, fmt.Errorf("application arguments must not contain NUL")
		}
	}
	return args, nil
}
