//go:build dev
// +build dev

package main

import (
	"fmt"
	"time"

	"github.com/wailsapp/wails/cmd"
	"gopkg.in/AlecAivazis/survey.v1"
)

var templateHelper = cmd.NewTemplateHelper()

var qs = []*survey.Question{
	{
		Name:     "Name",
		Prompt:   &survey.Input{Message: "Please enter the name of your template (eg: React/Webpack Basic):"},
		Validate: survey.Required,
	},
	{
		Name:     "ShortDescription",
		Prompt:   &survey.Input{Message: "Please enter a short description for the template (eg: React with Webpack 4):"},
		Validate: survey.Required,
	},
	{
		Name:     "Description",
		Prompt:   &survey.Input{Message: "Please enter a long description:"},
		Validate: survey.Required,
	},
	{
		Name:     "FrontendDir",
		Prompt:   &survey.Input{Message: "Please enter the name of the directory the frontend code resides (eg: frontend):"},
		Validate: survey.Required,
	},
	{
		Name:     "Install",
		Prompt:   &survey.Input{Message: "Please enter the install command (eg: npm install):"},
		Validate: survey.Required,
	},
	{
		Name:     "Build",
		Prompt:   &survey.Input{Message: "Please enter the build command (eg: npm run build):"},
		Validate: survey.Required,
	},
	{
		Name:     "Serve",
		Prompt:   &survey.Input{Message: "Please enter the serve command (eg: npm run serve):"},
		Validate: survey.Required,
	},
	{
		Name:     "Bridge",
		Prompt:   &survey.Input{Message: "Please enter the name of the directory to copy the wails bridge runtime (eg: src):"},
		Validate: survey.Required,
	},
}

func newTemplate(devCommand *cmd.Command) {

	commandDescription := `This command scaffolds everything needed to develop a new template.`
	newTemplate := devCommand.Command("newtemplate", "Generate a new template").
		LongDescription(commandDescription)

	newTemplate.Action(func() error {
		logger.PrintSmallBanner("Generating new project template")
		fmt.Println()

		var answers cmd.TemplateMetadata

		// perform the questions
		err := survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		dirname := templateHelper.SanitizeFilename(answers.Name)
		prompt := []*survey.Question{{
			Prompt: &survey.Input{
				Message: "Please enter a directory name for the template:",
				Default: dirname,
			},
			Validate: func(val interface{}) error {
				err := survey.Required(val)
				if err != nil {
					return err
				}
				if templateHelper.IsValidTemplate(val.(string)) {
					return fmt.Errorf("template directory already exists")
				}
				if templateHelper.SanitizeFilename(val.(string)) != val.(string) {
					return fmt.Errorf("invalid directory name '%s'", val.(string))
				}
				return nil
			},
		}}
		err = survey.Ask(prompt, &dirname)
		if err != nil {
			return err
		}

		answers.Version = "1.0.0"
		answers.Created = time.Now().String()

		// Get Author info from system info
		system := cmd.NewSystemHelper()
		author, err := system.GetAuthor()
		if err == nil {
			answers.Author = author
		}

		templateDirectory, err := templateHelper.CreateNewTemplate(dirname, &answers)
		if err != nil {
			return err
		}

		logger.Green("Created new template '%s' in directory '%s'", answers.Name, templateDirectory)

		return nil
	})
}
