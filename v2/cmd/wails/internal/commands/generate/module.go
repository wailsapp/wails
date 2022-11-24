package generate

import (
	"io"
	"os"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/internal/project"
	"github.com/wailsapp/wails/v2/pkg/commands/bindings"
	"github.com/wailsapp/wails/v2/pkg/commands/buildtags"
)

type generateFlags struct {
	tags string
}

// AddModuleCommand adds the `module` subcommand for the `generate` command
func AddModuleCommand(app *clir.Cli, parent *clir.Command, w io.Writer) error {

	command := parent.NewSubCommand("module", "Generate wailsjs modules")
	genFlags := generateFlags{}
	command.StringFlag("tags", "tags to pass to Go compiler (quoted and space separated)", &genFlags.tags)

	command.Action(func() error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		projectConfig, err := project.Load(cwd)
		if err != nil {
			return err
		}

		buildTags, err := buildtags.Parse(genFlags.tags)
		if err != nil {
			return err
		}

		_, err = bindings.GenerateBindings(bindings.Options{
			Tags:     buildTags,
			TsPrefix: projectConfig.Bindings.TsGeneration.Prefix,
			TsSuffix: projectConfig.Bindings.TsGeneration.Suffix,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return nil
}
