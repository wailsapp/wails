package assetserver

import (
	"bytes"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/gabriel-vasile/mimetype"
)

var (
	utf8BOM = []byte{0xEF, 0xBB, 0xBF}

	mimeCache = map[string]string{}
	mimeMutex sync.Mutex

	// The list of builtin mime-types by extension as defined by
	// the golang standard lib package "mime"
	// The standard lib also takes into account mime type definitions from
	// etc files like '/etc/apache2/mime.types' but we want to have the
	// same behavivour on all platforms and not depend on some external file.
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
	mimeMutex.Lock()
	defer mimeMutex.Unlock()

	result := mimeTypesByExt[filepath.Ext(filename)]
	if result != "" {
		return result
	}

	result = mimeCache[filename]
	if result != "" {
		return result
	}

	// Strip a leading UTF-8 BOM before content detection. The mimetype
	// detectors do not skip it for every format (e.g. SVG), and the charset
	// of BOM-prefixed text is still reported as UTF-8 without it.
	detect := mimetype.Detect(bytes.TrimPrefix(data, utf8BOM))
	if detect == nil {
		result = http.DetectContentType(data)
	} else {
		result = detect.String()
	}

	if result == "" {
		result = "application/octet-stream"
	}

	mimeCache[filename] = result
	return result
}
