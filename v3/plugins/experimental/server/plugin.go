package server

import (
	_ "embed"
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed plugin.js
var pluginJS string

//go:embed client.js
var clientJS string

type Config struct {
	Host    string
	Port    int
	Enabled bool
	Headers map[string]string
}

func (c Config) ListenAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type Plugin struct {
	config *Config
	server *Server
}

func NewPlugin(config *Config) *Plugin {
	return &Plugin{
		config: config,
		server: NewServer(config),
	}
}

func (s *Plugin) CallableByJS() []string {
	return []string{} // maybe # clients?
}

func (p *Plugin) InjectJS() string {
	return pluginJS
}

// Init is called when the plugin is loaded. It is passed the application.App
// instance. This is where you should do any setup.
func (p *Plugin) Init() error {
	p.server.app = application.Get()
	p.server.run()
	return nil
}

// Shutdown will stop the server
func (s *Plugin) Shutdown() {
	s.server.Shutdown()
}

// Name returns the name of the plugin.
func (s *Plugin) Name() string {
	return "github.com/wailsapp/wails/v3/plugins/server"
}
