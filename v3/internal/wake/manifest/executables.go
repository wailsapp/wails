package manifest

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func validateExecutables(config Config) error {
	if err := validateExecutable(config.Root, "project", config.Project.Package, config.Project.Kind, config.Project.CommitSymbol, config.Project.UPXLevel); err != nil {
		return err
	}
	if (config.Project.VersionFile == "") != (config.Project.VersionSymbol == "") {
		return fieldValidationError("project.version_file", "version_file and version_symbol must be configured together")
	}
	if config.Project.VersionFile != "" {
		if !token.IsIdentifier(config.Project.VersionSymbol) {
			return fieldValidationError("project.version_symbol", "must be a Go declaration name")
		}
		if _, err := ResolveProjectPath(config.Root, "version source", config.Project.VersionFile, true); err != nil {
			return fieldValidationCause("project.version_file", err, "%v", err)
		}
	}
	names := map[string]bool{strings.ToLower(config.Project.BinaryName): true}
	for _, binary := range config.Binaries {
		field := fmt.Sprintf("binaries[%q]", binary.Name)
		if !slugPattern.MatchString(binary.Name) || names[strings.ToLower(binary.Name)] {
			return fieldValidationError(field, "name must be a unique lowercase slug distinct from project.binary_name")
		}
		names[strings.ToLower(binary.Name)] = true
		if binary.Package == "" {
			return fieldValidationError(field+".package", "a local Go package is required")
		}
		if err := validateExecutable(config.Root, field, binary.Package, binary.Kind, binary.CommitSymbol, binary.UPXLevel); err != nil {
			return err
		}
	}
	return nil
}

func validateExecutable(root, field, entry, kind, commit string, upx int) error {
	if kind != "" && kind != "gui" && kind != "cli" && kind != "server" {
		return fieldValidationError(field+".kind", "must be gui, cli or server")
	}
	if upx < 0 || upx > 9 {
		return fieldValidationError(field+".upx_level", "must be 0 (disabled) or 1 through 9")
	}
	if entry != "" && entry != "." {
		clean := filepath.ToSlash(filepath.Clean(entry))
		if strings.ContainsAny(entry, "\x00\r\n") || strings.HasPrefix(entry, "-") || strings.Contains(clean, "...") {
			return fieldValidationError(field+".package", "must name one local Go package")
		}
		path, err := ResolveProjectPath(root, "Go package", entry, true)
		if err != nil {
			return fieldValidationCause(field+".package", err, "%v", err)
		}
		info, err := os.Stat(path)
		if err != nil {
			return fieldValidationCause(field+".package", err, "%v", err)
		}
		if !info.IsDir() {
			return fieldValidationError(field+".package", "must name a local Go package directory")
		}
	}
	if commit != "" {
		pkg, symbol, found := strings.LastIndex(commit, "."), "", false
		if pkg > 0 {
			symbol, found = commit[pkg+1:], true
		}
		if !found || !token.IsIdentifier(symbol) || strings.ContainsAny(commit, " \t\r\n\x00'\"") {
			return fieldValidationError(field+".commit_symbol", "must be a qualified Go string variable such as main.Commit")
		}
	}
	return nil
}
