package generate

import (
	"io"

	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/generate/template"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/parser"
)

// AddSubcommand adds the `generate` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer) error {

	command := app.NewSubCommand("generate", "Code Generation Tools")

	//AddModuleCommand(app, command, w)
	template.AddSubCommand(app, command, w)

	return nil
}

func logPackage(pkg *parser.Package, logger *clilogger.CLILogger) {

	logger.Println("Processed Go package '" + pkg.Gopackage.Name + "' as '" + pkg.Name + "'")
	for _, strct := range pkg.Structs() {
		logger.Println("")
		logger.Println("  Processed struct '" + strct.Name + "'")
		if strct.IsBound {
			for _, method := range strct.Methods {
				logger.Println("    Bound method '" + method.Name + "'")
			}
		}
		if strct.IsUsedAsData {
			for _, field := range strct.Fields {
				if !field.Ignored {
					logger.Print("    Processed ")
					if field.IsOptional {
						logger.Print("optional ")
					}
					logger.Println("field '" + field.Name + "' as '" + field.JSName() + "'")
				}
			}
		}
	}
	logger.Println("")
}
