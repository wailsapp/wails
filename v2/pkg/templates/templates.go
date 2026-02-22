package templates

import (
	"embed"
	"encoding/json"
	"fmt"
	gofs "io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"

	"github.com/leaanthony/debme"
	"github.com/leaanthony/gosod"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
)

//go:embed all:templates
var templates embed.FS

//go:embed all:ides/*
var ides embed.FS

// Cahce for the templates
// We use this because we need different views of the same data
var templateCache []Template = nil

// Data contains the data we wish to embed during template installation
type Data struct {
	ProjectName        string
	BinaryName         string
	WailsVersion       string
	NPMProjectName     string
	AuthorName         string
	AuthorEmail        string
	AuthorNameAndEmail string
	WailsDirectory     string
	GoSDKPath          string
	WindowsFlags       string
	CGOEnabled         string
	OutputFile         string
}

// Options for installing a template
type Options struct {
	ProjectName         string
	TemplateName        string
	BinaryName          string
	TargetDir           string
	Logger              *clilogger.CLILogger
	PathToDesktopBinary string
	PathToServerBinary  string
	InitGit             bool
	AuthorName          string
	AuthorEmail         string
	IDE                 string
	ProjectNameFilename string // The project name but as a valid filename
	WailsVersion        string
	GoSDKPath           string
	WindowsFlags        string
	CGOEnabled          string
	CGOLDFlags          string
	OutputFile          string
}

// Template holds data relating to a template
// including the metadata stored in template.json
type Template struct {
	// Template details
	Name        string `json:"name"`
	ShortName   string `json:"shortname"`
	Author      string `json:"author"`
	Description string `json:"description"`
	HelpURL     string `json:"helpurl"`

	// Other data
	FS gofs.FS `json:"-"`
}

