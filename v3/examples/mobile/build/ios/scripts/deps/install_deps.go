// install_deps.go - iOS development dependency checker
// This script checks for required iOS development tools.
// It's designed to be portable across different shells by using Go instead of shell scripts.
//
// Usage:
//   go run install_deps.go                      # Interactive mode
//   TASK_FORCE_YES=true go run install_deps.go  # Auto-accept prompts
//   CI=true go run install_deps.go              # CI mode (auto-accept)

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Dependency struct {
	Name       string
	CheckFunc  func() (bool, string) // Returns (success, details)
	Required   bool
	InstallCmd []string
	InstallMsg string
	SuccessMsg string
	FailureMsg string
}

func main() {
	fmt.Println("Checking iOS development dependencies...")
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println()

	hasErrors := false
	dependencies := []Dependency{
		{
			Name: "Xcode",
			CheckFunc: func() (bool, string) {
				// Check if xcodebuild exists
				if !checkCommand([]string{"xcodebuild", "-version"}) {
					return false, ""
				}
				// Get version info
				out, err := exec.Command("xcodebuild", "-version").Output()
				if err != nil {
					return false, ""
				}
				lines := strings.Split(string(out), "\n")
				if len(lines) > 0 {
					return true, strings.TrimSpace(lines[0])
				}
				return true, ""
			},
			Required:   true,
			InstallMsg: "Please install Xcode from the Mac App Store:\n   https://apps.apple.com/app/xcode/id497799835\n   Xcode is REQUIRED for iOS development (includes iOS SDKs, simulators, and frameworks)",
			SuccessMsg: "✅ Xcode found",
			FailureMsg: "❌ Xcode not found (REQUIRED)",
		},
		{
			Name: "Xcode Developer Path",
			CheckFunc: func() (bool, string) {
				// Check if xcode-select points to a valid Xcode path
				out, err := exec.Command("xcode-select", "-p").Output()
				if err != nil {
					return false, "xcode-select not configured"
				}
				path := strings.TrimSpace(string(out))

				// Check if path exists and is in Xcode.app
				if _, err := os.Stat(path); err != nil {
					return false, "Invalid Xcode path"
				}

				// Verify it's pointing to Xcode.app (not just Command Line Tools)
				if !strings.Contains(path, "Xcode.app") {
					return false, fmt.Sprintf("Points to %s (should be Xcode.app)", path)
				}

				return true, path
			},
			Required:   true,
			InstallCmd: []string{"sudo", "xcode-select", "-s", "/Applications/Xcode.app/Contents/Developer"},
			InstallMsg: "Xcode developer path needs to be configured",
			SuccessMsg: "✅ Xcode developer path configured",
			FailureMsg: "❌ Xcode developer path not configured correctly",
		},
		{
			Name: "iOS SDK",
			CheckFunc: func() (bool, string) {
				// Get the iOS Simulator SDK path
				cmd := exec.Command("xcrun", "--sdk", "iphonesimulator", "--show-sdk-path")
				output, err := cmd.Output()
				if err != nil {
					return false, "Cannot find iOS SDK"
				}
				sdkPath := strings.TrimSpace(string(output))

				// Check if the SDK path exists
				if _, err := os.Stat(sdkPath); err != nil {
					return false, "iOS SDK path not found"
				}

				// Check for UIKit framework (essential for iOS development)
				uikitPath := fmt.Sprintf("%s/System/Library/Frameworks/UIKit.framework", sdkPath)
				if _, err := os.Stat(uikitPath); err != nil {
					return false, "UIKit.framework not found"
				}

				// Get SDK version
				versionCmd := exec.Command("xcrun", "--sdk", "iphonesimulator", "--show-sdk-version")
				versionOut, _ := versionCmd.Output()
				version := strings.TrimSpace(string(versionOut))

				return true, fmt.Sprintf("iOS %s SDK", version)
			},
			Required:   true,
			InstallMsg: "iOS SDK comes with Xcode. Please ensure Xcode is properly installed.",
			SuccessMsg: "✅ iOS SDK found with UIKit framework",
			FailureMsg: "❌ iOS SDK not found or incomplete",
		},
		{
			Name: "iOS Simulator Runtime",
			CheckFunc: func() (bool, string) {
				if !checkCommand([]string{"xcrun", "simctl", "help"}) {
					return false, ""
				}
				// Check if we can list runtimes
				out, err := exec.Command("xcrun", "simctl", "list", "runtimes").Output()
				if err != nil {
					return false, "Cannot access simulator"
				}
				// Count iOS runtimes
				lines := strings.Split(string(out), "\n")
				count := 0
				var versions []string
				for _, line := range lines {
					if strings.Contains(line, "iOS") && !strings.Contains(line, "unavailable") {
						count++
						// Extract version number
						if parts := strings.Fields(line); len(parts) > 2 {
							for _, part := range parts {
								if strings.HasPrefix(part, "(") && strings.HasSuffix(part, ")") {
									versions = append(versions, strings.Trim(part, "()"))
									break
								}
							}
						}
					}
				}
				if count > 0 {
					return true, fmt.Sprintf("%d runtime(s): %s", count, strings.Join(versions, ", "))
				}
				return false, "No iOS runtimes installed"
			},
			Required:   true,
			InstallMsg: "iOS Simulator runtimes come with Xcode. You may need to download them:\n   Xcode → Settings → Platforms → iOS",
			SuccessMsg: "✅ iOS Simulator runtime available",
			FailureMsg: "❌ iOS Simulator runtime not available",
		},
	}

	// Check each dependency
	for _, dep := range dependencies {
		success, details := dep.CheckFunc()
		if success {
			msg := dep.SuccessMsg
			if details != "" {
				msg = fmt.Sprintf("%s (%s)", dep.SuccessMsg, details)
			}
			fmt.Println(msg)
		} else {
			fmt.Println(dep.FailureMsg)
			if details != "" {
				fmt.Printf("   Details: %s\n", details)
			}
			if dep.Required {
				hasErrors = true
				if len(dep.InstallCmd) > 0 {
					fmt.Println()
					fmt.Println("   " + dep.InstallMsg)
					fmt.Printf("   Fix command: %s\n", strings.Join(dep.InstallCmd, " "))
					if promptUser("Do you want to run this command?") {
						fmt.Println("Running command...")
						cmd := exec.Command(dep.InstallCmd[0], dep.InstallCmd[1:]...)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						cmd.Stdin = os.Stdin
						if err := cmd.Run(); err != nil {
							fmt.Printf("Command failed: %v\n", err)
							os.Exit(1)
						}
						fmt.Println("✅ Command completed. Please run this check again.")
					} else {
						fmt.Printf("   Please run manually: %s\n", strings.Join(dep.InstallCmd, " "))
					}
				} else {
					fmt.Println("   " + dep.InstallMsg)
				}
			}
		}
	}

	// Check for iPhone simulators
	fmt.Println()
	fmt.Println("Checking for iPhone simulator devices...")
	if !checkCommand([]string{"xcrun", "simctl", "list", "devices"}) {
		fmt.Println("❌ Cannot check for iPhone simulators")
		hasErrors = true
	} else {
		out, err := exec.Command("xcrun", "simctl", "list", "devices").Output()
		if err != nil {
			fmt.Println("❌ Failed to list simulator devices")
			hasErrors = true
		} else if !strings.Contains(string(out), "iPhone") {
			fmt.Println("⚠️  No iPhone simulator devices found")
			fmt.Println()

			// Get the latest iOS runtime
			runtimeOut, err := exec.Command("xcrun", "simctl", "list", "runtimes").Output()
			if err != nil {
				fmt.Println("   Failed to get iOS runtimes:", err)
			} else {
				lines := strings.Split(string(runtimeOut), "\n")
				var latestRuntime string
				for _, line := range lines {
					if strings.Contains(line, "iOS") && !strings.Contains(line, "unavailable") {
						// Extract runtime identifier
						parts := strings.Fields(line)
						if len(parts) > 0 {
							latestRuntime = parts[len(parts)-1]
						}
					}
				}

				if latestRuntime == "" {
					fmt.Println("   No iOS runtime found. Please install iOS simulators in Xcode:")
					fmt.Println("   Xcode → Settings → Platforms → iOS")
				} else {
					fmt.Println("   Would you like to create an iPhone 15 Pro simulator?")
					createCmd := []string{"xcrun", "simctl", "create", "iPhone 15 Pro", "iPhone 15 Pro", latestRuntime}
					fmt.Printf("   Command: %s\n", strings.Join(createCmd, " "))
					if promptUser("Create simulator?") {
						cmd := exec.Command(createCmd[0], createCmd[1:]...)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						if err := cmd.Run(); err != nil {
							fmt.Printf("   Failed to create simulator: %v\n", err)
						} else {
							fmt.Println("   ✅ iPhone 15 Pro simulator created")
						}
					} else {
						fmt.Println("   Skipping simulator creation")
						fmt.Printf("   Create manually: %s\n", strings.Join(createCmd, " "))
					}
				}
			}
		} else {
			// Count iPhone devices
			count := 0
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if strings.Contains(line, "iPhone") && !strings.Contains(line, "unavailable") {
					count++
				}
			}
			fmt.Printf("✅ %d iPhone simulator device(s) available\n", count)
		}
	}

	// Final summary
	fmt.Println()
	fmt.Println("=" + strings.Repeat("=", 50))
	if hasErrors {
		fmt.Println("❌ Some required dependencies are missing or misconfigured.")
		fmt.Println()
		fmt.Println("Quick setup guide:")
		fmt.Println("1. Install Xcode from Mac App Store (if not installed)")
		fmt.Println("2. Open Xcode once and agree to the license")
		fmt.Println("3. Install additional components when prompted")
		fmt.Println("4. Run: sudo xcode-select -s /Applications/Xcode.app/Contents/Developer")
		fmt.Println("5. Download iOS simulators: Xcode → Settings → Platforms → iOS")
		fmt.Println("6. Run this check again")
		os.Exit(1)
	} else {
		fmt.Println("✅ All required dependencies are installed!")
		fmt.Println("   You're ready for iOS development with Wails!")
	}
}

func checkCommand(args []string) bool {
	if len(args) == 0 {
		return false
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Run()
	return err == nil
}

func promptUser(question string) bool {
	// Check if we're in a non-interactive environment
	if os.Getenv("CI") != "" || os.Getenv("TASK_FORCE_YES") == "true" {
		fmt.Printf("%s [y/N]: y (auto-accepted)\n", question)
		return true
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", question)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}