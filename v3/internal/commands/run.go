package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/wailsapp/wails/v3/internal/features"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

type RunOptions struct {
	Device      string `name:"device" description:"Android adb serial or iOS device/simulator identifier"`
	Emulator    string `name:"emulator" description:"Android Virtual Device to start or reuse"`
	Destination string `name:"destination" description:"iOS destination: simulator (default) or device"`

	Tags   string `name:"tags" description:"Additional comma-separated Go build tags"`
	Target string `name:"target" description:"One platform/architecture to run (defaults to the host)"`
	Plan   bool   `name:"plan" description:"Show resolved build and launch settings without executing"`
}

// RunApplication builds and launches the local application. Absence of a
// manifest selects ordinary go run; invalid manifests never take that fallback.
func RunApplication(options *RunOptions, args []string) (resultErr error) {
	var cancellation *runCancellation
	defer func() {
		if errors.Is(resultErr, context.Canceled) && cancellation != nil {
			resultErr = &applicationExitError{err: resultErr, code: 128 + int(cancellation.signal.Load())}
			return
		}
		var exit *exec.ExitError
		if errors.As(resultErr, &exit) {
			resultErr = &applicationExitError{err: resultErr, code: exit.ExitCode()}
		}
	}()
	DisableFooter = true
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	ctx, stop := context.WithCancel(context.Background())
	cancellation = &runCancellation{}
	cancellation.signal.Store(int32(syscall.SIGINT))
	ctx = context.WithValue(ctx, runCancellationKey{}, cancellation)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)
	defer stop()
	go func() {
		select {
		case value := <-signals:
			if unixSignal, ok := value.(syscall.Signal); ok {
				cancellation.signal.Store(int32(unixSignal))
			}
			stop()
		case <-ctx.Done():
		}
	}()
	_, _, err = manifest.Discover(root)
	if errors.Is(err, fs.ErrNotExist) {
		if options.Device != "" || options.Emulator != "" || options.Destination != "" {
			return fmt.Errorf("mobile launch options require wails.hcl")
		}
		if options.Target != "" && options.Target != runtime.GOOS+"/"+runtime.GOARCH {
			return fmt.Errorf("go run fallback only supports the host target")
		}
		goArgs := []string{"run"}
		if options.Tags != "" {
			goArgs = append(goArgs, "-tags", options.Tags)
		}
		goArgs = append(goArgs, ".")
		goArgs = append(goArgs, args...)
		if options.Plan {
			return json.NewEncoder(os.Stdout).Encode(map[string]any{"command": append([]string{"go"}, goArgs...), "directory": root})
		}
		return runApplicationProcess(ctx, root, "go", goArgs, nil)
	}
	if err != nil {
		return err
	}
	if !features.WakeEnabled() {
		return fmt.Errorf("wails.hcl found: enable HCL with WAILS_EXP_USE_WAKE=1 to run this project")
	}
	loaded, err := manifest.Load(root, "")
	if err != nil {
		return err
	}
	root = loaded.Config.Root
	target := options.Target
	if target == "" {
		target = runtime.GOOS + "/" + runtime.GOARCH
	}
	goos, arch, ok := strings.Cut(target, "/")
	if !ok {
		return fmt.Errorf("run target must be platform/architecture")
	}
	if supported := loaded.Config.Project.SupportedPlatforms; len(supported) > 0 {
		allowed := false
		for _, platform := range supported {
			if platform == goos {
				allowed = true
			}
		}
		if !allowed {
			return fmt.Errorf("project %s does not support %s; supported platforms: %s", loaded.Config.Project.Name, goos, strings.Join(supported, ", "))
		}
	}
	settings, err := loaded.Config.RunForTarget(goos, arch)
	if err != nil {
		return err
	}
	if !options.Plan && goos != "android" && goos != "ios" && (goos != runtime.GOOS || (arch != runtime.GOARCH && !(goos == "darwin" && arch == "universal"))) {
		return fmt.Errorf("cannot run desktop target %s on %s/%s", target, runtime.GOOS, runtime.GOARCH)
	}
	if args != nil {
		settings.Args = append([]string{}, args...)
	}
	build := manifestRunOptions{Context: ctx, Verb: "run", Loaded: loaded, TargetOS: goos, TargetArch: arch, Tags: appendUniqueStrings(settings.Tags, splitComma(options.Tags)...)}
	if err := configureRunTarget(options, goos, arch, settings, &build); err != nil {
		return err
	}
	_, _, plan, err := resolveManifestPlan(build)
	if err != nil {
		return err
	}
	executable, err := runArtifactPath(root, plan, goos, loaded.Config.Project.BinaryName)
	if err != nil {
		return err
	}
	if options.Plan {
		keys := make([]string, 0, len(settings.Environment))
		for key := range settings.Environment {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		return json.NewEncoder(os.Stdout).Encode(map[string]any{"target": target, "executable": executable, "directory": root, "args": settings.Args, "environment_keys": keys, "compilers": resolvedPlanCompilers(plan)})
	}
	if goos == "android" || goos == "ios" {
		return runMobileApplication(ctx, options, settings, build, plan, executable, goos, arch)
	}
	if _, err := runManifestPipelineResult(build); err != nil {
		return err
	}
	if goos == "darwin" {
		if err := registerRunBundle(ctx, root, executable, runManifestTool); err != nil {
			return err
		}
	}
	return runApplicationProcess(ctx, root, executable, settings.Args, declaredEnvironment(nil, settings.Environment))
}

