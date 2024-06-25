package application

import (
	"reflect"

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

var pluginReflector = reflect.TypeOf((*Plugin)(nil)).Elem()

func (p *PluginManager) Init() []error {
	globalApplication.info("Initialising plugins", "count", len(p.plugins))
	for _, plugin := range p.plugins {
		globalApplication.info("Initialising plugin: " + plugin.Name())
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
		if reflect.TypeOf(service.instance).Implements(pluginReflector) {
			found := service.instance.(Plugin)
			p.AddPlugin(found)
			globalApplication.info("Plugin found: " + found.Name())
		}

	}
}

func (p *PluginManager) AddPlugin(plugin Plugin) {
	p.plugins[plugin.Name()] = plugin
}
