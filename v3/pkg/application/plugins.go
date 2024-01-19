package application

import (
	"github.com/wailsapp/wails/v3/internal/assetserver"
)

type Plugin interface {
	Name() string
	Init() error
	Shutdown()
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

func (p *PluginManager) Init() error {
	for _, plugin := range p.plugins {
		err := plugin.Init()
		if err != nil {
			globalApplication.error("Plugin failed to initialise:", "plugin", plugin.Name(), "error", err.Error())
			p.Shutdown()
			return err
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

func (p *PluginManager) Shutdown() {
	for _, plugin := range p.initialisedPlugins {
		plugin.Shutdown()
		globalApplication.debug("Plugin shutdown: " + plugin.Name())
	}
}
