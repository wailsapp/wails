package main

import (
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

// Logger struct
type Logger struct {
	runtime *wails.Runtime
}

func (l *Logger) WailsInit(runtime *wails.Runtime) error {

	runtime.Log.SetLogLevel(logger.TRACE)
	// runtime.Log.SetLogLevel(logger.DEBUG)
	// runtime.Log.SetLogLevel(logger.INFO)
	// runtime.Log.SetLogLevel(logger.WARNING)
	// runtime.Log.SetLogLevel(logger.ERROR)

	l.runtime = runtime
	return nil
}
