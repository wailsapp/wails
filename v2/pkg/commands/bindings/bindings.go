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
	Filename         string
	Tags             []string
	ProjectDirectory string
	GoModTidy        bool
	TsPrefix         string
	TsSuffix         string
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
		stdout, stderr, err = shell.RunCommand(workingDirectory, "go", "mod", "tidy")
		if err != nil {
			return stdout, fmt.Errorf("%s\n%s\n%s", stdout, stderr, err)
		}
	}

	stdout, stderr, err = shell.RunCommand(workingDirectory, "go", "build", "-tags", tagString, "-o", filename)
	if err != nil {
		return stdout, fmt.Errorf("%s\n%s\n%s", stdout, stderr, err)
	}

	defer func() {
		// Best effort removal of temp file
		_ = os.Remove(filename)
	}()

	// Set environment variables accordingly
	env := os.Environ()
	env = shell.SetEnv(env, "tsprefix", options.TsPrefix)
	env = shell.SetEnv(env, "tssuffix", options.TsSuffix)

	stdout, stderr, err = shell.RunCommandWithEnv(env, workingDirectory, filename)
	if err != nil {
		return stdout, fmt.Errorf("%s\n%s\n%s", stdout, stderr, err)
	}

	if stderr != "" {
		log.Println(colour.DarkYellow(stderr))
	}

	return stdout, nil
}
