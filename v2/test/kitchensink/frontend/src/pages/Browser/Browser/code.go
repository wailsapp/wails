package main

import "github.com/wailsapp/wails/v2"

type MyStruct struct {
	runtime *wails.Runtime
}

func (l *MyStruct) ShowHelp() {
	l.runtime.Browser.Open("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
}
