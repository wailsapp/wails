package templates

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/internal/buildinfo"
	"github.com/wailsapp/wails/v3/internal/s"
	"github.com/wailsapp/wails/v3/internal/term"
	"github.com/wailsapp/wails/v3/internal/version"
	"gopkg.in/yaml.v3"

	"errors"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/git"
	"github.com/wailsapp/wails/v3/internal/debug"

	"github.com/wailsapp/wails/v3/internal/flags"

	"github.com/wailsapp/wails/v3/internal/gosod"

	"github.com/wailsapp/wails/v3/internal/lo"
)

//go:embed *
var templates embed.FS

type TemplateData struct {
	Name        string
	Description string
	FS          fs.FS
}

var defaultTemplates = []TemplateData{}

func init() {
	dirs, err := templates.ReadDir(".")
	if err != nil {
		return
	}
	for _, dir := range dirs {
		if strings.HasPrefix(dir.Name(), "_") {
			continue
		}
		if dir.IsDir() {
			template, err := parseTemplate(templates, dir.Name())
			if err != nil {
				continue
			}
			defaultTemplates = append(defaultTemplates,
				TemplateData{
					Name:        dir.Name(),
					Description: template.Description,
					FS:          templates,
				})
		}
	}
}

func ValidTemplateName(name string) bool {
	return lo.ContainsBy(defaultTemplates, func(template TemplateData) bool {
		return template.Name == name
	})
}

func GetDefaultTemplates() []TemplateData {
	return defaultTemplates
}

// IsTypescript reports whether the named template produces a TypeScript project.
//
// Built-in templates declare this explicitly via `typescript:` in their
// template.yaml (TypeScript now owns the bare framework name, e.g. `react`,
// while JavaScript variants carry a `-js` suffix, e.g. `react-js`). For local
// and remote templates that predate the flag we fall back to the historical
// `-ts` suffix convention so community templates keep working.
func IsTypescript(name string) bool {
	if strings.HasSuffix(name, "-ts") {
		return true
	}
	if strings.HasSuffix(name, "-js") {
		return false
	}
	if ValidTemplateName(name) {
		if tmpl, err := getInternalTemplate(name); err == nil {
			return tmpl.Typescript
		}
	}
	return false
}

