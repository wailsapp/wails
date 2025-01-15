package fileserver

import (
	"context"
	"net/http"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// ---------------- Service Setup ----------------
// This is the main Service struct. It can be named anything you like.

type Config struct {
	RootPath string
}

type Service struct {
	config *Config
	fs     http.Handler
}

func New(config *Config) *Service {
	return &Service{
		config: config,
		fs:     http.FileServer(http.Dir(config.RootPath)),
	}
}

// ServiceShutdown is called when the app is shutting down
// You can use this to clean up any resources you have allocated
func (s *Service) ServiceShutdown() error {
	return nil
}

// ServiceName returns the name of the plugin.
// You should use the go module format e.g. github.com/myuser/myplugin
func (s *Service) ServiceName() string {
	return "github.com/wailsapp/wails/v3/services/fileserver"
}

// ServiceStartup is called when the app is starting up. You can use this to
// initialise any resources you need.
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// Any initialization code here
	return nil
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create a new file server rooted at the given path
	s.fs.ServeHTTP(w, r)
}
