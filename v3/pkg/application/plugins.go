package application

type Plugin interface {
	Name() string
	Init(app *App) error
	Shutdown()
	// Exported is a list of method names that should be exposed to the frontend
	Exported() []string
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

func (p *PluginManager) Shutdown() {
	for _, plugin := range p.plugins {
		plugin.Shutdown()
		globalApplication.info("Plugin '%s' shutdown", plugin.Name())
	}
}
