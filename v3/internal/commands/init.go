package commands

import (
	"fmt"
	"github.com/go-git/go-git/v5/config"
	"github.com/wailsapp/wails/v3/internal/term"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/templates"
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
		return fmt.Errorf("please use the -n flag to specify a project name")
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
