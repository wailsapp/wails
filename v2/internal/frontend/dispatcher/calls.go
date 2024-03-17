package dispatcher

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

type callMessage struct {
	Name       string            `json:"name"`
	Args       []json.RawMessage `json:"args"`
	CallbackID string            `json:"callbackID"`
}

func (d *Dispatcher) processCallMessage(message string, sender frontend.Frontend) (string, error) {
	var payload callMessage
	err := json.Unmarshal([]byte(message[1:]), &payload)
	if err != nil {
		return "", err
	}

	var result interface{}

	// Handle different calls
	switch true {
	case strings.HasPrefix(payload.Name, systemCallPrefix):
		result, err = d.processSystemCall(payload, sender)
	default:
		// Lookup method
		registeredMethod := d.bindingsDB.GetMethod(payload.Name)

		// Check we have it
		if registeredMethod == nil {
			return "", fmt.Errorf("method '%s' not registered", payload.Name)
		}

		args, err2 := registeredMethod.ParseArgs(payload.Args)
		if err2 != nil {
			errmsg := fmt.Errorf("error parsing arguments: %s", err2.Error())
			result, _ := d.NewErrorCallback(errmsg.Error(), payload.CallbackID)
			return result, errmsg
		}
		result, err = registeredMethod.Call(args)
	}

	callbackMessage := &CallbackMessage{
		CallbackID: payload.CallbackID,
	}
	if err != nil {
		// Use the error formatter if one was provided
		if d.errfmt != nil {
			callbackMessage.Err = d.errfmt(err)
		} else {
			callbackMessage.Err = err.Error()
		}
	} else {
		callbackMessage.Result = result
	}
	messageData, err := json.Marshal(callbackMessage)
	d.log.Trace("json call result data: %+v\n", string(messageData))
	if err != nil {
		// what now?
		d.log.Fatal(err.Error())
	}

	return "c" + string(messageData), nil
}

// CallbackMessage defines a message that contains the result of a call
type CallbackMessage struct {
	Result     interface{} `json:"result"`
	Err        any         `json:"error"`
	CallbackID string      `json:"callbackid"`
}

func (d *Dispatcher) NewErrorCallback(message string, callbackID string) (string, error) {
	result := &CallbackMessage{
		CallbackID: callbackID,
		Err:        message,
	}
	messageData, err := json.Marshal(result)
	d.log.Trace("json call result data: %+v\n", string(messageData))
	return string(messageData), err
}
