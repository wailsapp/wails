package generate

import (
	"io"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/pkg/parser"
)

// AddSubcommand adds the `dev` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer) error {

	command := app.NewSubCommand("generate", "Code Generation Tools")

	// Backend API
	backendAPI := command.NewSubCommand("api", "Generates a JS module for the frontend to interface with the backend")

	backendAPI.Action(func() error {
		return parser.GenerateWailsFrontendPackage()
	})
	return nil
}
