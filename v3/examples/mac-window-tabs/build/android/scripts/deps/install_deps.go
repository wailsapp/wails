package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	fmt.Println("Checking Android development dependencies...")
	fmt.Println()

	errors := []string{}

	// Check Go
	if !checkCommand("go", "version") {
		errors = append(errors, "Go is not installed. Install from https://go.dev/dl/")
	} else {
		fmt.Println("✓ Go is installed")
	}

	// Check ANDROID_HOME
	androidHome := os.Getenv("ANDROID_HOME")
	if androidHome == "" {
		androidHome = os.Getenv("ANDROID_SDK_ROOT")
	}
	if androidHome == "" {
		// Try common default locations
		home, _ := os.UserHomeDir()
		possiblePaths := []string{
			filepath.Join(home, "Android", "Sdk"),
			filepath.Join(home, "Library", "Android", "sdk"),
			"/usr/local/share/android-sdk",
		}
		for _, p := range possiblePaths {
			if _, err := os.Stat(p); err == nil {
				androidHome = p
				break
			}
		}
	}

	if androidHome == "" {
		errors = append(errors, "ANDROID_HOME not set. Install Android Studio and set ANDROID_HOME environment variable")
	} else {
		fmt.Printf("✓ ANDROID_HOME: %s\n", androidHome)
	}

	// Check adb
	if !checkCommand("adb", "version") {
		if androidHome != "" {
			platformTools := filepath.Join(androidHome, "platform-tools")
			errors = append(errors, fmt.Sprintf("adb not found. Add %s to PATH", platformTools))
		} else {
			errors = append(errors, "adb not found. Install Android SDK Platform-Tools")
		}
	} else {
		fmt.Println("✓ adb is installed")
	}

	// Check emulator
	if !checkCommand("emulator", "-list-avds") {
		if androidHome != "" {
			emulatorPath := filepath.Join(androidHome, "emulator")
			errors = append(errors, fmt.Sprintf("emulator not found. Add %s to PATH", emulatorPath))
		} else {
			errors = append(errors, "emulator not found. Install Android Emulator via SDK Manager")
		}
	} else {
		fmt.Println("✓ Android Emulator is installed")
	}

	// Check NDK
	ndkHome := os.Getenv("ANDROID_NDK_HOME")
	if ndkHome == "" && androidHome != "" {
		// Look for NDK in default location
		ndkDir := filepath.Join(androidHome, "ndk")
		if entries, err := os.ReadDir(ndkDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					ndkHome = filepath.Join(ndkDir, entry.Name())
					break
				}
			}
		}
	}

	if ndkHome == "" {
		errors = append(errors, "Android NDK not found. Install NDK via Android Studio > SDK Manager > SDK Tools > NDK (Side by side)")
	} else {
		fmt.Printf("✓ Android NDK: %s\n", ndkHome)
	}

	// Check Java
	if !checkCommand("java", "-version") {
		errors = append(errors, "Java not found. Install JDK 11+ (OpenJDK recommended)")
	} else {
		fmt.Println("✓ Java is installed")
	}

	// Check for AVD (Android Virtual Device)
	if checkCommand("emulator", "-list-avds") {
		cmd := exec.Command("emulator", "-list-avds")
		output, err := cmd.Output()
		if err == nil && len(strings.TrimSpace(string(output))) > 0 {
			avds := strings.Split(strings.TrimSpace(string(output)), "\n")
			fmt.Printf("✓ Found %d Android Virtual Device(s)\n", len(avds))
		} else {
			fmt.Println("⚠ No Android Virtual Devices found. Create one via Android Studio > Tools > Device Manager")
		}
	}

	fmt.Println()

	if len(errors) > 0 {
		fmt.Println("❌ Missing dependencies:")
		for _, err := range errors {
			fmt.Printf("   - %s\n", err)
		}
		fmt.Println()
		fmt.Println("Setup instructions:")
		fmt.Println("1. Install Android Studio: https://developer.android.com/studio")
		fmt.Println("2. Open SDK Manager and install:")
		fmt.Println("   - Android SDK Platform (API 34)")
		fmt.Println("   - Android SDK Build-Tools")
		fmt.Println("   - Android SDK Platform-Tools")
		fmt.Println("   - Android Emulator")
		fmt.Println("   - NDK (Side by side)")
		fmt.Println("3. Set environment variables:")
		if runtime.GOOS == "darwin" {
			fmt.Println("   export ANDROID_HOME=$HOME/Library/Android/sdk")
		} else {
			fmt.Println("   export ANDROID_HOME=$HOME/Android/Sdk")
		}
		fmt.Println("   export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/emulator")
		fmt.Println("4. Create an AVD via Android Studio > Tools > Device Manager")
		os.Exit(1)
	}

	fmt.Println("✓ All Android development dependencies are installed!")
}

func checkCommand(name string, args ...string) bool {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}
