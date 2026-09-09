package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/dev"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

// mobileDevProcess owns the installed application and its log attachment.
// The staged artifact survives Stop so a failed replacement can reinstall it.
type mobileDevProcess struct {
	discard func()
	dev.Process
	artifact string
	ready    *dev.Readiness
	stop     func(time.Duration)
}

func (p *mobileDevProcess) Stop(t time.Duration)             { p.stop(t) }
func (p *mobileDevProcess) NeedsStopBeforeReplacement() bool { return true }
func (p *mobileDevProcess) Discard() {
	if p.discard != nil {
		p.discard()
	}
}

func mobileDevBuildOptions(ctx context.Context, options *DevOptions, loaded *manifest.Loaded, goos, arch, url string, port int) manifestRunOptions {
	copyLoaded := *loaded
	copyLoaded.Config.Selected = manifest.Profile{}
	result := manifestRunOptions{Context: ctx, Verb: "build", Loaded: &copyLoaded, TargetOS: goos, TargetArch: arch, Development: true, Tags: manifestDevTags(options), Environment: []string{"FRONTEND_DEVSERVER_URL=" + url, wailsVitePort + "=" + strconv.Itoa(port)}}
	if goos == "android" {
		result.Formats = []string{"apk"}
	}
	if goos == "ios" && options.Destination == "device" {
		copyLoaded.Config.Selected = manifest.Profile{Name: "_wails_dev_device", Targets: []manifest.ProfileTarget{{Target: goos + "/" + arch, Destination: "device", Sign: true}}}
		result.TargetOS, result.TargetArch = "", ""
	}
	return result
}
func mobileDevArtifact(root string, run manifestPipelineRun, goos, arch string) (string, error) {
	for _, key := range run.Plan.Artifacts {
		node := run.Plan.Nodes[key]
		if node.Artifact.Format == "app" || node.Artifact.Format == "apk" {
			return filepath.Join(root, filepath.FromSlash(node.Output)), nil
		}
	}
	return "", fmt.Errorf("development build produced no installable %s artifact", goos)
}

