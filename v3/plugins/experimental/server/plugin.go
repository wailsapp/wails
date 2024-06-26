package server

import (
	_ "embed"
	"fmt"
	"io/fs"
)

//go:embed plugin.js
var pluginJS string

//go:embed ipc_websocket.js
var clientJS string

type Config struct {
	Host    string
	Port    int
	Enabled bool
	Headers map[string]string
	Assets  fs.FS
}

func (c Config) ListenAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type plugin struct {
	config *Config
	server *Server
}

func NewPlugin(config *Config) *plugin {
	return &plugin{
		config: config,
		server: NewServer(config),
	}
}

func (s *plugin) Assets() fs.FS {
	return nil
}

func (s *plugin) CallableByJS() []string {
	return []string{} // maybe # clients?
}

func (p *plugin) InjectJS() string {
	return ""
	//return clientJS
}

// Init is called when the plugin is loaded. It is passed the application.App
// instance. This is where you should do any setup.
func (p *plugin) Init() error {
	p.server.run()
	return nil
}

// Shutdown will stop the server
func (s *plugin) Shutdown() error {
	s.server.Shutdown()
	return nil
}

// Name returns the name of the plugin.
func (s *plugin) Name() string {
	return "server"
}
