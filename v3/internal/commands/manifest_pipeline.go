package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/report/pulse"
	"github.com/wailsapp/wails/v3/internal/term"
	"github.com/wailsapp/wails/v3/internal/version"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
	"gopkg.in/yaml.v3"
)

type manifestRunOptions struct {
	Context                             context.Context
	Verb, Profile, TargetOS, TargetArch string
	Targets                             []pipeline.Target
	Formats                             []string
	Development, Force, Obfuscated      bool
	Tags                                []string
}

func runManifestPipeline(options manifestRunOptions) error {
	started := time.Now()
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	loaded, err := manifest.Load(root, options.Profile)
	if err != nil {
		return err
	}
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: options.Verb, TargetOS: options.TargetOS, TargetArch: options.TargetArch, Targets: options.Targets, Formats: options.Formats, Development: options.Development, ExtraTags: options.Tags, Obfuscated: options.Obfuscated})
	if err != nil {
		return err
	}
	reporter := pulse.New(os.Stdout, report.Normal)
	term.Header(strings.ToUpper(options.Verb[:1]) + options.Verb[1:])
	report.SetActive(reporter)
	defer report.SetActive(nil)
	reporter.BuildStart(options.Verb, plan.Target, len(plan.Nodes))
	ctx := options.Context
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()
	}
	results, err := (pipeline.Executor{Handler: &manifestHandler{root: root, config: loaded.Config}}).Execute(ctx, plan, pipeline.ExecuteOptions{Root: root, Force: options.Force, Reporter: reporter})
	if err != nil {
		reporter.BuildEnd(time.Since(started), false)
		return err
	}
	resultKeys := make([]string, 0, len(results))
	for key := range results {
		resultKeys = append(resultKeys, string(key))
	}
	sort.Strings(resultKeys)
	for _, rawKey := range resultKeys {
		key := pipeline.NodeKey(rawKey)
		node := plan.Nodes[key]
		result := results[key]
		if result.Output != "" && node.ArtifactKind != "" {
			reporter.Artifact(report.Artifact{Path: filepath.Join(root, filepath.FromSlash(result.Output)), Kind: node.ArtifactKind})
		}
	}
	reporter.BuildEnd(time.Since(started), true)
	return nil
}

type manifestHandler struct {
	root   string
	config manifest.Config
}

func (h *manifestHandler) Identity(_ context.Context, node pipeline.Node) (string, error) {
	var tools []string
	if executable, err := os.Executable(); err == nil {
		tools = append(tools, executable)
	}
	switch node.Kind {
	case pipeline.InstallFrontendDependencies, pipeline.BuildFrontend:
		manager := h.config.Frontend.PackageManager
		path, err := exec.LookPath(manager)
		if err != nil {
			return "", fmt.Errorf("%s not found: %w", manager, err)
		}
		tools = append(tools, path)
		if manager == "npm" {
			if node, err := exec.LookPath("node"); err == nil {
				tools = append(tools, node)
			}
		}
	case pipeline.GenerateBindings, pipeline.CompileApplication:
		name := "go"
		if node.Kind == pipeline.CompileApplication {
			spec, err := manifestNodeSpec[pipeline.CompileSpec](node)
			if err != nil {
				return "", err
			}
			if spec.Obfuscated {
				name = "garble"
			}
		}
		path, err := exec.LookPath(name)
		if err != nil {
			return "", fmt.Errorf("%s not found: %w", name, err)
		}
		tools = append(tools, path)
	case pipeline.PackageArtifact:
		spec, err := manifestNodeSpec[pipeline.PackageSpec](node)
		if err != nil {
			return "", err
		}
		name := map[string]string{"nsis": "makensis", "msix": "MakeAppx.exe"}[spec.Format]
		if name != "" {
			path, err := exec.LookPath(name)
			if err != nil {
				return "", fmt.Errorf("%s packaging requires %s: %w", spec.Format, name, err)
			}
			tools = append(tools, path)
		}
	case pipeline.SignArtifact:
		// The signer is implemented by this CLI; its executable identity below
		// is the complete tool identity.
	case pipeline.RunHook:
		spec, err := manifestNodeSpec[pipeline.HookSpec](node)
		if err != nil {
			return "", err
		}
		tools = append(tools, filepath.Join(h.root, filepath.FromSlash(spec.Hook.Script)))
	default:
		// Built-in handlers are identified by the running CLI below.
	}
	identity, err := toolIdentity(tools...)
	if err != nil {
		return "", err
	}
	return "wails-" + version.String() + ":" + string(node.Kind) + "|" + identity + "|env:" + relevantEnvironment(node), nil
}

func relevantEnvironment(node pipeline.Node) string {
	keys := []string{"PATH"}
	switch node.Kind {
	case pipeline.InstallFrontendDependencies, pipeline.BuildFrontend:
		keys = append(keys, "CI", "NODE_ENV", "NPM_CONFIG_USERCONFIG", "NPM_CONFIG_REGISTRY", "PNPM_HOME", "YARN_CACHE_FOLDER")
	case pipeline.GenerateBindings, pipeline.CompileApplication:
		keys = append(keys, "GOENV", "GOFLAGS", "GOTOOLCHAIN", "GOWORK", "CGO_ENABLED", "CC", "CXX", "CGO_CFLAGS", "CGO_CPPFLAGS", "CGO_CXXFLAGS", "CGO_LDFLAGS")
	case pipeline.RunHook:
		if node.Cache == pipeline.CacheArtifact {
			values := os.Environ()
			sort.Strings(values)
			return strings.Join(values, "\x00")
		}
	}
	values := make([]string, 0, len(keys))
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			values = append(values, key+"="+value)
		}
	}
	sort.Strings(values)
	return strings.Join(values, "\x00")
}

