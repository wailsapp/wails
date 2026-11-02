//go:build darwin && !ios && !server

package application

func addDialogCallback(callback dialogResultCallback) int {
	mutex.Lock()
	defer mutex.Unlock()

	// Find the first free integer key
	var id int
	for {
		if _, exists := callbacks[id]; !exists {
			break
		}
		id++
	}

	// Save the function in the map using the integer key
	callbacks[id] = callback

	// Return the key
	return id
}

func newDialogImpl(d *MessageDialog) *macosDialog {
	return &macosDialog{
		dialog: d,
	}
}

func newOpenFileDialogImpl(d *OpenFileDialogStruct) *macosOpenFileDialog {
	return &macosOpenFileDialog{
		dialog: d,
	}
}

func newSaveFileDialogImpl(d *SaveFileDialogStruct) *macosSaveFileDialog {
	return &macosSaveFileDialog{
		dialog: d,
	}
}

func removeDialogCallback(id int) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(callbacks, id)
}

// handleDialogCallback delivers the button pressed in a native dialog to its callback.
func handleDialogCallback(id int, buttonPressed int) {
	mutex.Lock()
	callback, exists := callbacks[id]
	mutex.Unlock()

	if !exists {
		return
	}

	// Call the function with the button number
	callback(buttonPressed) // Replace nil with the actual slice of buttons
}
