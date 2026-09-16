package application

// SoundManager plays system sounds and alert tones.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type SoundManager struct {
	app *App
}

func newSoundManager(app *App) *SoundManager {
	return &SoundManager{app: app}
}
