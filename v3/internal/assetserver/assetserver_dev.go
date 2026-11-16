//go:build !production

package assetserver

import "embed"

//go:embed defaults
var defaultHTML embed.FS

func defaultIndexHTML(language string) []byte {
	// Get the 2 character language code
	lang := "en"
	if len(language) >= 2 {
		lang = language[:2]
	}
	indexBytes, err := defaultHTML.ReadFile("defaults/index." + lang + ".html")
	if err != nil {
		return []byte("index.html not found")
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
