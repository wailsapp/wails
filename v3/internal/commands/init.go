package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/templates"
	"github.com/wailsapp/wails/v3/internal/term"
)

var DisableFooter bool

// GitURLToModuleName converts a git URL to a Go module name by removing common prefixes
// and suffixes. It handles HTTPS, SSH, Git protocol, and filesystem URLs.
func GitURLToModuleName(gitURL string) string {
	moduleName := gitURL
	if strings.HasSuffix(moduleName, ".git") {
		moduleName = moduleName[:len(moduleName)-4]
	}
	// Handle various URL schemes
	for _, prefix := range []string{
		"https://",
		"http://",
		"git://",
		"ssh://",
		"file://",
	} {
		if strings.HasPrefix(moduleName, prefix) {
			moduleName = moduleName[len(prefix):]
			break
		}
	}
	// Handle SSH URLs (git@github.com:username/project.git)
	if strings.HasPrefix(moduleName, "git@") {
		// Remove the 'git@' prefix
		moduleName = moduleName[4:]
		// Replace ':' with '/' for proper module path
		moduleName = strings.Replace(moduleName, ":", "/", 1)
	}
	// Remove leading forward slash for file system paths
	moduleName = strings.TrimPrefix(moduleName, "/")
	return moduleName
}

func initGitRepository(projectDir string, gitURL string) error {
	// Initialize repository
	repo, err := git.PlainInit(projectDir, false)
	if err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// Create remote
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{gitURL},
	})
	if err != nil {
		return fmt.Errorf("failed to create git remote: %w", err)
	}

	// Update go.mod with the module name
	moduleName := GitURLToModuleName(gitURL)

	goModPath := filepath.Join(projectDir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	// Replace module name
	lines := strings.Split(string(content), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("go.mod is empty")
	}
	lines[0] = fmt.Sprintf("module %s", moduleName)
	newContent := strings.Join(lines, "\n")

	err = os.WriteFile(goModPath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write go.mod: %w", err)
	}

	// Stage all files
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get git worktree: %w", err)
	}

	_, err = worktree.Add(".")
	if err != nil {
		return fmt.Errorf("failed to stage files: %w", err)
	}

	return nil
}

func Init(options *flags.Init) error {

	if options.List {
		term.Header("Available templates")
		return printTemplates()
	}

	if options.Quiet {
		term.DisableOutput()
	}
	term.Header("Init project")

	// Check if the template is a typescript template
	isTypescript := false
	if strings.HasSuffix(options.TemplateName, "-ts") {
		isTypescript = true
	}

	if options.ProjectName == "" {
		interactiveOptions, err := startInteractive(options) // Call startInteractive, pass options
		if err != nil {
			return err // Return error from interactive form
		}
		options = interactiveOptions // Use options returned from interactive form
		if options.ProjectName == "" {
			return fmt.Errorf("project name is required") // Ensure project name is provided
		}
	}

	options.ProjectName = sanitizeFileName(options.ProjectName)

	err := templates.Install(options)
	if err != nil {
		return err
	}

	// Rename gitignore to .gitignore
	err = os.Rename(filepath.Join(options.ProjectDir, "gitignore"), filepath.Join(options.ProjectDir, ".gitignore"))
	if err != nil {
		return err
	}

	// Generate build assets
	buildAssetsOptions := &BuildAssetsOptions{
		Name:               options.ProjectName,
		Dir:                filepath.Join(options.ProjectDir, "build"),
		Silent:             true,
		ProductCompany:     options.ProductCompany,
		ProductName:        options.ProductName,
		ProductDescription: options.ProductDescription,
		ProductVersion:     options.ProductVersion,
		ProductIdentifier:  options.ProductIdentifier,
		ProductCopyright:   options.ProductCopyright,
		ProductComments:    options.ProductComments,
		Typescript:         isTypescript,
	}
	err = GenerateBuildAssets(buildAssetsOptions)
	if err != nil {
		return err
	}
	// Initialize git repository if URL is provided
	if options.Git != "" {
		err = initGitRepository(options.ProjectDir, options.Git)
		if err != nil {
			return err
		}
		if !options.Quiet {
			term.Infof("Initialized git repository with remote: %s\n", options.Git)
		}
	}
	return nil
}

func printTemplates() error {
	defaultTemplates := templates.GetDefaultTemplates()

	pterm.Println()
	table := pterm.TableData{{"Name", "Description"}}
	for _, template := range defaultTemplates {
		table = append(table, []string{template.Name, template.Description})
	}
	err := pterm.DefaultTable.WithHasHeader(true).WithBoxed(true).WithData(table).Render()
	pterm.Println()
	return err
}

func sanitizeFileName(fileName string) string {
	// Regular expression to match non-allowed characters in file names
	// You can adjust this based on the specific requirements of your file system
	reg := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)

	// Replace matched characters with an underscore or any other safe character
	return reg.ReplaceAllString(fileName, "_")
}

