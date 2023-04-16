package application

import "github.com/wailsapp/wails/v2/pkg/assetserver"

type Plugin interface {
	Name() string
	Init(app *App) error
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
		err := plugin.Init(globalApplication)
		if err != nil {
			globalApplication.error("Plugin '%s' failed to initialise: %s", plugin.Name(), err.Error())
			p.Shutdown()
			return err
		}
		p.initialisedPlugins = append(p.initialisedPlugins, plugin)
		injectJS := plugin.InjectJS()
		if injectJS != "" {
			p.assetServer.AddPluginScript(plugin.Name(), injectJS)
		}
		globalApplication.info("Plugin '%s' initialised", plugin.Name())
	}
	return nil
}

func (p *PluginManager) Shutdown() {
	for _, plugin := range p.initialisedPlugins {
		plugin.Shutdown()
		globalApplication.info("Plugin '%s' shutdown", plugin.Name())
	}
}
