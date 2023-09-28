package application

import (
	"net/http"
)

const (
	ClipboardSetText = 0
	ClipboardText    = 1
)

var clipboardMethods = map[int]string{
	ClipboardSetText: "SetText",
	ClipboardText:    "Text",
}

func (m *MessageProcessor) processClipboardMethod(method int, rw http.ResponseWriter, _ *http.Request, _ Window, params QueryParams) {

	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Unable to parse arguments: %s", err.Error())
		return
	}

	switch method {
	case ClipboardSetText:
		text := args.String("text")
		if text == nil {
			m.Error("SetText: text is required")
			return
		}
		globalApplication.Clipboard().SetText(*text)
		m.ok(rw)
		m.Info("Runtime Call:", "method", "Clipboard."+clipboardMethods[method], "text", *text)
	case ClipboardText:
		text, _ := globalApplication.Clipboard().Text()
		m.text(rw, text)
		m.Info("Runtime Call:", "method", "Clipboard."+clipboardMethods[method], "text", text)
	default:
		m.httpError(rw, "Unknown clipboard method: %d", method)
		return
	}

}
