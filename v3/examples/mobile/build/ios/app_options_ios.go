//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

// modifyOptionsForIOS adjusts the application options for iOS
func modifyOptionsForIOS(opts *application.Options) {
	// Disable signal handlers on iOS to prevent crashes
	opts.DisableDefaultSignalHandler = true

	// Enable native UITabBar in the iOS example by default
	opts.IOS.EnableNativeTabs = true
	// Configure example tab items (titles + SF Symbols)
	opts.IOS.NativeTabsItems = []application.NativeTabItem{
		{Title: "Bindings", SystemImage: "link"},
		{Title: "Go Runtime", SystemImage: "gearshape"},
		{Title: "JS Runtime", SystemImage: "chevron.left.slash.chevron.right"},
	}
}