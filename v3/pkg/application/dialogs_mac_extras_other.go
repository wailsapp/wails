//go:build !darwin || ios || server

package application

// Content type identifiers have no native panel support here; AddContentType
// translates them to extension filters instead (see dialogs_mac_extras.go).
const dialogContentTypesNative = false

func dialogPrompt(PromptOptions) (string, bool, error) {
	return "", false, ErrDialogNotSupported
}

func dialogPickColor(options ColorPickerOptions) (RGBA, bool, error) {
	return options.Initial, false, ErrDialogNotSupported
}

func dialogPickFont(FontPickerOptions) (FontDescriptor, bool, error) {
	return FontDescriptor{}, false, ErrDialogNotSupported
}
