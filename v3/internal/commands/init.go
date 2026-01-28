package commands

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5/config"
	"github.com/wailsapp/wails/v3/internal/defaults"
	"github.com/wailsapp/wails/v3/internal/term"

	"github.com/go-git/go-git/v5"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/templates"
)

var DisableFooter bool

// See https://github.com/git/git/blob/master/Documentation/urls.adoc
var (
	gitProtocolFormat = regexp.MustCompile(`^(?:ssh|git|https?|ftps?|file)://`)
	gitScpLikeGuard   = regexp.MustCompile(`^[^/:]+:`)
	gitScpLikeFormat  = regexp.MustCompile(`^(?:([^@/:]+)@)?([^@/:]+):([^\\].*)$`)
)

// gitURLToModulePath converts a git URL to a Go module name by removing common prefixes
// and suffixes. It handles HTTPS, SSH, Git protocol, and filesystem URLs.
func gitURLToModulePath(gitURL string) string {
	var path string

	if gitProtocolFormat.MatchString(gitURL) {
		// Standard URL
		parsed, err := url.Parse(gitURL)
		if err != nil {
			term.Warningf("invalid Git repository URL: %s; module path will default to 'changeme'", err)
			return "changeme"
		}

		path = parsed.Host + parsed.Path
	} else if gitScpLikeGuard.MatchString(gitURL) {
		// SCP-like URL
		match := gitScpLikeFormat.FindStringSubmatch(gitURL)
		if match != nil {
			sep := ""
			if !strings.HasPrefix(match[3], "/") {
				// Add slash between host and path if missing
				sep = "/"
			}

			path = match[2] + sep + match[3]
		}
	}

	if path == "" {
		// Filesystem path
		path = gitURL
	}

	if strings.HasSuffix(path, ".git") {
		path = path[:len(path)-4]
	}

	// Remove leading forward slash for file system paths
	return strings.TrimPrefix(path, "/")
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

func applyGlobalDefaults(options *flags.Init, globalDefaults defaults.GlobalDefaults) {
	templateName := globalDefaults.GetTemplateName()
	if options.TemplateName == "vanilla" && templateName != "" && templateName != "vanilla" {
		options.TemplateName = templateName
		options.TemplateFromDefaults = true
	}

	if options.ProductCompany == "My Company" && globalDefaults.Author.Company != "" {
		options.ProductCompany = globalDefaults.Author.Company
	}

	if options.ProductCopyright == "\u00a9 now, My Company" {
		options.ProductCopyright = globalDefaults.GenerateCopyright()
	}

	if options.ProductIdentifier == "" && globalDefaults.Project.ProductIdentifierPrefix != "" {
		options.ProductIdentifier = globalDefaults.GenerateProductIdentifier(options.ProjectName)
	}

	if options.ProductDescription == "My Product Description" && globalDefaults.Project.DescriptionTemplate != "" {
		options.ProductDescription = globalDefaults.GenerateDescription(options.ProjectName)
	}

	if options.ProductVersion == "0.1.0" && globalDefaults.Project.DefaultVersion != "" {
		options.ProductVersion = globalDefaults.GetDefaultVersion()
	}

	// Only apply UseInterfaces from defaults if not explicitly set via CLI flag
	// (default value is false, so we check if it's still false and defaults say true)
	if !options.UseInterfaces && globalDefaults.Project.UseInterfaces {
		options.UseInterfaces = globalDefaults.Project.UseInterfaces
		options.UseInterfacesFromDefaults = true
	}
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
		return errors.New("please use the -n flag to specify a project name")
	}

	options.ProjectName = sanitizeFileName(options.ProjectName)

	// Load and apply global defaults
	globalDefaults, err := defaults.Load()
	if err != nil {
		// Log warning but continue - global defaults are optional
		term.Warningf("Could not load global defaults: %v\n", err)
	} else {
		applyGlobalDefaults(options, globalDefaults)
	}

	if options.ModulePath == "" {
		if options.Git == "" {
			options.ModulePath = "changeme"
		} else {
			options.ModulePath = gitURLToModulePath(options.Git)
		}
	}

	err = templates.Install(options)
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
		UseInterfaces:      options.UseInterfaces,
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
