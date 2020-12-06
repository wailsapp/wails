package main

import "github.com/wailsapp/wails/v2"

type MyStruct struct {
	runtime *wails.Runtime
}

func (m *MyStruct) WailsInit(runtime *wails.Runtime) error {

	maxAttempts := 3
	runtime.Events.OnMultiple("unlock attempts", func(optionalData ...interface{}) {
		// Do something (at most) maxAttempts times
	}, maxAttempts)

	m.runtime = runtime
	return nil
}
