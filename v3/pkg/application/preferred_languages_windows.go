//go:build windows && !server

package application

import "golang.org/x/sys/windows"

// platformPreferredLanguages returns the display languages from Settings ›
// Time & language › Language & region, most preferred first.
// GetUserDefaultLocaleName would report the "Regional format" setting instead,
// which can differ from the display language.
func platformPreferredLanguages() []string {
	langs, err := windows.GetUserPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
	if err != nil {
		return nil
	}
	return langs
}