// NormalizeBinaryName converts a project name into a valid binary/package name:
// lowercased, with spaces replaced by dashes, and any remaining characters not
// in [a-z0-9-] replaced by dashes, collapsing runs and trimming edges.
func NormalizeBinaryName(name string) string {
	name = strings.ToLower(name)
	var b strings.Builder
	prevDash := false
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			prevDash = false
		} else if !prevDash {
			b.WriteRune('-')
			prevDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

type TemplateOptions struct {
	Cls string `description:"A helper for using close template tags safely }}" default:"}}"`
	Opn string `description:"A helper for using open template tags safely {{" default:"{{"`
	*flags.Init
	BinaryName      string
	LocalModulePath string
	UseTypescript   bool
	UseInterfaces   bool
	WailsVersion    string
}

type FileAssociation struct {
	Ext         string `yaml:"ext"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	IconName    string `yaml:"iconName"`
	Role        string `yaml:"role"`
	MimeType    string `yaml:"mimeType"`
}

type ProtocolConfig struct {
	Scheme      string `yaml:"scheme"                json:"scheme"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

func getInternalTemplate(templateName string) (*Template, error) {
	templateData, found := lo.Find(defaultTemplates, func(template TemplateData) bool {
		return template.Name == templateName
	})

	if !found {
		return nil, nil
	}

	template, err := parseTemplate(templateData.FS, templateData.Name)
	if err != nil {
		return nil, err
	}
	template.source = sourceInternal
	return &template, nil
}

func getLocalTemplate(templateName string) (*Template, error) {
	var template Template
	var err error
	_, err = os.Stat(templateName)
	if err != nil {
		return nil, nil
	}

	template, err = parseTemplate(os.DirFS(templateName), "")
	if err != nil {
		return nil, err
	}
	template.source = sourceLocal

	return &template, nil
}

type BaseTemplate struct {
	Name         string `json:"name"         yaml:"name"         description:"The name of the template"`
	ShortName    string `json:"shortname"    yaml:"shortname"    description:"The short name of the template"`
	Author       string `json:"author"       yaml:"author"       description:"The author of the template"`
	Description  string `json:"description"  yaml:"description"  description:"The template description"`
	HelpURL      string `json:"helpurl"      yaml:"helpurl"      description:"The help url for the template"`
	Version      string `json:"version"      yaml:"version"      description:"The version of the template" default:"v0.0.1"`
	WailsVersion uint8  `json:"wailsVersion" yaml:"wailsVersion" description:"The Wails major version this template targets"`
	Typescript   bool   `json:"typescript"   yaml:"typescript"   description:"Whether the template produces a TypeScript project"`
	Dir          string `json:"-"            yaml:"-"            description:"The directory to generate the template" default:"."`
	Frontend     string `json:"-"            yaml:"-"            description:"The frontend directory to migrate"`
}

type source int

const (
	sourceInternal source = 1
	sourceLocal    source = 2
	sourceRemote   source = 3
)

// Template holds data relating to a template including the metadata stored in template.yaml
type Template struct {
	BaseTemplate `yaml:",inline"`
	// Schema is kept for backwards-compatible reading of legacy template.json files.
	// New templates use template.yaml with WailsVersion instead.
	Schema uint8 `json:"schema" yaml:"-"`

	// Other data
	FS      fs.FS `json:"-" yaml:"-"`
	source  source
	tempDir string
}

// parseTemplate loads and validates a template's metadata.
//
// It first looks for template.yaml (the v3 native format). If found, wailsVersion
// must be present and equal to 3. If only template.json exists (legacy format),
// schema must be 3 for a v3 template; schema 0 means the template is for Wails v2.
func parseTemplate(templateFS fs.FS, templateName string) (Template, error) {
	var result Template

	prefix := ""
	if templateName != "" {
		prefix = templateName + "/"
	}

	// --- YAML path: preferred v3 format ---
	yamlData, yamlErr := fs.ReadFile(templateFS, prefix+"template.yaml")
	if yamlErr == nil {
		if err := yaml.Unmarshal(yamlData, &result); err != nil {
			return result, fmt.Errorf("error parsing template.yaml: %w", err)
		}
		result.FS = templateFS
		if result.WailsVersion == 0 {
			return result, fmt.Errorf("template.yaml must specify 'wailsVersion' (e.g. wailsVersion: 3)")
		}
		if result.WailsVersion != 3 {
			return result, fmt.Errorf("template targets Wails v%d and is not compatible with this version of Wails", result.WailsVersion)
		}
		return result, nil
	}

	// --- JSON path: legacy / backwards-compat ---
	jsonData, jsonErr := fs.ReadFile(templateFS, prefix+"template.json")
	if jsonErr != nil {
		if errors.Is(yamlErr, fs.ErrNotExist) && errors.Is(jsonErr, fs.ErrNotExist) {
			return result, fmt.Errorf("no template.yaml or template.json found in template")
		}
		if !errors.Is(yamlErr, fs.ErrNotExist) {
			return result, fmt.Errorf("error reading template.yaml: %w", yamlErr)
		}
		return result, fmt.Errorf("error reading template.json: %w", jsonErr)
	}

	if err := json.Unmarshal(jsonData, &result); err != nil {
		return result, fmt.Errorf("error parsing template.json: %w", err)
	}
	result.FS = templateFS

	if result.Schema == 0 {
		return result, fmt.Errorf("template not supported by Wails v3: no schema or wailsVersion found. This template is probably for Wails v2")
	}
	if result.Schema != 3 {
		return result, fmt.Errorf("template schema %d is not supported by Wails v3. Ensure 'schema' is set to 3 in template.json", result.Schema)
	}

	return result, nil
}

// gitclone clones uri into a temporary directory and returns its path.
func gitclone(uri string) (string, error) {
	dirname, err := os.MkdirTemp("", "wails-template-*")
	if err != nil {
		return "", err
	}

	parts := strings.SplitN(uri, "@", 2)
	url, tag := parts[0], ""
	if len(parts) > 1 {
		tag = parts[1]
	}

	return dirname, git.Clone(url, dirname, tag)
}

func getRemoteTemplate(uri string) (*Template, error) {
	tempDir, err := gitclone(uri)
	if err != nil {
		return nil, err
	}

	// cleanup is called on any error path so the temp dir is never leaked.
	cleanup := func() { os.RemoveAll(tempDir) }

	if err = os.RemoveAll(filepath.Join(tempDir, ".git")); err != nil {
		cleanup()
		return nil, err
	}

	parsedTemplate, err := parseTemplate(os.DirFS(tempDir), "")
	if err != nil {
		cleanup()
		return nil, err
	}

	parsedTemplate.tempDir = tempDir
	parsedTemplate.source = sourceRemote
	return &parsedTemplate, nil
}

func Install(options *flags.Init) error {
	var wd = lo.Must(os.Getwd())
	var projectDir string
	if options.ProjectDir == "." || options.ProjectDir == "" {
		projectDir = wd
	} else {
		projectDir = options.ProjectDir
	}
	var err error
	projectDir, err = filepath.Abs(filepath.Join(projectDir, options.ProjectName))
	if err != nil {
		return err
	}

	buildInfo, err := buildinfo.Get()
	if err != nil {
		return err
	}

	// Calculate relative path from project directory to LocalModulePath
	var localModulePath string

	// Use module path if it is set
	if buildInfo.Development {
		var relativePath string
		// Check if the project directory and LocalModulePath are in the same drive
		if filepath.VolumeName(wd) != filepath.VolumeName(debug.LocalModulePath) {
			relativePath = debug.LocalModulePath
		} else {
			relativePath, err = filepath.Rel(projectDir, debug.LocalModulePath)
		}
		if err != nil {
			return err
		}
		localModulePath = filepath.ToSlash(relativePath + "/")
	}
	UseTypescript := IsTypescript(options.TemplateName)

	templateData := TemplateOptions{
		Init:            options,
		BinaryName:      NormalizeBinaryName(options.ProjectName),
		LocalModulePath: localModulePath,
		UseTypescript:   UseTypescript,
		UseInterfaces:   options.UseInterfaces,
		WailsVersion:    version.String(),
		Opn:             "{{",
		Cls:             "}}",
	}

	defer func() {
		// Remove template metadata files from the generated project — they are
		// for the template system only, not part of the user's project.
		_ = os.Remove(filepath.Join(templateData.ProjectDir, "template.yaml"))
		_ = os.Remove(filepath.Join(templateData.ProjectDir, "template.json"))
	}()

	var template *Template

	if ValidTemplateName(options.TemplateName) {
		template, err = getInternalTemplate(options.TemplateName)
		if err != nil {
			return err
		}
	} else {
		template, err = getLocalTemplate(options.TemplateName)
		if err != nil {
			return err
		}
		if template == nil {
			template, err = getRemoteTemplate(options.TemplateName)
		}
	}

	if template == nil {
		return fmt.Errorf("invalid template name: %s. Use -l flag to view available templates or use a valid filepath / url to a template", options.TemplateName)
	}

	templateData.ProjectDir = projectDir

	// If project directory already exists and is not empty, error
	if _, err := os.Stat(templateData.ProjectDir); !os.IsNotExist(err) {
		// Check if the directory is empty
		files := lo.Must(os.ReadDir(templateData.ProjectDir))
		if len(files) > 0 {
			return fmt.Errorf("project directory '%s' already exists and is not empty", templateData.ProjectDir)
		}
	}

	if template.source == sourceRemote && !options.SkipWarning {
		var confirmed = confirmRemote(template)
		if !confirmed {
			return nil
		}
	}

	term.Section("Project")

	language := "JavaScript"
	if UseTypescript {
		language = "TypeScript"
	}

	framework := strings.TrimSuffix(strings.TrimSuffix(options.TemplateName, "-ts"), "-js")
	if len(framework) > 0 {
		framework = strings.ToUpper(framework[:1]) + framework[1:]
	}

	frameworkDisplay := framework
	languageDisplay := language
	if options.TemplateFromDefaults {
		frameworkDisplay += " (default)"
		languageDisplay += " (default)"
	}

	rows := [][]string{
		{"Name", options.ProjectName},
		{"Directory", filepath.FromSlash(options.ProjectDir)},
		{"Framework", frameworkDisplay},
		{"Language", languageDisplay},
	}

	if UseTypescript {
		bindingStyle := "Classes"
		if options.UseInterfaces {
			bindingStyle = "Interfaces"
		}
		if options.UseInterfacesFromDefaults {
			bindingStyle += " (default)"
		}
		rows = append(rows, []string{"Bindings", bindingStyle})
	}

	fmt.Print(term.RenderTable(rows))

	switch template.source {
	case sourceInternal:
		tfs, err := fs.Sub(template.FS, options.TemplateName)
		if err != nil {
			return err
		}
		common, err := fs.Sub(templates, "_common")
		if err != nil {
			return err
		}
		err = gosod.New(common).Extract(options.ProjectDir, templateData)
		if err != nil {
			return err
		}
		err = gosod.New(tfs).Extract(options.ProjectDir, templateData)
		if err != nil {
			return err
		}
	case sourceLocal, sourceRemote:
		publisher := fmt.Sprintf("CN=%s", options.ProductCompany)
		data := struct {
		TemplateOptions
		Dir                   string
		Name                  string
		BinaryName            string
		ProductName           string
		ProductDescription    string
		ProductVersion        string
		ProductCompany        string
		ProductCopyright      string
		ProductComments       string
		ProductIdentifier     string
		Publisher             string
		ProcessorArchitecture string
		ExecutableName        string
		ExecutablePath        string
		OutputPath            string
		CertificatePath       string
		FileAssociations      []FileAssociation
		Protocols             []ProtocolConfig
		Silent                bool
		Typescript            bool
	}{
		Name:                  options.ProjectName,
		BinaryName:            NormalizeBinaryName(options.ProjectName),
		Silent:                true,
		ProductCompany:        options.ProductCompany,
		ProductName:           options.ProductName,
		ProductDescription:    options.ProductDescription,
		ProductVersion:        options.ProductVersion,
		ProductIdentifier:     options.ProductIdentifier,
		ProductCopyright:      options.ProductCopyright,
		ProductComments:       options.ProductComments,
		Publisher:             publisher,
		ProcessorArchitecture: "x64",
		ExecutableName:        options.ProjectName,
		ExecutablePath:        options.ProjectName,
		OutputPath:            fmt.Sprintf("%s.msix", options.ProjectName),
		CertificatePath:       "",
		FileAssociations:      []FileAssociation{},
		Protocols:             []ProtocolConfig{},
		Typescript:            templateData.UseTypescript,
		TemplateOptions:       templateData,
	}
		// If options.ProjectDir does not exist, create it
		if _, err := os.Stat(options.ProjectDir); os.IsNotExist(err) {
			err = os.Mkdir(options.ProjectDir, 0755)
			if err != nil {
				return err
			}
		}
		err = gosod.New(template.FS).Extract(options.ProjectDir, data)
		if err != nil {
			return err
		}

		if template.tempDir != "" {
			s.RMDIR(template.tempDir)
		}
	}
	if !options.SkipGoModTidy {
		err = goModTidy(templateData.ProjectDir)
		if err != nil {
			return err
		}
	}

	// Change to project directory
	err = os.Chdir(templateData.ProjectDir)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(term.OkStyle.Render("✓") + " Project '" + options.ProjectName + "' created successfully.")

	return nil

}

func GenerateTemplate(options *BaseTemplate) error {
	if options.Name == "" {
		return fmt.Errorf("please provide a template name using the -name flag")
	}

	baseOutputDir, err := filepath.Abs(options.Dir)
	if err != nil {
		return err
	}
	outDir := filepath.Join(baseOutputDir, options.Name)

	if _, err := os.Stat(outDir); err == nil {
		return fmt.Errorf("directory '%s' already exists", outDir)
	} else if !os.IsNotExist(err) {
		return err
	}

	// Copy the common files (Go backend, Taskfile, go.mod, etc.) verbatim.
	// These files contain template variables like {{.ProjectName}} that must be
	// preserved so they are expanded when users later run `wails init -t <template>`.
	commonFS, err := fs.Sub(templates, "_common")
	if err != nil {
		return err
	}
	if err = os.CopyFS(outDir, commonFS); err != nil {
		return err
	}

	// Replace the placeholder frontend directory with the real frontend content.
	frontendDir := filepath.Join(outDir, "frontend")
	if err = os.RemoveAll(frontendDir); err != nil {
		return err
	}
	if options.Frontend != "" {
		if err = os.CopyFS(frontendDir, os.DirFS(options.Frontend)); err != nil {
			return fmt.Errorf("failed to copy frontend from '%s': %w", options.Frontend, err)
		}
	} else {
		baseFrontendFS, err := fs.Sub(templates, "base/frontend")
		if err != nil {
			return err
		}
		if err = os.CopyFS(frontendDir, baseFrontendFS); err != nil {
			return err
		}
	}

	// Copy NEXTSTEPS.md from the embedded base directory.
	nextstepsData, err := templates.ReadFile("base/NEXTSTEPS.md")
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(outDir, "NEXTSTEPS.md"), nextstepsData, 0644); err != nil {
		return err
	}

	// Write template.yaml with the provided metadata.
	// wailsVersion is mandatory in the YAML format and set to 3.
	options.WailsVersion = 3
	optionsYAML, err := yaml.Marshal(&options)
	if err != nil {
		return err
	}
	const modeline = "# yaml-language-server: $schema=https://v3.wails.io/schemas/template.v3.json\n"
	if err = os.WriteFile(filepath.Join(outDir, "template.yaml"), append([]byte(modeline), optionsYAML...), 0644); err != nil {
		return err
	}

	fmt.Printf("Template '%s' generated in %s\n", options.Name, outDir)
	fmt.Printf("See NEXTSTEPS.md for guidance on customising and publishing your template.\n")
	return nil
}

// stripUnsafe removes all control characters from untrusted strings before
// printing them to the terminal, preventing ANSI injection and multiline spoofing.
func stripUnsafe(s string) string {
	return strings.Map(func(r rune) rune {
		if r < 0x20 || (r >= 0x7f && r <= 0x9f) {
			return -1
		}
		return r
	}, s)
}

func confirmRemote(template *Template) bool {
	pterm.Println(pterm.LightRed("\n⚠  THIRD-PARTY TEMPLATE WARNING ⚠"))
	pterm.Println()
	pterm.Println(pterm.LightYellow("You are about to create a project from a remote template:"))
	pterm.Printf("  Name:   %s\n", stripUnsafe(template.Name))
	pterm.Printf("  Author: %s\n", stripUnsafe(template.Author))
	if template.HelpURL != "" {
		pterm.Printf("  URL:    %s\n", stripUnsafe(template.HelpURL))
	}
	pterm.Println()
	pterm.Println(pterm.LightYellow("Remote templates are third-party code. The Wails project does not review,"))
	pterm.Println(pterm.LightYellow("endorse, or accept any responsibility for their contents."))
	pterm.Println(pterm.LightYellow("Only proceed if you trust the source of this template."))
	pterm.Println()

	result, _ := pterm.DefaultInteractiveConfirm.
		WithDefaultText("Do you accept responsibility for using this third-party template?").
		WithConfirmText("y").
		WithRejectText("n").
		Show()

	return result
}

// goModTidy runs go mod tidy in the given project directory
// It returns an error if the command fails
func goModTidy(projectDir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w\n%s", err, string(output))
	}
	return nil
}
