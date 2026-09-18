package setupwizard

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// checkMobileDependencies returns dependency statuses for the requested mobile
// platforms. iOS checks are macOS-only; on other hosts they report that iOS
// builds require a Mac (see wizard_mobile_other.go).
func (w *Wizard) checkMobileDependencies(ios, android bool) []DependencyStatus {
	var deps []DependencyStatus
	if ios {
		deps = append(deps, checkXcodeApp(), checkIOSRuntime())
	}
	if android {
		deps = append(deps, checkJDK(), checkAndroidSDK(), checkAndroidNDK(), checkAndroidEmulator())
	}
	return deps
}

// execCombined runs a command and returns its combined stdout+stderr. Some tools
// (notably `java -version`) print to stderr.
func execCombined(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// androidSDKRoot returns the configured Android SDK location, or the OS default.
func androidSDKRoot() string {
	for _, env := range []string{"ANDROID_HOME", "ANDROID_SDK_ROOT"} {
		if v := strings.TrimSpace(os.Getenv(env)); v != "" {
			return v
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Android", "sdk")
	case "windows":
		return filepath.Join(home, "AppData", "Local", "Android", "Sdk")
	default:
		return filepath.Join(home, "Android", "Sdk")
	}
}

func dirExists(p string) bool {
	if p == "" {
		return false
	}
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

// hasSubdir reports whether parent contains at least one subdirectory.
func hasSubdir(parent string) bool {
	entries, err := os.ReadDir(parent)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() {
			return true
		}
	}
	return false
}

// jdkInstallCommand returns an auto-install command for a JDK, or "" if none is
// available for this host.
func jdkInstallCommand() string {
	switch runtime.GOOS {
	case "darwin":
		if commandExists("brew") {
			return "brew install openjdk@21"
		}
	case "linux":
		switch {
		case commandExists("apt-get"):
			return "sudo apt-get install -y openjdk-21-jdk"
		case commandExists("dnf"):
			return "sudo dnf install -y java-21-openjdk-devel"
		case commandExists("pacman"):
			return "sudo pacman -S --noconfirm jdk-openjdk"
		}
	}
	return ""
}

// findJDKHome locates a JDK that may be installed but not on PATH (e.g. a
// keg-only Homebrew openjdk). Returns the JAVA_HOME-style dir, or "".
func findJDKHome() string {
	if jh := strings.TrimSpace(os.Getenv("JAVA_HOME")); jh != "" && fileExists(filepath.Join(jh, "bin", "java")) {
		return jh
	}
	if runtime.GOOS == "darwin" {
		if out, err := execCommand("/usr/libexec/java_home", "-v", "17+"); err == nil && out != "" && fileExists(filepath.Join(out, "bin", "java")) {
			return out
		}
	}
	globs := []string{
		"/opt/homebrew/opt/openjdk*/libexec/openjdk.jdk/Contents/Home",
		"/opt/homebrew/Cellar/openjdk*/*/libexec/openjdk.jdk/Contents/Home",
		"/usr/local/opt/openjdk*/libexec/openjdk.jdk/Contents/Home",
		"/usr/local/Cellar/openjdk*/*/libexec/openjdk.jdk/Contents/Home",
		"/Library/Java/JavaVirtualMachines/*/Contents/Home",
		"/usr/lib/jvm/*",
	}
	for _, g := range globs {
		matches, _ := filepath.Glob(g)
		for _, m := range matches {
			if fileExists(filepath.Join(m, "bin", "java")) {
				return m
			}
		}
	}
	return ""
}

// inShellConfig reports whether the user's shell startup files contain the
// given assignment (e.g. "ANDROID_HOME="). The running wizard can't see env
// changes made after it launched, so to make "Re-check" work after someone
// pastes our export lines, we read the rc files directly — that's exactly where
// we told them to add it, and it's what future shells (and builds) will use.
func inShellConfig(assignment string) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	for _, f := range []string{".zshrc", ".zprofile", ".zshenv", ".bashrc", ".bash_profile", ".profile"} {
		data, err := os.ReadFile(filepath.Join(home, f))
		if err != nil {
			continue
		}
		if strings.Contains(string(data), assignment) {
			return true
		}
	}
	return false
}

func checkJDK() DependencyStatus {
	dep := DependencyStatus{Name: "Java (JDK 17+)", Required: true, HelpURL: "https://adoptium.net/"}

	// 1. On PATH — fully ready.
	if out, err := execCombined("java", "-version"); err == nil {
		// e.g. openjdk version "21.0.1" 2023-10-17
		if i := strings.Index(out, "\""); i >= 0 {
			if j := strings.Index(out[i+1:], "\""); j >= 0 {
				dep.Version = out[i+1 : i+1+j]
			}
		}
		dep.Installed = true
		dep.Status = "installed"
		if major := leadingInt(dep.Version); major != 0 && major < 17 {
			dep.Status = "needs_update"
			dep.Message = "JDK 17 or newer is recommended for Android builds (found " + dep.Version + ")"
			dep.InstallCommand = jdkInstallCommand()
		}
		return dep
	}

	// 2. Installed but not on this process's PATH (e.g. keg-only Homebrew openjdk).
	if home := findJDKHome(); home != "" {
		dep.Installed = true
		// If they've already added JAVA_HOME to their shell config, future shells
		// (and builds) will find it — treat it as configured even though THIS
		// process's stale env can't see it. Makes "Re-check" work after pasting.
		if inShellConfig("JAVA_HOME=") {
			dep.Status = "installed"
			if out, err := execCombined(filepath.Join(home, "bin", "java"), "-version"); err == nil {
				if i := strings.Index(out, "\""); i >= 0 {
					if j := strings.Index(out[i+1:], "\""); j >= 0 {
						dep.Version = out[i+1 : i+1+j]
					}
				}
			}
			return dep
		}
		// Otherwise tell them how to expose it (no reinstall needed).
		dep.Status = "needs_config"
		dep.HelpURL = "" // the copyable config below is the action, not a doc link
		dep.Message = "A JDK is installed but isn't on your PATH. Add these to your shell config (~/.zshrc or ~/.bashrc), then Re-check:"
		dep.ConfigCommand = "export JAVA_HOME=\"" + home + "\"\nexport PATH=\"$JAVA_HOME/bin:$PATH\""
		return dep
	}

	// 3. Genuinely missing.
	dep.Status = "not_installed"
	dep.Message = "A JDK (17 or newer) is required to build for Android"
	dep.InstallCommand = jdkInstallCommand()
	return dep
}

func checkAndroidSDK() DependencyStatus {
	dep := DependencyStatus{Name: "Android SDK", Required: true, HelpURL: "https://developer.android.com/studio"}
	root := androidSDKRoot()

	// The SDK is usable once it has the command-line tools (which provide
	// sdkmanager) or platform-tools.
	if root != "" && (dirExists(filepath.Join(root, "cmdline-tools")) || dirExists(filepath.Join(root, "platform-tools"))) {
		dep.Installed = true
		// Exported in this process, or added to the shell config (so future
		// shells/builds will have it — the running wizard can't see live env
		// changes, so check the rc files to make Re-check work after pasting).
		if os.Getenv("ANDROID_HOME") != "" || os.Getenv("ANDROID_SDK_ROOT") != "" || inShellConfig("ANDROID_HOME=") {
			dep.Status = "installed"
			dep.Version = root
		} else {
			// Present on disk but not exported — builds that rely on ANDROID_HOME
			// will fail, so flag it as needs-config rather than a green check.
			dep.Status = "needs_config"
			dep.HelpURL = "" // the copyable config below is the action, not a doc link
			dep.Message = "SDK found, but ANDROID_HOME isn't set. Add these to your shell config (~/.zshrc or ~/.bashrc), then Re-check:"
			dep.ConfigCommand = "export ANDROID_HOME=\"" + root + "\"\nexport PATH=\"$ANDROID_HOME/platform-tools:$PATH\""
		}
		return dep
	}

	dep.Status = "not_installed"
	dep.Message = "Android SDK command-line tools are required; then set ANDROID_HOME"
	if runtime.GOOS == "darwin" && commandExists("brew") {
		dep.InstallCommand = "brew install --cask android-commandlinetools"
	} else {
		dep.HelpURL = "https://developer.android.com/studio"
		dep.HelpLabel = "Download Android Studio for " + osLabel()
	}
	return dep
}

// osLabel returns a friendly OS name for help text.
func osLabel() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS"
	case "windows":
		return "Windows"
	default:
		return "Linux"
	}
}

