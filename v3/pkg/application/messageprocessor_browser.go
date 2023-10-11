package application

import (
	"github.com/pkg/browser"
	"net/http"
)

const (
	BrowserOpenURL = 0
)

var browserMethods = map[int]string{
	BrowserOpenURL: "OpenURL",
}

func (m *MessageProcessor) processBrowserMethod(method int, rw http.ResponseWriter, _ *http.Request, _ Window, params QueryParams) {

	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Unable to parse arguments: %s", err.Error())
		return
	}

	switch method {
	case BrowserOpenURL:
		url := args.String("url")
		if url == nil {
			m.Error("OpenURL: url is required")
			return
		}
		err := browser.OpenURL(*url)
		if err != nil {
			m.Error("OpenURL: %s", err.Error())
			return
		}
		m.ok(rw)
		m.Info("Runtime Call:", "method", "Browser."+browserMethods[method], "url", *url)
	default:
		m.httpError(rw, "Unknown browser method: %d", method)
		return
	}

}
