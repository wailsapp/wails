package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func configureRunTarget(options *RunOptions, goos, arch string, settings manifest.Run, build *manifestRunOptions) error {
	switch goos {
	case "android":
		if options.Destination != "" {
			return fmt.Errorf("--destination is only valid for iOS")
		}
		if options.Device != "" && options.Emulator != "" {
			return fmt.Errorf("--device and --emulator are mutually exclusive")
		}
		if len(settings.Args) != 0 || len(settings.Environment) != 0 {
			return fmt.Errorf("Android run does not support application arguments or runtime environment overrides; move shared defaults into desktop target run blocks")
		}
		build.Formats = []string{"apk"}
	case "ios":
		if !options.Plan && runtime.GOOS != "darwin" {
			return fmt.Errorf("iOS run requires macOS and Xcode")
		}
		if options.Emulator != "" {
			return fmt.Errorf("use --device to select an iOS simulator")
		}
		if options.Destination != "" && options.Destination != "device" && options.Destination != "simulator" {
			return fmt.Errorf("iOS --destination must be simulator or device")
		}
		if options.Destination == "device" {
			if options.Device == "" {
				return fmt.Errorf("physical iOS run requires --device")
			}
			// Keep credentials in the normal signing configuration; only select the outcome here.
			build.Loaded.Config.Selected = manifest.Profile{Name: "_wails_run_device", Targets: []manifest.ProfileTarget{{Target: goos + "/" + arch, Destination: "device", Sign: true}}}
			build.TargetOS, build.TargetArch = "", ""
		}
	default:
		if options.Device != "" || options.Emulator != "" || options.Destination != "" {
			return fmt.Errorf("device, emulator and destination options require a mobile target")
		}
	}
	return nil
}

func runMobileApplication(ctx context.Context, options *RunOptions, settings manifest.Run, build manifestRunOptions, plan pipeline.Plan, artifact, goos, arch string) error {
	root := build.Loaded.Config.Root
	if goos == "android" {
		operations := realAndroidDeployOperations()
		attach := operations.attach
		operations.attach = func(ctx context.Context, device androidDevice, identifier string) error {
			defer func() {
				cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_, _ = operations.adb(cleanup, "-s", device.Serial, "shell", "am", "force-stop", identifier)
			}()
			return attach(ctx, device, identifier)
		}

		operations.buildAPK = func(ctx context.Context, _ string, deviceArch string) (string, string, error) {
			if arch != "universal" && arch != deviceArch {
				return "", "", fmt.Errorf("device architecture is %s; use --target android/%s", deviceArch, deviceArch)
			}
			build.Context = ctx
			if _, err := runManifestPipelineResult(build); err != nil {
				return "", "", err
			}
			for _, node := range plan.Nodes {
				if spec, ok := node.Spec.(pipeline.PackageSpec); ok && spec.Format == "apk" {
					return artifact, spec.Project.Identifier, nil
				}
			}
			return "", "", fmt.Errorf("Android plan has no APK application identity")
		}
		return androidRunWithOperations(ctx, AndroidRunOptions{Device: options.Device, Emulator: options.Emulator, Logs: true}, "", operations)
	}
	device := options.Device
	var err error
	if options.Destination != "device" {
		device, err = selectIOSDevSimulator(ctx, device)
		if err != nil {
			return err
		}
	}
	if _, err = runManifestPipelineResult(build); err != nil {
		return err
	}
	identifier, err := runManifestTool(ctx, root, "/usr/libexec/PlistBuddy", "-c", "Print :CFBundleIdentifier", filepath.Join(artifact, "Info.plist"))
	if err != nil {
		return err
	}
	identifier = strings.TrimSpace(identifier)
	var install, launch []string
	if options.Destination == "device" {
		install = []string{"devicectl", "device", "install", "app", "--device", device, artifact}
		launch = iosDeviceRunCommand(device, identifier, settings)
	} else {
		install = []string{"simctl", "install", device, artifact}
		launch = []string{"simctl", "launch", "--console", "--terminate-running-process", device, identifier}
		defer func() {
			cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, _ = runManifestTool(cleanup, root, "xcrun", "simctl", "terminate", device, identifier)
		}()
	}
	if _, err := runManifestTool(ctx, root, "xcrun", install...); err != nil {
		return fmt.Errorf("install iOS application: %w", err)
	}
	if options.Destination != "device" {
		launch = append(launch, settings.Args...)
	}
	env := make([]string, 0, len(settings.Environment))
	if options.Destination != "device" {
		for key, value := range settings.Environment {
			env = append(env, "SIMCTL_CHILD_"+key+"="+value)
		}
	}
	return runApplicationProcess(ctx, root, "xcrun", launch, env)
}

func iosDeviceRunCommand(device, identifier string, settings manifest.Run) []string {
	args := []string{"devicectl", "device", "process", "launch", "--device", device, "--terminate-existing", "--console"}
	if len(settings.Environment) > 0 {
		// A map of strings is always JSON encodable. Pass it as one argument;
		// shell quoting would change the values delivered by devicectl.
		environment, _ := json.Marshal(settings.Environment)
		args = append(args, "--environment-variables", string(environment))
	}
	args = append(args, identifier)
	return append(args, settings.Args...)
}
