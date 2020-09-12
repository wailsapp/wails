package html

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"unsafe"
)

type assetTypes struct {
	JS      string
	CSS     string
	FAVICON string
	HTML    string
}

// AssetTypes is an enum for the asset type keys
var AssetTypes *assetTypes = &assetTypes{
	JS:      "javascript",
	CSS:     "css",
	FAVICON: "favicon",
	HTML:    "html",
}

// Asset describes an asset type and its path
type Asset struct {
	Type string
	Path string
	Data string
}

// Load the asset from disk
func (a *Asset) Load(basedirectory string) error {
	assetpath := filepath.Join(basedirectory, a.Path)
	data, err := ioutil.ReadFile(assetpath)
	if err != nil {
		return err
	}
	a.Data = string(data)
	return nil
}

// AsString returns the data as a READ ONLY string
func (a *Asset) AsString() string {
	return a.Data
}

// AsCHexData processes the asset data so it may be used by C
func (a *Asset) AsCHexData() string {

	// This will be our final string to hexify
	dataString := a.Data

	switch a.Type {
	case AssetTypes.HTML:

		// Escape HTML
		var re = regexp.MustCompile(`\s{2,}`)
		result := re.ReplaceAllString(a.Data, ``)
		result = strings.ReplaceAll(result, "\n", "")
		result = strings.ReplaceAll(result, "\r\n", "")
		result = strings.ReplaceAll(result, "\n", "")

		// Inject wailsloader code
		result = strings.Replace(result, `</body>`, `<script id='wailsloader'>window.wailsloader = { html: true, runtime: false, userjs: false, usercss: false };var self=document.querySelector('#wailsloader');self.parentNode.removeChild(self);</script></body>`, 1)

		url := url.URL{Path: result}
		urlString := strings.ReplaceAll(url.String(), "/", "%2f")

		// Save Data uRI string
		dataString = "data:text/html;charset=utf-8," + urlString

	case AssetTypes.CSS:

		// Escape CSS data
		var re = regexp.MustCompile(`\s{2,}`)
		result := re.ReplaceAllString(a.Data, ``)
		result = strings.ReplaceAll(result, "\n", "")
		result = strings.ReplaceAll(result, "\r\n", "")
		result = strings.ReplaceAll(result, "\n", "")
		result = strings.ReplaceAll(result, "\t", "")
		result = strings.ReplaceAll(result, `\`, `\\`)
		result = strings.ReplaceAll(result, `"`, `\"`)
		result = strings.ReplaceAll(result, `'`, `\'`)
		result = strings.ReplaceAll(result, ` {`, `{`)
		result = strings.ReplaceAll(result, `: `, `:`)
		dataString = fmt.Sprintf("window.wails._.InjectCSS(\"%s\");", result)
	}

	// Get byte data of the string
	bytes := *(*[]byte)(unsafe.Pointer(&dataString))

	// Create a strings builder
	var cdata strings.Builder

	// Set buffer size to 4k
	cdata.Grow(4096)

	// Convert each byte to hex
	for _, b := range bytes {
		cdata.WriteString(fmt.Sprintf("0x%x, ", b))
	}

	return cdata.String()
}

func (a *Asset) Dump() {
	fmt.Printf("{ Type: %s, Path: %s, Data: %+v }\n", a.Type, a.Path, a.Data[:10])
}