func configureMobileDev(ctx context.Context, options *DevOptions, output *dev.Output, ops *dev.Operations) (func(), error) {
	goos, arch, err := splitTarget(options.Target)
	if err != nil {
		return nil, err
	}
	if options.Profile != "" {
		return nil, fmt.Errorf("dev does not accept production profiles")
	}
	if goos == "ios" && runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("iOS development requires macOS and Xcode")
	}
	if goos == "ios" && arch != "arm64" {
		return nil, fmt.Errorf("iOS development requires ios/arm64")
	}
	if goos == "android" && options.Host != "" && options.Host != "127.0.0.1" {
		return nil, fmt.Errorf("Android dev uses adb reverse; --host must be 127.0.0.1")
	}
	if goos == "android" && options.Destination != "" {
		return nil, fmt.Errorf("--destination is only valid for iOS")
	}
	if goos == "ios" && options.Emulator != "" {
		return nil, fmt.Errorf("use --device for an iOS simulator")
	}
	if goos == "ios" && options.Destination != "" && options.Destination != "simulator" && options.Destination != "device" {
		return nil, fmt.Errorf("iOS --destination must be simulator or device")
	}
	if goos == "ios" && options.Destination == "device" && !options.Plan {
		ip := net.ParseIP(options.Host)
		if ip == nil || ip.IsLoopback() || ip.IsUnspecified() {
			return nil, fmt.Errorf("physical iOS dev requires --host with this Mac's reachable LAN IP and --device with the iPhone/iPad identifier")
		}
		if options.Device == "" {
			return nil, fmt.Errorf("physical iOS dev requires --device with the iPhone/iPad identifier")
		}
	}
	ops.ValidateTarget = func(string, string) error { return nil }
	ops.Plan = func(l *manifest.Loaded, o, a string) error {
		return printManifestPlan(mobileDevBuildOptions(ctx, options, l, o, a, "", 0), false)
	}
	ops.Build = func(c context.Context, l *manifest.Loaded, o, a, u string, p int) (manifestPipelineRun, error) {
		return runManifestPipelineResult(mobileDevBuildOptions(dev.WithOutput(c, output), options, l, o, a, u, p))
	}
	identifiers := map[string]string{}
	ops.BinaryPath = func(root string, run manifestPipelineRun, o, a string) (string, error) {
		path, err := mobileDevArtifact(root, run, o, a)
		if err != nil {
			return "", err
		}
		if o == "android" {
			for _, key := range run.Plan.Artifacts {
				if spec, ok := run.Plan.Nodes[key].Spec.(pipeline.PackageSpec); ok {
					identifiers[path] = spec.Project.Identifier
				}
			}
		}
		return path, nil
	}
	ops.BackendChanged = func(manifestPipelineRun, string, string) bool { return true }
	if options.Plan {
		return func() {}, nil
	}
	device := options.Device
	if goos == "android" {
		chosen, e := chooseAndroidDevice(ctx, AndroidRunOptions{Device: device, Emulator: options.Emulator}, realAndroidDeployOperations())
		if e != nil {
			return nil, e
		}
		device = chosen.Serial
		actual, e := androidDeviceABI(ctx, device)
		if e != nil {
			return nil, e
		}
		if actual != arch {
			return nil, fmt.Errorf("device architecture is %s; use --target android/%s", actual, actual)
		}
	} else if options.Destination != "device" {
		device, err = selectIOSDevSimulator(ctx, device)
		if err != nil {
			return nil, err
		}
	}
	staged := map[string]bool{}
	cleanup := func() {
		for p := range staged {
			os.RemoveAll(p)
		}
	}
	launch := func(root, artifact, url string, port int) (dev.Process, error) {
		stage, e := os.MkdirTemp(filepath.Join(root, ".wails", "dev"), "launch-")
		if e != nil {
			return nil, e
		}
		staged[stage] = true
		name := filepath.Base(artifact)
		if goos == "ios" {
			name = strings.TrimSuffix(name, ".signed")
		}
		saved := filepath.Join(stage, name)
		if e = copyManifestPath(artifact, saved); e != nil {
			os.RemoveAll(stage)
			delete(staged, stage)
			return nil, e
		}
		identifier := identifiers[artifact]
		identifiers[saved] = identifier
		child, err := launchMobileDev(ctx, options, output, device, root, saved, url, port, identifier)
		if err != nil {
			os.RemoveAll(stage)
			delete(staged, stage)
			delete(identifiers, saved)
			return nil, err
		}
		child.(*mobileDevProcess).discard = func() { os.RemoveAll(stage); delete(staged, stage); delete(identifiers, saved) }
		return child, nil
	}
	ops.StartApp = launch
	ops.WaitReady = func(c context.Context, p dev.Process, t time.Duration) error {
		return p.(*mobileDevProcess).ready.Wait(c, p, t)
	}
	ops.RestoreApp = func(c context.Context, root string, previous dev.Process, url string, port int) (dev.Process, error) {
		p, e := launch(root, previous.(*mobileDevProcess).artifact, url, port)
		if e != nil {
			return nil, e
		}
		if e = p.(*mobileDevProcess).ready.Wait(c, p, 30*time.Second); e != nil {
			p.Stop(time.Second)
			p.(*mobileDevProcess).Discard()
			return nil, e
		}
		return p, nil
	}
	return cleanup, nil
}
func selectIOSDevSimulator(ctx context.Context, device string) (string, error) {
	if device == "" {
		data, e := runManifestTool(ctx, "", "xcrun", "simctl", "list", "devices", "booted", "-j")
		if e != nil {
			return "", e
		}
		var list struct {
			Devices map[string][]struct {
				UDID string `json:"udid"`
			}
		}
		if e = json.Unmarshal([]byte(data), &list); e != nil {
			return "", e
		}
		for runtime, devices := range list.Devices {
			if strings.Contains(runtime, "iOS") {
				for _, d := range devices {
					if device != "" {
						return "", fmt.Errorf("multiple iOS simulators are booted; select one with --device")
					}
					device = d.UDID
				}
			}
		}
		if device == "" {
			return "", fmt.Errorf("no iOS simulator is booted; pass --device with a simulator UDID to boot it")
		}
	}
	// bootstatus starts a selected shutdown simulator and waits for usability.
	_, err := runManifestTool(ctx, "", "xcrun", "simctl", "bootstatus", device, "-b")
	return device, err
}
func launchMobileDev(ctx context.Context, options *DevOptions, output *dev.Output, device, root, artifact, url string, port int, applicationID string) (dev.Process, error) {
	host := "127.0.0.1"
	if strings.HasPrefix(options.Target, "ios/") && options.Destination == "device" {
		host = options.Host
	}
	ready, err := dev.NewReadinessOn(host)
	if err != nil {
		return nil, err
	}
	var success bool
	defer func() {
		if !success {
			ready.Close()
		}
	}()
	env := append(ready.Environment(), "FRONTEND_DEVSERVER_URL="+url, wailsVitePort+"="+strconv.Itoa(port))
	if strings.HasPrefix(options.Target, "android/") {
		p, e := launchAndroidDev(ctx, output, device, root, artifact, env, ready, port, applicationID)
		success = e == nil
		return p, e
	}
	identifier, e := runManifestTool(ctx, root, "/usr/libexec/PlistBuddy", "-c", "Print :CFBundleIdentifier", filepath.Join(artifact, "Info.plist"))
	if e != nil {
		return nil, e
	}
	identifier = strings.TrimSpace(identifier)
	var args, childEnv []string
	if options.Destination == "device" {
		if _, e = runManifestTool(ctx, root, "xcrun", "devicectl", "device", "install", "app", "--device", device, artifact); e != nil {
			return nil, e
		}
		for _, value := range env {
			childEnv = append(childEnv, "DEVICECTL_CHILD_"+value)
		}
		args = []string{"devicectl", "device", "process", "launch", "--device", device, "--terminate-existing", "--console", identifier}
	} else {
		if _, e = runManifestTool(ctx, root, "xcrun", "simctl", "install", device, artifact); e != nil {
			return nil, e
		}
		for _, value := range env {
			childEnv = append(childEnv, "SIMCTL_CHILD_"+value)
		}
		args = []string{"simctl", "launch", "--console", "--terminate-running-process", device, identifier}
	}
	child, e := startManifestProcessOutput([]*dev.Output{output}, root, "xcrun", childEnv, args...)
	if e != nil {
		return nil, e
	}
	var once sync.Once
	p := &mobileDevProcess{Process: child, artifact: artifact, ready: ready}
	p.stop = func(t time.Duration) {
		once.Do(func() {
			if options.Destination != "device" {
				c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				_, _ = runManifestTool(c, root, "xcrun", "simctl", "terminate", device, identifier)
				cancel()
			}
			child.Stop(t)
			ready.Close()
		})
	}
	success = true
	return p, nil
}

