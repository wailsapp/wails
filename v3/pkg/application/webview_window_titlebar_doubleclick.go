package application

func (w *WebviewWindow) handleTitlebarDoubleClick() {
	if w.IsFullscreen() {
		return
	}

	switch parseMacTitlebarDoubleClickAction(platformTitlebarDoubleClickPreference()) {
	case titlebarDoubleClickToggleMaximise:
		w.ToggleMaximise()
	case titlebarDoubleClickMinimise:
		w.Minimise()
	}
}
