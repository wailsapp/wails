package logutils

import (
	"fmt"

	"github.com/wailsapp/wails/v2/internal/tui"
)

func LogGreen(message string, args ...interface{}) {
	if len(message) == 0 {
		return
	}
	text := fmt.Sprintf(message, args...)
	println(tui.Green(text))
}

func LogRed(message string, args ...interface{}) {
	if len(message) == 0 {
		return
	}
	text := fmt.Sprintf(message, args...)
	println(tui.Red(text))
}

func LogDarkYellow(message string, args ...interface{}) {
	if len(message) == 0 {
		return
	}
	text := fmt.Sprintf(message, args...)
	println(tui.DarkYellow(text))
}
