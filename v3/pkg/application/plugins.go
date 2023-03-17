package application

import (
	"fmt"
)

type Plugin interface {
	Name() string
	Init(app *App) error
	Call(args []any) (any, error)
}

type PluginManager struct {
	plugins map[string]Plugin
}

func NewPluginManager(plugins map[string]Plugin) *PluginManager {
	return &PluginManager{
		plugins: plugins,
	}
}

func (p *PluginManager) Init() error {
	for _, plugin := range p.plugins {
		err := plugin.Init(globalApplication)
		if err != nil {
			return err
		}
		globalApplication.info("Plugin '%s' initialised", plugin.Name())
	}
	return nil
}

func (p *PluginManager) Call(name string, args []any) (any, error) {
	plugin, ok := p.plugins[name]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", name)
	}
	return plugin.Call(args)
}
