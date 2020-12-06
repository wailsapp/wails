package main

import "github.com/wailsapp/wails/v2"

type MyStruct struct {
	runtime *wails.Runtime
}

func (m *MyStruct) WailsInit(runtime *wails.Runtime) error {

	runtime.Events.Once("initialised", func(optionalData ...interface{}) {
		// Do something once
	})

	m.runtime = runtime
	return nil
}
