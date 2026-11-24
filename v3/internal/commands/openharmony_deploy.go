package commands

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pterm/pterm"
)

// OpenHarmonyDeployOptions installs and launches a HAP on a connected
// OpenHarmony device through hdc. It mirrors the install/launch half of the
// Android tooling without emulating ADB's reverse-port development server.
type OpenHarmonyDeployOptions struct {
	HAP      string `name:"hap" description:"Existing HAP to install; build one first with wails3 build"`
	Device   string `name:"device" description:"hdc target serial (default: the only connected target)"`
	Ability  string `name:"ability" description:"Ability to start (defaults to EntryAbility)"`
	NoLaunch bool   `name:"no-launch" description:"Install the HAP without launching it"`
	Logs     bool   `name:"logs" description:"Stream hilog after launch until Ctrl+C"`
}

// OpenHarmonyLogsOptions streams application logs from a connected device.
type OpenHarmonyLogsOptions struct {
	Device string `name:"device" description:"hdc target serial"`
	Tag    string `name:"tag" description:"hilog tag filter (default: Wails)"`
}

type openHarmonyTarget struct {
	Serial string
	State  string
}

func openHarmonyTargets(ctx context.Context) ([]openHarmonyTarget, error) {
	out, err := runManifestTool(ctx, "", "hdc", "list", "targets")
	if err != nil {
		return nil, err
	}
	var targets []openHarmonyTarget
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "[") {
			continue
		}
		fields := strings.Fields(line)
		target := openHarmonyTarget{Serial: fields[0]}
		if len(fields) > 1 {
			target.State = fields[1]
		}
		targets = append(targets, target)
	}
	return targets, nil
}

func resolveOpenHarmonyDevice(ctx context.Context, requested string) (string, error) {
	if requested != "" {
		return requested, nil
	}
	targets, err := openHarmonyTargets(ctx)
	if err != nil {
		return "", err
	}
	connected := make([]string, 0, len(targets))
	for _, target := range targets {
		if target.State == "" || target.State == "Connect" || target.State == "connected" {
			connected = append(connected, target.Serial)
		}
	}
	if len(connected) != 1 {
		return "", fmt.Errorf("expected exactly one connected hdc target, found %d; pass --device", len(connected))
	}
	return connected[0], nil
}

// OpenHarmonyDevices lists connected hdc targets.
func OpenHarmonyDevices(_ *OpenHarmonyDeployOptions) error {
	targets, err := openHarmonyTargets(context.Background())
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		pterm.Info.Println("No connected OpenHarmony devices")
		return nil
	}
	for _, target := range targets {
		pterm.Println(fmt.Sprintf("%s\t%s", target.Serial, target.State))
	}
	return nil
}

// OpenHarmonyDeploy installs a HAP and optionally launches it.
func OpenHarmonyDeploy(options *OpenHarmonyDeployOptions) error {
	if options.HAP == "" {
		return fmt.Errorf("--hap is required; build one first with `wails3 build -targets openharmony/<arch> -formats hap`")
	}
	bundle, err := openHarmonyBundleName(options.HAP)
	if err != nil {
		return err
	}
	device, err := resolveOpenHarmonyDevice(context.Background(), options.Device)
	if err != nil {
		return err
	}
	if _, err := runManifestTool(context.Background(), "", "hdc", "-t", device, "install", options.HAP); err != nil {
		return err
	}
	pterm.Success.Printfln("Installed %s on %s", bundle, device)
	if options.NoLaunch {
		return nil
	}
	ability := options.Ability
	if ability == "" {
		ability = "EntryAbility"
	}
	if _, err := runManifestTool(context.Background(), "", "hdc", "-t", device, "shell", "aa", "start", "-b", bundle, "-a", ability); err != nil {
		return err
	}
	pterm.Success.Printfln("Started %s/%s", bundle, ability)
	if options.Logs {
		return streamOpenHarmonyLogs(device, "Wails")
	}
	return nil
}

// OpenHarmonyLogs streams hilog from a connected device.
func OpenHarmonyLogs(options *OpenHarmonyLogsOptions) error {
	device, err := resolveOpenHarmonyDevice(context.Background(), options.Device)
	if err != nil {
		return err
	}
	tag := options.Tag
	if tag == "" {
		tag = "Wails"
	}
	return streamOpenHarmonyLogs(device, tag)
}

func streamOpenHarmonyLogs(device, tag string) error {
	cmd := exec.CommandContext(context.Background(), "hdc", "-t", device, "shell", "hilog", "-x", "-T", tag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// openHarmonyBundleName reads bundleName from the compiled module.json inside
// a HAP (a zip archive). No device connection is needed.
func openHarmonyBundleName(hap string) (string, error) {
	reader, err := zip.OpenReader(hap)
	if err != nil {
		return "", fmt.Errorf("open HAP %s: %w", hap, err)
	}
	defer reader.Close()
	for _, file := range reader.File {
		if file.Name != "module.json" {
			continue
		}
		stream, err := file.Open()
		if err != nil {
			return "", err
		}
		defer stream.Close()
		var document struct {
			App struct {
				BundleName string `json:"bundleName"`
			} `json:"app"`
		}
		if err := json.NewDecoder(stream).Decode(&document); err != nil {
			return "", fmt.Errorf("parse HAP module.json: %w", err)
		}
		if document.App.BundleName == "" {
			return "", fmt.Errorf("HAP module.json has no app.bundleName")
		}
		return document.App.BundleName, nil
	}
	return "", fmt.Errorf("HAP %s has no module.json", hap)
}
