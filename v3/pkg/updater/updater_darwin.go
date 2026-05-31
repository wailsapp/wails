//go:build darwin

package updater

import (
	"os"
	"path/filepath"
	"strings"
)

// bundleTarget returns the .app bundle path when exe lives inside one, or exe
// unchanged when it does not.
func bundleTarget(exe string) string {
	parts := strings.Split(filepath.Clean(exe), string(os.PathSeparator))
	for i, p := range parts {
		if strings.HasSuffix(p, ".app") {
			return string(os.PathSeparator) + filepath.Join(parts[1:i+1]...)
		}
	}
	return exe
}
