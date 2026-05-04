package commands

import (
	"fmt"
	"os"
	"os/exec"
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

	// Add Go module cache mount using GOPATH with fallback to ~/go
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			gopath = filepath.Join(home, "go")
		}
	}
	if gopath != "" {
		dockerPath := filepath.ToSlash(gopath)
		mounts = append(mounts, fmt.Sprintf("-v %s/pkg/mod:/go/pkg/mod", dockerPath))
	}

	// Parse go.mod for local replace directives and add volume mounts
	data, err := os.ReadFile("go.mod")
	if err == nil {
		f, err := modfile.Parse("go.mod", data, nil)
		if err == nil {
			for _, r := range f.Replace {
				// Only handle local directory replacements (no version = local path)
				if r.New.Version != "" {
					continue
				}
				path := r.New.Path
				if !filepath.IsAbs(path) {
					abs, err := filepath.Abs(path)
					if err != nil {
						continue
					}
					path = abs
				}
				if info, err := os.Stat(path); err == nil && info.IsDir() {
					dockerPath := filepath.ToSlash(path)
					mounts = append(mounts, fmt.Sprintf("-v %s:%s:ro", dockerPath, dockerPath))
				}
			}
		}
	}

	fmt.Print(strings.Join(mounts, " "))
	return nil
}
