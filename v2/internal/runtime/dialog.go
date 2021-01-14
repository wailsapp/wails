package runtime

import (
	"fmt"

	"github.com/wailsapp/wails/v2/internal/crypto"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	dialogoptions "github.com/wailsapp/wails/v2/pkg/options/dialog"
)

// Dialog defines all Dialog related operations
type Dialog interface {
	Open(dialogOptions *dialogoptions.OpenDialog) []string
	Save(dialogOptions *dialogoptions.SaveDialog) string
	Message(dialogOptions *dialogoptions.MessageDialog) string
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

// Open prompts the user to select a file
func (r *dialog) Open(dialogOptions *dialogoptions.OpenDialog) []string {

	// Create unique dialog callback
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:openselected:" + uniqueCallback
	dialogResponseChannel, err := r.bus.Subscribe(responseTopic)
	if err != nil {
		fmt.Printf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	message := "dialog:select:open:" + uniqueCallback
	r.bus.Publish(message, dialogOptions)

	// Wait for result
	var result *servicebus.Message = <-dialogResponseChannel

	// Delete subscription to response topic
	r.bus.UnSubscribe(responseTopic)

	return result.Data().([]string)
}

// Save prompts the user to select a file
func (r *dialog) Save(dialogOptions *dialogoptions.SaveDialog) string {

	// Create unique dialog callback
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:saveselected:" + uniqueCallback
	dialogResponseChannel, err := r.bus.Subscribe(responseTopic)
	if err != nil {
		fmt.Printf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	message := "dialog:select:save:" + uniqueCallback
	r.bus.Publish(message, dialogOptions)

	// Wait for result
	var result *servicebus.Message = <-dialogResponseChannel

	// Delete subscription to response topic
	r.bus.UnSubscribe(responseTopic)

	return result.Data().(string)
}

// Message show a message to the user
func (r *dialog) Message(dialogOptions *dialogoptions.MessageDialog) string {

	// Create unique dialog callback
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:messageselected:" + uniqueCallback
	dialogResponseChannel, err := r.bus.Subscribe(responseTopic)
	if err != nil {
		fmt.Printf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	message := "dialog:select:message:" + uniqueCallback
	r.bus.Publish(message, dialogOptions)

	// Wait for result
	var result *servicebus.Message = <-dialogResponseChannel

	// Delete subscription to response topic
	r.bus.UnSubscribe(responseTopic)

	return result.Data().(string)
}
