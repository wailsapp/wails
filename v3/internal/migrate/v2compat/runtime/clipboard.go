package runtime

import (
	"context"
	"errors"
)

// ClipboardGetText mirrors the v2 runtime.ClipboardGetText function.
// v3 equivalent: app.Clipboard.Text.
func ClipboardGetText(_ context.Context) (string, error) {
	a := app()
	if a == nil {
		return "", errNoApp
	}
	text, ok := a.Clipboard.Text()
	if !ok {
		return "", errors.New("no text on clipboard")
	}
	return text, nil
}

// ClipboardSetText mirrors the v2 runtime.ClipboardSetText function.
// v3 equivalent: app.Clipboard.SetText.
func ClipboardSetText(_ context.Context, text string) error {
	a := app()
	if a == nil {
		return errNoApp
	}
	if !a.Clipboard.SetText(text) {
		return errors.New("failed to set clipboard text")
	}
	return nil
}
