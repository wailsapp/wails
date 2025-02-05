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

var (
	DisableFooter    bool
	projectName      string
	projectDir       string // Shortened variable name
	templateName     string // Shortened variable name
	gitRepo          string
	confirmCreation  bool
	productCompany   string
	productName      string
	productDesc      string // Shortened variable name
	productVersion   string = "1.0.0"
	productID        string // Shortened variable name
	productCopyright string
	productComments  string
)

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

	options.ProjectName = projectName
	options.ProjectDir = projectDir
	options.TemplateName = templateName
	options.Git = gitRepo
	options.ProductCompany = productCompany
	options.ProductName = productName
	options.ProductDescription = productDesc // Use shortened variable
	options.ProductVersion = productVersion
	options.ProductIdentifier = productID // Use shortened variable
	options.ProductCopyright = productCopyright
	options.ProductComments = productComments

	if options.ProjectName == "" {
		startInteractive()
		return nil
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

func startInteractive() {
	var templateOptions []huh.Option[string]
	var templateSelect *huh.Select[string] // Declare templateSelect outside

	templateSelect = huh.NewSelect[string](). // Initialize templateSelect here
							Title("Template").
							Description("Project template to use (Enter to list)").
							Options(templateOptions...).
							Value(&templateName)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
				Description("Name of project").
				Value(&projectName),

			huh.NewInput().
				Title("Project Dir").
				Description("Target directory (empty for default)").
				Value(&projectDir),

			templateSelect, // Use templateSelect here, no assignment
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Git Repo URL (optional)").
				Description("Git repo to initialize (optional)").
				Value(&gitRepo),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Company (optional)").
				Value(&productCompany),
			huh.NewInput().
				Title("Product Name (optional)").
				Value(&productName),
			huh.NewInput().
				Title("Version (optional)").
				Value(&productVersion),
			huh.NewInput().
				Title("ID (optional)").
				Value(&productID),
			huh.NewInput().
				Title("Copyright (optional)").
				Value(&productCopyright),
			huh.NewText().
				Title("Description (optional)").
				Lines(2).
				Value(&productDesc),
			huh.NewText().
				Title("Comments (optional)").
				Lines(2).
				Value(&productComments),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Confirm?").
				Value(&confirmCreation),
		),
	)

	// Dynamically load and order templates when the "Template" select is rendered
	templateSelect.OptionsFunc(func() []huh.Option[string] {
		defaultTemplates := templates.GetDefaultTemplates()

		// Reorder templates: Ensure "vanilla" is first
		vanillaTemplate := templates.TemplateData{Name: "vanilla", Description: ""} // Create vanilla template for comparison
		orderedTemplates := []templates.TemplateData{}
		vanillaFound := false
		for _, t := range defaultTemplates {
			if t.Name == vanillaTemplate.Name {
				orderedTemplates = append([]templates.TemplateData{t}, orderedTemplates...) // Prepend vanilla
				vanillaFound = true
			} else {
				orderedTemplates = append(orderedTemplates, t) // Append other templates
			}
		}
		if !vanillaFound { // If "vanilla" template isn't found (unlikely, but for safety)
			orderedTemplates = append([]templates.TemplateData{vanillaTemplate}, orderedTemplates...)
		}

		templateOptions = nil // Clear existing options
		for _, t := range orderedTemplates {
			templateOptions = append(templateOptions, huh.NewOption(t.Name, t.Name))
		}

		// Set default template to "vanilla" - do this BEFORE running the form
		templateName = "vanilla" // Set default value here

		return templateOptions
	}, nil)

	formStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")). // A nice cyan color
		Padding(10, 2).
		Margin(10, 0) // Add a little margin above and below

	formRenderer := lipgloss.NewRenderer(os.Stdout) // Corrected: Pass os.Stdout to NewRenderer

	err := form.Run()
	if err != nil {
		term.Error(err)
	}

	formString := form.View()                                                   // Get the form's rendered output directly from form.View()
	styledForm := formRenderer.NewStyle().Inherit(formStyle).Render(formString) // Corrected: Use formRenderer.NewStyle().Render()

	fmt.Println(styledForm) // Print the styled form

	if confirmCreation {
		fmt.Println("Creating project...")
		// In real code, continue with project creation logic here
		fmt.Printf("Project Name: %s\n", projectName)
		fmt.Printf("Project Directory: %s\n", projectDir)
		fmt.Printf("Template: %s\n", templateName)
		fmt.Printf("Git Repo: %s\n", gitRepo)
		fmt.Printf("Product Company: %s\n", productCompany)
		fmt.Printf("Product Name: %s\n", productName)
		fmt.Printf("Product Description: %s\n", productDesc)
		fmt.Printf("Product Version: %s\n", productVersion)
		fmt.Printf("Product Identifier: %s\n", productID)
		fmt.Printf("Product Copyright: %s\n", productCopyright)
		fmt.Printf("Product Comments: %s\n", productComments)

		options := &flags.Init{
			ProjectName:        projectName,
			ProjectDir:         projectDir,
			TemplateName:       templateName,
			Git:                gitRepo,
			ProductCompany:     productCompany,
			ProductName:        productName,
			ProductDescription: productDesc,
			ProductVersion:     productVersion,
			ProductIdentifier:  productID,
			ProductCopyright:   productCopyright,
			ProductComments:    productComments,
		}

		// Call the original Init function with the populated options
		if err := Init(options); err != nil {
			term.Error(err)
		}

	} else {
		fmt.Println("Project creation cancelled.")
	}
}
