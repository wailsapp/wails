//go:build !darwin

package setupwizard

// On non-macOS hosts, iOS builds aren't possible. These report that clearly so
// the wizard can surface it rather than silently omitting iOS.

func checkXcodeApp() DependencyStatus {
	return DependencyStatus{
		Name:     "Xcode (iOS)",
		Required: false,
		Status:   "not_installed",
		Message:  "iOS apps can only be built on macOS with Xcode",
		HelpURL:  "https://v3.wails.io/guides/mobile/ios",
	}
}

func checkIOSRuntime() DependencyStatus {
	return DependencyStatus{
		Name:     "iOS Simulator Runtime",
		Required: false,
		Status:   "not_installed",
		Message:  "iOS development requires macOS",
		HelpURL:  "https://v3.wails.io/guides/mobile/ios",
	}
}
