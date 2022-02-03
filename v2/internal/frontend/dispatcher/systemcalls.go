package dispatcher

import (
	"fmt"
	"strings"

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
	default:
		return nil, fmt.Errorf("unknown systemcall message: %s", payload.Name)
	}

}
