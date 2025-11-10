package application

import (
	"github.com/wailsapp/wails/v3/pkg/errs"
)

type ContextMenuData struct {
	Id   string `json:"id"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Data string `json:"data"`
}

func (d ContextMenuData) clone() *ContextMenuData {
	return &ContextMenuData{
		Id:   d.Id,
		X:    d.X,
		Y:    d.Y,
		Data: d.Data,
	}
}

const (
	ContextMenuOpen = 0
)

var contextmenuMethodNames = map[int]string{
	ContextMenuOpen: "Open",
}

func (m *MessageProcessor) processContextMenuMethod(req *RuntimeRequest, window Window) (any, error) {
	switch req.Method {
	case ContextMenuOpen:
		var data ContextMenuData
		err := req.Args.ToStruct(&data)
		if err != nil {
			return nil, errs.WrapInvalidContextMenuCallErrorf(err, "error parsing parameters")
		}

		window.OpenContextMenu(&data)

		return unit, err
	default:
		return nil, errs.NewInvalidContextMenuCallErrorf("unknown method: %d", req.Method)
	}
}
