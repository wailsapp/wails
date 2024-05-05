//go:build windows

package start_at_login

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
	"strings"
)

func (p *Plugin) init() error {
	// TBD
	return nil
}

func (p *Plugin) getRegistryKey() (string, string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", "", fmt.Errorf("failed to get executable path: %s", err)
	}

	registryKey := p.options.RegistryKey
	if p.options.RegistryKey == "" {
		registryKey = strings.Split(filepath.Base(exePath), ".")[0]
	}

	return registryKey, exePath, nil
}

func openRegKey() (registry.Key, error) {
	// Open the registry key
	return registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
}

func (p *Plugin) IsStartAtLogin() (bool, error) {
	registryKey, exePath, err := p.getRegistryKey()
	if err != nil {
		return false, err
	}

	key, err := openRegKey()
	if err != nil {
		return false, err
	}
	defer key.Close()

	// Get the registry value
	value, _, err := key.GetStringValue(registryKey)
	if err != nil {
		return false, nil
	}
	return value == exePath, nil
}

func (p *Plugin) StartAtLogin(enabled bool) error {
	registryKey, exePath, err := p.getRegistryKey()
	if err != nil {
		return err
	}

	if enabled {
		// Open the registry key
		key, err := openRegKey()
		defer key.Close()

		// Set the registry value
		err = key.SetStringValue(registryKey, exePath)
		if err != nil {
			return fmt.Errorf("failed to set registry value: %s", err)
		}
	} else {
		// Remove registry key
		key, err := openRegKey()
		if err != nil {
			return fmt.Errorf("failed to open registry key: %s", err)
		}
		defer key.Close()

		// Remove the registry value
		err = key.DeleteValue(registryKey)
		if err != nil {
			return fmt.Errorf("failed to delete registry value: %s", err)
		}
	}
	return nil
}
