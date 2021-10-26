package assetserver

import (
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gabriel-vasile/mimetype"
)

var (
	cache = map[string]string{}
	mutex sync.Mutex
)

func GetMimetype(filename string, data []byte) string {
	mutex.Lock()
	defer mutex.Unlock()

	result := cache[filename]
	if result != "" {
		return result
	}

	detect := mimetype.Detect(data)
	if detect == nil {
		result = http.DetectContentType(data)
	} else {
		result = detect.String()
	}

	if filepath.Ext(filename) == ".css" && strings.HasPrefix(result, "text/plain") {
		result = strings.Replace(result, "text/plain", "text/css", 1)
	}

	if filepath.Ext(filename) == ".js" && strings.HasPrefix(result, "text/plain") {
		result = strings.Replace(result, "text/plain", "text/javascript", 1)
	}

	if result == "" {
		result = "application/octet-stream"
	}

	cache[filename] = result
	return result
}
