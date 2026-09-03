package bindings

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/internal/colour"
	"github.com/wailsapp/wails/v2/internal/shell"
	"github.com/wailsapp/wails/v2/pkg/commands/buildtags"
)

// Options for generating bindings
type Options struct {
	Filename          string
	Tags              []string
	ProjectDirectory  string
	Compiler          string
	GoModTidy         bool
	TsPrefix          string
	TsSuffix          string
	TsOutputType      string
	UseNullableSlices bool
}

// GenerateBindings generates bindings for the Wails project in the given ProjectDirectory.
// If no project directory is given then the current working directory is used.
func GenerateBindings(options Options) (string, error) {
	filename, _ := lo.Coalesce(options.Filename, "wailsbindings")
	if runtime.GOOS == "windows" {
		filename += ".exe"
	}

	// go build -tags bindings -o bindings.exe
	tempDir := os.TempDir()
	filename = filepath.Join(tempDir, filename)

	workingDirectory, _ := lo.Coalesce(options.ProjectDirectory, lo.Must(os.Getwd()))

	var stdout, stderr string
	var err error

	tags := append(options.Tags, "bindings")
	genModuleTags := lo.Without(tags, "desktop", "production", "debug", "dev")
	tagString := buildtags.Stringify(genModuleTags)

	if options.GoModTidy {
		stdout, stderr, err = shell.RunCommand(workingDirectory, options.Compiler, "mod", "tidy")
		if err != nil {
			return stdout, fmt.Errorf("%s\n%s\n%s", stdout, stderr, err)
		}
	}

	envBuild := os.Environ()
	envBuild = shell.SetEnv(envBuild, "GOOS", runtime.GOOS)
	envBuild = shell.SetEnv(envBuild, "GOARCH", runtime.GOARCH)
	// wailsbindings is executed on the build machine.
	// So, use the default C compiler, not the one set for cross compiling.
	envBuild = shell.RemoveEnv(envBuild, "CC")

	buildArgs := []string{"build", "-buildvcs=false", "-tags", tagString, "-o", filename}
	if runtime.GOOS == "windows" {
		// Go 1.25 switched to DWARF5 by default, which causes the internal
		// linker to emit malformed PE section headers on Windows — the binary
		// won't run ("not a valid Win32 application" / "not compatible with
		// the version of Windows"). Stripping debug info avoids DWARF5
		// entirely and is appropriate for a temporary tool binary.
		// See golang/go#75077, golang/go#75121, wails#4551, wails#4605.
		buildArgs = append(buildArgs, "-ldflags=-s -w")
	}
	stdout, stderr, err = shell.RunCommandWithEnv(envBuild, workingDirectory, options.Compiler, buildArgs...)
	if err != nil {
		return stdout, fmt.Errorf("%s\n%s\n%s", stdout, stderr, err)
	}

	if runtime.GOOS == "darwin" {
		// Remove quarantine attribute
		stdout, stderr, err = shell.RunCommand(workingDirectory, "/usr/bin/xattr", "-rc", filename)
		if err != nil {
			return stdout, fmt.Errorf("%s\n%s\n%s", stdout, stderr, err)
		}
	}

	defer func() {
		// Best effort removal of temp file
		_ = os.Remove(filename)
	}()

	// Set environment variables accordingly
	env := os.Environ()
	env = shell.SetEnv(env, "tsprefix", options.TsPrefix)
	env = shell.SetEnv(env, "tssuffix", options.TsSuffix)
	env = shell.SetEnv(env, "tsoutputtype", options.TsOutputType)
	if options.UseNullableSlices {
		env = shell.SetEnv(env, "usenullableslices", "true")
	}

	stdout, stderr, err = shell.RunCommandWithEnv(env, workingDirectory, filename)
	if err != nil {
		return stdout, fmt.Errorf("%s\n%s\n%s", stdout, stderr, err)
	}

	if stderr != "" {
		log.Println(colour.DarkYellow(stderr))
	}

	return stdout, nil
}
