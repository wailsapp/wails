package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/leaanthony/slicer"
)

// PackageManager indicates different package managers
type PackageManager int

const (
	// UNKNOWN package manager
	UNKNOWN PackageManager = iota
	// NPM package manager
	NPM
	// YARN package manager
	YARN
)

type author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type frontend struct {
	Dir     string `json:"dir"`
	Install string `json:"install"`
	Build   string `json:"build"`
	Bridge  string `json:"bridge"`
	Serve   string `json:"serve"`
}

type framework struct {
	Name     string            `json:"name"`
	BuildTag string            `json:"buildtag"`
	Options  map[string]string `json:"options,omitempty"`
}

// ProjectHelper is a helper struct for managing projects
type ProjectHelper struct {
	log       *Logger
	system    *SystemHelper
	templates *TemplateHelper
}

// NewProjectHelper creates a new Project helper struct
func NewProjectHelper() *ProjectHelper {
	return &ProjectHelper{
		log:       NewLogger(),
		system:    NewSystemHelper(),
		templates: NewTemplateHelper(),
	}
}

// GenerateProject generates a new project using the options given
func (ph *ProjectHelper) GenerateProject(projectOptions *ProjectOptions) error {

	// Calculate project path
	projectPath, err := filepath.Abs(projectOptions.OutputDirectory)
	if err != nil {
		return err
	}

	_ = projectPath

	if fs.DirExists(projectPath) {
		return fmt.Errorf("directory '%s' already exists", projectPath)
	}

	// Create project directory
	err = fs.MkDir(projectPath)
	if err != nil {
		return err
	}

	// Create and save project config
	err = projectOptions.WriteProjectConfig()
	if err != nil {
		return err
	}

	err = ph.templates.InstallTemplate(projectPath, projectOptions)
	if err != nil {
		return err
	}

	// // If we are on windows, dump a windows_resource.json
	// if runtime.GOOS == "windows" {
	// 	ph.GenerateWindowsResourceConfig(projectOptions)
	// }

	return nil
}

// // GenerateWindowsResourceConfig generates the default windows resource file
// func (ph *ProjectHelper) GenerateWindowsResourceConfig(po *ProjectOptions) {

// 	fmt.Println(buffer.String())

// 	// vi.Build()
// 	// vi.Walk()
// 	// err := vi.WriteSyso(outPath, runtime.GOARCH)
// }

// LoadProjectConfig loads the project config from the given directory
func (ph *ProjectHelper) LoadProjectConfig(dir string) (*ProjectOptions, error) {
	po := ph.NewProjectOptions()
	err := po.LoadConfig(dir)
	return po, err
}

// NewProjectOptions creates a new default set of project options
func (ph *ProjectHelper) NewProjectOptions() *ProjectOptions {
	result := ProjectOptions{
		Name:        "",
		Description: "Enter your project description",
		Version:     "0.1.0",
		BinaryName:  "",
		system:      ph.system,
		log:         ph.log,
		templates:   ph.templates,
		Author:      &author{},
	}

	// Populate system config
	config, err := ph.system.LoadConfig()
	if err == nil {
		result.Author.Name = config.Name
		result.Author.Email = config.Email
	}

	return &result
}

// ProjectOptions holds all the options available for a project
type ProjectOptions struct {
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	Author                 *author   `json:"author,omitempty"`
	Version                string    `json:"version"`
	OutputDirectory        string    `json:"-"`
	UseDefaults            bool      `json:"-"`
	Template               string    `json:"-"`
	BinaryName             string    `json:"binaryname"`
	FrontEnd               *frontend `json:"frontend,omitempty"`
	Tags                   string    `json:"tags"`
	NPMProjectName         string    `json:"-"`
	system                 *SystemHelper
	log                    *Logger
	templates              *TemplateHelper
	selectedTemplate       *TemplateDetails
	WailsVersion           string
	typescriptDefsFilename string
	Verbose                bool `json:"-"`
	CrossCompile           bool
	Platform               string
	Architecture           string
	LdFlags                string
	GoPath                 string
	UseFirebug             bool

	// Supported platforms
	Platforms []string `json:"platforms,omitempty"`
}

// PlatformSupported returns true if the template is supported
// on the current platform
func (po *ProjectOptions) PlatformSupported() bool {

	// Default is all platforms supported
	if len(po.Platforms) == 0 {
		return true
	}

	// Check that the platform is in the list
	platformsSupported := slicer.String(po.Platforms)
	return platformsSupported.Contains(runtime.GOOS)
}

// Defaults sets the default project template
func (po *ProjectOptions) Defaults() {
	po.Template = "vuebasic"
	po.WailsVersion = Version
}

// SetTypescriptDefsFilename indicates that we want to generate typescript bindings to the given file
func (po *ProjectOptions) SetTypescriptDefsFilename(filename string) {
	po.typescriptDefsFilename = filename
}

// GetNPMBinaryName returns the type of package manager used by the project
func (po *ProjectOptions) GetNPMBinaryName() (PackageManager, error) {
	if po.FrontEnd == nil {
		return UNKNOWN, fmt.Errorf("No frontend specified in project options")
	}

	if strings.Index(po.FrontEnd.Install, "npm") > -1 {
		return NPM, nil
	}

	if strings.Index(po.FrontEnd.Install, "yarn") > -1 {
		return YARN, nil
	}

	return UNKNOWN, nil
}

