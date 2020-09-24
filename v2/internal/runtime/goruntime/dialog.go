package goruntime

import (
	"fmt"

	"github.com/wailsapp/wails/v2/internal/crypto"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Dialog defines all Dialog related operations
type Dialog interface {
	SaveFile(params ...string) string
	SelectFile(params ...string) string
	SelectDirectory(params ...string) []string
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

// SelectFile prompts the user to select a file
func (r *dialog) SelectFile(params ...string) string {

	// Extract title + filter
	title, filter := r.processTitleAndFilter(params...)

	// Create unique dialog callback
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:fileselected:" + uniqueCallback
	dialogResponseChannel, err := r.bus.Subscribe(responseTopic)
	if err != nil {
		fmt.Printf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	// Publish dialog request
	message := "dialog:select:file:" + title
	if filter != "" {
		message += ":" + filter
	}
	r.bus.Publish(message, responseTopic)

	// Wait for result
	result := <-dialogResponseChannel

	// Delete subscription to response topic
	r.bus.UnSubscribe(responseTopic)

	return result.Data().(string)
}

// SaveFile prompts the user to select a file to save to
func (r *dialog) SaveFile(params ...string) string {

	// Extract title + filter
	title, filter := r.processTitleAndFilter(params...)

	// Create unique dialog callback
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:filesaveselected:" + uniqueCallback
	dialogResponseChannel, err := r.bus.Subscribe(responseTopic)
	if err != nil {
		fmt.Printf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	// Publish dialog request
	message := "dialog:select:filesave:" + title
	if filter != "" {
		message += ":" + filter
	}
	r.bus.Publish(message, responseTopic)

	// Wait for result
	result := <-dialogResponseChannel

	// Delete subscription to response topic
	r.bus.UnSubscribe(responseTopic)

	return result.Data().(string)
}

// SelectDirectory prompts the user to select a file
func (r *dialog) SelectDirectory(params ...string) []string {

	// Extract title + filter
	title, filter := r.processTitleAndFilter(params...)

	// Create unique dialog callback
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:directoryselected:" + uniqueCallback
	dialogResponseChannel, err := r.bus.Subscribe(responseTopic)
	if err != nil {
		fmt.Printf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	// Publish dialog request
	message := "dialog:select:directory:" + title
	if filter != "" {
		message += ":" + filter
	}
	r.bus.Publish(message, responseTopic)

	// Wait for result
	var result *servicebus.Message = <-dialogResponseChannel

	// Delete subscription to response topic
	r.bus.UnSubscribe(responseTopic)

	return result.Data().([]string)
}
