package build

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/html"
)

// ServerBuilder builds applications as a server
type ServerBuilder struct {
	*BaseBuilder
}

func newServerBuilder() *ServerBuilder {
	result := &ServerBuilder{
		BaseBuilder: NewBaseBuilder(),
	}
	return result
}

// BuildAssets builds the assets for the desktop application
func (s *ServerBuilder) BuildAssets(options *Options) error {
	var err error

	assets, err := s.BaseBuilder.ExtractAssets()
	if err != nil {
		return err
	}

	// Build embedded assets (HTML/JS/CSS/etc)
	err = s.BuildBaseAssets(assets)
	if err != nil {
		return err
	}

	// Build static assets
	err = s.buildStaticAssets(s.projectData)
	if err != nil {
		return err
	}

	return nil
}

// BuildBaseAssets builds the base assets
func (s *ServerBuilder) BuildBaseAssets(assets *html.AssetBundle) error {
	db, err := assets.ConvertToAssetDB()
	if err != nil {
		return err
	}

	// Fetch, update, and reinject index.html
	index, err := db.Read("index.html")
	if err != nil {
		return fmt.Errorf(`Failed to locate "index.html"`)
	}
	splits := strings.Split(string(index), "</body>")
	if len(splits) != 2 {
		return fmt.Errorf("Unable to locate a </body> tag in your frontend/index.html")
	}
	injectScript := `<script defer src="/wails.js"></script><script defer src="/bindings.js"></script>`
	result := []string{}
	result = append(result, splits[0])
	result = append(result, injectScript)
	result = append(result, "</body>")
	result = append(result, splits[1])

	db.Remove("index.html") // Remove the non-prefixed index.html
	db.AddAsset("/index.html", []byte(strings.Join(result, "")))

	// Add wails.js
	wailsjsPath := fs.RelativePath("../../../internal/runtime/assets/server.js")
	if rawData, err := ioutil.ReadFile(wailsjsPath); err == nil {
		db.AddAsset("/wails.js", rawData)
	}

	targetFile := fs.RelativePath("../../../internal/webserver/webassetsdb.go")
	s.addFileToDelete(targetFile)
	f, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(db.Serialize("db", "webserver"))

	return nil
}

// BuildRuntime builds the javascript runtime used by the HTML client to connect to the websocket
func (s *ServerBuilder) BuildRuntime(options *Options) error {

	sourceDir := fs.RelativePath("../../../internal/runtime/js")
	if err := s.NpmInstall(sourceDir); err != nil {
		return err
	}

	options.Logger.Print("    - Embedding Runtime...")
	envvars := []string{"WAILSPLATFORM=" + options.Platform}
	var err error
	if err = s.NpmRunWithEnvironment(sourceDir, "build:server", false, envvars); err != nil {
		return err
	}

	options.Logger.Println("done.")
	return nil
}
