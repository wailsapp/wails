package dispatcher

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func (d *Dispatcher) mustAtoI(input string) int {
	result, err := strconv.Atoi(input)
	if err != nil {
		d.log.Error("cannot convert %s to integer!", input)
	}
	return result
}

func (d *Dispatcher) processWindowMessage(message string, sender frontend.Frontend) (string, error) {
	if len(message) < 2 {
		return "", errors.New("Invalid Window Message: " + message)
	}

	switch message[1] {
	case 'A':
		switch message[2:] {
		case "SDT":
			go sender.WindowSetSystemDefaultTheme()
		case "LT":
			go sender.WindowSetLightTheme()
		case "DT":
			go sender.WindowSetDarkTheme()
		case "TP:0", "TP:1":
			if message[2:] == "TP:0" {
				go sender.WindowSetAlwaysOnTop(false)
			} else if message[2:] == "TP:1" {
				go sender.WindowSetAlwaysOnTop(true)
			}
		}
	case 'c':
		go sender.WindowCenter()
	case 'T':
		title := message[2:]
		go sender.WindowSetTitle(title)
	case 'F':
		go sender.WindowFullscreen()
	case 'f':
		go sender.WindowUnfullscreen()
	case 's':
		parts := strings.Split(message[3:], ":")
		w := d.mustAtoI(parts[0])
		h := d.mustAtoI(parts[1])
		go sender.WindowSetSize(w, h)
	case 'p':
		parts := strings.Split(message[3:], ":")
		x := d.mustAtoI(parts[0])
		y := d.mustAtoI(parts[1])
		go sender.WindowSetPosition(x, y)
	case 'H':
		go sender.WindowHide()
	case 'S':
		go sender.WindowShow()
	case 'R':
		go sender.WindowReloadApp()
	case 'r':
		var rgba options.RGBA
		err := json.Unmarshal([]byte(message[3:]), &rgba)
		if err != nil {
			return "", err
		}
		go sender.WindowSetBackgroundColour(&rgba)
	case 'M':
		go sender.WindowMaximise()
	case 't':
		go sender.WindowToggleMaximise()
	case 'U':
		go sender.WindowUnmaximise()
	case 'm':
		go sender.WindowMinimise()
	case 'u':
		go sender.WindowUnminimise()
	case 'Z':
		parts := strings.Split(message[3:], ":")
		w := d.mustAtoI(parts[0])
		h := d.mustAtoI(parts[1])
		go sender.WindowSetMaxSize(w, h)
	case 'z':
		parts := strings.Split(message[3:], ":")
		w := d.mustAtoI(parts[0])
		h := d.mustAtoI(parts[1])
		go sender.WindowSetMinSize(w, h)
	default:
		d.log.Error("unknown Window message: %s", message)
	}

	return "", nil
}
