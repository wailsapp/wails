package main

import (
	"fmt"

	"github.com/wailsapp/wails/v2"
)

type MyStruct struct {
	runtime *wails.Runtime
}

type Notes struct {
	lines []string
}

func (m *MyStruct) WailsInit(runtime *wails.Runtime) error {

	runtime.Events.On("notes updated", func(optionalData ...interface{}) {
		// Get notes
		notes := optionalData[0].(*Notes)
		// Save the notes to disk
		fmt.Printf("%+v\n", notes)
	})

	m.runtime = runtime
	return nil
}
