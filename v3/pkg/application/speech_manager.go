package application

// SpeechManager speaks text and recognises speech using platform services.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type SpeechManager struct {
	app *App
}

func newSpeechManager(app *App) *SpeechManager {
	return &SpeechManager{app: app}
}
