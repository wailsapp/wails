package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey"
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

	fs := NewFSHelper()
	exists, err := ph.templates.TemplateExists(projectOptions.Template)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("template '%s' is invalid", projectOptions.Template)
	}

	// Calculate project path
	projectPath, err := filepath.Abs(projectOptions.OutputDirectory)
	if err != nil {
		return err
	}

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
	ph.log.Yellow("Project '%s' generated in directory '%s'!", projectOptions.Name, projectOptions.OutputDirectory)
	ph.log.Yellow("To compile the project, run 'wails build' in the project directory.")
	return nil
}

// LoadProjectConfig loads the project config from the given directory
func (ph *ProjectHelper) LoadProjectConfig(dir string) (*ProjectOptions, error) {
	po := ph.NewProjectOptions()
	err := po.LoadConfig(dir)
	return po, err
}

// NewProjectOptions creates a new default set of project options
func (ph *ProjectHelper) NewProjectOptions() *ProjectOptions {
	result := ProjectOptions{
		Name:            "",
		Description:     "Enter your project description",
		Version:         "0.1.0",
		BinaryName:      "",
		system:          NewSystemHelper(),
		log:             NewLogger(),
		templates:       NewTemplateHelper(),
		templateNameMap: make(map[string]string),
		Author:          &author{},
	}

	// Populate system config
	config, err := ph.system.LoadConfig()
	if err == nil {
		result.Author.Name = config.Name
		result.Author.Email = config.Email
	}

	return &result
}

// SelectQuestion creates a new select type question for Survey
func SelectQuestion(name, message string, options []string, defaultValue string, required bool) *survey.Question {
	result := survey.Question{
		Name: name,
		Prompt: &survey.Select{
			Message: message,
			Options: options,
			Default: defaultValue,
		},
	}
	if required {
		result.Validate = survey.Required
	}
	return &result
}

// InputQuestion creates a new input type question for Survey
func InputQuestion(name, message string, defaultValue string, required bool) *survey.Question {
	result := survey.Question{
		Name: name,
		Prompt: &survey.Input{
			Message: message + ":",
			Default: defaultValue,
		},
	}
	if required {
		result.Validate = survey.Required
	}
	return &result
}

// ProjectOptions holds all the options available for a project
type ProjectOptions struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Author          *author   `json:"author,omitempty"`
	Version         string    `json:"version"`
	OutputDirectory string    `json:"-"`
	UseDefaults     bool      `json:"-"`
	Template        string    `json:"-"`
	BinaryName      string    `json:"binaryname"`
	FrontEnd        *frontend `json:"frontend,omitempty"`
	NPMProjectName  string    `json:"-"`
	system          *SystemHelper
	log             *Logger
	templates       *TemplateHelper
	templateNameMap map[string]string // Converts template prompt text to template name
}

// Defaults sets the default project template
func (po *ProjectOptions) Defaults() {
	po.Template = "basic"
}

// PromptForInputs asks the user to input project details
func (po *ProjectOptions) PromptForInputs() error {

	var questions []*survey.Question

	if po.Name == "" {
		questions = append(questions, InputQuestion("Name", "The name of the project", "My Project", true))
	} else {
		fmt.Println("Project Name: " + po.Name)
	}

	if po.BinaryName == "" {
		var binaryNameComputed = computeBinaryName(po.Name)
		questions = append(questions, InputQuestion("BinaryName", "The output binary name", binaryNameComputed, true))
	} else {
		fmt.Println("Output binary Name: " + po.BinaryName)
	}

	err := processOutputDirectory(po.OutputDirectory, &questions)
	if err != nil {
		return err
	}

	templateDetails, err := po.templates.GetTemplateDetails()
	if err != nil {
		return err
	}

	templates := []string{}
	// Add a Custom Template
	// templates = append(templates, "Custom - Choose your own CSS framework")
	for templateName, templateDetails := range templateDetails {
		templateText := templateName
		// Check if metadata json exists
		if templateDetails.Metadata != nil {
			shortdescription := templateDetails.Metadata["shortdescription"]
			if shortdescription != "" {
				templateText += " - " + shortdescription.(string)
			}
		}
		templates = append(templates, templateText)
		po.templateNameMap[templateText] = templateName
	}

	if po.Template != "" {
		if _, ok := templateDetails[po.Template]; !ok {
			po.log.Error("Template '%s' invalid.", po.Template)
			questions = append(questions, SelectQuestion("Template", "Select template", templates, templates[0], true))
		}
	} else {
		questions = append(questions, SelectQuestion("Template", "Select template", templates, templates[0], true))
	}

	err = survey.Ask(questions, po)
	if err != nil {
		return err
	}

	// Setup NPM Project name
	po.NPMProjectName = strings.ToLower(strings.Replace(po.Name, " ", "_", -1))

	// Fix template name
	if po.templateNameMap[po.Template] != "" {
		po.Template = po.templateNameMap[po.Template]
	}

	// Populate template details
	templateMetadata := templateDetails[po.Template].Metadata
	if templateMetadata["frontenddir"] != nil {
		po.FrontEnd = &frontend{}
		po.FrontEnd.Dir = templateMetadata["frontenddir"].(string)
	}
	if templateMetadata["install"] != nil {
		if po.FrontEnd == nil {
			return fmt.Errorf("install set in template metadata but not frontenddir")
		}
		po.FrontEnd.Install = templateMetadata["install"].(string)
	}
	if templateMetadata["build"] != nil {
		if po.FrontEnd == nil {
			return fmt.Errorf("build set in template metadata but not frontenddir")
		}
		po.FrontEnd.Build = templateMetadata["build"].(string)
	}

	if templateMetadata["bridge"] != nil {
		if po.FrontEnd == nil {
			return fmt.Errorf("bridge set in template metadata but not frontenddir")
		}
		po.FrontEnd.Bridge = templateMetadata["bridge"].(string)
	}

	if templateMetadata["serve"] != nil {
		if po.FrontEnd == nil {
			return fmt.Errorf("serve set in template metadata but not frontenddir")
		}
		po.FrontEnd.Serve = templateMetadata["serve"].(string)
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

func processOutputDirectory(outputDirectory string, questions *[]*survey.Question) error {

	if outputDirectory != "" {
		projectPath, err := filepath.Abs(outputDirectory)
		if err != nil {
			return err
		}

		if NewFSHelper().DirExists(projectPath) {
			return fmt.Errorf("directory '%s' already exists", projectPath)
		}

		fmt.Println("Project Directory: " + outputDirectory)
	} else {
		*questions = append(*questions, InputQuestion("OutputDirectory", "Project directory name", "", true))
	}
	return nil
}
