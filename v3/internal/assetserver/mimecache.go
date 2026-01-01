package assetserver

import (
	"net/http"
	"path/filepath"
	"sync"

	"github.com/wailsapp/mimetype"
)

var (
	// mimeCache uses sync.Map for better concurrent read performance
	// since reads are far more common than writes
	mimeCache sync.Map

	// The list of builtin mime-types by extension as defined by
	// the golang standard lib package "mime"
	// The standard lib also takes into account mime type definitions from
	// /etc files like '/etc/apache2/mime.types' but we want to have the
	// same behaviour on all platforms and not depend on some external file.
	mimeTypesByExt = map[string]string{
		".avif": "image/avif",
		".css":  "text/css; charset=utf-8",
		".gif":  "image/gif",
		".htm":  "text/html; charset=utf-8",
		".html": "text/html; charset=utf-8",
		".jpeg": "image/jpeg",
		".jpg":  "image/jpeg",
		".js":   "text/javascript; charset=utf-8",
		".json": "application/json",
		".mjs":  "text/javascript; charset=utf-8",
		".pdf":  "application/pdf",
		".png":  "image/png",
		".svg":  "image/svg+xml",
		".wasm": "application/wasm",
		".webp": "image/webp",
		".xml":  "text/xml; charset=utf-8",
	}
)

func GetMimetype(filename string, data []byte) string {
	// Fast path: check extension map first (no lock needed)
	if result := mimeTypesByExt[filepath.Ext(filename)]; result != "" {
		return result
	}

	// Check cache (lock-free read)
	if cached, ok := mimeCache.Load(filename); ok {
		return cached.(string)
	}

	// Slow path: detect and cache
	var result string
	detect := mimetype.Detect(data)
	if detect == nil {
		result = http.DetectContentType(data)
	} else {
		result = detect.String()
	}

	if result == "" {
		result = "application/octet-stream"
	}

	mimeCache.Store(filename, result)
	return result
}
