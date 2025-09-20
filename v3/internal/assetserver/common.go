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
		// Debug CSS serving with clear markers
		if strings.HasSuffix(filename, ".css") {
			fmt.Printf("\n🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨\n")
			fmt.Printf("CSS FILE BEING SERVED:\n")
			fmt.Printf("  Filename: %s\n", filename)
			fmt.Printf("  MimeType: %s\n", mimeType)
			fmt.Printf("  Size: %d bytes\n", len(blob))
			if len(blob) > 0 {
				preview := string(blob)
				if len(preview) > 100 {
					preview = preview[:100] + "..."
				}
				fmt.Printf("  Preview: %s\n", preview)
			}
			fmt.Printf("🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨🎨\n\n")
		}
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
