package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/leaanthony/debme"
	"github.com/leaanthony/gosod"
	"github.com/tidwall/sjson"
	"github.com/wailsapp/wails/v2/cmd/wails/flags"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/template"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/project"
	"github.com/wailsapp/wails/v2/internal/tui"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/commands/bindings"
	"github.com/wailsapp/wails/v2/pkg/commands/buildtags"
)

func generateModule(f *flags.GenerateModule) error {
	if f.NoColour {
		tui.SetNoColour()
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

	err = tui.WithSpinner("Generating bindings", func() error {
		_, err := bindings.GenerateBindings(bindings.Options{
			Compiler:     f.Compiler,
			Tags:         buildTags,
			TsPrefix:     projectConfig.Bindings.TsGeneration.Prefix,
			TsSuffix:     projectConfig.Bindings.TsGeneration.Suffix,
			TsOutputType: projectConfig.Bindings.TsGeneration.OutputType,
		})
		return err
	})
	return err
}

func generateTemplate(f *flags.GenerateTemplate) error {
	if f.NoColour {
		tui.SetNoColour()
	}

	quiet := f.Quiet
	logger := clilogger.New(os.Stdout)
	logger.Mute(quiet)
	_ = logger

	if f.Name == "" {
		return fmt.Errorf("please provide a template name using the -name flag")
	}

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

	tui.Section("Generating template")

	if !empty {
		templateDir = filepath.Join(cwd, f.Name)
		tui.BulletPoint("Creating new template directory: %s", f.Name)
		err = fs.Mkdir(templateDir)
		if err != nil {
			return err
		}
	}

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

	tui.BulletPoint("Extracting base template files...")

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

	if f.Frontend == "" {
		fmt.Println()
		fmt.Println()
		tui.Info("No frontend specified to migrate. Template created.")
		fmt.Println()
		return nil
	}

	frontendDir := filepath.Join(templateDir, "frontend")
	err = os.RemoveAll(frontendDir)
	if err != nil {
		return err
	}

	tui.BulletPoint("Migrating existing project files to frontend directory...")

	sourceDir, err := filepath.Abs(f.Frontend)
	if err != nil {
		return err
	}

	newFrontendDir := filepath.Join(templateDir, "frontend")
	err = fs.CopyDirExtended(sourceDir, newFrontendDir, []string{f.Name, "node_modules"})
	if err != nil {
		return err
	}

	err = processPackageJSON(frontendDir)
	if err != nil {
		return err
	}

	err = processPackageLockJSON(frontendDir)
	if err != nil {
		return err
	}

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

	tui.BulletPoint("Updating package.json data...")
	json, _ = sjson.SetBytes(json, "name", "{{.ProjectName}}")
	json, _ = sjson.SetBytes(json, "author", "{{.AuthorName}}")

	err = os.WriteFile(packageJSON, json, 0o644)
	if err != nil {
		return err
	}
	baseDir := filepath.Dir(packageJSON)
	tui.BulletPoint("Renaming package.json -> package.tmpl.json...")
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

	tui.BulletPoint("Updating package-lock.json data...")
	json, _ = sjson.Set(json, "name", "{{.ProjectName}}")

	err = os.WriteFile(filename, []byte(json), 0o644)
	if err != nil {
		return err
	}
	baseDir := filepath.Dir(filename)
	tui.BulletPoint("Renaming package-lock.json -> package-lock.tmpl.json...")
	err = os.Rename(filename, filepath.Join(baseDir, "package-lock.tmpl.json"))
	if err != nil {
		return err
	}
	return nil
}
