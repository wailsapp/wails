package generate

import (
	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/generate/template"
	"io"

	"github.com/leaanthony/clir"
)

// AddSubcommand adds the `generate` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer) error {

	command := app.NewSubCommand("generate", "Code Generation Tools")

	err := AddModuleCommand(app, command, w)
	if err != nil {
		return err
	}

	template.AddSubCommand(app, command, w)

	return nil
}
