package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/leaanthony/debme"
	"github.com/leaanthony/gosod"
	"github.com/pterm/pterm"
	"github.com/tidwall/sjson"
	"github.com/wailsapp/wails/v2/cmd/wails/flags"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/template"
	"github.com/wailsapp/wails/v2/internal/colour"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/project"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/commands/bindings"
	"github.com/wailsapp/wails/v2/pkg/commands/buildtags"
)

func generateModule(f *flags.GenerateModule) error {
	if f.NoColour {
		pterm.DisableColor()
		colour.ColourEnabled = false
	}

	quiet := f.Verbosity == flags.Quiet
	logger := clilogger.New(os.Stdout)
	logger.Mute(quiet)

	buildTags, err := buildtags.Parse(f.Tags)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	projectConfig, err := project.Load(cwd)
	if err != nil {
		return err
	}

	if projectConfig.Bindings.TsGeneration.OutputType == "" {
		projectConfig.Bindings.TsGeneration.OutputType = "classes"
	}

	_, err = bindings.GenerateBindings(bindings.Options{
		Compiler:     f.Compiler,
		Tags:         buildTags,
		TsPrefix:     projectConfig.Bindings.TsGeneration.Prefix,
		TsSuffix:     projectConfig.Bindings.TsGeneration.Suffix,
		TsOutputType: projectConfig.Bindings.TsGeneration.OutputType,
	})
	if err != nil {
		return err
	}
	return nil
}

func generateTemplate(f *flags.GenerateTemplate) error {
	if f.NoColour {
		pterm.DisableColor()
		colour.ColourEnabled = false
	}

	quiet := f.Quiet
	logger := clilogger.New(os.Stdout)
	logger.Mute(quiet)

	// name is mandatory
	if f.Name == "" {
		return fmt.Errorf("please provide a template name using the -name flag")
	}

	// If the current directory is not empty, we create a new directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	templateDir := filepath.Join(cwd, f.Name)
	if !fs.DirExists(templateDir) {
		err := os.MkdirAll(templateDir, 0o755)
		if err != nil {
			return err
		}
	}
	empty, err := fs.DirIsEmpty(templateDir)
	if err != nil {
		return err
	}

	pterm.DefaultSection.Println("Generating template")

	if !empty {
		templateDir = filepath.Join(cwd, f.Name)
		printBulletPoint("Creating new template directory:", f.Name)
		err = fs.Mkdir(templateDir)
		if err != nil {
			return err
		}
	}

	// Create base template
	baseTemplate, err := debme.FS(template.Base, "base")
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

	printBulletPoint("Extracting base template files...")

	err = g.Extract(templateDir, &templateData{
		Name:         f.Name,
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
	if f.Frontend == "" {
		pterm.Println()
		pterm.Println()
		pterm.Info.Println("No frontend specified to migrate. Template created.")
		pterm.Println()
		return nil
	}

	// Remove frontend directory
	frontendDir := filepath.Join(templateDir, "frontend")
	err = os.RemoveAll(frontendDir)
	if err != nil {
		return err
	}

	// Copy the files into a new frontend directory
	printBulletPoint("Migrating existing project files to frontend directory...")

	sourceDir, err := filepath.Abs(f.Frontend)
	if err != nil {
		return err
	}

	newFrontendDir := filepath.Join(templateDir, "frontend")
	err = fs.CopyDirExtended(sourceDir, newFrontendDir, []string{f.Name, "node_modules"})
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
}

func processPackageJSON(frontendDir string) error {
	var err error

	packageJSON := filepath.Join(frontendDir, "package.json")
	if !fs.FileExists(packageJSON) {
		return fmt.Errorf("no package.json found - cannot process")
	}

	json, err := os.ReadFile(packageJSON)
	if err != nil {
		return err
	}

	// We will ignore these errors - it's not critical
	printBulletPoint("Updating package.json data...")
	json, _ = sjson.SetBytes(json, "name", "{{.ProjectName}}")
	json, _ = sjson.SetBytes(json, "author", "{{.AuthorName}}")

	err = os.WriteFile(packageJSON, json, 0o644)
	if err != nil {
		return err
	}
	baseDir := filepath.Dir(packageJSON)
	printBulletPoint("Renaming package.json -> package.tmpl.json...")
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
		return fmt.Errorf("no package-lock.json found - cannot process")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	json := string(data)

	// We will ignore these errors - it's not critical
	printBulletPoint("Updating package-lock.json data...")
	json, _ = sjson.Set(json, "name", "{{.ProjectName}}")

	err = os.WriteFile(filename, []byte(json), 0o644)
	if err != nil {
		return err
	}
	baseDir := filepath.Dir(filename)
	printBulletPoint("Renaming package-lock.json -> package-lock.tmpl.json...")
	err = os.Rename(filename, filepath.Join(baseDir, "package-lock.tmpl.json"))
	if err != nil {
		return err
	}
	return nil
}
