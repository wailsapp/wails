package generate

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/common"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/internal/shell"
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

		tagList, err := common.ParseUserTags(tags)
		if err != nil {
			return err
		}
		tagList = append(tagList, "bindings")
		genModuleTags := lo.Without(tagList, "desktop", "production", "debug", "dev")
		genModuleTagsString := strings.Join(genModuleTags, ",")

		stdout, stderr, err := shell.RunCommand(cwd, "go", "build", "-tags", genModuleTagsString, "-o", filename)
		if err != nil {
			return fmt.Errorf("%s\n%s\n%s", stdout, stderr, err)
		}

		stdout, stderr, err = shell.RunCommand(cwd, filename)
		println(stdout)
		println(stderr)
		if err != nil {
			return fmt.Errorf("%s\n%s\n%s", stdout, stderr, err)
		}

		// Best effort removal of temp file
		_ = os.Remove(filename)

		return nil
	})
	return nil
}
