package binding

import (
	"fmt"
	"strings"

	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/lib/messages"
	"github.com/wailsapp/wails/runtime"
)

type internalMethods struct {
	log     *logger.CustomLogger
	browser *runtime.Browser
}

func newInternalMethods() *internalMethods {
	return &internalMethods{
		log:     logger.NewCustomLogger("InternalCall"),
		browser: runtime.NewBrowser(),
	}
}

func (i *internalMethods) processCall(callData *messages.CallData) (interface{}, error) {
	if !strings.HasPrefix(callData.BindingName, ".wails.") {
		return nil, fmt.Errorf("Invalid call signature '%s'", callData.BindingName)
	}

	// Strip prefix
	var splitCall = strings.Split(callData.BindingName, ".")[2:]
	if len(splitCall) != 2 {
		return nil, fmt.Errorf("Invalid call signature '%s'", callData.BindingName)
	}

	group := splitCall[0]
	switch group {
	case "Browser":
		return i.processBrowserCommand(splitCall[1], callData.Data)
	default:
		return nil, fmt.Errorf("Unknown internal command group '%s'", group)
	}
}

func (i *internalMethods) processBrowserCommand(command string, data interface{}) (interface{}, error) {
	switch command {
	case "OpenURL":
		url := data.(string)
		// Strip string quotes. Credit: https://stackoverflow.com/a/44222648
		if url[0] == '"' {
			url = url[1:]
		}
		if i := len(url) - 1; url[i] == '"' {
			url = url[:i]
		}
		i.log.Debugf("Calling Browser.OpenURL with '%s'", url)
		return nil, i.browser.OpenURL(url)
	case "OpenFile":
		filename := data.(string)
		// Strip string quotes. Credit: https://stackoverflow.com/a/44222648
		if filename[0] == '"' {
			filename = filename[1:]
		}
		if i := len(filename) - 1; filename[i] == '"' {
			filename = filename[:i]
		}
		i.log.Debugf("Calling Browser.OpenFile with '%s'", filename)
		return nil, i.browser.OpenFile(filename)
	default:
		return nil, fmt.Errorf("Unknown Browser command '%s'", command)
	}
}
