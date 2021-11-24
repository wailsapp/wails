package template

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/leaanthony/debme"

	"github.com/leaanthony/gosod"
	"github.com/wailsapp/wails/v2/internal/fs"

	"github.com/leaanthony/clir"
	"github.com/tidwall/sjson"
)

//go:embed base
var base embed.FS

func AddSubCommand(app *clir.Cli, parent *clir.Command, w io.Writer) {

	// command
	command := parent.NewSubCommand("template", "Generates a wails template")

	name := ""
	command.StringFlag("name", "The name of the template", &name)

	existingProjectDir := ""
	command.StringFlag("frontend", "A path to an existing frontend project to include in the template", &existingProjectDir)

	// Quiet Init
	quiet := false
	command.BoolFlag("q", "Suppress output to console", &quiet)

	command.Action(func() error {

		// name is mandatory
		if name == "" {
			command.PrintHelp()
			return fmt.Errorf("no template name given")
		}

		// If the current directory is not empty, we create a new directory
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		templateDir := filepath.Join(cwd, name)
		if !fs.DirExists(templateDir) {
			err := os.MkdirAll(templateDir, 0755)
			if err != nil {
				return err
			}
		}
		empty, err := fs.DirIsEmpty(templateDir)
		if err != nil {
			return err
		}
		if !empty {
			templateDir = filepath.Join(cwd, name)
			println("Creating new template directory:", name)
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
			Name         string
			Description  string
			TemplateDir  string
			WailsVersion string
		}

		println("Extracting base template files...")

		err = g.Extract(templateDir, &templateData{
			Name:         name,
			TemplateDir:  templateDir,
			WailsVersion: app.Version(),
		})
		if err != nil {
			return err
		}

		err = os.Chdir(cwd)
		if err != nil {
			return err
		}

		// If we aren't migrating the files, just exit
		if existingProjectDir == "" {
			return nil
		}

		// Remove frontend directory
		frontendDir := filepath.Join(templateDir, "frontend")
		err = os.RemoveAll(frontendDir)
		if err != nil {
			return err
		}

		// Copy the files into a new frontend directory
		println("Migrating existing project files to frontend directory...")

		sourceDir, err := filepath.Abs(existingProjectDir)
		if err != nil {
			return err
		}

		newFrontendDir := filepath.Join(templateDir, "frontend")
		err = fs.CopyDirExtended(sourceDir, newFrontendDir, []string{name, "node_modules"})
		if err != nil {
			return err
		}

		// Process package.json
		err = processPackageJSON(frontendDir)
		if err != nil {
			return err
		}

		// Process package-lock.json
		err = processPackageLockJSON(frontendDir)
		if err != nil {
			return err
		}

		// Remove node_modules - ignore error, eg it doesn't exist
		_ = os.RemoveAll(filepath.Join(frontendDir, "node_modules"))

		return nil

	})
}

func processPackageJSON(frontendDir string) error {
	var err error

	packageJSON := filepath.Join(frontendDir, "package.json")
	if !fs.FileExists(packageJSON) {
		println("No package.json found - cannot process.")
		return nil
	}

	json, err := os.ReadFile(packageJSON)
	if err != nil {
		return err
	}

	// We will ignore these errors - it's not critical
	println("Updating package.json data...")
	json, _ = sjson.SetBytes(json, "name", "{{.ProjectName}}")
	json, _ = sjson.SetBytes(json, "author", "{{.AuthorName}}")

	err = os.WriteFile(packageJSON, json, 0644)
	if err != nil {
		return err
	}
	baseDir := filepath.Dir(packageJSON)
	println("Renaming package.json -> package.tmpl.json...")
	err = os.Rename(packageJSON, filepath.Join(baseDir, "package.tmpl.json"))
	if err != nil {
		return err
	}
	return nil
}

func processPackageLockJSON(frontendDir string) error {
	var err error

	filename := filepath.Join(frontendDir, "package-lock.json")
	if !fs.FileExists(filename) {
		println("No package-lock.json found - cannot process.")
		return nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	json := string(data)

	// We will ignore these errors - it's not critical
	println("Updating package-lock.json data...")
	json, _ = sjson.Set(json, "name", "{{.ProjectName}}")

	err = os.WriteFile(filename, []byte(json), 0644)
	if err != nil {
		return err
	}
	baseDir := filepath.Dir(filename)
	println("Renaming package-lock.json -> package-lock.tmpl.json...")
	err = os.Rename(filename, filepath.Join(baseDir, "package-lock.tmpl.json"))
	if err != nil {
		return err
	}
	return nil
}
