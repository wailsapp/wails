package browser

import (
	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// ---------------- Plugin Setup ----------------
// This is the main plugin struct. It can be named anything you like.
// It must implement the application.Plugin interface.
// Both the Init() and Shutdown() methods are called synchronously when the app starts and stops.

type Plugin struct{}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Shutdown() {}

func (p *Plugin) Name() string {
	return "github.com/wailsapp/wails/v3/plugins/browser"
}

func (p *Plugin) Init(_ *application.App) error {
	return nil
}

func (p *Plugin) CallableByJS() []string {
	return []string{
		"OpenURL",
		"OpenFile",
	}
}

func (p *Plugin) InjectJS() string {
	return ""
}

// ---------------- Plugin Methods ----------------

func (p *Plugin) OpenURL(url string) error {
	return browser.OpenURL(url)
}

func (p *Plugin) OpenFile(path string) error {
	return browser.OpenFile(path)
}
