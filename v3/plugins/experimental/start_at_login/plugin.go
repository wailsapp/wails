package start_at_login

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"io/fs"
)

type Plugin struct {
	disabled bool
	options  Config
}

type Config struct {
	// RegistryKey is the key in the registry to use for storing the start at login setting.
	// This defaults to the name of the executable
	RegistryKey string
}

func NewPlugin(options Config) *Plugin {
	return &Plugin{
		options: options,
	}
}

// Shutdown is called when the app is shutting down
// You can use this to clean up any resources you have allocated
func (p *Plugin) Shutdown() error { return nil }

// Name returns the name of the plugin.
// You should use the go module format e.g. github.com/myuser/myplugin
func (p *Plugin) Name() string {
	return "github.com/wailsapp/wails/v3/plugins/start_at_login"
}

func (p *Plugin) Init(api application.PluginAPI) error {
	// OS specific initialiser
	err := p.init()
	if err != nil {
		return err
	}
	return nil
}

func (p *Plugin) CallableByJS() []string {
	return []string{
		"StartAtLogin",
		"IsStartAtLogin",
	}
}

func (p *Plugin) Assets() fs.FS {
	return nil
}
