package application

import (
	"github.com/pkg/errors"
)

type PluginAPI interface {
}

type Plugin interface {
	Name() string
	Init() error
	Shutdown() error
}

type PluginManager struct {
	plugins            map[string]Plugin
	initialisedPlugins []Plugin
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
	}
}

func (p *PluginManager) Init() []error {
	for _, plugin := range p.plugins {
		globalApplication.debug("Initialising plugin: " + plugin.Name())
		err := plugin.Init()
		if err != nil {
			globalApplication.error("Plugin failed to initialise:", "plugin", plugin.Name(), "error", err.Error())
			return p.Shutdown()
		}
		p.initialisedPlugins = append(p.initialisedPlugins, plugin)
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

func (p *PluginManager) ProcessPlugins(services []Service) {
	for _, service := range services {
		if found, ok := service.instance.(Plugin); ok {
			p.AddPlugin(found)
		}
	}
}

func (p *PluginManager) AddPlugin(plugin Plugin) {
	p.plugins[plugin.Name()] = plugin
}
