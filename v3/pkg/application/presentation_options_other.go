//go:build !darwin || ios || server

package application

// Off macOS the presentation options are documented no-ops:
// App.SetPresentationOptions returns ErrMacOnly and PresentationOptions
// reports MacPresentationDefault.

func macPresentationOptionsSupported() bool { return false }

func macSetPresentationOptions(MacPresentationOptions) error { return ErrMacOnly }

func macPresentationOptions() MacPresentationOptions { return MacPresentationDefault }