// startInteractive starts the interactive form and returns the populated flags.Init and error.
// It now accepts the flags.Init as input, to initialize default values if needed.
func startInteractive(initialOptions *flags.Init) (*flags.Init, error) {
	var templateOptions []huh.Option[string]
	var templateSelect *huh.Select[string] // Declare templateSelect outside
	confirmProjectCreation := false        // Local variable for confirmation

	options := initialOptions // Use the passed-in options, avoids shadowing

	if options == nil { // Defensive check in case nil options are passed
		options = &flags.Init{
			ProductVersion: "1.0.0", // Default value if no initial options are provided
		}
	} else if options.ProductVersion == "" {
		options.ProductVersion = "1.0.0" // Ensure default if not set in initial options
	}
	templateName := &options.TemplateName // keep pointer for default value setting

	templateSelect = huh.NewSelect[string](). // Initialize templateSelect here
							Title("Template").
							Description("Project template to use (Enter to list)").
							Options(templateOptions...).
							Value(templateName) // Bind to options.TemplateName

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
				Description("Name of project").
				Value(&options.ProjectName), // Bind to options.ProjectName

			huh.NewInput().
				Title("Project Dir").
				Description("Target directory (empty for default)").
				Value(&options.ProjectDir), // Bind to options.ProjectDir

			templateSelect, // Use templateSelect here, no assignment
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Git Repo URL (optional)").
				Description("Git repo to initialize (optional)").
				Value(&options.Git), // Bind to options.Git
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Company (optional)").
				Value(&options.ProductCompany), // Bind to options.ProductCompany
			huh.NewInput().
				Title("Product Name (optional)").
				Value(&options.ProductName), // Bind to options.ProductName
			huh.NewInput().
				Title("Version (optional)").
				Value(&options.ProductVersion), // Bind to options.ProductVersion
			huh.NewInput().
				Title("ID (optional)").
				Value(&options.ProductIdentifier), // Bind to options.ProductIdentifier
			huh.NewInput().
				Title("Copyright (optional)").
				Value(&options.ProductCopyright), // Bind to options.ProductCopyright
			huh.NewText().
				Title("Description (optional)").
				Lines(1).
				Value(&options.ProductDescription), // Bind to options.ProductDescription
			huh.NewText().
				Title("Comments (optional)").
				Lines(1).
				Value(&options.ProductComments), // Bind to options.ProductComments
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Confirm?").
				Value(&confirmProjectCreation), // Bind to local variable confirmProjectCreation
		),
	)

	// Dynamically load and order templates when the "Template" select is rendered
	templateSelect.OptionsFunc(func() []huh.Option[string] {
		defaultTemplates := templates.GetDefaultTemplates()

		// Reorder templates: Ensure "vanilla" is first
		vanillaTemplate := templates.TemplateData{Name: "vanilla", Description: ""} // Use TemplateData
		orderedTemplates := []templates.TemplateData{}                              // Use TemplateData
		vanillaFound := false
		for _, t := range defaultTemplates {
			if t.Name == vanillaTemplate.Name {
				orderedTemplates = append([]templates.TemplateData{t}, orderedTemplates...) // Prepend vanilla - use TemplateData
				vanillaFound = true
			} else {
				orderedTemplates = append(orderedTemplates, t) // Append other templates - use TemplateData
			}
		}
		if !vanillaFound { // If "vanilla" template isn't found (unlikely, but for safety)
			orderedTemplates = append([]templates.TemplateData{vanillaTemplate}, orderedTemplates...) // Use TemplateData
		}

		templateOptions = nil // Clear existing options
		for _, t := range orderedTemplates {
			templateOptions = append(templateOptions, huh.NewOption(t.Name, t.Name))
		}

		// Set default template to "vanilla" - do this BEFORE running the form
		*templateName = "vanilla" // Set default value here, using pointer

		return templateOptions
	}, nil)

	formStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")). // A nice cyan color
		Padding(1, 2).
		Margin(1, 0) // Add a little margin above and below

	formRenderer := lipgloss.NewRenderer(os.Stdout) // Corrected: Pass os.Stdout to NewRenderer

	err := form.Run()
	if err != nil {
		term.Error(err)
		return nil, err // Return error if form fails
	}

	formString := formRenderer.NewStyle().Render(form.View()) // Get the form's rendered output
	styledForm := formStyle.Render(formString)                // Apply lipgloss style
	fmt.Println(styledForm)                                   // Print the styled form

	if confirmProjectCreation { // Check local variable confirmProjectCreation
		fmt.Println("Creating project...")
		return options, nil // Return populated options
	} else {
		fmt.Println("Project creation cancelled.")
		return nil, fmt.Errorf("project creation cancelled by user") // Return cancellation error
	}
}
