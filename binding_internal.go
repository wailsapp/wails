package wails

import "strings"
import "fmt"

type internalMethods struct{
	log *CustomLogger
	browser *RuntimeBrowser
}

func newInternalMethods() *internalMethods {
	return &internalMethods{
		log: newCustomLogger("InternalCall"),
		browser: newRuntimeBrowser(),
	}
}

func (i *internalMethods) processCall(callData *callData) (interface{}, error) {
	if !strings.HasPrefix(callData.BindingName, ".wails.") {
		return nil, fmt.Errorf("Invalid call signature '%s'", callData.BindingName)
	}

	// Strip prefix
	var splitCall = strings.Split(callData.BindingName,".")[2:]
	if len(splitCall) != 2 {
		return nil, fmt.Errorf("Invalid call signature '%s'", callData.BindingName)
	}

	group := splitCall[0]
	switch group {
	case "browser":
		return i.processBrowserCommand(splitCall[1], callData.Data)
	default:
		return nil, fmt.Errorf("Unknown internal command group '%s'", group)
	}
}

func (i *internalMethods) processBrowserCommand(command string, data interface{}) (interface{}, error) {
	switch command {
	case "openURL": 
		url := data.(string)
		// Strip string quotes. Credit: https://stackoverflow.com/a/44222648
		if url[0] == '"' {
			url = url[1:]
		}
		if i := len(url)-1; url[i] == '"' {
				url = url[:i]
		}
		i.log.Debugf("Calling browser.openURL with '%s'", url)
		return nil, i.browser.OpenURL(url)
	default:
		return nil, fmt.Errorf("Unknown browser command '%s'", command)
	}
}