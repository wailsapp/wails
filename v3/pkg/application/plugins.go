package application

import (
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v3/internal/assetserver"
)

type PluginAPI interface {
}

type Plugin interface {
	Name() string
	Init(api PluginAPI) error
	Shutdown() error
	CallableByJS() []string
	InjectJS() string
}

type PluginManager struct {
	plugins            map[string]Plugin
	assetServer        *assetserver.AssetServer
	initialisedPlugins []Plugin
}

func NewPluginManager(plugins map[string]Plugin, assetServer *assetserver.AssetServer) *PluginManager {
	result := &PluginManager{
		plugins:     plugins,
		assetServer: assetServer,
	}
	return result
}

func (p *PluginManager) Init() []error {

	api := newPluginAPI()
	for _, plugin := range p.plugins {
		err := plugin.Init(api)
		if err != nil {
			globalApplication.error("Plugin failed to initialise:", "plugin", plugin.Name(), "error", err.Error())
			return p.Shutdown()
		}
		p.initialisedPlugins = append(p.initialisedPlugins, plugin)
		injectJS := plugin.InjectJS()
		if injectJS != "" {
			p.assetServer.AddPluginScript(plugin.Name(), injectJS)
		}
		globalApplication.debug("Plugin initialised: " + plugin.Name())
	}
	return nil
}

func (p *PluginManager) Shutdown() []error {
	var errs []error
	for _, plugin := range p.initialisedPlugins {
		err := plugin.Shutdown()
		globalApplication.debug("Plugin shutdown: " + plugin.Name())
		if err != nil {
			err = errors.Wrap(err, "Plugin failed to shutdown: "+plugin.Name())
			errs = append(errs, err)
		}
	}
	return errs
}
