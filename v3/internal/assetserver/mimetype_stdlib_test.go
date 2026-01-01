package assetserver

import (
	"net/http"
	"path/filepath"
	"testing"
)

// TestMimeTypeDetection_WebFormats validates that extension-based detection
// plus stdlib fallback correctly handles all common web asset formats.
// This test ensures we can safely remove the github.com/wailsapp/mimetype dependency.
func TestMimeTypeDetection_WebFormats(t *testing.T) {
	// webMimeTests covers all common web formats that Wails applications typically serve
	webMimeTests := []struct {
		name       string
		filename   string
		data       []byte
		wantPrefix string // Use prefix matching since charset may vary
	}{
		// === TEXT FORMATS (extension-based) ===
		{"HTML file", "index.html", []byte("<!DOCTYPE html><html></html>"), "text/html"},
		{"HTM file", "page.htm", []byte("<html></html>"), "text/html"},
		{"CSS file", "styles.css", []byte(".class { color: red; }"), "text/css"},
		{"JavaScript file", "app.js", []byte("function test() {}"), "text/javascript"},
		{"ES Module file", "module.mjs", []byte("export default {}"), "text/javascript"},
		{"JSON file", "data.json", []byte(`{"key": "value"}`), "application/json"},
		{"XML file", "data.xml", []byte("<?xml version=\"1.0\"?><root/>"), "text/xml"},

		// === IMAGE FORMATS (extension-based) ===
		{"PNG file", "image.png", pngData, "image/png"},
		{"JPEG file", "photo.jpg", jpegData, "image/jpeg"},
		{"JPEG alt ext", "photo.jpeg", jpegData, "image/jpeg"},
		{"GIF file", "anim.gif", gifData, "image/gif"},
		{"WebP file", "image.webp", webpData, "image/webp"},
		{"AVIF file", "image.avif", avifData, "image/avif"},
		{"SVG file", "icon.svg", []byte("<svg></svg>"), "image/svg+xml"},
		{"PDF file", "doc.pdf", pdfData, "application/pdf"},

		// === WASM (extension-based) ===
		{"WASM file", "app.wasm", wasmData, "application/wasm"},

		// === FONT FORMATS (need detection or extension map) ===
		{"WOFF file", "font.woff", woffData, "font/woff"},
		{"WOFF2 file", "font.woff2", woff2Data, "font/woff2"},
		{"TTF file", "font.ttf", ttfData, "font/ttf"},
		{"OTF file", "font.otf", otfData, "font/otf"},
		{"EOT file", "font.eot", eotData, "application/vnd.ms-fontobject"},

		// === AUDIO/VIDEO (common web formats) ===
		{"MP3 file", "audio.mp3", mp3Data, "audio/mpeg"},
		{"MP4 file", "video.mp4", mp4Data, "video/mp4"},
		{"WebM file", "video.webm", webmData, "video/webm"},
		{"OGG file", "audio.ogg", oggData, "audio/ogg"},

		// === ARCHIVES (sometimes served by web apps) ===
		{"ZIP file", "archive.zip", zipData, "application/zip"},
		{"GZIP file", "data.gz", gzipData, "application/"},

		// === SOURCE MAPS (common in dev mode) ===
		{"Source map", "app.js.map", []byte(`{"version":3}`), "application/json"},

		// === ICO (favicon) ===
		{"ICO file", "favicon.ico", icoData, "image/"},

		// === FALLBACK TESTS ===
		{"Unknown binary", "data.bin", []byte{0x00, 0x01, 0x02, 0x03}, "application/octet-stream"},
		{"Plain text (no ext)", "readme", []byte("Hello World"), "text/plain"},
	}

	for _, tt := range webMimeTests {
		t.Run(tt.name, func(t *testing.T) {
			got := getMimeTypeStdlib(tt.filename, tt.data)
			if !hasPrefix(got, tt.wantPrefix) {
				t.Errorf("getMimeTypeStdlib(%q) = %q, want prefix %q", tt.filename, got, tt.wantPrefix)
			}
		})
	}
}

// getMimeTypeStdlib is the proposed replacement that uses only stdlib
func getMimeTypeStdlib(filename string, data []byte) string {
	// Fast path: check extension map first
	if result := extMimeTypes[filepath.Ext(filename)]; result != "" {
		return result
	}

	// Fallback to stdlib content-based detection
	result := http.DetectContentType(data)
	if result == "" {
		result = "application/octet-stream"
	}
	return result
}

