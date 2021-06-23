package template

import (
	"embed"
	"io"
	"os"
	"path/filepath"

	"github.com/leaanthony/debme"

	"github.com/leaanthony/gosod"
	"github.com/wailsapp/wails/v2/internal/fs"

	"github.com/leaanthony/clir"
)

//go:embed base
var base embed.FS

func AddSubCommand(app *clir.Cli, parent *clir.Command, w io.Writer) {

	// command
	command := parent.NewSubCommand("template", "Generates a wails template")

	name := ""
	command.StringFlag("name", "The name of the template", &name)

	useLocalFilesAsFrontend := false
	command.BoolFlag("frontend", "This indicates that the current directory is a frontend project and should be used by the template", &useLocalFilesAsFrontend)

	// Quiet Init
	quiet := false
	command.BoolFlag("q", "Suppress output to console", &quiet)

	command.Action(func() error {

		// If the current directory is not empty, we create a new directory
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		templateDir := cwd
		empty, err := fs.DirIsEmpty(templateDir)
		if err != nil {
			return err
		}
		if !empty {
			templateDir = filepath.Join(cwd, name)
			err = fs.Mkdir(templateDir)
			if err != nil {
				return err
			}
		}

		// Create base template
		baseTemplate, err := debme.FS(base, "base")
		if err != nil {
			return err
		}
		g := gosod.New(baseTemplate)
		g.SetTemplateFilters([]string{".template"})

		err = os.Chdir(templateDir)
		if err != nil {
			return err
		}

		type templateData struct {
			Name        string
			Description string
		}

		err = g.Extract(templateDir, &templateData{
			Name: name,
		})
		if err != nil {
			return err
		}

		if useLocalFilesAsFrontend == false {
			return nil
		}

		// Remove frontend directory
		frontendDir := filepath.Join(templateDir, "frontend")
		err = os.RemoveAll(frontendDir)
		if err != nil {
			return err
		}

		err = fs.CopyDirExtended(cwd, frontendDir, []string{name})
		if err != nil {
			return err
		}

		//// Create logger
		//logger := clilogger.New(w)
		//logger.Mute(quiet)
		//
		//app.PrintBanner()
		//
		//logger.Print("Generating Javascript module for Go code...")
		//
		//// Start Time
		//start := time.Now()
		//
		//p, err := parser.GenerateWailsFrontendPackage()
		//if err != nil {
		//	return err
		//}
		//
		//logger.Println("done.")
		//logger.Println("")
		//
		//elapsed := time.Since(start)
		//packages := p.Packages
		//
		//// Print report
		//for _, pkg := range p.Packages {
		//	if pkg.ShouldGenerate() {
		//		generate.logPackage(pkg, logger)
		//	}
		//
		//}
		//
		//logger.Println("%d packages parsed in %s.", len(packages), elapsed)

		return nil

	})
}
