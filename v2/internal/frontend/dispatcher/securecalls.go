package dispatcher

import (
	"encoding/json"
	"fmt"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

type secureCallMessage struct {
	ID         int               `json:"id"`
	Args       []json.RawMessage `json:"args"`
	CallbackID string            `json:"callbackID"`
}

func (d *Dispatcher) processSecureCallMessage(message string, sender frontend.Frontend) (string, error) {
	var payload secureCallMessage
	err := json.Unmarshal([]byte(message[1:]), &payload)
	if err != nil {
		return "", err
	}

	var result interface{}

	// Lookup method
	registeredMethod := d.bindingsDB.GetObfuscatedMethod(payload.ID)

	// Check we have it
	if registeredMethod == nil {
		return "", fmt.Errorf("method '%d' not registered", payload.ID)
	}

	args, err2 := registeredMethod.ParseArgs(payload.Args)
	if err2 != nil {
		errmsg := fmt.Errorf("error parsing arguments: %s", err2.Error())
		result, _ := d.NewErrorCallback(errmsg.Error(), payload.CallbackID)
		return result, errmsg
	}
	result, err = registeredMethod.Call(args)

	callbackMessage := &CallbackMessage{
		CallbackID: payload.CallbackID,
	}
	if err != nil {
		callbackMessage.Err = err.Error()
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
