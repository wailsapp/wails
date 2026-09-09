//go:build darwin

package setupwizard

import "strings"

// checkXcodeApp verifies that the full Xcode app (not just the Command Line
// Tools) is installed and selected — required for iOS builds and simulators.
func checkXcodeApp() DependencyStatus {
	dep := DependencyStatus{
		Name:     "Xcode",
		Required: true,
		HelpURL:  "https://apps.apple.com/app/xcode/id497799835",
	}

	path, err := execCommand("xcode-select", "-p")
	if err != nil || strings.Contains(path, "CommandLineTools") {
		dep.Status = "not_installed"
		dep.Message = "Full Xcode is required for iOS (the Command Line Tools alone aren't enough). After installing, run: sudo xcode-select -s /Applications/Xcode.app"
		dep.HelpLabel = "Get Xcode from the App Store"
		return dep
	}

	ver, err := execCommand("xcodebuild", "-version")
	if err != nil {
		dep.Status = "not_installed"
		dep.Message = "Xcode is selected but xcodebuild is unavailable — open Xcode once to finish installing its components"
		return dep
	}

	dep.Installed = true
	dep.Status = "installed"
	if line := strings.SplitN(ver, "\n", 2)[0]; strings.HasPrefix(line, "Xcode") {
		dep.Version = strings.TrimSpace(strings.TrimPrefix(line, "Xcode"))
	}
	return dep
}

// checkIOSRuntime verifies that at least one iOS simulator runtime is installed.
func checkIOSRuntime() DependencyStatus {
	dep := DependencyStatus{
		Name:     "iOS Simulator Runtime",
		Required: true,
		HelpURL:  "https://developer.apple.com/documentation/xcode/installing-additional-simulator-runtimes",
	}

	out, err := execCommand("xcrun", "simctl", "list", "runtimes")
	if err == nil && strings.Contains(out, "SimRuntime.iOS") {
		dep.Installed = true
		dep.Status = "installed"
		return dep
	}

	dep.Status = "not_installed"
	dep.Message = "No iOS simulator runtime is installed"
	dep.InstallCommand = "xcodebuild -downloadPlatform iOS"
	return dep
}
