//go:build !darwin || ios || server

package application

import "os"

func platformAccessibilitySettings() AccessibilitySettings {
	return AccessibilitySettings{}
}

func platformKeyboardLayout() KeyboardLayout {
	return KeyboardLayout{Languages: []string{}}
}

func platformLocale() LocaleInfo {
	info := localeFromEnvironment(os.Getenv)
	if info.Preferred == nil {
		info.Preferred = []string{}
	}
	return info
}

func platformSetDefaultHandler(handler DefaultHandler) error {
	return ErrDefaultHandlerUnsupported
}

func platformDefaultHandler(handler DefaultHandler) (AppInfo, error) {
	return AppInfo{}, ErrDefaultHandlerUnsupported
}
