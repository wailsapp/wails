package flags

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/internal/project"
	"github.com/wailsapp/wails/v2/pkg/commands/build"
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
	devServerURL  *url.URL
	projectConfig *project.Project
}

func (*Dev) Default() *Dev {
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
	var err error
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	d.projectConfig, err = project.Load(cwd)
	if err != nil {
		return err
	}

	d.AssetDir, _ = lo.Coalesce(d.AssetDir, d.projectConfig.AssetDirectory)
	d.projectConfig.AssetDirectory = filepath.ToSlash(d.AssetDir)
	if d.AssetDir != "" {
		d.AssetDir, err = filepath.Abs(d.AssetDir)
		if err != nil {
			return err
		}
	}

	d.ReloadDirs, _ = lo.Coalesce(d.ReloadDirs, d.projectConfig.ReloadDirectories)
	d.projectConfig.ReloadDirectories = filepath.ToSlash(d.ReloadDirs)
	d.DevServer, _ = lo.Coalesce(d.DevServer, d.projectConfig.DevServer)
	d.projectConfig.DevServer = d.DevServer
	d.FrontendDevServerURL, _ = lo.Coalesce(d.FrontendDevServerURL, d.projectConfig.FrontendDevServerURL)
	d.projectConfig.FrontendDevServerURL = d.FrontendDevServerURL
	d.WailsJSDir, _ = lo.Coalesce(d.WailsJSDir, d.projectConfig.GetWailsJSDir(), d.projectConfig.GetFrontendDir())
	d.projectConfig.WailsJSDir = filepath.ToSlash(d.WailsJSDir)

	if d.Debounce == 100 && d.projectConfig.DebounceMS != 100 {
		if d.projectConfig.DebounceMS == 0 {
			d.projectConfig.DebounceMS = 100
		}
		d.Debounce = d.projectConfig.DebounceMS
	}
	d.projectConfig.DebounceMS = d.Debounce

	d.AppArgs, _ = lo.Coalesce(d.AppArgs, d.projectConfig.AppArgs)

	if d.Save {
		err = d.projectConfig.Save()
		if err != nil {
			return err
		}
	}

	return nil

}

// GenerateBuildOptions creates a build.Options using the flags
func (d *Dev) GenerateBuildOptions() *build.Options {
	result := &build.Options{
		OutputType:     "dev",
		Mode:           build.Dev,
		Arch:           runtime.GOARCH,
		Pack:           true,
		Platform:       runtime.GOOS,
		LDFlags:        d.LdFlags,
		Compiler:       d.Compiler,
		ForceBuild:     d.ForceBuild,
		IgnoreFrontend: d.SkipFrontend,
		SkipBindings:   d.SkipBindings,
		Verbosity:      d.Verbosity,
		WailsJSDir:     d.WailsJSDir,
		RaceDetector:   d.RaceDetector,
		ProjectData:    d.projectConfig,
	}

	return result
}

func (d *Dev) ProjectConfig() *project.Project {
	return d.projectConfig
}

func (d *Dev) DevServerURL() *url.URL {
	return d.devServerURL
}
