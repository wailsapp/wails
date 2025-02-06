package assetserver

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

const (
	HeaderHost          = "Host"
	HeaderContentType   = "Content-Type"
	HeaderContentLength = "Content-Length"
	HeaderUserAgent     = "User-Agent"
	// TODO: Is this needed?
	HeaderCacheControl = "Cache-Control"
	HeaderUpgrade      = "Upgrade"

	WailsUserAgentValue = "wails.io"
)

type assetServerLogger struct{}

var assetServerLoggerKey assetServerLogger

func ServeFile(rw http.ResponseWriter, filename string, blob []byte) error {
	header := rw.Header()
	header.Set(HeaderContentLength, fmt.Sprintf("%d", len(blob)))
	if mimeType := header.Get(HeaderContentType); mimeType == "" {
		mimeType = GetMimetype(filename, blob)
		header.Set(HeaderContentType, mimeType)
	}

	rw.WriteHeader(http.StatusOK)
	_, err := io.Copy(rw, bytes.NewReader(blob))
	return err
}

func isWebSocket(req *http.Request) bool {
	upgrade := req.Header.Get(HeaderUpgrade)
	return strings.EqualFold(upgrade, "websocket")
}

func contextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, assetServerLoggerKey, logger)
}

func logInfo(ctx context.Context, message string, args ...interface{}) {
	if logger, _ := ctx.Value(assetServerLoggerKey).(*slog.Logger); logger != nil {
		logger.Info(message, args...)
	}
}

func logError(ctx context.Context, message string, args ...interface{}) {
	if logger, _ := ctx.Value(assetServerLoggerKey).(*slog.Logger); logger != nil {
		logger.Error(message, args...)
	}
}
