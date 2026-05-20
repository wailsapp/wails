package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

// DockerMountsOptions holds options for the docker-mounts command.
type DockerMountsOptions struct{}

// HasCCOptions holds options for the has-cc command.
type HasCCOptions struct{}

// ToolHasCC checks if a C compiler (gcc or clang) is available in PATH.
// Outputs "true" or "false" for use in Taskfile sh: variables, replacing the
// bash-only `command -v gcc` pattern which fails on Windows.
func ToolHasCC(_ *HasCCOptions) error {
	DisableFooter = true
	_, gccErr := exec.LookPath("gcc")
	_, clangErr := exec.LookPath("clang")
	if gccErr == nil || clangErr == nil {
		fmt.Print("true")
	} else {
		fmt.Print("false")
	}
	return nil
}

// ToolDockerMounts outputs Docker volume mount flags for cross-compilation.
// It generates mounts for the Go module cache and any local replace directives
// in go.mod. This is a cross-platform replacement for the bash pipeline that
// was previously used in Taskfile templates.
//
// The project is assumed to be mounted at /app inside the container (matching
// the Taskfile `docker run -v "{{.ROOT_DIR}}:/app"` convention). Each -v flag
// is double-quoted so paths with spaces survive shell tokenisation in both
// POSIX sh and cmd.exe (go-task may use either on Windows).
func ToolDockerMounts(_ *DockerMountsOptions) error {
	DisableFooter = true

	var mounts []string

	// Add Go module cache mount. GOPATH may contain multiple entries
	// (os.PathListSeparator-separated); use only the first, falling back to ~/go.
	gopath := firstGOPATHEntry()
	if gopath != "" {
		hostPath := filepath.ToSlash(gopath)
		mounts = append(mounts, fmt.Sprintf(`-v "%s/pkg/mod:/go/pkg/mod"`, hostPath))
	}

	// Parse go.mod for local replace directives and add volume mounts.
	// Relative replace paths map under /app (the container project root);
	// absolute replace paths are mounted at the same path inside the container,
	// because Go inside the container resolves them literally.
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return fmt.Errorf("reading go.mod: %w", err)
	}
	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return fmt.Errorf("parsing go.mod: %w", err)
	}
	for _, r := range f.Replace {
		// Only handle local directory replacements (no version = local path)
		if r.New.Version != "" {
			continue
		}
		relPath := r.New.Path // forward-slash path as written in go.mod

		// Resolve absolute host path from the (possibly relative) replace path.
		hostAbsPath := relPath
		if !filepath.IsAbs(relPath) {
			abs, err := filepath.Abs(relPath)
			if err != nil {
				return fmt.Errorf("resolving replace path %q: %w", relPath, err)
			}
			hostAbsPath = abs
		}
		// Skip silently if the replace target doesn't exist locally —
		// a sibling may legitimately not be cloned on every machine, matching
		// the original Taskfile bash behaviour (`if [ -d "$path" ]`).
		if info, err := os.Stat(hostAbsPath); err != nil || !info.IsDir() {
			continue
		}
		hostDockerPath := filepath.ToSlash(hostAbsPath)

		// Compute the container-side destination.
		// Relative replace paths in go.mod are relative to the project root,
		// which maps to /app inside the container. path.Clean handles ".." correctly.
		var containerPath string
		if filepath.IsAbs(relPath) {
			// Windows drive-letter absolute paths (e.g. C:\vendor\lib) cannot be
			// mapped to valid Linux container destination paths — skip them.
			if len(relPath) >= 2 && relPath[1] == ':' {
				continue
			}
			// Mount at the same absolute path inside the container so Go
			// finds the module at the literal path written in go.mod.
			containerPath = hostDockerPath
		} else {
			containerPath = path.Clean("/app/" + relPath)
		}

		mounts = append(mounts, fmt.Sprintf(`-v "%s:%s:ro"`, hostDockerPath, containerPath))
	}

	fmt.Print(strings.Join(mounts, " "))
	return nil
}

// firstGOPATHEntry returns the first entry from GOPATH, or ~/go when unset.
func firstGOPATHEntry() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, "go")
	}
	entries := filepath.SplitList(gopath)
	if len(entries) == 0 {
		return ""
	}
	return entries[0]
}