// PromptForInputs asks the user to input project details
func (po *ProjectOptions) PromptForInputs() error {

	processProjectName(po)

	processBinaryName(po)

	err := processOutputDirectory(po)
	if err != nil {
		return err
	}

	// Process Templates
	templateList := slicer.Interface()
	options := slicer.String()
	templateDetails, err := po.templates.GetTemplateDetails()
	if err != nil {
		return err
	}

	if po.Template != "" {
		// Check template is valid if given
		if templateDetails[po.Template] == nil {
			keys := make([]string, 0, len(templateDetails))
			for k := range templateDetails {
				keys = append(keys, k)
			}
			return fmt.Errorf("invalid template name '%s'. Valid options: %s", po.Template, strings.Join(keys, ", "))
		}
		po.selectedTemplate = templateDetails[po.Template]
	} else {

		keys := make([]string, 0)
		for k := range templateDetails {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			templateDetail := templateDetails[k]
			templateList.Add(templateDetail)
			if !templateDetail.Metadata.PlatformSupported() {
				templateDetail.Metadata.Name = "* " + templateDetail.Metadata.Name
			}
			options.Add(fmt.Sprintf("%s - %s", templateDetail.Metadata.Name, templateDetail.Metadata.ShortDescription))
		}

		templateIndex := 0

		if len(options.AsSlice()) > 1 {
			templateIndex = PromptSelection("Please select a template (* means unsupported on current platform)", options.AsSlice(), 0)
		}

		if len(templateList.AsSlice()) == 0 {
			return fmt.Errorf("aborting: no templates found")
		}

		// After selection do this....
		po.selectedTemplate = templateList.AsSlice()[templateIndex].(*TemplateDetails)
	}

	po.selectedTemplate.Metadata.Name = strings.TrimPrefix(po.selectedTemplate.Metadata.Name, "* ")
	if !po.selectedTemplate.Metadata.PlatformSupported() {
		println("WARNING: This template is unsupported on this platform!")
	}
	fmt.Println("Template: " + po.selectedTemplate.Metadata.Name)

	// Setup NPM Project name
	po.NPMProjectName = strings.ToLower(strings.Replace(po.Name, " ", "_", -1))

	// Fix template name
	po.Template = strings.Split(po.selectedTemplate.Path, string(os.PathSeparator))[0]

	// // Populate template details
	templateMetadata := po.selectedTemplate.Metadata

	err = processTemplateMetadata(templateMetadata, po)
	if err != nil {
		return err
	}

	return nil
}

// WriteProjectConfig writes the project configuration into
// the project directory
func (po *ProjectOptions) WriteProjectConfig() error {
	targetDir, err := filepath.Abs(po.OutputDirectory)
	if err != nil {
		return err
	}

	targetFile := filepath.Join(targetDir, "project.json")
	filedata, err := json.MarshalIndent(po, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(targetFile, filedata, 0600)
}

// LoadConfig loads the project configuration file from the
// given directory
func (po *ProjectOptions) LoadConfig(projectDir string) error {
	targetFile := filepath.Join(projectDir, "project.json")
	rawBytes, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(rawBytes, po)
}

func computeBinaryName(projectName string) string {
	if projectName == "" {
		return ""
	}
	var binaryNameComputed = strings.ToLower(projectName)
	binaryNameComputed = strings.Replace(binaryNameComputed, " ", "-", -1)
	binaryNameComputed = strings.Replace(binaryNameComputed, string(filepath.Separator), "-", -1)
	binaryNameComputed = strings.Replace(binaryNameComputed, ":", "-", -1)
	return binaryNameComputed
}

func processOutputDirectory(po *ProjectOptions) error {
	// po.OutputDirectory
	if po.OutputDirectory == "" {
		po.OutputDirectory = PromptRequired("Project directory name", computeBinaryName(po.Name))
	}
	projectPath, err := filepath.Abs(po.OutputDirectory)
	if err != nil {
		return err
	}

	if NewFSHelper().DirExists(projectPath) {
		return fmt.Errorf("directory '%s' already exists", projectPath)
	}

	fmt.Println("Project Directory: " + po.OutputDirectory)
	return nil
}

func processProjectName(po *ProjectOptions) {
	if po.Name == "" {
		po.Name = Prompt("The name of the project", "My Project")
	}
	fmt.Println("Project Name: " + po.Name)
}

func processBinaryName(po *ProjectOptions) {
	if po.BinaryName == "" {
		var binaryNameComputed = computeBinaryName(po.Name)
		po.BinaryName = Prompt("The output binary name", binaryNameComputed)
	}
	fmt.Println("Output binary Name: " + po.BinaryName)
}

func processTemplateMetadata(templateMetadata *TemplateMetadata, po *ProjectOptions) error {
	if templateMetadata.FrontendDir != "" {
		po.FrontEnd = &frontend{}
		po.FrontEnd.Dir = templateMetadata.FrontendDir
	}
	if templateMetadata.Install != "" {
		if po.FrontEnd == nil {
			return fmt.Errorf("install set in template metadata but not frontenddir")
		}
		po.FrontEnd.Install = templateMetadata.Install
	}
	if templateMetadata.Build != "" {
		if po.FrontEnd == nil {
			return fmt.Errorf("build set in template metadata but not frontenddir")
		}
		po.FrontEnd.Build = templateMetadata.Build
	}

	if templateMetadata.Bridge != "" {
		if po.FrontEnd == nil {
			return fmt.Errorf("bridge set in template metadata but not frontenddir")
		}
		po.FrontEnd.Bridge = templateMetadata.Bridge
	}

	if templateMetadata.Serve != "" {
		if po.FrontEnd == nil {
			return fmt.Errorf("serve set in template metadata but not frontenddir")
		}
		po.FrontEnd.Serve = templateMetadata.Serve
	}

	// Save platforms
	po.Platforms = templateMetadata.Platforms

	return nil
}
