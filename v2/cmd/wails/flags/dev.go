package flags

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/internal/project"
	"net"
	"net/url"
	"os"
	"path/filepath"
)

type Dev struct {
	BuildCommon

	AssetDir             string `flag:"assetdir" description:"Serve assets from the given directory instead of using the provided asset FS"`
	Extensions           string `flag:"e" description:"Extensions to trigger rebuilds (comma separated) eg go"`
	ReloadDirs           string `flag:"reloaddirs" description:"Additional directories to trigger reloads (comma separated)"`
	Browser              bool   `flag:"browser" description:"Open the application in a browser"`
	NoReload             bool   `flag:"noreload" description:"Disable reload on asset change"`
	NoColour             bool   `flag:"nocolor" description:"Disable colour in output"`
	WailsJSDir           string `flag:"wailsjsdir" description:"Directory to generate the Wails JS modules"`
	LogLevel             string `flag:"loglevel" description:"LogLevel to use - Trace, Debug, Info, Warning, Error)"`
	ForceBuild           bool   `flag:"f" description:"Force build of application"`
	Debounce             int    `flag:"debounce" description:"The amount of time to wait to trigger a reload on change"`
	DevServer            string `flag:"devserver" description:"The address of the wails dev server"`
	AppArgs              string `flag:"appargs" description:"arguments to pass to the underlying app (quoted and space separated)"`
	Save                 bool   `flag:"save" description:"Save the given flags as defaults"`
	FrontendDevServerURL string `flag:"frontenddevserverurl" description:"The url of the external frontend dev server to use"`

	// Internal state
	devServerURL *url.URL
}

func Default() *Dev {
	result := &Dev{
		Extensions: "go",
		Debounce:   100,
	}
	result.BuildCommon = result.BuildCommon.Default()
	return result
}

func (d *Dev) Process() error {

	var err error
	err = d.loadAndMergeProjectConfig()
	if err != nil {
		return err
	}

	if _, _, err := net.SplitHostPort(d.DevServer); err != nil {
		return fmt.Errorf("DevServer is not of the form 'host:port', please check your wails.json")
	}

	d.devServerURL, err = url.Parse("http://" + d.DevServer)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dev) loadAndMergeProjectConfig() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	projectConfig, err := project.Load(cwd)
	if err != nil {
		return err
	}

	d.AssetDir, _ = lo.Coalesce(d.AssetDir, projectConfig.AssetDirectory)
	projectConfig.AssetDirectory = filepath.ToSlash(d.AssetDir)
	if d.AssetDir != "" {
		d.AssetDir, err = filepath.Abs(d.AssetDir)
		if err != nil {
			return err
		}
	}

	d.ReloadDirs, _ = lo.Coalesce(d.ReloadDirs, projectConfig.ReloadDirectories)
	projectConfig.ReloadDirectories = filepath.ToSlash(d.ReloadDirs)
	d.DevServer, _ = lo.Coalesce(d.DevServer, projectConfig.DevServer)
	projectConfig.DevServer = d.DevServer
	d.FrontendDevServerURL, _ = lo.Coalesce(d.FrontendDevServerURL, projectConfig.FrontendDevServerURL)
	projectConfig.FrontendDevServerURL = d.FrontendDevServerURL
	d.WailsJSDir, _ = lo.Coalesce(d.WailsJSDir, projectConfig.GetWailsJSDir(), projectConfig.GetFrontendDir())
	projectConfig.WailsJSDir = filepath.ToSlash(d.WailsJSDir)

	if d.Debounce == 100 && projectConfig.DebounceMS != 100 {
		if projectConfig.DebounceMS == 0 {
			projectConfig.DebounceMS = 100
		}
		d.Debounce = projectConfig.DebounceMS
	}
	projectConfig.DebounceMS = d.Debounce

	d.AppArgs, _ = lo.Coalesce(d.AppArgs, projectConfig.AppArgs)

	if d.Save {
		err = projectConfig.Save()
		if err != nil {
			return err
		}
	}

	return nil

}
