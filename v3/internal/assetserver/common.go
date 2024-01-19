package assetserver

import (
	"bytes"
	"fmt"
	"io"
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

func serveFile(rw http.ResponseWriter, filename string, blob []byte) error {
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
