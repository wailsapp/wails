package initialise

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/flytam/filenamify"
	"github.com/leaanthony/slicer"

	"github.com/wailsapp/wails/v2/pkg/buildassets"

	"github.com/wailsapp/wails/v2/pkg/templates"

	"github.com/leaanthony/clir"
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/git"
)

// AddSubcommand adds the `init` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer) error {

	command := app.NewSubCommand("init", "Initialise a new Wails project")

	// Setup template name flag
	templateName := "vanilla"
	description := "Name of built-in template to use, path to template or template url."
	command.StringFlag("t", description, &templateName)

	// Setup project name
	projectName := ""
	command.StringFlag("n", "Name of project", &projectName)

	// For CI
	ciMode := false
	command.BoolFlag("ci", "CI Mode", &ciMode).Hidden()

	// Setup project directory
	projectDirectory := ""
	command.StringFlag("d", "Project directory", &projectDirectory)

	// Quiet Init
	quiet := false
	command.BoolFlag("q", "Suppress output to console", &quiet)

	initGit := false
	gitInstalled := git.IsInstalled()
	if gitInstalled {
		// Git Init
		command.BoolFlag("g", "Initialise git repository", &initGit)
	}

	// VSCode project files
	ide := ""
	command.StringFlag("ide", "Generate IDE project files", &ide)

	// List templates
	list := false
	command.BoolFlag("l", "List templates", &list)

	command.Action(func() error {

		// Create logger
		logger := clilogger.New(w)
		logger.Mute(quiet)

		// Are we listing templates?
		if list {
			app.PrintBanner()
			err := templates.OutputList(logger)
			logger.Println("")
			return err
		}

		// Validate name
		if len(projectName) == 0 {
			logger.Println("ERROR: Project name required")
			logger.Println("")
			command.PrintHelp()
			return nil
		}

		// Validate IDE option
		supportedIDEs := slicer.String([]string{"vscode", "goland"})
		ide = strings.ToLower(ide)
		if ide != "" {
			if !supportedIDEs.Contains(ide) {
				return fmt.Errorf("ide '%s' not supported. Valid values: %s", ide, supportedIDEs.Join(" "))
			}
		}

		if !quiet {
			app.PrintBanner()
		}

		task := fmt.Sprintf("Initialising Project '%s'", projectName)
		logger.Println(task)
		logger.Println(strings.Repeat("-", len(task)))

		projectFilename, err := filenamify.Filenamify(projectName, filenamify.Options{
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
			ProjectName:         projectName,
			TargetDir:           projectDirectory,
			TemplateName:        templateName,
			Logger:              logger,
			IDE:                 ide,
			InitGit:             initGit,
			ProjectNameFilename: projectFilename,
			WailsVersion:        app.Version(),
			GoSDKPath:           goSDKPath,
		}

		// Try to discover author details from git config
		findAuthorDetails(options)

		return initProject(options, quiet, ciMode)
	})

	return nil
}

// initProject is our main init command
func initProject(options *templates.Options, quiet bool, ciMode bool) error {

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

	if !ciMode {
		// Run `go mod tidy` to ensure `go.sum` is up to date
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = options.TargetDir
		cmd.Stderr = os.Stderr
		if !quiet {
			println("")
			cmd.Stdout = os.Stdout
		}
		err = cmd.Run()
		if err != nil {
			return err
		}
	} else {
		// Update go mod
		workspace := os.Getenv("GITHUB_WORKSPACE")
		println("GitHub workspace:", workspace)
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
	options.Logger.Println("Project Name:      " + options.ProjectName)
	options.Logger.Println("Project Directory: " + options.TargetDir)
	options.Logger.Println("Project Template:  " + options.TemplateName)
	options.Logger.Println("Template Support:  " + template.HelpURL)

	// IDE message
	switch options.IDE {
	case "vscode":
		options.Logger.Println("VSCode config files generated.")
	case "goland":
		options.Logger.Println("Goland config files generated.")
	}

	if options.InitGit {
		options.Logger.Println("Git repository initialised.")
	}

	if remote {
		options.Logger.Println("\nNOTE: You have created a project using a remote template. The Wails project takes no responsibility for 3rd party templates. Only use remote templates that you trust.")
	}

	options.Logger.Println("")
	options.Logger.Println(fmt.Sprintf("Initialised project '%s' in %s.", options.ProjectName, elapsed.Round(time.Millisecond).String()))
	options.Logger.Println("")

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
	err = os.WriteFile(filepath.Join(options.TargetDir, ".gitignore"), []byte(strings.Join(ignore, "\n")), 0644)
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
		log.Fatal(err)
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for i, line := range lines {
		println(line)
		if strings.HasPrefix(line, "// replace") {
			println("Found replace line")
			splitLine := strings.Split(line, " ")
			splitLine[5] = targetPath + "/v2"
			lines[i] = strings.Join(splitLine[1:], " ")
			continue
		}
	}

	err = os.WriteFile("go.mod", []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
