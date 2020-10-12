package main

import wails "github.com/wailsapp/wails/v2"

type MyStruct struct {
	runtime *wails.Runtime
}

func (l *MyStruct) WailsInit(runtime *wails.Runtime) error {

	message := "Hello World!"
	runtime.Log.Print(message)
	// runtime.Log.Trace(message)
	// runtime.Log.Debug(message)
	// runtime.Log.Info(message)
	// runtime.Log.Warning(message)
	// runtime.Log.Error(message)
	// runtime.Log.Fatal(message)

	l.runtime = runtime
	return nil
}