func checkAndroidNDK() DependencyStatus {
	dep := DependencyStatus{Name: "Android NDK", Required: true, HelpURL: "https://developer.android.com/ndk"}
	root := androidSDKRoot()

	if root != "" && dirExists(filepath.Join(root, "ndk")) && hasSubdir(filepath.Join(root, "ndk")) {
		dep.Installed = true
		dep.Status = "installed"
		return dep
	}

	dep.Status = "not_installed"
	dep.Message = "The NDK is required to compile Go for Android"
	// Auto-installable once sdkmanager (from the cmdline-tools) is on PATH.
	if commandExists("sdkmanager") {
		dep.InstallCommand = "sdkmanager --install ndk;26.3.11579264"
	} else {
		dep.HelpURL = "https://developer.android.com/studio/projects/install-ndk"
		dep.HelpLabel = "Install via Android Studio's SDK Manager"
	}
	return dep
}

func checkAndroidEmulator() DependencyStatus {
	dep := DependencyStatus{Name: "Android Emulator", Required: false, HelpURL: "https://developer.android.com/studio/run/managing-avds"}
	root := androidSDKRoot()
	emulatorBin := filepath.Join(root, "emulator", "emulator")

	if root == "" || !fileExists(emulatorBin) {
		dep.Status = "not_installed"
		dep.Message = "Optional — needed to run your app on a virtual device"
		if commandExists("sdkmanager") {
			dep.InstallCommand = "sdkmanager --install emulator platform-tools"
		} else {
			dep.HelpLabel = "Install via Android Studio's SDK Manager"
		}
		return dep
	}

	// Emulator is installed; check for at least one AVD.
	avds, _ := execCommand(emulatorBin, "-list-avds")
	if strings.TrimSpace(avds) == "" {
		dep.Installed = true
		dep.Status = "needs_update"
		dep.Message = "Installed, but no virtual device (AVD) exists yet — create one in Android Studio or with avdmanager"
		return dep
	}

	dep.Installed = true
	dep.Status = "installed"
	dep.Message = "Ready"
	return dep
}

// leadingInt parses the leading integer of s (e.g. "21.0.1" -> 21).
func leadingInt(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}
