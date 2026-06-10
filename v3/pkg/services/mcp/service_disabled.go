//go:build !mcp

package mcp

import (
	"context"
	"net/http"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const enabled = false

// mcpServer is a stub. The real implementation is compiled in with the `mcp`
// build tag.
type mcpServer struct{}

func (s *Service) startup(_ context.Context, _ application.ServiceOptions) error {
	if app := application.Get(); app != nil {
		app.Logger.Debug("MCP service is disabled. Build with the 'mcp' tag (or set WAILS_MCP=1 when building) to enable it.")
	}
	return nil
}

func (s *Service) shutdown() error {
	return nil
}

func (s *Service) serveHTTP(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "MCP service is disabled. Build with the 'mcp' tag to enable it.", http.StatusNotFound)
}
