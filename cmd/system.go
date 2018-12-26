package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey"
	"github.com/leaanthony/spinner"
	homedir "github.com/mitchellh/go-homedir"
)

// SystemHelper - Defines everything related to the system
type SystemHelper struct {
	log               *Logger
	fs                *FSHelper
	configFilename    string
	homeDir           string
	wailsSystemDir    string
	wailsSystemConfig string
}

// NewSystemHelper - Creates a new System Helper
func NewSystemHelper() *SystemHelper {
	result := &SystemHelper{
		fs:             NewFSHelper(),
		log:            NewLogger(),
		configFilename: "wails.json",
	}
	result.setSystemDirs()
	return result
}

// Internal
// setSystemDirs calculates the system directories it is interested in
func (s *SystemHelper) setSystemDirs() {
	var err error
	s.homeDir, err = homedir.Dir()
	if err != nil {
		log.Fatal("Cannot find home directory! Please file a bug report!")
	}
	// TODO: A better config system
	s.wailsSystemDir = filepath.Join(s.homeDir, ".wails")
	s.wailsSystemConfig = filepath.Join(s.wailsSystemDir, s.configFilename)
}

// ConfigFileExists - Returns true if it does!
func (s *SystemHelper) ConfigFileExists() bool {
	return s.fs.FileExists(s.wailsSystemConfig)
}

// SystemDirExists - Returns true if it does!
func (s *SystemHelper) systemDirExists() bool {
	return s.fs.DirExists(s.wailsSystemDir)
}

// LoadConfig attempts to load the Wails system config
func (s *SystemHelper) LoadConfig() (*SystemConfig, error) {
	return NewSystemConfig(s.wailsSystemConfig)
}

// ConfigFileIsValid checks if the config file is valid
func (s *SystemHelper) ConfigFileIsValid() bool {
	_, err := NewSystemConfig(s.wailsSystemConfig)
	return err == nil
}

// BackupConfig attempts to backup the system config file
func (s *SystemHelper) BackupConfig() (string, error) {
	now := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	backupFilename := s.wailsSystemConfig + "." + now
	err := s.fs.CopyFile(s.wailsSystemConfig, backupFilename)
	if err != nil {
		return "", err
	}
	return backupFilename, nil
}

func (s *SystemHelper) setup() error {

	// Answers. We all need them.
	answers := &SystemConfig{}

	// Try to load current values - ignore errors
	config, err := s.LoadConfig()
	defaultName := ""
	defaultEmail := ""
	if config != nil {
		defaultName = config.Name
		defaultEmail = config.Email
	}
	// Questions
	var simpleQs = []*survey.Question{
		{
			Name: "Name",
			Prompt: &survey.Input{
				Message: "What is your name:",
				Default: defaultName,
			},
			Validate: survey.Required,
		},
		{
			Name: "Email",
			Prompt: &survey.Input{
				Message: "What is your email address:",
				Default: defaultEmail,
			},
			Validate: survey.Required,
		},
	}

	// ask the questions
	err = survey.Ask(simpleQs, answers)
	if err != nil {
		return err
	}

	// Create the directory
	err = s.fs.MkDirs(s.wailsSystemDir)
	if err != nil {
		return err
	}

	fmt.Println()
	s.log.White("Wails config saved to: " + s.wailsSystemConfig)
	s.log.White("Feel free to customise these settings.")
	fmt.Println()

	return answers.Save(s.wailsSystemConfig)
}

const introText = `
Wails is a lightweight framework for creating web-like desktop apps in Go. 
I'll need to ask you a few questions so I can fill in your project templates and then I will try and see if you have the correct dependencies installed. If you don't have the right tools installed, I'll try and suggest how to install them.
`

// CheckInitialised checks if the system has been set up
// and if not, runs setup
func (s *SystemHelper) CheckInitialised() error {
	if !s.systemDirExists() {
		s.log.Yellow("System not initialised. Running setup.")
		return s.setup()
	}
	return nil
}

