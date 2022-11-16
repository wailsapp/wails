package generate

import (
	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/pkg/commands/bindings"
	"github.com/wailsapp/wails/v2/pkg/commands/buildtags"
	"io"
)

type generateFlags struct {
	tags    string
	prefix  string
	suffix string
}

// AddModuleCommand adds the `module` subcommand for the `generate` command
func AddModuleCommand(app *clir.Cli, parent *clir.Command, w io.Writer) error {

	command := parent.NewSubCommand("module", "Generate wailsjs modules")
	genFlags := generateFlags{}
	command.StringFlag("tags", "tags to pass to Go compiler (quoted and space separated)", &genFlags.tags)

	command.StringFlag("tsprefix", "prefix for generated typescript entities", &genFlags.prefix)
	command.StringFlag("tssuffix", "suffix for generated typescript entities", &genFlags.suffix)

	command.Action(func() error {

		buildTags, err := buildtags.Parse(genFlags.tags)
		if err != nil {
			return err
		}

		_, err = bindings.GenerateBindings(bindings.Options{
			Tags: buildTags,
			TsPrefix: genFlags.prefix,
			TsSuffix: genFlags.suffix,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return nil
}
