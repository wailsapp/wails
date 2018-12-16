package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey"
	homedir "github.com/mitchellh/go-homedir"
)

type SystemHelper struct {
	log               *Logger
	fs                *FSHelper
	configFilename    string
	homeDir           string
	wailsSystemDir    string
	wailsSystemConfig string
}

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

// Initialise attempts to set up the Wails system.
// An error is returns if there is a problem
func (s *SystemHelper) Initialise() error {

	// System dir doesn't exist
	if !s.systemDirExists() {
		s.log.Green("Welcome to Wails!")
		s.log.Green("To get you set up, I'll need to ask you a few things...")
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

type SystemConfig struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewSystemConfig(filename string) (*SystemConfig, error) {
	result := &SystemConfig{}
	err := result.load(filename)
	return result, err
}

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
