package application

import (
	"net/http"
)

func (m *MessageProcessor) processClipboardMethod(method string, rw http.ResponseWriter, _ *http.Request, _ *WebviewWindow, params QueryParams) {

	switch method {
	case "SetText":
		text := params.String("text")
		if text == nil {
			m.Error("SetText: text is required")
			return
		}
		globalApplication.Clipboard().SetText(*text)
		m.ok(rw)
		m.Info("", "method", "Clipboard."+method, "text", *text)
	case "Text":
		text, _ := globalApplication.Clipboard().Text()
		m.text(rw, text)
		m.Info("", "method", "Clipboard."+method, "text", text)
	default:
		m.httpError(rw, "Unknown clipboard method: %s", method)
	}

}
