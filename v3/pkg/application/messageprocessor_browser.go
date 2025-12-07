package application

import (
	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v3/pkg/errs"
)

const (
	BrowserOpenURL = 0
)

var browserMethodNames = map[int]string{
	BrowserOpenURL: "OpenURL",
}

func (m *MessageProcessor) processBrowserMethod(req *RuntimeRequest) (any, error) {
	switch req.Method {
	case BrowserOpenURL:
		url := req.Args.AsMap().String("url")
		if url == nil {
			return nil, errs.NewInvalidBrowserCallErrorf("missing argument 'url'")
		}

		sanitizedURL, err := ValidateAndSanitizeURL(*url)
		if err != nil {
			return nil, errs.WrapInvalidBrowserCallErrorf(err, "invalid URL")
		}

		err = browser.OpenURL(sanitizedURL)
		if err != nil {
			m.Error("OpenURL: invalid URL - %s", err.Error())
			return nil, errs.WrapInvalidBrowserCallErrorf(err, "OpenURL failed")
		}

		return unit, nil
	default:
		return nil, errs.NewInvalidBrowserCallErrorf("unknown method: %d", req.Method)
	}
}
