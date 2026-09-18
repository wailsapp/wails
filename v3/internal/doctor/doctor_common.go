package doctor

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

func checkCommonDependencies(result map[string]string, ok *bool) {
	// Check for npm
	npmVersion := []byte("Not Installed. Requires npm >= 7.0.0")
	npmVersion, err := exec.Command("npm", "-v").Output()
	if err != nil {
		*ok = false
	} else {
		npmVersion = bytes.TrimSpace(npmVersion)
		// Check that it's at least version 7 by converting first byte to int and checking if it's >= 7
		// Parse the semver string
		semver := strings.Split(string(npmVersion), ".")
		if len(semver) > 0 {
			major, _ := strconv.Atoi(semver[0])
			if major < 7 {
				*ok = false
				npmVersion = append(npmVersion, []byte(". Installed, but requires npm >= 7.0.0")...)
			} else {
				*ok = true
			}
		}
	}
	result["npm"] = string(npmVersion)

	// Check for Docker (optional - used for macOS cross-compilation from Linux)
	checkDocker(result)

	// Android toolchain (optional - only needed for `wails3 task android:*`)
	checkAndroid(result)
}

// androidSDKRoot resolves the Android SDK location from the standard
// environment variables, falling back to the conventional install path.
func androidSDKRoot() string {
	for _, env := range []string{"ANDROID_HOME", "ANDROID_SDK_ROOT"} {
		if v := os.Getenv(env); v != "" {
			return v
		}
	}
	if home, err := os.UserHomeDir(); err == nil {
		// Conventional per-OS install locations
		candidates := []string{
			filepath.Join(home, "Library", "Android", "sdk"), // macOS
			filepath.Join(home, "Android", "Sdk"),            // Linux
		}
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			candidates = append(candidates, filepath.Join(localAppData, "Android", "Sdk")) // Windows
		}
		for _, def := range candidates {
			if _, err := os.Stat(def); err == nil {
				return def
			}
		}
	}
	return ""
}

// checkAndroid reports the Android SDK, NDK and Java state. All entries are
// optional (prefixed with "*"): they are only required for Android builds.
func checkAndroid(result map[string]string) {
	sdk := androidSDKRoot()
	if sdk == "" {
		result["*Android SDK"] = "Not found. Set ANDROID_HOME (install via Android Studio or the command-line tools)."
		return
	}
	result["*Android SDK"] = sdk

	// adb (platform-tools); the binary is adb.exe on Windows
	adbName := "adb"
	if runtime.GOOS == "windows" {
		adbName = "adb.exe"
	}
	adb := filepath.Join(sdk, "platform-tools", adbName)
	if _, err := os.Stat(adb); err == nil {
		result["*Android platform-tools"] = "Installed"
	} else if _, err := exec.LookPath("adb"); err == nil {
		result["*Android platform-tools"] = "Installed (on PATH)"
	} else {
		result["*Android platform-tools"] = "Not found. Install with: sdkmanager 'platform-tools'"
	}

	// NDK: prefer ANDROID_NDK_HOME, else newest under $SDK/ndk
	ndk := os.Getenv("ANDROID_NDK_HOME")
	if ndk == "" {
		if entries, err := os.ReadDir(filepath.Join(sdk, "ndk")); err == nil {
			var versions []string
			for _, e := range entries {
				if e.IsDir() {
					versions = append(versions, e.Name())
				}
			}
			if len(versions) > 0 {
				sort.Strings(versions)
				ndk = filepath.Join(sdk, "ndk", versions[len(versions)-1])
			}
		}
	}
	if ndk != "" {
		result["*Android NDK"] = ndk
	} else {
		result["*Android NDK"] = "Not found. Install with: sdkmanager 'ndk;26.3.11579264'"
	}

	// Java (required by Gradle)
	javaVersion := "Not found. Install a JDK (e.g. brew install openjdk@21) or set JAVA_HOME."
	if out, err := exec.Command("java", "-version").CombinedOutput(); err == nil {
		line := strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)[0]
		javaVersion = strings.TrimSpace(line)
	} else if javaHome := os.Getenv("JAVA_HOME"); javaHome != "" {
		if out, err := exec.Command(filepath.Join(javaHome, "bin", "java"), "-version").CombinedOutput(); err == nil {
			javaVersion = strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)[0]
		}
	}
	result["*Java (Android)"] = javaVersion
}

func checkDocker(result map[string]string) {
	dockerVersion, err := exec.Command("docker", "--version").Output()
	if err != nil {
		result["docker"] = "*Not installed (optional - for cross-compilation)"
		return
	}

	// Check if Docker daemon is running
	_, err = exec.Command("docker", "info").Output()
	if err != nil {
		version := strings.TrimSpace(string(dockerVersion))
		result["docker"] = "*" + version + " (daemon not running)"
		return
	}

	version := strings.TrimSpace(string(dockerVersion))

	// Check for the unified cross-compilation image
	imageCheck, _ := exec.Command("docker", "image", "inspect", "wails-cross").Output()
	if len(imageCheck) == 0 {
		result["docker"] = "*" + version + " (wails-cross image not built - run: wails3 task setup:docker)"
	} else {
		result["docker"] = "*" + version + " (cross-compilation ready)"
	}
}
