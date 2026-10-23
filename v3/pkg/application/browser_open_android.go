//go:build android && !server

package application

import "errors"

// internal/browser shells out to xdg-open, which does not exist on Android,
// so URLs are opened through the Java bridge (Intent.ACTION_VIEW) instead.
func openBrowserURL(url string) error {
	Android.OpenURL(url)
	return nil
}

func openBrowserFile(_ string) error {
	return errors.New("OpenFile is not supported on Android")
}
