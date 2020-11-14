package generate

import (
	"io"
	"time"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/parser"
)

// AddSubcommand adds the `dev` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer) error {

	command := app.NewSubCommand("generate", "Code Generation Tools")

	// Backend API
	backendAPI := command.NewSubCommand("module", "Generates a JS module for the frontend to interface with the backend")

	// Quiet Init
	quiet := false
	backendAPI.BoolFlag("q", "Supress output to console", &quiet)

	backendAPI.Action(func() error {

		// Create logger
		logger := clilogger.New(w)
		logger.Mute(quiet)

		app.PrintBanner()

		logger.Print("Generating Javascript module for Go code...")

		// Start Time
		start := time.Now()

		p, err := parser.GenerateWailsFrontendPackage()
		if err != nil {
			return err
		}

		logger.Println("done.")
		logger.Println("")

		elapsed := time.Since(start)
		packages := p.Packages

		// Print report
		for _, pkg := range p.Packages {
			if pkg.ShouldGenerate() {
				logPackage(pkg, logger)
			}

		}

		logger.Println("%d packages parsed in %s.", len(packages), elapsed)

		return nil

	})
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

	// logger.Println("  Original Go Package Path:", pkg.Gopackage.PkgPath)
	// logger.Println("  Original Go Package Path:", pkg.Gopackage.PkgPath)
}
