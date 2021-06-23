package generate

import (
	"io"
	"time"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/parser"
)

func AddModuleCommand(app *clir.Cli, parent *clir.Command, w io.Writer) {

	// Backend API
	backendAPI := parent.NewSubCommand("module", "Generates a JS module for the frontend to interface with the backend")

	// Quiet Init
	quiet := false
	backendAPI.BoolFlag("q", "Suppress output to console", &quiet)

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
}
