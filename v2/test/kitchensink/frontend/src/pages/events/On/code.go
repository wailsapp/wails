package main

import wails "github.com/wailsapp/wails/v2"

type MyStruct struct {
	runtime *wails.Runtime
}

func (m *MyStruct) WailsInit(runtime *wails.Runtime) error {

	runtime.Events.On("notes updated", func(optionalData ...interface{}) {
		// Get notes
		notes := optionalData[0].(*Notes)
		// Save the notes to disk
	})

	m.runtime = runtime
	return nil
}
