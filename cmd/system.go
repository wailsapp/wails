package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
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
	s.homeDir, err = os.UserHomeDir()
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

// GetAuthor returns a formatted string of the user's name and email
func (s *SystemHelper) GetAuthor() (string, error) {
	var config *SystemConfig
	config, err := s.LoadConfig()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s <%s>", config.Name, config.Email), nil
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

	systemConfig := make(map[string]string)

	// Try to load current values - ignore errors
	config, _ := s.LoadConfig()

	if config.Name != "" {
		systemConfig["name"] = PromptRequired("What is your name", config.Name)
	} else if n, err := getGitConfigValue("user.name"); err == nil && n != "" {
		systemConfig["name"] = PromptRequired("What is your name", n)
	} else {
		systemConfig["name"] = PromptRequired("What is your name")
	}

	if config.Email != "" {
		systemConfig["email"] = PromptRequired("What is your email address", config.Email)
	} else if e, err := getGitConfigValue("user.email"); err == nil && e != "" {
		systemConfig["email"] = PromptRequired("What is your email address", e)
	} else {
		systemConfig["email"] = PromptRequired("What is your email address")
	}

	// Create the directory
	err := s.fs.MkDirs(s.wailsSystemDir)
	if err != nil {
		return err
	}

	// Save
	configData, err := json.Marshal(&systemConfig)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(s.wailsSystemConfig, configData, 0755)
	if err != nil {
		return err
	}
	fmt.Println()
	s.log.White("Wails config saved to: " + s.wailsSystemConfig)
	s.log.White("Feel free to customise these settings.")
	fmt.Println()

	return nil
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

// SystemConfig - Defines system wide configuration data
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

// CheckDependenciesSilent checks for dependencies but
// only outputs if there's an error
func CheckDependenciesSilent(logger *Logger) (bool, error) {
	logger.SetErrorOnly(true)
	result, err := CheckDependencies(logger)
	logger.SetErrorOnly(false)
	return result, err
}

// CheckDependencies will look for Wails dependencies on the system
// Errors are reported in error and the bool return value is whether
// the dependencies are all installed.
func CheckDependencies(logger *Logger) (bool, error) {

	switch runtime.GOOS {
	case "darwin":
		logger.Yellow("Detected Platform: OSX")
	case "windows":
		logger.Yellow("Detected Platform: Windows")
	case "linux":
		logger.Yellow("Detected Platform: Linux")
	default:
		return false, fmt.Errorf("Platform %s is currently not supported", runtime.GOOS)
	}

	logger.Yellow("Checking for prerequisites...")
	// Check we have a cgo capable environment

	requiredPrograms, err := GetRequiredPrograms()
	if err != nil {
		return false, nil
	}
	errors := false
	programHelper := NewProgramHelper()
	for _, program := range *requiredPrograms {
		bin := programHelper.FindProgram(program.Name)
		if bin == nil {
			errors = true
			logger.Error("Program '%s' not found. %s", program.Name, program.Help)
		} else {
			logger.Green("Program '%s' found: %s", program.Name, bin.Path)
		}
	}

	// Linux has library deps
	if runtime.GOOS == "linux" {
		// Check library prerequisites
		requiredLibraries, err := GetRequiredLibraries()
		if err != nil {
			return false, err
		}

		var libraryChecker CheckPkgInstalled
		distroInfo := GetLinuxDistroInfo()

		switch distroInfo.Distribution {
		case Ubuntu, Debian, Zorin, Parrot, Linuxmint, Elementary, Kali, Neon, Deepin, Raspbian, PopOS:
			libraryChecker = DpkgInstalled
		case Arch, ArcoLinux, ArchLabs, Ctlos, Manjaro, ManjaroARM, EndeavourOS:
			libraryChecker = PacmanInstalled
		case CentOS, Fedora, Tumbleweed, Leap, RHEL:
			libraryChecker = RpmInstalled
		case Gentoo:
			libraryChecker = EqueryInstalled
		case VoidLinux:
			libraryChecker = XbpsInstalled
		case Solus:
			libraryChecker = EOpkgInstalled
		case Crux:
			libraryChecker = PrtGetInstalled
		default:
			return false, RequestSupportForDistribution(distroInfo)
		}

		for _, library := range *requiredLibraries {
			installed, err := libraryChecker(library.Name)
			if err != nil {
				return false, err
			}
			if !installed {
				errors = true
				logger.Error("Library '%s' not found. %s", library.Name, library.Help)
			} else {
				logger.Green("Library '%s' installed.", library.Name)
			}
		}
	}
	logger.White("")

	return !errors, err
}
