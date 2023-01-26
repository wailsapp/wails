package application

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/options"
)

func (m *MessageProcessor) mustAtoI(input string) int {
	result, err := strconv.Atoi(input)
	if err != nil {
		m.Error("cannot convert %s to integer!", input)
	}
	return result
}

func (m *MessageProcessor) processWindowMessage(message string, window *WebviewWindow) {
	if len(message) < 2 {
		m.Error("Invalid Window Message: " + message)
	}

	switch message[1] {
	case 'A':
		switch message[2:] {
		//case "SDT":
		//	go window.WindowSetSystemDefaultTheme()
		//case "LT":
		//	go window.SetLightTheme()
		//case "DT":
		//	go window.SetDarkTheme()
		case "TP:0", "TP:1":
			if message[2:] == "TP:0" {
				go window.SetAlwaysOnTop(false)
			} else if message[2:] == "TP:1" {
				go window.SetAlwaysOnTop(true)
			}
		}
	case 'c':
		go window.Center()
	case 'T':
		title := message[2:]
		go window.SetTitle(title)
	case 'F':
		go window.Fullscreen()
	case 'f':
		go window.UnFullscreen()
	case 's':
		parts := strings.Split(message[3:], ":")
		w := m.mustAtoI(parts[0])
		h := m.mustAtoI(parts[1])
		go window.SetSize(w, h)
	case 'p':
		parts := strings.Split(message[3:], ":")
		x := m.mustAtoI(parts[0])
		y := m.mustAtoI(parts[1])
		go window.SetPosition(x, y)
	case 'H':
		go window.Hide()
	case 'S':
		go window.Show()
	//case 'R':
	//	go window.ReloadApp()
	case 'r':
		var rgba options.RGBA
		err := json.Unmarshal([]byte(message[3:]), &rgba)
		if err != nil {
			m.Error("Invalid RGBA Message: %s", err.Error())
		}
		go window.SetBackgroundColour(&rgba)
	case 'M':
		go window.Maximise()
	//case 't':
	//	go window.ToggleMaximise()
	case 'U':
		go window.UnMaximise()
	case 'm':
		go window.Minimise()
	case 'u':
		go window.UnMinimise()
	case 'Z':
		parts := strings.Split(message[3:], ":")
		w := m.mustAtoI(parts[0])
		h := m.mustAtoI(parts[1])
		go window.SetMaxSize(w, h)
	case 'z':
		parts := strings.Split(message[3:], ":")
		w := m.mustAtoI(parts[0])
		h := m.mustAtoI(parts[1])
		go window.SetMinSize(w, h)
	default:
		m.Error("unknown Window message: %s", message)
	}
}
