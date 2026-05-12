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
	"github.com/wailsapp/wails/v2/cmd/wails/flags"
	"github.com/wailsapp/wails/v2/internal/tui"
	"github.com/wailsapp/wails/v2/pkg/buildassets"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/git"
	"github.com/wailsapp/wails/v2/pkg/templates"
)

func initProject(f *flags.Init) error {
	if f.NoColour {
		tui.SetNoColour()
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

		tui.Section("Available templates")

		table := [][]string{{"Template", "Short Name", "Description"}}
		for _, template := range templateList {
			table = append(table, []string{template.Name, template.ShortName, template.Description})
		}
		tui.HeaderTable(table)
		fmt.Println()
		return nil
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

	tui.Section(fmt.Sprintf("Initialising Project '%s'", f.ProjectName))

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

	goPath := filepath.ToSlash(filepath.Dir(goBinary))
	goSDKPath := strings.TrimSuffix(goPath, "/bin")

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

	findAuthorDetails(options)

	start := time.Now()

	var remote bool
	var template *templates.Template
	err = tui.WithSpinner("Installing template", func() error {
		var installErr error
		remote, template, installErr = templates.Install(options)
		return installErr
	})
	if err != nil {
		return err
	}

	err = buildassets.Install(options.TargetDir)
	if err != nil {
		return err
	}

	err = os.Chdir(options.TargetDir)
	if err != nil {
		return err
	}

	err = updateModuleNameToProjectName(options, quiet)
	if err != nil {
		return err
	}

	if !f.CIMode {
		err = tui.WithSpinner("Running go mod tidy", func() error {
			cmd := exec.Command("go", "mod", "tidy")
			cmd.Dir = options.TargetDir
			cmd.Stderr = os.Stderr
			if !quiet {
				cmd.Stdout = os.Stdout
			}
			return cmd.Run()
		})
		if err != nil {
			return err
		}
	} else {
		workspace := os.Getenv("GITHUB_WORKSPACE")
		fmt.Println("GitHub workspace:", workspace)
		if workspace == "" {
			os.Exit(1)
		}
		updateReplaceLine(workspace)
	}

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

	elapsed := time.Since(start)

	table := [][]string{
		{"Project Name", options.ProjectName},
		{"Project Directory", options.TargetDir},
		{"Template", template.Name},
		{"Template Source", template.HelpURL},
	}
	tui.Table(table)

	switch options.IDE {
	case "vscode":
		fmt.Println()
		tui.Info("VSCode config files generated.")
	case "goland":
		fmt.Println()
		tui.Info("Goland config files generated.")
	}

	if options.InitGit {
		tui.Info("Git repository initialised.")
	}

	if remote {
		tui.Warning("NOTE: You have created a project using a remote template. The Wails project takes no responsibility for 3rd party templates. Only use remote templates that you trust.")
	}

	fmt.Println()
	fmt.Printf("Initialised project '%s' in %s.\n", options.ProjectName, elapsed.Round(time.Millisecond).String())
	fmt.Println()

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
			fmt.Println("Found replace line")
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
