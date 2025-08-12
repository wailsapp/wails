// Package mac provides MacOS related utility functions for Wails applications
package mac

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/shell"
)

// ShowNotification will either add or remove this application to/from the login
// items, depending on the given boolean flag. The limitation is that the
// currently running app must be in an app bundle.
func ShowNotification(title string, subtitle string, message string, sound string) error {
	command := fmt.Sprintf("display notification \"%s\"", message)
	if len(title) > 0 {
		command += fmt.Sprintf(" with title \"%s\"", title)
	}
	if len(subtitle) > 0 {
		command += fmt.Sprintf(" subtitle \"%s\"", subtitle)
	}
	if len(sound) > 0 {
		command += fmt.Sprintf(" sound name \"%s\"", sound)
	}
	_, stde, err := shell.RunCommand("/tmp", "osascript", "-e", command)
	if err != nil {
		return errors.Wrap(err, stde)
	}
	return nil
}
