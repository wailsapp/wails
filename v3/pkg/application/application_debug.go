//go:build !production

package application

// We use this to patch the application to production mode.
func newApplication(options *Options) *App {
	result := &App{
		isDebugMode: true,
		options:     options.getOptions(true),
	}
	result.init()
	return result
}
