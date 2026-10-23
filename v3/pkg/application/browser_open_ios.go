//go:build ios && !server

package application

import "errors"

// internal/browser has no iOS backend, so URLs are opened through UIApplication.
func openBrowserURL(url string) error {
	IOS.OpenURL(url)
	return nil
}

func openBrowserFile(_ string) error {
	return errors.New("OpenFile is not supported on iOS")
}