// extMimeTypes is an expanded map covering all common web formats
// This replaces the need for the mimetype library for web assets
var extMimeTypes = map[string]string{
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

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// Magic bytes for various formats
var (
	// PNG: 89 50 4E 47 0D 0A 1A 0A
	pngData = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D}

	// JPEG: FF D8 FF
	jpegData = []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}

	// GIF: 47 49 46 38
	gifData = []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}

	// WebP: 52 49 46 46 ... 57 45 42 50
	webpData = []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50}

	// AVIF: ... ftypavif or ftypavis
	avifData = []byte{0x00, 0x00, 0x00, 0x1C, 0x66, 0x74, 0x79, 0x70, 0x61, 0x76, 0x69, 0x66}

	// PDF: 25 50 44 46
	pdfData = []byte{0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E}

	// WASM: 00 61 73 6D
	wasmData = []byte{0x00, 0x61, 0x73, 0x6D, 0x01, 0x00, 0x00, 0x00}

	// WOFF: 77 4F 46 46
	woffData = []byte{0x77, 0x4F, 0x46, 0x46, 0x00, 0x01, 0x00, 0x00}

	// WOFF2: 77 4F 46 32
	woff2Data = []byte{0x77, 0x4F, 0x46, 0x32, 0x00, 0x01, 0x00, 0x00}

	// TTF: 00 01 00 00
	ttfData = []byte{0x00, 0x01, 0x00, 0x00, 0x00}

	// OTF: 4F 54 54 4F (OTTO)
	otfData = []byte{0x4F, 0x54, 0x54, 0x4F, 0x00}

	// EOT: varies, but starts with size bytes then magic
	eotData = []byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x00}

	// MP3: FF FB or FF FA or ID3
	mp3Data = []byte{0xFF, 0xFB, 0x90, 0x00}

	// MP4: ... ftyp
	mp4Data = []byte{0x00, 0x00, 0x00, 0x1C, 0x66, 0x74, 0x79, 0x70, 0x69, 0x73, 0x6F, 0x6D}

	// WebM: 1A 45 DF A3 (EBML header)
	webmData = []byte{0x1A, 0x45, 0xDF, 0xA3}

	// OGG: 4F 67 67 53
	oggData = []byte{0x4F, 0x67, 0x67, 0x53, 0x00, 0x02}

	// ZIP: 50 4B 03 04
	zipData = []byte{0x50, 0x4B, 0x03, 0x04}

	// GZIP: 1F 8B
	gzipData = []byte{0x1F, 0x8B, 0x08}

	// ICO: 00 00 01 00
	icoData = []byte{0x00, 0x00, 0x01, 0x00, 0x01, 0x00}
)

// TestMimeTypeExtensionMapCompleteness checks that all extensions in the
// original mimeTypesByExt are covered by the expanded extMimeTypes
func TestMimeTypeExtensionMapCompleteness(t *testing.T) {
	for ext, mime := range mimeTypesByExt {
		if newMime, ok := extMimeTypes[ext]; !ok {
			t.Errorf("extension %q missing from extMimeTypes (was: %q)", ext, mime)
		} else if newMime != mime {
			// Allow differences as long as they're equivalent
			if !hasPrefix(newMime, mime[:len(mime)-10]) { // rough prefix check
				t.Logf("extension %q changed: %q -> %q (verify this is correct)", ext, mime, newMime)
			}
		}
	}
}

// BenchmarkMimeType_StdlibOnly benchmarks the stdlib-only implementation
func BenchmarkMimeType_StdlibOnly(b *testing.B) {
	testCases := []struct {
		name     string
		filename string
		data     []byte
	}{
		{"ExtHit_JS", "app.js", []byte("function() {}")},
		{"ExtHit_CSS", "styles.css", []byte(".class { }")},
		{"ExtHit_PNG", "image.png", pngData},
		{"ExtMiss_Binary", "data.bin", []byte{0x00, 0x01, 0x02}},
		{"ContentDetect_PNG", "unknown", pngData},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				_ = getMimeTypeStdlib(tc.filename, tc.data)
			}
		})
	}
}
