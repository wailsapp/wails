package assetserver

import (
	"net/http"
	"path/filepath"
	"sync"
)

var (
	// mimeCache uses sync.Map for better concurrent read performance
	// since reads are far more common than writes
	mimeCache sync.Map

	// mimeTypesByExt maps file extensions to MIME types for common web formats.
	// This approach is preferred over content-based detection because:
	// 1. Extension-based lookup is O(1) vs O(n) content scanning
	// 2. Web assets typically have correct extensions
	// 3. stdlib's http.DetectContentType handles remaining cases adequately
	// 4. Saves ~208KB binary size by not using github.com/wailsapp/mimetype
	mimeTypesByExt = map[string]string{
		// HTML
		".htm":  "text/html; charset=utf-8",
		".html": "text/html; charset=utf-8",

		// CSS/JS
		".css": "text/css; charset=utf-8",
		".js":  "text/javascript; charset=utf-8",
		".mjs": "text/javascript; charset=utf-8",
		".ts":  "text/typescript; charset=utf-8",
		".tsx": "text/typescript; charset=utf-8",
		".jsx": "text/javascript; charset=utf-8",

		// Data formats
		".json": "application/json",
		".xml":  "text/xml; charset=utf-8",
		".yaml": "text/yaml; charset=utf-8",
		".yml":  "text/yaml; charset=utf-8",
		".toml": "text/toml; charset=utf-8",

		// Images
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".webp": "image/webp",
		".avif": "image/avif",
		".svg":  "image/svg+xml",
		".ico":  "image/x-icon",
		".bmp":  "image/bmp",
		".tiff": "image/tiff",
		".tif":  "image/tiff",

		// Fonts
		".woff":  "font/woff",
		".woff2": "font/woff2",
		".ttf":   "font/ttf",
		".otf":   "font/otf",
		".eot":   "application/vnd.ms-fontobject",

		// Audio
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".ogg":  "audio/ogg",
		".m4a":  "audio/mp4",
		".aac":  "audio/aac",
		".flac": "audio/flac",
		".opus": "audio/opus",

		// Video
		".mp4":  "video/mp4",
		".webm": "video/webm",
		".ogv":  "video/ogg",
		".mov":  "video/quicktime",
		".avi":  "video/x-msvideo",
		".mkv":  "video/x-matroska",
		".m4v":  "video/mp4",

		// Documents
		".pdf": "application/pdf",
		".txt": "text/plain; charset=utf-8",
		".md":  "text/markdown; charset=utf-8",

		// Archives
		".zip": "application/zip",
		".gz":  "application/gzip",
		".tar": "application/x-tar",

		// WebAssembly
		".wasm": "application/wasm",

		// Source maps
		".map": "application/json",
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

	// Slow path: use stdlib content-based detection and cache
	result := http.DetectContentType(data)
	if result == "" {
		result = "application/octet-stream"
	}

	mimeCache.Store(filename, result)
	return result
}
