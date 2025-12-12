package application

import (
	"encoding/json"

	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v3/pkg/errs"
)

const (
	BrowserOpenURL  = 0
	BrowserSendData = 1
)

var browserMethodNames = map[int]string{
	BrowserOpenURL:  "OpenURL",
	BrowserSendData: "SendData",
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

	case BrowserSendData:
		// Handle browser data extraction results from browser mode windows
		dataStr := req.Args.AsMap().String("data")
		if dataStr == nil {
			return nil, errs.NewInvalidBrowserCallErrorf("missing argument 'data'")
		}

		var browserData BrowserData
		err := json.Unmarshal([]byte(*dataStr), &browserData)
		if err != nil {
			m.Error("SendData: failed to parse browser data - %s", err.Error())
			return nil, errs.WrapInvalidBrowserCallErrorf(err, "failed to parse browser data")
		}

		// Store in global browser data store
		GetBrowserDataStore().Store(browserData.WindowName, &browserData)

		// Call the OnDataExtracted callback if set
		window, exists := globalApplication.Window.GetByName(browserData.WindowName)
		if exists {
			if webviewWindow, ok := window.(*WebviewWindow); ok {
				if webviewWindow.options.BrowserMode != nil && webviewWindow.options.BrowserMode.OnDataExtracted != nil {
					webviewWindow.options.BrowserMode.OnDataExtracted(&browserData)
				}
			}
		}

		return unit, nil

	default:
		return nil, errs.NewInvalidBrowserCallErrorf("unknown method: %d", req.Method)
	}
}
