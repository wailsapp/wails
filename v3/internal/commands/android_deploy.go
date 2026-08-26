package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/term"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

// AndroidRunOptions controls development APK deployment. A profile selects
// project configuration; deployment always builds the selected device ABI as
// a development APK because production AABs cannot be installed directly.
type AndroidRunOptions struct {
	Profile  string `name:"profile" description:"Build configuration profile (or pass it as the first argument)"`
	Device   string `name:"device" description:"Connected adb device serial"`
	Emulator string `name:"emulator" description:"Android Virtual Device name to start or reuse"`
	APK      string `name:"apk" description:"Install an existing APK instead of building one"`
	NoLaunch bool   `name:"no-launch" description:"Install the APK without launching the app"`
}

type AndroidDevicesOptions struct{}

type androidDevice struct {
	Serial, State, Model, Product string
	Emulator                      bool
}

type androidSelection struct {
	Label, Value string
}

type androidDeployOperations struct {
	interactive func() bool
	listDevices func(context.Context) ([]androidDevice, error)
	deviceABI   func(context.Context, string) (string, error)
	buildAPK    func(context.Context, string, string) (string, string, error)
	packageID   func(string) (string, error)
	adb         func(context.Context, ...string) (string, error)
	selectValue func(string, string, []androidSelection) (string, error)
	listAVDs    func(context.Context) ([]string, error)
	startAVD    func(context.Context, string) (androidDevice, error)
}

func AndroidRun(options *AndroidRunOptions, arguments []string) error {
	profile, err := manifestProfile(options.Profile, arguments)
	if err != nil {
		return err
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	operations := realAndroidDeployOperations()
	return androidRunWithOperations(ctx, *options, profile, operations)
}

func AndroidDevices(_ *AndroidDevicesOptions) error {
	devices, err := listAndroidDevices(context.Background())
	if err != nil {
		return err
	}
	if len(devices) == 0 {
		fmt.Println("No Android devices or emulators found.")
		return nil
	}
	fmt.Println("SERIAL\tTYPE\tSTATE\tMODEL")
	for _, device := range devices {
		kind := "device"
		if device.Emulator {
			kind = "emulator"
		}
		fmt.Printf("%s\t%s\t%s\t%s\n", device.Serial, kind, device.State, device.Model)
	}
	return nil
}

func androidRunWithOperations(ctx context.Context, options AndroidRunOptions, profile string, operations androidDeployOperations) error {
	if options.Device != "" && options.Emulator != "" {
		return fmt.Errorf("--device and --emulator are mutually exclusive")
	}
	device, err := chooseAndroidDevice(ctx, options, operations)
	if err != nil {
		return err
	}
	arch, err := operations.deviceABI(ctx, device.Serial)
	if err != nil {
		return fmt.Errorf("inspect Android device %s ABI: %w", device.Serial, err)
	}
	apk := options.APK
	packageID := ""
	if apk == "" {
		apk, packageID, err = operations.buildAPK(ctx, profile, arch)
		if err != nil {
			return err
		}
	} else {
		absolute, resolveErr := filepath.Abs(apk)
		if resolveErr != nil {
			return resolveErr
		}
		if info, statErr := os.Stat(absolute); statErr != nil || info.IsDir() {
			return fmt.Errorf("APK %s is not a readable file", absolute)
		}
		apk = absolute
		packageID, err = operations.packageID(profile)
		if err != nil {
			return err
		}
	}
	term.Header("Android Deploy")
	pterm.Info.Printfln("Installing %s on %s", apk, androidDeviceLabel(device))
	if output, installErr := operations.adb(ctx, "-s", device.Serial, "install", "-r", apk); installErr != nil {
		return fmt.Errorf("install APK on %s: %s: %w", device.Serial, strings.TrimSpace(output), installErr)
	}
	if options.NoLaunch {
		pterm.Success.Println("APK installed")
		return nil
	}
	component := packageID + "/com.wails.app.MainActivity"
	if output, launchErr := operations.adb(ctx, "-s", device.Serial, "shell", "am", "start", "-n", component); launchErr != nil {
		return fmt.Errorf("launch %s on %s: %s: %w", packageID, device.Serial, strings.TrimSpace(output), launchErr)
	}
	pterm.Success.Printfln("Launched %s on %s", packageID, androidDeviceLabel(device))
	return nil
}

func chooseAndroidDevice(ctx context.Context, options AndroidRunOptions, operations androidDeployOperations) (androidDevice, error) {
	if options.Emulator != "" {
		return operations.startAVD(ctx, options.Emulator)
	}
	devices, err := operations.listDevices(ctx)
	if err != nil {
		return androidDevice{}, err
	}
	if options.Device != "" {
		for _, device := range devices {
			if device.Serial != options.Device {
				continue
			}
			if device.State != "device" {
				return androidDevice{}, fmt.Errorf("Android device %s is %s; authorize or reconnect it before deployment", device.Serial, device.State)
			}
			return device, nil
		}
		return androidDevice{}, fmt.Errorf("Android device %s is not connected", options.Device)
	}
	available := make([]androidDevice, 0, len(devices))
	for _, device := range devices {
		if device.State == "device" {
			available = append(available, device)
		}
	}
	if len(available) == 1 {
		return available[0], nil
	}
	if !operations.interactive() {
		return androidDevice{}, fmt.Errorf("found %d available Android targets; pass --device or --emulator for non-interactive deployment", len(available))
	}
	selections := make([]androidSelection, 0, len(available)+1)
	for _, device := range available {
		selections = append(selections, androidSelection{Label: androidDeviceLabel(device), Value: device.Serial})
	}
	selections = append(selections, androidSelection{Label: "Start an Android emulator…", Value: "@start-emulator"})
	selected, err := operations.selectValue("Deployment target", "Choose a connected device or start an emulator", selections)
	if err != nil {
		return androidDevice{}, err
	}
	if selected != "@start-emulator" {
		for _, device := range available {
			if device.Serial == selected {
				return device, nil
			}
		}
		return androidDevice{}, fmt.Errorf("selected Android device %s is no longer available", selected)
	}
	avds, err := operations.listAVDs(ctx)
	if err != nil {
		return androidDevice{}, err
	}
	if len(avds) == 0 {
		return androidDevice{}, fmt.Errorf("no Android Virtual Devices are configured; create one with Android Studio or avdmanager")
	}
	choices := make([]androidSelection, len(avds))
	for index, avd := range avds {
		choices[index] = androidSelection{Label: avd, Value: avd}
	}
	avd, err := operations.selectValue("Android emulator", "Choose an Android Virtual Device", choices)
	if err != nil {
		return androidDevice{}, err
	}
	return operations.startAVD(ctx, avd)
}

func realAndroidDeployOperations() androidDeployOperations {
	operations := androidDeployOperations{}
	operations.interactive = term.IsTerminal
	operations.listDevices = listAndroidDevices
	operations.deviceABI = androidDeviceABI
	operations.buildAPK = buildAndroidDevelopmentAPK
	operations.packageID = func(profile string) (string, error) {
		loaded, err := manifest.Load(".", profile)
		if err != nil {
			return "", err
		}
		return loaded.Config.Project.Identifier, nil
	}
	operations.adb = runADB
	operations.selectValue = selectAndroidValue
	operations.listAVDs = listAndroidAVDs
	operations.startAVD = func(ctx context.Context, name string) (androidDevice, error) {
		return startAndroidAVD(ctx, name, operations.listDevices, operations.adb)
	}
	return operations
}

func parseADBDevices(output string) []androidDevice {
	var result []androidDevice
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 2 || fields[0] == "List" || strings.HasPrefix(fields[0], "*") {
			continue
		}
		device := androidDevice{Serial: fields[0], State: fields[1], Emulator: strings.HasPrefix(fields[0], "emulator-")}
		for _, field := range fields[2:] {
			key, value, found := strings.Cut(field, ":")
			if !found {
				continue
			}
			switch key {
			case "model":
				device.Model = value
			case "product":
				device.Product = value
			}
		}
		result = append(result, device)
	}
	return result
}

