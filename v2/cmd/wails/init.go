package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/flytam/filenamify"
	"github.com/leaanthony/slicer"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v2/cmd/wails/flags"
	"github.com/wailsapp/wails/v2/internal/colour"
	"github.com/wailsapp/wails/v2/pkg/buildassets"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/git"
	"github.com/wailsapp/wails/v2/pkg/templates"
)

func initProject(f *flags.Init) error {
	if f.NoColour {
		pterm.DisableColor()
		colour.ColourEnabled = false
	}

	quiet := f.Quiet

	// Create logger
	logger := clilogger.New(os.Stdout)
	logger.Mute(quiet)

	// Are we listing templates?
	if f.List {
		app.PrintBanner()
		templateList, err := templates.List()
		if err != nil {
			return err
		}

		pterm.DefaultSection.Println("Available templates")

		table := pterm.TableData{{"Template", "Short Name", "Description"}}
		for _, template := range templateList {
			table = append(table, []string{template.Name, template.ShortName, template.Description})
		}
		err = pterm.DefaultTable.WithHasHeader(true).WithBoxed(true).WithData(table).Render()
		pterm.Println()
		return err
	}

	// Validate name
	if len(f.ProjectName) == 0 {
		return fmt.Errorf("please provide a project name using the -n flag")
	}

	// Validate IDE option
	supportedIDEs := slicer.String([]string{"vscode", "goland"})
	ide := strings.ToLower(f.IDE)
	if ide != "" {
		if !supportedIDEs.Contains(ide) {
			return fmt.Errorf("ide '%s' not supported. Valid values: %s", ide, supportedIDEs.Join(" "))
		}
	}

	if !quiet {
		app.PrintBanner()
	}

	pterm.DefaultSection.Printf("Initialising Project '%s'", f.ProjectName)

	projectFilename, err := filenamify.Filenamify(f.ProjectName, filenamify.Options{
		Replacement: "_",
		MaxLength:   255,
	})
	if err != nil {
		return err
	}
	goBinary, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("unable to find Go compiler. Please download and install Go: https://golang.org/dl/")
	}

	// Get base path and convert to forward slashes
	goPath := filepath.ToSlash(filepath.Dir(goBinary))
	// Trim bin directory
	goSDKPath := strings.TrimSuffix(goPath, "/bin")

	// Create Template Options
	options := &templates.Options{
		ProjectName:         f.ProjectName,
		TargetDir:           f.ProjectDir,
		TemplateName:        f.TemplateName,
		Logger:              logger,
		IDE:                 ide,
		InitGit:             f.InitGit,
		ProjectNameFilename: projectFilename,
		WailsVersion:        app.Version(),
		GoSDKPath:           goSDKPath,
	}

	// Try to discover author details from git config
	findAuthorDetails(options)

	// Start Time
	start := time.Now()

	// Install the template
	remote, template, err := templates.Install(options)
	if err != nil {
		return err
	}

	// Install the default assets
	err = buildassets.Install(options.TargetDir)
	if err != nil {
		return err
	}

	err = os.Chdir(options.TargetDir)
	if err != nil {
		return err
	}

	// Change the module name to project name
	err = updateModuleNameToProjectName(options, quiet)
	if err != nil {
		return err
	}

	if !f.CIMode {
		// Run `go mod tidy` to ensure `go.sum` is up to date
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = options.TargetDir
		cmd.Stderr = os.Stderr
		if !quiet {
			cmd.Stdout = os.Stdout
		}
		err = cmd.Run()
		if err != nil {
			return err
		}
	} else {
		// Update go mod
		workspace := os.Getenv("GITHUB_WORKSPACE")
		pterm.Println("GitHub workspace:", workspace)
		if workspace == "" {
			os.Exit(1)
		}
		updateReplaceLine(workspace)
	}

	// Remove the `.git`` directory in the template project
	err = os.RemoveAll(".git")
	if err != nil {
		return err
	}

	if options.InitGit {
		err = initGit(options)
		if err != nil {
			return err
		}
	}

	if quiet {
		return nil
	}

	// Output stats
	elapsed := time.Since(start)

	// Create pterm table
	table := pterm.TableData{
		{"Project Name", options.ProjectName},
		{"Project Directory", options.TargetDir},
		{"Template", template.Name},
		{"Template Source", template.HelpURL},
	}
	err = pterm.DefaultTable.WithData(table).Render()
	if err != nil {
		return err
	}

	// IDE message
	switch options.IDE {
	case "vscode":
		pterm.Println()
		pterm.Info.Println("VSCode config files generated.")
	case "goland":
		pterm.Println()
		pterm.Info.Println("Goland config files generated.")
	}

	if options.InitGit {
		pterm.Info.Println("Git repository initialised.")
	}

	if remote {
		pterm.Warning.Println("NOTE: You have created a project using a remote template. The Wails project takes no responsibility for 3rd party templates. Only use remote templates that you trust.")
	}

	pterm.Println("")
	pterm.Printf("Initialised project '%s' in %s.\n", options.ProjectName, elapsed.Round(time.Millisecond).String())
	pterm.Println("")

	return nil
}

func initGit(options *templates.Options) error {
	err := git.InitRepo(options.TargetDir)
	if err != nil {
		return errors.Wrap(err, "Unable to initialise git repository:")
	}

	ignore := []string{
		"build/bin",
		"frontend/dist",
		"frontend/node_modules",
	}
	err = os.WriteFile(filepath.Join(options.TargetDir, ".gitignore"), []byte(strings.Join(ignore, "\n")), 0o644)
	if err != nil {
		return errors.Wrap(err, "Unable to create gitignore")
	}

	return nil
}

// findAuthorDetails tries to find the user's name and email
// from gitconfig. If it finds them, it stores them in the project options
func findAuthorDetails(options *templates.Options) {
	if git.IsInstalled() {
		name, err := git.Name()
		if err == nil {
			options.AuthorName = strings.TrimSpace(name)
		}

		email, err := git.Email()
		if err == nil {
			options.AuthorEmail = strings.TrimSpace(email)
		}
	}
}

func updateReplaceLine(targetPath string) {
	file, err := os.Open("go.mod")
	if err != nil {
		fatal(err.Error())
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	err = file.Close()
	if err != nil {
		fatal(err.Error())
	}

	if err := scanner.Err(); err != nil {
		fatal(err.Error())
	}

	for i, line := range lines {
		println(line)
		if strings.HasPrefix(line, "// replace") {
			pterm.Println("Found replace line")
			splitLine := strings.Split(line, " ")
			splitLine[5] = targetPath + "/v2"
			lines[i] = strings.Join(splitLine[1:], " ")
			continue
		}
	}

	err = os.WriteFile("go.mod", []byte(strings.Join(lines, "\n")), 0o644)
	if err != nil {
		fatal(err.Error())
	}
}

func updateModuleNameToProjectName(options *templates.Options, quiet bool) error {
	cmd := exec.Command("go", "mod", "edit", "-module", options.ProjectName)
	cmd.Dir = options.TargetDir
	cmd.Stderr = os.Stderr
	if !quiet {
		cmd.Stdout = os.Stdout
	}

	return cmd.Run()
}
