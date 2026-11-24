package commands

import (
	"context"
	"fmt"

	"github.com/pterm/pterm"
)

// HarmonyOSDeployOptions installs and launches a commercial HarmonyOS HAP
// through hdc. The device protocol is shared with OpenHarmony; only the
// requested build target and user-facing labels differ.
type HarmonyOSDeployOptions struct {
	HAP      string `name:"hap" description:"Existing HAP to install; build one first with wails3 build"`
	Device   string `name:"device" description:"hdc target serial (default: the only connected target)"`
	Ability  string `name:"ability" description:"Ability to start (defaults to EntryAbility)"`
	NoLaunch bool   `name:"no-launch" description:"Install the HAP without launching it"`
	Logs     bool   `name:"logs" description:"Stream hilog after launch until Ctrl+C"`
}

// HarmonyOSLogsOptions streams application logs from a connected device.
type HarmonyOSLogsOptions struct {
	Device string `name:"device" description:"hdc target serial"`
	Tag    string `name:"tag" description:"hilog tag filter (default: Wails)"`
}

// HarmonyOSDevices lists connected hdc targets.
func HarmonyOSDevices(_ *HarmonyOSDeployOptions) error {
	targets, err := openHarmonyTargets(context.Background())
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		pterm.Info.Println("No connected HarmonyOS devices")
		return nil
	}
	for _, target := range targets {
		pterm.Println(fmt.Sprintf("%s\t%s", target.Serial, target.State))
	}
	return nil
}

// HarmonyOSDeploy installs a commercial HAP and optionally launches it.
func HarmonyOSDeploy(options *HarmonyOSDeployOptions) error {
	if options.HAP == "" {
		return fmt.Errorf("--hap is required; build one first with `wails3 build -targets harmonyos/<arch> -formats hos_hap`")
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

// HarmonyOSLogs streams hilog from a connected device.
func HarmonyOSLogs(options *HarmonyOSLogsOptions) error {
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
