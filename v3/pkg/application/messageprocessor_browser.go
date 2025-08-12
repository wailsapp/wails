package application

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pkg/browser"
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
		m.httpError(rw, "Invalid browser call:", fmt.Errorf("unable to parse arguments: %w", err))
		return
	}

	switch method {
	case BrowserOpenURL:
		url := args.String("url")
		if url == nil {
			m.httpError(rw, "Invalid browser call:", errors.New("missing argument 'url'"))
			return
		}

		sanitizedURL, err := ValidateAndSanitizeURL(*url)
		if err != nil {
			m.Error("OpenURL: invalid URL - %s", err.Error())
			m.httpError(rw, fmt.Sprintf("Invalid URL: %s", err.Error()), err)
			return
		}

		err = browser.OpenURL(sanitizedURL)
		if err != nil {
			m.httpError(rw, "OpenURL failed:", err)
			return
		}

		m.ok(rw)
		m.Info("Runtime call:", "method", "Browser."+browserMethods[method], "url", sanitizedURL)
	default:
		m.httpError(rw, "Invalid browser call:", fmt.Errorf("unknown method: %d", method))
		return
	}
}
