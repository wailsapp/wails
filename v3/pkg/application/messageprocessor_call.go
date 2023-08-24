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

func (m *MessageProcessor) processCallMethod(method string, rw http.ResponseWriter, r *http.Request, window *WebviewWindow, params QueryParams) {
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
		var boundMethod *BoundMethod
		methodID := r.Header.Get("x-wails-method-id")
		if methodID == "" {
			boundMethod = globalApplication.bindings.Get(&options)
			if boundMethod == nil {
				m.callErrorCallback(window, "Error getting binding for method: %s", callID, fmt.Errorf("method '%s' not found", options.Name()))
				return
			}
		} else {
			id, err := strconv.ParseUint(methodID, 10, 32)
			if err != nil {
				m.callErrorCallback(window, "Error parsing method id for call: %s", callID, err)
				return
			}
			boundMethod = globalApplication.bindings.GetByID(uint32(id))
		}
		if boundMethod == nil {
			m.callErrorCallback(window, "Error getting binding for method: %s", callID, fmt.Errorf("method ID '%s' not found", options.Name()))
			return
		}
		go func() {
			result, err := boundMethod.Call(options.Args)
			if err != nil {
				m.callErrorCallback(window, "Error calling method: %s", callID, err)
				return
			}
			var jsonResult = []byte("{}")
			if result != nil {
				// convert result to json
				jsonResult, err = json.Marshal(result)
				if err != nil {
					m.callErrorCallback(window, "Error converting result to json: %s", callID, err)
					return
				}
			}
			m.callCallback(window, callID, string(jsonResult), true)
			m.Info("Call Binding:", "method", boundMethod, "args", options.Args, "result", result)
		}()
		m.ok(rw)
	default:
		m.httpError(rw, "Unknown call method: %s", method)
	}

}
