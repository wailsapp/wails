package packagesign

import (
	"os"
	"path/filepath"
	"strings"
)

// KeyFile distinguishes exported files from legacy GnuPG key IDs. Explicit
// paths and conventional key extensions remain files even when unreadable or
// missing, so a broken key cannot silently select the external signer.
func KeyFile(root, reference string) (string, bool) {
	if reference == "" {
		return "", false
	}
	path := reference
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	_, err := os.Lstat(path)
	if err == nil || !os.IsNotExist(err) || filepath.IsAbs(reference) || strings.ContainsAny(reference, `/\`) {
		return path, true
	}
	switch strings.ToLower(filepath.Ext(reference)) {
	case ".asc", ".gpg", ".pgp", ".key", ".pem":
		return path, true
	}
	return "", false
}
