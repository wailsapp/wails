package runtime

import (
	"fmt"

	"github.com/wailsapp/wails/v2/internal/crypto"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	dialogoptions "github.com/wailsapp/wails/v2/pkg/options/dialog"
)

// Dialog defines all Dialog related operations
type Dialog interface {
	OpenFile(dialogOptions *dialogoptions.OpenDialog) (string, error)
	OpenMultipleFiles(dialogOptions *dialogoptions.OpenDialog) ([]string, error)
	OpenDirectory(dialogOptions *dialogoptions.OpenDialog) (string, error)
	SaveFile(dialogOptions *dialogoptions.SaveDialog) (string, error)
	Message(dialogOptions *dialogoptions.MessageDialog) (string, error)
}

// dialog exposes the Dialog interface
type dialog struct {
	bus *servicebus.ServiceBus
}

// newDialogs creates a new Dialogs struct
func newDialog(bus *servicebus.ServiceBus) Dialog {
	return &dialog{
		bus: bus,
	}
}

// processTitleAndFilter return the title and filter from the given params.
// title is the first string, filter is the second
func (r *dialog) processTitleAndFilter(params ...string) (string, string) {

	var title, filter string

	if len(params) > 0 {
		title = params[0]
	}

	if len(params) > 1 {
		filter = params[1]
	}

	return title, filter
}

// OpenDirectory prompts the user to select a directory
func (r *dialog) OpenDirectory(dialogOptions *dialogoptions.OpenDialog) (string, error) {

	// Create unique dialog callback
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:opendirectoryselected:" + uniqueCallback
	dialogResponseChannel, err := r.bus.Subscribe(responseTopic)
	if err != nil {
		return "", fmt.Errorf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	message := "dialog:select:directory:" + uniqueCallback
	r.bus.Publish(message, dialogOptions)

	// Wait for result
	var result = <-dialogResponseChannel

	// Delete subscription to response topic
	r.bus.UnSubscribe(responseTopic)

	return result.Data().(string), nil
}

// OpenFile prompts the user to select a file
func (r *dialog) OpenFile(dialogOptions *dialogoptions.OpenDialog) (string, error) {

	// Create unique dialog callback
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:openselected:" + uniqueCallback
	dialogResponseChannel, err := r.bus.Subscribe(responseTopic)
	if err != nil {
		return "", fmt.Errorf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	message := "dialog:select:open:" + uniqueCallback
	r.bus.Publish(message, dialogOptions)

	// Wait for result
	var result = <-dialogResponseChannel

	// Delete subscription to response topic
	r.bus.UnSubscribe(responseTopic)

	return result.Data().(string), nil
}

// OpenMultipleFiles prompts the user to select a file
func (r *dialog) OpenMultipleFiles(dialogOptions *dialogoptions.OpenDialog) ([]string, error) {

	// Create unique dialog callback
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:openmultipleselected:" + uniqueCallback
	dialogResponseChannel, err := r.bus.Subscribe(responseTopic)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	message := "dialog:select:openmultiple:" + uniqueCallback
	r.bus.Publish(message, dialogOptions)

	// Wait for result
	var result = <-dialogResponseChannel

	// Delete subscription to response topic
	r.bus.UnSubscribe(responseTopic)

	return result.Data().([]string), nil
}

// SaveFile prompts the user to select a file
func (r *dialog) SaveFile(dialogOptions *dialogoptions.SaveDialog) (string, error) {

	// Create unique dialog callback
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:saveselected:" + uniqueCallback
	dialogResponseChannel, err := r.bus.Subscribe(responseTopic)
	if err != nil {
		return "", fmt.Errorf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	message := "dialog:select:save:" + uniqueCallback
	r.bus.Publish(message, dialogOptions)

	// Wait for result
	var result = <-dialogResponseChannel

	// Delete subscription to response topic
	r.bus.UnSubscribe(responseTopic)

	return result.Data().(string), nil
}

// Message show a message to the user
func (r *dialog) Message(dialogOptions *dialogoptions.MessageDialog) (string, error) {

	// Create unique dialog callback
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:messageselected:" + uniqueCallback
	dialogResponseChannel, err := r.bus.Subscribe(responseTopic)
	if err != nil {
		return "", fmt.Errorf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	message := "dialog:select:message:" + uniqueCallback
	r.bus.Publish(message, dialogOptions)

	// Wait for result
	var result = <-dialogResponseChannel

	// Delete subscription to response topic
	r.bus.UnSubscribe(responseTopic)

	return result.Data().(string), nil
}
