package main

import (
	"github.com/wailsapp/wails/v2"
)

type MyStruct struct {
	runtime *wails.Runtime
}

// ShowPlatformHelp shows specific help for the platform
func (l *MyStruct) ShowPlatformHelp() {
	l.runtime.Browser.Open("https://wails.app/gettingstarted/" + l.runtime.System.Platform())
}
