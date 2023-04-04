package single_instance

import (
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/application"
	"os"
	"path/filepath"
)

type Config struct {
	// Add any configuration options here
	LockFileName                  string
	LockFilePath                  string
	ActivateAppOnSubsequentLaunch bool
}

type Plugin struct {
	config   *Config
	app      *application.App
	lockfile *os.File
}

func (p *Plugin) CallableByJS() []string {
	return []string{}
}

func (p *Plugin) InjectJS() string {
	return ""
}

func NewPlugin(config *Config) *Plugin {
	if config.LockFilePath == "" {
		// Use the system default temp directory
		config.LockFilePath = os.TempDir()
	}
	if config.LockFileName == "" {
		// Use the executable name
		config.LockFileName = filepath.Base(os.Args[0]) + ".lock"
	}
	return &Plugin{
		config: config,
	}
}

// Shutdown is called when the app is shutting down
func (p *Plugin) Shutdown() {
	p.lockfile.Close()
}

// Name returns the name of the plugin.
func (p *Plugin) Name() string {
	return "github.com/wailsapp/wails/v3/plugins/single-instance"
}

// Init is called when the app is starting up. You can use this to
// initialise any resources you need. You can also access the application
// instance via the app property.
func (p *Plugin) Init(app *application.App) error {
	p.app = app

	var err error
	lockfileName := p.config.LockFilePath + "/" + p.config.LockFileName
	p.lockfile, err = CreateLockFile(lockfileName, p.app.GetPID())
	if err != nil {
		if p.config.ActivateAppOnSubsequentLaunch {
			pid, err := GetLockFilePid(lockfileName)
			if err != nil {
				return err
			}
			err = p.activeInstance(pid)
			if err != nil {
				return err
			}
		}
		return fmt.Errorf("another instance of this application is already running")
	}
	return nil
}

// Exported returns a list of exported methods that can be called from the frontend
func (p *Plugin) Exported() []string {
	return []string{}
}