func listAndroidDevices(ctx context.Context) ([]androidDevice, error) {
	output, err := runADB(ctx, "devices", "-l")
	if err != nil {
		return nil, fmt.Errorf("list Android devices: %s: %w", strings.TrimSpace(output), err)
	}
	return parseADBDevices(output), nil
}

func runADB(ctx context.Context, arguments ...string) (string, error) {
	output, err := exec.CommandContext(ctx, "adb", arguments...).CombinedOutput()
	return string(output), err
}

func androidDeviceABI(ctx context.Context, serial string) (string, error) {
	output, err := runADB(ctx, "-s", serial, "shell", "getprop", "ro.product.cpu.abilist")
	if err != nil {
		return "", err
	}
	abis := strings.Split(strings.TrimSpace(output), ",")
	for _, candidate := range []struct{ abi, arch string }{{"arm64-v8a", "arm64"}, {"x86_64", "amd64"}} {
		for _, abi := range abis {
			if strings.TrimSpace(abi) == candidate.abi {
				return candidate.arch, nil
			}
		}
	}
	return "", fmt.Errorf("unsupported ABI list %q; Wails supports arm64-v8a and x86_64", strings.TrimSpace(output))
}

func buildAndroidDevelopmentAPK(ctx context.Context, profile, arch string) (string, string, error) {
	loaded, err := manifest.Load(".", profile)
	if err != nil {
		return "", "", err
	}
	if profile != "" {
		androidSelected := false
		for _, target := range loaded.Config.Selected.Targets {
			if strings.HasPrefix(target.Target, "android/") {
				androidSelected = true
				break
			}
		}
		if !androidSelected {
			return "", "", fmt.Errorf("profile %q does not select an Android target", profile)
		}
	}
	development := *loaded
	development.Config.Selected = manifest.Profile{}
	run, err := runManifestPipelineResult(manifestRunOptions{
		Context: ctx, Verb: "package", Loaded: &development,
		TargetOS: "android", TargetArch: arch, Formats: []string{"apk"}, Development: true,
	})
	if err != nil {
		return "", "", err
	}
	for _, key := range run.Plan.Artifacts {
		node := run.Plan.Nodes[key]
		if node.Artifact.Format == "apk" {
			return filepath.Join(development.Config.Root, filepath.FromSlash(node.Output)), development.Config.Project.Identifier, nil
		}
	}
	return "", "", fmt.Errorf("Android development build produced no APK")
}

