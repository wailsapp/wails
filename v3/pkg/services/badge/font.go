package badge

import (
	"errors"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// FontManager handles font discovery on Windows with minimal caching
type FontManager struct {
	fontCache     map[string]string // Maps only requested font filenames to paths
	fontDirs      []string          // Directories to search for fonts
	mu            sync.RWMutex      // Mutex for thread-safe access to the cache
	registryPaths []string          // Registry paths to search for fonts
}

// NewFontManager creates a new FontManager instance
func NewFontManager() *FontManager {
	return &FontManager{
		fontCache: make(map[string]string),
		fontDirs: []string{
			filepath.Join(os.Getenv("windir"), "Fonts"),
			filepath.Join(os.Getenv("localappdata"), "Microsoft", "Windows", "Fonts"),
		},
		registryPaths: []string{
			`SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`,
		},
	}
}

// FindFont searches for a font by filename and returns its full path
// Only caches fonts that are found
func (fm *FontManager) FindFont(fontFilename string) (string, error) {
	fontKey := strings.ToLower(fontFilename)

	// Check if already in cache
	fm.mu.RLock()
	if path, exists := fm.fontCache[fontKey]; exists {
		fm.mu.RUnlock()
		return path, nil
	}
	fm.mu.RUnlock()

	// If not in cache, search for the font
	fontPath, err := fm.searchForFont(fontFilename)
	if err != nil {
		return "", err
	}

	// Add to cache only if found
	fm.mu.Lock()
	fm.fontCache[fontKey] = fontPath
	fm.mu.Unlock()

	return fontPath, nil
}

// searchForFont looks for a font in all known locations
func (fm *FontManager) searchForFont(fontFilename string) (string, error) {
	fontFileLower := strings.ToLower(fontFilename)

	// 1. Direct file check in font directories (fastest approach)
	for _, dir := range fm.fontDirs {
		fontPath := filepath.Join(dir, fontFilename)
		if fileExists(fontPath) {
			return fontPath, nil
		}
	}

	// 2. Search in registry (can find fonts with different paths)
	// System fonts (HKEY_LOCAL_MACHINE)
	for _, regPath := range fm.registryPaths {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath, registry.QUERY_VALUE)
		if err == nil {
			defer k.Close()

			// Look for the specific font in registry values
			fontPath, found := fm.findFontInRegistry(k, fontFileLower, fm.fontDirs[0])
			if found {
				return fontPath, nil
			}
		}
	}

	// 3. User fonts (HKEY_CURRENT_USER)
	for _, regPath := range fm.registryPaths {
		k, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.QUERY_VALUE)
		if err == nil {
			defer k.Close()

			// Look for the specific font in registry values
			fontPath, found := fm.findFontInRegistry(k, fontFileLower, fm.fontDirs[1])
			if found {
				return fontPath, nil
			}
		}
	}

	return "", errors.New("font not found: " + fontFilename)
}

// findFontInRegistry searches for a specific font in a registry key
func (fm *FontManager) findFontInRegistry(k registry.Key, fontFileLower string, defaultDir string) (string, bool) {
	valueNames, err := k.ReadValueNames(0)
	if err != nil {
		return "", false
	}

	for _, name := range valueNames {
		value, _, err := k.GetStringValue(name)
		if err != nil {
			continue
		}

		// Check if this registry entry corresponds to our font
		valueLower := strings.ToLower(value)
		if strings.HasSuffix(valueLower, fontFileLower) {
			// If it's a relative path, assume it's in the default font directory
			if !strings.Contains(value, "\\") {
				value = filepath.Join(defaultDir, value)
			}

			if fileExists(value) {
				return value, true
			}
		}
	}

	return "", false
}

func (fm *FontManager) FindFontOrDefault(name string) string {
	fontsToFind := []string{name, "segoeuib.ttf", "arialbd.ttf"}
	for _, font := range fontsToFind {
		path, err := fm.FindFont(font)
		if err == nil {
			return path
		}
	}
	return ""
}

// Helper functions
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
