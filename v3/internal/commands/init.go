package commands

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/defaults"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/git"
	"github.com/wailsapp/wails/v3/internal/setupwizard"
	"github.com/wailsapp/wails/v3/internal/templates"
	"github.com/wailsapp/wails/v3/internal/term"
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
	if err := git.Init(projectDir); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}
	if err := git.RemoteAdd(projectDir, "origin", gitURL); err != nil {
		return fmt.Errorf("failed to create git remote: %w", err)
	}
	if err := git.AddAll(projectDir); err != nil {
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

	// Apply UseInterfaces from global defaults only when the user hasn't explicitly
	// disabled it via --useinterfaces=false. The CLI default is true, so if it's
	// still true here the user didn't override it and we can apply the configured default.
	if options.UseInterfaces {
		options.UseInterfaces = globalDefaults.Project.UseInterfaces
	}
	options.UseInterfacesFromDefaults = true
}

// runInitWizard launches the browser-based project wizard, seeds it with the
// global defaults + the available templates, and writes the user's choices back
// into options. options.ProjectName is left empty if the user closes the wizard
// without creating.
func runInitWizard(options *flags.Init) error {
	globalDefaults, err := defaults.Load()
	if err != nil {
		// Don't apply on a load failure — that could reset fields (e.g. flip
		// UseInterfaces) from a half-loaded config.
		term.Warningf("Could not load global defaults: %v\n", err)
	} else {
		applyGlobalDefaults(options, globalDefaults)
	}

	// With no -n provided, the name-derived seeds are bogus (e.g. "A  application").
	// Clear them so the wizard derives sensible values from the entered name.
	if options.ProjectName == "" {
		options.ProductDescription = ""
		options.ProductIdentifier = ""
	}

	var initTemplates []setupwizard.InitTemplate
	for _, t := range templates.GetDefaultTemplates() {
		initTemplates = append(initTemplates, setupwizard.InitTemplate{Name: t.Name, Description: t.Description})
	}

	defaultTemplate := globalDefaults.GetTemplateName()
	if defaultTemplate == "" {
		defaultTemplate = options.TemplateName
	}

	baseDir, err := filepath.Abs(options.ProjectDir)
	if err != nil {
		baseDir = options.ProjectDir
	}

	productName := options.ProductName
	if productName == "" || productName == "My Product" {
		productName = options.ProjectName
	}

	data := setupwizard.InitData{
		ProjectName:        options.ProjectName,
		TemplateName:       defaultTemplate,
		ProductName:        productName,
		ProductCompany:     options.ProductCompany,
		ProductIdentifier:  options.ProductIdentifier,
		ProductDescription: options.ProductDescription,
		ProductVersion:     options.ProductVersion,
		ProductCopyright:   options.ProductCopyright,
		ProductComments:    options.ProductComments,
		UseInterfaces:      options.UseInterfaces,
		BaseDir:            baseDir,
		Templates:          initTemplates,
		DefaultTemplate:    defaultTemplate,
	}

	result, err := setupwizard.NewInitWizard(data).RunInit()
	if err != nil {
		return err
	}
	if result == nil {
		// User closed the wizard without creating; signal via empty name.
		options.ProjectName = ""
		return nil
	}

	options.ProjectName = result.ProjectName
	options.TemplateName = result.TemplateName
	options.ProductName = result.ProductName
	options.ProductCompany = result.ProductCompany
	options.ProductIdentifier = result.ProductIdentifier
	options.ProductDescription = result.ProductDescription
	options.ProductVersion = result.ProductVersion
	options.ProductCopyright = result.ProductCopyright
	options.ProductComments = result.ProductComments
	options.UseInterfaces = result.UseInterfaces
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

	// Interactive UI: collect the project name, template and config values in the
	// browser wizard, then fall through to the normal scaffolding below. Runs
	// before the project-name check because the wizard is where the name is set.
	if options.UI {
		if err := runInitWizard(options); err != nil {
			return err
		}
		if options.ProjectName == "" {
			// User closed the wizard without creating.
			return nil
		}
	}

	if options.ProjectName == "" {
		return errors.New("please use the -n flag to specify a project name")
	}

	options.ProjectName = sanitizeFileName(options.ProjectName)

	// Load and apply global defaults. Skipped in UI mode: the wizard already
	// seeded from defaults and the returned values are the user's explicit
	// choices, which a second pass could clobber (e.g. vanilla -> default).
	if !options.UI {
		globalDefaults, err := defaults.Load()
		if err != nil {
			// Log warning but continue - global defaults are optional
			term.Warningf("Could not load global defaults: %v\n", err)
		} else {
			applyGlobalDefaults(options, globalDefaults)
		}
	}

	// Determine the binding language AFTER global defaults are applied: when no
	// -t is given, applyGlobalDefaults may have just set options.TemplateName
	// from the wizard's configured default template, and that must drive the
	// TypeScript-vs-JavaScript bindings choice.
	isTypescript := templates.IsTypescript(options.TemplateName)

	if options.ModulePath == "" {
		if options.Git == "" {
			options.ModulePath = "changeme"
		} else {
			options.ModulePath = gitURLToModulePath(options.Git)
		}
	}

	err := templates.Install(options)
	if err != nil {
		return err
	}

	// Rename gitignore to .gitignore
	err = os.Rename(filepath.Join(options.ProjectDir, "gitignore"), filepath.Join(options.ProjectDir, ".gitignore"))
	if err != nil {
		return err
	}

	// Rename frontend/npmrc to frontend/.npmrc. Dotfiles can't be embedded, so
	// it ships without the leading dot. Optional, so failure is non-fatal.
	npmrcSrc := filepath.Join(options.ProjectDir, "frontend", "npmrc")
	if _, statErr := os.Stat(npmrcSrc); statErr == nil {
		_ = os.Rename(npmrcSrc, filepath.Join(options.ProjectDir, "frontend", ".npmrc"))
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

	// In UI mode the wizard is explicitly about the project config, so write the
	// chosen values into build/config.yml itself (the generated assets already
	// carry them; the source config.yml is otherwise a static copy).
	if options.UI {
		if err := writeProjectConfigYML(options); err != nil {
			term.Warningf("Could not update build/config.yml: %v\n", err)
		}
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

// writeProjectConfigYML rewrites the `info:` values in the freshly scaffolded
// build/config.yml from the chosen options, preserving the file's comments. The
// info keys map 1:1 onto the Product* fields.
func writeProjectConfigYML(options *flags.Init) error {
	path := filepath.Join(options.ProjectDir, "build", "config.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)

	values := map[string]string{
		"companyName":       options.ProductCompany,
		"productName":       options.ProductName,
		"productIdentifier": options.ProductIdentifier,
		"description":       options.ProductDescription,
		"copyright":         options.ProductCopyright,
		"comments":          options.ProductComments,
		"version":           options.ProductVersion,
	}
	for key, val := range values {
		if val == "" {
			continue
		}
		// Replace only the quoted value on the `  key: "..."` line, keeping any
		// trailing comment. Anchored to start-of-line so it won't touch the
		// commented ios: overrides.
		re := regexp.MustCompile(`(?m)^(\s*` + regexp.QuoteMeta(key) + `:\s*")[^"]*(")`)
		// Escape for a YAML double-quoted scalar — backslash first (so the escapes
		// added next aren't doubled), then quotes; collapse newlines to keep it a
		// single-line value; finally escape `$` for the regexp replacement template.
		repl := val
		repl = strings.ReplaceAll(repl, `\`, `\\`)
		repl = strings.ReplaceAll(repl, `"`, `\"`)
		repl = strings.ReplaceAll(repl, "\r", "")
		repl = strings.ReplaceAll(repl, "\n", " ")
		repl = strings.ReplaceAll(repl, `$`, `$$`)
		content = re.ReplaceAllString(content, `${1}`+repl+`${2}`)
	}

	return os.WriteFile(path, []byte(content), 0o644)
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