func runManifestTool(ctx context.Context, root, name string, args ...string) (string, error) {
	text, err := runManifestCommand(ctx, root, nil, name, args...)
	if err != nil {
		return text, &dev.Diagnostic{Component: name, Phase: "deploy", Command: append([]string{name}, args...), Output: text, ExitCode: -1, Err: fmt.Errorf("%s: %w", strings.TrimSpace(text), err)}
	}
	return text, nil
}
func launchAndroidDev(ctx context.Context, output *dev.Output, device, root, artifact string, env []string, ready *dev.Readiness, port int, identifier string) (dev.Process, error) {
	var err error
	if identifier == "" {
		return nil, fmt.Errorf("Android development plan has no application identifier")
	}
	adb := func(c context.Context, args ...string) (string, error) {
		return runManifestTool(c, root, "adb", append([]string{"-s", device}, args...)...)
	}
	if _, err = adb(ctx, "install", "-r", artifact); err != nil {
		return nil, err
	}
	// Own only new reverse mappings; --no-rebind never replaces another session.
	var mappings []string
	removeMappings := func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for _, p := range mappings {
			_, _ = adb(c, "reverse", "--remove", "tcp:"+p)
		}
	}
	for _, p := range []string{strconv.Itoa(port), ready.Port()} {
		if _, err = adb(ctx, "reverse", "--no-rebind", "tcp:"+p, "tcp:"+p); err != nil {
			removeMappings()
			return nil, err
		}
		mappings = append(mappings, p)
	}
	args := []string{"shell", "am", "start", "-S", "-n", identifier + "/com.wails.app.MainActivity"}
	for _, value := range env {
		k, v, _ := strings.Cut(value, "=")
		args = append(args, "--es", k, v)
	}
	if _, err = adb(ctx, args...); err != nil {
		removeMappings()
		return nil, err
	}
	childCtx, cancel := context.WithCancel(ctx)
	child := &manifestProcess{done: make(chan struct{})}
	stdout, stderr := output.Writer("android", "stdout"), output.Writer("android", "stderr")
	go func() {
		err := attachAndroidApplicationWithOperations(childCtx, androidDevice{Serial: device}, identifier, androidAttachOperations{
			deviceState: androidDeviceConnectionState, processID: androidApplicationProcessID,
			startLogs: func(c context.Context, s, p string) (<-chan error, error) {
				return startAndroidLogcatWithWriters(c, s, p, stdout, stderr)
			}, pollInterval: time.Second, startupWait: 15 * time.Second,
		})
		stdout.Close()
		stderr.Close()
		child.mu.Lock()
		child.err = err
		child.mu.Unlock()
		close(child.done)
	}()
	var once sync.Once
	p := &mobileDevProcess{Process: child, artifact: artifact, ready: ready}
	p.stop = func(time.Duration) {
		once.Do(func() {
			cancel()
			c, stop := context.WithTimeout(context.Background(), 5*time.Second)
			_, _ = adb(c, "shell", "am", "force-stop", identifier)
			stop()
			<-child.done
			removeMappings()
			ready.Close()
		})
	}
	return p, nil
}

func prepareAndroidDevNetwork(root string) error {
	files := map[string]string{
		"app/src/debug/AndroidManifest.xml":           `<manifest xmlns:android="http://schemas.android.com/apk/res/android" xmlns:tools="http://schemas.android.com/tools"><application android:networkSecurityConfig="@xml/wails_dev_network" tools:replace="android:networkSecurityConfig" /></manifest>`,
		"app/src/debug/res/xml/wails_dev_network.xml": `<network-security-config><base-config cleartextTrafficPermitted="false"/><domain-config cleartextTrafficPermitted="true"><domain>127.0.0.1</domain></domain-config></network-security-config>`,
	}
	for path, text := range files {
		p := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(p, []byte(text), 0644); err != nil {
			return err
		}
	}
	return nil
}
