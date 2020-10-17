package main

import (
	wails "github.com/wailsapp/wails/v2"
)

// Events struct
type Events struct {
	runtime *wails.Runtime
}

// WailsInit is called at application startup
func (e *Events) WailsInit(runtime *wails.Runtime) error {
	// Perform your setup here
	e.runtime = runtime
	return nil
}

// Subscribe will subscribe
func (e *Events) Subscribe(eventName string) {
	e.runtime.Events.On(eventName, func(args ...interface{}) {
		type callbackData struct {
			Name string
			Data []interface{}
		}
		result := callbackData{Name: eventName, Data: args}
		e.runtime.Events.Emit("event fired by go subscriber", result)
	})
}
