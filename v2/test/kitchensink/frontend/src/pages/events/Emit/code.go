package main

import (
	"io/ioutil"

	"github.com/wailsapp/wails/v2"
)

type MyStruct struct {
	runtime *wails.Runtime
}

func (m *MyStruct) WailsInit(runtime *wails.Runtime) error {

	// Load notes
	data, err := ioutil.ReadFile("notes.txt")
	if err != nil {
		return err
	}

	// Emit an event with the loaded data
	runtime.Events.Emit("notes loaded", string(data))

	m.runtime = runtime
	return nil
}
