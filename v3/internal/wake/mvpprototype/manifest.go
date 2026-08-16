// Package mvpprototype is a throwaway vertical slice for validating the
// manifest-driven Wake cache. It is deliberately not a production API.
package mvpprototype

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/BurntSushi/toml"
)

// Manifest is the deliberately tiny subset needed by the prototype.
type Manifest struct {
	Project Project `toml:"project"`
}

type Project struct {
	Name        string `toml:"name"`
	ProductName string `toml:"product_name"`
	Identifier  string `toml:"identifier"`
	Version     string `toml:"version"`
	BinaryName  string `toml:"binary_name"`
}

func LoadManifest(root string) (Manifest, error) {
	var result Manifest
	path := filepath.Join(root, "wails.toml")
	if _, err := toml.DecodeFile(path, &result); err != nil {
		return Manifest{}, fmt.Errorf("wake mvp: read manifest: %w", err)
	}
	if result.Project.Name == "" || result.Project.ProductName == "" ||
		result.Project.Identifier == "" || result.Project.Version == "" {
		return Manifest{}, fmt.Errorf("wake mvp: [project] requires name, product_name, identifier, and version")
	}
	if result.Project.BinaryName == "" {
		result.Project.BinaryName = deriveBinaryName(result.Project.Name)
	}
	if result.Project.BinaryName == "" {
		return Manifest{}, fmt.Errorf("wake mvp: project name does not produce a binary name")
	}
	if result.Project.BinaryName != filepath.Base(result.Project.BinaryName) ||
		result.Project.BinaryName == "." || result.Project.BinaryName == ".." ||
		strings.ContainsAny(result.Project.BinaryName, `/\`) {
		return Manifest{}, fmt.Errorf("wake mvp: binary_name must be a plain file name, got %q", result.Project.BinaryName)
	}
	return result, nil
}

func deriveBinaryName(name string) string {
	var out strings.Builder
	dash := false
	for _, r := range strings.ToLower(strings.TrimSpace(name)) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			if dash && out.Len() > 0 {
				out.WriteByte('-')
			}
			dash = false
			out.WriteRune(r)
		default:
			dash = true
		}
	}
	return strings.Trim(out.String(), "-")
}
