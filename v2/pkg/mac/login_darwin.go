// Package mac provides MacOS related utility functions for Wails applications
package mac

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/leaanthony/slicer"
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/shell"
)

// StartAtLogin will either add or remove this application to/from the login
// items, depending on the given boolean flag. The limitation is that the
// currently running app must be in an app bundle.
func StartAtLogin(enabled bool) error {
	exe, err := os.Executable()
	if err != nil {
		return errors.Wrap(err, "Error running os.Executable:")
	}
	binName := filepath.Base(exe)
	if !strings.HasSuffix(exe, "/Contents/MacOS/"+binName) {
		return fmt.Errorf("app needs to be running as package.app file to start at login")
	}
	appPath := strings.TrimSuffix(exe, "/Contents/MacOS/"+binName)
	var command string
	if enabled {
		command = fmt.Sprintf("tell application \"System Events\" to make login item at end with properties {name: \"%s\",path:\"%s\", hidden:false}", binName, appPath)
	} else {
		command = fmt.Sprintf("tell application \"System Events\" to delete login item \"%s\"", binName)
	}
	_, stde, err := shell.RunCommand("/tmp", "osascript", "-e", command)
	if err != nil {
		return errors.Wrap(err, stde)
	}
	return nil
}

// StartsAtLogin will indicate if this application is in the login
// items. The limitation is that the currently running app must be
// in an app bundle.
func StartsAtLogin() (bool, error) {
	exe, err := os.Executable()
	if err != nil {
		return false, err
	}
	binName := filepath.Base(exe)
	if !strings.HasSuffix(exe, "/Contents/MacOS/"+binName) {
		return false, fmt.Errorf("app needs to be running as package.app file to start at login")
	}
	results, stde, err := shell.RunCommand("/tmp", "osascript", "-e", `tell application "System Events" to get the name of every login item`)
	if err != nil {
		return false, errors.Wrap(err, stde)
	}
	results = strings.TrimSpace(results)
	startupApps := slicer.String(strings.Split(results, ", "))
	return startupApps.Contains(binName), nil
}
