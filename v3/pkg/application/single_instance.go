package application

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var alreadyRunningError = errors.New("application is already running")

// SecondInstanceData contains information about the second instance launch
type SecondInstanceData struct {
	Args           []string          `json:"args"`
	WorkingDir     string            `json:"workingDir"`
	AdditionalData map[string]string `json:"additionalData,omitempty"`
}

// SingleInstanceOptions defines options for single instance functionality
type SingleInstanceOptions struct {
	// UniqueID is used to identify the application instance
	// This should be unique per application, e.g. "com.myapp.myapplication"
	UniqueID string

	// OnSecondInstanceLaunch is called when a second instance of the application is launched
	// The callback receives data about the second instance launch
	OnSecondInstanceLaunch func(data SecondInstanceData)

	// AdditionalData allows passing custom data from second instance to first
	AdditionalData map[string]string

	// ExitCode is the exit code to use when the second instance exits
	ExitCode int
}

// platformLock is the interface that platform-specific lock implementations must implement
type platformLock interface {
	// acquire attempts to acquire the lock
	acquire(uniqueID string) error
	// release releases the lock and cleans up resources
	release()
	// notify sends data to the first instance
	notify(data string) error
}

// singleInstanceManager handles the single instance functionality
type singleInstanceManager struct {
	options *SingleInstanceOptions
	lock    platformLock
	app     *App
}

func newSingleInstanceManager(app *App, options *SingleInstanceOptions) (*singleInstanceManager, error) {
	if options == nil {
		return nil, nil
	}

	manager := &singleInstanceManager{
		options: options,
		app:     app,
	}

	// Create platform-specific lock
	lock, err := newPlatformLock(manager)
	if err != nil {
		return nil, err
	}

	manager.lock = lock

	// Try to acquire the lock
	err = lock.acquire(options.UniqueID)
	if err != nil {
		return manager, err
	}

	return manager, nil
}

func (m *singleInstanceManager) cleanup() {
	if m == nil || m.lock == nil {
		return
	}
	m.lock.release()
}

// notifyFirstInstance sends data to the first instance of the application
func (m *singleInstanceManager) notifyFirstInstance() error {
	data := SecondInstanceData{
		Args:           os.Args,
		WorkingDir:     getCurrentWorkingDir(),
		AdditionalData: m.options.AdditionalData,
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return m.lock.notify(string(serialized))
}

func getCurrentWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// getLockPath returns the path to the lock file for Unix systems
func getLockPath(uniqueID string) string {
	// Use system temp directory
	tmpDir := os.TempDir()
	lockFileName := uniqueID + ".lock"
	return filepath.Join(tmpDir, lockFileName)
}
