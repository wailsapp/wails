package start_at_login

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Plugin struct {
	app      *application.App
	disabled bool
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

// Shutdown is called when the app is shutting down
// You can use this to clean up any resources you have allocated
func (p *Plugin) Shutdown() {}

// Name returns the name of the plugin.
// You should use the go module format e.g. github.com/myuser/myplugin
func (p *Plugin) Name() string {
	return "github.com/wailsapp/wails/v3/plugins/start_at_login"
}

func (p *Plugin) Init(app *application.App) error {
	p.app = app
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

func (p *Plugin) InjectJS() string {
	return ""
}
