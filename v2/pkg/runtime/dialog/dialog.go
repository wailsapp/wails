// +build !experimental

package dialog

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/internal/crypto"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// FileFilter defines a filter for dialog boxes
type FileFilter struct {
	DisplayName string // Filter information EG: "Image Files (*.jpg, *.png)"
	Pattern     string // semi-colon separated list of extensions, EG: "*.jpg;*.png"
}

// OpenDialogOptions contains the options for the OpenDialogOptions runtime method
type OpenDialogOptions struct {
	DefaultDirectory           string
	DefaultFilename            string
	Title                      string
	Filters                    []FileFilter
	AllowFiles                 bool
	AllowDirectories           bool
	ShowHiddenFiles            bool
	CanCreateDirectories       bool
	ResolvesAliases            bool
	TreatPackagesAsDirectories bool
}

// SaveDialogOptions contains the options for the SaveDialog runtime method
type SaveDialogOptions struct {
	DefaultDirectory           string
	DefaultFilename            string
	Title                      string
	Filters                    []FileFilter
	ShowHiddenFiles            bool
	CanCreateDirectories       bool
	TreatPackagesAsDirectories bool
}

type DialogType string

const (
	InfoDialog     DialogType = "info"
	WarningDialog  DialogType = "warning"
	ErrorDialog    DialogType = "error"
	QuestionDialog DialogType = "question"
)

// MessageDialogOptions contains the options for the Message dialogs, EG Info, Warning, etc runtime methods
type MessageDialogOptions struct {
	Type          DialogType
	Title         string
	Message       string
	Buttons       []string
	DefaultButton string
	CancelButton  string
	Icon          string
}

// processTitleAndFilter return the title and filter from the given params.
// title is the first string, filter is the second
func processTitleAndFilter(params ...string) (string, string) {

	var title, filter string

	if len(params) > 0 {
		title = params[0]
	}

	if len(params) > 1 {
		filter = params[1]
	}

	return title, filter
}

type Dialog struct{}

// OpenDirectory prompts the user to select a directory
func (d *Dialog) OpenDirectory(ctx context.Context, dialogOptions OpenDialogOptions) (string, error) {

	bus := servicebus.ExtractBus(ctx)

	// Create unique dialog callback
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:opendirectoryselected:" + uniqueCallback
	dialogResponseChannel, err := bus.Subscribe(responseTopic)
	if err != nil {
		return "", fmt.Errorf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	message := "dialog:select:directory:" + uniqueCallback
	bus.Publish(message, dialogOptions)

	// Wait for result
	var result = <-dialogResponseChannel

	// Delete subscription to response topic
	bus.UnSubscribe(responseTopic)

	return result.Data().(string), nil
}

// OpenFile prompts the user to select a file
func (d *Dialog) OpenFile(ctx context.Context, dialogOptions OpenDialogOptions) (string, error) {

	bus := servicebus.ExtractBus(ctx)

	// Create unique dialog callback
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:openselected:" + uniqueCallback
	dialogResponseChannel, err := bus.Subscribe(responseTopic)
	if err != nil {
		return "", fmt.Errorf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	message := "dialog:select:open:" + uniqueCallback
	bus.Publish(message, dialogOptions)

	// Wait for result
	var result = <-dialogResponseChannel

	// Delete subscription to response topic
	bus.UnSubscribe(responseTopic)

	return result.Data().(string), nil
}

// OpenMultipleFiles prompts the user to select a file
func (d *Dialog) OpenMultipleFiles(ctx context.Context, dialogOptions OpenDialogOptions) ([]string, error) {

	bus := servicebus.ExtractBus(ctx)
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:openmultipleselected:" + uniqueCallback
	dialogResponseChannel, err := bus.Subscribe(responseTopic)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	message := "dialog:select:openmultiple:" + uniqueCallback
	bus.Publish(message, dialogOptions)

	// Wait for result
	var result = <-dialogResponseChannel

	// Delete subscription to response topic
	bus.UnSubscribe(responseTopic)

	return result.Data().([]string), nil
}

// SaveFile prompts the user to select a file
func (d *Dialog) SaveFile(ctx context.Context, dialogOptions SaveDialogOptions) (string, error) {

	bus := servicebus.ExtractBus(ctx)
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:saveselected:" + uniqueCallback
	dialogResponseChannel, err := bus.Subscribe(responseTopic)
	if err != nil {
		return "", fmt.Errorf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	message := "dialog:select:save:" + uniqueCallback
	bus.Publish(message, dialogOptions)

	// Wait for result
	var result = <-dialogResponseChannel

	// Delete subscription to response topic
	bus.UnSubscribe(responseTopic)

	return result.Data().(string), nil
}

// Message show a message to the user
func (d *Dialog) Message(ctx context.Context, dialogOptions MessageDialogOptions) (string, error) {

	bus := servicebus.ExtractBus(ctx)

	// Create unique dialog callback
	uniqueCallback := crypto.RandomID()

	// Subscribe to the respose channel
	responseTopic := "dialog:messageselected:" + uniqueCallback
	dialogResponseChannel, err := bus.Subscribe(responseTopic)
	if err != nil {
		return "", fmt.Errorf("ERROR: Cannot subscribe to bus topic: %+v\n", err.Error())
	}

	message := "dialog:select:message:" + uniqueCallback
	bus.Publish(message, dialogOptions)

	// Wait for result
	var result = <-dialogResponseChannel

	// Delete subscription to response topic
	bus.UnSubscribe(responseTopic)

	return result.Data().(string), nil
}
