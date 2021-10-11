package generate

import (
	"fmt"
	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/cmd/wails/internal"
	"github.com/wailsapp/wails/v2/internal/shell"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// AddModuleCommand adds the `module` subcommand for the `generate` command
func AddModuleCommand(app *clir.Cli, parent *clir.Command, w io.Writer) error {

	command := parent.NewSubCommand("module", "Generate wailsjs modules")
	var tags string
	command.StringFlag("tags", "tags to pass to Go compiler (quoted and space separated)", &tags)

	command.Action(func() error {

		filename := "wailsbindings"
		if runtime.GOOS == "windows" {
			filename += ".exe"
		}
		// go build -tags bindings -o bindings.exe
		tempDir := os.TempDir()
		filename = filepath.Join(tempDir, filename)

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		tagList := internal.ParseUserTags(tags)
		tagList = append(tagList, "bindings")

		stdout, stderr, err := shell.RunCommand(cwd, "go", "build", "-tags", strings.Join(tagList, ","), "-o", filename)
		if err != nil {
			return fmt.Errorf("%s\n%s\n%s", stdout, stderr, err)
		}

		stdout, stderr, err = shell.RunCommand(cwd, filename)
		if err != nil {
			return fmt.Errorf("%s\n%s\n%s", stdout, stderr, err)
		}

		err = os.Remove(filename)
		if err != nil {
			return err
		}

		return nil
	})
	return nil
}