func parseTemplate(template gofs.FS) (Template, error) {
	var result Template
	data, err := gofs.ReadFile(template, "template.json")
	if err != nil {
		return result, errors.Wrap(err, "Error parsing template")
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	result.FS = template
	return result, nil
}

// List returns the list of available templates
func List() ([]Template, error) {
	// If the cache isn't loaded, load it
	if templateCache == nil {
		err := loadTemplateCache()
		if err != nil {
			return nil, err
		}
	}

	return templateCache, nil
}

// getTemplateByShortname returns the template with the given short name
func getTemplateByShortname(shortname string) (Template, error) {
	var result Template

	// If the cache isn't loaded, load it
	if templateCache == nil {
		err := loadTemplateCache()
		if err != nil {
			return result, err
		}
	}

	for _, template := range templateCache {
		if template.ShortName == shortname {
			return template, nil
		}
	}

	return result, fmt.Errorf("shortname '%s' is not a valid template shortname", shortname)
}

// Loads the template cache
func loadTemplateCache() error {
	templatesFS, err := debme.FS(templates, "templates")
	if err != nil {
		return err
	}

	// Get directories
	files, err := templatesFS.ReadDir(".")
	if err != nil {
		return err
	}

	// Reset cache
	templateCache = []Template{}

	for _, file := range files {
		if file.IsDir() {
			templateFS, err := templatesFS.FS(file.Name())
			if err != nil {
				return err
			}
			template, err := parseTemplate(templateFS)
			if err != nil {
				// Cannot parse this template, continue
				continue
			}
			templateCache = append(templateCache, template)
		}
	}

	return nil
}

// Install the given template. Returns true if the template is remote.
func Install(options *Options) (bool, *Template, error) {
	// Get cwd
	cwd, err := os.Getwd()
	if err != nil {
		return false, nil, err
	}

	// Did the user want to install in current directory?
	if options.TargetDir == "" {
		options.TargetDir = filepath.Join(cwd, options.ProjectName)
		if fs.DirExists(options.TargetDir) {
			return false, nil, fmt.Errorf("cannot create project directory. Dir exists: %s", options.TargetDir)
		}
	} else {
		// Get the absolute path of the given directory
		targetDir, err := filepath.Abs(options.TargetDir)
		if err != nil {
			return false, nil, err
		}
		options.TargetDir = targetDir
		if !fs.DirExists(options.TargetDir) {
			err := fs.Mkdir(options.TargetDir)
			if err != nil {
				return false, nil, err
			}
		}
	}

	// Flag to indicate remote template
	remoteTemplate := false

	// Is this a shortname?
	template, err := getTemplateByShortname(options.TemplateName)
	if err != nil {
		// Is this a filepath?
		templatePath, err := filepath.Abs(options.TemplateName)
		if fs.DirExists(templatePath) {
			templateFS := os.DirFS(templatePath)
			template, err = parseTemplate(templateFS)
			if err != nil {
				return false, nil, errors.Wrap(err, "Error installing template")
			}
		} else {
			// git clone to temporary dir
			tempdir, err := gitclone(options)
			defer func(path string) {
				err := os.RemoveAll(path)
				if err != nil {
					log.Fatal(err)
				}
			}(tempdir)
			if err != nil {
				return false, nil, err
			}
			// Remove the .git directory
			err = os.RemoveAll(filepath.Join(tempdir, ".git"))
			if err != nil {
				return false, nil, err
			}

			templateFS := os.DirFS(tempdir)
			template, err = parseTemplate(templateFS)
			if err != nil {
				return false, nil, err
			}
			remoteTemplate = true
		}
	}

	// Use Gosod to install the template
	installer := gosod.New(template.FS)

	// Ignore template.json files
	installer.IgnoreFile("template.json")

	// Setup the data.
	// We use the directory name for the binary name, like Go
	BinaryName := filepath.Base(options.TargetDir)
	NPMProjectName := strings.ToLower(strings.ReplaceAll(BinaryName, " ", ""))
	localWailsDirectory := fs.RelativePath("../../../../../..")

	templateData := &Data{
		ProjectName:    options.ProjectName,
		BinaryName:     filepath.Base(options.TargetDir),
		NPMProjectName: NPMProjectName,
		WailsDirectory: localWailsDirectory,
		AuthorEmail:    options.AuthorEmail,
		AuthorName:     options.AuthorName,
		WailsVersion:   options.WailsVersion,
		GoSDKPath:      options.GoSDKPath,
	}

	// Create a formatted name and email combo.
	if options.AuthorName != "" {
		templateData.AuthorNameAndEmail = options.AuthorName + " "
	}
	if options.AuthorEmail != "" {
		templateData.AuthorNameAndEmail += "<" + options.AuthorEmail + ">"
	}
	templateData.AuthorNameAndEmail = strings.TrimSpace(templateData.AuthorNameAndEmail)

	installer.RenameFiles(map[string]string{
		"gitignore.txt": ".gitignore",
	})

	// Extract the template
	err = installer.Extract(options.TargetDir, templateData)
	if err != nil {
		return false, nil, err
	}

	err = generateIDEFiles(options)
	if err != nil {
		return false, nil, err
	}

	return remoteTemplate, &template, nil
}

// Clones the given uri and returns the temporary cloned directory
func gitclone(options *Options) (string, error) {
	// Create temporary directory
	dirname, err := os.MkdirTemp("", "wails-template-*")
	if err != nil {
		return "", err
	}

	// Parse remote template url and version number
	templateInfo := strings.Split(options.TemplateName, "@")
	cloneOption := &git.CloneOptions{
		URL: templateInfo[0],
	}
	if len(templateInfo) > 1 {
		cloneOption.ReferenceName = plumbing.NewTagReferenceName(templateInfo[1])
	}

	_, err = git.PlainClone(dirname, false, cloneOption)

	return dirname, err
}

func generateIDEFiles(options *Options) error {
	switch options.IDE {
	case "vscode":
		return generateVSCodeFiles(options)
	case "goland":
		return generateGolandFiles(options)
	}

	return nil
}

type ideOptions struct {
	name         string
	targetDir    string
	options      *Options
	renameFiles  map[string]string
	ignoredFiles []string
}

func generateGolandFiles(options *Options) error {
	ideoptions := ideOptions{
		name:      "goland",
		targetDir: filepath.Join(options.TargetDir, ".idea"),
		options:   options,
		renameFiles: map[string]string{
			"projectname.iml": options.ProjectNameFilename + ".iml",
			"gitignore.txt":   ".gitignore",
			"name":            ".name",
		},
	}
	if !options.InitGit {
		ideoptions.ignoredFiles = []string{"vcs.xml"}
	}
	err := installIDEFiles(ideoptions)
	if err != nil {
		return errors.Wrap(err, "generating Goland IDE files")
	}

	return nil
}

func generateVSCodeFiles(options *Options) error {
	ideoptions := ideOptions{
		name:      "vscode",
		targetDir: filepath.Join(options.TargetDir, ".vscode"),
		options:   options,
	}
	return installIDEFiles(ideoptions)
}

func installIDEFiles(o ideOptions) error {
	source, err := debme.FS(ides, "ides/"+o.name)
	if err != nil {
		return err
	}

	// Use gosod to install the template
	installer := gosod.New(source)

	if o.renameFiles != nil {
		installer.RenameFiles(o.renameFiles)
	}

	for _, ignoreFile := range o.ignoredFiles {
		installer.IgnoreFile(ignoreFile)
	}

	binaryName := filepath.Base(o.options.TargetDir)
	o.options.WindowsFlags = ""
	o.options.CGOEnabled = "1"

	switch runtime.GOOS {
	case "windows":
		binaryName += ".exe"
		o.options.WindowsFlags = " -H windowsgui"
		o.options.CGOEnabled = "0"
	case "darwin":
		o.options.CGOLDFlags = "-framework UniformTypeIdentifiers"
	}

	o.options.PathToDesktopBinary = filepath.ToSlash(filepath.Join("build", "bin", binaryName))

	err = installer.Extract(o.targetDir, o.options)
	if err != nil {
		return err
	}

	return nil
}
