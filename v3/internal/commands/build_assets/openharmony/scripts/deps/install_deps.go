// install_deps.go - OpenHarmony development dependency checker.
// Reports the tools and SDK paths required by the OpenHarmony Taskfile tasks.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	fmt.Println("Checking OpenHarmony development dependencies...")
	fmt.Println()

	problems := []string{}

	if _, err := exec.LookPath("go"); err != nil {
		problems = append(problems, "Go is not installed. Install from https://go.dev/dl/ (an OpenHarmony-enabled toolchain is required for cross-compilation)")
	} else {
		fmt.Println("✓ Go is installed")
	}

	sdk := os.Getenv("DEVECO_SDK_HOME")
	if sdk == "" {
		problems = append(problems, "DEVECO_SDK_HOME is not set. Install DevEco Studio and point it at its sdk directory")
	} else {
		fmt.Printf("✓ DEVECO_SDK_HOME: %s\n", sdk)

		native := filepath.Join(sdk, "default", "openharmony", "native")
		clang := filepath.Join(native, "llvm", "bin", "clang")
		if runtime.GOOS == "windows" {
			clang += ".exe"
		}
		if info, err := os.Stat(clang); err != nil || info.IsDir() {
			problems = append(problems, fmt.Sprintf("OpenHarmony clang not found at %s", clang))
		} else {
			fmt.Println("✓ OpenHarmony clang is installed")
		}

		sysroot := filepath.Join(native, "sysroot")
		if info, err := os.Stat(sysroot); err != nil || !info.IsDir() {
			problems = append(problems, fmt.Sprintf("OpenHarmony sysroot not found at %s", sysroot))
		} else {
			fmt.Println("✓ OpenHarmony sysroot is installed")
		}

		hvigorw := filepath.Join(filepath.Dir(filepath.Clean(sdk)), "tools", "hvigor", "bin", "hvigorw")
		if runtime.GOOS == "windows" {
			hvigorw += ".bat"
		}
		if info, err := os.Stat(hvigorw); err != nil || info.IsDir() {
			problems = append(problems, fmt.Sprintf("hvigorw not found at %s (DevEco Studio tools are missing)", hvigorw))
		} else {
			fmt.Println("✓ hvigorw is installed")
		}
	}

	fmt.Println()
	if len(problems) > 0 {
		fmt.Println("❌ Missing dependencies:")
		for _, problem := range problems {
			fmt.Printf("   - %s\n", problem)
		}
		fmt.Println()
		fmt.Println("Install DevEco Studio, then set DEVECO_SDK_HOME to its sdk directory.")
		os.Exit(1)
	}

	fmt.Println("✓ All OpenHarmony development dependencies are installed!")
}
