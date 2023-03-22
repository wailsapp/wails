package application

type Plugin interface {
	Name() string
	Init(app *App) error
	Shutdown()
	CallableByJS() []string
	InjectJS() string
}

type PluginManager struct {
	plugins map[string]Plugin
}

func NewPluginManager(plugins map[string]Plugin) *PluginManager {
	result := &PluginManager{
		plugins: plugins,
	}
	globalApplication.OnWindowCreation(result.onWindowCreation)
	return result
}

func (p *PluginManager) Init() error {
	for _, plugin := range p.plugins {
		err := plugin.Init(globalApplication)
		if err != nil {
			return err
		}
		injectJS := plugin.InjectJS()
		if injectJS != "" {

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

func (p *PluginManager) onWindowCreation(window *WebviewWindow) {
	for _, plugin := range p.plugins {
		injectJS := plugin.InjectJS()
		if injectJS != "" {
			window.ExecJS(injectJS)
		}
	}
}
