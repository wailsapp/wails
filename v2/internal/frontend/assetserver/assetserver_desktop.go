//go:build desktop
// +build desktop

package assetserver

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/leaanthony/debme"
	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

type AssetServer struct {
	assets    debme.Debme
	indexFile []byte
	runtimeJS []byte
}

func NewAssetServer(assets embed.FS, bindingsJSON string) (*AssetServer, error) {
	result := &AssetServer{}
	var buffer bytes.Buffer
	buffer.Write(runtime.IPCJS)
	buffer.WriteString(`window.wailsbindings='` + bindingsJSON + `';` + "\n")
	buffer.Write(runtime.RuntimeJS)
	result.runtimeJS = buffer.Bytes()
	err := result.init(assets)
	return result, err
}

func (a *AssetServer) IndexHTML() string {
	return string(a.indexFile)
}

func injectScript(input string, script string) ([]byte, error) {
	splits := strings.Split(input, "<head>")
	if len(splits) != 2 {
		return nil, fmt.Errorf("unable to locate a </body> tag in your html")
	}

	var result bytes.Buffer
	result.WriteString(splits[0])
	result.WriteString("<head>")
	result.WriteString(script)
	result.WriteString(splits[1])
	return result.Bytes(), nil
}

func processAssets(assets embed.FS) (debme.Debme, error) {

	result, err := debme.FS(assets, ".")
	if err != nil {
		return result, err
	}
	// Find index.html
	stat, err := fs.Stat(assets, "index.html")
	if stat != nil {
		return debme.FS(assets, ".")
	}
	var indexFiles slicer.StringSlicer
	err = fs.WalkDir(result, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, "index.html") {
			indexFiles.Add(path)
		}
		return nil
	})
	if err != nil {
		return debme.Debme{}, err
	}

	if indexFiles.Length() > 1 {
		return debme.Debme{}, fmt.Errorf("multiple 'index.html' files found in assets")
	}

	path, _ := filepath.Split(indexFiles.AsSlice()[0])
	return debme.FS(assets, path)
}

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
	a.indexFile, err = injectScript(string(indexHTML), `<script src="/wails/runtime.js"></script>`)
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
	case "/wails/runtime.js":
		content = a.runtimeJS
	case "/wails/ipc.js":
		content = runtime.IPCJS
	default:
		content, err = a.assets.ReadFile(filename)
	}
	if err != nil {
		return nil, "", err
	}
	mimeType := http.DetectContentType(content)
	return content, mimeType, nil
}
