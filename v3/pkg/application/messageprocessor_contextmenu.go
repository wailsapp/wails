package application

import (
	"fmt"
	"net/http"
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

func (m *MessageProcessor) processContextMenuMethod(method int, rw http.ResponseWriter, _ *http.Request, window Window, params QueryParams) {
	switch method {
	case ContextMenuOpen:
		var data ContextMenuData
		err := params.ToStruct(&data)
		if err != nil {
			m.httpError(rw, "Invalid contextmenu call:", fmt.Errorf("error parsing parameters: %w", err))
			return
		}

		window.OpenContextMenu(&data)

		m.ok(rw)
		m.Info("Runtime call:", "method", "ContextMenu."+contextmenuMethodNames[method], "id", data.Id, "x", data.X, "y", data.Y, "data", data.Data)
	default:
		m.httpError(rw, "Invalid contextmenu call:", fmt.Errorf("unknown method: %d", method))
		return
	}
}
