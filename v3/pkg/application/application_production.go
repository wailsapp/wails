//go:build production

package application

func newApplication(options *Options) *App {
	result := &App{
		isDebugMode: false,
		options:     options.getOptions(false),
	}
	result.init()
	return result
}
