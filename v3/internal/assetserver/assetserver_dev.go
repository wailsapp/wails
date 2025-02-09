//go:build !production

package assetserver

import (
	"embed"
	"io"
	iofs "io/fs"
)

//go:embed defaults
var defaultHTML embed.FS

func defaultIndexHTML(language string) []byte {
	result := []byte("index.html not found")
	// Create an fs.Sub in the defaults directory
	defaults, err := iofs.Sub(defaultHTML, "defaults")
	if err != nil {
		return result
	}
	// Get the 2 character language code
	lang := "en"
	if len(language) >= 2 {
		lang = language[:2]
	}
	// Now we can read the index.html file in the format
	// index.<lang>.html.

	indexFile, err := defaults.Open("index." + lang + ".html")
	if err != nil {
		return result
	}

	indexBytes, err := io.ReadAll(indexFile)
	if err != nil {
		return result
	}
	return indexBytes
}

func (a *AssetServer) LogDetails() {
	var info = []any{
		"middleware", a.options.Middleware != nil,
		"handler", a.options.Handler != nil,
	}
	if devServerURL := GetDevServerURL(); devServerURL != "" {
		info = append(info, "devServerURL", devServerURL)
	}
	a.options.Logger.Info("AssetServer Info:", info...)
}
