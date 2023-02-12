package application

import (
	"net/http"
)

type ContextMenuData struct {
	Id   string `json:"id"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Data any    `json:"data"`
}

func (m *MessageProcessor) processContextMenuMethod(method string, rw http.ResponseWriter, _ *http.Request, window *WebviewWindow, params QueryParams) {

	switch method {
	case "OpenContextMenu":
		var data ContextMenuData
		err := params.ToStruct(&data)
		if err != nil {
			m.httpError(rw, "error parsing contextmenu message: %s", err.Error())
			return
		}
		window.openContextMenu(&data)
		m.ok(rw)
	default:
		m.httpError(rw, "Unknown clipboard method: %s", method)
	}

}
