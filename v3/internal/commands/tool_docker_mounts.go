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
func ToolDockerMounts(_ *DockerMountsOptions) error {
	DisableFooter = true

	var mounts []string

	// Add Go module cache mount. GOPATH may contain multiple entries
	// (os.PathListSeparator-separated); use only the first, falling back to ~/go.
	gopath := firstGOPATHEntry()
	if gopath != "" {
		hostPath := filepath.ToSlash(gopath)
		mounts = append(mounts, fmt.Sprintf("-v '%s/pkg/mod:/go/pkg/mod'", hostPath))
	}

	// Parse go.mod for local replace directives and add volume mounts.
	// The container project root is /app; replace paths must be remapped accordingly.
	data, err := os.ReadFile("go.mod")
	if err == nil {
		f, err := modfile.Parse("go.mod", data, nil)
		if err == nil {
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
						continue
					}
					hostAbsPath = abs
				}
				if info, err := os.Stat(hostAbsPath); err != nil || !info.IsDir() {
					continue
				}
				hostDockerPath := filepath.ToSlash(hostAbsPath)

				// Compute the container-side destination.
				// Relative replace paths in go.mod are relative to the project root,
				// which maps to /app inside the container. path.Clean handles ".." correctly.
				var containerPath string
				if filepath.IsAbs(relPath) {
					// Absolute host paths can't be reliably remapped; use as-is.
					containerPath = filepath.ToSlash(relPath)
				} else {
					containerPath = path.Clean("/app/" + relPath)
				}

				mounts = append(mounts, fmt.Sprintf("-v '%s:%s:ro'", hostDockerPath, containerPath))
			}
		}
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
