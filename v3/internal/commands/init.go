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
		interactiveOptions, err := startInteractive(options)
		if err != nil {
			return err
		}
		options = interactiveOptions
		if options.ProjectName == "" {
			return fmt.Errorf("project name is required")
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

func startInteractive(initialOptions *flags.Init) (*flags.Init, error) {
	var templateOptions []huh.Option[string]
	var templateSelect *huh.Select[string]
	confirmProjectCreation := false

	options := initialOptions

	if options == nil {
		options = &flags.Init{
			ProductVersion: "1.0.0",
		}
	} else if options.ProductVersion == "" {
		options.ProductVersion = "1.0.0"
	}
	templateName := &options.TemplateName

	templateSelect = huh.NewSelect[string]().
		Title("Template").
		Description("Project template to use (Enter to list)").
		Options(templateOptions...).
		Value(templateName)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
				Description("Name of project").
				Value(&options.ProjectName),

			huh.NewInput().
				Title("Project Dir").
				Description("Target directory (empty for default)").
				Value(&options.ProjectDir),

			templateSelect,
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Git Repo URL (optional)").
				Description("Git repo to initialize (optional)").
				Value(&options.Git),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Company (optional)").
				Value(&options.ProductCompany),
			huh.NewInput().
				Title("Product Name (optional)").
				Value(&options.ProductName),
			huh.NewInput().
				Title("Version (optional)").
				Value(&options.ProductVersion),
			huh.NewInput().
				Title("ID (optional)").
				Value(&options.ProductIdentifier),
			huh.NewInput().
				Title("Copyright (optional)").
				Value(&options.ProductCopyright),
			huh.NewText().
				Title("Description (optional)").
				Lines(1).
				Value(&options.ProductDescription),
			huh.NewText().
				Title("Comments (optional)").
				Lines(1).
				Value(&options.ProductComments),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Confirm?").
				Value(&confirmProjectCreation),
		),
	)

	templateSelect.OptionsFunc(func() []huh.Option[string] {
		defaultTemplates := templates.GetDefaultTemplates()

		vanillaTemplate := templates.TemplateData{Name: "vanilla", Description: ""}
		orderedTemplates := []templates.TemplateData{}
		vanillaFound := false
		for _, t := range defaultTemplates {
			if t.Name == vanillaTemplate.Name {
				orderedTemplates = append([]templates.TemplateData{t}, orderedTemplates...)
				vanillaFound = true
			} else {
				orderedTemplates = append(orderedTemplates, t)
			}
		}
		if !vanillaFound {
			orderedTemplates = append([]templates.TemplateData{vanillaTemplate}, orderedTemplates...)
		}

		templateOptions = nil
		for _, t := range orderedTemplates {
			templateOptions = append(templateOptions, huh.NewOption(t.Name, t.Name))
		}

		*templateName = "vanilla"

		return templateOptions
	}, nil)

	formStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Margin(1, 0)

	formRenderer := lipgloss.NewRenderer(os.Stdout)

	err := form.Run()
	if err != nil {
		term.Error(fmt.Errorf("Interactive form failed: %w", err))
		return nil, fmt.Errorf("form execution failed: %w", err)
	}

	formString := formRenderer.NewStyle().Render(form.View())
	styledForm := formStyle.Render(formString)
	fmt.Println(styledForm)

	if confirmProjectCreation {
		fmt.Println("Creating project...")
		return options, nil
	} else {
		fmt.Println("Project creation cancelled.")
		return nil, fmt.Errorf("project creation cancelled by user")
	}
}
