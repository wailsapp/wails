package application

import (
	"github.com/wailsapp/wails/v3/pkg/errs"
)

const (
	ClipboardSetText = 0
	ClipboardText    = 1
)

var clipboardMethods = map[int]string{
	ClipboardSetText: "SetText",
	ClipboardText:    "Text",
}

func (m *MessageProcessor) processClipboardMethod(req *RuntimeRequest) (any, error) {
	args := req.Args.AsMap()

	var text string

	switch req.Method {
	case ClipboardSetText:
		textp := args.String("text")
		if textp == nil {
			return nil, errs.NewInvalidClipboardCallErrorf("missing argument 'text'")
		}
		text = *textp
		globalApplication.Clipboard.SetText(text)
		return unit, nil
	case ClipboardText:
		text, _ = globalApplication.Clipboard.Text()
		return text, nil
	default:
		return nil, errs.NewInvalidClipboardCallErrorf("unknown method: %d", req.Method)
	}
}
