package fileserver

import (
	"net/http"
	"sync/atomic"
)

type Config struct {
	// RootPath specifies the filesystem path from which requests are to be served.
	RootPath string
}

type Service struct {
	fs atomic.Pointer[http.Handler]
}

// New initialises an unconfigured fileserver. See [Configure] for details.
func New() *Service {
	return NewWithConfig(nil)
}

// New initialises and optionally configures a fileserver. See [Service.Configure] for details.
func NewWithConfig(config *Config) *Service {
	result := &Service{}
	result.Configure(config)
	return result
}

// ServiceName returns the name of the plugin.
// You should use the go module format e.g. github.com/myuser/myplugin
func (s *Service) ServiceName() string {
	return "github.com/wailsapp/wails/v3/services/fileserver"
}

// Configure reconfigures the fileserver.
// If config is nil, then every request will receive a 503 Service Unavailable response.
//
//wails:ignore
func (s *Service) Configure(config *Config) {
	if config == nil {
		s.fs.Store(&dummyHandler)
	} else {
		var fs http.Handler = http.FileServer(http.Dir(config.RootPath))
		s.fs.Store(&fs)
	}
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	(*s.fs.Load()).ServeHTTP(w, r)
}

var dummyHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Fileserver service has not been configured yet", http.StatusServiceUnavailable)
})
