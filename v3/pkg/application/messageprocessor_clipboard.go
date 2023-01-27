package application

import (
	"net/http"
)

func (m *MessageProcessor) processClipboardMethod(method string, rw http.ResponseWriter, r *http.Request, window *WebviewWindow, params QueryParams) {

	switch method {
	case "SetText":
		title := params.String("text")
		if title == nil {
			m.Error("SetText: text is required")
			return
		}
		globalApplication.Clipboard().SetText(*title)
		m.ok(rw)
	case "GetText":
		text := globalApplication.Clipboard().Text()
		m.text(rw, text)
	default:
		m.httpError(rw, "Unknown clipboard method: %s", method)
	}

}