func (h *manifestHandler) Run(ctx context.Context, node pipeline.Node) (pipeline.RunResult, error) {
	switch node.Kind {
	case pipeline.InstallFrontendDependencies:
		spec, err := manifestNodeSpec[pipeline.InstallSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.install(ctx, spec)
	case pipeline.GenerateBindings:
		spec, err := manifestNodeSpec[pipeline.BindingsSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.bindings(spec, node.Output)
	case pipeline.BuildFrontend:
		spec, err := manifestNodeSpec[pipeline.FrontendSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.frontend(ctx, spec)
	case pipeline.CompileApplication:
		spec, err := manifestNodeSpec[pipeline.CompileSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.compile(ctx, spec)
	case pipeline.GeneratePlatformAssets:
		spec, err := manifestNodeSpec[pipeline.AssetsSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.assets(spec)
	case pipeline.PackageArtifact:
		spec, err := manifestNodeSpec[pipeline.PackageSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.packageArtifact(ctx, spec)
	case pipeline.SignArtifact:
		spec, err := manifestNodeSpec[pipeline.SignSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.sign(ctx, spec)
	case pipeline.RunHook:
		spec, err := manifestNodeSpec[pipeline.HookSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.hook(ctx, spec)
	default:
		return pipeline.RunResult{}, fmt.Errorf("wake: no handler for %s", node.Kind)
	}
}

func manifestNodeSpec[T pipeline.NodeSpec](node pipeline.Node) (T, error) {
	spec, ok := node.Spec.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("wake: %s node carries %T, expected %T", node.Kind, node.Spec, zero)
	}
	return spec, nil
}

func (h *manifestHandler) install(ctx context.Context, s pipeline.InstallSpec) (pipeline.RunResult, error) {
	args := []string{s.Command}
	if s.Manager == "npm" {
		args = append(args, "--no-audit", "--no-fund")
	}
	output, err := runPackageManager(ctx, filepath.Join(h.root, s.Directory), nil, s.Manager, args...)
	return pipeline.RunResult{Detail: output}, err
}

func (h *manifestHandler) bindings(s pipeline.BindingsSpec, output string) (pipeline.RunResult, error) {
	// GenerateBindings sets the process-wide footer flag; restore the caller's
	// value after using the legacy command adapter.
	previous := DisableFooter
	defer func() { DisableFooter = previous }()
	flagsString := ""
	if len(s.Tags) > 0 {
		flagsString = "-tags " + strings.Join(s.Tags, ",")
	}
	err := GenerateBindings(&flags.GenerateBindingsOptions{BuildFlagsString: flagsString, OutputDirectory: filepath.Join(h.root, filepath.FromSlash(output)), ModelsFilename: s.Config.ModelsFilename, IndexFilename: s.Config.IndexFilename, TimeType: s.Config.TimeType, TS: s.Config.TypeScript, UseInterfaces: s.Config.Interfaces, Clean: true, Silent: true, Obfuscated: s.Obfuscated}, nil)
	return pipeline.RunResult{}, err
}

func (h *manifestHandler) frontend(ctx context.Context, s pipeline.FrontendSpec) (pipeline.RunResult, error) {
	var args []string
	switch s.Manager {
	case "npm", "pnpm", "bun":
		args = []string{"run", s.Command}
	case "yarn":
		args = []string{s.Command}
	}
	if s.Manager == "npm" {
		args = append(args, "--silent")
	}
	env := []string{"PRODUCTION=" + fmt.Sprint(s.Production)}
	output, err := runPackageManager(ctx, filepath.Join(h.root, s.Directory), env, s.Manager, args...)
	return pipeline.RunResult{Detail: output}, err
}

func (h *manifestHandler) compile(ctx context.Context, s pipeline.CompileSpec) (pipeline.RunResult, error) {
	switch s.TargetOS {
	case "android":
		return h.compileAndroid(ctx, s)
	case "ios":
		return h.compileIOS(ctx, s)
	}
	if s.TargetOS == "darwin" && s.TargetArch == "universal" {
		return h.compileDarwinUniversal(ctx, s)
	}
	output := filepath.Join(h.root, filepath.FromSlash(s.Output))
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	tool, args := compileGoArgs(s)
	if s.TargetOS == "windows" {
		overlay, err := h.windowsResourceOverlay(s)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		args = append(args, "-overlay", overlay)
	}
	args = append(args, "-o", output, ".")
	env := []string{"GOOS=" + s.TargetOS, "GOARCH=" + s.TargetArch}
	if s.TargetOS == "darwin" && s.MinimumVersion != "" {
		minimumFlag := "-mmacosx-version-min=" + s.MinimumVersion
		env = append(env, "MACOSX_DEPLOYMENT_TARGET="+s.MinimumVersion, "CGO_CFLAGS="+minimumFlag, "CGO_LDFLAGS="+minimumFlag)
	}
	result, err := runManifestCommand(ctx, h.root, env, tool, args...)
	if err != nil && strings.Contains(result, "updates to go.mod needed") {
		return pipeline.RunResult{}, fmt.Errorf("Go module files need maintenance; run `go mod tidy` manually: %w", err)
	}
	return pipeline.RunResult{Detail: result}, err
}

func (h *manifestHandler) compileDarwinUniversal(ctx context.Context, s pipeline.CompileSpec) (pipeline.RunResult, error) {
	profile := h.config.Profile
	if profile == "" {
		profile = "default"
	}
	workspace := filepath.ToSlash(filepath.Join(".wails", "build", profile, "darwin-universal", "binaries"))
	inputs := make([]string, 0, 2)
	var details []string
	for _, arch := range []string{"amd64", "arm64"} {
		child := s
		child.TargetArch = arch
		child.Output = filepath.ToSlash(filepath.Join(workspace, h.config.Project.BinaryName+"-"+arch))
		result, err := h.compile(ctx, child)
		if err != nil {
			return result, err
		}
		if result.Detail != "" {
			details = append(details, result.Detail)
		}
		inputs = append(inputs, filepath.Join(h.root, filepath.FromSlash(child.Output)))
	}
	output := filepath.Join(h.root, filepath.FromSlash(s.Output))
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	previous := DisableFooter
	defer func() { DisableFooter = previous }()
	err := ToolLipo(&flags.Lipo{Inputs: inputs, Output: output})
	return pipeline.RunResult{Detail: strings.Join(details, "\n")}, err
}

func compileGoArgs(s pipeline.CompileSpec) (string, []string) {
	tool := "go"
	args := []string{"build"}
	if s.Obfuscated {
		tool = "garble"
		args = append(append([]string{}, s.GarbleArgs...), "build")
	}
	if len(s.Tags) > 0 {
		args = append(args, "-tags", strings.Join(s.Tags, ","))
	}
	if s.TrimPath {
		args = append(args, "-trimpath")
	}
	args = append(args, "-buildvcs=false")
	if len(s.CompilerFlags) > 0 {
		args = append(args, "-gcflags", strings.Join(s.CompilerFlags, " "))
	}
	ld := append([]string{}, s.LinkerFlags...)
	if s.Strip {
		ld = appendUniqueStrings(ld, "-w", "-s")
	}
	if len(ld) > 0 {
		args = append(args, "-ldflags", strings.Join(ld, " "))
	}
	return tool, args
}

func (h *manifestHandler) windowsResourceOverlay(s pipeline.CompileSpec) (string, error) {
	assets := filepath.Join(h.root, filepath.FromSlash(s.Assets), "windows")
	generated := filepath.Join(h.root, ".wails", "generated", "windows", s.TargetArch)
	if err := os.MkdirAll(generated, 0o755); err != nil {
		return "", err
	}
	syso := filepath.Join(generated, "wails_windows_"+s.TargetArch+".syso")
	if err := GenerateSyso(&SysoOptions{Manifest: filepath.Join(assets, "wails.exe.manifest"), Info: filepath.Join(assets, "info.json"), Icon: filepath.Join(assets, "icon.ico"), Out: syso, Arch: s.TargetArch}); err != nil {
		return "", err
	}
	virtual := filepath.Join(h.root, "wails_windows_"+s.TargetArch+".syso")
	data, err := json.Marshal(struct {
		Replace map[string]string `json:"Replace"`
	}{map[string]string{virtual: syso}})
	if err != nil {
		return "", err
	}
	overlay := filepath.Join(generated, "overlay.json")
	return overlay, os.WriteFile(overlay, data, 0o644)
}

func (h *manifestHandler) compileAndroid(ctx context.Context, s pipeline.CompileSpec) (pipeline.RunResult, error) {
	sdk := firstNonempty(os.Getenv("ANDROID_HOME"), os.Getenv("ANDROID_SDK_ROOT"))
	ndk := os.Getenv("ANDROID_NDK_HOME")
	if ndk == "" && sdk != "" {
		matches, _ := filepath.Glob(filepath.Join(sdk, "ndk", "*"))
		sort.Strings(matches)
		if len(matches) > 0 {
			ndk = matches[len(matches)-1]
		}
	}
	if ndk == "" {
		return pipeline.RunResult{}, fmt.Errorf("Android NDK not found; set ANDROID_NDK_HOME or install an NDK under ANDROID_HOME")
	}
	hostTag := "linux-x86_64"
	switch runtime.GOOS {
	case "darwin":
		hostTag = "darwin-x86_64"
	case "windows":
		hostTag = "windows-x86_64"
	}
	triple := "aarch64-linux-android"
	jni := "arm64-v8a"
	if s.TargetArch == "amd64" {
		triple = "x86_64-linux-android"
		jni = "x86_64"
	}
	if s.TargetArch != "arm64" && s.TargetArch != "amd64" {
		return pipeline.RunResult{}, fmt.Errorf("unsupported Android architecture %s", s.TargetArch)
	}
	minSDK := s.MinimumVersion
	if minSDK == "" {
		minSDK = "21"
	}
	toolchain := filepath.Join(ndk, "toolchains", "llvm", "prebuilt", hostTag, "bin")
	cc := filepath.Join(toolchain, triple+minSDK+"-clang")
	cxx := filepath.Join(toolchain, triple+minSDK+"-clang++")
	if runtime.GOOS == "windows" {
		cc += ".cmd"
		cxx += ".cmd"
	}
	if _, err := os.Stat(cc); err != nil {
		return pipeline.RunResult{}, fmt.Errorf("Android compiler not found: %s", cc)
	}
	output := filepath.Join(h.root, s.Output)
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	s.Tags = appendUniqueStrings(s.Tags, "android")
	tool, args := compileGoArgs(s)
	args = append(args, "-buildmode=c-shared", "-overlay", filepath.Join(h.root, s.Assets, "android", "overlay.json"), "-o", output, ".")
	env := []string{"GOOS=android", "GOARCH=" + s.TargetArch, "CGO_ENABLED=1", "CC=" + cc, "CXX=" + cxx, "WAILS_ANDROID_JNI=" + jni}
	result, err := runManifestCommand(ctx, h.root, env, tool, args...)
	return pipeline.RunResult{Detail: result}, err
}

func (h *manifestHandler) compileIOS(ctx context.Context, s pipeline.CompileSpec) (pipeline.RunResult, error) {
	if runtime.GOOS != "darwin" {
		return pipeline.RunResult{}, fmt.Errorf("iOS builds require macOS and Xcode")
	}
	variant := s.Variant
	if variant == "" {
		variant = "simulator"
	}
	sdk := "iphonesimulator"
	suffix := "-simulator"
	minFlag := "-mios-simulator-version-min="
	if variant == "device" {
		sdk = "iphoneos"
		suffix = ""
		minFlag = "-miphoneos-version-min="
	}
	sdkPath, err := commandOutput(ctx, h.root, "xcrun", "--sdk", sdk, "--show-sdk-path")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	clang, err := commandOutput(ctx, h.root, "xcrun", "--sdk", sdk, "--find", "clang")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	minimum := s.MinimumVersion
	if minimum == "" {
		minimum = "15.0"
	}
	target := s.TargetArch + "-apple-ios" + minimum + suffix
	output := filepath.Join(h.root, s.Output)
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	s.Tags = appendUniqueStrings(s.Tags, "ios")
	tool, args := compileGoArgs(s)
	args = append(args, "-buildmode=c-archive", "-overlay", filepath.Join(h.root, s.Assets, "ios", "xcode", "overlay.json"), "-o", output, ".")
	flags := "-isysroot " + sdkPath + " -target " + target + " " + minFlag + minimum
	env := []string{"GOOS=ios", "GOARCH=" + s.TargetArch, "CGO_ENABLED=1", "CC=" + clang, "CGO_CFLAGS=" + flags, "CGO_LDFLAGS=-isysroot " + sdkPath + " -target " + target}
	result, err := runManifestCommand(ctx, h.root, env, tool, args...)
	return pipeline.RunResult{Detail: result}, err
}

func commandOutput(ctx context.Context, dir, name string, args ...string) (string, error) {
	output, err := runManifestCommand(ctx, dir, nil, name, args...)
	if err != nil {
		return "", fmt.Errorf("%s: %s: %w", name, output, err)
	}
	return strings.TrimSpace(output), nil
}
func firstNonempty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func (h *manifestHandler) assets(s pipeline.AssetsSpec) (pipeline.RunResult, error) {
	output := filepath.Join(h.root, filepath.FromSlash(s.Directory))
	if err := removeGenerated(output); err != nil {
		return pipeline.RunResult{}, err
	}
	parent := filepath.Dir(output)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	tmp, err := os.MkdirTemp(parent, ".assets-*")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	defer os.RemoveAll(tmp)
	options := &BuildAssetsOptions{Dir: tmp, Name: s.Project.Name, BinaryName: s.Project.BinaryName, ProductName: s.Project.ProductName, ProductDescription: s.Project.Description, ProductVersion: s.Project.Version, ProductCompany: s.Project.CompanyName, ProductIdentifier: s.Project.Identifier, ProductCopyright: s.Project.Copyright, ProductComments: s.Project.Comments, Silent: true, Typescript: h.config.Frontend.Bindings.TypeScript, UseInterfaces: h.config.Frontend.Bindings.Interfaces}
	if err := GenerateBuildAssets(options); err != nil {
		return pipeline.RunResult{}, err
	}
	if s.Project.Icon != "" {
		source := filepath.Join(h.root, filepath.FromSlash(s.Project.Icon))
		icon := filepath.Join(tmp, "appicon.png")
		if err := copyManifestPath(source, icon); err != nil {
			return pipeline.RunResult{}, err
		}
		if err := GenerateIcons(&IconsOptions{Input: icon, WindowsFilename: filepath.Join(tmp, "windows", "icon.ico"), MacFilename: filepath.Join(tmp, "darwin", "icons.icns")}); err != nil {
			return pipeline.RunResult{}, err
		}
	}
	configPath, err := writeGeneratedConfig(tmp, s.Project, s.Associations, s.Protocols)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	if err := UpdateBuildAssets(&UpdateBuildAssetsOptions{Dir: tmp, Name: s.Project.Name, BinaryName: s.Project.BinaryName, Config: configPath, Silent: true}); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	source := filepath.Join(tmp, s.TargetOS)
	if err := copyManifestPath(source, filepath.Join(output, s.TargetOS)); err != nil {
		return pipeline.RunResult{}, err
	}
	_ = os.Remove(filepath.Join(output, s.TargetOS, "Taskfile.yml"))
	if err := copyManifestPath(filepath.Join(tmp, "appicon.png"), filepath.Join(output, "appicon.png")); err != nil {
		return pipeline.RunResult{}, err
	}
	resolvedConfig := filepath.Join(output, "config.yml")
	if err := copyManifestPath(configPath, resolvedConfig); err != nil {
		return pipeline.RunResult{}, err
	}
	switch s.TargetOS {
	case "ios":
		xcode := filepath.Join(output, "ios", "xcode")
		if err := IOSOverlayGen(&IOSOverlayGenOptions{Out: filepath.Join(xcode, "overlay.json"), Config: resolvedConfig}); err != nil {
			return pipeline.RunResult{}, err
		}
		if err := IOSXcodeGen(&IOSXcodeGenOptions{OutDir: xcode, Config: resolvedConfig}); err != nil {
			return pipeline.RunResult{}, err
		}
	case "android":
		if err := AndroidOverlayGen(&AndroidOverlayGenOptions{Out: filepath.Join(output, "android", "overlay.json"), Config: resolvedConfig}); err != nil {
			return pipeline.RunResult{}, err
		}
	}
	if err := applyGeneratedTargetSettings(output, s); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{}, nil
}

func applyGeneratedTargetSettings(output string, spec pipeline.AssetsSpec) error {
	if spec.Project.BuildNumber > 0 {
		value := strconv.Itoa(spec.Project.BuildNumber)
		for _, path := range []string{
			filepath.Join(output, "darwin", "Info.plist"),
			filepath.Join(output, "ios", "xcode", "main", "Info.plist"),
		} {
			if err := replacePlistString(path, "CFBundleVersion", value); err != nil {
				return err
			}
		}
		gradle := filepath.Join(output, "android", "app", "build.gradle")
		if err := replaceGeneratedPattern(gradle, regexp.MustCompile(`(?m)^(\s*versionCode\s+)\d+`), "${1}"+value); err != nil {
			return err
		}
		nsis := filepath.Join(output, "windows", "nsis", "project.nsi")
		if err := replaceGeneratedLiteral(nsis, `${INFO_PRODUCTVERSION}.0`, `${INFO_PRODUCTVERSION}.`+value); err != nil {
			return err
		}
	}
	if spec.MinimumVersion != "" {
		keys := map[string]string{"darwin": "LSMinimumSystemVersion", "ios": "MinimumOSVersion"}
		if key := keys[spec.TargetOS]; key != "" {
			path := filepath.Join(output, spec.TargetOS, "Info.plist")
			if spec.TargetOS == "ios" {
				path = filepath.Join(output, "ios", "xcode", "main", "Info.plist")
			}
			if err := replacePlistString(path, key, spec.MinimumVersion); err != nil {
				return err
			}
		}
	}
	return nil
}

func replacePlistString(path, key, value string) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	marker := "<key>" + key + "</key>"
	start := bytes.Index(data, []byte(marker))
	if start < 0 {
		return nil
	}
	stringStart := bytes.Index(data[start+len(marker):], []byte("<string>"))
	if stringStart < 0 {
		return nil
	}
	stringStart += start + len(marker) + len("<string>")
	stringEnd := bytes.Index(data[stringStart:], []byte("</string>"))
	if stringEnd < 0 {
		return nil
	}
	stringEnd += stringStart
	result := append(append(append([]byte(nil), data[:stringStart]...), value...), data[stringEnd:]...)
	return os.WriteFile(path, result, 0o644)
}

func replaceGeneratedPattern(path string, pattern *regexp.Regexp, replacement string) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return os.WriteFile(path, pattern.ReplaceAll(data, []byte(replacement)), 0o644)
}

func replaceGeneratedLiteral(path, old, replacement string) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes.ReplaceAll(data, []byte(old), []byte(replacement)), 0o644)
}

func (h *manifestHandler) packageArtifact(ctx context.Context, s pipeline.PackageSpec) (pipeline.RunResult, error) {
	switch s.Format {
	case "app":
		if s.TargetOS == "ios" {
			return h.packageIOS(ctx, s, false)
		}
		return h.packageApp(s)
	case "dmg":
		if _, err := h.packageApp(pipeline.PackageSpec{Project: s.Project, Binary: s.Binary, Assets: s.Assets, Output: strings.TrimSuffix(s.Output, ".dmg") + ".app", TargetOS: s.TargetOS, TargetArch: s.TargetArch, Format: "app"}); err != nil {
			return pipeline.RunResult{}, err
		}
		err := ToolPackage(&flags.ToolPackage{Format: "dmg", ExecutableName: s.Project.BinaryName, Out: filepath.Dir(filepath.Join(h.root, s.Output))})
		return pipeline.RunResult{}, err
	case "appimage":
		return h.packageAppImage(s)
	case "deb", "rpm", "archlinux":
		return h.packageLinux(s)
	case "nsis":
		return h.packageNSIS(ctx, s)
	case "msix":
		return h.packageMSIX(s)
	case "ipa":
		return h.packageIOS(ctx, s, true)
	case "apk", "aab":
		return h.packageAndroid(ctx, s)
	default:
		return pipeline.RunResult{}, fmt.Errorf("unsupported package format %q", s.Format)
	}
}

func (h *manifestHandler) packageApp(s pipeline.PackageSpec) (pipeline.RunResult, error) {
	output := filepath.Join(h.root, filepath.FromSlash(s.Output))
	if err := removeGenerated(output); err != nil {
		return pipeline.RunResult{}, err
	}
	macos := filepath.Join(output, "Contents", "MacOS")
	resources := filepath.Join(output, "Contents", "Resources")
	if err := os.MkdirAll(macos, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := os.MkdirAll(resources, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := copyManifestPath(filepath.Join(h.root, s.Binary), filepath.Join(macos, s.Project.BinaryName)); err != nil {
		return pipeline.RunResult{}, err
	}
	assets := filepath.Join(h.root, s.Assets, "darwin")
	for _, name := range []string{"Info.plist", "icons.icns", "Assets.car"} {
		source := filepath.Join(assets, name)
		if _, err := os.Stat(source); err == nil {
			dest := filepath.Join(output, "Contents", name)
			if name == "icons.icns" || name == "Assets.car" {
				dest = filepath.Join(resources, name)
			}
			if err := copyManifestPath(source, dest); err != nil {
				return pipeline.RunResult{}, err
			}
		}
	}
	return pipeline.RunResult{}, nil
}

func (h *manifestHandler) packageAppImage(s pipeline.PackageSpec) (pipeline.RunResult, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return pipeline.RunResult{}, err
	}
	defer func() { _ = os.Chdir(workingDirectory) }()
	assets := filepath.Join(h.root, s.Assets)
	desktop := filepath.Join(assets, s.Project.BinaryName+".desktop")
	if err := GenerateDotDesktop(&DotDesktopOptions{OutputFile: desktop, Type: "Application", Name: s.Project.ProductName, Exec: s.Project.BinaryName, Icon: s.Project.BinaryName, Comment: s.Project.Description, Categories: "Utility;", Version: "1.0"}); err != nil {
		return pipeline.RunResult{}, err
	}
	buildDir := filepath.Join(h.root, ".wails", "package", "appimage", s.TargetArch)
	outDir := filepath.Dir(filepath.Join(h.root, s.Output))
	if err := GenerateAppImage(&GenerateAppImageOptions{Binary: filepath.Join(h.root, s.Binary), Icon: filepath.Join(assets, "appicon.png"), DesktopFile: desktop, OutputDir: outDir, BuildDir: buildDir}); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{}, ensureExpectedAppImage(outDir, filepath.Join(h.root, s.Output))
}

func (h *manifestHandler) packageLinux(s pipeline.PackageSpec) (pipeline.RunResult, error) {
	config, err := h.prepareLinuxPackageConfig(s)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	if s.Config.Template != "" {
		config = filepath.Join(h.root, filepath.FromSlash(s.Config.Template))
	}
	packageRoot := filepath.Join(h.root, ".wails", "package")
	if err := os.MkdirAll(packageRoot, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	outDir, err := os.MkdirTemp(packageRoot, "linux-output-*")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	defer os.RemoveAll(outDir)
	if err := ToolPackage(&flags.ToolPackage{Format: s.Format, ExecutableName: s.Project.BinaryName, ConfigPath: config, Out: outDir}); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{}, findAndMovePackage(outDir, s.Project.BinaryName, s.Format, filepath.Join(h.root, s.Output))
}

func (h *manifestHandler) prepareLinuxPackageConfig(s pipeline.PackageSpec) (string, error) {
	if s.Config.Template != "" {
		return filepath.Join(h.root, filepath.FromSlash(s.Config.Template)), nil
	}
	assets := filepath.Join(h.root, s.Assets)
	workspace := filepath.Join(h.root, ".wails", "package", "linux", s.TargetArch, s.Format)
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		return "", err
	}
	desktop := filepath.Join(workspace, s.Project.BinaryName+".desktop")
	if err := GenerateDotDesktop(&DotDesktopOptions{OutputFile: desktop, Type: "Application", Name: s.Project.ProductName, Exec: s.Project.BinaryName, Icon: s.Project.BinaryName, Comment: s.Project.Description, Categories: "Utility;", Version: "1.0"}); err != nil {
		return "", err
	}
	source := filepath.Join(assets, "linux", "nfpm", "nfpm.yaml")
	data, err := os.ReadFile(source)
	if err != nil {
		return "", err
	}
	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err != nil {
		return "", err
	}
	config["name"] = s.Project.BinaryName
	config["arch"] = s.TargetArch
	config["platform"] = "linux"
	config["version"] = s.Project.Version
	if contents, ok := config["contents"].([]any); ok {
		for _, raw := range contents {
			entry, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			source, _ := entry["src"].(string)
			switch {
			case strings.Contains(source, "/bin/"):
				entry["src"] = filepath.Join(h.root, s.Binary)
			case strings.HasSuffix(source, "appicon.png"):
				entry["src"] = filepath.Join(assets, "appicon.png")
			case strings.HasSuffix(source, ".desktop"):
				entry["src"] = desktop
			}
		}
	}
	if scripts, ok := config["scripts"].(map[string]any); ok {
		for name, raw := range scripts {
			path, ok := raw.(string)
			if ok {
				scripts[name] = filepath.Join(assets, "linux", "nfpm", "scripts", filepath.Base(path))
			}
		}
	}
	resolved, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}
	path := filepath.Join(workspace, "nfpm.yaml")
	return path, os.WriteFile(path, resolved, 0o644)
}

func (h *manifestHandler) packageNSIS(ctx context.Context, s pipeline.PackageSpec) (pipeline.RunResult, error) {
	dir := filepath.Join(h.root, s.Assets, "windows", "nsis")
	if s.Config.Template != "" {
		if err := copyManifestPath(filepath.Join(h.root, filepath.FromSlash(s.Config.Template)), filepath.Join(dir, "project.nsi")); err != nil {
			return pipeline.RunResult{}, err
		}
	}
	if err := GenerateWebView2Bootstrapper(&GenerateWebView2Options{Directory: dir}); err != nil {
		return pipeline.RunResult{}, err
	}
	flag := "AMD64"
	if s.TargetArch == "arm64" {
		flag = "ARM64"
	}
	binary, err := filepath.Abs(filepath.Join(h.root, s.Binary))
	if err != nil {
		return pipeline.RunResult{}, err
	}
	result, err := runManifestCommand(ctx, dir, nil, "makensis", "-DARG_WAILS_"+flag+"_BINARY="+binary, "project.nsi")
	if err != nil {
		return pipeline.RunResult{Detail: result}, err
	}
	generated := filepath.Join(h.root, filepath.Dir(s.Assets), "bin", s.Project.Name+"-"+s.TargetArch+"-installer.exe")
	if err := copyManifestPath(generated, filepath.Join(h.root, s.Output)); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{Detail: result}, nil
}

func (h *manifestHandler) packageMSIX(s pipeline.PackageSpec) (pipeline.RunResult, error) {
	if runtime.GOOS != "windows" {
		return pipeline.RunResult{}, fmt.Errorf("MSIX packaging requires Windows")
	}
	workspace := filepath.Join(h.root, ".wails", "package", "msix", s.TargetArch)
	if err := removeGenerated(workspace); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	var associations []map[string]any
	for _, association := range h.config.Associations {
		for _, extension := range association.Extensions {
			associations = append(associations, map[string]any{"ext": extension, "name": association.Name, "description": association.Description, "iconName": association.Icon, "role": association.Role, "mimeType": association.MIMEType})
		}
	}
	config := map[string]any{"info": map[string]any{"companyName": s.Project.CompanyName, "productName": s.Project.ProductName, "productIdentifier": s.Project.Identifier, "description": s.Project.Description, "copyright": s.Project.Copyright, "comments": s.Project.Comments, "version": s.Project.Version}, "fileAssociations": associations}
	data, err := json.Marshal(config)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	configPath := filepath.Join(workspace, "config.json")
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		return pipeline.RunResult{}, err
	}
	arch := map[string]string{"amd64": "x64", "arm64": "arm64", "386": "x86"}[s.TargetArch]
	if arch == "" {
		return pipeline.RunResult{}, fmt.Errorf("unsupported MSIX architecture %s", s.TargetArch)
	}
	signing := h.config.Signing.Windows
	publisher := "CN=" + s.Project.CompanyName
	if signing.Identity != "" {
		publisher = signing.Identity
	}
	err = ToolMSIX(&flags.ToolMSIX{ConfigPath: configPath, Publisher: publisher, CertificatePath: signing.Certificate, Arch: arch, ExecutableName: s.Project.BinaryName + ".exe", ExecutablePath: filepath.Join(h.root, s.Binary), OutputPath: filepath.Join(h.root, s.Output), UseMakeAppx: true})
	return pipeline.RunResult{}, err
}

func (h *manifestHandler) packageAndroid(ctx context.Context, s pipeline.PackageSpec) (pipeline.RunResult, error) {
	workspace := filepath.Join(h.root, ".wails", "package", "android", s.TargetArch, s.Format)
	if err := removeGenerated(workspace); err != nil {
		return pipeline.RunResult{}, err
	}
	source := filepath.Join(h.root, s.Assets, "android")
	if s.Config.Template != "" {
		source = filepath.Join(h.root, filepath.FromSlash(s.Config.Template))
	}
	if err := copyManifestPath(source, workspace); err != nil {
		return pipeline.RunResult{}, err
	}
	abi := "arm64-v8a"
	if s.TargetArch == "amd64" {
		abi = "x86_64"
	}
	jni := filepath.Join(workspace, "app", "src", "main", "jniLibs", abi, "libwails.so")
	if err := copyManifestPath(filepath.Join(h.root, s.Binary), jni); err != nil {
		return pipeline.RunResult{}, err
	}
	gradlew := filepath.Join(workspace, "gradlew")
	if err := os.Chmod(gradlew, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	task := "assembleRelease"
	produced := filepath.Join(workspace, "app", "build", "outputs", "apk", "release", "app-release.apk")
	if s.Format == "aab" {
		task = "bundleRelease"
		produced = filepath.Join(workspace, "app", "build", "outputs", "bundle", "release", "app-release.aab")
	}
	output, err := runManifestCommand(ctx, workspace, nil, gradlew, task)
	if err != nil {
		return pipeline.RunResult{Detail: output}, err
	}
	if err := copyManifestPath(produced, filepath.Join(h.root, s.Output)); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{Detail: output}, nil
}

func (h *manifestHandler) packageIOS(ctx context.Context, s pipeline.PackageSpec, ipa bool) (pipeline.RunResult, error) {
	if runtime.GOOS != "darwin" {
		return pipeline.RunResult{}, fmt.Errorf("iOS packaging requires macOS and Xcode")
	}
	variant := s.Variant
	if variant == "" {
		variant = "simulator"
	}
	if ipa && variant != "device" {
		return pipeline.RunResult{}, fmt.Errorf("IPA packaging requires [targets.ios.%s] variant = \"device\"", s.TargetArch)
	}
	sdk := "iphonesimulator"
	suffix := "-simulator"
	if variant == "device" {
		sdk = "iphoneos"
		suffix = ""
	}
	minimum := s.MinimumVersion
	if minimum == "" {
		minimum = "15.0"
	}
	target := s.TargetArch + "-apple-ios" + minimum + suffix
	sdkPath, err := commandOutput(ctx, h.root, "xcrun", "--sdk", sdk, "--show-sdk-path")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	workspace := filepath.Join(h.root, ".wails", "package", "ios", s.TargetArch)
	if err := removeGenerated(workspace); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	executable := filepath.Join(workspace, s.Project.BinaryName)
	mainM := filepath.Join(h.root, s.Assets, "ios", "xcode", "main", "main.m")
	frameworks := []string{"Foundation", "UIKit", "WebKit", "Security", "CoreFoundation", "UniformTypeIdentifiers", "LocalAuthentication", "UserNotifications", "AVFoundation", "CoreLocation", "CoreMotion", "SystemConfiguration"}
	args := []string{"-sdk", sdk, "clang", "-target", target, "-isysroot", sdkPath}
	for _, framework := range frameworks {
		args = append(args, "-framework", framework)
	}
	args = append(args, "-lresolv", "-o", executable, mainM, "-Wl,-force_load,"+filepath.Join(h.root, s.Binary))
	linkOutput, err := runManifestCommand(ctx, h.root, nil, "xcrun", args...)
	if err != nil {
		return pipeline.RunResult{Detail: linkOutput}, err
	}
	appOutput := filepath.Join(h.root, s.Output)
	if ipa {
		appOutput = filepath.Join(workspace, s.Project.BinaryName+".app")
	}
	if err := removeGenerated(appOutput); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := os.MkdirAll(appOutput, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := copyManifestPath(executable, filepath.Join(appOutput, s.Project.BinaryName)); err != nil {
		return pipeline.RunResult{}, err
	}
	info := filepath.Join(h.root, s.Assets, "ios", "xcode", "main", "Info.plist")
	if err := copyManifestPath(info, filepath.Join(appOutput, "Info.plist")); err != nil {
		return pipeline.RunResult{}, err
	}
	assetInput := filepath.Join(h.root, s.Assets, "ios", "xcode", "main", "Assets.xcassets")
	assetTemp := filepath.Join(workspace, "compiled-assets")
	if _, err := os.Stat(assetInput); err == nil {
		if err := os.MkdirAll(assetTemp, 0o755); err != nil {
			return pipeline.RunResult{}, err
		}
		actool := []string{"actool", "--compile", assetTemp, "--app-icon", "AppIcon", "--platform", sdk, "--minimum-deployment-target", minimum, "--product-type", "com.apple.product-type.application", "--target-device", "iphone", "--target-device", "ipad", "--output-partial-info-plist", filepath.Join(appOutput, "assetcatalog_generated_info.plist"), assetInput}
		if output, err := runManifestCommand(ctx, h.root, nil, "xcrun", actool...); err != nil {
			return pipeline.RunResult{Detail: output}, err
		}
		if _, err := os.Stat(filepath.Join(assetTemp, "Assets.car")); err == nil {
			if err := copyManifestPath(filepath.Join(assetTemp, "Assets.car"), filepath.Join(appOutput, "Assets.car")); err != nil {
				return pipeline.RunResult{}, err
			}
		}
	}
	identity := "-"
	signing := h.config.Signing.IOS
	if signing.Enabled && signing.Identity != "" {
		identity = signing.Identity
	}
	signArgs := []string{"--force", "--sign", identity}
	if variant == "device" && signing.Entitlements != "" {
		signArgs = append(signArgs, "--entitlements", filepath.Join(h.root, signing.Entitlements))
	}
	signArgs = append(signArgs, appOutput)
	if output, err := runManifestCommand(ctx, h.root, nil, "codesign", signArgs...); err != nil {
		return pipeline.RunResult{Detail: output}, err
	}
	if !ipa {
		return pipeline.RunResult{Detail: linkOutput}, nil
	}
	payload := filepath.Join(workspace, "Payload")
	if err := os.MkdirAll(payload, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := copyManifestPath(appOutput, filepath.Join(payload, filepath.Base(appOutput))); err != nil {
		return pipeline.RunResult{}, err
	}
	outputPath := filepath.Join(h.root, s.Output)
	_ = os.Remove(outputPath)
	zipOutput, err := runManifestCommand(ctx, workspace, nil, "zip", "-qry", outputPath, "Payload")
	return pipeline.RunResult{Detail: strings.TrimSpace(linkOutput + "\n" + zipOutput)}, err
}

func (h *manifestHandler) sign(ctx context.Context, s pipeline.SignSpec) (pipeline.RunResult, error) {
	if !s.Config.Enabled {
		return pipeline.RunResult{}, fmt.Errorf("signing is not enabled for %s", s.TargetOS)
	}
	input := filepath.Join(h.root, s.Input)
	output := filepath.Join(h.root, s.Input+".signed")
	if s.TargetOS == "ios" {
		if s.Config.Identity == "" {
			return pipeline.RunResult{}, fmt.Errorf("iOS signing requires signing.ios.identity")
		}
		// iOS bundles must be signed while being assembled. packageIOS used the
		// configured identity; retain the normal signed-artifact output contract.
		return pipeline.RunResult{}, copyManifestPath(input, output)
	}
	if s.TargetOS == "android" {
		if s.Config.Certificate == "" || s.Config.Identity == "" || s.Config.Credential == "" {
			return pipeline.RunResult{}, fmt.Errorf("Android signing requires certificate (keystore), identity (key alias), and credential (password environment variable name)")
		}
		if os.Getenv(s.Config.Credential) == "" {
			return pipeline.RunResult{}, fmt.Errorf("Android signing password environment variable %s is empty", s.Config.Credential)
		}
		keystore := projectOrAbsolutePath(h.root, s.Config.Certificate)
		var detail string
		var err error
		if s.Format == "apk" {
			detail, err = runManifestCommand(ctx, h.root, nil, "apksigner", "sign", "--ks", keystore, "--ks-key-alias", s.Config.Identity, "--ks-pass", "env:"+s.Config.Credential, "--out", output, input)
		} else {
			detail, err = runManifestCommand(ctx, h.root, nil, "jarsigner", "-keystore", keystore, "-storepass:env", s.Config.Credential, "-signedjar", output, input, s.Config.Identity)
		}
		return pipeline.RunResult{Detail: detail}, err
	}
	err := Sign(&flags.Sign{Input: input, Certificate: s.Config.Certificate, Thumbprint: s.Config.Thumbprint, Timestamp: s.Config.TimestampServer, Identity: s.Config.Identity, Entitlements: s.Config.Entitlements, Notarize: s.Config.Notarize, KeychainProfile: s.Config.Credential, PGPKey: chooseString(s.TargetOS == "linux", s.Config.Certificate, ""), Role: chooseString(s.TargetOS == "linux", s.Config.Identity, "")})
	if err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{}, copyManifestPath(input, output)
}

func (h *manifestHandler) hook(ctx context.Context, s pipeline.HookSpec) (pipeline.RunResult, error) {
	script, err := pathInsideProject(h.root, s.Hook.Script)
	if err != nil {
		return pipeline.RunResult{}, fmt.Errorf("hook script: %w", err)
	}
	info, err := os.Stat(script)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	if info.IsDir() {
		return pipeline.RunResult{}, fmt.Errorf("hook script is a directory: %s", s.Hook.Script)
	}
	dir := h.root
	if s.Hook.Directory != "" {
		dir, err = pathInsideProject(h.root, s.Hook.Directory)
		if err != nil {
			return pipeline.RunResult{}, fmt.Errorf("hook directory: %w", err)
		}
	}
	env := []string{"WAILS_PROJECT_DIR=" + h.root, "WAILS_TARGET_OS=" + s.TargetOS, "WAILS_TARGET_ARCH=" + s.TargetArch, "WAILS_PROFILE=" + s.Profile, "WAILS_OUTPUT=" + s.Output, "WAILS_PIPELINE_VERSION=1"}
	output, err := runManifestCommand(ctx, dir, env, script)
	return pipeline.RunResult{Detail: output}, err
}

func runManifestCommand(ctx context.Context, dir string, extraEnv []string, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), extraEnv...)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Run()
	return strings.TrimSpace(output.String()), err
}

func runPackageManager(ctx context.Context, dir string, extraEnv []string, manager string, args ...string) (string, error) {
	if manager != "npm" {
		return runManifestCommand(ctx, dir, extraEnv, manager, args...)
	}
	npmPath, err := exec.LookPath("npm")
	if err != nil {
		return "", err
	}
	script, err := filepath.EvalSymlinks(npmPath)
	if err != nil {
		script = npmPath
	}
	nodePath, err := exec.LookPath("node")
	if err != nil {
		return "", err
	}
	return runManifestCommand(ctx, dir, extraEnv, nodePath, append([]string{script}, args...)...)
}
func toolIdentity(paths ...string) (string, error) {
	parts := make([]string, 0, len(paths))
	for _, path := range paths {
		resolved, err := filepath.EvalSymlinks(path)
		if err == nil {
			path = resolved
		}
		info, err := os.Stat(path)
		if err != nil {
			return "", err
		}
		parts = append(parts, fmt.Sprintf("%s:%d:%d:%o", path, info.Size(), info.ModTime().UnixNano(), info.Mode()))
	}
	sort.Strings(parts)
	return strings.Join(parts, "|"), nil
}
func appendUniqueStrings(values []string, extra ...string) []string {
	seen := map[string]bool{}
	for _, v := range values {
		seen[v] = true
	}
	for _, v := range extra {
		if !seen[v] {
			values = append(values, v)
			seen[v] = true
		}
	}
	return values
}
func chooseString(condition bool, value, fallback string) string {
	if condition {
		return value
	}
	return fallback
}
func projectOrAbsolutePath(root, value string) string {
	if filepath.IsAbs(value) {
		return value
	}
	return filepath.Join(root, filepath.FromSlash(value))
}
func pathInsideProject(root, value string) (string, error) {
	value = strings.ReplaceAll(value, `\`, "/")
	if filepath.IsAbs(value) {
		return "", fmt.Errorf("path must be project-relative: %s", value)
	}
	resolved := filepath.Clean(filepath.Join(root, filepath.FromSlash(value)))
	relative, err := filepath.Rel(root, resolved)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes the project: %s", value)
	}
	return resolved, nil
}

func writeGeneratedConfig(dir string, project manifest.Project, associations []manifest.Association, protocols []manifest.Protocol) (string, error) {
	type info struct {
		CompanyName       string `yaml:"companyName"`
		ProductName       string `yaml:"productName"`
		ProductIdentifier string `yaml:"productIdentifier"`
		Description       string `yaml:"description"`
		Copyright         string `yaml:"copyright"`
		Comments          string `yaml:"comments"`
		Version           string `yaml:"version"`
	}
	var files []FileAssociation
	for _, a := range associations {
		for _, ext := range a.Extensions {
			files = append(files, FileAssociation{Ext: ext, Name: a.Name, Description: a.Description, IconName: a.Icon, Role: a.Role, MimeType: a.MIMEType})
		}
	}
	var schemes []ProtocolConfig
	for _, p := range protocols {
		schemes = append(schemes, ProtocolConfig{Scheme: p.Scheme, Description: p.Description})
	}
	data, err := yaml.Marshal(struct {
		Info             info              `yaml:"info"`
		FileAssociations []FileAssociation `yaml:"fileAssociations,omitempty"`
		Protocols        []ProtocolConfig  `yaml:"protocols,omitempty"`
	}{Info: info{project.CompanyName, project.ProductName, project.Identifier, project.Description, project.Copyright, project.Comments, project.Version}, FileAssociations: files, Protocols: schemes})
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "config.yml")
	return path, os.WriteFile(path, data, 0o644)
}

func removeGenerated(path string) error {
	clean := filepath.Clean(path)
	if !strings.Contains(filepath.ToSlash(clean), "/.wails/") && !strings.HasSuffix(filepath.ToSlash(clean), ".app") {
		if _, err := os.Stat(clean); errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("refusing to replace non-generated path %s", path)
	}
	return os.RemoveAll(clean)
}
func copyManifestPath(source, destination string) error {
	info, err := os.Lstat(source)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(source)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return err
		}
		if err := os.Remove(destination); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return os.Symlink(target, destination)
	}
	if info.IsDir() {
		if err := os.MkdirAll(destination, info.Mode().Perm()); err != nil {
			return err
		}
		entries, err := os.ReadDir(source)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if err := copyManifestPath(filepath.Join(source, entry.Name()), filepath.Join(destination, entry.Name())); err != nil {
				return err
			}
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode().Perm())
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}
func ensureExpectedAppImage(dir, expected string) error {
	if _, err := os.Stat(expected); err == nil {
		return nil
	}
	matches, err := filepath.Glob(filepath.Join(dir, "*.AppImage"))
	if err != nil {
		return err
	}
	if len(matches) != 1 {
		return fmt.Errorf("AppImage packager produced %d candidates, expected %s", len(matches), expected)
	}
	return os.Rename(matches[0], expected)
}
func findAndMovePackage(dir, name, format, expected string) error {
	extension := "." + format
	if format == "archlinux" {
		extension = ".pkg.tar.zst"
	}
	matches, err := filepath.Glob(filepath.Join(dir, name+"*"+extension))
	if err != nil {
		return err
	}
	if len(matches) != 1 {
		return fmt.Errorf("packager produced %d %s candidates", len(matches), format)
	}
	if matches[0] == expected {
		return nil
	}
	if err := os.Rename(matches[0], expected); err == nil {
		return nil
	}
	// Windows does not replace an existing file with Rename. Build outputs are
	// disposable, so retry after removing only the exact expected destination.
	if err := os.Remove(expected); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return os.Rename(matches[0], expected)
}

var _ = runtime.GOOS
