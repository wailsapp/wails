package reportcopy

import "github.com/wailsapp/wails/v3/pkg/w32"

// Write copies text using the existing native Windows clipboard support.
func Write(text string) error { return w32.SetClipboardText(text) }
