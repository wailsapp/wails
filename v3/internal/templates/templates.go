package templates

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v3/internal/buildinfo"
	"github.com/wailsapp/wails/v3/internal/version"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/debug"

	"github.com/wailsapp/wails/v3/internal/flags"

	"github.com/leaanthony/gosod"

	"github.com/samber/lo"
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

type TemplateOptions struct {
	*flags.Init
	LocalModulePath string
	UseTypescript   bool
	WailsVersion    string
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

	return &template, nil
}

func getLocalTemplate(templateName string) (*Template, error) {
	var template Template
	var err error
	_, err = os.Stat(templateName)
	if err != nil {
		return nil, nil
	}

	template, err = parseTemplate(os.DirFS(templateName), templateName)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

// Template holds data relating to a template including the metadata stored in template.yaml
type Template struct {

	// Template details
	Name        string `json:"name"`
	ShortName   string `json:"shortname"`
	Author      string `json:"author"`
	Description string `json:"description"`
	HelpURL     string `json:"helpurl"`
	Version     int8   `json:"version"`

	// Other data
	FS fs.FS `json:"-"`
}

func parseTemplate(template fs.FS, templateName string) (Template, error) {
	var result Template
	jsonFile := "template.json"
	if templateName != "" {
		jsonFile = templateName + "/template.json"
	}
	data, err := fs.ReadFile(template, jsonFile)
	if err != nil {
		return result, errors.Wrap(err, "Error parsing template")
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	result.FS = template

	// We need to do a version check here
	if result.Version == 0 {
		return result, fmt.Errorf("template not supported by wails 3. This template is probably for wails 2")
	}
	if result.Version != 3 {
		return result, fmt.Errorf("template version %d is not supported by wails 3. Ensure 'version' is set to 3 in the `template.json` file", result.Version)
	}

	return result, nil
}

// Clones the given uri and returns the temporary cloned directory
func gitclone(uri string) (string, error) {
	// Create temporary directory
	dirname, err := os.MkdirTemp("", "wails-template-*")
	if err != nil {
		return "", err
	}

	// Parse remote template url and version number
	templateInfo := strings.Split(uri, "@")
	cloneOption := &git.CloneOptions{
		URL: templateInfo[0],
	}
	if len(templateInfo) > 1 {
		cloneOption.ReferenceName = plumbing.NewTagReferenceName(templateInfo[1])
	}

	_, err = git.PlainClone(dirname, false, cloneOption)

	return dirname, err

}

func getRemoteTemplate(uri string) (template *Template, err error) {
	// git clone to temporary dir
	var tempDir string
	tempDir, err = gitclone(uri)

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tempDir)

	if err != nil {
		return
	}

	// Remove the .git directory
	err = os.RemoveAll(filepath.Join(tempDir, ".git"))
	if err != nil {
		return
	}

	templateFS := os.DirFS(tempDir)
	var parsedTemplate Template
	parsedTemplate, err = parseTemplate(templateFS, "")
	if err != nil {
		return
	}
	return &parsedTemplate, nil
}

func Install(options *flags.Init) error {
	var wd string = lo.Must(os.Getwd())
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
	UseTypescript := strings.HasSuffix(options.TemplateName, "-ts")

	templateData := TemplateOptions{
		Init:            options,
		LocalModulePath: localModulePath,
		UseTypescript:   UseTypescript,
		WailsVersion:    version.VersionString,
	}

	defer func() {
		// if `template.json` exists, remove it
		_ = os.Remove(filepath.Join(templateData.ProjectDir, "template.json"))
	}()

	var template *Template
	template, err = getInternalTemplate(options.TemplateName)
	if err != nil {
		return err
	}
	if template == nil {
		template, err = getLocalTemplate(options.TemplateName)
	}
	if err != nil {
		return err
	}
	if template == nil {
		template, err = getRemoteTemplate(options.TemplateName)
	}
	if err != nil {
		return err
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

	pterm.Printf("Creating project\n")
	pterm.Printf("----------------\n\n")
	table := pterm.TableData{
		{"Project Name", options.ProjectName},
		{"Project Directory", filepath.FromSlash(options.ProjectDir)},
		{"Template", template.Name},
		{"Template Source", template.HelpURL},
	}
	err = pterm.DefaultTable.WithData(table).Render()
	if err != nil {
		return err
	}
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

	// Change to project directory
	err = os.Chdir(templateData.ProjectDir)
	if err != nil {
		return err
	}
	// Run `go mod tidy`
	err = exec.Command("go", "mod", "tidy").Run()
	if err != nil {
		return err
	}

	pterm.Printf("\nProject '%s' created successfully.\n", options.ProjectName)

	return nil

}
