package application

import (
	"errors"
	"fmt"
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
		m.httpError(rw, "Invalid clipboard call:", fmt.Errorf("unable to parse arguments: %w", err))
		return
	}

	var text string

	switch method {
	case ClipboardSetText:
		textp := args.String("text")
		if textp == nil {
			m.httpError(rw, "Invalid clipboard call:", errors.New("missing argument 'text'"))
			return
		}
		text = *textp
		globalApplication.Clipboard().SetText(text)
		m.ok(rw)
	case ClipboardText:
		text, _ = globalApplication.Clipboard().Text()
		m.text(rw, text)
	default:
		m.httpError(rw, "Invalid clipboard call:", fmt.Errorf("unknown method: %d", method))
		return
	}

	m.Info("Runtime call:", "method", "Clipboard."+clipboardMethods[method], "text", text)
}
