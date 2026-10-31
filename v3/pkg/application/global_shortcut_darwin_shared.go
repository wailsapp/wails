//go:build darwin && !ios && !server

package application

import "fmt"

func (g *macosGlobalShortcuts) unregisterAll() error {
	var firstErr error
	for id := range g.refs {
		if err := g.unregister(id); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func macKeyCodeFor(key string) (int, error) {
	if len(key) == 1 && key[0] >= 'a' && key[0] <= 'z' {
		if code := keyCodeForChar(uint16(key[0])); code >= 0 {
			return code, nil
		}
	}
	if code, ok := macKeyCodes[key]; ok {
		return code, nil
	}
	return 0, fmt.Errorf("key %q is not supported as a global shortcut on macOS", key)
}
