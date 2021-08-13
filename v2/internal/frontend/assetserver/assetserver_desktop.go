// +build desktop

package assetserver

import (
	"embed"
	"net/http"
)

func (a *AssetServer) init(assets embed.FS) error {

	var err error
	a.assets, err = processAssets(assets)
	if err != nil {
		return err
	}
	indexHTML, err := a.assets.ReadFile("index.html")
	if err != nil {
		return err
	}
	a.indexFile, err = injectScript(string(indexHTML), "<script>"+a.runtimeJS+"</script>")
	if err != nil {
		return err
	}
	return nil
}

func (a *AssetServer) Load(filename string) ([]byte, string, error) {
	var content []byte
	var err error
	switch filename {
	case "/":
		content = a.indexFile
	default:
		content, err = a.assets.ReadFile(filename)
	}
	if err != nil {
		return nil, "", err
	}
	mimeType := http.DetectContentType(content)
	return content, mimeType, nil
}