func runArtifactPath(root string, plan pipeline.Plan, goos, binary string) (string, error) {
	for _, key := range plan.Artifacts {
		node := plan.Nodes[key]
		if node.Artifact.Kind != pipeline.ArtifactBinary && node.Artifact.Format != "app" && node.Artifact.Format != "apk" {
			continue
		}
		path := filepath.Join(root, filepath.FromSlash(node.Output))
		if goos == "darwin" && node.Artifact.Format == "app" {
			path = filepath.Join(path, "Contents", "MacOS", binary)
		}
		return path, nil
	}
	return "", fmt.Errorf("build plan produced no runnable artifact for %s", goos)
}

func runApplicationProcess(ctx context.Context, root, executable string, args, environment []string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	process, err := startManifestProcess(root, executable, environment, args...)
	if err != nil {
		return err
	}
	select {
	case <-process.done:
		return process.waitError()
	case <-ctx.Done():
		termination := os.Signal(os.Interrupt)
		if cancellation, ok := ctx.Value(runCancellationKey{}).(*runCancellation); ok {
			termination = syscall.Signal(cancellation.signal.Load())
		}
		process.stopWithSignal(2*time.Second, termination)
		return ctx.Err()
	}
}

// ApplicationExitCode preserves desktop application exit status at the CLI.
type applicationExitError struct {
	err  error
	code int
}

func (e *applicationExitError) Error() string { return e.err.Error() }
func (e *applicationExitError) Unwrap() error { return e.err }
func ApplicationExitCode(err error) (int, bool) {
	var exit *applicationExitError
	if errors.As(err, &exit) && exit.code >= 0 {
		return exit.code, true
	}
	return 0, false
}

type runCancellationKey struct{}
type runCancellation struct{ signal atomic.Int32 }

// Register the assembled bundle before running its binary so Launch Services
// can deliver file associations and custom URL schemes to this build.
func registerRunBundle(ctx context.Context, root, executable string, runTool func(context.Context, string, string, ...string) (string, error)) error {
	bundle := filepath.Dir(filepath.Dir(filepath.Dir(executable)))
	const registrar = "/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister"
	if _, err := runTool(ctx, root, registrar, "-f", bundle); err != nil {
		return fmt.Errorf("register application bundle: %w", err)
	}
	return nil
}
