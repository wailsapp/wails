//go:build ignore

// verify-manifest-build-system runs real host toolchains against an already
// migrated disposable Wails project. It is intentionally separate from unit
// tests because native packaging is slow and host/tool dependent.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/wake/benchmark"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

type acceptanceCommand struct {
	args        []string
	verifyCache bool
}

func main() {
	wails := flag.String("wails", "wails3", "Wails CLI executable")
	project := flag.String("project", ".", "migrated disposable project")
	targetsFlag := flag.String("targets", "", "comma-separated host-platform targets; defaults to the current host target")
	appImage := flag.Bool("appimage", false, "include the slow/networked Linux AppImage acceptance run")
	android := flag.Bool("android", false, "include Android amd64, arm64 and universal AAB runs when the SDK/NDK are installed")
	ios := flag.Bool("ios", false, "include iOS simulator app packaging on macOS")
	iosDevice := flag.Bool("ios-device", false, "include iOS device IPA packaging on macOS with a device target profile")
	iosDeviceProfile := flag.String("ios-device-profile", "ios-device", "complete HCL profile used by -ios-device; it must select ios/arm64, destination=device, IPA and signing")
	sign := flag.Bool("sign", false, "run host signing with credentials already declared in wails.hcl")
	verifyCache := flag.Bool("verify-cache", true, "rerun cacheable commands and require a zero-work result")
	timeout := flag.Duration("timeout", 30*time.Minute, "complete verification timeout")
	flag.Parse()
	root, err := filepath.Abs(*project)
	if err != nil {
		fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, manifest.Filename)); err != nil {
		fatal(fmt.Errorf("%s is not a manifest project: %w", root, err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	if *ios && *iosDevice {
		fatal(fmt.Errorf("choose either simulator -ios or device -ios-device acceptance"))
	}
	targets := []string{runtime.GOOS + "/" + runtime.GOARCH}
	if *targetsFlag != "" {
		targets = nil
		for _, target := range strings.Split(*targetsFlag, ",") {
			target = strings.TrimSpace(target)
			if target != "" {
				targets = append(targets, target)
			}
		}
		if len(targets) == 0 {
			fatal(fmt.Errorf("-targets must contain at least one OS/architecture target"))
		}
	}
	commands := make([]acceptanceCommand, 0, len(targets)*3+4)
	hostFormats := ""
	switch runtime.GOOS {
	case "linux":
		hostFormats = "deb,rpm,archlinux"
		if *appImage {
			hostFormats += ",appimage"
		}
	case "windows":
		hostFormats = "nsis,msix"
	case "darwin":
		// The runnable .app is the base output from the preceding build. DMG is
		// the only explicit macOS package format.
		hostFormats = "dmg"
	default:
		fatal(fmt.Errorf("unsupported native acceptance host %s", runtime.GOOS))
	}
	for _, target := range targets {
		platform, _, ok := strings.Cut(target, "/")
		if !ok || platform != runtime.GOOS {
			fatal(fmt.Errorf("host acceptance target %q must use %s/<arch>", target, runtime.GOOS))
		}
		commands = append(commands,
			acceptanceCommand{args: []string{"build", "--targets", target}, verifyCache: true},
			acceptanceCommand{args: []string{"build", "--targets", target, "--formats", hostFormats}, verifyCache: true},
		)
		if *sign {
			commands = append(commands, acceptanceCommand{args: []string{"sign", "--targets", target, "--formats", hostFormats}})
		}
	}
	if *ios {
		if runtime.GOOS != "darwin" {
			fatal(fmt.Errorf("iOS acceptance requires macOS"))
		}
		// iOS app assembly invokes codesign, including for the simulator, and is
		// deliberately non-reusable. The base build already produces the .app;
		// "app" is not an explicit --formats value.
		commands = append(commands, acceptanceCommand{args: []string{"build", "--targets", "ios/arm64"}})
	}
	if *iosDevice {
		if runtime.GOOS != "darwin" {
			fatal(fmt.Errorf("iOS acceptance requires macOS"))
		}
		if strings.TrimSpace(*iosDeviceProfile) == "" {
			fatal(fmt.Errorf("-ios-device-profile must name a complete device profile"))
		}
		// Destination and signing are profile intent, not anonymous CLI
		// overrides. This run is non-reusable because it signs the device app.
		commands = append(commands, acceptanceCommand{args: []string{"build", *iosDeviceProfile}})
	}
	if *android {
		for _, target := range []string{"android/amd64", "android/arm64", "android/universal"} {
			commands = append(commands, acceptanceCommand{args: []string{"build", "--targets", target, "--formats", "aab"}, verifyCache: true})
		}
	}
	for _, item := range commands {
		fmt.Printf("\n==> %s %s\n", *wails, strings.Join(item.args, " "))
		command := exec.CommandContext(ctx, *wails, item.args...)
		command.Dir = root
		command.Env = os.Environ()
		command.Stdout, command.Stderr = os.Stdout, os.Stderr
		if err := command.Run(); err != nil {
			fatal(err)
		}
		verifyReceipt(root)
		if *verifyCache && item.verifyCache {
			verifyWarmCache(ctx, root, *wails, item.args)
		}
	}
	fmt.Printf("\nmanifest acceptance passed on %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func verifyWarmCache(ctx context.Context, root, wails string, args []string) {
	command := append([]string{wails}, args...)
	result, err := benchmark.Run(ctx, benchmark.Config{
		Scenario:         "native-warm-" + args[0],
		Command:          command,
		WorkingDirectory: root,
		Environment:      os.Environ(),
		Samples:          1,
	})
	if err != nil {
		fatal(err)
	}
	// Reproducible final outputs and their receipt are content-addressed. A warm
	// native command must therefore execute no handlers at all.
	zero := 0
	if err := benchmark.CheckBudget(result, nil, benchmark.Budget{MinSamples: 1, ExpectedExecutedSteps: &zero}); err != nil {
		fatal(err)
	}
	sample := result.Samples[0]
	if sample.CachedSteps <= 0 {
		fatal(fmt.Errorf("%s did not report any cached steps", strings.Join(args, " ")))
	}
	verifyReceipt(root)
	fmt.Printf("    warm cache verified: %.1fms wall, %d cached\n", sample.WallMS, sample.CachedSteps)
}

func verifyReceipt(root string) {
	receipt := ".wails/artifacts/receipt.json"
	verified, err := pipeline.VerifyArtifactReceipt(root, receipt)
	if err != nil {
		fatal(err)
	}
	fmt.Printf("    receipt verified: %d artifact(s)\n", len(verified.Artifacts))
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
