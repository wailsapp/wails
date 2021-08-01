package assetserver

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/leaanthony/debme"
	"github.com/leaanthony/slicer"
	"io/fs"
	"path/filepath"
	"strings"
)

type AssetServer struct {
	assets    debme.Debme
	indexFile []byte
}

func (a *AssetServer) IndexHTML() string {
	return string(a.indexFile)
}

func NewAssetServer(assets embed.FS) (*AssetServer, error) {
	result := &AssetServer{}
	err := result.init(assets)
	return result, err
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
