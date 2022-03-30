package assetserver

import (
	"bufio"
	"bytes"
	"net/http"
	"path/filepath"
	"sync"
	"unicode"

	"github.com/gabriel-vasile/mimetype"
)

var (
	cache           = map[string]string{}
	mutex           sync.Mutex
	bomUTF8         = []byte{0xef, 0xbb, 0xbf}
	startComment    = []byte{0x3c, 0x21, 0x2d, 0x2d}
	lenStartComment = len(startComment)
	endComment      = []byte{0x2d, 0x2d, 0x3e}
	lenComment      = len(endComment)
)

func GetMimetype(filename string, data []byte) string {
	mutex.Lock()
	defer mutex.Unlock()

	// short-circuit .js, .css and check .html and .svg to ensure the
	// browser evaluates them in the right context
	startDetectAt := 0
	switch filepath.Ext(filename) {
	case ".js":
		return "application/javascript"
	case ".css":
		return "text/css; charset=utf-8"
	case ".html":
		// some HTML files will get classified as text/plain in some circumstances
		// for instance when starting with a BOM.
		if len(data) > 4 && data[0] == bomUTF8[0] && data[1] == bomUTF8[1] && data[2] == bomUTF8[2] {
			startDetectAt = 3
		}
		// we could return early but might as well let the "mimetype.Detect" function take care of it in
		// case there's invalid content in the file
	case ".svg":
		// some SVG files will get classified as text/html in some circumstances
		// for instance when starting with comments such '<!-- my comment -->'
		startDetectAt = getSvgStart(data)
		// we could return early but might as well let the "mimetype.Detect" function take care of it in
		// case there's invalid content in the file
	}

	result := cache[filename]
	if result != "" {
		return result
	}

	detect := mimetype.Detect(data[startDetectAt:])
	if detect == nil {
		result = http.DetectContentType(data[startDetectAt:])
	} else {
		result = detect.String()
	}

	if result == "" {
		result = "application/octet-stream"
	}

	cache[filename] = result
	return result
}

func getSvgStart(data []byte) int {
	lenData := len(data)
	if lenData > lenStartComment+lenComment {
		var reader *bytes.Reader
		pos := 0
		if len(data) > 4 && data[0] == bomUTF8[0] && data[1] == bomUTF8[1] && data[2] == bomUTF8[2] {
			reader = bytes.NewReader(data[4:])
		} else {
			reader = bytes.NewReader(data)
		}
		scanner := bufio.NewScanner(reader)
	EndScan:
		for scanner.Scan() {
			for _, ch := range scanner.Text() {
				if !unicode.IsSpace(ch) {
					break EndScan
				} else {
					pos++
				}
			}
		}

		startComment := bytes.Index(data, startComment)
		if startComment > -1 && startComment >= pos {
			endComment := bytes.Index(data[startComment+1:], endComment)
			if endComment > -1 && endComment+lenComment+1 < lenData {
				return endComment + lenComment + 1
			}
		}
	}

	return 0
}
