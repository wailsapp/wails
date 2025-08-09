package dispatcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

const systemCallPrefix = ":wails:"

type position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type size struct {
	W int `json:"w"`
	H int `json:"h"`
}

func (d *Dispatcher) processSystemCall(payload callMessage, sender frontend.Frontend) (interface{}, error) {
	// Strip prefix
	name := strings.TrimPrefix(payload.Name, systemCallPrefix)

	switch name {
	case "WindowGetPos":
		x, y := sender.WindowGetPosition()
		return &position{x, y}, nil
	case "WindowGetSize":
		w, h := sender.WindowGetSize()
		return &size{w, h}, nil
	case "ScreenGetAll":
		return sender.ScreenGetAll()
	case "WindowIsMaximised":
		return sender.WindowIsMaximised(), nil
	case "WindowIsMinimised":
		return sender.WindowIsMinimised(), nil
	case "WindowIsNormal":
		return sender.WindowIsNormal(), nil
	case "WindowIsFullscreen":
		return sender.WindowIsFullscreen(), nil
	case "Environment":
		return runtime.Environment(d.ctx), nil
	case "ClipboardGetText":
		t, err := sender.ClipboardGetText()
		return t, err
	case "ClipboardSetText":
		if len(payload.Args) < 1 {
			return false, errors.New("empty argument, cannot set clipboard")
		}
		var arg string
		if err := json.Unmarshal(payload.Args[0], &arg); err != nil {
			return false, err
		}
		if err := sender.ClipboardSetText(arg); err != nil {
			return false, err
		}
		return true, nil
	default:
		return nil, fmt.Errorf("unknown systemcall message: %s", payload.Name)
	}
}
