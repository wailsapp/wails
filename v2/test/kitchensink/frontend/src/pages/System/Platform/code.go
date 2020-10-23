package main

import (
	wails "github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/internal/runtime"
)

type MyStruct struct {
	runtime *wails.Runtime
}

func (l *MyStruct) ShowPlatformHelp() {
	l.runtime.Browser.Open("https://wails.app/gettingstarted/" + runtime.System.Platform())
}