// Initialise attempts to set up the Wails system.
// An error is returns if there is a problem
func (s *SystemHelper) Initialise() error {

	// System dir doesn't exist
	if !s.systemDirExists() {
		s.log.Green("Welcome to Wails!")
		s.log.Green(introText)
		return s.setup()
	}

	// Config doesn't exist
	if !s.ConfigFileExists() {
		s.log.Green("Looks like the system config is missing.")
		s.log.Green("To get you back on track, I'll need to ask you a few things...")
		return s.setup()
	}

	// Config exists but isn't valid.
	if !s.ConfigFileIsValid() {
		s.log.Green("Looks like the system config got corrupted.")
		backupFile, err := s.BackupConfig()
		if err != nil {
			s.log.Green("I tried to backup your config file but got this error: %s", err.Error())
		} else {
			s.log.Green("Just in case you needed it, I backed up your config file here: %s", backupFile)
		}
		s.log.Green("To get you back on track, I'll need to ask you a few things...")
		return s.setup()
	}

	return s.setup()
}

// SystemConfig - Defines system wode configuration data
type SystemConfig struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// NewSystemConfig - Creates a new SystemConfig helper object
func NewSystemConfig(filename string) (*SystemConfig, error) {
	result := &SystemConfig{}
	err := result.load(filename)
	return result, err
}

// Save - Saves the system config to the given filename
func (sc *SystemConfig) Save(filename string) error {
	// Convert config to JSON string
	theJSON, err := json.MarshalIndent(sc, "", "  ")
	if err != nil {
		return err
	}

	// Write it out to the config file
	return ioutil.WriteFile(filename, theJSON, 0644)
}

func (sc *SystemConfig) load(filename string) error {
	configData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	// Load and unmarshall!
	err = json.Unmarshal(configData, &sc)
	if err != nil {
		return err
	}
	return nil
}

// CheckDependencies will look for Wails dependencies on the system
// Errors are reported in error and the bool return value is whether
// the dependencies are all installed.
func CheckDependencies(logger *Logger) (error, bool) {

	switch runtime.GOOS {
	case "darwin":
		logger.Yellow("Detected Platform: OSX")
	case "windows":
		logger.Yellow("Detected Platform: Windows")
	case "linux":
		logger.Yellow("Detected Platform: Linux")
	default:
		return fmt.Errorf("Platform %s is currently not supported", runtime.GOOS), false
	}

	logger.Yellow("Checking for prerequisites...")
	// Check we have a cgo capable environment

	requiredPrograms, err := GetRequiredPrograms()
	if err != nil {
		return nil, false
	}
	errors := false
	spinner := spinner.New()
	programHelper := NewProgramHelper()
	for _, program := range *requiredPrograms {
		spinner.Start("Looking for program '%s'", program.Name)
		bin := programHelper.FindProgram(program.Name)
		if bin == nil {
			errors = true
			spinner.Errorf("Program '%s' not found. %s", program.Name, program.Help)
		} else {
			spinner.Successf("Program '%s' found: %s", program.Name, bin.Path)
		}
	}

	// Linux has library deps
	if runtime.GOOS == "linux" {
		// Check library prerequisites
		requiredLibraries, err := GetRequiredLibraries()
		if err != nil {
			return err, false
		}
		distroInfo := GetLinuxDistroInfo()
		for _, library := range *requiredLibraries {
			spinner.Start()
			switch distroInfo.Distribution {
			case Ubuntu:
				installed, err := DpkgInstalled(library.Name)
				if err != nil {
					return err, false
				}
				if !installed {
					errors = true
					spinner.Errorf("Library '%s' not found. %s", library.Name, library.Help)
				} else {
					spinner.Successf("Library '%s' installed.", library.Name)
				}
			default:
				return fmt.Errorf("unable to check libraries on distribution '%s'. Please ensure that the '%s' equivalent is installed", distroInfo.DistributorID, library.Name), false
			}
		}
	}

	return err, !errors
}
