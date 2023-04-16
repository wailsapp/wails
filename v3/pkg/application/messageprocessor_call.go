package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (m *MessageProcessor) callErrorCallback(window *WebviewWindow, message string, callID *string, err error) {
	errorMsg := fmt.Sprintf(message, err)
	m.Error(errorMsg)
	msg := "_wails.callErrorCallback('" + *callID + "', " + strconv.Quote(errorMsg) + ");"
	window.ExecJS(msg)
}

func (m *MessageProcessor) callCallback(window *WebviewWindow, callID *string, result string, isJSON bool) {
	msg := fmt.Sprintf("_wails.callCallback('%s', %s, %v);", *callID, strconv.Quote(result), isJSON)
	window.ExecJS(msg)
}

func (m *MessageProcessor) processCallMethod(method string, rw http.ResponseWriter, _ *http.Request, window *WebviewWindow, params QueryParams) {
	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Unable to parse arguments: %s", err)
		return
	}
	callID := args.String("call-id")
	if callID == nil {
		m.Error("call-id is required")
		return
	}
	switch method {
	case "Call":
		var options CallOptions
		err := params.ToStruct(&options)
		if err != nil {
			m.callErrorCallback(window, "Error parsing call options: %s", callID, err)
			return
		}
		bindings := globalApplication.bindings.Get(&options)
		if bindings == nil {
			m.callErrorCallback(window, "Error getting binding for method: %s", callID, fmt.Errorf("'%s' not found", options.MethodName))
			return
		}
		go func() {
			result, err := bindings.Call(options.Args)
			if err != nil {
				m.callErrorCallback(window, "Error calling method: %s", callID, err)
				return
			}
			// convert result to json
			jsonResult, err := json.Marshal(result)
			if err != nil {
				m.callErrorCallback(window, "Error converting result to json: %s", callID, err)
				return
			}
			m.callCallback(window, callID, string(jsonResult), true)
		}()
		m.ok(rw)
	default:
		m.httpError(rw, "Unknown dialog method: %s", method)
	}

}