func listAndroidAVDs(ctx context.Context) ([]string, error) {
	output, err := exec.CommandContext(ctx, "emulator", "-list-avds").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("list Android emulators: %s: %w", strings.TrimSpace(string(output)), err)
	}
	var result []string
	for _, line := range strings.Split(string(output), "\n") {
		if name := strings.TrimSpace(line); name != "" {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result, nil
}

func startAndroidAVD(ctx context.Context, name string, listDevices func(context.Context) ([]androidDevice, error), adb func(context.Context, ...string) (string, error)) (androidDevice, error) {
	devices, err := listDevices(ctx)
	if err != nil {
		return androidDevice{}, fmt.Errorf("list Android devices before starting emulator: %w", err)
	}
	for _, device := range devices {
		if !device.Emulator || device.State != "device" {
			continue
		}
		output, err := adb(ctx, "-s", device.Serial, "emu", "avd", "name")
		if err == nil && parseAndroidAVDName(output) == name {
			return device, nil
		}
	}
	root, err := os.Getwd()
	if err != nil {
		return androidDevice{}, err
	}
	logDirectory := filepath.Join(root, ".wails", "android")
	if err := os.MkdirAll(logDirectory, 0o755); err != nil {
		return androidDevice{}, err
	}
	logPath := filepath.Join(logDirectory, "emulator-"+strings.NewReplacer("/", "-", "\\", "-").Replace(name)+".log")
	log, err := os.Create(logPath)
	if err != nil {
		return androidDevice{}, err
	}
	command := exec.Command("emulator", "-avd", name, "-no-snapshot-save")
	command.Stdout, command.Stderr = log, log
	if err := command.Start(); err != nil {
		_ = log.Close()
		return androidDevice{}, fmt.Errorf("start Android emulator %q: %w", name, err)
	}
	_ = command.Process.Release()
	_ = log.Close()
	spinner := term.StartSpinner("Starting Android emulator " + name)
	defer term.StopSpinner(spinner)
	deadline := time.Now().Add(3 * time.Minute)
	var lastDiscoveryErr error
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return androidDevice{}, ctx.Err()
		case <-time.After(time.Second):
		}
		devices, lastDiscoveryErr = listDevices(ctx)
		if lastDiscoveryErr != nil {
			continue
		}
		for _, device := range devices {
			if !device.Emulator || device.State != "device" {
				continue
			}
			output, queryErr := adb(ctx, "-s", device.Serial, "emu", "avd", "name")
			if queryErr == nil && parseAndroidAVDName(output) == name {
				return device, nil
			}
		}
	}
	if lastDiscoveryErr != nil {
		return androidDevice{}, fmt.Errorf("Android emulator %q did not become ready within 3 minutes; last device discovery failed: %w; see %s", name, lastDiscoveryErr, logPath)
	}
	return androidDevice{}, fmt.Errorf("Android emulator %q did not become ready within 3 minutes; see %s", name, logPath)
}

func parseAndroidAVDName(output string) string {
	for _, line := range strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && line != "OK" {
			return line
		}
	}
	return ""
}

func selectAndroidValue(title, description string, selections []androidSelection) (string, error) {
	options := make([]huh.Option[string], len(selections))
	for index, selection := range selections {
		options[index] = huh.NewOption(selection.Label, selection.Value)
	}
	var value string
	err := huh.NewForm(huh.NewGroup(huh.NewSelect[string]().Title(title).Description(description).Options(options...).Value(&value))).Run()
	return value, err
}

func androidDeviceLabel(device androidDevice) string {
	kind := "Device"
	if device.Emulator {
		kind = "Emulator"
	}
	model := device.Model
	if model == "" {
		model = device.Product
	}
	if model == "" {
		return fmt.Sprintf("%s · %s", kind, device.Serial)
	}
	return fmt.Sprintf("%s · %s · %s", kind, strings.ReplaceAll(model, "_", " "), device.Serial)
}
