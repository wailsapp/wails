//go:build darwin

package start_at_login

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/pkg/logger"
	"github.com/wailsapp/wails/v3/pkg/mac"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (p *Plugin) init() error {
	bundleID := mac.GetBundleID()
	if bundleID == "" {
		p.app.Log(&logger.Message{
			Level:   "INFO",
			Message: "Application is not in bundle. StartAtLogin will not work.",
		})
		p.disabled = true
	}
	return nil
}

func (p *Plugin) StartAtLogin(enabled bool) error {
	if p.disabled {
		return nil
	}
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

	cmd := exec.Command("osascript", "-e", command)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func (p *Plugin) IsStartAtLogin() (bool, error) {
	if p.disabled {
		return false, nil
	}
	exe, err := os.Executable()
	if err != nil {
		return false, err
	}
	binName := filepath.Base(exe)
	if !strings.HasSuffix(exe, "/Contents/MacOS/"+binName) {
		return false, fmt.Errorf("app needs to be running as package.app file to start at login")
	}
	appPath := strings.TrimSuffix(exe, "/Contents/MacOS/"+binName)
	appName := strings.TrimSuffix(filepath.Base(appPath), ".app")
	cmd := exec.Command("osascript", "-e", `tell application "System Events" to get the name of every login item`)
	results, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	resultsString := strings.TrimSpace(string(results))
	startupApps := strings.Split(resultsString, ", ")
	result := lo.Contains(startupApps, appName)
	return result, nil
}
