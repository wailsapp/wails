//go:build production

package application

func newApplication(options Options) *App {
	result := &App{
		isDebugMode: false,
		options:     options,
	}
	result.init()
	return result
}

func (a *App) logStartup() {}

func (a *App) preRun() error { return nil }

func (a *App) postQuit() error { return nil }
