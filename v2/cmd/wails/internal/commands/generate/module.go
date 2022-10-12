package generate

import (
	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/pkg/commands/bindings"
	"github.com/wailsapp/wails/v2/pkg/commands/buildtags"
	"io"
)

// AddModuleCommand adds the `module` subcommand for the `generate` command
func AddModuleCommand(app *clir.Cli, parent *clir.Command, w io.Writer) error {

	command := parent.NewSubCommand("module", "Generate wailsjs modules")
	var tags string
	command.StringFlag("tags", "tags to pass to Go compiler (quoted and space separated)", &tags)

	command.Action(func() error {

		buildTags, err := buildtags.Parse(tags)
		if err != nil {
			return err
		}

		_, err = bindings.GenerateBindings(bindings.Options{
			Tags: buildTags,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return nil
}
